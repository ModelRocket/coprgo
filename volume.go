package coprhd

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	CreateVolumeUri   = "block/volumes.json"
	QueryVolumeUriTpl = "block/volumes/%s.json"
	SearchVolumeUri   = "block/volumes/search.json?"
	ListVolumesUri    = "block/volumes/bulk.json"
	DeleteVolUriTpl   = "block/volumes/%s/deactivate.json"
)

var (
	ErrCreateResponse = errors.New("Invalid create response received")
)

type (
	// VolumeService is used to create, search, and query for volumes
	VolumeService struct {
		*Client
		id      string
		array   string
		pool    string
		group   string
		project string
	}

	// Volume is a complete coprhd volume object
	Volume struct {
		BaseObject          `json:",inline"`
		WWN                 string      `json:"wwn"`
		Protocols           []string    `json:"protocols"`
		Protection          interface{} `json:"protection"`
		ConsistencyGroup    string      `json:"consistency_group,omitempty"`
		StorageController   string      `json:"storage_controller"`
		DeviceLabel         string      `json:"device_label"`
		NativeId            string      `json:"native_id"`
		ProvisionedCapacity string      `json:"provisioned_capacity_gb"`
		AllocatedCapacity   string      `json:"allocated_capacity_gb"`
		RequestedCapacity   string      `json:"requested_capacity_gb"`
		PreAllocationSize   string      `json:"pre_allocation_size_gb"`
		IsComposite         bool        `json:"is_composite"`
		ThinlyProvisioned   bool        `json:"thinly_provisioned"`
		HABackingVolumes    []string    `json:"high_availability_backing_volumes"`
		AccessState         string      `json:"access_state"`
		StoragePool         Resource    `json:"storage_pool"`
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

	// ListVolumesRes is the reply to geting a list of volumes
	ListVolumesRes struct {
		Volumes []string `json:"id"`
	}
)

func (this *Client) Volume() *VolumeService {
	return &VolumeService{
		Client: this.Copy(),
	}
}

func (this *VolumeService) Id(id string) *VolumeService {
	// make sure Volume is capitalized
	this.id = strings.Replace(id, "volume", "Volume", 1)
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
func (this *VolumeService) Create(name string, size uint64) (*Volume, error) {
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
		return nil, err
	}

	if len(res.Task) != 1 {
		return nil, ErrCreateResponse
	}

	task := res.Task[0]

	// wait for the task to complete
	err = this.Task().WaitDone(task.Id, TaskStateReady, time.Second*180)
	if err != nil {
		return nil, err
	}

	this.id = task.Resource.Id

	return this.Query()
}

func (this *VolumeService) Query() (*Volume, error) {
	path := fmt.Sprintf(QueryVolumeUriTpl, this.id)
	vol := Volume{}

	err := this.Get(path, nil, &vol)
	if err != nil {
		return nil, err
	}

	return &vol, nil
}

func (this *VolumeService) Search(query string) (*Volume, error) {
	path := SearchVolumeUri + query

	res, err := this.Client.Search(path)
	if err != nil {
		return nil, err
	}

	this.id = res[0].Id

	return this.Query()
}

func (this *VolumeService) List() ([]string, error) {

	res := ListVolumesRes{}

	err := this.Get(ListVolumesUri, nil, &res)
	if err != nil {
		return nil, err
	}
	return res.Volumes, nil
}

func (this *VolumeService) Delete(force bool) error {
	path := fmt.Sprintf(DeleteVolUriTpl, this.id)

	if force {
		path = path + "?force=true"
	}

	task := Task{}

	err := this.Post(path, nil, &task)
	if err != nil {
		return err
	}

	return this.Task().WaitDone(task.Id, TaskStateReady, time.Second*180)
}
