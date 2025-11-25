package client

import (
	"reflect"
	"testing"
)

func TestAddBgpInstanceV7AndDeleteBgpInstanceV7(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	bgpInstanceName := "test-instance-v7"
	routerID := "10.255.255.1"

	expectedBgpInstance := &BgpInstanceV7{
		Name:                     bgpInstanceName,
		AS:                       65530,
		RouterID:                 routerID,
		ClientToClientReflection: true,
		RedistributeConnected:    true,
		RedistributeStatic:       true,
	}

	bgpInstance, err := c.AddBgpInstanceV7(expectedBgpInstance)
	if err != nil {
		t.Fatalf("Error creating BGP instance v7: %v", err)
	}

	expectedBgpInstance.Id = bgpInstance.Id

	if !reflect.DeepEqual(bgpInstance, expectedBgpInstance) {
		t.Errorf("BGP instance v7 does not match expected. actual: %v expected: %v", bgpInstance, expectedBgpInstance)
	}

	err = c.DeleteBgpInstanceV7(bgpInstance.Name)
	if err != nil {
		t.Errorf("Error deleting BGP instance v7: %v", err)
	}
}

func TestAddAndUpdateBgpInstanceV7WithOptionalFieldsAndDeleteBgpInstanceV7(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	bgpInstanceName := "test-instance-v7-update"
	routerID := "10.255.255.2"

	expectedBgpInstance := &BgpInstanceV7{
		Name:                  bgpInstanceName,
		AS:                    65531,
		RouterID:              routerID,
		RedistributeConnected: false,
		RedistributeStatic:    false,
		RedistributeOspf:      false,
		Comment:               "Initial config",
	}

	bgpInstance, err := c.AddBgpInstanceV7(expectedBgpInstance)
	if err != nil {
		t.Fatalf("Error creating BGP instance v7: %v", err)
	}

	expectedBgpInstance.Id = bgpInstance.Id

	if !reflect.DeepEqual(bgpInstance, expectedBgpInstance) {
		t.Errorf("BGP instance v7 does not match expected. actual: %v expected: %v", bgpInstance, expectedBgpInstance)
	}

	// Update fields
	bgpInstance.RedistributeConnected = true
	bgpInstance.RedistributeStatic = true
	bgpInstance.Comment = "Updated config"

	updatedBgpInstance, err := c.UpdateBgpInstanceV7(bgpInstance)
	if err != nil {
		t.Errorf("Error updating BGP instance v7: %v", err)
	}

	if !reflect.DeepEqual(updatedBgpInstance, bgpInstance) {
		t.Errorf("Updated BGP instance v7 does not match expected. actual: %v expected: %v", updatedBgpInstance, bgpInstance)
	}

	err = c.DeleteBgpInstanceV7(bgpInstance.Name)
	if err != nil {
		t.Errorf("Error deleting BGP instance v7: %v", err)
	}
}

func TestFindBgpInstanceV7(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	bgpInstanceName := "test-instance-v7-find"
	routerID := "10.255.255.3"

	expectedBgpInstance := &BgpInstanceV7{
		Name:     bgpInstanceName,
		AS:       65532,
		RouterID: routerID,
	}

	bgpInstance, err := c.AddBgpInstanceV7(expectedBgpInstance)
	if err != nil {
		t.Fatalf("Error creating BGP instance v7: %v", err)
	}
	defer func() {
		_ = c.DeleteBgpInstanceV7(bgpInstanceName)
	}()

	foundBgpInstance, err := c.FindBgpInstanceV7(bgpInstanceName)
	if err != nil {
		t.Errorf("Error finding BGP instance v7: %v", err)
	}

	if foundBgpInstance.Name != bgpInstance.Name {
		t.Errorf("Found BGP instance v7 name does not match. actual: %s expected: %s", foundBgpInstance.Name, bgpInstance.Name)
	}
}
