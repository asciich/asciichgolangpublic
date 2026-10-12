# 🔁 hostreconciler

> Bring the GitOps approach to your hosts.

`hostreconciler` periodically reads a declarative configuration and
reconciles the current host to match it: **desired state in, converged host out.**

## Config structure

### Minimal config

```text
<config-dir>/
└── host.yaml
```

**`host.yaml`**

```yaml
---
hostname: host.example.com
```

### Reference: `host.yaml`

| Field      | Type     | Required | Description                              |
|------------|----------|:--------:|------------------------------------------|
| `hostname` | `string` |    ✅    | Fully qualified hostname of the host.    |
