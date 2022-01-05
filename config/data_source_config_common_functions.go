package config

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
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
				"excel": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"password": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"worksheet": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"yaml": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"json": {
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

func buildConfigDataSourceLookup(set *schema.Set) ([]map[string]interface{}, error) {
	var lookup []map[string]interface{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		mvalue := make(map[string]interface{})
		mvalue["Column"] = m["column"].(string)

		source := 0
		if m["worksheet"].(string) != "" {
			source++
		}
		if m["json"].(string) != "" {
			source++
		}
		if m["yaml"].(string) != "" {
			source++
		}
		if source > 1 {
			return nil, fmt.Errorf("only 1 type of lookup source is required (worksheet/json/yaml)")
		}

		if m["worksheet"].(string) != "" {
			if m["excel"].(string) != "" {
				mvalue["Excel"] = m["excel"].(string)
			}
			if m["password"].(string) != "" {
				mvalue["Password"] = m["password"].(string)
			}
			mvalue["Worksheet"] = m["worksheet"].(string)
		} else {
			mvalue["Excel"] = nil
			mvalue["Worksheet"] = nil
			mvalue["Password"] = nil
		}
		if m["json"].(string) != "" {
			mvalue["Json"] = m["json"].(string)
		} else {
			mvalue["Json"] = nil
		}
		if m["yaml"].(string) != "" {
			mvalue["Yaml"] = m["yaml"].(string)
		} else {
			mvalue["Yaml"] = nil
		}

		mvalue["Key"] = m["key_column"].(string)
		mvalue["Value"] = m["value_column"].(string)
		lookup = append(lookup, mvalue)
	}
	return lookup, nil
}

func checkLookupValue(lookup []map[string]interface{}, key string) bool {
	for _, lv := range lookup {
		if lv["Column"].(string) == key {
			return true
		}
	}
	return false
}

func getLookupValue(lookup []map[string]interface{}, default_excel string, default_password string, default_worksheet string, key string, value string) (string, error) {
	var lookupValue = ""
	for _, lv := range lookup {
		if lv["Column"].(string) == key {
			if lv["Json"] != nil {
				jd, _ := stringToInterface(lv["Json"].(string))
				jsondata := jd.(map[interface{}]interface{})
				if jsondata[value] != nil {
					lookupValue = jsondata[value].(string)
				} else {
					lookupValue = ""
				}
			} else if lv["Yaml"] != nil {
				yd, _ := stringToInterface(lv["Yaml"].(string))
				yamldata := yd.(map[interface{}]interface{})
				if yamldata[value] != nil {
					lookupValue = yamldata[value].(string)
				} else {
					lookupValue = ""
				}
			} else if lv["Worksheet"] != nil {
				excel_file := default_excel
				excel_pass := default_password
				if lv["Excel"] != nil {
					excel_file = lv["Excel"].(string)
				}
				if lv["Password"] != nil {
					excel_pass = lv["Password"].(string)
				}
				f, err := excelize.OpenFile(excel_file, excelize.Options{Password: excel_pass})
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
					lookupValue = ""
				}
			} else {
				lookupValue = ""
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

func stringToInterface(s string) (interface{}, error) {
	var v interface{}

	// Try if the string is yaml
	err := yaml.Unmarshal([]byte(s), &v)
	if err != nil {
		// Try if the string is json
		err = json.Unmarshal([]byte(s), &v)
		if err != nil {
			return nil, fmt.Errorf("unable to parse string using yaml or json")
		}
	}
	return v, nil
}

func stringToMap(s string) ([]map[string]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows, nil
}
