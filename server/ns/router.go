package ns

import (
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/sat20-labs/name-dns/common"
	serverCommon "github.com/sat20-labs/name-dns/server/define"
	"go.etcd.io/bbolt"
)

//go:embed static/*
var staticFiles embed.FS

const (
	SITE_MAP_INDEX_ITEM_COUNT = 500
	GEN_SITE_MAP_TIME         = 10 * time.Minute
)

type Service struct {
	RpcConfig     *serverCommon.Rpc
	OrdxRpcConfig *serverCommon.OrdxRpc
	DB            *bbolt.DB
	siteMapIndex  *SiteMapIndex
}

func New(
	rpcConfig *serverCommon.Rpc,
	ordxRpcConfig *serverCommon.OrdxRpc,
	db *bbolt.DB) *Service {

	return &Service{
		RpcConfig:     rpcConfig,
		OrdxRpcConfig: ordxRpcConfig,
		DB:            db,
	}
}

func (s *Service) Init(r *gin.Engine) (err error) {
	err = common.InitBucket(s.DB, BUCKET_NAME_COUNT)
	if err != nil {
		return
	}
	err = common.InitBucket(s.DB, BUCKET_COMMON_SUMMARY)
	if err != nil {
		return
	}

	err = s.initSiteMap()
	if err != nil {
		return
	}

	// gen sitemap index
	go func() {
		ticker := time.NewTicker(GEN_SITE_MAP_TIME)
		defer ticker.Stop()

		s.genSiteMapIndex()
		for range ticker.C {
			s.genSiteMapIndex()
		}
	}()

	s.initRouter(r)
	return
}

func (s *Service) initSiteMap() (err error) {
	err = os.MkdirAll(s.RpcConfig.SiteMap.Path, os.ModePerm)
	if err != nil {
		return
	}
	sitemapIndexPath := fmt.Sprintf("%s/sitemap_index.xml", s.RpcConfig.SiteMap.Path)
	if _, err := os.Stat(sitemapIndexPath); os.IsNotExist(err) {
		standardSitemapIndex := `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
</sitemapindex>`
		err = os.WriteFile(sitemapIndexPath, []byte(standardSitemapIndex), 0644)
		if err != nil {
			return fmt.Errorf("failed to create sitemap_index.xml: %v", err)
		}
	}
	return nil
}

func (s *Service) genSiteMapIndex() {
	totalNameCount, err := s.getTotalNameCount()
	if err != nil {
		common.Log.Error(err)
		return
	}
	if s.siteMapIndex == nil {
		s.siteMapIndex = &SiteMapIndex{}
		sitemapIndexPath := fmt.Sprintf("%s/sitemap_index.xml", s.RpcConfig.SiteMap.Path)
		xmlFile, err := os.Open(sitemapIndexPath)
		if err != nil {
			common.Log.Error(err)
			return
		}
		defer xmlFile.Close()

		byteValue, err := io.ReadAll(xmlFile)
		if err != nil {
			common.Log.Error(err)
			return
		}

		err = xml.Unmarshal(byteValue, s.siteMapIndex)
		if err != nil {
			common.Log.Error(err)
			return
		}
	}

	page := uint64(0)
	limit := 1
	url := fmt.Sprintf(s.OrdxRpcConfig.NsStatus, page, limit)
	resp, _, err := common.ApiRequest(url, "GET")
	if err != nil {
		common.Log.Error(err)
		return
	}
	var nameStatusResp NameStatusResp
	err = json.Unmarshal(resp, &nameStatusResp)
	if err != nil {
		common.Log.Error(err)
		return
	}

	if totalNameCount >= nameStatusResp.Data.Total {
		return
	}

	totalNameDomainCount, err := s.getTotalNameDomainCount()
	if err != nil {
		common.Log.Error(err)
		return
	}
	total := nameStatusResp.Data.Total / SITE_MAP_INDEX_ITEM_COUNT
	index := uint64(8)
	if len(s.siteMapIndex.SiteMapItemList) > 0 {
		re := regexp.MustCompile(`https://[^/]+/sitemap/(\d+)\.xml`)
		lastIndex := uint64(len(s.siteMapIndex.SiteMapItemList) - 1)
		loc := s.siteMapIndex.SiteMapItemList[lastIndex].Loc
		matches := re.FindStringSubmatch(loc)
		if len(matches) < 2 {
			common.Log.Errorf("failed to find string submatch: %s", loc)
			return
		}
		index, err = strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			common.Log.Error(err)
			return
		}

		sitemapItemPath := fmt.Sprintf("%s/sitemap_%d.xml", s.RpcConfig.SiteMap.Path, index)
		xmlFile, err := os.Open(sitemapItemPath)
		if err != nil {
			common.Log.Error(err)
			return
		}
		defer xmlFile.Close()

		byteValue, err := io.ReadAll(xmlFile)
		if err != nil {
			common.Log.Error(err)
			return
		}

		var siteMapItem SiteMapItem
		err = xml.Unmarshal(byteValue, &siteMapItem)
		if err != nil {
			common.Log.Error(err)
			return
		}
		totalNameDomainCount -= uint64(len(siteMapItem.URLs))
		s.siteMapIndex.SiteMapItemList[lastIndex].LastMod = time.Now().Format("2006-01-02")
		s.siteMapIndex.SiteMapItemList[lastIndex].Loc = fmt.Sprintf("https://%s/sitemap/%d.xml", s.RpcConfig.Host, index)
		genCount := s.genSiteMapFile(index)
		totalNameDomainCount += genCount
		index++
	}

	for ; index <= total; index++ {
		siteMapItem := &SiteMapIndexItem{
			Loc:     fmt.Sprintf("https://%s/sitemap/%d.xml", s.RpcConfig.Host, index),
			LastMod: time.Now().Format("2006-01-02"),
		}
		genCount := s.genSiteMapFile(index)
		if genCount == 0 {
			continue
		}
		totalNameDomainCount += genCount
		s.siteMapIndex.SiteMapItemList = append(s.siteMapIndex.SiteMapItemList, siteMapItem)
	}

	xmlData, _ := xml.MarshalIndent(s.siteMapIndex, "", "  ")
	sitemapIndexPath := fmt.Sprintf("%s/sitemap_index.xml", s.RpcConfig.SiteMap.Path)
	err = os.WriteFile(sitemapIndexPath, xmlData, 0644)
	if err != nil {
		common.Log.Error(err)
	}
	err = s.setTotalNameDomainCount(totalNameDomainCount)
	if err != nil {
		common.Log.Error(err)
		return
	}
	err = s.setTotalNameCount(nameStatusResp.Data.Total)
	if err != nil {
		common.Log.Error(err)
		return
	}
}

func (s *Service) genSiteMapFile(siteMapIndex uint64) uint64 {
	start := siteMapIndex * SITE_MAP_INDEX_ITEM_COUNT
	limit := SITE_MAP_INDEX_ITEM_COUNT
	url := fmt.Sprintf(s.OrdxRpcConfig.NsStatus, start, limit)
	resp, _, err := common.ApiRequest(url, "GET")
	if err != nil {
		common.Log.Error(err)
		return 0
	}

	var nameStatusResp NameStatusResp
	err = json.Unmarshal(resp, &nameStatusResp)
	if err != nil {
		common.Log.Error(err)
		return 0
	}

	urlSet := SiteMapItem{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]*SiteMapItemURL, 0),
	}

	index := -1
	for _, nameItem := range nameStatusResp.Data.Names {
		index++
		resp, _, err := common.ApiRequest(fmt.Sprintf(s.OrdxRpcConfig.NsRouting, nameItem.Name), "GET")
		if err != nil {
			common.Log.Error(err)
			continue
		}

		var nameRoutingResp NameRoutingResp
		err = json.Unmarshal(resp, &nameRoutingResp)
		if err != nil {
			common.Log.Error(err)
			continue
		}
		if nameRoutingResp.Code != 0 {
			common.Log.Warn(nameRoutingResp.Msg)
			continue
		}

		if nameRoutingResp.Data.Index == "" {
			common.Log.Debugf("no find name routeing, Data.Index is empty, sitemap index:%d, index:%d, name: %s", siteMapIndex, index, nameRoutingResp.Data.Name)
			continue
		}

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
	if len(urlSet.URLs) == 0 {
		return 0
	}
	xmlData, _ := xml.MarshalIndent(urlSet, "", "  ")
	sitemapIndexPath := fmt.Sprintf("%s/%d.xml", s.RpcConfig.SiteMap.Path, siteMapIndex)
	err = os.WriteFile(sitemapIndexPath, xmlData, 0644)
	if err != nil {
		common.Log.Error(err)
	}

	return uint64(len(urlSet.URLs))
}

func (s *Service) initRouter(r *gin.Engine) {
	staticServer := http.FS(staticFiles)
	r.StaticFS("/static", staticServer)
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("static/favicon.ico", staticServer)
	})

	r.GET("/robots.txt", s.robots)
	r.Static("/sitemap", s.RpcConfig.SiteMap.Path)
	r.GET("/", s.content)
	r.GET("/summary/name-count", s.countHtml)
}
