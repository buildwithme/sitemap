package domain

import (
	"encoding/xml"
	"log"
	"sort"

	"github.com/buildwithme/sitemap/pkg/models"
	"github.com/buildwithme/sitemap/pkg/utils"
)

func saveSitemapToFile(filepath string, urlDetails []models.UrlResult) error {
	sort.Slice(urlDetails, func(i, j int) bool {
		return urlDetails[i].URL < urlDetails[j].URL
	})

	sitemapXML := generateSitemapFile(urlDetails)

	data, err := xml.MarshalIndent(sitemapXML, "", "   ")
	if err != nil {
		log.Println(err)
		return err
	}

	var filename = filepath + ".xml"

	log.Printf("saveSitemapToFile: Save sitemap file to %s\n", filename)

	return utils.WriteToFile(filename, []byte(xml.Header), data)
}

func generateSitemapFile(urlDetails []models.UrlResult) *models.SitemapXML {
	var sitemapXML = &models.SitemapXML{Version: "http://www.sitemaps.org/schemas/sitemap/0.9", XSI: "http://www.w3.org/2001/XMLSchema-instance"}

	for i := range urlDetails {
		details := urlDetails[i]
		urlLink := details.URL
		sitemapXML.URL = append(sitemapXML.URL, models.UrlXML{Loc: &urlLink})
	}

	return sitemapXML
}
