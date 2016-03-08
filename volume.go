package coprhd

import (
	"errors"
	"fmt"
)

const (
	CreateVolumeUri = "block/volumes.json"
	ListVolumesUri  = "block/volumes/bulk.json"
	DeleteVolUrlTpl = "block/volumes/%s/deactivate.json"
)

var (
	ErrCreateResponse = errors.New("Invalid create response received")
)

type (
	VolumeService struct {
		*Client

		// id is the volume id
		id string

		// project is the project for the volume command
		project string

		// varray is the Virtual Array for the volume command
		array string

		// vpool is the Virtual Pool for the volume command
		pool string

		// group is the consistency group id for the command
		group string
	}

	// CreateVolumeReq represents the json parameters for the create volume REST call
	CreateVolumeReq struct {
		ConsistencyGroup string `json:"consistency_group,omitempty"`
		Count            int    `json:"count"`
		Name             string `json:"name"`
		Project          string `json:"project"`
		Size             string `json:"size"`
		VArray           string `json:"varray"`
		VPool            string `json:"vpool"`
	}

	// CreateVolumeRes is the reply from the create volume REST call
	CreateVolumeRes struct {
		Task []Task `json:"task"`
	}

	VolumeId string

	ListVolumesRes struct {
		Volumes []VolumeId `json:"id"`
	}
)

func (this *Client) Volume() *VolumeService {
	return &VolumeService{
		Client: this,
	}
}

func (this *VolumeService) Id(id string) *VolumeService {
	this.id = id
	return this
}

func (this *VolumeService) Array(array string) *VolumeService {
	this.array = array
	return this
}

func (this *VolumeService) Pool(pool string) *VolumeService {
	this.pool = pool
	return this
}

func (this *VolumeService) Group(group string) *VolumeService {
	this.group = group
	return this
}

func (this *VolumeService) Project(project string) *VolumeService {
	this.project = project
	return this
}

// CreateVolume creates a new volume with the specified name using the volume service
func (this *VolumeService) Create(name string, size int64) (string, error) {
	sz := float64(size / (1024 * 1024 * 1000))

	req := CreateVolumeReq{
		Count:   1,
		Name:    name,
		Project: this.project,
		VArray:  this.array,
		VPool:   this.pool,
		Size:    fmt.Sprintf("%.6fGB", sz),
	}

	if this.group != "" {
		req.ConsistencyGroup = this.group
	}

	res := CreateVolumeRes{}

	err := this.Post(CreateVolumeUri, &req, &res)
	if err != nil {
		return "", err
	}

	if len(res.Task) != 1 {
		return "", ErrCreateResponse
	}

	task := res.Task[0]
	this.id = task.Resource.Id

	err = this.Task().WaitDone(task.Id, TaskStateReady)

	return this.id, err
}

func (this *VolumeService) List() ([]VolumeId, error) {

	res := ListVolumesRes{}

	err := this.Get(ListVolumesUri, nil, &res)
	if err != nil {
		return nil, err
	}
	return res.Volumes, nil
}

func (this *VolumeService) Delete(force bool) error {
	path := fmt.Sprintf(DeleteVolUrlTpl, this.id)

	if force {
		path = path + "?force=true"
	}

	return this.Post(path, nil, nil)
}
