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
  base_endpoint = "airflow.net"
  oauth2_token  = "token"
}

resource "airflow_variable" "default" {
  key   = "foo"
  value = "bar"
}
```

## Authentication

An OAUTH2 token must be passed to the provider block.

### Google Composer Example

```terraform
data "http" "client_id" {
  url = "composer-url"
}

resource "google_service_account" "example" {
  account_id = "example"
}

data "google_service_account_access_token" "impersonated" {
  target_service_account = google_service_account.example.email
  delegates              = []
  scopes                 = ["userinfo-email", "cloud-platform"]
  lifetime               = "300s"
}

provider "google" {
  alias        = "impersonated"
  access_token = data.google_service_account_access_token.impersonated.access_token
}

data "google_service_account_id_token" "oidc" {
  provider               = google.impersonated
  target_service_account = google_service_account.example.email
  delegates              = []
  include_email          = true
  target_audience        = regex("[A-Za-z0-9-]*\\.apps\\.googleusercontent\\.com", data.http.client_id.body)
}

provider "airflow" {
  base_endpoint = data.http.client_id.url
  oauth2_token  = data.google_service_account_id_token.oidc.id_token
}
```

## Argument Reference

- `base_endpoint` - (Required) The Airflow API endpoint.
- `oauth2_token` - (Required) An OAUTH2 identity token used to authenticate against an Airflow server.
