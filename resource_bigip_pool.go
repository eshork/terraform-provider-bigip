package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/bitly/go-simplejson"
	"log"
	"fmt"
	// "strconv"
	"strings"
	"errors"
)


func resourcePool() *schema.Resource {
	return &schema.Resource{
		Create: resourcePoolCreate,
		Read:   resourcePoolRead,
		Update: resourcePoolUpdate,
		Delete: resourcePoolDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"balancing_method": &schema.Schema{
				Type:     schema.TypeString,
				// options are: 
				//   round_robin
				//   ratio_member
				//   least_connections_member
				//   observed_member
				//   predictive_member
				//   ratio_node
				//   least_connections_node
				//   fastest_node
				//   observed_node
				//   predictive_node
				//   dynamic_ratio_node
				//   fastest_application
				//   least_sessions
				//   dynamic_ratio_member
				//   weighted_least_connections_member
				//   weighted_least_connections_node
				//   ratio_session
				//   ratio_least_connections_member
				//   ratio_least_connections_node
				Optional: true,
				Default: "round_robin",
				// REST param is "loadBalancingMode"
			},

			// "address": &schema.Schema{
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// },
			// "enabled": &schema.Schema{
			// 	Type:     schema.TypeBool,
			// 	Optional: true,
			// 	Default: true,
			// },
			// "connection_limit": &schema.Schema{
			// 	Type:     schema.TypeInt,
			// 	Optional: true,
			// 	Default: 0,
			// },
			// "connection_rate_limit": &schema.Schema{
			// 	Type:     schema.TypeInt,
			// 	Optional: true,
			// 	Default: 0,
			// },
		},
	}
}


func resourcePoolCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourcePoolCreate")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_pool
	log.Println("[BIGIP] resourcePoolCreate restURL : " + restURL)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	description := d.Get("description").(string)
	balancing_method := d.Get("balancing_method").(string)
	destID := "~" + partition + "~" + name

	log.Println("[BIGIP] resourcePoolCreate name : " + name)
	log.Println("[BIGIP] resourcePoolCreate partition : " + partition)
	log.Println("[BIGIP] resourcePoolCreate description : " + description)
	log.Println("[BIGIP] resourcePoolCreate balancing_method : " + balancing_method)

	json := simplejson.New()
	json.Set("name", name)
	json.Set("partition", partition)
	json.Set("description", description)
	json.Set("loadBalancingMode", xlate_toPoolBalanceMode(balancing_method))

	jsonPostString := JSONtoString(json)
	jsonPostString = "{" + jsonPostString + "}"

	log.Println("[BIGIP] resourcePoolCreate jsonPostString : " + jsonPostString)

	jsonRet, err := F5Post(restURL, jsonPostString, *client)
	if err != nil {
		return err
	}

	log.Println("[BIGIP] resourcePoolCreate destID : " + jsonRet.MustString() )
	d.SetId(destID)

	return nil

	return errors.New("NYI")
}


func resourcePoolRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourcePoolRead")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_pool
	log.Println("[BIGIP] resourcePoolCreate restURL : " + restURL)

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	description := d.Get("description").(string)
	balancing_method := d.Get("balancing_method").(string)
	destID := d.Id()

	restURL = restURL + "/" + destID
	log.Println("[BIGIP] resourcePoolRead restURL : " + restURL)

	log.Println("[BIGIP] resourcePoolRead name : " + name)
	log.Println("[BIGIP] resourcePoolRead partition : " + partition)
	log.Println("[BIGIP] resourcePoolRead description : " + description)
	log.Println("[BIGIP] resourcePoolRead balancing_method : " + balancing_method)

	myJson, f5err := F5Get(restURL, *client)
	if f5err != nil {
		log.Println("[BIGIP] resourcePoolRead if f5err != nil : ")
		if strings.Contains( f5err.Error() , "404") {
			log.Println("[BIGIP] resourcePoolRead NODE WAS DELETED : " + destID)
			d.SetId("")
			return nil
		}
		fmt.Println(f5err)
		return f5err
	}
	log.Println("[BIGIP] resourcePoolRead JSON : " + myJson.MustString() )

	/*
	{
	kind: "tm:ltm:pool:poolstate",
	name: "tfTestPool",
	partition: "Website",
	fullPath: "/Website/tfTestPool",
	generation: 62812786,
	selfLink: "https://localhost/mgmt/tm/ltm/pool/~Website~tfTestPool?ver=11.5.1",
	allowNat: "yes",
	allowSnat: "yes",
	description: "test pool",
	ignorePersistedWeight: "disabled",
	ipTosToClient: "pass-through",
	ipTosToServer: "pass-through",
	linkQosToClient: "pass-through",
	linkQosToServer: "pass-through",
	loadBalancingMode: "round-robin",
	minActiveMembers: 0,
	minUpMembers: 0,
	minUpMembersAction: "failover",
	minUpMembersChecking: "disabled",
	queueDepthLimit: 0,
	queueOnConnectionLimit: "disabled",
	queueTimeLimit: 0,
	reselectTries: 0,
	slowRampTime: 10,
		membersReference: {
		link: "https://localhost/mgmt/tm/ltm/pool/~Website~tfTestPool/members?ver=11.5.1",
		isSubcollection: true
		}
	}
	*/

	resKind := myJson.Get("kind").MustString("")
	resName := myJson.Get("name").MustString("")
	resPartition := myJson.Get("partition").MustString("")
	resFullPath := myJson.Get("fullPath").MustString("")
	resGeneration := myJson.Get("generation").MustString("")
	resSelfLink := myJson.Get("selfLink").MustString("")
	resDescription := myJson.Get("description").MustString("")
	resBalancingMethod := myJson.Get("loadBalancingMode").MustString("")

	log.Println("[BIGIP] resourcePoolRead resKind : " + resKind)
	log.Println("[BIGIP] resourcePoolRead resName : " + resName)
	log.Println("[BIGIP] resourcePoolRead resPartition : " + resPartition)
	log.Println("[BIGIP] resourcePoolRead resFullPath : " + resFullPath)
	log.Println("[BIGIP] resourcePoolRead resGeneration : " + resGeneration)
	log.Println("[BIGIP] resourcePoolRead resSelfLink : " + resSelfLink)
	log.Println("[BIGIP] resourcePoolRead resDescription : " + resDescription)
	log.Println("[BIGIP] resourcePoolRead resBalancingMethod : " + resBalancingMethod)

	d.Set("description",resDescription)
	d.Set("balancing_method", xlate_fromPoolBalanceMode(resBalancingMethod))

	return nil
}


func resourcePoolUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourcePoolUpdate : ")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_pool
	destID := d.Id()
	restURL = restURL + "/" + destID
	log.Println("[BIGIP] resourcePoolUpdate restURL : " + restURL)

	json := simplejson.New()

	var needsFullReCreate bool = false

	// try to figure out what parts need to be updated
	if d.HasChange("name") {
		log.Println("[BIGIP] resourcePoolUpdate : name CHANGED")
		json.Set("name", d.Get("name"))
		// THIS WILL NOT WORK - IS READ-ONLY IN WEB UI - NEEDS DESTROY/CREATE TO PERFORM
		needsFullReCreate = true // and this makes DESTROY/CREATE a reality!
	}
	if d.HasChange("description") {
		log.Println("[BIGIP] resourcePoolUpdate : description CHANGED")
		json.Set("description", d.Get("description"))
	}
	if d.HasChange("balancing_method") {
		log.Println("[BIGIP] resourcePoolUpdate : balancing_method CHANGED")
		json.Set("loadBalancingMode", xlate_toPoolBalanceMode(d.Get("balancing_method").(string)))
	}

	jsonPostString := JSONtoString(json)
	jsonPostString = "{" + jsonPostString + "}"

	log.Println("[BIGIP] resourcePoolUpdate jsonPutString : " + jsonPostString)

	if needsFullReCreate {
		// first delete the resource
		rdelerr := resourcePoolDelete(d, m)
		if rdelerr != nil {
			return rdelerr
		}
		rcreateerr := resourcePoolCreate(d, m)
		if rcreateerr != nil {
			return rcreateerr
		}
	} else {
		// a simple update is okay this time
		jsonRet, err := F5Put(restURL, jsonPostString, *client)
		if err != nil {
			return err
		}
		log.Println("[BIGIP] resourcePoolUpdate JSON RET : " + jsonRet.MustString() )
	}

	return nil
}


func resourcePoolDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourcePoolDelete : ")

	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_pool

	name := d.Get("name").(string)
	partition := d.Get("partition").(string)
	destID := d.Id()
	restURL = restURL + "/" + destID

	log.Println("[BIGIP] resourcePoolDelete restURL : " + restURL)
	log.Println("[BIGIP] resourcePoolDelete name : " + name)
	log.Println("[BIGIP] resourcePoolDelete partition : " + partition)

	_, f5err := F5Delete(restURL, *client)
	
	if f5err != nil {
		log.Println("[BIGIP] resourcePoolDelete if f5err != nil : ")
		if strings.Contains( f5err.Error() , "404") {
			log.Println("[BIGIP] resourcePoolDelete NODE NOT FOUND : " + destID)
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

	return errors.New("NYI") // yep, we never hit this anymore, but Go has pedantic imports, so... #reasons
}

// do we really need to define brand new (and only slightly different) setting strings here? nope; but we did, so live with it -- at least we did it with a map[string]string
var resourcePoolBalanceModeMap map[string]string = map[string]string{
	// terraform : F5REST
	"round_robin" : "round-robin",
	"ratio_member" : "ratio-member",
	"least_connections_member" : "least-connections-member",
	"observed_member" : "observed-member",
	"predictive_member" : "predictive-member",
	"ratio_node" : "ratio-node",
	"least_connections_node" : "least-connections-node",
	"fastest_node" : "fastest-node",
	"observed_node" : "observed-node",
	"predictive_node" : "predictive-node",
	"dynamic_ratio_node" : "dynamic-ratio-node",
	"fastest_application" : "fastest-app-response",
	"least_sessions" : "least-sessions",
	"dynamic_ratio_member" : "dynamic-ratio-member",
	"weighted_least_connections_member" : "weighted-least-connections-member",
	"weighted_least_connections_node" : "weighted-least-connections-node",
	"ratio_session" : "ratio-session",
	"ratio_least_connections_member" : "ratio-least-connections-member",
	"ratio_least_connections_node" : "ratio-least-connections-node",
}

// pass in the user-provided balancing_method string, and get back the F5-expected value for RESTAPI calls
func xlate_toPoolBalanceMode(userMode string) string {
	if val, ok := resourcePoolBalanceModeMap[userMode]; ok {
		//if we have the index, return the value
		return val
	}
	// well that didn't work as we planned, so: ROUND ROBIN!
	log.Println("[BIGIP] xlate_toPoolBalanceMode FAILED! defaulting to round-robin! : ")
	return "round-robin"
}

// pass in the F5 RESTAPI balancing_method string, and get back the user-expected value
func xlate_fromPoolBalanceMode(userMode string) string {
	//reverse the map and do a lookup
	reversed_resourcePoolBalanceModeMap := make(map[string]string)
	for k, v := range resourcePoolBalanceModeMap{
		reversed_resourcePoolBalanceModeMap[v] = k
	}
	if val, ok := reversed_resourcePoolBalanceModeMap[userMode]; ok {
		//if we have the index, return the value
		return val
	}

	// well that didn't work as we planned, so: ROUND ROBIN!
	log.Println("[BIGIP] xlate_fromPoolBalanceMode FAILED! defaulting to round_robin! : ")
	return "round_robin"
}
