terraform {
  required_version = ">= 1.3.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
  default_tags {
    tags = {
      Project     = var.project_name
      Environment = "Production"
      ManagedBy   = "Terraform"
    }
  }
}

data "aws_availability_zones" "available" {}

# ------------------------------------------------------------------------------
# VPC (Network Foundation)
# ------------------------------------------------------------------------------
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = "${var.project_name}-vpc"
  cidr = var.vpc_cidr

  azs             = slice(data.aws_availability_zones.available.names, 0, 3)
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true
}

# ------------------------------------------------------------------------------
# AWS SQS Queue
# ------------------------------------------------------------------------------
resource "aws_sqs_queue" "anomaly_queue" {
  name                       = "${var.project_name}-queue"
  visibility_timeout_seconds = 30
  message_retention_seconds  = 345600 # 4 days
}

# ------------------------------------------------------------------------------
# ElastiCache Redis (Sliding Window State)
# ------------------------------------------------------------------------------
resource "aws_elasticache_subnet_group" "redis_subnet" {
  name       = "${var.project_name}-redis-subnet"
  subnet_ids = module.vpc.private_subnets
}

resource "aws_elasticache_replication_group" "redis" {
  replication_group_id       = "${var.project_name}-redis"
  description                = "Redis cluster for anomaly rate limiting"
  node_type                  = "cache.t4g.micro"
  port                       = 6379
  parameter_group_name       = "default.redis7"
  subnet_group_name          = aws_elasticache_subnet_group.redis_subnet.name
  security_group_ids         = [aws_security_group.redis_sg.id]
  automatic_failover_enabled = true
  num_cache_clusters         = 2 # Primary + 1 Replica
}

resource "aws_security_group" "redis_sg" {
  name_prefix = "${var.project_name}-redis-sg"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block] # Only allow access from within the VPC
  }
}

# ------------------------------------------------------------------------------
# Aurora PostgreSQL (Anomaly Storage)
# ------------------------------------------------------------------------------
module "aurora" {
  source  = "terraform-aws-modules/rds-aurora/aws"
  version = "~> 8.0"

  name           = "${var.project_name}-db"
  engine         = "aurora-postgresql"
  engine_version = "15.3"
  
  # Using Serverless v2 to auto-scale based on threat spikes
  serverlessv2_scaling_configuration = {
    min_capacity = 0.5
    max_capacity = 4.0
  }

  instance_class = "db.serverless"
  instances = {
    one = {}
  }

  vpc_id               = module.vpc.vpc_id
  db_subnet_group_name = module.vpc.database_subnet_group_name
  security_group_rules = {
    vpc_ingress = {
      cidr_blocks = [module.vpc.vpc_cidr_block]
    }
  }

  master_username = "admin"
  master_password = var.db_password
  manage_master_user_password = false
  skip_final_snapshot = true
}

# ------------------------------------------------------------------------------
# EKS Cluster (Worker & BFF Engine)
# ------------------------------------------------------------------------------
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = "${var.project_name}-cluster"
  cluster_version = "1.28"

  vpc_id                         = module.vpc.vpc_id
  subnet_ids                     = module.vpc.private_subnets
  cluster_endpoint_public_access = true

  eks_managed_node_groups = {
    default = {
      min_size     = 2
      max_size     = 5
      desired_size = 2

      instance_types = ["t3.medium"]
      capacity_type  = "ON_DEMAND"
    }
  }
}