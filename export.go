package coprhd

import ()

const (
	CreateExportUri = "block/exports.json"

	ExportTypeExclusive = "Exclusive"
)

type (
	ExportService struct {
		*Client
		itrs       []string
		project    string
		exportType ExportType
		array      string
		volumes    []ExportVolume
	}

	CreateExportReq struct {
		Initiators []string       `json:"initiators"`
		Name       string         `json:"name"`
		Project    string         `json:"project"`
		Type       ExportType     `json:"type"`
		VArray     string         `json:"varray"`
		Volumes    []ExportVolume `json:"volumes"`
	}

	CreateExportRes struct {
		Name string `json:"name"`
	}

	ExportVolume struct {
		Id string `json:"id"`
	}

	ExportType string
)

func (this *Client) Export() *ExportService {
	return &ExportService{
		Client:     this,
		itrs:       make([]string, 0),
		volumes:    make([]ExportVolume, 0),
		exportType: ExportTypeExclusive,
	}
}

func (this *ExportService) Initiators(itrs ...string) *ExportService {
	this.itrs = append(this.itrs, itrs...)
	return this
}

func (this *ExportService) Volumes(vols ...string) *ExportService {
	for _, v := range vols {
		this.volumes = append(this.volumes, ExportVolume{v})
	}
	return this
}

func (this *ExportService) Project(project string) *ExportService {
	this.project = project
	return this
}

func (this *ExportService) Array(array string) *ExportService {
	this.array = array
	return this
}

func (this *ExportService) Type(t ExportType) *ExportService {
	this.exportType = t
	return this
}

func (this *ExportService) Create(name string) (string, error) {
	req := CreateExportReq{
		Name:       name,
		Initiators: this.itrs,
		Project:    this.project,
		Type:       this.exportType,
		VArray:     this.array,
		Volumes:    this.volumes,
	}

	task := Task{}

	err := this.Post(CreateExportUri, &req, &task)
	if err != nil {
		return "", err
	}

	err = this.Task().WaitDone(task.Id, TaskStateReady)

	return task.Resource.Id, err
}
