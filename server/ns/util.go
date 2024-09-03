package ns

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

func getSubdomain(c *gin.Context) string {
	pattern := `^((?:[a-zA-Z0-9-_]+\.)+)[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(c.Request.Host)
	if len(matches) > 1 {
		subdomain := matches[1]
		if subdomain[len(subdomain)-1] == '.' {
			subdomain = subdomain[:len(subdomain)-1]
		}
		return subdomain
	} else {
		return ""
	}
}
