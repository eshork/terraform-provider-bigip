package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/bitly/go-simplejson"
	"strconv"
	"strings"
	"log"
	"fmt"
	"errors"
)

// const (
//     TypeInvalid ValueType = iota
//     TypeBool
//     TypeInt
//     TypeFloat
//     TypeString
//     TypeList
//     TypeMap
//     TypeSet
// )

func resourceVServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVServerCreate,
		Read:   resourceVServerRead,
		Update: resourceVServerUpdate,
		Delete: resourceVServerDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"partition": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dest_port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"dest_mask": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "255.255.255.255",
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "0.0.0.0/0",
			},
			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "",
			},
			"protocol": &schema.Schema{ //"ipProtocol" = "any" or "tcp" or "udp"
				Type:     schema.TypeString,
				Optional: true,
				Default: "tcp",
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"snat_automap": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
		},
	}
}


func resourceVServerCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourceVServerCreate")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_virtual
	log.Println("[BIGIP] resourceVServerCreate restURL : " + restURL)



	name := d.Get("name").(string)
	description := d.Get("description").(string)
	partition := d.Get("partition").(string)
	dest_ip := d.Get("dest_ip").(string)
	dest_port := d.Get("dest_port").(int)
	dest_port_str := strconv.Itoa(dest_port)
	destination := dest_ip + ":" + dest_port_str
	source := d.Get("source").(string)
	pool := d.Get("pool").(string)
	dest_mask := d.Get("dest_mask").(string)
	protocol := d.Get("protocol").(string)
	enabled := d.Get("enabled").(bool)
	enabled_str := "true"
	if !enabled {
		enabled_str = "false"
	}
	snat_automap := d.Get("snat_automap").(bool)
	snat_automap_str := "true"
	if !snat_automap {
		snat_automap_str = "false"
	}

	destID := "~" + partition + "~" + name


	log.Println("[BIGIP] resourceVServerCreate destID : " + destID)
	log.Println("[BIGIP] resourceVServerCreate name : " + name)
	log.Println("[BIGIP] resourceVServerCreate description : " + description)
	log.Println("[BIGIP] resourceVServerCreate partition : " + partition)
	log.Println("[BIGIP] resourceVServerCreate dest_ip : " + dest_ip)
	log.Println("[BIGIP] resourceVServerCreate dest_port : " + dest_port_str)
	log.Println("[BIGIP] resourceVServerCreate destination : " + destination)
	log.Println("[BIGIP] resourceVServerCreate dest_mask : " + dest_mask)
	log.Println("[BIGIP] resourceVServerCreate source : " + source)
	log.Println("[BIGIP] resourceVServerCreate protocol : " + protocol)
	log.Println("[BIGIP] resourceVServerCreate enabled : " + enabled_str)
	log.Println("[BIGIP] resourceVServerCreate pool : " + pool)
	log.Println("[BIGIP] resourceVServerCreate snat_automap : " + snat_automap_str)


	
	json := simplejson.New()
	json.Set("name", name)
	if description != "" {
		json.Set("description", description)
	}
	json.Set("partition", partition)
	json.Set("destination", destination)
	json.Set("mask", dest_mask)
	json.Set("source", source)
	json.Set("ipProtocol", protocol)
	if enabled {
		json.Set("enabled", true)
	}else{
		json.Set("disabled", true)
	}

	if pool != "" {
		json.Set("pool", pool)
	}

	if snat_automap {
		json.SetPath([]string{"sourceAddressTranslation","type"}, "automap")
	}


	jsonPostString := JSONtoString(json)
	jsonPostString = "{" + jsonPostString + "}"

	log.Println("[BIGIP] resourceVServerCreate jsonPostString : " + jsonPostString)


	jsonRet, err := F5Post(restURL, jsonPostString, *client)
	if err != nil {
		return err
	}

	log.Println("[BIGIP] resourceVServerCreate JSON RET : " + jsonRet.MustString() )
	d.SetId(destID)

	return resourceVServerRead(d,m)
	return nil
	return errors.New("NYI")
}


func resourceVServerRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourceVServerRead")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_virtual
	log.Println("[BIGIP] resourceVServerRead restURL : " + restURL)
	destID := d.Id()
	restURL = restURL + "/" + destID
	log.Println("[BIGIP] resourceVServerRead restURL : " + restURL)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	partition := d.Get("partition").(string)
	dest_ip := d.Get("dest_ip").(string)
	dest_port := d.Get("dest_port").(int)
	dest_port_str := strconv.Itoa(dest_port)
	destination := dest_ip + ":" + dest_port_str
	source := d.Get("source").(string)
	pool := d.Get("pool").(string)
	dest_mask := d.Get("dest_mask").(string)
	protocol := d.Get("protocol").(string)
	enabled := d.Get("enabled").(bool)
	enabled_str := "true"
	if !enabled {
		enabled_str = "false"
	}
	snat_automap := d.Get("snat_automap").(bool)
	snat_automap_str := "true"
	if !snat_automap {
		snat_automap_str = "false"
	}


	log.Println("[BIGIP] resourceVServerRead destID : " + destID)
	log.Println("[BIGIP] resourceVServerRead name : " + name)
	log.Println("[BIGIP] resourceVServerRead description : " + description)
	log.Println("[BIGIP] resourceVServerRead partition : " + partition)
	log.Println("[BIGIP] resourceVServerRead dest_ip : " + dest_ip)
	log.Println("[BIGIP] resourceVServerRead dest_port : " + dest_port_str)
	log.Println("[BIGIP] resourceVServerRead destination : " + destination)
	log.Println("[BIGIP] resourceVServerRead dest_mask : " + dest_mask)
	log.Println("[BIGIP] resourceVServerRead source : " + source)
	log.Println("[BIGIP] resourceVServerRead protocol : " + protocol)
	log.Println("[BIGIP] resourceVServerRead enabled : " + enabled_str)
	log.Println("[BIGIP] resourceVServerRead pool : " + pool)
	log.Println("[BIGIP] resourceVServerRead snat_automap : " + snat_automap_str)


	myJson, f5err := F5Get(restURL, *client)
	if f5err != nil {
		log.Println("[BIGIP] resourceVServerRead if f5err != nil : ")
		if strings.Contains( f5err.Error() , "404") {
			log.Println("[BIGIP] resourceVServerRead NODE WAS DELETED : " + destID)
			d.SetId("")
			return nil
		}
		fmt.Println(f5err)
		return f5err
	}
	log.Println("[BIGIP] resourceVServerRead JSON : " + JSONtoString(myJson) )

/*
{
kind: "tm:ltm:virtual:virtualstate",
name: "tfTestVirtualServer",
partition: "Website",
fullPath: "/Website/tfTestVirtualServer",
generation: 63146816,
selfLink: "https://localhost/mgmt/tm/ltm/virtual/~Website~tfTestVirtualServer?ver=11.5.1",
addressStatus: "yes",
autoLasthop: "default",
cmpEnabled: "yes",
connectionLimit: 0,
description: "just a test virtual server",
destination: "/Website/1.1.1.1:80",
disabled: true,
gtmScore: 0,
ipProtocol: "tcp",
mask: "255.255.255.255",
mirror: "disabled",
mobileAppTunnel: "disabled",
nat64: "disabled",
pool: "/Website/tfTestPool",
rateLimit: "disabled",
rateLimitDstMask: 0,
rateLimitMode: "object",
rateLimitSrcMask: 0,
source: "0.0.0.0/0",
sourceAddressTranslation: {
type: "automap"
},
sourcePort: "preserve",
synCookieStatus: "not-activated",
translateAddress: "enabled",
translatePort: "enabled",
vlansDisabled: true,
vsIndex: 138,
policiesReference: {
link: "https://localhost/mgmt/tm/ltm/virtual/~Website~tfTestVirtualServer/policies?ver=11.5.1",
isSubcollection: true
},
profilesReference: {
link: "https://localhost/mgmt/tm/ltm/virtual/~Website~tfTestVirtualServer/profiles?ver=11.5.1",
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
	resSource := myJson.Get("source").MustString("")
	resMask := myJson.Get("mask").MustString("")

	resPool := myJson.Get("pool").MustString("")
	resProtocol := myJson.Get("ipProtocol").MustString("")

	var resEnabled bool
	resEnabled = true
	resEnabled_str := "true"
	if _, err := myJson.Get("enabled").Bool(); err != nil {
		resEnabled = false
		resEnabled_str = "false"
	}


	var resSNatAutomap bool = false
	var resSNatAutomap_str string = "false"
	resSNAT_map := myJson.Get("sourceAddressTranslation").MustMap()
	if resSNAT_map["type"] == "automap" {
		resSNatAutomap = true
		resSNatAutomap_str = "true"
	}
	resSNAT := fmt.Sprint(resSNAT_map)

	resDestination := myJson.Get("destination").MustString("")

	destStrPart, portIntPart := grokVServerDestination(resDestination)




	log.Println("[BIGIP] resourceVServerRead resKind : " + resKind)
	log.Println("[BIGIP] resourceVServerRead resName : " + resName)
	log.Println("[BIGIP] resourceVServerRead resPartition : " + resPartition)
	log.Println("[BIGIP] resourceVServerRead resFullPath : " + resFullPath)
	log.Println("[BIGIP] resourceVServerRead resGeneration : " + resGeneration)
	log.Println("[BIGIP] resourceVServerRead resSelfLink : " + resSelfLink)
	log.Println("[BIGIP] resourceVServerRead resDescription : " + resDescription)
	log.Println("[BIGIP] resourceVServerRead resSource : " + resSource)
	log.Println("[BIGIP] resourceVServerRead resMask : " + resMask)
	log.Println("[BIGIP] resourceVServerRead resPool : " + resPool)
	log.Println("[BIGIP] resourceVServerRead resProtocol : " + resProtocol)
	log.Println("[BIGIP] resourceVServerRead resEnabled : " + resEnabled_str)
	log.Println("[BIGIP] resourceVServerRead resSNAT : " + resSNAT)
	log.Println("[BIGIP] resourceVServerRead resNatAutomap : " + resSNatAutomap_str)
	log.Println("[BIGIP] resourceVServerRead resDestination : " + resDestination)
	log.Println("[BIGIP] resourceVServerRead destStrPart : " + destStrPart)
	log.Println("[BIGIP] resourceVServerRead portIntPart : " + strconv.Itoa(portIntPart))


	d.Set("description",resDescription)
	d.Set("name",resName)

	d.Set("partition",resPartition)
	d.Set("source",resSource)
	d.Set("dest_mask",resMask)
	d.Set("pool",resPool)
	d.Set("protocol",resProtocol)
	d.Set("enabled",resEnabled)
	d.Set("snat_automap",resSNatAutomap)

	d.Set("dest_ip",destStrPart)
	d.Set("dest_port",portIntPart)


	return nil
	return errors.New("NYI")
}


func resourceVServerUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourceVServerUpdate")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_virtual
	destID := d.Id()
	restURL = restURL + "/" + destID
	log.Println("[BIGIP] resourceVServerUpdate restURL : " + restURL)


	json := simplejson.New()

	var needsFullReCreate bool = false

	// try to figure out what parts need to be updated

	if d.HasChange("name") {
		log.Println("[BIGIP] resourceVServerUpdate : name CHANGED")
		json.Set("name", d.Get("name"))
		// THIS WILL NOT WORK - IS READ-ONLY IN WEB UI - NEEDS DESTROY/CREATE TO PERFORM
		needsFullReCreate = true
	}
	if d.HasChange("partition") {
		log.Println("[BIGIP] resourceVServerUpdate : partition CHANGED")
		json.Set("partition", d.Get("partition"))
		// THIS WILL NOT WORK - IS READ-ONLY IN WEB UI - NEEDS DESTROY/CREATE TO PERFORM
		needsFullReCreate = true
	}
	if d.HasChange("description") {
		log.Println("[BIGIP] resourceVServerUpdate : description CHANGED")
		json.Set("description", d.Get("description"))
	}
	if d.HasChange("source") {
		log.Println("[BIGIP] resourceVServerUpdate : source CHANGED")
		json.Set("source", d.Get("source"))
	}
	if d.HasChange("enabled") {
		log.Println("[BIGIP] resourceVServerUpdate : enabled CHANGED")
		newEnabled := d.Get("enabled").(bool)

		if newEnabled {
			json.Set("enabled", true)
		}else{
			json.Set("disabled", true)
		}
	}
	if d.HasChange("pool") {
		log.Println("[BIGIP] resourceVServerUpdate : pool CHANGED")
		json.Set("pool", d.Get("pool"))
	}
	if d.HasChange("protocol") {
		log.Println("[BIGIP] resourceVServerUpdate : protocol CHANGED")
		json.Set("ipProtocol", d.Get("protocol"))
	}
	if d.HasChange("dest_mask") {
		log.Println("[BIGIP] resourceVServerUpdate : dest_mask CHANGED")
		json.Set("mask", d.Get("dest_mask"))
	}

	if d.HasChange("snat_automap") {
		log.Println("[BIGIP] resourceVServerUpdate : snat_automap CHANGED")
		newSnatAutomap := d.Get("snat_automap").(bool)
		if newSnatAutomap {
			json.SetPath([]string{"sourceAddressTranslation","type"}, "automap")
		} else {
			log.Println("[BIGIP] resourceVServerUpdate : snat_automap er-me-gerd-no-love...")
			// json.SetPath([]string{"sourceAddressTranslation"}, "")
			json.SetPath([]string{"sourceAddressTranslation","type"}, "none")
		}
	}

	if d.HasChange("dest_ip") || d.HasChange("dest_port") {
		newDestIp := d.Get("dest_ip").(string)
		newDestPort := d.Get("dest_port").(int)
		newDestPortStr := strconv.Itoa(newDestPort)
		destination := newDestIp + ":" + newDestPortStr
		json.Set("destination", destination)
	}


// 


// dest_ip
// dest_port
// destination




	// do work, son
	if needsFullReCreate {
		log.Println("[BIGIP] resourceVServerUpdate : FULL RECREATE")
		rdelerr := resourceVServerDelete(d, m)
		if rdelerr != nil {
			return rdelerr
		}
		rcreateerr := resourceVServerCreate(d, m)
		if rcreateerr != nil {
			return rcreateerr
		}
	} else {
		log.Println("[BIGIP] resourceVServerUpdate : UPDATE OKAY")
		jsonPostString := JSONtoString(json)
		jsonPostString = "{" + jsonPostString + "}"
		log.Println("[BIGIP] resourceVServerUpdate jsonPutString : " + jsonPostString)

		// a simple update is okay this time
		jsonRet, err := F5Put(restURL, jsonPostString, *client)
		if err != nil {
			return err
		}
		log.Println("[BIGIP] resourceVServerUpdate JSON RET : " + JSONtoString(jsonRet) )

	}

	return nil
	return errors.New("NYI")
}


func resourceVServerDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[BIGIP] resourceVServerDelete")
	client := m.(*BIGIPClient)
	restURL := "https://" + client.RestIP + client.ModuleURL_virtual
	destID := d.Id()
	restURL = restURL + "/" + destID
	log.Println("[BIGIP] resourceVServerDelete restURL : " + restURL)


	_, f5err := F5Delete(restURL, *client)
	
	if f5err != nil {
		log.Println("[BIGIP] resourceVServerDelete if f5err != nil : ")
		if strings.Contains( f5err.Error() , "404") {
			log.Println("[BIGIP] resourceVServerDelete NODE NOT FOUND : " + destID)
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


// stupid helper function to undo F5 string/value embedding
func grokVServerDestination(destString string) (string, int) {
	// find last :
	i_colon := strings.LastIndex(destString, ":")
	if i_colon == -1 {
		return "", 0
	}
	port_part := destString[i_colon+1:len(destString)]
	rest_part := destString[0:i_colon]
	log.Println("[BIGIP] grokVServerDestination : port_part = " + port_part)
	log.Println("[BIGIP] grokVServerDestination : rest_part = " + rest_part)

	var portPartInt int
	var portPartErr error
	if portPartInt, portPartErr = strconv.Atoi(port_part); portPartErr != nil {
		return "", 0
	}
	log.Println("[BIGIP] grokVServerDestination : port_part (int) = " + strconv.Itoa(portPartInt))

	i_slash := strings.LastIndex(rest_part, "/")
	addr_part := destString[i_slash+1:len(rest_part)]
	log.Println("[BIGIP] grokVServerDestination : addr_part = " + addr_part)

	return addr_part, portPartInt
	return "", 0 // this doesn't happen - it isn't real - look away!
}





