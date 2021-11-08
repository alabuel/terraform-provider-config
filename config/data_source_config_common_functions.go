package config

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFilterSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func dataSourceLookupSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"column": {
					Type:     schema.TypeString,
					Required: true,
				},
				"worksheet": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"key_column": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value_column": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func buildConfigDataSourceFilters(set *schema.Set) []map[string]interface{} {
	var filters []map[string]interface{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, fmt.Sprintf("%v", e))
		}
		mvalue := make(map[string]interface{})
		mvalue["Name"] = m["name"].(string)
		mvalue["Values"] = filterValues
		filters = append(filters, mvalue)
	}
	return filters
}

func buildConfigDataSourceLookup(set *schema.Set) []map[string]interface{} {
	var lookup []map[string]interface{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		mvalue := make(map[string]interface{})
		mvalue["Column"] = m["column"].(string)
		if m["worksheet"] != nil {
			mvalue["Worksheet"] = m["worksheet"].(string)
		} else {
			mvalue["Worksheet"] = nil
		}
		mvalue["Key"] = m["key_column"].(string)
		mvalue["Value"] = m["value_column"].(string)
		lookup = append(lookup, mvalue)
	}
	return lookup
}

func checkLookupValue(lookup []map[string]interface{}, key string) bool {
	for _, lv := range lookup {
		if lv["Column"].(string) == key {
			return true
		}
	}
	return false
}

func getLookupValue(lookup []map[string]interface{}, excel_file, default_worksheet string, key string, value string) (string, error) {
	var lookupValue = ""
	for _, lv := range lookup {
		if lv["Column"].(string) == key {
			f, err := excelize.OpenFile(excel_file)
			if err != nil {
				return "", err
			}
			worksheet := ""
			if lv["Worksheet"].(string) != "" {
				worksheet = lv["Worksheet"].(string)
			} else {
				worksheet = default_worksheet
			}
			rows, err := f.GetRows(worksheet)
			if err != nil {
				return "", fmt.Errorf(fmt.Sprintf("%v", rows))
			}

			columns := len(rows[0])

			// get column of key
			header := rows[0]
			column_key := 0
			column_value := -1
			for i := 0; i < columns; i++ {
				if header[i] == key {
					column_key = i
				}
				if header[i] == lv["Value"] {
					column_value = i
				}
			}
			// get row of key
			if column_value >= 0 {
				for _, row := range rows {
					if row[column_key] == value {
						lookupValue = row[column_value]
					}
				}
			} else {
				return "", fmt.Errorf("lookup value not found")
			}
		}
	}
	return lookupValue, nil
}

func checkFiltersForItem(filters []map[string]interface{}, key string, value string) bool {
	for _, fv := range filters {
		if fv["Name"] == key {
			for _, d := range fv["Values"].([]string) {
				if d == value {
					return true
				}
			}
		}
	}
	return false
}

func stringInList(s string, list []string) bool {
	for _, b := range list {
		if b == s {
			return true
		}
	}
	return false
}
