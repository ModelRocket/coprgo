package coprhd

const (
	CreateVolumeUri = "block/volumes.json"
	ListVolumesUri  = "block/volumes/bulk.json"
)

type (
	VolumeService struct {
		// client is a pointer to the coprhd client
		*Client

		// project is the project for the volume command
		project string

		// varray is the Virtual Array for the volume command
		varray string

		// vpool is the Virtual Pool for the volume command
		vpool string

		// group is the consistency group id for the command
		group string
	}

	// CreateVolumeArgs represents the json parameters for the create volume REST call
	CreateVolumeArgs struct {
		ConsistencyGroup string `json:"consistency_group"`
		Count            int    `json:"count"`
		Name             string `json:"name"`
		Project          string `json:"project"`
		Size             string `json:"size"`
		VArray           string `json:"varray"`
		VPool            string `json:"vpool"`
	}

	// CreateVolumeReply is the reply from the create volume REST call
	CreateVolumeReply struct {
		Task []struct {
			Resource struct {
				Name string `json:"name"`
				Id   string `json:"id"`
			} `json:"resource"`
		} `json:"task"`
	}

	VolumeId string

	ListVolumesResponse struct {
		Volumes []VolumeId `json:"id"`
	}
)

func (this *Client) Volume() *VolumeService {
	return &VolumeService{
		Client: this,
	}
}

func (this *VolumeService) VArray(array string) *VolumeService {
	this.varray = array
	return this
}

func (this *VolumeService) VPool(pool string) *VolumeService {
	this.vpool = pool
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
	return "", nil
}

func (this *VolumeService) List() ([]VolumeId, error) {

	res := ListVolumesResponse{}

	err := this.Get(ListVolumesUri, nil, &res)
	if err != nil {
		return nil, err
	}
	return res.Volumes, nil
}
