output "sqs_queue_url" {
  description = "The URL of the SQS queue for the Go worker"
  value       = aws_sqs_queue.anomaly_queue.url
}

output "redis_endpoint" {
  description = "The endpoint of the ElastiCache Redis cluster"
  value       = aws_elasticache_replication_group.redis.primary_endpoint_address
}

output "aurora_endpoint" {
  description = "The writer endpoint for the Aurora PostgreSQL database"
  value       = module.aurora.cluster_endpoint
}

output "eks_cluster_name" {
  description = "The name of the EKS cluster"
  value       = module.eks.cluster_name
}

output "configure_kubectl" {
  description = "Command to configure kubectl for your new cluster"
  value       = "aws eks --region ${var.region} update-kubeconfig --name ${module.eks.cluster_name}"
}