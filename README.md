# scm-controller

A Kubernetes native way of interacting with source code management providers
(GitHub, GitLab, BitBucket, etc).

## Getting Started

### Installation

1. Install SCM into the cluster

```bash
kubectl apply -f  <release.yaml>
```

2. Watch the progress

```bash
kubectl -n scm-system get pods --watch
```

### Tutorial

