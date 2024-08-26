package define

type Rpc struct {
	Addr    string `yaml:"addr"`
	Domain  string `yaml:"domain"`
	LogPath string `yaml:"log_path"`
}

type OrdxRpc struct {
	NsRouting          string `yaml:"ns_routing"`
	InscriptionContent string `yaml:"inscription_content"`
}
