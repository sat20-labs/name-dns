package ns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sat20-labs/name-dns/common"
)

func (s *Service) content(c *gin.Context) {
	isMatch := false
	name := ""
	for _, domain := range s.RpcConfig.DomainList {
		if !strings.Contains(c.Request.Host, domain) {
			continue
		}
		if c.Request.Host == domain {
			continue
		}
		name = strings.TrimSuffix(c.Request.Host, "."+domain)
		if name == "" {
			continue
		}
		isMatch = true
		break
	}
	if !isMatch {
		msg := fmt.Sprintf("no match domain %s", c.Request.Host)
		c.String(http.StatusBadRequest, msg)
		return
	}

	startTime := time.Now()
	nsRoutingResp, _, err := common.ApiRequest(s.OrdxRpcConfig.NsRouting, name, "GET")
	elapsed := time.Since(startTime)
	common.Log.Info(fmt.Sprintf("call: %s, elapsed time: %s", s.OrdxRpcConfig.NsRouting+name, elapsed))
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var nameRoutingResp NameRoutingResp
	err = json.Unmarshal(nsRoutingResp, &nameRoutingResp)
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
	inscriptionContent, header, err := common.HtmlRequest(s.OrdxRpcConfig.InscriptionContent, nameRoutingResp.Data.Index)
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
	if err := s.incTotalNameCount(); err != nil {
		common.Log.Error(err)
	}
}
