// Package config provides configuration loading and validation for driftwatch.
//
// Configuration is expressed as a YAML file. The Load function resolves the
// file path using the following precedence:
//
//  1. An explicit path passed by the caller.
//  2. The DRIFTWATCH_CONFIG environment variable.
//  3. A set of well-known default locations (e.g. driftwatch.yaml in the
//     current directory, /etc/driftwatch/config.yaml).
//
// Example configuration:
//
//	paths:
//	  - /etc/nginx
//	  - /etc/ssl/certs
//	interval: 30s
//	snapshot:
//	  dir: /var/driftwatch/snapshots
//	reporter:
//	  format: json
//	  output: stdout
//
// The Validate method is called automatically by LoadFromFile and returns an
// error if required fields are absent or values are out of range.
package config
