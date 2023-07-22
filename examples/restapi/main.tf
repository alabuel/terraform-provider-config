terraform {
  required_providers {
    config = {
      version = "0.2.8"
      source = "alabuel/config"
    }
  }
}

provider "config" {}

data "config_rest" "test" {
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

data "config_rest" "post" {
  uri = "http://localhost:3000/posts"
  header {
    key = "Content-Type"
    value = "application/json"
  }
  method = "POST"
  payload = <<-EOT
  {
    "title":"test-post",
    "author":"me"
  }
  EOT
}

output "response" {
    value = jsondecode(data.config_restapi_get.test.response)
}

output "post" {
  value = jsondecode(data.config_restapi_get.post.response)
}