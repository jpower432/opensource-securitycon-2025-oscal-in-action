package main
import rego.v1

# METADATA
# title: Minimum Approvals for Main Branch
# description: >-
#   Verifies that the branch protection rule for the 'main' branch
#   has at least the configured minimum number of required approving reviews.
# custom:
#   short_name: github_branch_protection
deny contains result if {
    # Check if a pull request rule exists
    not has_pull_request_rule

    chain := rego.metadata.chain()
    annotations := chain[0].annotations

    result := {
        "short_name": annotations.custom.short_name,
        "msg": "A branch protection rule of type 'pull_request' is required but was not found."
    }
}

# METADATA
# title: Minimum Approvals for Main Branch
# description: >-
#   Verifies that the branch protection rule for the 'main' branch
#   has at least the configured minimum number of required approving reviews.
# custom:
#   short_name: github_branch_protection
deny contains result if {
    # Check if the required number of approvals is met
    has_pull_request_rule

    required_count := 2
    some rule in input.values
    rule.type == "pull_request"
    rule.parameters.required_approving_review_count < required_count

     chain := rego.metadata.chain()
    annotations := chain[0].annotations

    result := {
        "short_name": annotations.custom.short_name,
        "msg":  sprintf("Branch protection for 'main' requires pull request reviews but has less than the configured minimum of %v.", [required_count])
    }
}

has_pull_request_rule if {
    some rule in input.values
    rule.type == "pull_request"
}
