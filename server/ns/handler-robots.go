package ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) robots(c *gin.Context) {
	host := c.Request.Host
	sitemapURL := "https://" + host + "/sitemap/sitemap_index.xml"
	robotsContent := "User-agent: *\n" +
		"Disallow: /private/\n\n" +
		"Disallow: /admin/\n\n" +
		"Disallow: /login/\n" +
		"Sitemap: " + sitemapURL

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, robotsContent)
}
