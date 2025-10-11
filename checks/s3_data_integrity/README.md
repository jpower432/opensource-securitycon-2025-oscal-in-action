# S3 Data Integrity Policy - CNSCC-STO-03.01 Compliance

This directory contains the implementation for CNSCC-STO-03.01: "Hashing and checksums are added to blocks, objects or files to detect and recover from corrupted data and provide protection against tampering."

## Files

- `policy/s3_data_integrity_policy.rego` - OPA policy that enforces CNSCC-STO-03.01 compliance
- `policy/s3_data_integrity_policy_test.rego` - Comprehensive test suite
- `main.tf` - Terraform configuration demonstrating STO-03 compliance
- `README.md` - This file

## CNSCC-STO-03.01 Requirements

The policy enforces the following CNSCC-STO-03.01 requirements:

1. **Versioning Enabled**: S3 bucket versioning must be enabled to maintain multiple versions of objects for corruption recovery
2. **Server-Side Encryption**: AES256 encryption must be enabled to protect data integrity
3. **Secure Transport Enforcement**: All data transmission must use TLS/HTTPS to prevent data corruption during transit
4. **Lifecycle Configuration**: Proper lifecycle management to maintain data integrity over time
5. **Monitoring and Logging**: Bucket notifications and CloudWatch logging to detect integrity issues

## How S3 Implements Data Integrity

### Built-in Checksums
- **ETags**: S3 automatically generates MD5 checksums for each object (stored as ETags)
- **Content-MD5**: Clients can provide MD5 checksums for validation during upload
- **SHA-256**: S3 uses SHA-256 for internal integrity verification

### Versioning for Recovery
- **Object Versions**: Multiple versions of objects are maintained
- **Corruption Detection**: Can compare versions to detect corruption
- **Recovery**: Previous versions can be restored if corruption is detected

### Encryption for Protection
- **AES-256**: Server-side encryption protects data at rest
- **Key Management**: AWS KMS integration for encryption key management
- **Tamper Detection**: Encrypted data changes are easily detectable

## Policy Implementation

The Rego policy validates the following Terraform resources:

### Required Resources
1. **aws_s3_bucket_versioning** - Must have status "Enabled"
2. **aws_s3_bucket_server_side_encryption_configuration** - Must use AES256
3. **aws_s3_bucket_policy** - Must enforce secure transport (HTTPS only)
4. **aws_s3_bucket_lifecycle_configuration** - Must be present and enabled
5. **aws_s3_bucket_notification** - Must be configured for monitoring

### Validation Rules
- **Versioning Check**: Ensures versioning is enabled for data recovery
- **Encryption Check**: Validates AES256 encryption is configured
- **Transport Security**: Verifies HTTPS-only policy is enforced
- **Lifecycle Management**: Confirms lifecycle configuration exists
- **Monitoring**: Ensures bucket notifications are configured

## Testing the Policy

Testing is performed in CI using conftest. The policy can be tested locally using OPA:

### 1. Run Unit Tests

```bash
cd policy
opa test s3_data_integrity_policy_test.rego s3_data_integrity_policy.rego
```

### 2. Test with Terraform Configuration

```bash
cd policy
opa eval --data s3_data_integrity_policy.rego --input ../main.tf 'data.terraform.aws.s3.deny'
```

## Expected Behavior

When you run the test, the Rego policy should:

### For Compliant Configuration:
- Return empty deny array `[]`
- All STO-03.01 requirements are met

### For Non-Compliant Configuration:
- Return specific violation messages
- Identify which data integrity controls are missing

## Compliance Configuration

To achieve CNSCC-STO-03.01 compliance, configure S3 with:

- **Versioning**: `versioning_configuration.status = "Enabled"`
- **Encryption**: `sse_algorithm = "AES256"`
- **Secure Transport**: Policy denying non-HTTPS access
- **Lifecycle**: Proper lifecycle configuration for data management
- **Monitoring**: CloudWatch logging and bucket notifications

## Security Benefits

This implementation provides:

1. **Data Corruption Detection**: Versioning and checksums detect corrupted data
2. **Recovery Capability**: Previous versions can be restored
3. **Tamper Protection**: Encryption makes tampering detectable
4. **Transit Security**: HTTPS-only ensures data integrity during transmission
5. **Monitoring**: Continuous monitoring for integrity issues

## Technical Details

### Checksum Mechanisms
- **ETag Generation**: S3 automatically creates MD5 checksums
- **Client Validation**: Content-MD5 header for upload validation
- **Internal Verification**: SHA-256 for AWS internal integrity checks

### Versioning Benefits
- **Corruption Recovery**: Access previous uncorrupted versions
- **Change Tracking**: Monitor object changes over time
- **Audit Trail**: Complete history of object modifications

### Encryption Protection
- **At-Rest Security**: AES-256 encryption protects stored data
- **Key Management**: AWS KMS integration for key security
- **Tamper Evidence**: Encrypted data changes are obvious

## Security Note

This implementation demonstrates how OSCAL compliance controls can be enforced through policy-as-code to ensure proper data integrity mechanisms and protect against data corruption and tampering in cloud storage systems.