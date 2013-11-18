package sitemap

import (
	"github.com/redneckbeard/gadget"
)

type SitemapController struct {
	*gadget.DefaultController
}

func (c *SitemapController) IdPattern() string { return `[\w-]+` }

func (c *SitemapController) Index(r *gadget.Request) (int, interface{}) {
	index := NewSitemapIndex()
	response := gadget.NewResponse(index)
	response.Headers.Set("Content-Type", "text/xml")
	return 200, response
}

func (c *SitemapController) Show(r *gadget.Request) (int, interface{}) {
	if urlset, found := urlSets[r.UrlParams["sitemap_id"]]; !found {
		return 404, ""
	} else {
		response := gadget.NewResponse(urlset)
		response.Headers.Set("Content-Type", "text/xml")
		return 200, response
	}
}
