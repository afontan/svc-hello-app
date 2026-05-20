# svc-hello-app

Tiny Go HTTP service used for testing the local ArgoCD + Datadog setup.

Endpoints:
- `GET /healthz` — readiness/liveness probe
- `GET /hello?name=alejandro` — returns a JSON greeting and the current build version

## Local quickstart

```sh
make tidy
make test
make run            # listens on :8080
curl localhost:8080/hello?name=world
```

## Build and push to GHCR

```sh
make docker push OWNER=afontan
```

## CI

`.github/workflows/ci.yml` runs on every push:
1. `go test` + JUnit upload → Datadog CI Visibility
2. (on `main`) build + push image to `ghcr.io/<owner>/svc-hello:<sha>`
3. opens a commit on `svc-hello-deploy` bumping the image tag — ArgoCD picks it up

### Required GitHub secrets
- `DD_API_KEY` — Datadog API key (for `datadog-ci`)
- `DEPLOY_REPO_TOKEN` — PAT with `contents:write` scope on the `svc-hello-deploy` repo
