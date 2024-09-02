package ns

import (
	serverOrdx "github.com/sat20-labs/name-dns/server/define"
)

type NameCount struct {
	Name  string
	Count int
}

type NameRoutingResp struct {
	serverOrdx.BaseResp
	Data *NameRouting `json:"data"`
}

type NameRouting struct {
	Holder        string `json:"holder"`
	InscriptionId string `json:"inscription_id"`
	P             string `json:"p"`
	Op            string `json:"op"`
	Name          string `json:"name"`
	Handle        string `json:"ord_handle"`
	Index         string `json:"ord_index"`
}
