terraform {
  required_providers {
    config = {
      version = "0.2.2"
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
}

data "config_workbook" "vm" {
  csv = file("files/deployVM.csv")
  configuration_item = "virtual_machine"
}

data "config_workbook" "excel" {
  excel = "files/data.xlsx"
  worksheet = "Config"
  headers = ["config_name","data2","data3"]
}

data "config_workbook" "vexcel" {
  excel = "files/data.xlsx"
  worksheet = "Vert"
  orientation = "vertical"
  configuration_item = "my_vertical"
}

data "config_workbook" "vert" {
  excel = "files/event.xlsx"
  password = "1234"
  worksheet = "event_target"
  orientation = "vertical"
  configuration_item = "event_target"
}

data "config_workbook" "lookup" {
  excel = "files/event.xlsx"
  password = "1234"
  worksheet = "cloudwatch_event_rule"
  configuration_item = "cloudwatch_event_rule"

  lookup {
    excel = "files/event.xlsx"
    password = "1234"
    worksheet = "event_target"
    column = "command"
    key_column = "name"
    value_column = "script"
  }

  lookup {
    ini = file("files/event.ini")
    section = "event"
    column = "dependents"
    key_column = "name"
    value_column = "script"
  }
}

data "config_ini" "ini" {
  ini = file("files/event.ini")
}

# output "horiz" {
#   value = jsondecode(data.config_workbook.excel.json)
# }

# output "vert" {
#   value = jsondecode(data.config_workbook.vexcel.json)
# }

output "lookup" {
  value = jsondecode(data.config_workbook.lookup.json)
}

# output "events" {
#   value = jsondecode(data.config_workbook.vert.json)
# }

# output "csv" {
#   value = jsondecode(data.config_workbook.csv.json)
# }

# output "vm" {
#   value = jsondecode(data.config_workbook.vm.json)
# }

output "ini" {
  value = jsondecode(data.config_ini.ini.json)
}