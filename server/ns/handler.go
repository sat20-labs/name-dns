package ns

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sat20-labs/name-ns/common"
)

func (s *Service) getContent(c *gin.Context) {
	hostParts := strings.Split(c.Request.Host, ".")
	if len(hostParts) < 3 {
		common.Log.Error("invalid host")
		return
	}

	if !strings.Contains(c.Request.Host, s.rpcConfig.Domain) {
		c.String(http.StatusNotFound, "invalid host")
		return
	}
	name := strings.TrimSuffix(c.Request.Host, s.rpcConfig.Domain)
	if name == "" {
		c.String(http.StatusNotFound, "invalid host")
		return
	}

	startTime := time.Now()
	nsRoutingResp, _, err := common.RpcRequest(s.ordxRpcConfig.NsRouting, name, "GET")
	elapsed := time.Since(startTime)
	common.Log.Info(fmt.Sprintf("call: %s, elapsed time: %s", s.ordxRpcConfig.NsRouting+name, elapsed))
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
	inscriptionContent, header, err := common.RpcRequest(s.ordxRpcConfig.InscriptionContent, nameRoutingResp.Data.InscriptionId, "GET")
	elapsed = time.Since(startTime)
	common.Log.Info(fmt.Sprintf("call: %s, elapsed time: %s", s.ordxRpcConfig.InscriptionContent+nameRoutingResp.Data.InscriptionId, elapsed))
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
}
