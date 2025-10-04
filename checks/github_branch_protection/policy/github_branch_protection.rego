package main
import rego.v1

# METADATA
# title: GitHub Branch Protection Policy - Missing Pull Request Rule
# description: >-
#   Verifies that a branch protection rule of type 'pull_request' exists
#   for the 'main' branch. This is a prerequisite for enforcing code review.
# custom:
#   short_name: github_branch_protection
deny contains result if {
    not has_pull_request_rule
    
    chain := rego.metadata.chain()
    annotations := chain[0].annotations
    
    result := {
        "short_name": annotations.custom.short_name,
        "msg": "A branch protection rule of type 'pull_request' is required but was not found."
    }
}

# METADATA
# title: GitHub Branch Protection Policy - Minimum Approvals
# description: >-
#   Verifies that the branch protection rule for the 'main' branch
#   has at least the configured minimum number of required approving reviews.
#   This policy ensures proper code review practices are enforced.
# custom:
#   short_name: github_branch_protection
deny contains result if {
    has_pull_request_rule
    
    # Configuration constant for minimum required approvals
    # This checks for conformance to our organization policy
    MINIMUM_REQUIRED_APPROVALS := 2

    some rule in input.values
    rule.type == "pull_request"
    rule.parameters.required_approving_review_count < MINIMUM_REQUIRED_APPROVALS

    # Ensure the short name is added to the resutls
    # for tracebility
    chain := rego.metadata.chain()
    annotations := chain[0].annotations
    
    result := {
        "short_name": annotations.custom.short_name,
        "msg": sprintf("Branch protection for 'main' requires pull request reviews but has less than the configured minimum of %v.", [MINIMUM_REQUIRED_APPROVALS])
    }
}

# Helper function to check if a pull request rule exists
has_pull_request_rule if {
    some rule in input.values
    rule.type == "pull_request"
}
