package g

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sat20-labs/name-dns/common"
	mainCommon "github.com/sat20-labs/name-dns/main/common"
	"github.com/sat20-labs/name-dns/server"
	serverCommon "github.com/sat20-labs/name-dns/server/define"
	"github.com/sirupsen/logrus"
)

var (
	rpc    *server.Rpc
	SigInt chan os.Signal
)

func InitLog() error {
	var writers []io.Writer
	logPath := ""
	var lvl logrus.Level
	if mainCommon.YamlCfg != nil {
		logPath = mainCommon.YamlCfg.Log.Path
		var err error
		lvl, err = logrus.ParseLevel(mainCommon.YamlCfg.Log.Level)
		if err != nil {
			return fmt.Errorf("failed to parse log level: %s", err)
		}
	} else {
		return fmt.Errorf("cfg is not set")
	}
	if logPath != "" {
		exePath, _ := os.Executable()
		executableName := filepath.Base(exePath)
		if strings.Contains(executableName, "debug") {
			executableName = "debug"
		}
		fileHook, err := rotatelogs.New(
			logPath+"/"+executableName+".%Y%m%d.log",
			rotatelogs.WithLinkName(logPath+"/"+executableName+".log"),
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(1*time.Hour),
		)
		if err != nil {
			return fmt.Errorf("failed to create RotateFile hook, error: %s", err)
		}
		writers = append(writers, fileHook)
	}
	writers = append(writers, os.Stdout)
	common.Log.SetOutput(io.MultiWriter(writers...))
	common.Log.SetLevel(lvl)
	return nil
}

func InitSigInt() {
	SigInt = make(chan os.Signal, 100)
	signal.Notify(SigInt, os.Interrupt)
	go func() {
		for {
			<-SigInt
			common.Log.Infof("Received SIGINT (CTRL+C) and exit")
			os.Exit(0)
		}
	}()
}

func InitRpc() error {
	rpcConfig, err := serverCommon.ParseRpcConfig(mainCommon.YamlCfg.Rpc)
	if err != nil {
		return err
	}
	ordxRpcConfig, err := serverCommon.ParseOrdxRpcConfig(mainCommon.YamlCfg.OrdxRpc)
	if err != nil {
		return err
	}
	rpc = server.NewRpc(rpcConfig, ordxRpcConfig, store)
	return nil
}

func RunRpc() error {
	err := rpc.Start()
	if err != nil {
		return err
	}
	return nil
}
