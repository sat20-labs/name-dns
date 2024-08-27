package ns

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sat20-labs/name-ns/common"
)

//go:embed templates/*
var templatesRes embed.FS

func (s *Service) countHtml(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	nameCounts, total, err := getNameCounts(s.DB, page, pageSize)
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	templateFile := "templates/name_count.html"
	tmpl, err := template.ParseFS(templatesRes, templateFile)
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	data := struct {
		NameCounts []NameCount
		Page       int
		PrevPage   int
		NextPage   int
		PageSize   int
		TotalPages int
	}{
		NameCounts: nameCounts,
		Page:       page,
		PrevPage:   page - 1,
		NextPage:   page + 1,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(c.Writer, data); err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
	}
}
