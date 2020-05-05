provider "airflow" {
  variables_output_path = "./variables.json"
}

resource "airflow_variable" "foo" {
  key   = "foo"
  value = "bar"
}
