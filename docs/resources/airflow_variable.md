---
layout: "airflow"
page_title: "Airflow: airflow_variable"
sidebar_current: "docs-airflow-resource-variable"
description: |-
  Provides an Airflow variable
---

# airflow_variable

Provides an Airflow variable.

## Example Usage

```hcl
resource airflow_variable "example" {
  key   = "example"
  value = "example"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) The variable key.
* `value` - (Required) The variable value.

## Attributes Reference

This resource exports the following attributes:

* `id` - The variable key.

## Import

Variables can be imported using the variable key.

```terraform
terraform import airflow_variable.default example
```
