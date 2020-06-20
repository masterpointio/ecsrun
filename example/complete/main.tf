provider "aws" {
  region = var.region
}

module "label" {
  source    = "git::https://github.com/cloudposse/terraform-null-label.git?ref=tags/0.16.0"
  namespace = var.namespace
  stage     = var.stage
  name      = var.name
}

module "task_label" {
  source     = "git::https://github.com/cloudposse/terraform-null-label.git?ref=tags/0.16.0"
  context    = module.label.context
  attributes = compact(concat(var.attributes, ["task"]))
}

module "exec_label" {
  source     = "git::https://github.com/cloudposse/terraform-null-label.git?ref=tags/0.16.0"
  context    = module.label.context
  attributes = compact(concat(var.attributes, ["exec"]))
}

module "vpc" {
  source     = "git::https://github.com/cloudposse/terraform-aws-vpc.git?ref=tags/0.10.0"
  namespace  = var.namespace
  stage      = var.stage
  name       = var.name
  cidr_block = "10.0.0.0/16"
}

module "subnets" {
  source               = "git::https://github.com/cloudposse/terraform-aws-dynamic-subnets.git?ref=tags/0.19.0"
  availability_zones   = var.availability_zones
  namespace            = var.namespace
  stage                = var.stage
  vpc_id               = module.vpc.vpc_id
  igw_id               = module.vpc.igw_id
  cidr_block           = module.vpc.vpc_cidr_block
  nat_gateway_enabled  = var.nat_gateway_enabled
  nat_instance_enabled = ! var.nat_gateway_enabled
}

module "container_definition" {
  source          = "git::https://github.com/cloudposse/terraform-aws-ecs-container-definition.git?ref=tags/0.23.0"
  container_name  = module.label.id
  container_image = "bashell/alpine-bash:latest"
  log_configuration = {
    logDriver = "awslogs"
    options = {
      awslogs-create-group  = true
      awslogs-group         = module.label.id,
      awslogs-region        = var.region,
      awslogs-stream-prefix = var.name
    }
    secretOptions = null
  }
}

data "aws_iam_policy_document" "ecs_task_exec" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_exec" {
  name               = module.exec_label.id
  assume_role_policy = data.aws_iam_policy_document.ecs_task_exec.json
  tags               = module.exec_label.tags
}

data "aws_iam_policy_document" "ecs_exec" {
  statement {
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "ecr:GetAuthorizationToken",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
      "logs:CreateLogStream",
      "logs:CreateLogGroup",
      "logs:PutLogEvents"
    ]
  }
}

resource "aws_iam_role_policy" "ecs_exec" {
  name   = module.exec_label.id
  policy = data.aws_iam_policy_document.ecs_exec.json
  role   = aws_iam_role.ecs_exec.id
}

data "aws_iam_policy_document" "ecs_task" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_task" {
  name               = module.task_label.id
  assume_role_policy = data.aws_iam_policy_document.ecs_task.json
  tags               = module.task_label.tags
}

resource "aws_ecs_task_definition" "default" {
  family                   = module.label.id
  container_definitions    = module.container_definition.json
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_exec.arn
  task_role_arn            = aws_iam_role.ecs_task.arn
  tags                     = module.label.tags
}

resource "aws_ecs_cluster" "default" {
  name = module.label.id
  tags = module.label.tags
}
