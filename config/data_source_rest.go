package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RequestParameters struct {
	uri           string
	params        []map[string]interface{}
	headers       []map[string]interface{}
	user          string
	password      string
	authorization string
	method        string
	payload       string
}

type TokenRequestParameters struct {
	uri           string
	params        []map[string]interface{}
	headers       []map[string]interface{}
	user          string
	password      string
	client_id     string
	client_secret string
	grant_type    string
	payload       string
}

func dataSourceRestApiGet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRestApiRead,
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token_uri": {
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
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"grant_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "password",
			},
			"authorization": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "basic",
			},
			"param":  dataSourceKeyValueSchema(),
			"header": dataSourceKeyValueSchema(),
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "get",
			},
			"payload": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
	reqparm.uri = strings.TrimSpace(d.Get("uri").(string))
	reqparm.user = d.Get("user").(string)
	reqparm.password = d.Get("password").(string)
	reqparm.authorization = d.Get("authorization").(string)
	reqparm.method = d.Get("method").(string)
	reqparm.payload = d.Get("payload").(string)

	tokenparm := new(TokenRequestParameters)
	tokenparm.uri = strings.TrimSpace(d.Get("token_uri").(string))
	tokenparm.user = d.Get("user").(string)
	tokenparm.password = d.Get("password").(string)
	tokenparm.client_id = d.Get("client_id").(string)
	tokenparm.client_secret = d.Get("client_secret").(string)
	tokenparm.grant_type = d.Get("grant_type").(string)

	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(reqparm.uri)) {
		return diag.FromErr(fmt.Errorf("uri cannot contain whitespace. Got \"%s\"", reqparm.uri))
	}
	if whiteSpace.Match([]byte(tokenparm.uri)) {
		return diag.FromErr(fmt.Errorf("token_uri cannot contain whitespace. Got \"%s\"", reqparm.uri))
	}

	if v, ok := d.GetOk("param"); ok {
		reqparm.params = buildConfigDataSourceParams(v.(*schema.Set))
	}
	if v, ok := d.GetOk("header"); ok {
		reqparm.headers = buildConfigDataSourceParams(v.(*schema.Set))
	}

	if reqparm.authorization == "oauth2" {
		if tokenparm.user == "" {
			return diag.FromErr(fmt.Errorf("missing user for oath2 authentication"))
		}
		if tokenparm.password == "" {
			return diag.FromErr(fmt.Errorf("missing password for oath2 authentication"))
		}
		if tokenparm.client_id == "" {
			return diag.FromErr(fmt.Errorf("missing client_id for oath2 authentication"))
		}
		if tokenparm.client_secret == "" {
			return diag.FromErr(fmt.Errorf("missing client_secret for oath2 authentication"))
		}
		if tokenparm.grant_type == "" {
			return diag.FromErr(fmt.Errorf("missing grant_type for oath2 authentication"))
		}
		if tokenparm.uri == "" {
			return diag.FromErr(fmt.Errorf("missing token_uri for oath2 authentication"))
		}

	}

	data, err := createR
	equest(reqparm)
	if err != nil {
		return diag.FromErr(err)
	}
	if e := d.Set("response", data); e != nil {
		return diag.FromErr(e)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func createRequest(args *RequestParameters) (string, error) {
	param := url.Values{}
	for _, p := range args.params {
		param.Add(p["Key"].(string), p["Value"].(string))
	}

	url := args.uri
	if len(args.params) > 0 {
		url = args.uri + "?" + param.Encode()
	}

	method := args.method
	payload := strings.NewReader(``)
	if args.payload != "" {
		payload = strings.NewReader(args.payload)
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}

	// add basic authentication
	if args.authorization == "basic" && args.user != "" && args.password != "" {
		plainCred := args.user + ":" + args.password
		base64Cred := base64.StdEncoding.EncodeToString([]byte(plainCred))
		req.Header.Add("Authorization", "Basic "+base64Cred)
	}

	// add headers
	for _, h := range args.headers {
		req.Header.Set(h["Key"].(string), h["Value"].(string))
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
