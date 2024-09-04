package ns

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sat20-labs/name-dns/common"
)

func (s *Service) robots(c *gin.Context) {
	host := c.Request.Host
	sitemapURL := "https://" + host + "/sitemap_index.xml"
	robotsContent := "User-agent: *\n" +
		"Disallow: /private/\n\n" +
		"Disallow: /admin/\n\n" +
		"Disallow: /login/\n" +
		"Sitemap: " + sitemapURL
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, robotsContent)
}

func (s *Service) siteMapIndex(c *gin.Context) {
	totalNameCount, err := s.getTotalNameCount()
	if err != nil {
		common.Log.Error(err)
		return
	}
	nameListResp, err := s.ReqNameList(0, 1)
	if err != nil {
		common.Log.Error(err)
		return
	}

	if totalNameCount < nameListResp.Data.Total {
		siteMapIndex = &SiteMapIndex{
			XMLNS:           "http://www.sitemaps.org/schemas/sitemap/0.9",
			SiteMapItemList: []*SiteMapIndexItem{},
		}

		for index := uint64(0); index <= nameListResp.Data.Total/SITE_MAP_ITEM_COUNT; index++ {
			siteMapItem := &SiteMapIndexItem{
				Loc:     fmt.Sprintf("https://%s/%d.xml", s.RpcConfig.Host, index),
				LastMod: time.Now().Format("2006-01-02"),
			}
			siteMapIndex.SiteMapItemList = append(siteMapIndex.SiteMapItemList, siteMapItem)
		}
	}

	xmlData, _ := xml.MarshalIndent(siteMapIndex, "", "  ")
	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, string(xmlData))

	err = s.setTotalNameCount(nameListResp.Data.Total)
	if err != nil {
		common.Log.Error(err)
	}

	// err = os.WriteFile("./sitemap_index.xml", xmlData, 0644)
	// if err != nil {
	// 	common.Log.Errorf("save sitemap index error: %s", err.Error())
	// }
}

func (s *Service) siteMapItem(c *gin.Context) {
	siteMapItem := &SiteMapItem{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]*SiteMapItemURL, 0),
	}

	siteMapIndexStr := c.Param("index.xml")
	re := regexp.MustCompile(`(\d+)\.xml`)
	matches := re.FindStringSubmatch(siteMapIndexStr)

	if len(matches) < 2 {
		common.Log.Info("Invalid siteMapIndexStr: ", siteMapIndexStr)
		xmlData, _ := xml.MarshalIndent(siteMapItem, "", "  ")
		c.Header("Content-Type", "application/xml")
		c.String(http.StatusOK, string(xmlData))
		return
	}

	siteMapIndexStr = matches[1]
	siteMapIndex, err := strconv.ParseUint(siteMapIndexStr, 10, 64)
	if err != nil {
		common.Log.Errorf("Invalid number: %s\n", siteMapIndexStr)
		xmlData, _ := xml.MarshalIndent(siteMapItem, "", "  ")
		c.Header("Content-Type", "application/xml")
		c.String(http.StatusOK, string(xmlData))
		return
	}

	start := siteMapIndex * SITE_MAP_ITEM_COUNT
	nameListResp, err := s.ReqNameList(start, SITE_MAP_ITEM_COUNT)
	if err != nil {
		common.Log.Error(err)
		xmlData, _ := xml.MarshalIndent(siteMapItem, "", "  ")
		c.Header("Content-Type", "application/xml")
		c.String(http.StatusOK, string(xmlData))
		return
	}
	if nameListResp.Code != 0 {
		common.Log.Error(nameListResp.Msg)
		xmlData, _ := xml.MarshalIndent(siteMapItem, "", "  ")
		c.Header("Content-Type", "application/xml")
		c.String(http.StatusOK, string(xmlData))
		return
	}

	for _, nameItem := range nameListResp.Data.List {
		for key := range nameItem.KVs {
			if key == "ord_index" {
				lastmodTime := time.Unix(int64(nameItem.BlockTimestamp), 0)
				lastmod := lastmodTime.Format("2006-01-02")
				url := &SiteMapItemURL{
					Loc:        fmt.Sprintf("https://%s.%s", nameItem.Name, s.RpcConfig.Host),
					LastMod:    lastmod,
					ChangeFreq: "daily",
					Priority:   "0.8",
				}
				siteMapItem.URLs = append(siteMapItem.URLs, url)
				break
			}
		}
	}

	xmlData, _ := xml.MarshalIndent(siteMapItem, "", "  ")
	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, string(xmlData))
}
