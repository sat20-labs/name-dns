package ns

import (
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
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
		ticker := time.NewTicker(1 * time.Minute)
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

		pageCount := totalNameCount / SITE_MAP_INDEX_ITEM_COUNT
		listLen := uint64(len(s.siteMapIndex.SiteMapItemList))
		if pageCount != listLen {
			s.siteMapIndex = &SiteMapIndex{
				XMLNS:           "http://www.sitemaps.org/schemas/sitemap/0.9",
				SiteMapItemList: make([]*SiteMapIndexItem, 0),
			}
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

	total := nameStatusResp.Data.Total / SITE_MAP_INDEX_ITEM_COUNT
	index := totalNameCount / SITE_MAP_INDEX_ITEM_COUNT
	div1 := totalNameCount % SITE_MAP_INDEX_ITEM_COUNT
	common.Log.Infof("total: %d, index: %d, div1: %d", total, index, div1)
	if len(s.siteMapIndex.SiteMapItemList) != 0 && totalNameCount%SITE_MAP_INDEX_ITEM_COUNT != 0 {
		lastIndex := uint64(len(s.siteMapIndex.SiteMapItemList) - 1)
		if lastIndex != index {
			common.Log.Error("lastIndex != index")
			return
		}
		s.siteMapIndex.SiteMapItemList[index].LastMod = time.Now().Format("2006-01-02")
		s.siteMapIndex.SiteMapItemList[index].Loc = fmt.Sprintf("https://%s/sitemap/%d.xml", s.RpcConfig.Host, totalNameCount/SITE_MAP_INDEX_ITEM_COUNT)
		index++
	}

	for ; index <= total; index++ {
		siteMapItem := &SiteMapIndexItem{
			Loc:     fmt.Sprintf("https://%s/sitemaploc/%d", s.RpcConfig.Host, index),
			LastMod: time.Now().Format("2006-01-02"),
		}
		s.siteMapIndex.SiteMapItemList = append(s.siteMapIndex.SiteMapItemList, siteMapItem)
	}

	xmlData, _ := xml.MarshalIndent(s.siteMapIndex, "", "  ")
	sitemapIndexPath := fmt.Sprintf("%s/sitemap_index.xml", s.RpcConfig.SiteMap.Path)
	err = os.WriteFile(sitemapIndexPath, xmlData, 0644)
	if err != nil {
		common.Log.Error(err)
	}
	s.setTotalNameCount(nameStatusResp.Data.Total)
}

func (s *Service) initRouter(r *gin.Engine) {
	staticServer := http.FS(staticFiles)
	r.StaticFS("/static", staticServer)
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("static/favicon.ico", staticServer)
	})

	r.GET("/robots.txt", s.robots)
	r.Static("/sitemap", s.RpcConfig.SiteMap.Path)
	r.GET("/sitemaploc/:index", s.sitemapFile)

	r.GET("/", s.content)
	r.GET("/summary/name-count", s.countHtml)
}
