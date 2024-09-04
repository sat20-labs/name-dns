package flag

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sat20-labs/name-dns/main/conf"
	serverCommon "github.com/sat20-labs/name-dns/server/define"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func LoadConf(cfgPath string) (*conf.Conf, error) {
	confFile, err := os.Open(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open cfg: %s, error: %s", cfgPath, err)
	}
	defer confFile.Close()

	ret := &conf.Conf{}
	decoder := yaml.NewDecoder(confFile)
	err = decoder.Decode(ret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cfg: %s, error: %s", cfgPath, err)
	}

	if ret.DB.Path == "" {
		ret.DB.Path = "./db/"
	}
	ret.DB.Path = filepath.FromSlash(ret.DB.Path)
	if ret.DB.Path[len(ret.DB.Path)-1] != filepath.Separator {
		ret.DB.Path += string(filepath.Separator)
	}

	_, err = logrus.ParseLevel(ret.Log.Level)
	if err != nil {
		ret.Log.Level = "info"
	}
	if ret.Log.Path == "" {
		ret.Log.Path = "log"
	}
	ret.Log.Path = filepath.FromSlash(ret.Log.Path)
	if ret.Log.Path[len(ret.Log.Path)-1] != filepath.Separator {
		ret.Log.Path += string(filepath.Separator)
	}

	rpcService, err := serverCommon.ParseRpcConfig(ret.Rpc)
	if err != nil {
		return nil, err
	}
	if rpcService.Host == "" {
		rpcService.Host = "dkvs.xyz"
	}
	if rpcService.Addr == "" {
		rpcService.Addr = "0.0.0.0:9006"
	}
	if rpcService.LogPath == "" {
		rpcService.LogPath = "log/testnet4"
	}
	ret.Rpc = rpcService
	return ret, nil
}

func NewDefaultConf() (*conf.Conf, error) {
	ret := &conf.Conf{
		DB: conf.DB{
			Path: "db",
		},
		Log: conf.Log{
			Level: "debug",
			Path:  "log/testnet4",
		},
		Rpc: serverCommon.Rpc{
			Host:    "dkvs.xyz",
			Addr:    "0.0.0.0:9006",
			LogPath: "log/testnet4",
		},
		OrdxRpc: serverCommon.OrdxRpc{
			NameList:           "https://apiprd.sat20.org/testnet4/ns/status/",
			NsRouting:          "https://apiprd.sat20.org/testnet4/ns/name/",
			InscriptionContent: "https://apiprd.sat20.org/testnet4/ord/content/",
		},
	}
	return ret, nil
}

func SaveConf(conf *conf.Conf, filePath string) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
