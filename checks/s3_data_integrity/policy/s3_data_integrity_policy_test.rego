package terraform.aws.s3

import rego.v1

# Test cases for CNSCC-STO-03.01 compliance

# Test 1: Compliant configuration should pass
test_compliant_configuration {
    input := {
        "resource_changes": [
            {
                "type": "aws_s3_bucket",
                "change": {
                    "actions": ["create"],
                    "after": {
                        "bucket": "compliant-bucket"
                    }
                }
            },
            {
                "type": "aws_s3_bucket_versioning",
                "change": {
                    "after": {
                        "bucket": "compliant-bucket",
                        "versioning_configuration": [{"status": "Enabled"}]
                    }
                }
            },
            {
                "type": "aws_s3_bucket_server_side_encryption_configuration",
                "change": {
                    "after": {
                        "bucket": "compliant-bucket",
                        "rule": [{
                            "apply_server_side_encryption_by_default": [{"sse_algorithm": "AES256"}]
                        }]
                    }
                }
            },
            {
                "type": "aws_s3_bucket_policy",
                "change": {
                    "after": {
                        "bucket": "compliant-bucket",
                        "policy": json.marshal({
                            "Statement": [{
                                "Effect": "Deny",
                                "Condition": {
                                    "Bool": {
                                        "aws:SecureTransport": "false"
                                    }
                                }
                            }]
                        })
                    }
                }
            },
            {
                "type": "aws_s3_bucket_lifecycle_configuration",
                "change": {
                    "after": {
                        "bucket": "compliant-bucket",
                        "rule": [{"status": "Enabled"}]
                    }
                }
            },
            {
                "type": "aws_s3_bucket_notification",
                "change": {
                    "after": {
                        "bucket": "compliant-bucket"
                    }
                }
            }
        ]
    }
    
    deny(input) == []
}

# Test 2: Missing versioning should fail
test_missing_versioning {
    input := {
        "resource_changes": [
            {
                "type": "aws_s3_bucket",
                "change": {
                    "actions": ["create"],
                    "after": {
                        "bucket": "no-versioning-bucket"
                    }
                }
            }
        ]
    }
    
    count(deny(input)) == 1
    deny(input)[0] contains "versioning must be enabled"
}

# Test 3: Missing encryption should fail
test_missing_encryption {
    input := {
        "resource_changes": [
            {
                "type": "aws_s3_bucket",
                "change": {
                    "actions": ["create"],
                    "after": {
                        "bucket": "no-encryption-bucket"
                    }
                }
            },
            {
                "type": "aws_s3_bucket_versioning",
                "change": {
                    "after": {
                        "bucket": "no-encryption-bucket",
                        "versioning_configuration": [{"status": "Enabled"}]
                    }
                }
            }
        ]
    }
    
    count(deny(input)) == 1
    deny(input)[0] contains "server-side encryption must be enabled"
}

# Test 4: Missing secure transport enforcement should fail
test_missing_secure_transport {
    input := {
        "resource_changes": [
            {
                "type": "aws_s3_bucket",
                "change": {
                    "actions": ["create"],
                    "after": {
                        "bucket": "no-secure-transport-bucket"
                    }
                }
            },
            {
                "type": "aws_s3_bucket_versioning",
                "change": {
                    "after": {
                        "bucket": "no-secure-transport-bucket",
                        "versioning_configuration": [{"status": "Enabled"}]
                    }
                }
            },
            {
                "type": "aws_s3_bucket_server_side_encryption_configuration",
                "change": {
                    "after": {
                        "bucket": "no-secure-transport-bucket",
                        "rule": [{
                            "apply_server_side_encryption_by_default": [{"sse_algorithm": "AES256"}]
                        }]
                    }
                }
            }
        ]
    }
    
    count(deny(input)) == 1
    deny(input)[0] contains "secure transport must be enforced"
}

# Test 5: Missing lifecycle configuration should fail
test_missing_lifecycle {
    input := {
        "resource_changes": [
            {
                "type": "aws_s3_bucket",
                "change": {
                    "actions": ["create"],
                    "after": {
                        "bucket": "no-lifecycle-bucket"
                    }
                }
            },
            {
                "type": "aws_s3_bucket_versioning",
                "change": {
                    "after": {
                        "bucket": "no-lifecycle-bucket",
                        "versioning_configuration": [{"status": "Enabled"}]
                    }
                }
            },
            {
                "type": "aws_s3_bucket_server_side_encryption_configuration",
                "change": {
                    "after": {
                        "bucket": "no-lifecycle-bucket",
                        "rule": [{
                            "apply_server_side_encryption_by_default": [{"sse_algorithm": "AES256"}]
                        }]
                    }
                }
            },
            {
                "type": "aws_s3_bucket_policy",
                "change": {
                    "after": {
                        "bucket": "no-lifecycle-bucket",
                        "policy": json.marshal({
                            "Statement": [{
                                "Effect": "Deny",
                                "Condition": {
                                    "Bool": {
                                        "aws:SecureTransport": "false"
                                    }
                                }
                            }]
                        })
                    }
                }
            }
        ]
    }
    
    count(deny(input)) == 1
    deny(input)[0] contains "lifecycle configuration must be present"
}

# Test 6: Missing monitoring should fail
test_missing_monitoring {
    input := {
        "resource_changes": [
            {
                "type": "aws_s3_bucket",
                "change": {
                    "actions": ["create"],
                    "after": {
                        "bucket": "no-monitoring-bucket"
                    }
                }
            },
            {
                "type": "aws_s3_bucket_versioning",
                "change": {
                    "after": {
                        "bucket": "no-monitoring-bucket",
                        "versioning_configuration": [{"status": "Enabled"}]
                    }
                }
            },
            {
                "type": "aws_s3_bucket_server_side_encryption_configuration",
                "change": {
                    "after": {
                        "bucket": "no-monitoring-bucket",
                        "rule": [{
                            "apply_server_side_encryption_by_default": [{"sse_algorithm": "AES256"}]
                        }]
                    }
                }
            },
            {
                "type": "aws_s3_bucket_policy",
                "change": {
                    "after": {
                        "bucket": "no-monitoring-bucket",
                        "policy": json.marshal({
                            "Statement": [{
                                "Effect": "Deny",
                                "Condition": {
                                    "Bool": {
                                        "aws:SecureTransport": "false"
                                    }
                                }
                            }]
                        })
                    }
                }
            },
            {
                "type": "aws_s3_bucket_lifecycle_configuration",
                "change": {
                    "after": {
                        "bucket": "no-monitoring-bucket",
                        "rule": [{"status": "Enabled"}]
                    }
                }
            }
        ]
    }
    
    count(deny(input)) == 1
    deny(input)[0] contains "bucket monitoring must be enabled"
}