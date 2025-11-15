<!--
DiliVet – ML-DSA diagnostics and vetting toolkit
Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)
-->

# Kubernetes Deployment

This directory contains Kubernetes manifests for deploying dilivet as a containerized service.

## Quick Start

```bash
# Apply all manifests
kubectl apply -k k8s/

# Check deployment status
kubectl get deployment dilivet
kubectl get pods -l app=dilivet

# View logs
kubectl logs -l app=dilivet -f

# Run dilivet command in pod
kubectl exec -it deployment/dilivet -- /dilivet -version
```

## Manifests

- **deployment.yaml** - Deployment with security hardening (non-root, read-only filesystem)
- **service.yaml** - ClusterIP service (modify for LoadBalancer/NodePort if needed)
- **configmap.yaml** - Runtime configuration for ML-DSA parameters
- **kustomization.yaml** - Kustomize configuration for environment-specific overlays

## Security Features

- Runs as non-root user (UID 65534)
- Read-only root filesystem
- No privilege escalation
- Minimal capabilities (all dropped)
- Resource limits enforced

## Customization

### Production Overlay

Create `k8s/overlays/production/kustomization.yaml`:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: production

resources:
- ../../base

replicas:
- name: dilivet
  count: 3

images:
- name: ghcr.io/codethor0/dilivet
  newTag: v0.2.0
```

Apply with: `kubectl apply -k k8s/overlays/production/`

### Environment Variables

Add to deployment.yaml under `spec.template.spec.containers[0]`:

```yaml
env:
- name: DILIVET_LOG_LEVEL
  valueFrom:
    configMapKeyRef:
      name: dilivet-config
      key: log-level
```

## Health Checks

If dilivet exposes HTTP endpoints, add liveness/readiness probes:

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
```

## Scaling

```bash
# Manual scaling
kubectl scale deployment dilivet --replicas=3

# Horizontal Pod Autoscaler
kubectl autoscale deployment dilivet --cpu-percent=70 --min=1 --max=10
```

## Cleanup

```bash
kubectl delete -k k8s/
```
