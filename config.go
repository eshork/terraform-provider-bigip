package main

import (
	"log"
	"net/url"
	// "fmt"
	"errors"
	"strings"
	// "github.com/parnurzeal/gorequest"
	// "crypto/tls"
	// "github.com/bitly/go-simplejson"
)


type Config struct {
	RestUsername  string
	RestPassword  string
	RestIP        string
}


type BIGIPClient struct {
	RestUsername  string
	RestPassword  string
	RestIP        string
	URLBase       string

	ModuleURL_ltm string
	ModuleURL_persistence string
	ModuleURL_monitor string
	ModuleURL_defaultnodemonitor string
	ModuleURL_profile string
	ModuleURL_node string
	ModuleURL_pool string
	ModuleURL_virtual string
	ModuleURL_virtualaddress string
}

func (c *BIGIPClient) HasModule_persistence() bool {
	return c.ModuleURL_persistence != ""
}
func (c *BIGIPClient) HasModule_monitor() bool {
	return c.ModuleURL_monitor != ""
}
func (c *BIGIPClient) HasModule_defaultnodemonitor() bool {
	return c.ModuleURL_monitor != ""
}
func (c *BIGIPClient) HasModule_profile() bool {
	return c.ModuleURL_profile != ""
}
func (c *BIGIPClient) HasModule_node() bool {
	return c.ModuleURL_node != ""
}
func (c *BIGIPClient) HasModule_pool() bool {
	return c.ModuleURL_pool != ""
}
func (c *BIGIPClient) HasModule_virtual() bool {
	return c.ModuleURL_virtual != ""
}
func (c *BIGIPClient) HasModule_virtualaddress() bool {
	return c.ModuleURL_virtualaddress != ""
}



func (c *Config) Client() (interface{}, error) {
	log.Println("[BIGIP] Client() ")

	var client BIGIPClient
	client.RestUsername = c.RestUsername
	client.RestPassword = c.RestPassword
	client.RestIP = c.RestIP
	client.URLBase = "https://" + c.RestIP + "/mgmt/tm/ltm"


	// Get URLs for all the modules - also validates we have working credentials
	baseurl := client.URLBase
	myJson, f5err := F5Get(baseurl, client)
	
	if f5err != nil {
		return &client, f5err
	} else {
		resultType := myJson.Get("kind").MustString("")
		selfLink := myJson.Get("selfLink").MustString("")
		if resultType == "tm:ltm:ltmcollectionstate" {
			if selfLink != "" {
				tUrl, urlErr := url.Parse(selfLink)
				if urlErr != nil { return &client, urlErr }
				client.ModuleURL_ltm = tUrl.Path

				jsonItems := myJson.Get("items")

				for i := 0 ; true ; i++ {
					refLink := jsonItems.GetIndex(i).Get("reference").Get("link").MustString("")
					if refLink == "" {log.Println("[BIGIP] End!") ;  break }
					//
					tUrl, urlErr := url.Parse(refLink)
					if urlErr == nil {
						modulepath := tUrl.Path
						log.Println("[BIGIP] Module path: " + modulepath)
						switch true {
							case strings.Contains(modulepath, "/ltm/monitor"):
								log.Println("[BIGIP] Module path (monitor): " + modulepath)
								client.ModuleURL_monitor = modulepath
							case strings.Contains(modulepath, "/ltm/persistence"):
								log.Println("[BIGIP] Module path (persistence): " + modulepath)
								client.ModuleURL_persistence = modulepath
							case strings.Contains(modulepath, "/ltm/default-node-monitor"):
								log.Println("[BIGIP] Module path (defaultnodemonitor): " + modulepath)
								client.ModuleURL_defaultnodemonitor = modulepath
							case strings.Contains(modulepath, "/ltm/profile"):
								log.Println("[BIGIP] Module path (profile): " + modulepath)
								client.ModuleURL_profile = modulepath
							case strings.Contains(modulepath, "/ltm/node"):
								log.Println("[BIGIP] Module path (node): " + modulepath)
								client.ModuleURL_node = modulepath
							case strings.Contains(modulepath, "/ltm/pool"):
								log.Println("[BIGIP] Module path (pool): " + modulepath)
								client.ModuleURL_pool = modulepath
							case strings.Contains(modulepath, "/ltm/virtual-address"):
								log.Println("[BIGIP] Module path (virtualaddress): " + modulepath)
								client.ModuleURL_virtualaddress = modulepath
							case strings.Contains(modulepath, "/ltm/virtual"):
								log.Println("[BIGIP] Module path (virtual): " + modulepath)
								client.ModuleURL_virtual = modulepath
						}
					}
				}
			}
		} else {
			return &client, errors.New("bad response from device")
		}
	}

	if client.HasModule_persistence() { log.Println("[BIGIP] Module FOUND): " + "persistence") }
	if client.HasModule_monitor() { log.Println("[BIGIP] Module FOUND): " + "monitor") }
	if client.HasModule_defaultnodemonitor() { log.Println("[BIGIP] Module FOUND): " + "defaultnodemonitor") }
	if client.HasModule_profile() { log.Println("[BIGIP] Module FOUND): " + "profile") }
	if client.HasModule_node() { log.Println("[BIGIP] Module FOUND): " + "node") }
	if client.HasModule_pool() { log.Println("[BIGIP] Module FOUND): " + "pool") }
	if client.HasModule_virtual() { log.Println("[BIGIP] Module FOUND): " + "virtual") }
	if client.HasModule_virtualaddress() { log.Println("[BIGIP] Module FOUND): " + "virtualaddress") }


	return &client, nil
}
