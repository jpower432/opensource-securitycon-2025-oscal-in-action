package main
import rego.v1

test_github_branch_protections if {
   cfg := parse_config_file("example.json")
    deny with input as cfg
}
