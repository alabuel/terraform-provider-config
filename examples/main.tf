terraform {
  required_providers {
    config = {
      version = "0.1.6"
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

data "config_workbook" "vexcel" {
  excel = "files/data.xlsx"
  worksheet = "Vert"
  type = "vertical"
  configuration_item = "my_vertical"
}

output "horiz" {
  value = jsondecode(data.config_workbook.excel.json)
}

output "vert" {
  value = jsondecode(data.config_workbook.vexcel.json)
}
