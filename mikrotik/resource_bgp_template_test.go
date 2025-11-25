package mikrotik

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMikrotikBgpTemplate_basic(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-tpl")
	as := acctest.RandIntRange(64512, 65535)

	resourceName := "mikrotik_bgp_template.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpTemplateBasic(name, as),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpTemplateExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "as", strconv.Itoa(as)),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMikrotikBgpTemplate_update(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-tpl-upd")
	as := acctest.RandIntRange(64512, 65535)

	resourceName := "mikrotik_bgp_template.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpTemplateBasic(name, as),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpTemplateExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "multihop", "false"),
					resource.TestCheckResourceAttr(resourceName, "route_reflect", "false"),
				),
			},
			{
				Config: testAccBgpTemplateUpdated(name, as),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpTemplateExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "multihop", "true"),
					resource.TestCheckResourceAttr(resourceName, "route_reflect", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_bfd", "true"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Route Reflector template"),
				),
			},
		},
	})
}

func TestAccMikrotikBgpTemplate_withFilters(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-tpl-flt")
	as := acctest.RandIntRange(64512, 65535)

	resourceName := "mikrotik_bgp_template.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpTemplateWithFilters(name, as),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpTemplateExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "input_accept_nlri", "unicast"),
					resource.TestCheckResourceAttr(resourceName, "input_accept_communities", "standard,extended"),
					resource.TestCheckResourceAttr(resourceName, "output_redistribute", "connected,static"),
				),
			},
		},
	})
}

func TestAccMikrotikBgpTemplate_withGracefulRestart(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-tpl-gr")
	as := acctest.RandIntRange(64512, 65535)

	resourceName := "mikrotik_bgp_template.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpTemplateWithGracefulRestart(name, as),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpTemplateExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "graceful_restart", "yes"),
				),
			},
		},
	})
}

func testAccBgpTemplateBasic(name string, as int) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_template" "test" {
  name = "%s"
  as   = %d
}
`, name, as)
}

func testAccBgpTemplateUpdated(name string, as int) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_template" "test" {
  name           = "%s"
  as             = %d
  multihop       = true
  route_reflect  = true
  use_bfd        = true
  comment        = "Route Reflector template"
}
`, name, as)
}

func testAccBgpTemplateWithFilters(name string, as int) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_template" "test" {
  name                      = "%s"
  as                        = %d
  input_accept_nlri         = "unicast"
  input_accept_communities  = "standard,extended"
  output_redistribute       = "connected,static"
}
`, name, as)
}

func testAccBgpTemplateWithGracefulRestart(name string, as int) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_template" "test" {
  name              = "%s"
  as                = %d
  graceful_restart  = "yes"
}
`, name, as)
}

func testAccBgpTemplateExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for BGP Template")
		}

		c := client.NewClient(client.GetConfigFromEnv())
		tpl, err := c.FindBgpTemplate(rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error finding BGP Template: %v", err)
		}

		if tpl == nil {
			return fmt.Errorf("BGP Template not found")
		}

		return nil
	}
}

func testAccCheckMikrotikBgpTemplateDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_bgp_template" {
			continue
		}

		tpl, err := c.FindBgpTemplate(rs.Primary.Attributes["name"])
		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if tpl != nil {
			return fmt.Errorf("BGP Template (%s) still exists", tpl.Name)
		}
	}
	return nil
}
