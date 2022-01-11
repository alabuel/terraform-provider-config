package config

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RequestParameters struct {
	uri      string
	params   []map[string]interface{}
	headers  []map[string]interface{}
	user     string
	password string
}

func dataSourceRestApiGet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRestApiRead,
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"param":  dataSourceKeyValueSchema(),
			"header": dataSourceKeyValueSchema(),
			"response": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceRestApiRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	reqparm := new(RequestParameters)
	reqparm.uri = d.Get("uri").(string)
	reqparm.user = d.Get("user").(string)
	reqparm.password = d.Get("password").(string)

	if v, ok := d.GetOk("param"); ok {
		reqparm.params = buildConfigDataSourceParams(v.(*schema.Set))
	}
	if v, ok := d.GetOk("header"); ok {
		reqparm.headers = buildConfigDataSourceParams(v.(*schema.Set))
	}

	data, err := getRequest(reqparm)
	if err != nil {
		return diag.FromErr(err)
	}
	if e := d.Set("response", data); e != nil {
		return diag.FromErr(e)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func getRequest(args *RequestParameters) (string, error) {
	param := url.Values{}
	for _, p := range args.params {
		param.Add(p["Key"].(string), p["Value"].(string))
	}

	url := args.uri
	if len(args.params) > 0 {
		url = args.uri + "?" + param.Encode()
	}

	method := "GET"
	payload := strings.NewReader(``)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}

	// add basic authentication
	if args.user != "" && args.password != "" {
		plainCred := args.user + ":" + args.password
		base64Cred := base64.StdEncoding.EncodeToString([]byte(plainCred))
		req.Header.Add("Authorization", "Basic "+base64Cred)
	}

	// add headers
	for _, h := range args.headers {
		req.Header.Add(h["Key"].(string), h["Value"].(string))
	}

	// Send http request
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}
