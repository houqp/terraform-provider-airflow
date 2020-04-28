provider "airflow" {
}

resource "airflow_variable" "foo" {
  key   = "foo"
  value = "bar"
}
