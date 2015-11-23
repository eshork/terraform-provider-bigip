package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/bitly/go-simplejson"
	"log"
	"fmt"
	"strconv"
	"strings"
	// "errors"
)

func resourceNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceNodeCreate,
		Read:   resourceNodeRead,
		Update: resourceNodeUpdate,
		Delete: resourceNodeDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"connection_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default: 0,
			},
			"connection_rate_limit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default: 0,
			},
		},
	}
}

func resourceNodeCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourceNodeCreateresourceNodeCreateresourceNodeCreateresourceNodeCreateresourceNodeCreateresourceNodeCreate")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_node


	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	address := d.Get("address").(string)
	enabled := d.Get("enabled").(bool)
	description := d.Get("description").(string)
	connection_limit := d.Get("connection_limit")
	connection_rate_limit := d.Get("connection_rate_limit")
	destID := "~" + partition + "~" + name

	json := simplejson.New()
	json.Set("name", name)
	json.Set("partition", partition)
	json.Set("address", address)
	json.Set("description", description)
	json.Set("connectionLimit", connection_limit)
	json.Set("rateLimit", connection_rate_limit)

	if !enabled {
		// the default is to create in an enabled state, so we only do this if we want it created in disabled state
		json.Set("session", "user-disabled")
	}

	jsonPostString := JSONtoString(json)
	jsonPostString = "{" + jsonPostString + "}"

	log.Println("[BIGIP] resourceNodeCreate restURL : " + restURL)
	log.Println("[BIGIP] resourceNodeCreate destID : " + destID)
	log.Println("[BIGIP] resourceNodeCreate jsonPostString : " + jsonPostString)

	jsonRet, err := F5Post(restURL, jsonPostString, *client)
	if err != nil {
		return err
	}

	log.Println("[BIGIP] resourceNodeCreate destID : " + jsonRet.MustString() )
	d.SetId(destID)

	return nil
}

func resourceNodeRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourceNodeReadresourceNodeReadresourceNodeReadresourceNodeRead : ")


	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_node


	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	address := d.Get("address").(string)
	enabled := d.Get("enabled").(bool)
	var enabledStr string
	if enabled {
		enabledStr = "true"
	} else {
		enabledStr = "false"
	}
	description := d.Get("description").(string)
	connection_limit := strconv.Itoa( d.Get("connection_limit").(int) )
	connection_rate_limit := strconv.Itoa( d.Get("connection_rate_limit").(int) )
	// destID := "~" + partition + "~" + name
	destID := d.Id()

	restURL = restURL + "/" + destID
	log.Println("[BIGIP] resourceNodeRead restURL : " + restURL)


	myJson, f5err := F5Get(restURL, *client)

	if f5err != nil {
		log.Println("[BIGIP] resourceNodeRead if f5err != nil : ")
		if strings.Contains( f5err.Error() , "404") {
			log.Println("[BIGIP] resourceNodeRead NODE WAS DELETED : " + destID)
			d.SetId("")
			return nil
		}
		fmt.Println(f5err)
		return f5err
	}

	log.Println("[BIGIP] resourceNodeRead JSON : " + myJson.MustString() )

	log.Println("[BIGIP] resourceNodeRead name : " + name)
	log.Println("[BIGIP] resourceNodeRead partition : " + partition)
	log.Println("[BIGIP] resourceNodeRead address : " + address)
	log.Println("[BIGIP] resourceNodeRead description : " + description)
	log.Println("[BIGIP] resourceNodeRead connection_limit : " + connection_limit)
	log.Println("[BIGIP] resourceNodeRead connection_rate_limit : " + connection_rate_limit)
	log.Println("[BIGIP] resourceNodeRead enabled : " + enabledStr)

	/*
	{
	kind: "tm:ltm:node:nodestate",
	name: "testNode",
	partition: "Website",
	fullPath: "/Website/testNode",
	generation: 59825295,
	selfLink: "https://localhost/mgmt/tm/ltm/node/~Website~testNode?ver=11.5.1",
	address: "1.1.1.1",
	connectionLimit: 0,
	description: "just a test",
	dynamicRatio: 1,

	logging: "disabled",
	monitor: "default",
	rateLimit: "disabled",
	ratio: 1,
	session: "monitor-enabled",
	state: "down"
	}
	*/

	resKind := myJson.Get("kind").MustString("")
	resName := myJson.Get("name").MustString("")
	resPartition := myJson.Get("partition").MustString("")
	resFullPath := myJson.Get("fullPath").MustString("")
	resGeneration := myJson.Get("generation").MustString("")
	resSelfLink := myJson.Get("selfLink").MustString("")
	resAddress := myJson.Get("address").MustString("")

	resConnectionLimit, err := myJson.Get("connectionLimit").Int()
	if err != nil {
		resConnectionLimit = 0
	}

	resRatio, err := myJson.Get("ratio").Int()
	if err != nil {
		resRatio = 0
	}
	

	// rateLimit is a weird Int as a string value with 0 == 'disabled'
	resRateLimit := myJson.Get("rateLimit").MustString("")
	if resRateLimit == "disabled" {
		resRateLimit = "0"
	}
	resRateLimitInt, err := strconv.Atoi(resRateLimit)
	if err != nil {
		resRateLimitInt = 0
	}


	resDescription := myJson.Get("description").MustString("")
	resDynamicRatio := myJson.Get("dynamicRatio").MustString("")
	resLogging := myJson.Get("logging").MustString("")
	resMonitor := myJson.Get("monitor").MustString("")
	resSession := myJson.Get("session").MustString("")
	resState := myJson.Get("state").MustString("")


	log.Println("[BIGIP] resourceNodeRead resKind : " + resKind)
	log.Println("[BIGIP] resourceNodeRead resName : " + resName)
	log.Println("[BIGIP] resourceNodeRead resPartition : " + resPartition)
	log.Println("[BIGIP] resourceNodeRead resFullPath : " + resFullPath)
	log.Println("[BIGIP] resourceNodeRead resGeneration : " + resGeneration)
	log.Println("[BIGIP] resourceNodeRead resSelfLink : " + resSelfLink)
	log.Println("[BIGIP] resourceNodeRead resAddress : " + resAddress)
	log.Println("[BIGIP] resourceNodeRead resConnectionLimit : " + strconv.Itoa(resConnectionLimit))
	log.Println("[BIGIP] resourceNodeRead resDescription : " + resDescription)
	log.Println("[BIGIP] resourceNodeRead resDynamicRatio : " + resDynamicRatio)
	log.Println("[BIGIP] resourceNodeRead resLogging : " + resLogging)
	log.Println("[BIGIP] resourceNodeRead resMonitor : " + resMonitor)
	log.Println("[BIGIP] resourceNodeRead resRateLimit : " + strconv.Itoa(resRateLimitInt))
	log.Println("[BIGIP] resourceNodeRead resRatio : " + strconv.Itoa(resRatio))
	log.Println("[BIGIP] resourceNodeRead resSession : " + resSession)
	log.Println("[BIGIP] resourceNodeRead resState : " + resState)

	d.Set("address",resAddress)
	d.Set("description",resDescription)
	d.Set("connection_limit",resConnectionLimit)
	d.Set("connection_rate_limit", resRateLimitInt)
	d.Set("ratio", resRatio)
	if resSession == "user-disabled" {
		d.Set("enabled", false)
	}
	if resSession == "monitor-enabled" {
		d.Set("enabled", true)
	}

	return nil
}

func resourceNodeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourceNodeUpdate : ")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_node
	destID := d.Id()
	restURL = restURL + "/" + destID
	log.Println("[BIGIP] resourceNodeUpdate restURL : " + restURL)

	json := simplejson.New()

	var needsFullReCreate bool = false

	// try to figure out what parts need to be updated
	if d.HasChange("name") {
		log.Println("[BIGIP] resourceNodeUpdate : name CHANGED")
		json.Set("name", d.Get("name"))
		// THIS WILL NOT WORK - IS READ-ONLY IN WEB UI - NEEDS DESTROY/CREATE TO PERFORM
		needsFullReCreate = true
	}
	if d.HasChange("address") {
		log.Println("[BIGIP] resourceNodeUpdate : address CHANGED")
		json.Set("address", d.Get("address"))
		// THIS WILL NOT WORK - IS READ-ONLY IN WEB UI - NEEDS DESTROY/CREATE TO PERFORM
		needsFullReCreate = true
	}
	if d.HasChange("description") {
		log.Println("[BIGIP] resourceNodeUpdate : description CHANGED")
		json.Set("description", d.Get("description"))
	}
	if d.HasChange("connection_limit") {
		log.Println("[BIGIP] resourceNodeUpdate : connection_limit CHANGED")
		json.Set("connectionLimit", d.Get("connection_limit"))
	}
	if d.HasChange("connection_rate_limit") {
		log.Println("[BIGIP] resourceNodeUpdate : connection_rate_limit CHANGED")
		if d.Get("connection_rate_limit") == 0 {
			json.Set("rateLimit", "disabled")
		} else {
			json.Set("rateLimit", d.Get("connection_rate_limit"))
		}
	}
	if d.HasChange("ratio") {
		log.Println("[BIGIP] resourceNodeUpdate : ratio CHANGED")
		json.Set("ratio", d.Get("ratio"))
	}
	if d.HasChange("enabled") {
		log.Println("[BIGIP] resourceNodeUpdate : enabled CHANGED")
		if d.Get("enabled") == true {
			json.Set("session", "user-enabled")
		} else {
			json.Set("session", "user-disabled")
		}
	}

	jsonPostString := JSONtoString(json)
	jsonPostString = "{" + jsonPostString + "}"

	log.Println("[BIGIP] resourceNodeUpdate jsonPutString : " + jsonPostString)

	if needsFullReCreate {
		// first delete the resource
		rdelerr := resourceNodeDelete(d, m)
		if rdelerr != nil {
			return rdelerr
		}
		rcreateerr := resourceNodeCreate(d, m)
		if rcreateerr != nil {
			return rcreateerr
		}
	} else {
		// a simple update is okay this time
		jsonRet, err := F5Put(restURL, jsonPostString, *client)
		if err != nil {
			return err
		}
		log.Println("[BIGIP] resourceNodeUpdate JSON RET : " + jsonRet.MustString() )
	}

	return nil
}

func resourceNodeDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourceNodeDelete : ")


	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_node

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	
	destID := d.Id()

	restURL = restURL + "/" + destID
	log.Println("[BIGIP] resourceNodeDelete restURL : " + restURL)
	log.Println("[BIGIP] resourceNodeDelete name : " + name)
	log.Println("[BIGIP] resourceNodeDelete partition : " + partition)

	_, f5err := F5Delete(restURL, *client)

	if f5err != nil {
		log.Println("[BIGIP] resourceNodeDelete if f5err != nil : ")
		if strings.Contains( f5err.Error() , "404") {
			log.Println("[BIGIP] resourceNodeDelete NODE NOT FOUND : " + destID)
			// technically, if it is already missing, that's good enough for success!
			d.SetId("")
			return nil
		}
		fmt.Println(f5err)
		return f5err
	}

	// we did it! yay!
	d.SetId("")
	return nil
}

// func readFromConfig( d *schema.ResourceData, resourceType string, resourceID string, property string ) interface{} {
// 	if v, ok := d.GetOk(resourceType); ok {
// 		vL := v.(*schema.Map).List()
// 		for _, v := range vL {
// 		}
// 	}
// }


