package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_User(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	email := "test2@example.com"
	username := "test_user"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: userConfig(name1, name1, email, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("airflow_user.test", "first_name", name1),
					resource.TestCheckResourceAttr("airflow_user.test", "last_name", name1),
					resource.TestCheckResourceAttr("airflow_user.test", "email", email),
					resource.TestCheckResourceAttr("airflow_user.test", "username", username),
				),
			},
			{
				Config: userConfig(name2, name2, email, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("airflow_user.test", "first_name", name2),
					resource.TestCheckResourceAttr("airflow_user.test", "last_name", name2),
					resource.TestCheckResourceAttr("airflow_user.test", "email", email),
					resource.TestCheckResourceAttr("airflow_user.test", "username", username),
				),
			},
		},
	})
}

func userConfig(firstName, lastName, email, username string) string {
	return fmt.Sprintf(`
	resource "airflow_user" "test" {
		first_name = "%s"
		last_name = "%s"
		email = "%s"
		username = "%s"
		roles = ["Public"]
	}
	`, firstName, lastName, email, username)
}
