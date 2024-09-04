package define

type Rpc struct {
	Host    string `yaml:"host"`
	Addr    string `yaml:"addr"`
	LogPath string `yaml:"log_path"`
}

type OrdxRpc struct {
	NameList           string `yaml:"name_list"`
	NsRouting          string `yaml:"ns_routing"`
	InscriptionContent string `yaml:"inscription_content"`
}
