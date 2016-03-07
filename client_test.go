package coprhd

import (
	"fmt"
	"testing"
)

const (
	TestUserName = "root"
	TestPassword = "D1g1tal*23"
	TestHost     = "https://54.153.46.171:4443/"
)

var (
	proxyToken string
)

func init() {
}

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
