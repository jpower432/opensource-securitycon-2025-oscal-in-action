package main
import rego.v1

# METADATA
# title: Branch Protection Rules Present
# description: >-
#   Confirms that branch protection rules are present in the input.
# custom:
#   short_name: rules_present
#   solution: >-
#     Configure at least one branch protection rule for the primary branch.
deny contains msg if {
    # Check if the overall rules array exists
    not input.values
    msg := "No branch protection rules found in the input."
}


# METADATA
# title: Pull Request Rule Required
# description: >-
#   Confirms that a branch protection rule of type 'pull_request' is present.
#   This is a prerequisite for checking approval counts.
# custom:
#   short_name: pull_request_rule_required
#   solution: >-
#     Add a branch protection rule with type 'pull_request' to your branch protection settings.
#   depends_on:
#   - github_branch_protection.rules_present
deny contains msg if {
    not _has_pull_request_rule
    msg := "A branch protection rule of type 'pull_request' is required for the primary branch but was not found."
}


# METADATA
# title: Minimum Approvals for Main Branch
# description: >-
#   Verifies that the branch protection rule for the 'main' branch
#   has at least the configured minimum number of required approving reviews.
# custom:
#   short_name: min_approvals_check
#   solution: >-
#     Increase the 'required_approving_review_count' in the branch protection settings to meet or exceed the policy's minimum.
#   depends_on:
#   - github_branch_protection.rules_present
#   - github_branch_protection.pull_request_rule_required
deny contains msg if {
    required_count := data.rule_data__configuration__main_branch_min_approvals
    some rule in input.values
    rule.type == "pull_request"
    rule.parameters.required_approving_review_count < required_count

    msg := sprintf("Branch protection for 'main' requires pull request reviews but has less than the configured minimum of %v required approving reviews.", [required_count])
}

_has_pull_request_rule if {
    some rule in input.values
    rule.type == "pull_request"
}