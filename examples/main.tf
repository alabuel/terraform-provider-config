terraform {
  required_providers {
    config = {
      version = "0.0.1"
      source = "cloud/common/config"
    }
  }
}

provider "config" {}

data "config_workbook" "this" {
  csv = file("files/test.csv")
  schema = file("files/config.yaml")
}

output "data" {
  value = jsondecode(data.config_workbook.this.json)
}