# gwm

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![AppVersion: 0.1.0](https://img.shields.io/badge/AppVersion-0.1.0-informational?style=flat-square)

Gitea Webhook Modifier (GWM) helps to modify Gitea webhook request before it reaches a GitOps system.

## Installation

```shell
helm upgrade --install --namespace gwm --create-namespace gwm ghcr.io/jtyr/helm/gwm
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Pod affinity. |
| autoscaling.enabled | bool | `false` | Whether the Horizontal Pod Autoscaling is enabled or not. |
| autoscaling.replicas.max | int | `10` | Maximum number of Pods to scale up. |
| autoscaling.replicas.min | int | `1` | Minimum number of Pods to scale down. |
| autoscaling.targetUtilization.cpu | int | `80` | Target average CPU utilization when the autoscaling is triggered. |
| autoscaling.targetUtilization.memory | string | `nil` | Target average memory utilization when the autoscaling is triggered. |
| fullnameOverride | string | `nil` | Full name override. |
| gwm.forward | string | `"http://argocd-server.argocd/api/webhook"` | URL where GWM forwards requests. |
| gwm.host | string | `"0.0.0.0"` | Host on which GWM will run. |
| gwm.logLevel | string | `"info"` | The GWM log level. |
| gwm.port | int | `80` | Port number on which GWM will run. |
| gwm.replace | string | `"http://gitea-http.gitea:3000/"` | Regexp used to replace the found string. |
| gwm.search | string | `"^http://git.localhost/"` | Regexp to search for a specific string to be replaced. |
| image.pullPolicy | string | `"IfNotPresent"` | Docker image pull policy. |
| image.repository | string | `"ghcr.io/jtyr/docker/gwm"` | Docker image. |
| image.tag | string | `""` | Docker image tag. Overrides the image tag whose default is the chart `appVersion`. |
| imagePullSecrets | list | `[]` | Image pull secrets. |
| ingress.annotations | object | `{}` | Ingress annotations. |
| ingress.className | string | `nil` | Ingress class name. Leave empty to use the default Ingress controller class. |
| ingress.enabled | bool | `false` |  |
| ingress.hosts[0].host | string | `"gwm.localhost"` | Ingress host name. |
| ingress.hosts[0].paths[0].path | string | `"/"` | Path where the ingress points. |
| ingress.hosts[0].paths[0].pathType | string | `"ImplementationSpecific"` | Type of the Ingress path. |
| ingress.tls | list | `[]` |  |
| nameOverride | string | `nil` | Name override. |
| nodeSelector | object | `{}` | Node selector. |
| podAnnotations | object | `{}` | Pod annotations. |
| podSecurityContext | object | `{}` | Pod security context. |
| replicas | int | `1` | Number of replicas. Ignored when `autoscaling.enabled: true`. |
| resources.limits.cpu | string | `"10m"` | Maximum amount of CPU the Pod is limited to use. |
| resources.limits.memory | string | `"20Mi"` | Maximum amount of memory the Pod is limited to use. |
| resources.requests.cpu | string | `"5m"` | Minimum amount of CPU the Pod requires to be scheduled. |
| resources.requests.memory | string | `"10Mi"` | Minimum amount of memory the Pod requires to be scheduled. |
| securityContext.allowPrivilegeEscalation | bool | `false` | Whether to allow privilege escallation. |
| securityContext.capabilities.drop | list | `["ALL"]` | List of capabilities to drop. |
| securityContext.readOnlyRootFilesystem | bool | `true` | Whether to mount the filesystem as read-only. |
| securityContext.runAsNonRoot | bool | `true` | Whether to run as non-root user. |
| securityContext.runAsUser | int | `65534` | User name or ID to to use to run the container. |
| service.port | int | `80` | Service port. |
| service.type | string | `"ClusterIP"` | Service type. |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account. |
| serviceAccount.create | bool | `true` | Whether a Service Account should be created. |
| serviceAccount.name | string | `nil` | The name of the Service account to use. If not set and `create: true`, a name is generated using the `fullname` template. |
| tolerations | list | `[]` | Pod tolerations. |

## Author

Jiri Tyr
