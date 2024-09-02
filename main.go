package main

import (
	"github.com/sat20-labs/name-dns/common"
	"github.com/sat20-labs/name-dns/main/flag"
	"github.com/sat20-labs/name-dns/main/g"
)

func init() {
	flag.ParseCmdParams()
	g.InitSigInt()

	err := g.InitDB()
	if err != nil {
		common.Log.Fatal(err)
	}
	common.Log.Info("init db")

	err = g.InitRpc()
	if err != nil {
		common.Log.Fatal(err)
	}
	common.Log.Info("init rpc")
}

func main() {
	common.Log.Info("starting...")
	defer func() {
		g.ReleaseDB()
		common.Log.Info("release db")
		common.Log.Info("shutdown...")
	}()
	err := g.RunRpc()
	if err != nil {
		common.Log.Fatal(err)
	}
	common.Log.Info("rpc started")
	select {}
}
