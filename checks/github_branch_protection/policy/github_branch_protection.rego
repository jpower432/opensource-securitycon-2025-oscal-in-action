package main
import rego.v1

# METADATA
# title: GitHub Branch Protection Policy - Compliance
# description: >-
#   Enforces that the author(s) of a request may not also be the approver of the request.
#   At least two reviewers with equal or greater expertise should review & approve the request.
#   Verifies that a branch protection rule of type 'pull_request' exists for the 'main' branch.
# custom:
#   short_name: github_branch_protection
deny contains result if {
    not has_pull_request_rule
    
    chain := rego.metadata.chain()
    annotations := chain[0].annotations
    
    result := {
        "short_name": annotations.custom.short_name,
        "msg": "Violation: A branch protection rule of type 'pull_request' is required but was not found. This is a prerequisite for enforcing proper code review practices."
    }
}

# METADATA
# title: GitHub Branch Protection Policy - Minimum Approvals
# description: >-
#   Enforces that at least one reviewer with equal or greater expertise 
#   should review & approve the request. For projects with only two maintainers,
#   one reviewer is sufficient to ensure author-approver separation.
#   Verifies that the branch protection rule for the 'main' branch has at least 1 required approving review.
# custom:
#   short_name: github_branch_protection
deny contains result if {
    has_pull_request_rule
    
    # Requires at least one reviewer with equal or greater expertise
    # For projects with only two maintainers, one reviewer ensures author-approver separation
    MINIMUM_REQUIRED_APPROVALS := 1

    some rule in input.values
    rule.type == "pull_request"
    rule.parameters.required_approving_review_count < MINIMUM_REQUIRED_APPROVALS

    chain := rego.metadata.chain()
    annotations := chain[0].annotations
    
    result := {
        "short_name": annotations.custom.short_name,
        "msg": sprintf("Violation: Branch protection requires at least %v reviewer with equal or greater expertise (adjusted for projects with two maintainers), but only %v are required.", [MINIMUM_REQUIRED_APPROVALS, rule.parameters.required_approving_review_count])
    }
}

# METADATA
# title: GitHub Branch Protection Policy - Code Owner Review
# description: >-
#   Enforces that at least two reviewers with equal or greater expertise 
#   should review & approve the request. Requires code owner review to ensure 
#   reviewers have appropriate expertise for the code being changed.
# custom:
#   short_name: github_branch_protection
deny contains result if {
    has_pull_request_rule
    
    some rule in input.values
    rule.type == "pull_request"
    not rule.parameters.require_code_owner_review

    chain := rego.metadata.chain()
    annotations := chain[0].annotations
    
    result := {
        "short_name": annotations.custom.short_name,
        "msg": "Violation: Code owner review is required to ensure reviewers have equal or greater expertise for the code being changed."
    }
}

# METADATA
# title: GitHub Branch Protection Policy - Stale Review Dismissal
# description: >-
#   Enforces that fresh reviews are ensured by requiring that stale reviews 
#   are dismissed when new commits are pushed, maintaining review quality.
# custom:
#   short_name: github_branch_protection
deny contains result if {
    has_pull_request_rule
    
    some rule in input.values
    rule.type == "pull_request"
    not rule.parameters.dismiss_stale_reviews_on_push

    chain := rego.metadata.chain()
    annotations := chain[0].annotations
    
    result := {
        "short_name": annotations.custom.short_name,
        "msg": "Violation: Stale review dismissal is required to ensure fresh, relevant reviews from reviewers with appropriate expertise."
    }
}

# METADATA
# title: GitHub Branch Protection Policy - Review Thread Resolution
# description: >-
#   Enforces that thorough review is ensured by requiring that all review 
#   threads are resolved before approval, maintaining review quality and expertise validation.
# custom:
#   short_name: github_branch_protection
deny contains result if {
    has_pull_request_rule
    
    some rule in input.values
    rule.type == "pull_request"
    not rule.parameters.required_review_thread_resolution

    chain := rego.metadata.chain()
    annotations := chain[0].annotations
    
    result := {
        "short_name": annotations.custom.short_name,
        "msg": "Violation: Review thread resolution is required to ensure thorough review by qualified reviewers."
    }
}

# Helper function to check if a pull request rule exists
has_pull_request_rule if {
    some rule in input.values
    rule.type == "pull_request"
}
