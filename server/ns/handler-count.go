package ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sat20-labs/name-dns/common"
)

func (s *Service) nameAccessCount(c *gin.Context) {
	resp := &NameCountListResp{
		BaseResp: BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: &NameCountListData{
			ListResp: ListResp{
				Total: 0,
			},
			List: make([]*NameCount, 0),
		},
	}

	req := RangeReq{Cursor: 0, Size: 100}
	if err := c.ShouldBindQuery(&req); err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	if req.Cursor < 0 {
		req.Cursor = 0
	}

	if req.Size < 1 || req.Size > 1000 {
		req.Size = 100
	}

	list, total, err := s.getAccessNameCountList(req.Cursor, req.Size)
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	resp.Data.Total = uint64(total)
	resp.Data.List = list
	c.JSON(http.StatusOK, resp)
}

func (s *Service) summary(c *gin.Context) {
	resp := &SummaryResp{
		BaseResp: BaseResp{
			Code: 0,
			Msg:  "ok",
		},
		Data: &SummaryData{
			TotalNameAccessCount: 0,
			IndexHtmlAccessCount: 0,
		},
	}
	totalNameAccessCount, err := s.getTotalNameAccessCount()
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	indexHtmlAccessCount, err := s.getIndexHtmlAccessCount()
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	resp.Data = &SummaryData{
		TotalNameAccessCount: totalNameAccessCount,
		IndexHtmlAccessCount: indexHtmlAccessCount,
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Service) addIndexHtmlAccessCount(c *gin.Context) {
	resp := BaseResp{
		Code: 0,
		Msg:  "ok",
	}

	err := s.incIndexHtmlAccessCount()
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
