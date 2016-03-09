package coprhd

import (
	"fmt"
)

const (
	ProjectQueryUriTpl = "projects/%s.json"
	ProjectSearchUri   = "projects/search.json?"
)

type (
	ProjectService struct {
		*Client

		id string
	}

	Project struct {
		BaseObject `json:",inline"`
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

func (this *ProjectService) Query() (*Project, error) {
	path := fmt.Sprintf(ProjectQueryUriTpl, this.id)
	proj := Project{}

	err := this.Get(path, nil, &proj)
	if err != nil {
		return nil, err
	}

	return &proj, nil
}

func (this *ProjectService) Search(query string) (*Project, error) {
	path := ProjectSearchUri + query

	res, err := this.Client.Search(path)
	if err != nil {
		return nil, err
	}

	this.id = res[0].Id

	return this.Query()
}
