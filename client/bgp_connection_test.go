package client

import (
	"reflect"
	"testing"
)

func TestAddBgpConnectionAndDeleteBgpConnection(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	// Create BGP instance first
	instanceName := "conn-test-instance"
	_, err := c.AddBgpInstanceV7(&BgpInstanceV7{
		Name:     instanceName,
		AS:       65530,
		RouterID: "10.255.255.10",
	})
	if err != nil {
		t.Fatalf("unable to create BGP instance v7: %v", err)
	}
	defer func() {
		_ = c.DeleteBgpInstanceV7(instanceName)
	}()

	// Create connection
	connName := "test-conn"
	expectedConn := &BgpConnection{
		Name:          connName,
		AS:            65530,
		Instance:      instanceName,
		RemoteAddress: "192.168.1.1",
		RemoteAS:      65531,
		Multihop:      false,
		UseBFD:        false,
	}

	conn, err := c.AddBgpConnection(expectedConn)
	if err != nil {
		t.Fatalf("Error creating BGP connection: %v", err)
	}

	expectedConn.Id = conn.Id

	if !reflect.DeepEqual(conn, expectedConn) {
		t.Errorf("BGP connection does not match expected. actual: %v expected: %v", conn, expectedConn)
	}

	err = c.DeleteBgpConnection(conn.Name)
	if err != nil {
		t.Errorf("Error deleting BGP connection: %v", err)
	}
}

func TestAddAndUpdateBgpConnectionWithOptionalFields(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	// Create BGP instance first
	instanceName := "conn-update-test-instance"
	_, err := c.AddBgpInstanceV7(&BgpInstanceV7{
		Name:     instanceName,
		AS:       65532,
		RouterID: "10.255.255.11",
	})
	if err != nil {
		t.Fatalf("unable to create BGP instance v7: %v", err)
	}
	defer func() {
		_ = c.DeleteBgpInstanceV7(instanceName)
	}()

	// Create connection
	connName := "test-conn-update"
	expectedConn := &BgpConnection{
		Name:          connName,
		AS:            65532,
		Instance:      instanceName,
		RemoteAddress: "192.168.2.1",
		RemoteAS:      65533,
		Multihop:      false,
		UseBFD:        false,
		Comment:       "Initial",
	}

	conn, err := c.AddBgpConnection(expectedConn)
	if err != nil {
		t.Fatalf("Error creating BGP connection: %v", err)
	}
	defer func() {
		_ = c.DeleteBgpConnection(connName)
	}()

	// Update fields
	conn.Multihop = true
	conn.UseBFD = true
	conn.TTL = "255"
	conn.Comment = "Updated"

	updatedConn, err := c.UpdateBgpConnection(conn)
	if err != nil {
		t.Errorf("Error updating BGP connection: %v", err)
	}

	if !reflect.DeepEqual(updatedConn, conn) {
		t.Errorf("Updated BGP connection does not match expected. actual: %v expected: %v", updatedConn, conn)
	}
}

func TestFindBgpConnection(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	// Create BGP instance first
	instanceName := "conn-find-test-instance"
	_, err := c.AddBgpInstanceV7(&BgpInstanceV7{
		Name:     instanceName,
		AS:       65534,
		RouterID: "10.255.255.12",
	})
	if err != nil {
		t.Fatalf("unable to create BGP instance v7: %v", err)
	}
	defer func() {
		_ = c.DeleteBgpInstanceV7(instanceName)
	}()

	// Create connection
	connName := "test-conn-find"
	expectedConn := &BgpConnection{
		Name:          connName,
		AS:            65534,
		Instance:      instanceName,
		RemoteAddress: "192.168.3.1",
		RemoteAS:      65535,
	}

	conn, err := c.AddBgpConnection(expectedConn)
	if err != nil {
		t.Fatalf("Error creating BGP connection: %v", err)
	}
	defer func() {
		_ = c.DeleteBgpConnection(connName)
	}()

	foundConn, err := c.FindBgpConnection(connName)
	if err != nil {
		t.Errorf("Error finding BGP connection: %v", err)
	}

	if foundConn.Name != conn.Name {
		t.Errorf("Found BGP connection name does not match. actual: %s expected: %s", foundConn.Name, conn.Name)
	}
}
