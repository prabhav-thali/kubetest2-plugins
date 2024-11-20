#
# (C) Copyright IBM Corp. 2022.
# SPDX-License-Identifier: Apache-2.0
#

data "ibm_is_vpc" "vpc" {
  count = var.vpc_name == "" ? 0 : 1
  name  = var.vpc_name
}

data "ibm_is_subnet" "subnet" {
  count = var.subnet_name == "" ? 0 : 1
  name  = var.subnet_name
}

data "ibm_resource_group" "default_group" {
  name = var.resource_group
}

module "vpc" {
  # Create new vpc ans subnet only if vpc_name is not set
  count        = var.vpc_name == "" ? 1 : 0
  source       = "./vpc"
  cluster_name = var.cluster_name
  zone         = var.zone
  resource_group = data.ibm_resource_group.default_group.id
}

locals {
  vpc_id            = var.vpc_name == "" ? module.vpc[0].vpc_id : data.ibm_is_vpc.vpc[0].id
  subnet_id         = var.vpc_name == "" ? module.vpc[0].subnet_id : data.ibm_is_subnet.subnet[0].id
  security_group_id = var.vpc_name == "" ? module.vpc[0].security_group_id : data.ibm_is_vpc.vpc[0].default_security_group
}

data "ibm_is_image" "node_image" {
  name = var.node_image
}

resource "ibm_is_ssh_key" "created_ssh_key" {
  # Create the ssh key only if the public key is set
  count      = var.ssh_pub_key == "" ? 0 : 1
  name       = var.ssh_key_name
  public_key = var.ssh_pub_key
  resource_group = data.ibm_resource_group.default_group
}

data "ibm_is_ssh_key" "ssh_key" {
  # Wait if the key needs creating first
  depends_on = [ibm_is_ssh_key.created_ssh_key]
  name       = var.ssh_key_name
}

resource "ibm_is_instance_template" "node_template" {
  name    = "${var.cluster_name}-node-template"
  image   = data.ibm_is_image.node_image.id
  profile = var.node_profile
  vpc     = local.vpc_id
  zone    = var.zone
  resource_group = data.ibm_resource_group.default_group.id
  keys    = [data.ibm_is_ssh_key.ssh_key.id]

  primary_network_interface {
    subnet          = local.subnet_id
    security_groups = [local.security_group_id]
  }
}

module "master" {
  source                    = "./node"
  node_name                 = "${var.cluster_name}-master"
  node_instance_template_id = ibm_is_instance_template.node_template.id
  resource_group = data.ibm_resource_group.default_group.id
}

module "workers" {
  source                    = "./node"
  count                     = var.workers_count
  node_name                 = "${var.cluster_name}-worker-${count.index}"
  node_instance_template_id = ibm_is_instance_template.node_template.id
  resource_group = data.ibm_resource_group.default_group.id
}

resource "null_resource" "wait-for-master-completes" {
  connection {
    type = "ssh"
    user = "root"
    host = module.master.public_ip
    private_key = file(var.ssh_private_key)
    timeout = "20m"
  }
  provisioner "remote-exec" {
    inline = [
      "cloud-init status -w"
    ]
  }
}

resource "null_resource" "wait-for-workers-completes" {
  count = var.workers_count
  connection {
    type = "ssh"
    user = "root"
    host = module.workers[count.index].public_ip
    private_key = file(var.ssh_private_key)
    timeout = "15m"
  }
  provisioner "remote-exec" {
    inline = [
      "cloud-init status -w"
    ]
  }
}

#resource "local_file" "inventory" {
#  content = templatefile("${path.module}/k8s-ansible/inventory.tmpl",
#    {
#      ip_addrs = [for k, w in module.nodes : w.public_ip],
#    }
#  )
#  filename = "${path.module}/ansible/inventory"
#}


#resource "null_resource" "ansible" {
#  triggers = {
#    inventory = resource.local_file.inventory.content
#  }
#  provisioner "local-exec" {
#    working_dir = "./ansible"
#    command     = "ansible-playbook -i inventory -u root containerd.yaml -e containerd_release_version=${var.containerd_version} -e kube_version=${var.kube_version}"
#  }
#}

#resource "null_resource" "kubeadm" {
#  depends_on = [
#    null_resource.ansible
#  ]
#  provisioner "local-exec" {
#    command = "./kube-init.sh ${module.nodes[0].private_ip}"
#  }
#}

#resource "null_resource" "label_nodes" {
#  depends_on = [
#    null_resource.kubeadm
#  ]
#  provisioner "local-exec" {
#    command = "./label-nodes.sh ${var.region} ${var.zone} ${local.subnet_id}"
#  }
#}
