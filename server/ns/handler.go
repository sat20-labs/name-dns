package ns

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sat20-labs/name-ns/common"
)

func (s *Service) getNameCount(c *gin.Context) {
	name := c.Param("name")
	count, err := getNameCount(s.DB, name)
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": name, "count": count})
}

func (s *Service) getContent(c *gin.Context) {
	isMatch := false
	name := ""
	for _, domain := range s.RpcConfig.DomainList {
		if !strings.Contains(c.Request.Host, domain) {
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
	nsRoutingResp, _, err := common.RpcRequest(s.OrdxRpcConfig.NsRouting, name, "GET")
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

	startTime = time.Now()
	inscriptionContent, header, err := common.RpcRequest(s.OrdxRpcConfig.InscriptionContent, nameRoutingResp.Data.Index, "GET")
	elapsed = time.Since(startTime)
	common.Log.Info(fmt.Sprintf("call: %s, elapsed time: %s", s.OrdxRpcConfig.InscriptionContent+nameRoutingResp.Data.Index, elapsed))
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for key, values := range header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	contentType := header.Get("Content-Type")
	c.Data(http.StatusOK, contentType, inscriptionContent)

	if err := incrementNameCount(s.DB, name); err != nil {
		common.Log.Error(err)
	}
}

func (s *Service) proxyReq(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is missing"})
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}
