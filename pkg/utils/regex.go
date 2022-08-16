package utils

import (
	"regexp"

	"github.com/buildwithme/sitemap/pkg/constants"
)

func ValidURL(url string) (bool, error) {
	matched, err := regexp.Match(constants.UrlRegex, []byte(url))
	if err != nil {
		return false, err
	}
	return matched, nil
}
