package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAirflowConnection_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resourceName := "airflow_connection.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAirflowConnectionCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAirflowConnectionConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "connection_id", rName),
					resource.TestCheckResourceAttr(resourceName, "conn_type", "http"),
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

func TestAccAirflowConnection_full(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rNameUpdated := acctest.RandomWithPrefix("tf-acc-test")

	resourceName := "airflow_connection.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAirflowConnectionCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAirflowConnectionConfigFull(rName, rName, 443),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "connection_id", rName),
					resource.TestCheckResourceAttr(resourceName, "conn_type", "http"),
					resource.TestCheckResourceAttr(resourceName, "host", rName),
					resource.TestCheckResourceAttr(resourceName, "login", rName),
					resource.TestCheckResourceAttr(resourceName, "schema", rName),
					resource.TestCheckResourceAttr(resourceName, "port", "443"),
					resource.TestCheckResourceAttr(resourceName, "extra", rName),
					resource.TestCheckResourceAttr(resourceName, "password", rName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccAirflowConnectionConfigFull(rName, rNameUpdated, 80),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "connection_id", rName),
					resource.TestCheckResourceAttr(resourceName, "conn_type", "http"),
					resource.TestCheckResourceAttr(resourceName, "host", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "login", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "schema", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "port", "80"),
					resource.TestCheckResourceAttr(resourceName, "extra", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "password", rName),
				),
			},
		},
	})
}

func testAccCheckAirflowConnectionCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "airflow_connection" {
			continue
		}

		conn, res, err := client.ApiClient.ConnectionApi.GetConnection(client.AuthContext, rs.Primary.ID).Execute()
		if err == nil {
			if *conn.ConnectionId == rs.Primary.ID {
				return fmt.Errorf("Airflow Connection (%s) still exists.", rs.Primary.ID)
			}
		}

		if res != nil && res.StatusCode == 404 {
			continue
		}
	}

	return nil
}

func testAccAirflowConnectionConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "airflow_connection" "test" {
  connection_id = %[1]q
  conn_type     = "http"
}
`, rName)
}

func testAccAirflowConnectionConfigFull(rName, rName2 string, port int) string {
	return fmt.Sprintf(`
resource "airflow_connection" "test" {
  connection_id = %[1]q
  conn_type     = "http"
  host          = %[2]q
  login         = %[2]q
  schema        = %[2]q
  port          = %[3]d
  password      = %[2]q
  extra         = %[2]q
}
`, rName, rName2, port)
}
