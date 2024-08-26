package main

import (
	"github.com/sat20-labs/name-ns/common"
	"github.com/sat20-labs/name-ns/main/flag"
	"github.com/sat20-labs/name-ns/main/g"
)

func init() {
	flag.ParseCmdParams()
	g.InitSigInt()
}

func main() {
	common.Log.Info("Starting...")
	defer func() {
		common.Log.Info("shut down")
	}()
	err := g.InitRpc()
	if err != nil {
		common.Log.Error(err)
		return
	}

	select {}
}
