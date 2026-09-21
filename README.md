# terraform-provider-servicenow

Terraform provider for querying ServiceNow records.

## Supported Data Sources

- `servicenow_service_offering`
- `servicenow_user_group`

## Provider Configuration

```hcl
provider "servicenow" {
  base_url = "https://example.service-now.com"
}
```

Credentials can be configured with environment variables:

```bash
export SERVICENOW_USERNAME="..."
export SERVICENOW_PASSWORD="..."
```

The ServiceNow instance URL can also be configured with:

```bash
export SERVICENOW_BASE_URL="https://example.service-now.com"
```

## Example

```hcl
data "servicenow_service_offering" "example" {
  name = "Example Service Offering"
  view = "automate"
}

data "servicenow_user_group" "example" {
  name = "Example User Group"
}

output "service_offering_sys_id" {
  value = data.servicenow_service_offering.example.sys_id
}

output "user_group_sys_id" {
  value = data.servicenow_user_group.example.sys_id
}
```

More examples are available in the `examples/` directory.

## Local Development

Build the provider:

```bash
go build -o terraform-provider-servicenow
```

Create a local Terraform CLI config outside the repository, for example:

```hcl
provider_installation {
  dev_overrides {
    "achachw/servicenow" = "/absolute/path/to/terraform-provider-servicenow"
  }

  direct {}
}
```

Then run Terraform with:

```bash
TF_CLI_CONFIG_FILE=/absolute/path/to/dev.tfrc terraform plan
```

## Testing

```bash
go test ./...
```
