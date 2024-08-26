package conf

type Conf struct {
	Log     Log `yaml:"log"`
	Rpc     any `yaml:"rpc"`
	OrdxRpc any `yaml:"ordx_rpc"`
}

type Log struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}
