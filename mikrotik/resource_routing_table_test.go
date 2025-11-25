package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRoutingTableResource_basic(t *testing.T) {
	resourceName := "mikrotik_routing_table.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRoutingTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingTableConfig("customer_a", "main", false, "Customer A VRF"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRoutingTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "customer_a"),
					resource.TestCheckResourceAttr(resourceName, "fib", "main"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Customer A VRF"),
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

func TestAccRoutingTableResource_update(t *testing.T) {
	resourceName := "mikrotik_routing_table.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRoutingTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingTableConfig("test_vrf", "main", false, "Initial comment"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRoutingTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Initial comment"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
				),
			},
			{
				Config: testAccRoutingTableConfig("test_vrf", "main", true, "Updated comment"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRoutingTableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", "Updated comment"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
				),
			},
		},
	})
}

func TestAccRoutingTableResource_multipleVRFs(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRoutingTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingTableMultipleVRFsConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("mikrotik_routing_table.customer_a", "name", "customer_a"),
					resource.TestCheckResourceAttr("mikrotik_routing_table.customer_b", "name", "customer_b"),
					resource.TestCheckResourceAttr("mikrotik_routing_table.management", "name", "management"),
				),
			},
		},
	})
}

func testAccRoutingTableConfig(name, fib string, disabled bool, comment string) string {
	return fmt.Sprintf(`
resource "mikrotik_routing_table" "test" {
  name     = %q
  fib      = %q
  disabled = %t
  comment  = %q
}
`, name, fib, disabled, comment)
}

func testAccRoutingTableMultipleVRFsConfig() string {
	return `
resource "mikrotik_routing_table" "customer_a" {
  name    = "customer_a"
  fib     = "main"
  comment = "Customer A VRF"
}

resource "mikrotik_routing_table" "customer_b" {
  name    = "customer_b"
  fib     = "main"
  comment = "Customer B VRF"
}

resource "mikrotik_routing_table" "management" {
  name    = "management"
  fib     = "main"
  comment = "Management VRF"
}
`
}

func testAccCheckRoutingTableExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for %s", resourceName)
		}

		c := client.NewClient(client.GetConfigFromEnv())

		_, err := c.FindRoutingTable(rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error finding routing table: %s", err)
		}

		return nil
	}
}

func testAccCheckRoutingTableDestroy(s *terraform.State) error {
	c := client.NewClient(client.GetConfigFromEnv())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_routing_table" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		_, err := c.FindRoutingTable(name)
		if err == nil {
			return fmt.Errorf("Routing table %s still exists", name)
		}

		if !client.IsNotFoundError(err) {
			return err
		}
	}

	return nil
}
