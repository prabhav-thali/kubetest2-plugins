package vpc

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/prabhav-thali/kubetest2-plugins/pkg/providers"
	"github.com/prabhav-thali/kubetest2-plugins/pkg/tfvars/vpc"
	"github.com/spf13/pflag"
)

const (
	Name = "vpc"
)

var _ providers.Provider = &Provider{}

var VPCProvider = &Provider{}

type Provider struct {
	vpc.TFVars
}

func (p *Provider) Initialize() error {
	return nil
}

func (p *Provider) BindFlags(flags *pflag.FlagSet) {
	flags.StringVar(
		&p.Apikey, "vpc-api-key", "", "IBM Cloud API Key used for accessing the APIs",
	)
	flags.StringVar(
		&p.SSHKey, "vpc-ssh-key", "", "VPC SSH Key to authenticate LPARs",
	)
	flags.StringVar(
		&p.Region, "vpc-region", "", "IBM Cloud VPC region name",
	)
	flags.StringVar(
		&p.Zone, "vpc-zone", "", "IBM Cloud VPC zone name",
	)
	flags.StringVar(
		&p.ResourceGroup, "vpc-resource-group", "Default", "IBM Cloud resource group name(command: ibmcloud resource groups)",
	)
	flags.StringVar(
		&p.NodeImageName, "vpc-node-image-name", "", "Image ID(command: ibmcloud vsi imgs)",
	)
	flags.StringVar(
		&p.NodeProfile, "vpc-node-profile", "", "Image ID(command: ibmcloud vsi profiles)",
	)
	flags.IntVar(
		&p.Nodes, "vpc-nodes", 3, "Numbers of workers in the k8s cluster",
	)
	flags.StringVar(
		&p.KubeVersion, "vpc-kube-version", "", "Image ID(command: ibmcloud kubernetes version)",
	)
	flags.StringVar(
		&p.ContVersion, "vpc-cont-version", "", "Image ID(command: ibmcloud containerd version)",
	)
}

func (p *Provider) DumpConfig(dir string) error {
	filename := path.Join(dir, Name+".auto.tfvars.json")
	config, err := json.MarshalIndent(p.TFVars, "", "  ")
	if err != nil {
		return fmt.Errorf("errored file converting config to json: %v", err)
	}
	err = os.WriteFile(filename, config, 0644)
	if err != nil {
		return fmt.Errorf("failed to dump the json config to: %s, err: %v", filename, err)
	}
	return nil
}
