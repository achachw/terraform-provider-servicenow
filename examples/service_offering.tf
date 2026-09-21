terraform {
  required_providers {
    servicenow = {
      source = "achachw/servicenow"
    }
  }
}

provider "servicenow" {
  base_url = "https://example.service-now.com"
  # username can be set with SERVICENOW_USERNAME
  # password can be set with SERVICENOW_PASSWORD
}

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
