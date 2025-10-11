# S3 Data Integrity Policy - CNSCC-STO-03.01 Compliance (INTENTIONALLY FAILING)
# This demonstrates S3 bucket configuration that VIOLATES data integrity requirements
# by missing critical security features for checksums and hashing mechanisms.

# S3 bucket WITHOUT data integrity features (INTENTIONALLY INSECURE)
resource "aws_s3_bucket" "data_integrity_bucket" {
  bucket = "my-insecure-bucket-${random_string.bucket_suffix.result}"
  
  tags = {
    Name        = "Insecure Data Bucket"
    Purpose     = "CNSCC-STO-03.01 Violation Demo"
    Environment = "demo"
  }
}

# Random string for unique bucket naming
resource "random_string" "bucket_suffix" {
  length  = 8
  special = false
  upper   = false
}

# S3 bucket versioning DISABLED (INTENTIONALLY INSECURE for STO-03.01 violation)
# resource "aws_s3_bucket_versioning" "data_integrity_versioning" {
#   bucket = aws_s3_bucket.data_integrity_bucket.id
#   versioning_configuration {
#     status = "Enabled"
#   }
# }

# S3 bucket server-side encryption DISABLED (INTENTIONALLY INSECURE for STO-03.01 violation)
# resource "aws_s3_bucket_server_side_encryption_configuration" "data_integrity_encryption" {
#   bucket = aws_s3_bucket.data_integrity_bucket.id
#
#   rule {
#     apply_server_side_encryption_by_default {
#       sse_algorithm = "AES256"
#     }
#     bucket_key_enabled = true
#   }
# }

# S3 bucket public access block for security
resource "aws_s3_bucket_public_access_block" "data_integrity_pab" {
  bucket = aws_s3_bucket.data_integrity_bucket.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# S3 bucket policy REMOVED (INTENTIONALLY INSECURE for STO-03.01 violation)
# No secure transport enforcement or checksum validation
# resource "aws_s3_bucket_policy" "data_integrity_policy" {
#   bucket = aws_s3_bucket.data_integrity_bucket.id
#   ...
# }

# Data source for current AWS account
data "aws_caller_identity" "current" {}

# S3 bucket notification REMOVED (INTENTIONALLY INSECURE for STO-03.01 violation)
# No monitoring or integrity tracking
# resource "aws_s3_bucket_notification" "data_integrity_notification" {
#   ...
# }

# CloudWatch log group REMOVED (INTENTIONALLY INSECURE for STO-03.01 violation)
# No monitoring or logging
# resource "aws_cloudwatch_log_group" "data_integrity_logs" {
#   ...
# }

# S3 bucket lifecycle configuration REMOVED (INTENTIONALLY INSECURE for STO-03.01 violation)
# No lifecycle management or data integrity controls
# resource "aws_s3_bucket_lifecycle_configuration" "data_integrity_lifecycle" {
#   ...
# }

# Output values for reference (FAILING CONFIGURATION)
output "bucket_name" {
  description = "Name of the S3 bucket WITHOUT data integrity features (STO-03.01 VIOLATION)"
  value       = aws_s3_bucket.data_integrity_bucket.bucket
}

output "bucket_arn" {
  description = "ARN of the S3 bucket WITHOUT data integrity features (STO-03.01 VIOLATION)"
  value       = aws_s3_bucket.data_integrity_bucket.arn
}

output "versioning_status" {
  description = "Versioning status - DISABLED (STO-03.01 VIOLATION)"
  value       = "Disabled - No versioning configured"
}