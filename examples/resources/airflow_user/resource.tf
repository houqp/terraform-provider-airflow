resource "airflow_user" "example_user" {
  username   = "test_user"
  email      = "test@example.com"
  first_name = "test"
  last_name  = "user"
}
