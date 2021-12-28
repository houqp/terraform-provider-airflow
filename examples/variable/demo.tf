terraform {
  required_providers {
    airflow = {
      source  = "drfaust92/airflow"
      version = "0.2.9"
    }
  }
}

provider "airflow" {
  base_endpoint = "http://localhost:28080/"
  oauth2_token  = "some-token"
}

resource "airflow_variable" "foo" {
  key   = "foo"
  value = "bar"
}

resource "airflow_variable" "hello" {
  key   = "hello"
  value = "world"
}
