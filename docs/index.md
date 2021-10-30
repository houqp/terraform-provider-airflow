---
layout: "airflow"
page_title: "Provider: Airflow"
sidebar_current: "docs-airflow-index"
description: |-
  The Airflow provider is used to interact with Airflow.
---

# Airflow Provider

The Airflow provider is used to interact with the Airflow. The
provider needs to be configured with the proper credentials before it can be
used.

Use the navigation to the left to read about the available data sources.

## Example Usage

```hcl
provider "airflow" {
  base_endpoint  = "airflow.net"
  oauth2_token   = "token"
}

resource "airflow_variable" "default" {
  key = "foo"
  value = "bar"
}
```

## Authentication

An OAUTH2 token must be passed to the provider block.

## Argument Reference

- `base_endpoint` - (Required) The Airflow API endpoint.
- `oauth2_token` - (Required) An OAUTH2 identity token used to authenticate against an Airflow server.
