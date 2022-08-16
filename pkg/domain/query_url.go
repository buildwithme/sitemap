package domain

import (
	"sync"

	"github.com/buildwithme/sitemap/pkg"
	"github.com/buildwithme/sitemap/pkg/models"
)

func queryURL(wg *sync.WaitGroup, url string, level int, chanCallURLs chan<- models.UrlDetails, foundUrls chan<- map[string]models.UrlDetails, chanFailedURLs chan<- models.FailedURL, options *pkg.ParameterOptions) {
	defer wg.Done()

	if options.MaxDepth < level {
		return
	}

	sitemap := newindividualURLSitemap(url, level, chanCallURLs, chanFailedURLs, options)

	sitemap.request()

	sitemap.sendFoundURLs(foundUrls)
}
