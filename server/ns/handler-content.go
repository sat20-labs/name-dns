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

	contentType := header.Get("Content-Type")
	c.Data(http.StatusOK, contentType, inscriptionContent)

	if err := s.incrementNameCount(name); err != nil {
		common.Log.Error(err)
	}
	if err := s.incrementTotalNameCount(); err != nil {
		common.Log.Error(err)
	}
}
