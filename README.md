# driftwatch

Lightweight daemon that detects infrastructure config drift by comparing live state against version-controlled snapshots.

---

## Installation

```bash
go install github.com/yourorg/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/driftwatch.git && cd driftwatch && go build -o driftwatch .
```

---

## Usage

Initialize a snapshot of your current infrastructure state:

```bash
driftwatch snapshot --config ./driftwatch.yaml --out ./snapshots/
```

Run the daemon to continuously monitor for drift:

```bash
driftwatch watch --snapshot ./snapshots/baseline.json --interval 60s
```

When drift is detected, `driftwatch` reports the differing fields and their expected vs. live values:

```
[DRIFT] resource: aws_security_group.web
  - ingress.0.cidr_blocks: expected ["10.0.0.0/8"], got ["0.0.0.0/0"]
```

### Example Config (`driftwatch.yaml`)

```yaml
providers:
  - aws
  - kubernetes
snapshot_dir: ./snapshots
alert:
  webhook: https://hooks.example.com/notify
```

---

## Requirements

- Go 1.21+
- Credentials configured for your target infrastructure providers

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss significant changes.

---

## License

[MIT](LICENSE)