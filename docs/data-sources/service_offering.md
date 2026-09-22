---
page_title: "servicenow_service_offering Data Source"
subcategory: "Service Catalog"
description: |-
  Fetches one ServiceNow service offering by name.
---

# servicenow_service_offering

Fetches one ServiceNow service offering by name from the `service_offering` table.

The provider expects the name lookup to return exactly one record. If the exact name is not found, the provider searches for similar names and includes them in the diagnostic message.

## Example Usage

```terraform
data "servicenow_service_offering" "example" {
  name = "Example Service Offering"
  view = "automate"
}

output "service_offering_sys_id" {
  value = data.servicenow_service_offering.example.sys_id
}
```

## Schema

### Required

- `name` (String) Service offering name.

### Optional

- `view` (String) ServiceNow view used for the table request. Defaults to `automate`.

### Read-Only

- `id` (String) Terraform datasource ID. Same value as `sys_id`.
- `sys_id` (String) ServiceNow `sys_id`.
- `short_description` (String) Service offering short description.
- `operational_status` (String) Service offering operational status.
- `owned_by` (String) Service offering owner display value or `sys_id`.
- `managed_by` (String) Service offering manager display value or `sys_id`.
- `raw_json` (String) Raw ServiceNow record as JSON.
