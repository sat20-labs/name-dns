package flag

import (
	"flag"
	"os"

	"github.com/sat20-labs/name-dns/common"
	mainCommon "github.com/sat20-labs/name-dns/main/common"
	"github.com/sat20-labs/name-dns/main/g"
)

func ParseCmdParams() {
	init := flag.Bool("init", false, "generate config file in current dir")
	config := flag.String("config", "config.yaml", "env config file, default ./config.yaml")
	help := flag.Bool("help", false, "show help.")
	flag.Parse()

	if *help {
		common.Log.Info("name-ns server help:")
		common.Log.Info("Usage: 'name-ns -init'")
		common.Log.Info("Usage: 'name-ns -config ./config.yaml")
		common.Log.Info("Options:")
		common.Log.Info("  run service ->")
		common.Log.Info("    -init: init config in current dir")
		common.Log.Info("    -config: load config, default ./config.yaml")
		os.Exit(0)
	}

	if *init {
		err := generateConfigFile()
		if err != nil {
			common.Log.Fatal(err)
		}
		os.Exit(0)
	}

	err := InitConf(*config)
	if err != nil {
		common.Log.Fatal(err)
	}
	err = g.InitLog()
	if err != nil {
		common.Log.Fatal(err)
	}
}

func generateConfigFile() error {
	cfg, err := NewDefaultConf()
	if err != nil {
		return err
	}
	cfgPath, err := os.Getwd()
	if err != nil {
		return err
	}
	return SaveConf(cfg, cfgPath+"/config.yaml")
}

func InitConf(cfgPath string) (err error) {
	mainCommon.YamlCfg, err = LoadConf(cfgPath)
	if err != nil {
		return err
	}

	dbPath := mainCommon.YamlCfg.DB.Path
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		err := os.MkdirAll(dbPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return
}
