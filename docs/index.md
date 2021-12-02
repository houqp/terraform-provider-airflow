---
page_title: "Provider: Airflow"
description: Manage Airflow with Terraform.
---

# Airflow Provider

This is a terraform provider plugin for managing [Apache Airflow](https://airflow.apache.org/).

## Example Provider Configuration

```terraform
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
```