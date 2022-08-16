package models

import (
	"encoding/xml"
)

type URLPriority string

const (
	URLPriorityAlways  URLPriority = "always"
	URLPriorityHourly  URLPriority = "hourly"
	URLPriorityDaily   URLPriority = "daily"
	URLPriorityWeekly  URLPriority = "weekly"
	URLPriorityMonthly URLPriority = "monthly"
	URLPriorityYearly  URLPriority = "yearly"
	URLPriorityNever   URLPriority = "never"
)

type SitemapXML struct {
	XMLName xml.Name `xml:"urlset"`
	Version string   `xml:"xmlns,attr"`
	XSI     string   `xml:"xmlns:xsi,attr"`
	Text    string   `xml:",chardata"`
	URL     []UrlXML `xml:"url"`
}
type UrlXML struct {
	Loc        *string      `xml:"loc,omitempty"`
	LastMod    *string      `xml:"lastmod,omitempty"`
	Changefreq *string      `xml:"changefreq,omitempty"`
	Priority   *URLPriority `xml:"priority,omitempty"`
}
