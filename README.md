# Gitea Webhook Modifier (GWM)

This is an application that helps to modify [Gitea](https://gitea.com)
webhook request before it reaches a GitOps system (e.g.
[ArgoCD](https://argoproj.github.io/cd) or [FluxCD](https://fluxcd.io)).
It replaces certain string pattern specified by a regular expression
with another string anywhere in the JSON payload. By default, this is
used to replace the ingress URL (e.g. `http://git.localhost/`) with
in-cluster service URL (e.g. `http://gitea-http.gitea:3000/`).

## Local testing and development

Run the server:

```shell
GWM_LOG_LEVEL=trace GWM_FORWARD='http://localhost:8080' go run ./cmd/gwm/main.go
```

Send data:

```shell
curl -v -d '{"url":"http://git.localhost/myorg/tenant1"}' localhost:8080/webhook
```

## Usage

Install the GWM onto the cluster:

```shell
helm install -n gwm gwm oci://ghcr.io/jtyr/helm/gwm
```

Go to your local [Gitea Web UI](http://git.localhost) and define a GOGS webhook
to be sent to `http://gwm.gwm/webhook`.

## Configuration

There is several environment variables that can change the default behavior of
the application:

- `GWM_HOST` - Host name on which the HTTP server runs. Defaults to `0.0.0.0`.
- `GWM_PORT` - Port number on which the HTTP server runs. Defaults to `8080`.
- `GWM_SEARCH` - Regexp to search for a specific string to be replaced. Defaults to `^http://git.localhost/`.
- `GWM_REPLACE` - Regexp used to replace the found string. Defaults to `http://gitea-http.gitea:3000/`.
- `GWM_FORWARD` - URL to which to forward the modified request. Defaults to `http://argocd-server.argocd/api/webhook`.
- `GWM_LOG_LEVEL` - Log level. Possible values are `panic`, `fatal`, `error`, `warn`, `info`, `debug` and `trace`. Defaults to `info`.

## License

MIT

## Author

Jiri Tyr
