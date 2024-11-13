output "vpc_id" { value = local.vpc_id }
output "ssh_key_id" { value = data.ibm_is_ssh_key.ssh_key.id }
output "subnet_id" { value = local.subnet_id }
output "security_group_id" { value = local.security_group_id }
output "region" { value = var.region }
output "zone" { value = var.zone }
output "resource_group_id" { value = data.ibm_resource_group.default_group.id }
output "masters" {
  value = module.master.public_ips[*][0]
  description = "k8s master node IP addresses"
}

output "workers" {
  value = module.workers.public_ips[*][0]
  description = "k8s worker node IP addresses"
}

output "masters_private" {
  value = module.master.private_ips[*][0]
  description = "k8s master nodes private IP addresses"
}

output "workers_private" {
  value = module.workers.private_ips[*][0]
  description = "k8s worker nodes private IP addresses"
}