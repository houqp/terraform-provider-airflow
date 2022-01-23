---
layout: "airflow"
page_title: "Airflow: airflow_dag"
sidebar_current: "docs-airflow-resource-dag"
description: |-
  Provides an Airflow dag
---

# airflow_dag

Provides an Airflow DAG.

> Note this resource adpots an existing DAG and does not create a one, Also on delete the resource is only deleted from state and not acutally deleted. This behavior my change in the future

## Example Usage

```hcl
resource airflow_dag "example" {
  dag_id    = "example"
  is_paused = false
}
```

## Argument Reference

The following arguments are supported:

* `dag_id` - (Required) The ID of the DAG.
* `is_paused` - (Required) Whether the DAG is paused.

## Attributes Reference

This resource exports the following attributes:

* `id` - The ID of the DAG.
* `is_active` - Whether the DAG is currently seen by the scheduler(s).
* `is_subdag` - Whether the DAG is SubDAG.
* `description` - User-provided DAG description, which can consist of several sentences or paragraphs that describe DAG contents.
* `fileloc` - The absolute path to the file.
* `file_token` - The key containing the encrypted path to the file. Encryption and decryption take place only on the server. This prevents the client from reading an non-DAG file.
* `root_dag_id` - If the DAG is SubDAG then it is the top level DAG identifier. Otherwise, null.

## Import

DAGs can be imported using the DAG Id.

```terraform
terraform import airflow_dag.default example
```
