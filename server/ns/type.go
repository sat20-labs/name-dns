package ns

import (
	"encoding/xml"

	serverOrdx "github.com/sat20-labs/name-dns/server/define"
)

type BaseResp struct {
	Code int    `json:"code" example:"0"`
	Msg  string `json:"msg" example:"ok"`
}

type RangeReq struct {
	Cursor int `form:"cursor" binding:"omitempty"`
	Size   int `form:"size" binding:"omitempty"`
}

type ListResp struct {
	Total uint64 `json:"total" example:"9992"`
}

type NameCount struct {
	Name       string `json:"name"`
	ClickCount uint64 `json:"clickCount"`
}

type NameRoutingResp struct {
	serverOrdx.BaseResp
	Data *NameRouting `json:"data"`
}

type NameRouting struct {
	Holder        string `json:"holder"`
	InscriptionId string `json:"inscription_id"`
	P             string `json:"p"`
	Op            string `json:"op"`
	Name          string `json:"name"`
	Handle        string `json:"ord_handle"`
	Index         string `json:"ord_index"`
}

type KeyValueInDB struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	InscriptionId string `json:"inscriptionId"`
}

type NameInfo struct {
	InscriptionNumber  int64  `json:"inscriptionNumber"`
	Name               string `json:"name"`
	Sat                int64  `json:"sat"`
	Address            string `json:"address"`
	InscriptionId      string `json:"inscriptionId"`
	Utxo               string `json:"utxo"`
	Value              int64  `json:"value"`
	BlockHeight        int64  `json:"height"`
	BlockTimestamp     int64  `json:"timestamp"`
	InscriptionAddress string `json:"inscriptionAddress"`
	Preview            string `json:"preview"`
	KVs                map[string]*KeyValueInDB
}

type NameListResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total uint64     `json:"total"`
		Start uint64     `json:"start"`
		List  []NameInfo `json:"list"`
	} `json:"data"`
}

type SiteMapIndexItem struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type SiteMapIndex struct {
	XMLName         xml.Name            `xml:"sitemapindex"`
	XMLNS           string              `xml:"xmlns,attr"`
	SiteMapItemList []*SiteMapIndexItem `xml:"sitemap"`
}

type SiteMapItemURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

type SiteMapItem struct {
	XMLName xml.Name          `xml:"urlset"`
	XMLNS   string            `xml:"xmlns,attr"`
	URLs    []*SiteMapItemURL `xml:"url"`
}

type NameCountListData struct {
	ListResp
	List []*NameCount `json:"list"`
}

type NameCountListResp struct {
	BaseResp
	Data *NameCountListData `json:"data"`
}

type SummaryData struct {
	TotalNameClickCount uint64 `json:"totalNameClickCount"`
}

type SummaryResp struct {
	BaseResp
	Data *SummaryData `json:"data"`
}
