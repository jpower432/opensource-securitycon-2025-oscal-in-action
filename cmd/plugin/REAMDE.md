# OPA/Loki Plugin

This plugin is a hybrid that creates OPA policy bundles and can pull evidence from Loki by policy ID.

## Features

- **Policy Generation**: Creates OPA policy bundles using the existing OPA plugin functionality
- **Evidence Collection**: Pulls evidence from Loki logs by policy ID
- **Evidence Processing**: Converts Loki log entries into evidence links and observations
- **OSCAL Integration**: Returns results in OSCAL-compatible format

## Configuration

The plugin accepts the following configuration parameters:

- `loki_url`: The base URL of your Loki instance (e.g., `http://localhost:3100`)

## Usage

### Configuration Example

```json
{
  "loki_url": "http://localhost:3100"
}
```

### Policy ID Mapping

The plugin queries Loki for logs with the label `policy_id` matching the policy being evaluated. Each policy ID corresponds to a rule, and the evidence is collected from Loki logs.

### Evidence Processing

The `GetResults()` function:

1. **Queries Loki**: Uses the Loki API to fetch logs matching the policy ID
2. **Processes Evidence**: Converts log entries into evidence links
3. **Creates Observations**: Generates observations for each control in the policy
4. **Returns Results**: Provides structured results in OSCAL format

### Evidence Links

Each log entry becomes an evidence link with:
- **Description**: Human-readable description of the evidence
- **Href**: Direct link to the Loki query for that specific log entry
- **Properties**: Metadata including timestamp, message, labels, and source

### Observations

For each policy, the plugin creates:
- **ObservationByCheck**: Contains all evidence for a specific policy check
- **Subjects**: Individual log entries as subjects with their properties
- **RelevantEvidences**: Links to the evidence in Loki

## Loki Query Format

The plugin queries Loki using the following format:
```
{service_name=~".+"} } policy_id="<policy_id>"
```

This assumes your logs are labeled with `policy_id` to identify which policy they relate to.

## Error Handling

- If Loki client is not configured, returns empty results
- If Loki query fails for one policy, continues with other policies
- Logs all operations for debugging and monitoring