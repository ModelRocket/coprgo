package coprhd

import (
	"fmt"
)

const (
	VPoolQueryUriTpl = "block/vpools/%s.json"
	VPoolSearchUri   = "block/vpools/search.json?"
)

type (
	VPoolService struct {
		*Client

		id string
	}

	VPool struct {
		BaseObject `json:",inline"`
		Protocols  []string `json:"protocols"`
	}
)

func (this *Client) VPool() *VPoolService {
	return &VPoolService{
		Client: this,
	}
}

func (this *VPoolService) Id(id string) *VPoolService {
	this.id = id
	return this
}

func (this *VPoolService) Query() (*VPool, error) {
	path := fmt.Sprintf(VPoolQueryUriTpl, this.id)
	v := VPool{}

	err := this.Get(path, nil, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (this *VPoolService) Search(query string) (*VPool, error) {
	path := VPoolSearchUri + query

	res, err := this.Client.Search(path)
	if err != nil {
		return nil, err
	}

	this.id = res[0].Id

	return this.Query()
}
