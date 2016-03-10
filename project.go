package coprhd

import (
	"fmt"
)

const (
	queryProjectUriTpl = "projects/%s.json"
	searchProjectUri   = "projects/search.json?"
)

type (
	ProjectService struct {
		*Client

		id   string
		name string
	}

	Project struct {
		StorageObject `json:",inline"`
	}
)

func (this *Client) Project() *ProjectService {
	return &ProjectService{
		Client: this,
	}
}

func (this *ProjectService) Id(id string) *ProjectService {
	this.id = id
	return this
}

func (this *ProjectService) Name(name string) *ProjectService {
	this.name = name
	return this
}

func (this *ProjectService) Query() (*Project, error) {
	if !isStorageOsUrn(this.id) {
		return this.Search("name=" + this.name)
	}

	path := fmt.Sprintf(queryProjectUriTpl, this.id)
	proj := Project{}

	err := this.get(path, nil, &proj)
	if err != nil {
		return nil, err
	}

	return &proj, nil
}

func (this *ProjectService) Search(query string) (*Project, error) {
	path := searchProjectUri + query

	res, err := this.Client.Search(path)
	if err != nil {
		return nil, err
	}

	this.id = res[0].Id

	return this.Query()
}
