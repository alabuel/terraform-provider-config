terraform {
  required_providers {
    config = {
      version = "0.2.2"
      source = "aa/test/config"
    }
  }
}

provider "config" {}

data "config_restapi_get" "test" {
    uri = "http://localhost:3000/posts"
    header {
      key = "Content-Type"
      value = "application/json"
    }
    param {
      key = "id"
      value = "1"
    }
    param {
      key = "author"
      value = "typicode"
    }
}

output "response" {
    value = jsondecode(data.config_restapi_get.test.response)
}
