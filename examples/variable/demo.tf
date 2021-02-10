terraform {
  required_providers {
    airflow = {
      source  = "houqp/airflow"
      version = "0.2.1"
    }
  }
}

provider "airflow" {
  base_endpoint = "http://localhost:28080/"
  username      = "test"
  password      = "test"
}

resource "airflow_variable" "foo" {
  key   = "foo"
  value = "bar"
}

resource "airflow_variable" "hello" {
  key   = "hello"
  value = "world"
}
