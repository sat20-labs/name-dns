package conf

type Conf struct {
	DB          DB  `yaml:"db"`
	Log         Log `yaml:"log"`
	Rpc         any `yaml:"rpc"`
	OrdxRpc     any `yaml:"ordx_rpc"`
	OrdinalsRpc any `yaml:"ordinals_rpc"`
}

type DB struct {
	Path string `yaml:"path"`
}

type Log struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}
