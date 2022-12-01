# kola: Kubernetes Operator Lifecycle Assistant

## Usage

```
Interact with OLM package manifests

Usage:
  kola [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  dump        Dump details about a package
  help        Help about any command
  list        List available packages
  show        Show details about a package
  subscribe   Generate a Subscription for a package
  version     Show command version

Flags:
  -h, --help                help for kola
  -k, --kubeconfig string   Path to kubernetes client configuration
  -v, --verbose count       Increase log verbosity

Use "kola [command] --help" for more information about a command.
```

### List command

```
List available packages

Usage:
  kola list [flags]

Flags:
  -c, --catalogSource string   Match string in package catalog source
  -C, --certified              Match only certified packages
  -d, --description string     Match string in package description
  -g, --glob                   Arguments are glob patterns instead of substrings
  -h, --help                   help for list
  -m, --installMode string     Match package supported install mode
  -w, --keyword strings        Match package keyword
  -U, --uncertified            Match only certified packages

Global Flags:
  -k, --kubeconfig string   Path to kubernetes client configuration
  -v, --verbose count       Increase log verbosity
```

### Show command

```
Show details about a package

Usage:
  kola show [flags]

Flags:
  -d, --description   Include description in output
  -h, --help          help for show

Global Flags:
  -k, --kubeconfig string   Path to kubernetes client configuration
  -v, --verbose count       Increase log verbosity
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
  -k, --kubeconfig string   Path to kubernetes client configuration
  -v, --verbose count       Increase log verbosity
```
