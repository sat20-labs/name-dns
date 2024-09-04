package ns

import (
	"encoding/json"
	"fmt"

	"github.com/sat20-labs/name-dns/common"
)

func (s *Service) ReqNameList(start, limit uint64) (*NameListResp, error) {
	url := fmt.Sprintf(s.OrdxRpcConfig.NameList, start, limit)
	resp, _, err := common.ApiRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	var nameListResp NameListResp
	err = json.Unmarshal(resp, &nameListResp)
	if err != nil {
		return nil, err
	}
	return &nameListResp, nil
}
