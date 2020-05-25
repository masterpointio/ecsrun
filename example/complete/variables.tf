variable "stage" {
  type        = string
  description = "The environment that this infrastrcuture is being deployed to e.g. dev, stage, or prod"
}

variable "namespace" {
  type        = string
  description = "Namespace, which could be your organization name or abbreviation, e.g. 'eg' or 'cp'"
}

variable "name" {
  type        = string
  description = "The name for the ECR Repo. This will be added to the label to create the final repo name."
}

variable "attributes" {
  type        = list(string)
  default     = []
  description = "Additional attributes (e.g. `1`)"
}

variable "region" {
  default     = "us-east-1"
  type        = string
  description = "The AWS Region to deploy these resources to."
}

variable "availability_zones" {
  default     = ["us-east-1a"]
  type        = list(string)
  description = "List of Availability Zones where subnets will be created"
}

variable "nat_gateway_enabled" {
  default     = false
  type        = bool
  description = "Whether to enable NAT Gateways. If false, then the application uses NAT Instances, which are much cheaper."
}
