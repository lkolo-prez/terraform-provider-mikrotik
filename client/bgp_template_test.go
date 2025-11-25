package client

import (
	"reflect"
	"testing"
)

func TestAddBgpTemplateAndDeleteBgpTemplate(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	tplName := "test-template"
	expectedTpl := &BgpTemplate{
		Name:         tplName,
		AS:           65530,
		Multihop:     false,
		UseBFD:       false,
		RouteReflect: false,
	}

	tpl, err := c.AddBgpTemplate(expectedTpl)
	if err != nil {
		t.Fatalf("Error creating BGP template: %v", err)
	}

	expectedTpl.Id = tpl.Id

	if !reflect.DeepEqual(tpl, expectedTpl) {
		t.Errorf("BGP template does not match expected. actual: %v expected: %v", tpl, expectedTpl)
	}

	err = c.DeleteBgpTemplate(tpl.Name)
	if err != nil {
		t.Errorf("Error deleting BGP template: %v", err)
	}
}

func TestAddAndUpdateBgpTemplateWithOptionalFields(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	tplName := "test-template-update"
	expectedTpl := &BgpTemplate{
		Name:         tplName,
		AS:           65531,
		Multihop:     false,
		UseBFD:       false,
		RouteReflect: false,
		Comment:      "Initial",
	}

	tpl, err := c.AddBgpTemplate(expectedTpl)
	if err != nil {
		t.Fatalf("Error creating BGP template: %v", err)
	}
	defer func() {
		_ = c.DeleteBgpTemplate(tplName)
	}()

	// Update fields
	tpl.Multihop = true
	tpl.UseBFD = true
	tpl.RouteReflect = true
	tpl.Comment = "Updated"

	updatedTpl, err := c.UpdateBgpTemplate(tpl)
	if err != nil {
		t.Errorf("Error updating BGP template: %v", err)
	}

	if !reflect.DeepEqual(updatedTpl, tpl) {
		t.Errorf("Updated BGP template does not match expected. actual: %v expected: %v", updatedTpl, tpl)
	}
}

func TestFindBgpTemplate(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	tplName := "test-template-find"
	expectedTpl := &BgpTemplate{
		Name: tplName,
		AS:   65532,
	}

	tpl, err := c.AddBgpTemplate(expectedTpl)
	if err != nil {
		t.Fatalf("Error creating BGP template: %v", err)
	}
	defer func() {
		_ = c.DeleteBgpTemplate(tplName)
	}()

	foundTpl, err := c.FindBgpTemplate(tplName)
	if err != nil {
		t.Errorf("Error finding BGP template: %v", err)
	}

	if foundTpl.Name != tpl.Name {
		t.Errorf("Found BGP template name does not match. actual: %s expected: %s", foundTpl.Name, tpl.Name)
	}
}

func TestBgpTemplateWithGracefulRestart(t *testing.T) {
	SkipIfRouterOSV6OrEarlier(t, sysResources)
	c := NewClient(GetConfigFromEnv())

	tplName := "test-template-gr"
	expectedTpl := &BgpTemplate{
		Name:            tplName,
		AS:              65533,
		GracefulRestart: "yes",
		Comment:         "Graceful restart enabled",
	}

	tpl, err := c.AddBgpTemplate(expectedTpl)
	if err != nil {
		t.Fatalf("Error creating BGP template with graceful restart: %v", err)
	}
	defer func() {
		_ = c.DeleteBgpTemplate(tplName)
	}()

	if tpl.GracefulRestart != expectedTpl.GracefulRestart {
		t.Errorf("Graceful restart does not match. actual: %s expected: %s", tpl.GracefulRestart, expectedTpl.GracefulRestart)
	}
}
