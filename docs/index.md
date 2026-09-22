---
page_title: "ServiceNow Provider"
subcategory: ""
description: |-
  The ServiceNow provider queries records from ServiceNow table APIs.
---

# ServiceNow Provider

The ServiceNow provider queries records from ServiceNow table APIs.

## Example Usage

```terraform
terraform {
  required_providers {
    servicenow = {
      source  = "achachw/servicenow"
      version = "~> 0.1"
    }
  }
}

provider "servicenow" {
  base_url = "https://example.service-now.com"
}
```

Credentials can be configured with environment variables:

```shell
export SERVICENOW_USERNAME="..."
export SERVICENOW_PASSWORD="..."
```

The instance URL can also be configured with:

```shell
export SERVICENOW_BASE_URL="https://example.service-now.com"
```

## Schema

### Optional

- `base_url` (String) ServiceNow instance base URL. Can also be set with `SERVICENOW_BASE_URL`.
- `username` (String, Sensitive) ServiceNow username. Can also be set with `SERVICENOW_USERNAME`.
- `password` (String, Sensitive) ServiceNow password. Can also be set with `SERVICENOW_PASSWORD`.

## Supported Data Sources

- `servicenow_service_offering`
- `servicenow_user_group`
