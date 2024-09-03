package ns

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sat20-labs/name-dns/common"
)

func (s *Service) sitemapFile(c *gin.Context) {
	index, err := strconv.Atoi(c.Param("index"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid page parameter")
		return
	}
	start := index * SITE_MAP_INDEX_ITEM_COUNT
	limit := SITE_MAP_INDEX_ITEM_COUNT
	url := fmt.Sprintf(s.OrdxRpcConfig.NsStatus, start, limit)
	resp, _, err := common.ApiRequest(url, "GET")
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var nameStatusResp NameStatusResp
	err = json.Unmarshal(resp, &nameStatusResp)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse nsstatus api response, error: %s", err.Error())
		return
	}

	urlSet := SiteMapItem{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]*SiteMapItemURL, 0),
	}

	for _, nameItem := range nameStatusResp.Data.Names {
		lastmodTime := time.Unix(int64(nameItem.Time), 0)
		lastmod := lastmodTime.Format("2006-01-02")
		url := &SiteMapItemURL{
			Loc:        fmt.Sprintf("https://%s.%s", nameItem.Name, s.RpcConfig.Host),
			LastMod:    lastmod,
			ChangeFreq: "daily",
			Priority:   "0.8",
		}
		urlSet.URLs = append(urlSet.URLs, url)
	}
	xmlData, _ := xml.MarshalIndent(urlSet, "", "  ")

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, string(xmlData))
}

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
