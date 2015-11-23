package main

import (
	"log"
	"fmt"
	"reflect"
	"strconv"
	"errors"
	"strings"
	"github.com/parnurzeal/gorequest"
	"crypto/tls"
	"github.com/bitly/go-simplejson"
)

func F5Delete(url string, clientconfig BIGIPClient) (*simplejson.Json, error) {
	log.Println("[BIGIP] F5Delete:  " + url)
	emptyJson := simplejson.New()

	username := clientconfig.RestUsername
	password := clientconfig.RestPassword
	request := gorequest.New().SetBasicAuth(username, password).TLSClientConfig(&tls.Config{ InsecureSkipVerify: true})
	resp, body, err := request.Delete(url).End()
	if err != nil {
		log.Println("[BIGIP] err: something broke ")
		fmt.Println(err)
		return emptyJson, err[0]
	} else {
		if resp.Status != "200 OK" {
			//
			return emptyJson, errors.New("Did not get 200 OK - got " + resp.Status)
		} else {
			log.Println("[BIGIP] RESP: " + resp.Status)
			log.Println("[BIGIP] RESP: " + body)
			return emptyJson, nil
		}
	}

	return emptyJson, errors.New("unknown error")
}


func F5Get(url string, clientconfig BIGIPClient) (*simplejson.Json, error) {
	log.Println("[BIGIP] F5Get:  " + url)
	emptyJson := simplejson.New()

	username := clientconfig.RestUsername
	password := clientconfig.RestPassword
	request := gorequest.New().SetBasicAuth(username, password).TLSClientConfig(&tls.Config{ InsecureSkipVerify: true})
	resp, body, err := request.Get(url).End()
	if err != nil {
		log.Println("[BIGIP] err: something broke ")
		fmt.Println(err)
		return emptyJson, err[0]
	} else {
		if resp.Status != "200 OK" {
			//
			return emptyJson, errors.New("Did not get 200 OK - got " + resp.Status)
		} else {
			log.Println("[BIGIP] RESP: " + resp.Status)
			log.Println("[BIGIP] RESP: " + body)

			myJson, jsonErr := simplejson.NewFromReader( strings.NewReader(body) )
			if jsonErr != nil {
				log.Println("[BIGIP] jsonErr: something broke ")
				fmt.Println(jsonErr)
				return emptyJson, jsonErr
			} else {
				return myJson, nil
				//fmt.Println( myJson.Get("kind") )
			}
		}
	}

	return emptyJson, errors.New("unknown error")
}

func F5Put(url string, jsonPutString string, clientconfig BIGIPClient) (*simplejson.Json, error) {
	log.Println("[BIGIP] F5Put:  " + url)
	log.Println("[BIGIP] F5Put: " + jsonPutString)
	emptyJson := simplejson.New()

	username := clientconfig.RestUsername
	password := clientconfig.RestPassword
	request := gorequest.New().SetBasicAuth(username, password).TLSClientConfig(&tls.Config{ InsecureSkipVerify: true})
	resp, body, err := request.Put(url).
		Send(jsonPutString).End()


	if err != nil {
		log.Println("[BIGIP] err: something broke ")
		fmt.Println(err)
		return emptyJson, err[0]
	} else {
		log.Println("[BIGIP] " + resp.Status)
		log.Println(body)
		if resp.Status == "200 OK" {
			myJson, jsonErr := simplejson.NewFromReader( strings.NewReader(body) )
			if jsonErr != nil {
				return emptyJson, errors.New("JSON failed to parse: " + body)				
			} else {
				return myJson, nil
			}

		} else {
			myJson, jsonErr := simplejson.NewFromReader( strings.NewReader(body) )
			if jsonErr != nil {
				return emptyJson, errors.New("bad response from F5: " + resp.Status + " | " + body)				
			} else {
				msg := JSONtoString(myJson)
				return emptyJson, errors.New(msg)
			}
		}
	}

	return emptyJson, errors.New("F5Put: unknown error")
}

func F5Post(url string, jsonPostString string, clientconfig BIGIPClient) (*simplejson.Json, error) {
	log.Println("[BIGIP] F5Post:  " + url)
	log.Println("[BIGIP] F5Post: " + jsonPostString)
	emptyJson := simplejson.New()

	username := clientconfig.RestUsername
	password := clientconfig.RestPassword
	request := gorequest.New().SetBasicAuth(username, password).TLSClientConfig(&tls.Config{ InsecureSkipVerify: true})
	resp, body, err := request.Post(url).
		Send(jsonPostString).End()
	if err != nil {
		log.Println("[BIGIP] err: something broke ")
		fmt.Println(err)
		return emptyJson, err[0]
	} else {
		log.Println("[BIGIP] " + resp.Status)
		log.Println(body)
		if resp.Status == "200 OK" {
			myJson, jsonErr := simplejson.NewFromReader( strings.NewReader(body) )
			if jsonErr != nil {
				return emptyJson, errors.New("JSON failed to parse: " + body)				
			} else {
				return myJson, nil
			}

		} else {
			myJson, jsonErr := simplejson.NewFromReader( strings.NewReader(body) )
			if jsonErr != nil {
				return emptyJson, errors.New("bad response from F5: " + resp.Status + " | " + body)				
			} else {
				msg := JSONtoString(myJson)
				return emptyJson, errors.New(msg)
			}
		}
	}

	return emptyJson, errors.New("unknown error")
}


func JSONtoString( j *simplejson.Json ) string {
	log.Println("[BIGIP] JSONtoStringJSONtoStringJSONtoStringJSONtoStringJSONtoString ")
	var outstr string = ""

	mymap, err := j.Map()
	if err != nil {
		fmt.Println(j.Interface())
		fmt.Println(mymap)
		return "no map"
	}
	
	var needComma bool = false
	for k,v := range mymap {
		fmt.Println(k)
		fmt.Println(v)

		if needComma {
			outstr += ",\"" + k + "\":"
		} else {
			outstr += "\"" + k + "\":"
			needComma = true
		}

		if (reflect.ValueOf(v).Kind() == reflect.Map) {
			fmt.Println("has map!")
			fmt.Println(v)
				outstr += "{"
				next := j.Get(k)
				outstr += JSONtoString(next)
				outstr += "}"
		} else if (reflect.ValueOf(v).Kind() == reflect.Slice) || (reflect.ValueOf(v).Kind() == reflect.Array) {
			fmt.Println("has slice!")
			
			outstr += "["
			slen := reflect.ValueOf(v).Len()
			var inComma bool = false
			for i := 0; i < slen ; i++ {
				if inComma {
					outstr += ","
				} else {
					inComma = true
				}
				//
				tval := reflect.ValueOf(v).Index(i)

				if (tval.Kind() == reflect.Map) {
				//if str, ok := v.(mymap); ok {
					fmt.Println("has map!")
					fmt.Println(v)
						outstr += "{"
						outstr += "}"
				} else {
					if tval.Kind() == reflect.String {
						outstr += "\"" + tval.String() + "\""
					} else
					if tval.Kind() == reflect.Int {
						outstr += "\"" + strconv.Itoa(int(tval.Int())) + "\""
					}
				}
			}
			outstr += "]"
		} else {
			if str, ok := v.(string); ok {
				outstr += "\"" + str + "\""
			} else
			if num, ok := v.(int); ok {
				outstr += "\"" + strconv.Itoa(num) + "\""
			}
			if mybool, ok := v.(bool); ok {
				if mybool {
					outstr += "true"
				} else {
					outstr += "false"
				}
			}
		}
		
	}

	return outstr
}



