package config

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIni() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIniRead,
		Schema: map[string]*schema.Schema{
			"ini": {
				Type:     schema.TypeString,
				Required: true,
			},
			"section": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"json": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceIniRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	section := d.Get("section").(string)

	// var ini map[string]map[string]interface{}
	ini := iniParser(d.Get("ini"))
	if section != "" {
		data, _ := json.Marshal(ini[section])
		if e := d.Set("json", string(data)); e != nil {
			return diag.FromErr(e)
		}
	} else {
		data, _ := json.Marshal(ini)
		if e := d.Set("json", string(data)); e != nil {
			return diag.FromErr(e)
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
