package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAirflowUser_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rNameUpdated := acctest.RandomWithPrefix("tf-acc-test")

	resourceName := "airflow_user.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAirflowUserCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAirflowUserConfigBasic(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "email", rName),
					resource.TestCheckResourceAttr(resourceName, "first_name", rName),
					resource.TestCheckResourceAttr(resourceName, "last_name", rName),
					resource.TestCheckResourceAttr(resourceName, "username", rName),
					resource.TestCheckResourceAttr(resourceName, "password", rName),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "roles.*", "airflow_role.test", "name"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccAirflowUserConfigBasic(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "email", rName),
					resource.TestCheckResourceAttr(resourceName, "first_name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "last_name", rName),
					resource.TestCheckResourceAttr(resourceName, "username", rName),
					resource.TestCheckResourceAttr(resourceName, "password", rName),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "roles.*", "airflow_role.test", "name"),
				),
			},
		},
	})
}

func testAccCheckAirflowUserCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "airflow_user" {
			continue
		}

		user, res, err := client.ApiClient.UserApi.GetUser(client.AuthContext, rs.Primary.ID).Execute()
		if err == nil {
			if *user.Username == rs.Primary.ID {
				return fmt.Errorf("Airflow User (%s) still exists.", rs.Primary.ID)
			}
		}

		if res != nil && res.StatusCode == 404 {
			continue
		}
	}

	return nil
}

func testAccAirflowUserConfigBasic(rName, fName string) string {
	return fmt.Sprintf(`
resource "airflow_role" "test" {
  name   = %[1]q

  action {
    action   = "can_read"
	resource = "Audit Logs"
  } 
}

resource "airflow_user" "test" {
  email      = %[1]q
  first_name = %[2]q
  last_name  = %[1]q
  username   = %[1]q
  password   = %[1]q
  roles      = [airflow_role.test.name]
}
`, rName, fName)
}
