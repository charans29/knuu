![knuu-logo](./docs/knuu-logo.png)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/celestiaorg/knuu)
![GitHub Release](https://img.shields.io/github/v/release/celestiaorg/knuu)
[![CodeQL](https://github.com/celestiaorg/knuu/workflows/CodeQL/badge.svg)](https://github.com/celestiaorg/knuu/actions) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![OpenSSF Best Practices](https://bestpractices.coreinfrastructure.org/projects/7475/badge)](https://bestpractices.coreinfrastructure.org/projects/7475)

---

## Description

The goal of Knuu is to provide a framework for writing integration tests.
The framework is written in Go and is designed to be used in Go projects.
The idea is to provide a framework that uses the power of containers and Kubernetes without the test writer having to know the details of how to use them.

We invite you to explore our codebase, contribute, and join us in developing a framework to help projects write integration tests.

---

## Features

Knuu is designed around `Instances`, which you can create, start, control, communicate with other Instances, stop, and destroy.

Some of the features of Knuu are:

- Initialize an Instance from a Container/Docker image
- Configure startup commands
- Configure Networking
  - What ports to expose
  - Disable networking to simulate network outages
- Configure Storage
- Execute Commands
- Configure HW resources
- Create a pool of Instances and control them as a group
- Allow AddFile after Commit via ConfigMap
- Implement a TTL value for Pod cleanup
- Add a timeout variable
- See this issue for more upcoming features: [#91](https://github.com/celestiaorg/knuu/issues/91)

> If you have feedback on the framework, want to report a bug, or suggest an improvement, please create an issue [here](https://github.com/celestiaorg/knuu/issues/new/choose).

---

## Getting Started

This section will guide you on how to set up and run **Knuu**.

---

### Prerequisites

1. **Kubernetes cluster**: Set up access to a Kubernetes cluster using a `kubeconfig`.
   > In case you have no Kubernetes cluster running yet, you can get more information [here](https://kubernetes.io/docs/setup/).

2. **Docker**: Knuu uses Docker by default. If `KNUU_BUILDER` is not explicitly set to `kubernetes`, Docker is required to run Knuu.
   > You can install Docker by following the instructions [here](https://docs.docker.com/get-docker/).

### Writing Tests

The documentation you can find [here](https://pkg.go.dev/github.com/celestiaorg/knuu).

Simple example:

1. Add the following to your `go.mod` file:

    ```go
    require (
        github.com/celestiaorg/knuu v0.8.2
        github.com/stretchr/testify v1.8.4
    )
    ```

2. Run `go mod tidy` to download the dependencies.

3. Create a file called `main_test.go` with the following content to initialize Knuu:

    ```go
    package main

    import (
        "fmt"
        "os"
        "testing"
        "time"

        "github.com/celestiaorg/knuu/pkg/knuu"
    )

    func TestMain(m *testing.M) {
        err := knuu.Initialize()
        if err != nil {
           log.Fatalf("Error initializing knuu: %v:", err)
        }
        exitVal := m.Run()
        os.Exit(exitVal)
    }
    ```

4. Create a file called `example_test.go` with the following content:

    ```go
   package main

    import (
        "os"
        "testing"

        "github.com/celestiaorg/knuu/pkg/knuu"
        "github.com/stretchr/testify/assert"
    )

    func TestBasic(t *testing.T) {
        t.Parallel()
        // Setup

        instance, err := knuu.NewInstance("alpine")
        if err != nil {
            t.Fatalf("Error creating instance '%v':", err)
        }
        err = instance.SetImage("docker.io/alpine:latest")
        if err != nil {
            t.Fatalf("Error setting image: %v", err)
        }
        err = instance.SetCommand("sleep", "infinity")
        if err != nil {
            t.Fatalf("Error setting command: %v", err)
        }
        err = instance.Commit()
        if err != nil {
            t.Fatalf("Error committing instance: %v", err)
        }

        t.Cleanup(func() {
            // Cleanup
            if os.Getenv("KNUU_SKIP_CLEANUP") == "true" {
                t.Log("Skipping cleanup")
                return
            }

            err = instance.Destroy()
            if err != nil {
                t.Fatalf("Error destroying instance: %v", err)
            }
        })

        // Test logic

        err = instance.Start()
        if err != nil {
            t.Fatalf("Error starting instance: %v", err)
        }
        err = instance.WaitInstanceIsRunning()
        if err != nil {
            t.Fatalf("Error waiting for instance to be running: %v", err)
        }
        wget, err := instance.ExecuteCommand("echo", "Hello World!")
        if err != nil {
            t.Fatalf("Error executing command '%v':", err)
        }

        assert.Contains(t, wget, "Hello World!")
    }
    ```

You can find more examples in the following repositories:

- [celestiaorg/knuu](https://github.com/celestiaorg/knuu/e2e)
- [celestiaorg/celestia-app](https://github.com/celestiaorg/celestia-app/tree/main/test/e2e)

---

### Running Tests

You can use the built-in `go test` command to run the tests.

To run all tests in the current directory, you can run:

```shell
go test -v ./...
```

#### Environment Variables

You can set the following environment variables to change the behavior of knuu:

| Environment Variable | Description | Possible Values | Default |
| --- | --- | --- | --- |
| `KNUU_TIMEOUT` | The timeout for the tests. | Any valid duration | `60m` |
| `KNUU_BUILDER` | The builder to use for building images. | `docker`, `kubernetes` | `docker` |
| `LOG_LEVEL` | The debug level. | `debug`, `info`, `warn`, `error` | `info` |

---

# E2E

In the folder `e2e`, you will find some examples of how to use the [knuu](https://github.com/celestiaorg/knuu) Integration Test Framework.

## Setup

1. Install [Docker](https://docs.docker.com/get-docker/).

2. Set up access to a Kubernetes cluster using your `kubeconfig` and create the `test` namespace.

> **Note:** The used namespace can be changed by setting the `KNUU_NAMESPACE` environment variable.

## Write Tests

You can find the relevant documentation in the `pkg/knuu` package at: https://pkg.go.dev/github.com/celestiaorg/knuu

## Run

You can use the Makefile commands to easily target whatever test by setting the pkg, run, or count flags.

Targeting a directory

```shell
make test pkgs=./e2e/basic
```

Targeting a Test in a directory

```shell
make test pkgs=./e2e/basic run=TestJustThisTest
```

Run a test in a loop to debug

```shell
make test pkgs=./e2e/basic run=TestJustThisTest10Times count=10
```

---

## Contributing

We warmly welcome and appreciate contributions.

By participating in this project, you agree to abide by the [CNCF Code of Conduct](https://github.com/cncf/foundation/blob/main/code-of-conduct.md).

See the [Contributing Guide](./CONTRIBUTING.md) for more information.

To ensure that your contribution is working as expected, please run the tests in the `e2e` folder.

<!---
## Governance

[Describe the governance model for your project. Reference the GOVERNANCE.md file.]

## Adopters

[Provide information about the public adopters of your project. Reference the ADOPTERS.md file.]

## Security and Disclosure Information

See [SECURITY.md](SECURITY.md) for security and disclosure information.
--->

## Licensing

Knuu is licensed under the [Apache License 2.0](LICENSE).

<!---
## Contact

[Provide contact information for the project maintainers.]
--->
