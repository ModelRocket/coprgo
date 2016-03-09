package test

import (
	"fmt"
	. "github.com/ModelRocket/coprhd"
	"testing"
)

const (
	TestUserName   = "root"
	TestPassword   = "@ChangeMe1"
	TestHost       = "https://54.153.46.171:4443/"
	TestProject    = "default"
	TestArray      = "default"
	TestPool       = "default"
	TestVolumeName = "test_00"
	TestVolumeSize = 1024 * 1024 * 1024

	TestHostName  = "Test Host"
	TestInitiator = "iqn.1994-05.com.redhat:98d5cd397a18"
)

var (
	proxyToken string
	testVolume string
	testExport string
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
		Name(TestVolumeName).
		Project(TestProject).
		Array(TestArray).
		Pool(TestPool)

	vol, err := vs.Create(TestVolumeSize)
	if err != nil {
		t.Fatal("Failed to create volume:", err.Error())
	}

	fmt.Printf("Created volume %s %s\n", TestVolumeName, vol.Id)

	testVolume = vol.Id
}

func TestExportVolume(t *testing.T) {
	client := NewClient(TestHost, proxyToken)
	es := client.Export().
		Name(TestVolumeName).
		Initiators(TestInitiator).
		Volumes(testVolume).
		Project(TestProject).
		Array(TestArray)

	export, err := es.Create()

	if err != nil {
		t.Fatal("Failed to export volume:", err.Error())
	}

	fmt.Printf("Created export group %s\n", export.Id)

	testExport = export.Id
}

func TestExportDelete(t *testing.T) {
	client := NewClient(TestHost, proxyToken)
	err := client.Export().
		Delete(testExport)

	if err != nil {
		t.Fatal("Failed to delete export:", err.Error())
	}
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

func TestQueryHost(t *testing.T) {
	client := NewClient(TestHost, proxyToken)
	host, err := client.Host().
		Name(TestHostName).
		Query()

	if err != nil {
		t.Fatal("Failed to query host", err.Error())
	}

	fmt.Printf("Got host %s [%s]\n", host.Id, host.Name)
}

func TestQueryHostInitiators(t *testing.T) {
	client := NewClient(TestHost, proxyToken)
	itrs, err := client.Host().
		Name(TestHostName).
		Initiators()

	if err != nil {
		t.Fatal("Failed to query host initiators:", err.Error())
	}

	if len(itrs) > 0 {
		fmt.Printf("Got %d initiators => %s\n", len(itrs), itrs[0].Port)
	}
}
