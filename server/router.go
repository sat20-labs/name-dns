package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rs/zerolog"
	"github.com/sat20-labs/gzip"
	serverCommon "github.com/sat20-labs/name-ns/server/define"
	"github.com/sat20-labs/name-ns/server/ns"
	"go.etcd.io/bbolt"
)

const (
	STRICT_TRANSPORT_SECURITY   = "strict-transport-security"
	CONTENT_SECURITY_POLICY     = "content-security-policy"
	CACHE_CONTROL               = "cache-control"
	VARY                        = "vary"
	ACCESS_CONTROL_ALLOW_ORIGIN = "access-control-allow-origin"
	TRANSFER_ENCODING           = "transfer-encoding"
	CONTENT_ENCODING            = "content-encoding"
)

const (
	CONTEXT_TYPE_TEXT = "text/html; charset=utf-8"
	CONTENT_TYPE_JSON = "application/json"
)

type Rpc struct {
	nsService *ns.Service
}

func NewRpc(
	rpcConfig *serverCommon.Rpc,
	ordxRpcConfig *serverCommon.OrdxRpc,
	db *bbolt.DB,
) *Rpc {
	return &Rpc{
		nsService: ns.New(rpcConfig, ordxRpcConfig, db),
	}
}

func (s *Rpc) Start() error {
	rpcUrl := s.nsService.RpcConfig.Addr
	rpcLogFile := s.nsService.RpcConfig.LogPath
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	var writers []io.Writer
	if rpcLogFile != "" {
		exePath, _ := os.Executable()
		executableName := filepath.Base(exePath)
		if strings.Contains(executableName, "debug") {
			executableName = "debug"
		}
		executableName += ".rpc"
		fileHook, err := rotatelogs.New(
			rpcLogFile+"/"+executableName+".%Y%m%d.log",
			rotatelogs.WithLinkName(rpcLogFile+"/"+executableName+".log"),
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(1*time.Hour),
		)
		if err != nil {
			return fmt.Errorf("failed to create RotateFile hook, error %s", err)
		}
		writers = append(writers, fileHook)
	}
	writers = append(writers, os.Stdout)
	gin.DefaultWriter = io.MultiWriter(writers...)
	r.Use(logger.SetLogger(
		logger.WithLogger(logger.Fn(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
			if c.Request.Header["Authorization"] == nil {
				return l
			}
			return l.With().
				Str("Authorization", c.Request.Header["Authorization"][0]).
				Logger()
		})),
	))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.OptionsResponseStatusCode = 200
	r.Use(cors.New(config))

	// common header
	// r.Use(func(c *gin.Context) {
	// 	c.Writer.Header().Set(VARY, "Origin")
	// 	c.Writer.Header().Add(VARY, "Access-Control-Request-Method")
	// 	c.Writer.Header().Add(VARY, "Access-Control-Request-Headers")

	// 	c.Writer.Header().Del(CONTENT_SECURITY_POLICY)
	// 	c.Writer.Header().Set(
	// 		CONTENT_SECURITY_POLICY,
	// 		"default-src 'self'",
	// 	)

	// 	c.Writer.Header().Set(
	// 		STRICT_TRANSPORT_SECURITY,
	// 		"max-age=31536000; includeSubDomains; preload",
	// 	)

	// 	c.Writer.Header().Set(
	// 		ACCESS_CONTROL_ALLOW_ORIGIN,
	// 		"*",
	// 	)
	// 	c.Next()
	// })

	// zip
	r.Use(
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPathsRegexs(
				[]string{
					// `.*\/btc\/.*`,
				},
			),
		),
	)

	// router
	err := s.nsService.Init(r)
	if err != nil {
		return err
	}

	parts := strings.Split(rpcUrl, ":")
	if len(parts) < 2 {
		rpcUrl += ":80"
	}

	go r.Run(rpcUrl)
	return nil
}
