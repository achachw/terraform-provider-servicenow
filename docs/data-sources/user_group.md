---
page_title: "servicenow_user_group Data Source"
subcategory: "Identity"
description: |-
  Fetches one ServiceNow user group by name.
---

# servicenow_user_group

Fetches one ServiceNow user group by name from the `sys_user_group` table.

The provider expects the name lookup to return exactly one record. If the exact name is not found, the provider searches for similar names and includes them in the diagnostic message.

## Example Usage

```terraform
data "servicenow_user_group" "example" {
  name = "Example User Group"
}

output "user_group_sys_id" {
  value = data.servicenow_user_group.example.sys_id
}
```

## Schema

### Required

- `name` (String) ServiceNow user group name.

### Optional

- `view` (String) Optional ServiceNow view used for the table request.

### Read-Only

- `id` (String) Terraform datasource ID. Same value as `sys_id`.
- `sys_id` (String) ServiceNow `sys_id`.
- `description` (String) ServiceNow user group description.
- `email` (String) ServiceNow user group email.
- `manager` (String) ServiceNow user group manager display value or `sys_id`.
- `parent` (String) ServiceNow parent group display value or `sys_id`.
- `active` (String) Whether the ServiceNow user group is active.
- `raw_json` (String) Raw ServiceNow record as JSON.
