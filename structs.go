package main



//Used for the API token
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//To hold the token for the API
type Token struct {
	Token string `json:"token"`
}

//Struct for managing the hostname, time remaining, and expiration date.
type VM struct{
	Hostname  string
	TimeLeft float64
	Expiration string
	Environment string
}


//Struct to house the json which contains the individual VM details.  All specific information from the individual VM
//From the detailed response from CloudBolt
type HostDetail struct {
	Links struct {
		Actions []struct {
			PowerOn struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"power_on,omitempty"`
			PowerOff struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"power_off,omitempty"`
			Reboot struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"reboot,omitempty"`
			RefreshInfo struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"refresh_info,omitempty"`
			Snapshot struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"snapshot,omitempty"`
			ForcePowerOffReset struct {
				Href  string `json:"href"`
				Title string `json:"title"`
			} `json:"Force power-off / reset,omitempty"`
		} `json:"actions"`
		Environment struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"environment"`
		Group struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"group"`
		History struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"history"`
		Jobs struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"jobs"`
		OsBuild struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"os-build"`
		Owner struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"owner"`
		ProvisionJob struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"provision-job"`
		ResourceHandler struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"resource-handler"`
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
	} `json:"_links"`
	CPUCnt      int `json:"cpu-cnt"`
	Credentials struct {
		Key      string `json:"key"`
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"credentials"`
	DateAddedToCloudbolt string `json:"date-added-to-cloudbolt"`
	DiskSize             string `json:"disk-size"`
	Disks                []struct {
		Datastore        string `json:"datastore"`
		DiskSize         int    `json:"disk-size"`
		Name             string `json:"name"`
		ProvisioningType string `json:"provisioning-type"`
		UUID             string `json:"uuid"`
	} `json:"disks"`
	Hostname string        `json:"hostname"`
	IP       string        `json:"ip"`
	Labels   []interface{} `json:"labels"`
	Mac      string        `json:"mac"`
	MemSize  string        `json:"mem-size"`
	Networks []struct {
		IP        string `json:"ip"`
		Mac       string `json:"mac"`
		Name      string `json:"name"`
		Network   string `json:"network"`
		PrivateIP string `json:"private-ip"`
	} `json:"networks"`
	Notes      string `json:"notes"`
	OsFamily   string `json:"os-family"`
	Parameters struct {
		CaaSDatacenter string `json:"CaaS-Datacenter"`
	} `json:"parameters"`
	PowerStatus   string `json:"power-status"`
	RateBreakdown struct {
		Extra    string `json:"extra"`
		Hardware string `json:"hardware"`
		Software string `json:"software"`
		Total    string `json:"total"`
	} `json:"rate-breakdown"`
	Status              string `json:"status"`
	TechSpecificDetails struct {
		VmwareCluster     string `json:"vmware-cluster"`
		VmwareLinkedClone bool   `json:"vmware-linked-clone"`
	} `json:"tech-specific-details"`
}

//Minor entries from Cloudbolt API
type HostLookup struct {
	Links struct {
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
		Next struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"next"`
	} `json:"_links"`
	Total    int `json:"total"`
	Count    int `json:"count"`
	Embedded []Hostref `json:"_embedded"`

}


type Hostref struct {
	Links struct {
		Self struct {
			Href  string `json:"href"`
			Title string `json:"title"`
		} `json:"self"`
	} `json:"_links"`
	Hostname string `json:"hostname"`
}
