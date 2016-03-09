package coprhd

import (
	"fmt"
)

const (
	InitiatorSearchUri   = "compute/initiators/search.json?"
	InitiatorQueryUriTpl = "compute/initiators/%s.json"

	InitiatorTypeISCSI InitiatorType = "iSCSI"
	InitiatorTypeFC    InitiatorType = "FC"
)

type (
	InitiatorService struct {
		*Client
		id       string
		protocol InitiatorType
		node     string
		port     string
	}

	Initiator struct {
		BaseObject `json:",inline"`
		Host       Resource      `json:"host"`
		Protocol   InitiatorType `json:"protocol"`
		Status     string        `json:"registration_status"`
		Hostname   string        `json:"hostname"`
		Node       string        `json:"initiator_node"`
		Port       string        `json:"initiator_port"`
	}

	InitiatorType string
)

func (this *Client) Initiator() *InitiatorService {
	return &InitiatorService{
		Client: this,
	}
}

func (this *InitiatorService) Id(id string) *InitiatorService {
	this.id = id
	return this
}

func (this *InitiatorService) Port(port string) *InitiatorService {
	this.port = port
	return this
}

func (this *InitiatorService) Query() (*Initiator, error) {
	if !isStorageOsUrn(this.id) {
		return this.Search("initiator_port=" + this.port)
	}

	path := fmt.Sprintf(InitiatorQueryUriTpl, this.id)
	itr := Initiator{}

	err := this.Get(path, nil, &itr)
	if err != nil {
		return nil, err
	}

	return &itr, nil
}

func (this *InitiatorService) Search(query string) (*Initiator, error) {

	path := InitiatorSearchUri + query

	res, err := this.Client.Search(path)
	if err != nil {
		return nil, err
	}

	this.id = res[0].Id

	return this.Query()
}
