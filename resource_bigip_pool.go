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
	destID := "~" + partition + "~" + name

	log.Println("[BIGIP] resourcePoolCreate name : " + name)
	log.Println("[BIGIP] resourcePoolCreate partition : " + partition)
	log.Println("[BIGIP] resourcePoolCreate description : " + description)

	json := simplejson.New()
	json.Set("name", name)
	json.Set("partition", partition)
	json.Set("description", description)

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
	destID := d.Id()

	restURL = restURL + "/" + destID
	log.Println("[BIGIP] resourcePoolRead restURL : " + restURL)

	log.Println("[BIGIP] resourcePoolRead name : " + name)
	log.Println("[BIGIP] resourcePoolRead partition : " + partition)
	log.Println("[BIGIP] resourcePoolRead description : " + description)


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

	log.Println("[BIGIP] resourcePoolRead resKind : " + resKind)
	log.Println("[BIGIP] resourcePoolRead resName : " + resName)
	log.Println("[BIGIP] resourcePoolRead resPartition : " + resPartition)
	log.Println("[BIGIP] resourcePoolRead resFullPath : " + resFullPath)
	log.Println("[BIGIP] resourcePoolRead resGeneration : " + resGeneration)
	log.Println("[BIGIP] resourcePoolRead resSelfLink : " + resSelfLink)
	log.Println("[BIGIP] resourcePoolRead resDescription : " + resDescription)

	d.Set("description",resDescription)

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
		needsFullReCreate = true
	}
	if d.HasChange("description") {
		log.Println("[BIGIP] resourcePoolUpdate : description CHANGED")
		json.Set("description", d.Get("description"))
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

	return errors.New("NYI")
}




