# kola: Kubernetes Operator Lifecycle Assistant

## Usage

```
Interact with OLM package manifests

Usage:
  kola [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List available packages
  show        Show details about a package
  subscribe   Generate a Subscription for a package
  version     Show command version

Flags:
      --cache-lifetime duration   Set cache lifetime (default 10m0s)
  -h, --help                      help for kola
  -k, --kubeconfig string         Path to kubernetes client configuration
      --no-cache                  Disable local caching of results
  -v, --verbose count             Increase output verbosity

Use "kola [command] --help" for more information about a command.
```

### List command

```
List available packages

Usage:
  kola list [flags]

Flags:
  -c, --catalog-source string   Match string in package catalog source
  -C, --certified               Match only certified packages
  -d, --description string      Match string in package description
  -g, --glob                    Arguments are glob patterns instead of substrings
  -h, --help                    help for list
  -m, --install-mode string     Match package supported install mode
  -w, --keyword strings         Match package keyword

Global Flags:
      --cache-lifetime duration   Set cache lifetime (default 10m0s)
  -k, --kubeconfig string         Path to kubernetes client configuration
      --no-cache                  Disable local caching of results
  -v, --verbose count             Increase output verbosity
```

### Show command

```
Show details about a package

Usage:
  kola show [flags]

Flags:
  -h, --help   help for show

Global Flags:
      --cache-lifetime duration   Set cache lifetime (default 10m0s)
  -k, --kubeconfig string         Path to kubernetes client configuration
      --no-cache                  Disable local caching of results
  -v, --verbose count             Increase output verbosity
```

### Subscribe command

```
Generate a Subscription for a package

Usage:
  kola subscribe [flags]

Aliases:
  subscribe, sub

Flags:
  -a, --approval string    Set install plan approval for subscription (default "Automatic")
  -c, --channel string     Set channel for subscription
  -h, --help               help for subscribe
  -n, --namespace string   Set namespace for subscription

Global Flags:
      --cache-lifetime duration   Set cache lifetime (default 10m0s)
  -k, --kubeconfig string         Path to kubernetes client configuration
      --no-cache                  Disable local caching of results
  -v, --verbose count             Increase output verbosity
```

## Examples

### List all packages relating to "gitops"

```
$ kola list -w gitops
2022/12/01 15:15:32 found 11 packages
flux
resource-locker-operator
tf-controller
argocd-operator
gitops-primer
gitwebhook-operator
flux
vault-config-operator
devopsinabox
openshift-gitops-operator
patch-operator
```

### Show package details

```
$ kola show external-secrets-operator
Name: external-secrets-operator
Catalog source: Community Operators (community-operators)
Publisher: Red Hat
Provider: External Secrets
Channels:
- alpha (external-secrets-operator.v0.7.0-rc1)
- stable (external-secrets-operator.v0.7.0-rc1)
```

### Subscribe to a package

```
$ kola subscribe external-secrets-operator
apiVersion: Subscription
kind: operators.coreos.com/v1alpha1
metadata:
  creationTimestamp: null
  name: external-secrets-operator
spec:
  channel: alpha
  installPlanApproval: Automatic
  name: external-secrets-operator
  source: community-operators
  sourceNamespace: openshift-marketplace
status:
  lastUpdated: null
```
