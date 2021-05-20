package config

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"
)

func dataSourceConfigurationWorkbook() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationItemRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"csv": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration_item": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceConfigurationItemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// get all arguments passed to the resource
	csv_string := d.Get("csv").(string)
	schema := d.Get("schema").(string)
	configuration_item := d.Get("configuration_item").(string)

	// convert the schema to map
	map_yaml, err := stringToInterface(schema)
	if err != nil {
		return diag.FromErr(err)
	}
	maps := map_yaml.(map[interface{}]interface{})

	// convert the csv to map
	csv, err := stringToMap(csv_string)
	if err != nil {
		return diag.FromErr(err)
	}

	// remap all csv headers based on mapping configuration
	records := reMapData(csv, maps["configuration_workbook_mapping"])

	// get all unique configuration items
	items := unique(getConfigurationItems(records))

	// get the transformed data
	data := getItemData(records, items, configuration_item)

	// set the data to the attribute json
	if err := d.Set("json", data); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func getConfigurationItems(csv []map[string]interface{}) []string {
	var items []string
	for _, value := range csv {
		for k, v := range value {
			if k == "configuration_item" {
				items = append(items, v.(string))
			}
		}
	}
	return items
}

func getItemData(csv []map[string]interface{}, items []string, configuration_item string) string {
	listitem := make(map[string][]map[string]interface{})
	for _, item := range items {
		var itemdatalist []map[string]interface{}
		for _, value := range csv {
			if value["configuration_item"] == item {
				itemdata := make(map[string]interface{})
				for k, v := range value {
					if k != "configuration_item" {
						itemdata[k] = v
					}
				}
				itemdatalist = append(itemdatalist, itemdata)
				if configuration_item != "" {
					if configuration_item == item {
						listitem[item] = itemdatalist
					}
				} else {
					listitem[item] = itemdatalist
				}

			}
		}
	}
	j, _ := json.Marshal(listitem)
	return string(j)
}

func unique(items []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range items {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func reMapData(csv []map[string]string, mapping interface{}) []map[string]interface{} {
	new_csv := make([]map[string]interface{}, len(csv))
	for key, value := range csv {
		item_key := ""
		for k, v := range value {
			if k == "configuration_item" {
				item_key = v
			}
		}
		new_value := make(map[string]interface{})
		new_tag := make(map[string]string)
		tags := make(map[string]string)
		for k, v := range value {
			_ = v
			if strings.HasPrefix(k, "attr") {
				var new_key, new_type string
				new_key, new_type = getMapValue(mapping, item_key, k)
				if new_key != "" {
					if new_type == "string" {
						new_value[new_key] = value[k]
					} else if new_type == "number" || new_type == "numeric" {
						n, _ := strconv.ParseFloat(value[k], 64)
						new_value[new_key] = n
					} else if new_type == "boolean" || new_type == "bool" {
						val, _ := strconv.ParseBool(value[k])
						new_value[new_key] = val
					} else if new_type == "list" {
						new_value[new_key] = strings.Split(value[k], ",")
					} else if new_type == "map" || new_type == "hash" {
						vlist := strings.Split(value[k], ",")
						vmap := make(map[string]string)
						for _, vl := range vlist {
							mlist := strings.Split(vl, "=")
							vmap[mlist[0]] = mlist[1]
						}
						new_value[new_key] = vmap
					} else if new_type == "tag" {
						new_tag[strings.Title(new_key)] = value[k]
					} else {
						new_value[new_key] = value[k]
					}
				}
			} else if k == "configuration_item" {
				new_value[k] = value[k]
			} else if strings.HasPrefix(k, "t_") {
				if value[k] != "" {
					new_key := strings.Title(strings.Replace(k, "t_", "", -1))
					new_tag[new_key] = value[k]
				}
			} else if strings.HasPrefix(k, "n_") {
				if value[k] != "" {
					new_key := strings.Replace(k, "n_", "", -1)
					n, _ := strconv.ParseFloat(value[k], 64)
					new_value[new_key] = n
				}
			} else if strings.HasPrefix(k, "b_") {
				if value[k] != "" {
					new_key := strings.Replace(k, "b_", "", -1)
					val, _ := strconv.ParseBool(value[k])
					new_value[new_key] = val
				}
			} else if strings.HasPrefix(k, "l_") {
				if value[k] != "" {
					new_key := strings.Replace(k, "l_", "", -1)
					new_value[new_key] = strings.Split(value[k], ",")
				}
			} else if strings.HasPrefix(k, "m_") {
				if value[k] != "" {
					new_key := strings.Replace(k, "m_", "", -1)
					vlist := strings.Split(value[k], ",")
					vmap := make(map[string]string)
					for _, vl := range vlist {
						mlist := strings.Split(vl, "=")
						vmap[mlist[0]] = mlist[1]
					}
					new_value[new_key] = vmap
				}
			}
		}
		for k, v := range new_tag {
			tags[k] = v
		}
		new_value["tags"] = tags
		new_csv[key] = new_value
	}
	return new_csv
}

func getMapValue(config interface{}, config_item string, config_key string) (string, string) {
	return_value := ""
	return_type := "string"
	item := config.(map[interface{}]interface{})
	for _, key := range reflect.ValueOf(item).MapKeys() {
		if key.Interface().(string) == config_item {
			item_map := item[key.Interface()].(map[interface{}]interface{})
			for _, ik := range reflect.ValueOf(item_map).MapKeys() {
				if ik.Interface().(string) == config_key {
					return_type = "string"
					kv, ok := item_map[ik.Interface()].(map[interface{}]interface{})
					if !ok {
						return_value = item_map[ik.Interface()].(string)
					} else {
						for _, ikv := range reflect.ValueOf(kv).MapKeys() {
							if ikv.Interface().(string) == "name" {
								return_value = kv[ikv.Interface()].(string)
							}
							if ikv.Interface().(string) == "type" {
								return_type = kv[ikv.Interface()].(string)
							}
						}
					}
				}
			}
		}
	}
	return return_value, return_type
}

// func readCsvFile(filePath string) ([]map[string]string, error) {
// 	f, err := os.Open(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("Unable to read file %s", filePath)
// 	}
// 	defer f.Close()

// 	csvmap, err := csvToMap(f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return csvmap, nil
// }

// func csvToMap(reader io.Reader) ([]map[string]string, error) {
// 	r := csv.NewReader(reader)
// 	rows := []map[string]string{}
// 	var header []string
// 	for {
// 		record, err := r.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		if header == nil {
// 			header = record
// 		} else {
// 			dict := map[string]string{}
// 			for i := range header {
// 				dict[header[i]] = record[i]
// 			}
// 			rows = append(rows, dict)
// 		}
// 	}
// 	return rows, nil
// }

// func readYamlFile(filePath string) (interface{}, error) {
// 	f, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("Unable to read file %s", filePath)
// 	}

// 	var v interface{}
// 	err = yaml.Unmarshal(f, &v)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return v, nil
// }

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
