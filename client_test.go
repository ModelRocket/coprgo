package coprhd

import (
	"fmt"
	"testing"
)

const (
	TestUserName   = "root"
	TestPassword   = "@ChangeMe1"
	TestHost       = "https://54.153.46.171:4443/"
	TestProject    = "urn:storageos:Project:2dd5e0b5-4434-405a-8c5a-828daad17b3a:global"
	TestArray      = "urn:storageos:VirtualArray:7245ed6d-a4b9-4adf-9caa-901767586e1c:vdc1"
	TestPool       = "urn:storageos:VirtualPool:f943f75a-d610-4cea-9319-26f03c924e85:vdc1"
	TestVolumeName = "test_00"
	TestVolumeSize = 1024 * 1024 * 1024
)

var (
	proxyToken string
	testVolume string
)

func TestProxyToken(t *testing.T) {
	token, err := GetProxyToken(TestHost, TestUserName, TestPassword)
	if err != nil {
		t.Fatal("Failed to get a proxy token:", err.Error())
	}

	fmt.Printf("Proxytoken: %s", token)

	proxyToken = token
}

func TestEnumVolumes(t *testing.T) {
	client := NewClient(TestHost, proxyToken)

	vols, err := client.Volume().List()
	if err != nil {
		t.Fatal("Failed to get volume list:", err.Error())
	}

	for i, vol := range vols {
		fmt.Printf("Volume %d: %s\n", i, vol)
	}
}

func TestCreateVolume(t *testing.T) {
	client := NewClient(TestHost, proxyToken)
	vs := client.Volume().
		Project(TestProject).
		Array(TestArray).
		Pool(TestPool)

	id, err := vs.Create(TestVolumeName, TestVolumeSize)
	if err != nil {
		t.Fatal("Failed to create volume:", err.Error())
	}

	fmt.Printf("Created volume %s [%s]\n", id, TestVolumeName)

	testVolume = id
}

func TestDeleteVolume(t *testing.T) {
	client := NewClient(TestHost, proxyToken)
	err := client.Volume().
		Id(testVolume).
		Delete(true)

	if err != nil {
		t.Fatal("Failed to delete volume:", err.Error())
	}
}
