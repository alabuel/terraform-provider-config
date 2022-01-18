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

data "config_restapi_get" "post" {
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