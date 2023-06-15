@TODO followed the https://github.com/golang/go/wiki/CodeReviewComments
@TODO explain about performance in read and write in the file storage based solution, use jsonl?


Tech debts
- generalize errors
- be possible to generalize file storage type (yaml, json) using interface to communicate with FS
- Encapsulate fmt.Printf


Focused in legibility, not hard performance (avoid branching, etc)

Didnt used filewriter to storage.local propositally

Only helm v3
Check if inside pod

Didnt have time for finalize the get of yq .a.b.c images from values.yaml.
    Tried recursive search in manifests
    Render with dryrun and get images: %s

    Generate container images during chart addition and store it?

# README

This repository contains a codebase organized in a specific directory structure and follows certain design patterns. Here's an overview of the code organization and design patterns present in this project:

## Directory Structure

- `cmd`: This directory contains the command-line interface (CLI) code. Each `.go` file represents a specific command or subcommand of the CLI. For example, `add.go` handles the "add" command, `images.go` handles the "images" command, and so on.
- `internal`: This directory holds internal packages and modules of the application.
  - `models`: Contains models representing different entities used in the application. `github.go` and `helm.go` are examples of such models.
  - `repository`: Contains the repository layer responsible for data access and manipulation.
    - `container_images`: Handles container image-related operations.
    - `file_reader`: Provides functionality to read files. It has subdirectories for different file systems like `filesystem` and `github`.
    - `file_writer`: Handles writing data to files. It has a subdirectory for YAML file handling.
    - `helm`: Deals with Helm-related operations.
    - `kubernetes`: Provides functionality related to Kubernetes, including an `internal.go` file that holds internal utility functions.
    - `storage`: Handles storage-related operations. The `local` subdirectory contains code specific to local storage, including `internal.go`, `local.go`, and `models.go`.
  - `service`: Contains the service layer of the application. Each `.go` file represents a specific service. For example, `add.go` represents the "add" service, `images.go` represents the "images" service, and so on.
- `go.mod` and `go.sum`: These files are related to Go modules and dependencies management.
- `main.go`: This is the entry point of the application, where the CLI commands are registered and executed.
- `LICENSE`: This file contains the license information for the codebase.
- `README.md`: The file you are currently reading, which provides an overview of the code organization and design patterns.

## Design Patterns

Based on the directory structure and the information provided, it's difficult to determine specific design patterns used in this codebase. However, the organization of the code into separate layers (CLI, models, repository, service) suggests a possible layered architecture or separation of concerns. Additionally, the presence of models, repositories, and services indicates a potential implementation of the repository pattern and service layer pattern for handling data access and business logic separately.

It's important to note that without further information or code analysis, the design patterns used cannot be definitively identified. The information provided in this README is a general overview based on the directory structure and common design patterns used in software development.
