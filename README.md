# Terraform Configuration Workbook Provider

This is the repository for the Terraform Configuration Workbook Provider which you can use to parse a CSV file or an Excel file into a list of maps.

You can also provide a schema to set the key names of each map for your configurations.

# Using the Provider

## Configuring the Provider

```hcl
terraform {
  required_providers {
    config = {
      source = "alabuel/config/workbook"
      version = "0.1.0"
    }
  }
}

provider "config" {}
```

### Example - Using a CSV

```hcl
data "config_workbook "csv" {
  csv = <<-EOT
  configuration_item,attr1,attr2,attr3
  vpc,my_vpc,1,"10.0.0.0/16"
  EOT
}

# reading from a csv file
data "config_workbook" "csv_file" {
  csv = file("filename.csv")
}
```

### Example - Using an Excel file

```hcl
data "config_workbook" "excel" {
  excel = "filename.xlsx"
  worksheet = "Sheet1"
}
```

### Example - Using a CSV with a schema

```hcl
data "config_workbook" "csv_using_yaml" {
  csv = file("filename.csv")
  schema = file("schema.yaml")
}

data "config_workbook" "csv_using_json" {
  csv = file("filename.csv")
  schema = file("schema.json")
}
```

### Example - Using an Excel with a schema
```hcl
data "config_workbook" "excel_using_yaml" {
  excel = "filename.xlsx"
  worksheet = "Sheet1"
  schema = file("schema.yaml")
}

data "config_workbook" "excel_using_json" {
  excel = "filename.xlsx"
  worksheet = "Sheet1"
  schema = file("schema.json")
}
```

### Schema format - Example 1
```yaml
# you can set the attribute types
schema_config:
  vpc:
    attr1:
      name: name
      type: string
    attr2:
      name: create
      type: bool
    attr3:
      name: cidr_block
      type: string

# this format assumes that all attributes are of "string" types
schema_config:
  vpc:
    attr1: name
    attr2: create
    attr3: cidr_block

```

### Schema format - Example 3
```json
// you can set the attribute types
{
  "schema_config": {
    "vpc": {
      "attr1": {
        "name": "name",
        "type": "string"
      },
      "attr2": {
        "name": "create",
        "type": "bool"
      },
      "attr3": {
        "name": "cidr_block",
        "type": "string"
      }
    }
  }
} 

// this format assumes that all attributes are of "string" types
{
  "schema_config": {
    "vpc": {
      "attr1": "name",
      "attr2": "create",
      "attr3": "cidr_block"
    }
  }
}
```

## Valid attribute types
- string
- number/numeric
- bool/boolean
- map
- list

## Attribute naming convention
1. Should have a column name of "`configuration_item`".  This will identify the item you need to configure
2. Attributes starting with "`attr`" will be substituted with the correct attribute name using the provided schema.
3. You can preset the type of the attribute using prefixes.
    - `s_` or `string_`
    - `n_` or `num_` or `number_` or `numeric_`
    - `b_` or `bool_` or `boolean`
    - `m_` or `map_`
    - `l_` or `list_`
    - `t_` or `tag_`
    Attributes without prefixes will be treated as string.  Boolean values are (1,yes,true = True; 0,no,false = False)

## Example 1
|configuration_item|attr1|attr2|attr3|
|------------------|-----|-----|-----|
|vpc|my_vpc|1|"10.0.0.0/16"|

## Example 2
|configuration_item|name|b_create|cidr_block|
|------------------|-----|-----|-----|
|vpc|my_vpc|1|"10.0.0.0/16"|

### Example 1 and 2 results:
```json
{
  "vpc": [
    {
      "name": "my_vpc",
      "create": true,
      "cidr_block": "10.0.0.0/16"
    }
  ]
}
```

## Example 3
|configuration_item|attr1|attr2|attr3|t_environment|t_alias|t_purpose|
|------------------|-----|-----|-----|-------------|-------|---------|
|vpc|my_vpc|1|"10.0.0.0/16"|dev|vpc1|"application vpc"|

### Example 3 results will be:
```json
{
  "vpc": [
    {
      "name": "my_vpc",
      "create": true,
      "cidr_block": "10.0.0.0/16",
      "tags": {
        "Environment": "dev",
        "Alias": "vpc1",
        "Purpose": "application vpc"
      }
    }
  ]
}
```
