output "private_subnet_ids" {
  value = module.subnets.private_subnet_ids
}

output "default_security_group_id" {
  value = module.vpc.vpc_default_security_group_id
}

output "cluster_name" {
  value = module.label.id
}

output "task_definition_name" {
  value = module.label.id
}
