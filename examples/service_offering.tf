data "servicenow_service_offering" "example" {
  name = "Example Service Offering"
  view = "automate"
}

data "servicenow_user_group" "example" {
  name = "Example User Group"
}

