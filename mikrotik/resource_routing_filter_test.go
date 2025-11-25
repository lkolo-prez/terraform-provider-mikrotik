package mikrotik

import (
	"fmt"
	"testing"

	"github.com/ddelnano/terraform-provider-mikrotik/client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRoutingFilterRule_basic(t *testing.T) {
	resourceName := "mikrotik_routing_filter_rule.deny_default"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRoutingFilterRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingFilterRuleConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRoutingFilterRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "chain", "bgp-in-test"),
					resource.TestCheckResourceAttr(resourceName, "rule", "if (dst == 0.0.0.0/0) { reject }"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Deny default route from BGP"),
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

func TestAccRoutingFilterRule_bgpCommunity(t *testing.T) {
	resourceName := "mikrotik_routing_filter_rule.accept_customer"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRoutingFilterRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingFilterRuleConfigBGPCommunity(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRoutingFilterRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "chain", "bgp-in-test"),
					resource.TestCheckResourceAttr(resourceName, "rule", "if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
				),
			},
		},
	})
}

func TestAccRoutingFilterRule_setLocalPref(t *testing.T) {
	resourceName := "mikrotik_routing_filter_rule.set_localpref"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRoutingFilterRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingFilterRuleConfigSetLocalPref(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRoutingFilterRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "chain", "bgp-in-test"),
					resource.TestCheckResourceAttr(resourceName, "rule", "if (bgp-communities includes 65001:200) { set bgp-local-pref 200; accept }"),
				),
			},
		},
	})
}

func TestAccRoutingFilterRule_update(t *testing.T) {
	resourceName := "mikrotik_routing_filter_rule.test_update"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRoutingFilterRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingFilterRuleConfigUpdateBefore(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRoutingFilterRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule", "if (dst == 192.168.0.0/16) { reject }"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Block private networks"),
				),
			},
			{
				Config: testAccRoutingFilterRuleConfigUpdateAfter(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRoutingFilterRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule", "if (dst in 192.168.0.0/16) { reject }"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Block private networks - disabled for testing"),
				),
			},
		},
	})
}

func TestAccRoutingFilterChain_basic(t *testing.T) {
	resourceName := "mikrotik_routing_filter_chain.test_chain"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRoutingFilterChainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingFilterChainConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRoutingFilterChainExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "test-filter-chain"),
					resource.TestCheckResourceAttr(resourceName, "dynamic", "false"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Test filter chain"),
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

func TestAccRoutingFilterChain_withRules(t *testing.T) {
	chainResource := "mikrotik_routing_filter_chain.bgp_filtering"
	rule1Resource := "mikrotik_routing_filter_rule.deny_default"
	rule2Resource := "mikrotik_routing_filter_rule.accept_customer"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRoutingFilterChainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingFilterChainConfigWithRules(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check chain
					testAccCheckRoutingFilterChainExists(chainResource),
					resource.TestCheckResourceAttr(chainResource, "name", "bgp-in-complete"),
					resource.TestCheckResourceAttr(chainResource, "dynamic", "true"),
					
					// Check rules
					testAccCheckRoutingFilterRuleExists(rule1Resource),
					resource.TestCheckResourceAttr(rule1Resource, "chain", "bgp-in-complete"),
					resource.TestCheckResourceAttr(rule1Resource, "rule", "if (dst == 0.0.0.0/0) { reject }"),
					
					testAccCheckRoutingFilterRuleExists(rule2Resource),
					resource.TestCheckResourceAttr(rule2Resource, "chain", "bgp-in-complete"),
					resource.TestCheckResourceAttr(rule2Resource, "rule", "if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }"),
				),
			},
		},
	})
}

// Helper functions

func testAccCheckRoutingFilterRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		c := testAccProvider.Meta().(*client.Mikrotik)
		_, err := c.FindRoutingFilterRuleById(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Routing filter rule does not exist: %s", err)
		}

		return nil
	}
}

func testAccCheckRoutingFilterRuleDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client.Mikrotik)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_routing_filter_rule" {
			continue
		}

		_, err := c.FindRoutingFilterRuleById(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Routing filter rule still exists")
		}
	}

	return nil
}

func testAccCheckRoutingFilterChainExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No name is set")
		}

		c := testAccProvider.Meta().(*client.Mikrotik)
		_, err := c.FindRoutingFilterChain(rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Routing filter chain does not exist: %s", err)
		}

		return nil
	}
}

func testAccCheckRoutingFilterChainDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client.Mikrotik)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mikrotik_routing_filter_chain" {
			continue
		}

		_, err := c.FindRoutingFilterChain(rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Routing filter chain still exists")
		}
	}

	return nil
}

// Test configurations

func testAccRoutingFilterRuleConfigBasic() string {
	return `
resource "mikrotik_routing_filter_rule" "deny_default" {
  chain    = "bgp-in-test"
  rule     = "if (dst == 0.0.0.0/0) { reject }"
  disabled = false
  comment  = "Deny default route from BGP"
}
`
}

func testAccRoutingFilterRuleConfigBGPCommunity() string {
	return `
resource "mikrotik_routing_filter_rule" "accept_customer" {
  chain    = "bgp-in-test"
  rule     = "if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }"
  disabled = false
  comment  = "Accept customer routes with community tag"
}
`
}

func testAccRoutingFilterRuleConfigSetLocalPref() string {
	return `
resource "mikrotik_routing_filter_rule" "set_localpref" {
  chain    = "bgp-in-test"
  rule     = "if (bgp-communities includes 65001:200) { set bgp-local-pref 200; accept }"
  disabled = false
  comment  = "Set local-pref for preferred routes"
}
`
}

func testAccRoutingFilterRuleConfigUpdateBefore() string {
	return `
resource "mikrotik_routing_filter_rule" "test_update" {
  chain    = "bgp-in-test"
  rule     = "if (dst == 192.168.0.0/16) { reject }"
  disabled = false
  comment  = "Block private networks"
}
`
}

func testAccRoutingFilterRuleConfigUpdateAfter() string {
	return `
resource "mikrotik_routing_filter_rule" "test_update" {
  chain    = "bgp-in-test"
  rule     = "if (dst in 192.168.0.0/16) { reject }"
  disabled = true
  comment  = "Block private networks - disabled for testing"
}
`
}

func testAccRoutingFilterChainConfigBasic() string {
	return `
resource "mikrotik_routing_filter_chain" "test_chain" {
  name     = "test-filter-chain"
  dynamic  = false
  disabled = false
  comment  = "Test filter chain"
}
`
}

func testAccRoutingFilterChainConfigWithRules() string {
	return `
resource "mikrotik_routing_filter_chain" "bgp_filtering" {
  name     = "bgp-in-complete"
  dynamic  = true
  disabled = false
  comment  = "Complete BGP input filtering chain"
}

resource "mikrotik_routing_filter_rule" "deny_default" {
  chain    = mikrotik_routing_filter_chain.bgp_filtering.name
  rule     = "if (dst == 0.0.0.0/0) { reject }"
  disabled = false
  comment  = "Deny default route"
  
  depends_on = [mikrotik_routing_filter_chain.bgp_filtering]
}

resource "mikrotik_routing_filter_rule" "accept_customer" {
  chain    = mikrotik_routing_filter_chain.bgp_filtering.name
  rule     = "if (dst in 10.0.0.0/8 && bgp-communities includes 65001:100) { accept }"
  disabled = false
  comment  = "Accept customer routes"
  
  depends_on = [mikrotik_routing_filter_chain.bgp_filtering]
}
`
}
