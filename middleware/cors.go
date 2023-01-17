package middleware

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"app-invite-service/config"
)

// CORSMiddleware is a middleware to add cors headers
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if isValidOrigin(origin, cfg.AllowedOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With,  X-authorizer-url")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func isValidOrigin(url string, allowedOriginsStr string) bool {
	var allowedOrigins []string
	if allowedOriginsStr == "" {
		allowedOrigins = []string{"*"}
	} else {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
	}
	if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
		return true
	}

	hasValidURL := false
	hostName, port := getHostParts(url)
	currentOrigin := hostName + ":" + port

	for _, origin := range allowedOrigins {
		replacedString := origin
		// if it has regex whitelisted domains
		if strings.Contains(origin, "*") {
			replacedString = strings.ReplaceAll(origin, ".", "\\.")
			replacedString = strings.ReplaceAll(replacedString, "*", ".*")

			if strings.HasPrefix(replacedString, ".*") {
				replacedString += "\\b"
			}

			if strings.HasSuffix(replacedString, ".*") {
				replacedString = "\\b" + replacedString
			}
		}

		if matched, _ := regexp.MatchString(replacedString, currentOrigin); matched {
			hasValidURL = true
			break
		}
	}

	return hasValidURL
}

// GetHostParts function returns hostname and port
func getHostParts(uri string) (string, string) {
	tempURI := uri
	if !strings.HasPrefix(tempURI, "http://") && !strings.HasPrefix(tempURI, "https://") {
		tempURI = "https://" + tempURI
	}

	u, err := url.Parse(tempURI)
	if err != nil {
		return "localhost", "8080"
	}

	host := u.Hostname()
	port := u.Port()

	return host, port
}
