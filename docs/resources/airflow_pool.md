---
layout: "airflow"
page_title: "Airflow: airflow_pool"
sidebar_current: "docs-airflow-resource-pool"
description: |-
  Provides an Airflow pool
---

# airflow_pool

Provides an Airflow pool.

## Example Usage

```hcl
resource airflow_pool "example" {
  name  = "example"
  slots = 2
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of pool.
* `slots` - (Required) The maximum number of slots that can be assigned to tasks. One job may occupy one or more slots. 

## Attributes Reference

This resource exports the following attributes:

* `id` - The pool name.
* `occupied_slots` - The number of slots used by running/queued tasks at the moment.
* `used_slots` - The number of slots used by running tasks at the moment.
* `queued_slots` - The number of slots used by queued tasks at the moment.
* `open_slots` - The number of free slots at the moment.

## Import

Content can be imported using the pool name.

```terraform
terraform import airflow_pool.default example
```
