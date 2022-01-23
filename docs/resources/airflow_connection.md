---
layout: "airflow"
page_title: "Airflow: airflow_connection"
sidebar_current: "docs-airflow-resource-connection"
description: |-
  Provides an Airflow connection
---

# airflow_connection

Provides an Airflow connection.

## Example Usage

```hcl
resource airflow_connection "example" {
  connection_id = "example"
  conn_type     = "example"
}
```

## Argument Reference

The following arguments are supported:

* `connection_id` - (Required) The connection ID.
* `conn_type` - (Required) The connection type.
* `host` - (Optional) The host of the connection.
* `login` - (Optional) The login of the connection.
* `schema` - (Optional) The schema of the connection.
* `port` - (Optional) The port of the connection.
* `password` - (Optional) The paasword of the connection.
* `extra` - (Optional) Other values that cannot be put into another field, e.g. RSA keys.

## Attributes Reference

This resource exports the following attributes:

* `id` - The connection id.

## Import

Connections can be imported using the connection key.

```terraform
terraform import airflow_connection.default example
```
