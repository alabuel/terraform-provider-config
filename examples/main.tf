terraform {
  required_providers {
    config = {
      version = "0.0.1"
      source = "aa/common/config"
    }
  }
}

provider "config" {}

data "config_workbook" "csv" {
  csv = file("files/test.csv")
  schema = file("files/config.yaml")
}

data "config_workbook" "excel" {
  excel = "files/data.xlsx"
  worksheet = "Config"
}

output "csv" {
  value = jsondecode(data.config_workbook.csv.json)
}

output "excel" {
  value = jsondecode(data.config_workbook.excel.json)
}