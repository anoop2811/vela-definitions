# vela-definitions

A collection of KubeVela X-Definitions.

## Overview

This module contains Go-based KubeVela X-Definitions that can be applied to any KubeVela cluster.

## Directory Structure

- **components/** - ComponentDefinitions for workload types
- **traits/** - TraitDefinitions for operational behaviors
- **policies/** - PolicyDefinitions for application policies
- **workflows/** - WorkflowStepDefinitions for delivery workflows

## Usage

### Apply all definitions

```bash
vela def apply-module github.com/anoop2811/vela-definitions
```

### List definitions

```bash
vela def list-module github.com/anoop2811/vela-definitions
```

### Validate definitions

```bash
vela def validate-module github.com/anoop2811/vela-definitions
```

### Apply with namespace

```bash
vela def apply-module github.com/anoop2811/vela-definitions --namespace my-namespace
```

### Dry-run (preview without applying)

```bash
vela def apply-module github.com/anoop2811/vela-definitions --dry-run
```

## Adding New Definitions

1. Create a new Go file in the appropriate directory
2. Use the defkit package to define your component/trait/policy/workflow-step
3. Run `go mod tidy` to update dependencies
4. Validate with `vela def validate-module .`

Example component definition:

```go
package components

import "github.com/oam-dev/kubevela/pkg/definition/defkit"

func init() {
    defkit.Register(MyComponent)
}

func MyComponent() *defkit.ComponentDefinition {
    image := defkit.Param("image", "Container image").Required().String()

    return defkit.Component("my-component", "1.0.0").
        Description("My custom component").
        WithParameter(image).
        Output(defkit.K8sResource("deployment", "apps/v1", "Deployment").
            Set("metadata.name", defkit.Context("name")).
            Set("spec.template.spec.containers[0].image", image))
}
```

## Version

v0.1.0

## License

Apache License 2.0
