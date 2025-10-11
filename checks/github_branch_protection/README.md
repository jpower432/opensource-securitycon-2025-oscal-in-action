# GitHub Branch Protection Policy - CNSCC-SSC-09.01 Compliance

This directory contains the evaluation procedure for CNSCC-SSC-09.01: "The author(s) of a request may not also be the approver of the request. At least two reviewers with equal or greater expertise should review & approve the request."

**Note**: This policy has been adjusted for projects with only two maintainers, requiring a minimum of one reviewer to ensure author-approver separation while maintaining the core principle of the control.

## Files

- `policy/github_branch_protection.rego` - OPA policy that enforces CNSCC-SSC-09.01 compliance
- `example.json` - Example GitHub branch protection configuration (intentionally non-compliant)
- `README.md` - This file

## Policy Requirements

The policy enforces the following CNSCC-SSC-09.01 requirements (adjusted for projects with two maintainers):

1. **Pull Request Rule Exists**: A branch protection rule of type 'pull_request' must exist
2. **Minimum Reviewers**: At least 1 reviewer with equal or greater expertise must review requests (ensures author-approver separation with two maintainers)
3. **Code Owner Review**: Code owner review is required to ensure reviewers have appropriate expertise
4. **Stale Review Dismissal**: Stale reviews must be dismissed on new commits to maintain review quality
5. **Review Thread Resolution**: All review threads must be resolved before approval

## Testing the Policy

Testing is performed in CI using conftest. The policy can be tested locally using OPA:

1. Install OPA (Open Policy Agent):
   ```bash
   # On Ubuntu/Debian
   sudo apt-get install opa
   
   # On macOS
   brew install opa
   
   # Or download from https://www.openpolicyagent.org/
   ```

2. Test the policy with the example data:
   ```bash
   opa eval --data policy/github_branch_protection.rego --input example.json 'data.main.deny'
   ```

3. The policy should identify multiple CNSCC-SSC-09.01 violations in the example configuration.

## Expected Behavior

When you run the test, the Rego policy should identify the following CNSCC-SSC-09.01 violations (adjusted for projects with two maintainers):
- Stale review dismissal not enabled (needed for fresh reviews)
- Note: The example configuration now shows a compliant setup with 1 reviewer for two maintainers

## Compliance Configuration

To achieve CNSCC-SSC-09.01 compliance (adjusted for projects with two maintainers), configure GitHub branch protection with:
- `required_approving_review_count`: 1 (ensures author-approver separation with two maintainers)
- `require_code_owner_review`: true
- `dismiss_stale_reviews_on_push`: true
- `required_review_thread_resolution`: true

This ensures that:
- Authors cannot approve their own requests (achieved with one reviewer when there are only two maintainers)
- At least one qualified reviewer reviews each request
- Reviewers have appropriate expertise for the code being changed
- Reviews are fresh and thorough

## Security Note

This example demonstrates how OSCAL compliance controls can be enforced through policy-as-code to ensure proper code review practices and prevent security vulnerabilities from being introduced through inadequate review processes.