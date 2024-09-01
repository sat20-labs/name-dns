package define

import (
	"gopkg.in/yaml.v2"
)

func ParseRpcConfig(data any) (*Rpc, error) {
	rpcServiceRaw, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}
	ret := &Rpc{}
	err = yaml.Unmarshal(rpcServiceRaw, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func ParseOrdxRpcConfig(data any) (*OrdxRpc, error) {
	rpcServiceRaw, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}
	ret := &OrdxRpc{}
	err = yaml.Unmarshal(rpcServiceRaw, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func ParseOrdinalsRpcConfig(data any) (*OrdinalsRpc, error) {
	rpcServiceRaw, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}
	ret := &OrdinalsRpc{}
	err = yaml.Unmarshal(rpcServiceRaw, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
