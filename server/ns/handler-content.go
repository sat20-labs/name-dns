package ns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sat20-labs/name-dns/common"
	"golang.org/x/net/idna"
)

func (s *Service) nameContent(c *gin.Context) {
	name := getSubdomain(c)
	if name == "" {
		data, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "index.html not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
		return
	}
	name, err := idna.ToUnicode(name)
	if err != nil {
		fmt.Println("Error decoding from Punycode:", err)
		return
	}

	startTime := time.Now()
	url := fmt.Sprintf(s.OrdxRpcConfig.NsRouting, name)
	resp, _, err := common.ApiRequest(url, "GET")
	elapsed := time.Since(startTime)
	common.Log.Info(fmt.Sprintf("call: %s, elapsed time: %s", s.OrdxRpcConfig.NsRouting+name, elapsed))
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var nameRoutingResp NameRoutingResp
	err = json.Unmarshal(resp, &nameRoutingResp)
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if nameRoutingResp.Code != 0 {
		c.String(http.StatusNotFound, nameRoutingResp.Msg)
		return
	}
	if nameRoutingResp.Data.Index == "" {
		c.String(http.StatusNotFound, "no find name routeing, Data.Index is empty")
		return
	}

	startTime = time.Now()
	url = fmt.Sprintf(s.OrdxRpcConfig.InscriptionContent, nameRoutingResp.Data.Index)
	inscriptionContent, header, err := common.HtmlRequest(url)
	elapsed = time.Since(startTime)

	common.Log.Info(fmt.Sprintf("call: %s, elapsed time: %s", s.OrdxRpcConfig.InscriptionContent+nameRoutingResp.Data.Index, elapsed))
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	contentType := header.Get("Content-Type")
	c.Writer.Header().Set("content-encoding", header.Get("content-encoding"))
	c.Writer.Header().Add("access-control-allow-origin:", "*")
	c.Writer.Header().Add("cache-control", "public, max-age=1209600, immutable")
	c.Writer.Header().Add("connection", "keep-alive")
	c.Writer.Header().Add("transfer-encoding", "chunked")
	c.Writer.Header().Add("Vary", "origin")
	c.Writer.Header().Add("Vary", "access-control-request-method")
	c.Writer.Header().Add("Vary", "access-control-request-headers")

	c.Data(http.StatusOK, contentType, inscriptionContent)
	if err := s.incNameCount(name); err != nil {
		common.Log.Error(err)
	}
	if err := s.incTotalNameAccessCount(); err != nil {
		common.Log.Error(err)
	}
}
