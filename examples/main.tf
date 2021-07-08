terraform {
  required_providers {
    config = {
      version = "0.1.4"
      source = "alabuel/config"
    }
  }
}

provider "config" {}

data "config_workbook" "csv" {
  csv = file("files/test.csv")
  schema = file("files/config.yaml")
  filter {
    name = "name"
    values = ["item_name1"]
  }
  configuration_item = "configuratio_item"
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
