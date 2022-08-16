package domain

import (
	"log"
	"sync"
	"time"

	"github.com/buildwithme/sitemap/pkg"
	"github.com/buildwithme/sitemap/pkg/models"
	"github.com/buildwithme/sitemap/pkg/utils"
)

func GenerateSitemap(url string) {
	options := pkg.GetParameterOptions()

	urlDetails, failedURLs := StartSitemapGeneration(url, options)

	utils.EnsureDirectories(options.OutputFile)

	saveSitemapToFile(options.OutputFile, urlDetails)
	saveFailedURLs(options.OutputFile, failedURLs)
}

func StartSitemapGeneration(url string, options *pkg.ParameterOptions) ([]models.UrlResult, []models.FailedURL) {
	chanCallURLs := make(chan models.UrlDetails, 11)
	foundUrls := make(chan map[string]models.UrlDetails, 10)
	chanFailedURLs := make(chan models.FailedURL, 10)

	var wg sync.WaitGroup
	uniqueURLs := make(map[string]struct{})

	callExecuteSitemap := func(details models.UrlDetails, canAddDelta bool) {
		if _, ok := uniqueURLs[details.URL]; ok {
			return
		} else {
			uniqueURLs[details.URL] = struct{}{}
		}
		log.Printf("Request - [%d] - %q\n", details.Level, details.URL)
		if canAddDelta {
			wg.Add(1)
		}

		go queryURL(&wg, details.URL, details.Level+1, chanCallURLs, foundUrls, chanFailedURLs, options)
	}

	finishReceive := make(chan struct{})

	wg.Add(1)
	callExecuteSitemap(models.UrlDetails{URL: url}, false)

	go func() {
		for detailsValue := range chanCallURLs {
			callExecuteSitemap(detailsValue, true)
		}
	}()

	var uniqueUrlDetails = make(map[string]struct{})
	var urlDetails []models.UrlDetails
	go func() {
		for urlsDetails := range foundUrls {
			for key, value := range urlsDetails {
				if _, ok := uniqueUrlDetails[key]; ok {
					continue
				}
				uniqueUrlDetails[key] = struct{}{}
				urlDetails = append(urlDetails, value)
			}
		}
		finishReceive <- struct{}{}
	}()

	var failedURLs []models.FailedURL
	go func() {
		for failedURL := range chanFailedURLs {
			failedURLs = append(failedURLs, failedURL)
		}
		finishReceive <- struct{}{}
	}()

	wg.Wait()
	time.Sleep(time.Second)
	wg.Wait()
	close(chanCallURLs)
	close(foundUrls)
	close(chanFailedURLs)

	<-finishReceive
	<-finishReceive

	urlDetails = excludeFailedURLs(urlDetails, failedURLs)

	urlResults := convertToUrlResults(urlDetails)

	return urlResults, failedURLs
}

func convertToUrlResults(urlDetails []models.UrlDetails) (result []models.UrlResult) {
	for _, details := range urlDetails {
		result = append(result, models.UrlResult{
			URL:   details.URL,
			Title: details.Title,
		})
	}
	return
}

func excludeFailedURLs(urlDetails []models.UrlDetails, failedURLs []models.FailedURL) []models.UrlDetails {
	mapFailed := make(map[string]struct{})
	for i := range failedURLs {
		mapFailed[failedURLs[i].URL] = struct{}{}
	}
	var result []models.UrlDetails
	for i := range urlDetails {
		details := urlDetails[i]

		if _, ok := mapFailed[details.URL]; ok {
			continue
		}

		result = append(result, details)
	}

	return result
}
