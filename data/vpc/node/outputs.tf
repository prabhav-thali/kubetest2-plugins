#
# (C) Copyright IBM Corp. 2023.
# SPDX-License-Identifier: Apache-2.0
#

output "public_ips" {
  value = ibm_is_floating_ip.node.address
}
output "private_ips" {
  value = ibm_is_instance.node.primary_network_interface.0.primary_ip.0.address
}
