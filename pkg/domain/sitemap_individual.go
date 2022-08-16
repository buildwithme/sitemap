package domain

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/buildwithme/sitemap/pkg"
	"github.com/buildwithme/sitemap/pkg/models"
)

type individualURLSitemap struct {
	Level          int
	Url            string
	uniqueURLs     map[string]models.UrlDetails
	wg             sync.WaitGroup
	chanCallURLs   chan<- models.UrlDetails
	mutexURLs      sync.Mutex
	chanFailedURLs chan<- models.FailedURL
	Options        *pkg.ParameterOptions
}

func newindividualURLSitemap(url string, level int, chanCallURLs chan<- models.UrlDetails, chanFailedURLs chan<- models.FailedURL,
	options *pkg.ParameterOptions) *individualURLSitemap {
	return &individualURLSitemap{
		uniqueURLs:     make(map[string]models.UrlDetails),
		Url:            url,
		Level:          level,
		chanCallURLs:   chanCallURLs,
		chanFailedURLs: chanFailedURLs,
		Options:        options,
	}
}

func (s *individualURLSitemap) sendFoundURLs(foundUrls chan<- map[string]models.UrlDetails) {
	s.mutexURLs.Lock()
	if len(s.uniqueURLs) > 0 {
		foundUrls <- s.uniqueURLs
	}
	s.mutexURLs.Unlock()
}

func (s *individualURLSitemap) request() {
	c := http.Client{
		Timeout: time.Duration(s.Options.Timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var err error
	var resp *http.Response
	var failedReason string

	for i := 0; i <= s.Options.MaxRetry; i++ {
		if i != 0 {
			time.Sleep(time.Second * time.Duration(s.Options.RetryDetay))
			log.Printf("Retry [%d] to request %q\n", i, s.Url)
		}

		resp, err = c.Get(s.Url)
		if err != nil {
			log.Printf("failed:request - for url %q - level [%d]: %s\n", s.Url, s.Level, err.Error())
			failedReason = err.Error()
			continue
		}

		defer resp.Body.Close()

		codeType := resp.StatusCode / 100
		if codeType == 2 {
			break
		} else if codeType == 4 {
			log.Println("failed:request - ", resp.Status, resp.StatusCode, s.Url)
			continue
		} else {
			log.Println("failed:request - ", resp.Status, resp.StatusCode, s.Url)
			break
		}
	}

	if resp == nil {
		s.chanFailedURLs <- models.FailedURL{
			URL:    s.Url,
			Reason: failedReason,
		}
		return
	} else if resp.StatusCode/100 != 2 {
		s.chanFailedURLs <- models.FailedURL{
			URL:    s.Url,
			Reason: resp.Status,
		}
		return
	}

	if err = s.processResponse(resp.Body, s.Level); err != nil {
		log.Printf("failed:request:processResponse - for url %q - level [%d]: %s\n", s.Url, s.Level, err.Error())
		return
	}
}

func (s *individualURLSitemap) processResponse(r io.Reader, level int) (err error) {
	chanLineText := make(chan string, 10)
	chanRawURLs := make(chan models.UrlDetails, 10)

	s.wg.Add(2)
	go s.parseHTML(chanLineText, chanRawURLs)
	go s.processURLs(chanRawURLs)

	scanner := bufio.NewScanner(r)

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		chanLineText <- line
	}

	close(chanLineText)

	if err = scanner.Err(); err != nil {
		log.Println("Error: ", err)
		return
	}

	s.wg.Wait()
	time.Sleep(time.Second)
	s.wg.Wait()

	return
}

func (s *individualURLSitemap) parseHTML(chanLineText <-chan string, chanRawURLs chan<- models.UrlDetails) {
	defer s.wg.Done()

	var openATag string
	var hasOpenATag bool
	var hasHref bool

	var href string
	var urlFound string
	var endFirstAPart bool
	var name string

	var endATag string

	var illegalCharacterInName bool

	var foundBody bool

	var base string
	for text := range chanLineText {
		if !foundBody {
			if strings.Contains(text, "<body") {
				foundBody = true
				continue
			}
			if strings.Contains(text, "<base") {
				indexStart := strings.Index(text, "href=\"") + 6
				indexEnd := strings.Index(text[indexStart:], "\"")
				base = text[indexStart : indexStart+indexEnd]
			}
		}
		for i := 0; i < len(text); i++ {
			currentCharacter := string(text[i])
			if hasOpenATag {
				if endFirstAPart {
					if endATag == "" {
						if currentCharacter == "<" {
							endATag += currentCharacter
							continue
						}
					}

					if endATag == "<" {
						if currentCharacter == "/" {
							endATag += currentCharacter
							continue
						} else {
							illegalCharacterInName = true
							endATag = ""
						}
					}

					if endATag == "</" {
						if currentCharacter == "a" {
							endATag += currentCharacter
							continue
						} else {
							illegalCharacterInName = true
							endATag = ""
						}
					}

					if endATag == "</a" {
						if currentCharacter == ">" {
							urlFound = strings.Trim(urlFound, " ")
							urlFound = strings.Trim(urlFound, "\t")
							if len(urlFound) > 0 {
								chanRawURLs <- models.UrlDetails{
									Base:  base,
									URL:   urlFound,
									Title: name,
									Level: s.Level,
								}
							}
							name = ""
							href = ""
							openATag = ""
							urlFound = ""
							hasOpenATag = false
							hasHref = false
							endFirstAPart = false
							illegalCharacterInName = false
						} else {
							illegalCharacterInName = true
						}
						endATag = ""
						continue
					}

					if !illegalCharacterInName {
						name += currentCharacter
					}
					continue
				}
				if hasHref {
					if currentCharacter == ">" {
						endFirstAPart = true
					}
					continue
				}
				if href == "" {
					if currentCharacter == "h" {
						href += currentCharacter
						continue
					}
				}
				if href == "h" {
					if currentCharacter == "r" {
						href += currentCharacter
						continue
					} else {
						href = ""
					}
				}
				if href == "hr" {
					if currentCharacter == "e" {
						href += currentCharacter
						continue
					} else {
						href = ""
					}
				}
				if href == "hre" {
					if currentCharacter == "f" {
						href += currentCharacter
						continue
					} else {
						href = ""
					}
				}
				if href == "href" {
					if currentCharacter == "=" {
						href += currentCharacter
						continue
					} else {
						href = ""
					}
				}
				if href == "href=" {
					if currentCharacter == "\"" {
						href += currentCharacter
						continue
					} else {
						href = ""
					}
				}
				if href == "href=\"" {
					if currentCharacter == "\"" {
						hasHref = true
						continue
					} else {
						urlFound += currentCharacter
					}
				} else {
					href = ""
				}
				continue
			}

			if openATag == "" {
				if currentCharacter == "<" {
					openATag += string(currentCharacter)
				} else {
					continue
				}
			} else if openATag == "<" {
				if currentCharacter == "a" {
					openATag += currentCharacter
				} else {
					openATag = ""
					continue
				}
			} else if openATag == "<a" {
				if currentCharacter != " " {
					openATag = ""
					hasOpenATag = false
				} else {
					hasOpenATag = true
				}
				continue
			}
		}
	}
	close(chanRawURLs)
}

func (s *individualURLSitemap) processURLs(chanRawURLs <-chan models.UrlDetails) {
	defer s.wg.Done()

	parsedURL, err := url.Parse(s.Url)
	if err != nil {
		panic(err)
	}

	searchedHost := parsedURL.Host

	baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	for urlDetailsValue := range chanRawURLs {
		urlDetails := urlDetailsValue

		urlFound := urlDetails.URL

		if strings.Contains(urlFound, searchedHost) {
			if len(urlFound) > 2 && urlFound[:2] == "//" {
				urlFound = fmt.Sprintf("%s:%s", parsedURL.Scheme, urlFound)
			} else if strings.Contains(urlFound, "mailto:") {
				continue
			}
		} else if urlFound[0] == '#' || len(urlFound) > 1 && urlFound[:2] == "/#" {
			continue
		} else if urlFound[0] == '/' {
			base := urlDetails.Base
			if base == "" || base == "/" {
				base = baseURL
			}
			urlFound = base + urlFound
		} else {
			continue
		}

		urlDetails.URL = urlFound

		if !strings.Contains(urlDetails.URL, baseURL) {
			continue
		}

		s.mutexURLs.Lock()
		if val, ok := s.uniqueURLs[urlDetails.URL]; ok {
			if val.Level > urlDetails.Level {
				s.uniqueURLs[urlDetails.URL] = urlDetails
			}
			s.mutexURLs.Unlock()
			continue
		}

		s.uniqueURLs[urlDetails.URL] = urlDetails
		s.mutexURLs.Unlock()

		s.chanCallURLs <- urlDetails
	}
}
