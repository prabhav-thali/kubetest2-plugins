package vpc

type TFVars struct {
	Apikey        string  `json:"ibmcloud_api_key,omitempty"`
	DNSName       string  `json:"vpc_dns"`
	DNSZone       string  `json:"vpc_dns_zone"`
	SSHKey        string  `json:"ssh_key_name"`
	Region        string  `json:"region"`
	Zone          string  `json:"zone"`
	ResourceGroup string  `json:"resource_group"`
	NodeImageName string  `json:"node_image"`
	NodeProfile   string  `json:"node_profile"`
	ContVersion   string  `json:"cont_version"`
	KubeVersion   string  `json:"kube_version"`
	//DNSName       string  `json:"vpc_dns"`
	//DNSZone       string  `json:"vpc_dns_zone"`
	//ServiceID     string  `json:"vpc_service_id"`
	//NetworkName   string  `json:"vpc_network_name"`
	//Memory        float64 `json:"vpc_memory"`
	//Processors    float64 `json:"vpc_processors"`
}
