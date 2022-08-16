package main

import (
	"fmt"
	"os"

	"github.com/buildwithme/sitemap/pkg"
	"github.com/buildwithme/sitemap/pkg/constants"
	"github.com/buildwithme/sitemap/pkg/domain"
	"github.com/buildwithme/sitemap/pkg/utils"
)

func init() {
	pkg.GetParameterOptions()
}

func main() {
	url := retrieveURL()

	domain.GenerateSitemap(url)
}

func retrieveURL() string {
	arguments := os.Args[1:]

	var url string

	for i := 0; i < len(arguments); i++ {
		if arguments[i][0] == constants.Hyphen {
			continue
		}

		url = arguments[i]

		if valid, _ := utils.ValidURL(url); !valid {
			panic(fmt.Errorf("The URL provided doesn't match the regex %q", constants.UrlRegex))
		}

		break
	}

	if url == "" {
		panic("No URL provided")
	}

	return url
}
