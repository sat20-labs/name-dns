package flag

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sat20-labs/name-ns/main/conf"
	serverCommon "github.com/sat20-labs/name-ns/server/define"
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
	if rpcService.Addr == "" {
		rpcService.Addr = "0.0.0.0:80"
	}
	if rpcService.LogPath == "" {
		rpcService.LogPath = "log"
	}
	ret.Rpc = rpcService
	return ret, nil
}

func NewDefaultYamlConf() (*conf.Conf, error) {
	ret := &conf.Conf{
		Log: conf.Log{
			Level: "error",
			Path:  "log",
		},
		Rpc: serverCommon.Rpc{
			Addr:    "0.0.0.0:80",
			Domain:  "ordx.space",
			LogPath: "log",
		},
		OrdxRpc: serverCommon.OrdxRpc{
			NsRouting:          "https://apiprd.ordx.space/testnet4/ns/name/",
			InscriptionContent: "https://apiprd.ordx.space/testnet4/ord/content/",
		},
	}

	return ret, nil
}

func SaveYamlConf(conf *conf.Conf, filePath string) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
