config Provider
==================

- Website: https://registry.terraform.io
- Documentation: https://registry.terraform.io/providers/alabuel/config/latest/docs

Introduction
------------

This is a Terraform provider for the following:
- Read Excel file as a data source
- Read INI configuration file as a data source
- Get response from API servers as a data source

Requirements
------------

* [Terraform 0.13 or above](https://www.terraform.io/downloads.html)
* [Go Language 1.16.4 or above](https://golang.org/dl)


Using the Provider
------------

See the [config documentation](https://registry.terraform.io/providers/alabuel/config/latest/docs)

See config data source examples [here](examples/README.md)

## Execution
These are the Terraform commands that can be used for the config plugin:
* `terraform init` - The init command is used to initialize a working directory containing Terraform configuration files.
* `terraform plan` - Plan command shows plan for resources like how many resources will be provisioned and how many will be destroyed.
* `terraform apply` - apply is responsible to execute actual calls to provision resources.
* `terraform refresh` - By using the refresh command you can check the status of the request.
* `terraform show` - show will set a console output for resource configuration and request status.
* `terraform destroy` - destroy command will destroy all the  resources present in terraform configuration file.

Navigate to the location where `main.tf` and binary are placed and use the above commands as needed.

Upgrading the provider
----------------------

The config provider doesn't upgrade automatically once you've started using it. After a new release you can run 

```bash
terraform init -upgrade
```

to upgrade to the latest stable version of the config provider. See the [Terraform website](https://www.terraform.io/docs/configuration/providers.html#provider-versions)
for more information on provider upgrades, and how to set version constraints on your provider.


### A sample main.tf is as follows

```hcl
terraform {
  required_providers {
    config = {
      source = "alabuel/config"
      version = "0.2.4"
    }
  }
}

provider "config" {}

# reading Excel worksheet
# ------------------------
data "config_workbook" "excel" {
  excel = "filename.xlsx"
  worksheet = "Sheet1"
}

# getting api response
# -----------------------
data "config_restapi_get "apidata" {
  uri = "http://localhost:3000/posts"
}

# reading ini configuration file
# -------------------------------
data "config_ini" "cfg" {
  ini = file("configuration.ini")
}
```
