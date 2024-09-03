package ns

import (
	"encoding/xml"

	serverOrdx "github.com/sat20-labs/name-dns/server/define"
)

type NameCount struct {
	Name  string
	Count uint64
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

type NameStatusResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Version string `json:"version"`
		Total   uint64 `json:"total"`
		Start   uint64 `json:"start"`
		Names   []struct {
			Id                 uint64 `json:"id"`
			Name               string `json:"name"`
			Sat                uint64 `json:"sat"`
			Address            string `json:"address"`
			InscriptionId      string `json:"inscription_id"`
			Utxo               string `json:"utxo"`
			Value              int64  `json:"value"`
			Height             int64  `json:"height"`
			Time               int64  `json:"time"`
			InscriptionAddress string `json:"inscription_address"`
		} `json:"names"`
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
