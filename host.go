package coprhd

import ()

const (
	CreateHostUri = "compute/hosts.json"

	HostTypeLinux   HostType = "Linux"
	HostTypeWindows HostType = "Windows"
)

type (
	HostService struct {
		*Client
		name     string
		hostType HostType
		hostName string
	}

	HostType string
)

func (this *Client) Host() *HostService {
	return &HostService{
		Client: this,
	}
}
