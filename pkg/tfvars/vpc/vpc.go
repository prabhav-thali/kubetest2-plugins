package vpc

type TFVars struct {
	Apikey        string  `json:"vsi_api_key,omitempty"`
	SSHKey        string  `json:"vsi_ssh_key"`
	ClusterName   string  `json:"vsi_cluster_name"`
	Region        string  `json:"vsi_region"`
	Zone          string  `json:"vsi_zone"`
	ResourceGroup string  `json:"vsi_resource_group"`
	NodeImageName string  `json:"vsi_image_name"`
	NodeProfile   string  `json:"vsi_node_profile"`
	Nodes         int	  `json:"vsi_nodes"`
	ContVersion   string  `json:"vsi_cont_version"`
	KubeVersion   string  `json:"vsi_kube_version"`
	//DNSName       string  `json:"vsi_dns"`
	//DNSZone       string  `json:"vsi_dns_zone"`
	//ServiceID     string  `json:"vsi_service_id"`
	//NetworkName   string  `json:"vsi_network_name"`
	//Memory        float64 `json:"vsi_memory"`
	//Processors    float64 `json:"vsi_processors"`
}
