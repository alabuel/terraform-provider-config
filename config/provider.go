package config

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"config_workbook":    dataSourceConfigurationWorkbook(),
			"config_ini":         dataSourceIni(),
			"config_restapi_get": dataSourceRest(),
			"config_rest":        dataSourceRest(),
		},
	}
}
