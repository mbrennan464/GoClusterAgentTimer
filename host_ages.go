package main

import (
	"bytes"
	"fmt"
	"encoding/json"
	"log"
	"strconv"
	"time"
	"os"
	"gopkg.in/alecthomas/kingpin.v2"
)



const timeLayout = "2006-01-02T15:04:05"

//const postTitle = "CaaS Host Ages and Expiration "
const data_file = "/tmp/dat1"
//const slack_hook="https://hooks.slack.com/services/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" //CAAS-MONITORING

var (
	app = kingpin.New("caas-host-ages", "CaaS Host Ages").Author("THD Engineering & Operations")

	// Cloudbolt configuration
	cloudboltUrl      = app.Flag("cloudbolt-address", "The url for Cloudbolt").Envar("CB_ADDRESS").Default("https://server.homedepot.com").String()
	cbUser     = app.Flag("cloudbolt-user", "CloudBolt Auth username").Envar("CB_USER").Required().String()
	cbPassword = app.Flag("cloudbolt-password", "CloudBolt Auth password").Envar("CB_PASSWORD").Required().String()


	// slack configuration
	postTitle = app.Flag("slack-post-title", "Slack Post Title").Envar("SLACK_POST_TITLE").Default("CaaS Host Ages and Expiration ").String()
	slackBotName = app.Flag("slack-bot-name", "Slack Bot Name").Envar("SLACK_BOT_NAME").Default("ose-platform-bot").String()
	slackChannel = app.Flag("slack-channel", "Slack Post Channel").Envar("SLACK_CHANNEL").Default("#caas-monitoring").String()
	slackHook    = app.Flag("slack-hook", "Slack API Hook").Envar("SLACK_HOOK").Required().String()

)

func main(){

	// Parse app flags
	kingpin.MustParse(app.Parse(os.Args[1:]))
	// set logging format
	log.SetFlags(log.LstdFlags)
	log.SetOutput(os.Stdout)

	group := [2]int{20,22}
	for _,g := range group {

		var cluster []VM

		var sortedGroup [7][]VM
		cluster = searchGroup(g, cluster)
		sortedGroup = sortCluster(cluster)

		prepareMessage(sortedGroup)
		sendMessage()
	}
}

func searchGroup(group int, cluster []VM )[]VM {
	var buffer bytes.Buffer
	var bufferDet bytes.Buffer
	var cburl string
	var token string


	creds := Credentials{*cbUser,*cbPassword}
	b, err := json.Marshal(creds)
	if err != nil {
		fmt.Println(err)
	}
	cburl = *cloudboltUrl
	token = getToken(cburl, b)

	//Run through the list of VMs and calculate the time left:  Filtering by group number and power on
	fmt.Println("Preparing to check hosts....")
	buffer.Reset()
	buffer.WriteString(cburl)
	buffer.WriteString("/api/v2/servers/?filter=group:")
	buffer.WriteString(strconv.Itoa(group))
	buffer.WriteString(";power-status:POWERON")
	////Get the url/page for the details of each hostname, searching through servers by hostname
	body := httpGetSecure(buffer.String(), token)
	hostResults := new(HostLookup)
	err2 := json.Unmarshal(body, &hostResults)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(strconv.Itoa(hostResults.Total) + " hosts in group " + strconv.Itoa(group))
	page := 2
	for i :=0; i<=hostResults.Total; i++ {
		for _, vm := range hostResults.Embedded {
			fmt.Println(vm.Hostname + "  " + strconv.Itoa(i))
			i = i + 1
			bufferDet.Reset()
			bufferDet.WriteString(cburl)
			bufferDet.WriteString(vm.Links.Self.Href)
			bufferDet.WriteString("/")
			//Pull the VM details page from cloudbolt - which includes the dateadded field
			body = httpGetSecure(bufferDet.String(), token)
			hostdetails := new(HostDetail)
			err2 = json.Unmarshal(body, &hostdetails)
			if err2 != nil {
				log.Fatal(err2)
			}
			//Use the group name to determine the environment we are in
			env := hostdetails.Links.Group.Title
			t, _ := time.Parse(timeLayout, hostdetails.DateAddedToCloudbolt)
			dateExpiration := t.AddDate(0, 0, 60)
			//Build our VM array with the details we need for the report
			vm := VM{hostdetails.Hostname, converttime(hostdetails.DateAddedToCloudbolt), dateExpiration.Format(timeLayout),
				env}
			cluster = append(cluster,vm)
		}
		buffer.Reset()
		buffer.WriteString(cburl)
		buffer.WriteString("/api/v2/servers/?filter=group:")
		buffer.WriteString(strconv.Itoa(group))
		buffer.WriteString(";power-status:POWERON")
		buffer.WriteString("&page=")
		buffer.WriteString(strconv.Itoa(page))
		body = httpGetSecure(buffer.String(), token)
		hostResults = new(HostLookup)
		err2 = json.Unmarshal(body, &hostResults)
		if err2 != nil {
			log.Fatal(err2)
		}
		page = page + 1
	}

	return cluster
	
}








//Sorts the groups so as to print the report in attachments
func sortCluster(cluster []VM)([7][]VM){
	var sortedGroup [7][]VM

	fmt.Println("Sorting VMs into age groups.")

	for _,member := range cluster{
		if member.TimeLeft > 30 {
			sortedGroup[0] = append(sortedGroup[0],member)
		} else {
			if member.TimeLeft > 20 {
				sortedGroup[1] = append(sortedGroup[1],member)
			}else {
				if member.TimeLeft > 10 {
					sortedGroup[2] = append(sortedGroup[2],member)
				}else {
					if member.TimeLeft > 0 {
						sortedGroup[3] = append(sortedGroup[3],member)
					} else {
						if member.TimeLeft > -10 {
							sortedGroup[4] = append(sortedGroup[4],member)
						} else {
							if member.TimeLeft > -20 {
								sortedGroup[5] = append(sortedGroup[5],member)
							} else {
								sortedGroup[6] = append(sortedGroup[6],member)
							}

						}
					}
				}
			}
		}

	}

	return sortedGroup
}
