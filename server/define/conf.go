package define

type SiteMap struct {
	Path string `yaml:"path"`
}
type Rpc struct {
	Host    string  `yaml:"host"`
	Addr    string  `yaml:"addr"`
	LogPath string  `yaml:"log_path"`
	SiteMap SiteMap `yaml:"site_map"`
}

type OrdxRpc struct {
	NsStatus           string `yaml:"ns_status"`
	NsRouting          string `yaml:"ns_routing"`
	InscriptionContent string `yaml:"inscription_content"`
}
