package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAirflowRole_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resourceName := "airflow_role.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAirflowRoleCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAirflowRoleConfigBasic(rName, "can_read", "Audit Logs"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "action.*", map[string]string{
						"action":   "can_read",
						"resource": "Audit Logs",
					}),
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

func testAccCheckAirflowRoleCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "airflow_role" {
			continue
		}

		variable, res, err := client.ApiClient.RoleApi.GetRole(client.AuthContext, rs.Primary.ID).Execute()
		if err == nil {
			if *variable.Name == rs.Primary.ID {
				return fmt.Errorf("Airflow Role (%s) still exists.", rs.Primary.ID)
			}
		}

		if res != nil && res.StatusCode == 404 {
			continue
		}
	}

	return nil
}

func testAccAirflowRoleConfigBasic(rName, action, resource string) string {
	return fmt.Sprintf(`
resource "airflow_role" "test" {
  name   = %[1]q

  action {
    action   = %[2]q
	resource = %[3]q
  } 
}
`, rName, action, resource)
}
