package sitemap

import (
	"encoding/xml"
	"net/url"
	"path"
	"time"
)

const NAMESPACE = "http://www.sitemaps.org/schemas/sitemap/0.9"

var (
	Host    string
	urlSets = make(map[string]*UrlSet)
)

func Add(name string, urls ...*Url) {
	urlSet := NewUrlSet(urls...)
	urlSets[name] = urlSet
}

type SitemapIndex struct {
	XMLName   xml.Name   `xml:"sitemapindex"`
	Namespace string     `xml:"xmlns,attr"`
	Sitemaps  []*Sitemap `xml:"sitemap"`
}

func NewSitemapIndex() *SitemapIndex {
	index := &SitemapIndex{Namespace: NAMESPACE}
	u, _ := url.Parse(Host)
	for k, urls := range urlSets {
		u.Path = path.Join("sitemaps", k)
		index.Sitemaps = append(index.Sitemaps, &Sitemap{
			Loc:     u.String(),
			LastMod: urls.LastMod,
		})
	}
	return index
}

type Sitemap struct {
	Loc     string    `xml:"loc"`
	LastMod time.Time `xml:"lastmod,omitempty"`
}

type UrlSet struct {
	XMLName   xml.Name  `xml:"urlset"`
	Namespace string    `xml:"xmlns,attr"`
	LastMod   time.Time `xml:"-"`
	Urls      []*Url    `xml:"url"`
}

func NewUrlSet(urls ...*Url) *UrlSet {
	urlSet := &UrlSet{Namespace: NAMESPACE}
	var latest time.Time
	for _, u := range urls {
		loc, _ := url.Parse(u.Loc)
		if !loc.IsAbs() {
			loc, _ = url.Parse(Host)
			loc.Path = u.Loc
		}
		u.Loc = loc.String()
		if u.LastMod.After(latest) {
			latest = u.LastMod
		}
		urlSet.Urls = append(urlSet.Urls, u)
	}
	urlSet.LastMod = latest
	return urlSet
}

type Url struct {
	Loc        string    `xml:"loc"`
	LastMod    time.Time `xml:"lastmod,omitempty"`
	ChangeFreq string    `xml:"changefreq,omitempty"`
	Priority   float64   `xml:"priority,omitempty"`
}
