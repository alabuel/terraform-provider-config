terraform {
  required_providers {
    config = {
      version = "0.2.1"
      source = "aa/test/config"
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
  orientation = "vertical"
  configuration_item = "my_vertical"
}

data "config_workbook" "lkexcel" {
  excel = "files/event.xlsx"
  worksheet = "cloudwatch_event_rule"
  configuration_item = "cloudwatch_event_rule"

  lookup {
    column = "command"
    worksheet = "event_target"
    key_column = "name"
    value_column = "script"
  }

  lookup {
    column = "dependents"
    worksheet = "event_target"
    key_column = "name"
    value_column = "script"
  }
}

output "horiz" {
  value = jsondecode(data.config_workbook.excel.json)
}

output "vert" {
  value = jsondecode(data.config_workbook.vexcel.json)
}

output "lookup" {
  value = jsondecode(data.config_workbook.lkexcel.json)
}