package main

import (
	"strconv"
	"fmt"
	"bytes"
	"os"
	"crypto/tls"
	"io/ioutil"
	"time"
	"net/http"
	"log"

)

//Reads the /tmp/dat1 file, which contains all agent hosts, with days left until expiration.
//Then posts that data to the slack_hook
func sendMessage(){
	var buffer bytes.Buffer

	f,_ := os.Open(data_file)
	buffer.ReadFrom(f)
	f.Close()

	//proxyStr := os.Getenv("HTTPS_PROXY")
	//proxyURL, err := url.Parse(proxyStr)
	//if err != nil {
	//	fmt.Println(err)
	//}
	fmt.Println("Begin transmitted curl string:")
	fmt.Println(buffer.String())
	fmt.Println(":End transmitted curl string")
	var jsonStr = buffer.Bytes()
	req, err := http.NewRequest("POST", *slackHook, bytes.NewReader(jsonStr))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//Proxy: http.ProxyURL(proxyURL),

	}
	client := &http.Client{Transport: tr}
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	resp, err := client.Do(req)
	if err != nil{
		log.Fatal(err)
	}
	os.Remove(data_file)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response Status: ", resp.Status)
	fmt.Println("response Headers: ", resp.Header)
	fmt.Println("response Body: ", string(body))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")


}

//Wrapper function that prepares the /tmp/dat1 file that holds the message, and writes the beginning and end
// syntax formatting for the slack post
func prepareMessage(sortedGroup [7][]VM){
	var buffer bytes.Buffer

	fmt.Println("Preparing message to post.")
	first := true
	file, err := os.OpenFile(data_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // For read access.
	if err != nil {
		log.Fatal(err)
	}
	buffer.WriteString("{\"channel\": " + *slackChannel +
		"\"username\":" + *slackBotName +
		"\"attachments\": [ ")
	file.WriteString(buffer.String())


	//Build each section of the message
	for i,group := range sortedGroup {
		if len(group) != 0 {
			postGroup(group, i, first)
			first = false
		}
	}

	file.WriteString("]}")



}


func postGroup(group []VM, i int, first bool){
	var buffer bytes.Buffer
	expirationString := ""
	color := ""
	subtitle := ""
	pretext := ""
	fallback := ""



	file, err := os.OpenFile(data_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	if len(group) != 0 {
		if i == 0 {
			expirationString = "30+ days until expiration."
			color = "good"
		} else {
			if i == 1 {
				expirationString = "20+ days until expiration"
			} else {
				if i == 2 {
					expirationString = "10+ days until expiration"
					color = "good"
				} else {
					if i == 3 {
						expirationString = "Less than 10 days until expiration."
						color = "warning"
					} else {
						if i == 4 {
							expirationString = "Overdue."
							color = "danger"
						} else {
							if i == 5 {
								expirationString = "10+ days overdue."
								color = "danger"
							} else {
								if i == 6 {
									expirationString = "20+ days overdue"
									color = "danger"
								}
							}
						}
					}
				}
			}
		}
		if first {

			//Create title for beginning of section - with the Date
			now := time.Now()
			pretext = fallback
			//add environment to post
			fmt.Println(group[0].Environment)
			if group[0].Environment == "CAAS-Test" {
				subtitle = "DEVELOPMENT"
			}
			if group[0].Environment == "CAAS-NP" {
				subtitle = "NON PRODUCTION"
			}
			if group[0].Environment == "CAAS-PR" {
				subtitle = "PRODUCTION"
			}

			buffer.WriteString(*postTitle)
			buffer.WriteString(now.Format(timeLayout))
			fallback :=  buffer.String()
			buffer.Reset()
			pretext = fallback


		} else {
			buffer.WriteString(",")
		}

		//Part of the json payload for the slack message
		buffer.WriteString("{\"fallback\": \"" + fallback + "\"," +
			"\"color\": \"" + color + "\"," +
			"\"pretext\": \"" + pretext + "\"," +
			"\"author_icon\": \":stars:\"," +
			"\"title\":\"" + subtitle + "\", \"text\": \"\",\"fields\": [")
		buffer.WriteString("{\"title\":\"" + expirationString + "\"")
		buffer.WriteString(",\"value\":\"")
		days := 0.0
		//the groups that are closer to the expiration date - they are printed individually, with more details.
		for _, vm := range group {
			buffer.WriteString(vm.Hostname)
			buffer.WriteString(" ")
			days = vm.TimeLeft
			if i == 4 {
				//convert VMs that are overdue to an expired message - with days measured and date
				days = days * -1
				if days < 1 {
					buffer.WriteString( "Expired today. ")
				} else {
					buffer.WriteString(strconv.FormatFloat(days, 'f', 2, 64))
					buffer.WriteString(" days past expiration date. ")
				}
				buffer.WriteString(vm.Expiration)
				buffer.WriteString("\\n")
			}
			if i == 3 {
				if days < 1 {
					buffer.WriteString( "Expiring today. ")

				} else {
					//VMs that are due within 10 days
					buffer.WriteString(strconv.FormatFloat(days, 'f', 2, 64))
					buffer.WriteString(" days until expiration date. ")
				}
				buffer.WriteString(vm.Expiration)
				buffer.WriteString("\\n")
			}

		}
		buffer.WriteString("\"}]}")
		fmt.Println(buffer.String())
		file.WriteString(buffer.String())
		buffer.Reset()
		file.Close()

	}
}
