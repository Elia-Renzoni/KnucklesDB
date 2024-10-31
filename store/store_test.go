package store_test

import (
	"knucklesdb/store"
	"net"
	"testing"
)

func TestSetWithIpOnly(t *testing.T) {
	ip := net.IPv4(129, 35, 67, 89)
	intance := store.NewDBValues(ip, 4560, 30, "")
	storeInstance := store.NewKnucklesDB()

	if err := storeInstance.SetWithIpAddressOnly(ip.String(), intance); err != nil {
		t.Fail()
	}
}

func TestSetWithEndpointOnly(t *testing.T) {
	endpoint := "/delete"
	instance := store.NewDBValues(nil, 8080, 30, endpoint)
	storeInstance := store.NewKnucklesDB()

	if err := storeInstance.SetWithEndpointOnly(endpoint, instance); err != nil {
		t.Fail()
	}
}

func TestGetWithIpOnly(t *testing.T) {
	ip := net.IPv4(192, 56, 0, 0)
	values := store.NewDBValues(ip, 4040, 20, "")
	storeInstance := store.NewKnucklesDB()

	storeInstance.SetWithIpAddressOnly(ip.String(), values)

	values2, _ := storeInstance.SearchWithIpOnly(ip.String())

	// random values to check
	if !(values.GetListenPort() == values2.GetListenPort()) {
		t.Fail()
	}
}

func TestGetWithEndpointOnly(t *testing.T) {
	endpoint := "/test"
	values := store.NewDBValues(nil, 4040, 20, endpoint)
	storeInstance := store.NewKnucklesDB()

	storeInstance.SetWithEndpointOnly(endpoint, values)
	values2, _ := storeInstance.SearchWithEndpointOnly(endpoint)

	if !(values.GetOptionalEndpoint() == values2.GetOptionalEndpoint()) {
		t.Fail()
	}
}

func TestReturnEntries(t *testing.T) {
	ip := net.IPv4(13, 66, 77, 00)
	endpoint := "/update"
	ip2 := net.IPv4(55, 66, 44, 22)
	endpoint2 := "/remove"

	valuesIP1 := store.NewDBValues(ip, 8080, 30, "")
	vluesIP2 := store.NewDBValues(ip2, 8089, 23, "")
	valuesEndpoint1 := store.NewDBValues(nil, 7070, 2, endpoint)
	valuesEndpoint2 := store.NewDBValues(nil, 4040, 21, endpoint2)

	storeInstance := store.NewKnucklesDB()
	storeInstance.SetWithIpAddressOnly(ip.String(), valuesIP1)
	storeInstance.SetWithIpAddressOnly(ip2.String(), vluesIP2)
	storeInstance.SetWithEndpointOnly(endpoint, valuesEndpoint1)
	storeInstance.SetWithEndpointOnly(endpoint2, valuesEndpoint2)
	
	entries := storeInstance.ReturnEntries()

	if len(entries) <= 0 {
		t.Fail()
	}
}
