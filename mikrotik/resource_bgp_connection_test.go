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

func TestAccMikrotikBgpConnection_basic(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-conn")
	instanceName := acctest.RandomWithPrefix("tf-acc-inst")
	routerId := internal.GetNewIpAddr()
	remoteAddr := internal.GetNewIpAddr()
	localAs := acctest.RandIntRange(64512, 65535)
	remoteAs := acctest.RandIntRange(64512, 65535)

	resourceName := "mikrotik_bgp_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpConnectionBasic(name, instanceName, localAs, remoteAs, routerId, remoteAddr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccBgpConnectionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "as", strconv.Itoa(localAs)),
					resource.TestCheckResourceAttr(resourceName, "remote_address", remoteAddr),
					resource.TestCheckResourceAttr(resourceName, "remote_as", strconv.Itoa(remoteAs)),
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

func TestAccMikrotikBgpConnection_update(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-conn-upd")
	instanceName := acctest.RandomWithPrefix("tf-acc-inst")
	routerId := internal.GetNewIpAddr()
	remoteAddr := internal.GetNewIpAddr()
	localAs := acctest.RandIntRange(64512, 65535)
	remoteAs := acctest.RandIntRange(64512, 65535)

	resourceName := "mikrotik_bgp_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpConnectionBasic(name, instanceName, localAs, remoteAs, routerId, remoteAddr),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "multihop", "false"),
					resource.TestCheckResourceAttr(resourceName, "use_bfd", "false"),
				),
			},
			{
				Config: testAccBgpConnectionUpdated(name, instanceName, localAs, remoteAs, routerId, remoteAddr),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "multihop", "true"),
					resource.TestCheckResourceAttr(resourceName, "use_bfd", "true"),
					resource.TestCheckResourceAttr(resourceName, "ttl", "255"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Updated connection"),
				),
			},
		},
	})
}

func TestAccMikrotikBgpConnection_withFilters(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-conn-flt")
	instanceName := acctest.RandomWithPrefix("tf-acc-inst")
	routerId := internal.GetNewIpAddr()
	remoteAddr := internal.GetNewIpAddr()
	localAs := acctest.RandIntRange(64512, 65535)
	remoteAs := acctest.RandIntRange(64512, 65535)

	resourceName := "mikrotik_bgp_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckMikrotikBgpConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpConnectionWithFilters(name, instanceName, localAs, remoteAs, routerId, remoteAddr),
				Check: resource.ComposeTestCheckFunc(
					testAccBgpConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "input_accept_nlri", "unicast"),
					resource.TestCheckResourceAttr(resourceName, "output_redistribute", "connected,static"),
				),
			},
		},
	})
}

func testAccBgpConnectionBasic(name, instanceName string, localAs, remoteAs int, routerId, remoteAddr string) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance_v7" "instance" {
  name      = "%s"
  as        = %d
  router_id = "%s"
}

resource "mikrotik_bgp_connection" "test" {
  name           = "%s"
  as             = %d
  instance       = mikrotik_bgp_instance_v7.instance.name
  remote_address = "%s"
  remote_as      = %d
}
`, instanceName, localAs, routerId, name, localAs, remoteAddr, remoteAs)
}

func testAccBgpConnectionUpdated(name, instanceName string, localAs, remoteAs int, routerId, remoteAddr string) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance_v7" "instance" {
  name      = "%s"
  as        = %d
  router_id = "%s"
}

resource "mikrotik_bgp_connection" "test" {
  name           = "%s"
  as             = %d
  instance       = mikrotik_bgp_instance_v7.instance.name
  remote_address = "%s"
  remote_as      = %d
  multihop       = true
  use_bfd        = true
  ttl            = "255"
  comment        = "Updated connection"
}
`, instanceName, localAs, routerId, name, localAs, remoteAddr, remoteAs)
}

func testAccBgpConnectionWithFilters(name, instanceName string, localAs, remoteAs int, routerId, remoteAddr string) string {
	return fmt.Sprintf(`
resource "mikrotik_bgp_instance_v7" "instance" {
  name      = "%s"
  as        = %d
  router_id = "%s"
}

resource "mikrotik_bgp_connection" "test" {
  name                 = "%s"
  as                   = %d
  instance             = mikrotik_bgp_instance_v7.instance.name
  remote_address       = "%s"
  remote_as            = %d
  input_accept_nlri    = "unicast"
  output_redistribute  = "connected,static"
}
`, instanceName, localAs, routerId, name, localAs, remoteAddr, remoteAs)
}

func testAccBgpConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for BGP Connection")
		}

		c := client.NewClient(client.GetConfigFromEnv())
		conn, err := c.FindBgpConnection(rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error finding BGP Connection: %v", err)
		}

		if conn == nil {
			return fmt.Errorf("BGP Connection not found")
		}

		return nil
	}
}

func testAccCheckMikrotikBgpConnectionDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_bgp_connection" {
			continue
		}

		conn, err := c.FindBgpConnection(rs.Primary.Attributes["name"])
		if !client.IsNotFoundError(err) && err != nil {
			return err
		}

		if conn != nil {
			return fmt.Errorf("BGP Connection (%s) still exists", conn.Name)
		}
	}
	return nil
}
