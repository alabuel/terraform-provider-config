package config

import (
	"fmt"

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
