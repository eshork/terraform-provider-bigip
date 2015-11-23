package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/bitly/go-simplejson"
	"log"
	"fmt"
	"strconv"
	"strings"
	"errors"
)


func resourcePoolMember() *schema.Resource {
	return &schema.Resource{
		Create: resourcePoolMemberCreate,
		Read:   resourcePoolMemberRead,
		Update: resourcePoolMemberUpdate,
		Delete: resourcePoolMemberDelete,

		Schema: map[string]*schema.Schema{
			"node_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pool_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"node_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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


func resourcePoolMemberCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourcePoolMemberCreate")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_pool

	pool_id := d.Get("pool_id").(string)
	node_id := d.Get("node_id").(string)
	node_name := d.Get("node_name").(string)
	description := d.Get("description").(string)
	partition := d.Get("partition").(string)
	address := d.Get("address").(string)
	port := d.Get("port").(int)
	portStr := strconv.Itoa(port)

	name := node_name + ":" + portStr
	d.Set("name",name)


	destID := "~" + partition + "~" + name

	restURL = restURL + "/" + pool_id + "/members/"

	log.Println("[BIGIP] resourcePoolMemberCreate restURL : " + restURL)

	log.Println("[BIGIP] resourcePoolMemberCreate name : " + name)
	log.Println("[BIGIP] resourcePoolMemberCreate partition : " + partition)
	log.Println("[BIGIP] resourcePoolMemberCreate description : " + description)
	log.Println("[BIGIP] resourcePoolMemberCreate pool_id : " + pool_id)
	log.Println("[BIGIP] resourcePoolMemberCreate node_id : " + node_id)
	log.Println("[BIGIP] resourcePoolMemberCreate node_name : " + node_name)
	log.Println("[BIGIP] resourcePoolMemberCreate address : " + address)
	log.Println("[BIGIP] resourcePoolMemberCreate port : " + portStr)


	json := simplejson.New()
	// json.Set("name", name)
	json.Set("partition", partition)
	json.Set("address", address)
	json.Set("name", name)
	json.Set("description", description)

	jsonPostString := JSONtoString(json)
	jsonPostString = "{" + jsonPostString + "}"

	log.Println("[BIGIP] resourcePoolMemberCreate jsonPostString : " + jsonPostString)

	jsonRet, err := F5Post(restURL, jsonPostString, *client)
	if err != nil {
		return err
	}

	log.Println("[BIGIP] resourcePoolMemberCreate destID : " + jsonRet.MustString() )
	d.SetId(destID)

	return resourcePoolMemberRead(d,m)

	return errors.New("NYI")
}


func resourcePoolMemberRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourcePoolMemberRead")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_pool
	log.Println("[BIGIP] resourcePoolMemberRead restURL : " + restURL)

	pool_id := d.Get("pool_id").(string)
	node_id := d.Get("node_id").(string)
	node_name := d.Get("node_name").(string)
	description := d.Get("description").(string)
	partition := d.Get("partition").(string)
	address := d.Get("address").(string)
	port := d.Get("port").(int)
	portStr := strconv.Itoa(port)
	name := d.Get("name").(string)
	destID := d.Id()

	restURL = restURL + "/" + pool_id + "/members/" + destID

	log.Println("[BIGIP] resourcePoolMemberRead restURL : " + restURL)

	log.Println("[BIGIP] resourcePoolMemberRead name : " + name)
	log.Println("[BIGIP] resourcePoolMemberRead partition : " + partition)
	log.Println("[BIGIP] resourcePoolMemberRead description : " + description)
	log.Println("[BIGIP] resourcePoolMemberRead pool_id : " + pool_id)
	log.Println("[BIGIP] resourcePoolMemberRead node_id : " + node_id)
	log.Println("[BIGIP] resourcePoolMemberRead node_name : " + node_name)
	log.Println("[BIGIP] resourcePoolMemberRead address : " + address)
	log.Println("[BIGIP] resourcePoolMemberRead port : " + portStr)




	myJson, f5err := F5Get(restURL, *client)
	if f5err != nil {
		log.Println("[BIGIP] resourcePoolMemberRead if f5err != nil : ")
		if strings.Contains( f5err.Error() , "404") {
			log.Println("[BIGIP] resourcePoolMemberRead NODE WAS DELETED : " + destID)
			d.SetId("")
			return nil
		}
		fmt.Println(f5err)
		return f5err
	}
	log.Println("[BIGIP] resourcePoolMemberRead JSON : " + myJson.MustString() )

	/*
	{
	kind: "tm:ltm:pool:members:membersstate",
	name: "tfTestNode:80",
	partition: "Website",
	fullPath: "/Website/tfTestNode:80",
	generation: 63102316,
	selfLink: "https://localhost/mgmt/tm/ltm/pool/~Website~tfTestPool/members/~Website~tfTestNode:80?ver=11.5.1",
	address: "1.1.1.1",
	connectionLimit: 0,
	description: "test pool member!",
	dynamicRatio: 1,
	inheritProfile: "enabled",
	logging: "disabled",
	monitor: "default",
	priorityGroup: 0,
	rateLimit: "disabled",
	ratio: 1,
	session: "user-enabled",
	state: "unchecked"
	}
	*/



	resKind := myJson.Get("kind").MustString("")
	resName := myJson.Get("name").MustString("")
	resPartition := myJson.Get("partition").MustString("")
	resFullPath := myJson.Get("fullPath").MustString("")
	resGeneration := myJson.Get("generation").MustString("")
	resSelfLink := myJson.Get("selfLink").MustString("")
	resDescription := myJson.Get("description").MustString("")
	resAddress := myJson.Get("address").MustString("")


	log.Println("[BIGIP] resourcePoolMemberRead resKind : " + resKind)
	log.Println("[BIGIP] resourcePoolMemberRead resName : " + resName)
	log.Println("[BIGIP] resourcePoolMemberRead resPartition : " + resPartition)
	log.Println("[BIGIP] resourcePoolMemberRead resFullPath : " + resFullPath)
	log.Println("[BIGIP] resourcePoolMemberRead resGeneration : " + resGeneration)
	log.Println("[BIGIP] resourcePoolMemberRead resSelfLink : " + resSelfLink)
	log.Println("[BIGIP] resourcePoolMemberRead resDescription : " + resDescription)
	log.Println("[BIGIP] resourcePoolMemberRead resAddress : " + resAddress)

	if resAddress == "" {
		// no address means no pool member!
		d.SetId("")
		return nil
	}

	d.Set("description",resDescription)
	d.Set("address",resAddress)
	d.Set("name",resName)

	return nil
	return errors.New("NYI")
}


func resourcePoolMemberUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourcePoolMemberUpdate")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_pool
	log.Println("[BIGIP] resourcePoolMemberUpdate restURL : " + restURL)

	pool_id := d.Get("pool_id").(string)
	destID := d.Id()

	restURL = restURL + "/" + pool_id + "/members/" + destID

	log.Println("[BIGIP] resourcePoolMemberUpdate restURL : " + restURL)

	// return errors.New("NYI")


	// first delete the resource
	rdelerr := resourcePoolMemberDelete(d, m)
	if rdelerr != nil {
		return rdelerr
	}
	rcreateerr := resourcePoolMemberCreate(d, m)
	if rcreateerr != nil {
		return rcreateerr
	}


	return nil
}


func resourcePoolMemberDelete(d *schema.ResourceData, m interface{}) error {

	log.Println("[BIGIP] resourcePoolMemberDelete")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_pool
	log.Println("[BIGIP] resourcePoolMemberDelete restURL : " + restURL)

	pool_id := d.Get("pool_id").(string)
	node_id := d.Get("node_id").(string)
	node_name := d.Get("node_name").(string)
	description := d.Get("description").(string)
	partition := d.Get("partition").(string)
	address := d.Get("address").(string)
	port := d.Get("port").(int)
	portStr := strconv.Itoa(port)
	name := d.Get("name").(string)
	destID := d.Id()

	restURL = restURL + "/" + pool_id + "/members/" + destID

	log.Println("[BIGIP] resourcePoolMemberDelete restURL : " + restURL)

	log.Println("[BIGIP] resourcePoolMemberDelete name : " + name)
	log.Println("[BIGIP] resourcePoolMemberDelete partition : " + partition)
	log.Println("[BIGIP] resourcePoolMemberDelete description : " + description)
	log.Println("[BIGIP] resourcePoolMemberDelete pool_id : " + pool_id)
	log.Println("[BIGIP] resourcePoolMemberDelete node_id : " + node_id)
	log.Println("[BIGIP] resourcePoolMemberDelete node_name : " + node_name)
	log.Println("[BIGIP] resourcePoolMemberDelete address : " + address)
	log.Println("[BIGIP] resourcePoolMemberDelete port : " + portStr)



	// return errors.New("NYI")

	_, f5err := F5Delete(restURL, *client)
	
	if f5err != nil {
		log.Println("[BIGIP] resourcePoolMemberDelete if f5err != nil : ")
		if strings.Contains( f5err.Error() , "404") {
			log.Println("[BIGIP] resourcePoolMemberDelete NODE NOT FOUND : " + destID)
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




