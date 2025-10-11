package terraform.aws.s3

# This policy enforces CNSCC-STO-03.01: Hashing and checksums are added to blocks, 
# objects or files to detect and recover from corrupted data and provide protection 
# against tampering. It ensures S3 buckets are configured with proper data integrity 
# mechanisms.

# Deny if S3 bucket doesn't have versioning enabled (required for data integrity)
deny[msg] {
    # Check for S3 bucket resources
    resource := input.resource_changes[_]
    resource.type == "aws_s3_bucket"
    resource.change.actions[_] == "create"
    
    bucket_name := resource.change.after.bucket
    
    # Check if versioning is not enabled
    not has_versioning_enabled(bucket_name)
    
    msg := sprintf("S3 bucket '%s' violates CNSCC-STO-03.01: Versioning must be enabled for data integrity and corruption recovery", [bucket_name])
}

# Deny if S3 bucket doesn't have server-side encryption enabled
deny[msg] {
    # Check for S3 bucket resources
    resource := input.resource_changes[_]
    resource.type == "aws_s3_bucket"
    resource.change.actions[_] == "create"
    
    bucket_name := resource.change.after.bucket
    
    # Check if server-side encryption is not configured
    not has_server_side_encryption(bucket_name)
    
    msg := sprintf("S3 bucket '%s' violates CNSCC-STO-03.01: Server-side encryption must be enabled for data integrity protection", [bucket_name])
}

# Deny if S3 bucket doesn't have secure transport enforcement
deny[msg] {
    # Check for S3 bucket resources
    resource := input.resource_changes[_]
    resource.type == "aws_s3_bucket"
    resource.change.actions[_] == "create"
    
    bucket_name := resource.change.after.bucket
    
    # Check if bucket policy doesn't enforce secure transport
    not has_secure_transport_enforcement(bucket_name)
    
    msg := sprintf("S3 bucket '%s' violates CNSCC-STO-03.01: Secure transport must be enforced to protect data integrity during transmission", [bucket_name])
}

# Deny if S3 bucket doesn't have lifecycle configuration for data integrity
deny[msg] {
    # Check for S3 bucket resources
    resource := input.resource_changes[_]
    resource.type == "aws_s3_bucket"
    resource.change.actions[_] == "create"
    
    bucket_name := resource.change.after.bucket
    
    # Check if lifecycle configuration is missing
    not has_lifecycle_configuration(bucket_name)
    
    msg := sprintf("S3 bucket '%s' violates CNSCC-STO-03.01: Lifecycle configuration must be present to maintain data integrity over time", [bucket_name])
}

# Deny if S3 bucket doesn't have monitoring/logging enabled
deny[msg] {
    # Check for S3 bucket resources
    resource := input.resource_changes[_]
    resource.type == "aws_s3_bucket"
    resource.change.actions[_] == "create"
    
    bucket_name := resource.change.after.bucket
    
    # Check if bucket notification is missing (indicates no monitoring)
    not has_bucket_notification(bucket_name)
    
    msg := sprintf("S3 bucket '%s' violates CNSCC-STO-03.01: Bucket monitoring must be enabled to detect data integrity issues", [bucket_name])
}

# Helper function to check if versioning is enabled
has_versioning_enabled(bucket_name) {
    versioning := input.resource_changes[_]
    versioning.type == "aws_s3_bucket_versioning"
    versioning.change.after.bucket == bucket_name
    versioning.change.after.versioning_configuration[0].status == "Enabled"
}

# Helper function to check if server-side encryption is configured
has_server_side_encryption(bucket_name) {
    encryption := input.resource_changes[_]
    encryption.type == "aws_s3_bucket_server_side_encryption_configuration"
    encryption.change.after.bucket == bucket_name
    encryption.change.after.rule[0].apply_server_side_encryption_by_default[0].sse_algorithm == "AES256"
}

# Helper function to check if secure transport is enforced
has_secure_transport_enforcement(bucket_name) {
    policy := input.resource_changes[_]
    policy.type == "aws_s3_bucket_policy"
    policy.change.after.bucket == bucket_name
    
    policy_doc := json.unmarshal(policy.change.after.policy)
    statement := policy_doc.Statement[_]
    statement.Effect == "Deny"
    statement.Condition.Bool["aws:SecureTransport"] == "false"
}

# Helper function to check if lifecycle configuration exists
has_lifecycle_configuration(bucket_name) {
    lifecycle := input.resource_changes[_]
    lifecycle.type == "aws_s3_bucket_lifecycle_configuration"
    lifecycle.change.after.bucket == bucket_name
    lifecycle.change.after.rule[0].status == "Enabled"
}

# Helper function to check if bucket notification exists
has_bucket_notification(bucket_name) {
    notification := input.resource_changes[_]
    notification.type == "aws_s3_bucket_notification"
    notification.change.after.bucket == bucket_name
}