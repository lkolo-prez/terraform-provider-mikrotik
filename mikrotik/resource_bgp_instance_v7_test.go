package mikrotik

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMikrotikBgpInstanceV7_basic(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-bgpv7")
	routerId := internal.GetNewIpAddr()
	as := acctest.RandIntRange(64512, 65535)

	resourceName := "mikrotik_bgp_instance_v7.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpInstanceV7Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstanceV7Basic(name, as, routerId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpInstanceV7Exists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "as", strconv.Itoa(as)),
					resource.TestCheckResourceAttr(resourceName, "router_id", routerId),
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

func TestAccMikrotikBgpInstanceV7_update(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-bgpv7-upd")
	routerId := internal.GetNewIpAddr()
	routerIdUpdated := internal.GetNewIpAddr()
	as := acctest.RandIntRange(64512, 65535)

	resourceName := "mikrotik_bgp_instance_v7.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpInstanceV7Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstanceV7Basic(name, as, routerId),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpInstanceV7Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "router_id", routerId),
					resource.TestCheckResourceAttr(resourceName, "client_to_client_reflection", "true"),
				),
			},
			{
				Config: testAccBgpInstanceV7Updated(name, as, routerIdUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpInstanceV7Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "router_id", routerIdUpdated),
					resource.TestCheckResourceAttr(resourceName, "client_to_client_reflection", "false"),
					resource.TestCheckResourceAttr(resourceName, "redistribute_connected", "true"),
					resource.TestCheckResourceAttr(resourceName, "redistribute_static", "true"),
				),
			},
		},
	})
}

func TestAccMikrotikBgpInstanceV7_withVrf(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-bgpv7-vrf")
	vrfName := acctest.RandomWithPrefix("test-vrf")
	routerId := internal.GetNewIpAddr()
	as := acctest.RandIntRange(64512, 65535)

	resourceName := "mikrotik_bgp_instance_v7.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpInstanceV7Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpInstanceV7WithVrf(name, as, routerId, vrfName),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpInstanceV7Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vrf", vrfName),
					resource.TestCheckResourceAttr(resourceName, "routing_table", vrfName),
				),
			},
		},
	})
}

func testAccBgpInstanceV7Basic(name string, as int, routerId string) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance_v7" "test" {
  name       = "%s"
  as         = %d
  router_id  = "%s"
}
`, name, as, routerId)
}

func testAccBgpInstanceV7Updated(name string, as int, routerId string) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance_v7" "test" {
  name                        = "%s"
  as                          = %d
  router_id                   = "%s"
  client_to_client_reflection = false
  redistribute_connected      = true
  redistribute_static         = true
  comment                     = "Updated instance"
}
`, name, as, routerId)
}

func testAccBgpInstanceV7WithVrf(name string, as int, routerId, vrfName string) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance_v7" "test" {
  name          = "%s"
  as            = %d
  router_id     = "%s"
  vrf           = "%s"
  routing_table = "%s"
}
`, name, as, routerId, vrfName, vrfName)
}

func testAccBgpInstanceV7Exists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for BGP Instance V7")
		}

		c := client.NewClient(client.GetConfigFromEnv())
		instance, err := c.FindBgpInstanceV7(rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error finding BGP Instance V7: %v", err)
		}

		if instance == nil {
			return fmt.Errorf("BGP Instance V7 not found")
		}

		return nil
	}
}

func testAccCheckMikrotikBgpInstanceV7Destroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_bgp_instance_v7" {
			continue
		}

		instance, err := c.FindBgpInstanceV7(rs.Primary.Attributes["name"])
		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if instance != nil {
			return fmt.Errorf("BGP Instance V7 (%s) still exists", instance.Name)
		}
	}
	return nil
}
