package coprhd

import (
	"fmt"
	"time"
)

const (
	CreateHostUri      = "compute/hosts.json"
	QueryHostItrUriTpl = "compute/hosts/%s/initiators.json"
	SearchHostUri      = "compute/hosts/search.json?"
	QueryHostUriTpl    = "compute/hosts/%s.json"

	HostTypeLinux   HostType = "Linux"
	HostTypeWindows HostType = "Windows"
	HostTypeHPUX    HostType = "HPUX"
	HostTypeEsx     HostType = "Esx"
	HostTypeOther   HostType = "Other"
)

type (
	HostService struct {
		*Client
		id     string
		name   string
		typ    HostType
		os     string
		tenant string
	}

	Host struct {
		BaseObject         `json:",inline"`
		Type               HostType `json:"type"`
		OSVersion          string   `json:"os_version,omitempty"`
		HostName           string   `json:"host_name"`
		Port               int      `json:"port_number,omitempty"`
		Username           string   `json:"user_name,omitempty"`
		SSL                bool     `json:"use_ssl,omitempty"`
		Discoverable       bool     `json:"discoverable"`
		RegistrationStatus string   `json:"registration_status"`
		Tenant             Resource `json:"tenant"`
		Cluster            Resource `json:"cluster,omitempty"`
	}

	CreateHostReq struct {
		Name         string   `json:"name"`
		Type         HostType `json:"type"`
		OSVersion    string   `json:"os_version,omitempty"`
		HostName     string   `json:"host_name"`
		Port         int      `json:"port_number,omitempty"`
		Tenant       string   `json:"tenant"`
		SSL          bool     `json:"use_ssl,omitempty"`
		Discoverable bool     `json:"discoverable"`
		Username     string   `json:"user_name"`
		Password     string   `json:"password"`
	}

	QueryHostItrRes struct {
		Initiators []NamedResource `json:"initiator"`
	}

	HostType string
)

func (this *Client) Host() *HostService {
	return &HostService{
		Client: this.Copy(),
	}
}

func (this *HostService) Id(id string) *HostService {
	this.id = id
	return this
}

func (this *HostService) Name(name string) *HostService {
	this.name = name
	return this
}

func (this *HostService) Tenant(id string) *HostService {
	this.tenant = id
	return this
}

func (this *HostService) Type(t HostType) *HostService {
	this.typ = t
	return this
}

func (this *HostService) OSVersion(v string) *HostService {
	this.os = v
	return this
}

// Create creates a new host with the name and host
func (this *HostService) Create(host string) (*Host, error) {
	req := CreateHostReq{
		Name:         this.name,
		HostName:     host,
		Discoverable: false,
		Type:         this.typ,
		OSVersion:    this.os,
		Tenant:       this.tenant,
	}

	task := Task{}

	err := this.Post(CreateHostUri, &req, &task)
	if err != nil {
		if this.LastError().IsCreateHostDup() {
			return this.Query()
		}
		return nil, err
	}

	err = this.Task().WaitDone(task.Id, TaskStateReady, time.Second*180)
	if err != nil {
		return nil, err
	}

	this.id = task.Resource.Id

	return this.Query()
}

// Discover creates and attempts to discover a new host
func (this *HostService) Discover(host, username, password string, port int, ssl bool) (*Host, error) {
	req := CreateHostReq{
		Name:         this.name,
		HostName:     host,
		Port:         port,
		Discoverable: true,
		Username:     username,
		Password:     password,
		SSL:          ssl,
		Type:         this.typ,
		OSVersion:    this.os,
		Tenant:       this.tenant,
	}

	task := Task{}

	err := this.Post(CreateHostUri, &req, &task)
	if err != nil {
		if this.LastError().IsCreateHostDup() {
			return this.Query()
		}
		return nil, err
	}

	err = this.Task().WaitDone(task.Id, TaskStateReady, time.Second*180)
	if err != nil {
		return nil, err
	}

	this.id = task.Resource.Id

	return this.Query()
}

func (this *HostService) Query() (*Host, error) {
	if !isStorageOsUrn(this.id) {
		return this.Search("name=" + this.name)
	}

	path := fmt.Sprintf(QueryHostUriTpl, this.id)
	host := Host{}

	err := this.Get(path, nil, &host)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func (this *HostService) Search(query string) (*Host, error) {
	path := SearchHostUri + query

	res, err := this.Client.Search(path)
	if err != nil {
		return nil, err
	}

	this.id = res[0].Id

	return this.Query()
}

func (this *HostService) Delete(id string) error {
	path := fmt.Sprintf(DeleteExportUriTpl, id)

	task := Task{}

	err := this.Post(path, nil, &task)
	if err != nil {
		return err
	}

	return this.Task().WaitDone(task.Id, TaskStateReady, time.Second*180)
}

func (this *HostService) Initiators() ([]Initiator, error) {
	if err := this.queryHostByName(); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(QueryHostItrUriTpl, this.id)
	res := QueryHostItrRes{}
	itrs := make([]Initiator, 0)

	err := this.Get(path, nil, &res)
	if err != nil {
		return nil, err
	}

	for _, i := range res.Initiators {
		itr, err := this.Initiator().
			Id(i.Id).
			Query()

		if err != nil {
			return itrs, err
		}

		itrs = append(itrs, *itr)
	}

	return itrs, nil
}

func (this *HostService) queryHostByName() error {
	if !isStorageOsUrn(this.id) {
		host, err := this.Query()
		if err != nil {
			return err
		}
		this.id = host.Id
		this.name = host.Name
	}
	return nil
}
