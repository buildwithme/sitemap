package domain

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/buildwithme/sitemap/pkg/models"
	"github.com/buildwithme/sitemap/pkg/utils"
)

func saveFailedURLs(filepath string, failedURLs []models.FailedURL) error {
	sort.Slice(failedURLs, func(i, j int) bool {
		return failedURLs[i].URL < failedURLs[j].URL
	})

	data, err := json.MarshalIndent(failedURLs, "", "\t")
	if err != nil {
		log.Println(err)
		return err
	}

	var fileName = filepath + "_bad_urls.json"

	log.Printf("saveFailedURLs: Save URLs that have failed to respond with 2XX status to %s\n", fileName)

	return utils.WriteToFile(fileName, data)
}
