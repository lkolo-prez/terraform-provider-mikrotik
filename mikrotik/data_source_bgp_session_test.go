package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/ddelnano/terraform-provider-mikrotik/mikrotik/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMikrotikBgpSessionDataSource_basic(t *testing.T) {
	client.SkipIfRouterOSV6OrEarlier(t, sysResources)
	name := acctest.RandomWithPrefix("tf-acc-sess")
	instanceName := acctest.RandomWithPrefix("tf-acc-inst")
	routerId := internal.GetNewIpAddr()
	remoteAddr := internal.GetNewIpAddr()
	localAs := acctest.RandIntRange(64512, 65535)
	remoteAs := acctest.RandIntRange(64512, 65535)

	dataSourceName := "data.mikrotik_bgp_session.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpSessionDataSourceBasic(name, instanceName, localAs, remoteAs, routerId, remoteAddr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttrSet(dataSourceName, "established"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "remote_address"),
				),
			},
		},
	})
}

func testAccBgpSessionDataSourceBasic(name, instanceName string, localAs, remoteAs int, routerId, remoteAddr string) string {
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

data "mikrotik_bgp_session" "test" {
  name = mikrotik_bgp_connection.test.name
}
`, instanceName, localAs, routerId, name, localAs, remoteAddr, remoteAs)
}
