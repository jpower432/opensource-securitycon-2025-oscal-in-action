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
- `grafana-cloud-endpoint`: The Grafana Cloud Loki endpoint URL
- `grafana-cloud-instance-id`: Your Grafana Cloud instance ID (used as username for basic auth)
- `grafana-cloud-api-key`: Your Grafana Cloud API key (used as password for basic auth)

### Environment Variables

For security, credentials can be provided via environment variables instead of configuration:

- `GRAFANA_CLOUD_ENDPOINT`: Grafana Cloud Loki endpoint URL
- `GRAFANA_CLOUD_INSTANCE_ID`: Your Grafana Cloud instance ID
- `GRAFANA_CLOUD_API_KEY`: Your Grafana Cloud API key

**Security Note**: Environment variables take precedence over configuration values and are recommended for production deployments.

## Usage

### Configuration Examples

#### Local Loki Instance
```json
{
  "loki_url": "http://localhost:3100"
}
```

#### Grafana Cloud (Configuration)
```json
{
  "grafana-cloud-endpoint": "https://logs-prod-us-central1.grafana.net",
  "grafana-cloud-instance-id": "your-instance-id",
  "grafana-cloud-api-key": "your-api-key"
}
```

#### Grafana Cloud (Environment Variables - Recommended)
```bash
export GRAFANA_CLOUD_ENDPOINT="https://logs-prod-us-central1.grafana.net"
export GRAFANA_CLOUD_INSTANCE_ID="your-instance-id"
export GRAFANA_CLOUD_API_KEY="your-api-key"
```

```json
{
  "policy-results": "/path/to/results"
}
```

### Authentication Methods

The plugin uses **Basic Authentication** for Grafana Cloud:

- **Username**: Your Grafana Cloud instance ID (`grafana-cloud-instance-id`)
- **Password**: Your Grafana Cloud API key (`grafana-cloud-api-key`)

This follows Grafana Cloud's standard authentication pattern where the instance ID and API key are used together for basic authentication.

### Security Best Practices

1. **Use Environment Variables**: Store sensitive credentials in environment variables rather than configuration files
2. **Avoid Logging Secrets**: The plugin does not log instance IDs or API keys to prevent PII exposure
3. **Secure Storage**: Ensure environment variables are stored securely in your deployment environment
4. **Rotation**: Regularly rotate your Grafana Cloud API keys
5. **Least Privilege**: Use API keys with minimal required permissions

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