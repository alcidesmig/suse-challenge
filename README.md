# suse-cli-challenge

This project implements a technical challenge proposed by SUSE. It aims to provide a command-line interface (CLI) tool for managing and deploying Helm charts in a Kubernetes environment.

The time for the development of the challenge was 5 days, but due to other demands I was only able to work in parts of 2 and a half days.

## Challenge proposal

The initial proposal is described below.

### Instructions

Create a Golang CLI that will take as input a list of Helm Charts GitHub repo url or a local folder path (Ex: https://github.com/epinio/helm-charts/tree/main/chart/epinio or c:/epinio/helm-
charts/chart).

- The application should behave as described in the Expected Output section below.
- The application should be using go library clients/packages wherever possible.
- You can make any assumptions as you like and please state them if any.
- When `go build` is run at the root of your code base, it should provide us with the binary.

### Expected Output

The binary should behave in this manner.

```bash
Command - <binary> add <chart-location>
Result - Adds the helm chart information to the CLI’s internal list and a storage of your choice.

Command - <binary> index
Result – Generates Helm Repo Index file

Command - <binary> install chart <chart-name>
Result - Installs the given helm chart in the current Kubernetes cluster. Installation of the helm chart must happen inside a Kubernetes pod.

Command - <binary> images
Result - Provides a list of all the container images used in all the charts added.
```

> "Bonus Point: Write tests for your Golang application wherever deemed necessary."


## Assumptions

### "Must happen inside a pod"

The "Installation of the helm chart must happen inside a Kubernetes pod." requirement was understood as:
"It should only be possible to run the <bin> install chart command from within Kubernetes pods".

That said, there are several ways to check if the command was run from within a pod, but I haven't found any that can't be bypassed.

Thus, the existence of the `KUBERNETES_SERVICE_HOST` variable was used to verify the condition, being aware that this can be forged by simply assigning any value to the variable.

Tech debts
- generalize errors
- be possible to generalize file storage type (yaml, json) using interface to communicate with FS
- Encapsulate fmt.Printf


### Command "images"

I'm not satisfied with the current implementation of this requirement, but I haven't had enough time to implement it better.

The current implementation is very weak, and follows a non-recommended approach of calling commands directly from the system ("os/exec") to generate chart manifests and filter by image tags.

Ideas for more robust implementations:
- Do the same thing but using `helm.sh/helm/v3` packages instead of interacting directly with the shell
- Recursively search for images within chart files - and retrieve referenced values when they appear.
- Do a dry-run install and filter the images.

Another important point is that the generation happens every time the command is executed, and this can be optimized by porting the generation to the moment the chart is added.

### CLI's charts Storage

It was not specified where the chart information should be stored.

Therefore, an interface that can be implemented to change the storage location using dependency injection was created.

The current implementation, to avoid dependence on external systems, saves the information in the filesystem itself (in a yaml file).

```golang
type ChartStorageRepository interface {
	Init() error
	List(ctx context.Context) ([]models.ChartVersions, error)
	Get(ctx context.Context, chartName string) (map[string]models.ChartMetadata, error)
	Append(ctx context.Context, chart models.ChartMetadata, upsert bool) error
}
```

> Tech debt: choosing to use yaml implies in adding new charts to cause the entire data file to be read and written. This could be resolved using other formats such as `jsonlines` or hashing/tree structures.

### Performance

The performance-related interpretation was taken in two ways.
1. Code readability is better than performance micro-optimizations
2. This CLI is not highly sensitive to response time. More costly actions have not yet been optimally optimized.

### Helm version

The CLI was developed having in mind Helm 3.X.Y versions.

### Flags

Some optional flags have been added for certain commands, and these will be explained in the walkthrough.

### Windows and Linux

The code was developed aiming to having all functions running in Linux and Windows, but I haven't tested it deeply on windows because I use Linux as my main system.

### Tech Debts

- Use `FileWriterInterface` in Local Storage
- Optimize chart storage read/write (actually uses the entire file in all operations)
- Improve error messages and catching
- Encapsule stdout communication (actually uses `fmt.Print{f,ln}`)

## Directory Structure

- `cmd`: This directory contains the command-line interface (CLI) code. Each `.go` file represents a specific command or subcommand of the CLI. For example, `add.go` handles the "add" command, `images.go` handles the "images" command, and so on.
- `internal`: This directory holds internal packages and modules of the application.
  - `models`: Contains models representing different entities used in the application. `github.go` and `helm.go` are examples of such models.
  - `repository`: Contains the repository layer responsible for data access and manipulation.
    - `container_images`: Handles container image-related operations.
    - `file_reader`: Provides functionality to read files. It has subdirectories for different implementations like `filesystem` and `github`.
    - `file_writer`: Handles writing data to files. It has a subdirectory for YAML file handling (actually unique implementation).
    - `helm`: Deals with Helm-related operations.
    - `kubernetes`: Provides functionality related to Kubernetes.
    - `storage`: Handles local chart info storage-related operations.
  - `service`: Contains the service layer of the application. Each `.go` file represents a specific service. For example, `add.go` represents the "add" service, `images.go` represents the "images" service, and so on.
- `main.go`: This is the entry point of the application, where the CLI commands are registered and executed.

The directories can include an `internal.go` file that holds internal utility functions.

## Design Patterns

The organization of the code separate layers (CLI, models, repository, service) for layered architecture with separation of concerns.

The repository pattern and service layer pattern for handling data access and business logic separately were adopted.

The SOLID principles were followed reasonably, but not exhaustively for reasons of agility (implementation time).

All concrete implementations implements specified interfaces, so they can be easily mocked for unit tests with tools as [mockgen](https://pkg.go.dev/github.com/golang/mock/mockgen).


## Tests

I didn't have time to start the testing part, as much as I wanted to develop at least unit and integration tests.

## Walkthrough

The walkthrough presented here is a sequence of commands to test the project's functionalities and validate them.

### Menu

```bash
$ go build .
$ ./suse-cli-challenge -h
Technical exercise for showing Golang, Helm and Kubernetes knowledgement.

Usage:
  suse-cli-challenge [command]

Available Commands:
  add         Adds a Helm chart to the CLI's internal list and chosen storage.
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  images      Provides a list of container images used in all added charts.
  index       Generates a Helm repository index file.
  install     Installs a specific resource in the current Kubernetes cluster.

Flags:
  -h, --help   help for suse-cli-challenge

Use "suse-cli-challenge [command] --help" for more information about a command.
```

### Add
```bash
$ ./suse-cli-challenge add -h
This command allows you to add a Helm chart to the CLI's internal list for further processing.
You can provide the chart location as a GitHub repository URL or a local folder path.
The CLI will retrieve the chart from the specified location and store its information for future use.

Usage:
  suse-cli-challenge add [flags]

Flags:
  -h, --help     help for add
      --upsert   Update the chart local data if it already exists
```

```bash
$ ./suse-cli-challenge add https://github.com/epinio/helm-charts/tree/main/chart/application
Cloning repository for packaging chart...
Enumerating objects: 5215, done.
Counting objects: 100% (1117/1117), done.
Compressing objects: 100% (414/414), done.
Total 5215 (delta 925), reused 811 (delta 703), pack-reused 4098
Chart saved successfully.

$ ls -R ~/.config/suse-cli-challenge/
/home/alcides/.config/suse-cli-challenge/:
charts  charts.yaml

/home/alcides/.config/suse-cli-challenge/charts:
epinio-application-0.1.26.tgz
```

```bash
$ ./suse-cli-challenge add https://github.com/epinio/helm-charts/tree/main/chart/application
Cloning repository for packaging chart...
Enumerating objects: 5215, done.
Counting objects: 100% (1117/1117), done.
Compressing objects: 100% (414/414), done.
Total 5215 (delta 925), reused 811 (delta 703), pack-reused 4098
Chart saved successfully.

$ ls -R ~/.config/suse-cli-challenge/
/home/alcides/.config/suse-cli-challenge/:
charts  charts.yaml

/home/alcides/.config/suse-cli-challenge/charts:
epinio-application-0.1.26.tgz
```


```bash
$ git clone https://github.com/epinio/helm-charts.git
Cloning into 'helm-charts'...
remote: Enumerating objects: 5215, done.
remote: Counting objects: 100% (1036/1036), done.
remote: Compressing objects: 100% (401/401), done.
remote: Total 5215 (delta 855), reused 742 (delta 635), pack-reused 4179
Receiving objects: 100% (5215/5215), 2.37 MiB | 3.97 MiB/s, done.
Resolving deltas: 100% (3425/3425), done.

$ ./suse-cli-challenge add helm-charts/chart/upgrade-responder/
Chart saved successfully.
```

```bash
$ ./suse-cli-challenge add helm-charts/chart/upgrade-responder
Error: Version already exists. Please use --upsert if you want to override it.
Exception: error version already exists: (version 0.1.7) already exists

$ ./suse-cli-challenge add helm-charts/chart/upgrade-responder --upsert
Chart saved successfully.
```

### Index

```bash
./suse-cli-challenge index -h
With this command, you can generate a Helm repository index file.
The CLI will scan the internal list of Helm charts and create an index file that contains metadata about each chart.

Usage:
  suse-cli-challenge index [flags]
  suse-cli-challenge index [command]

Available Commands:
  print       Print the local helm index

Flags:
      --file string   Name for the file to write the charts index. (default "charts_index.yaml")
  -h, --help          help for index

Use "suse-cli-challenge index [command] --help" for more information about a command.
```


```bash
./suse-cli-challenge index prnt -h
With this command, you can generate a Helm repository index file.
The CLI will scan the internal list of Helm charts and create an index file that contains metadata about each chart.

Usage:
  suse-cli-challenge index [flags]
  suse-cli-challenge index [command]

Available Commands:
  print       Print the local helm index

Flags:
      --file string   Name for the file to write the charts index. (default "charts_index.yaml")
  -h, --help          help for index

Use "suse-cli-challenge index [command] --help" for more information about a command.

$ ./suse-cli-challenge index print -h
Print the local helm index to shell. This commands do not generate any files.

Usage:
  suse-cli-challenge index print [flags]

Flags:
  -h, --help   help for print
```

```bash
$ ./suse-cli-challenge index
Index successfully generated in "charts_index.yaml"
$ cat charts_index.yaml
- name: epinio-application
  data:
  - description: The helm chart for epinio applications to be deployed by
    version: 0.1.26
    url: https://github.com/epinio/helm-charts/tree/main/chart/application
    packaged_local_path: /home/alcides/.config/suse-cli-challenge/charts/epinio-application-0.1.26.tgz
- name: upgrade-responder
  data:
  - description: A Helm chart for Kubernetes
    version: 0.1.7
    url: /home/user/suse-challenge/helm-charts/chart/upgrade-responder
    packaged_local_path: /home/alcides/.config/suse-cli-challenge/charts/upgrade-responder-0.1.7.tgz

$ ./suse-cli-challenge index print
Chart Name: upgrade-responder
Versions:
- Version: 0.1.7
  Description: A Helm chart for Kubernetes
  URL: /home/alcides/suse-challenge/helm-charts/chart/upgrade-responder

Chart Name: epinio-application
Versions:
- Version: 0.1.26
  Description: The helm chart for epinio applications to be deployed by
  URL: https://github.com/epinio/helm-charts/tree/main/chart/application
```

### Images

```bash
$ ./suse-cli-challenge images
Chart: epinio-application
- Version: 0.1.26
  Images:

Chart: upgrade-responder
- Version: 0.1.7
  Images:
  - bats/bats:v1.4.1
  - curlimages/curl:7.85.0
  - grafana/grafana:9.4.7
  - influxdb:1.8.10-alpine
  - longhornio/upgrade-responder:v0.1.5
  - quay.io/kiwigrid/k8s-sidecar:1.22.0
```

### Install

```bash
$ ./suse-cli-challenge install chart epinio-application
To continue, please connect to a Kubernetes pod and execute the commands inside it. The installation process must be performed within the pod environment.
exit status 1
```

```bash
$ ./suse-cli-challenge chart inexistent-chart
No charts were found with the name "inexistent-chart". Please use the "add" command to add it.
exit status 1
```

```bash
$ ./suse-cli-challenge install chart epinio-application
2023/06/15 15:37:54 creating 2 resource(s)
Error: Failed while installing chart.
Exception: error installing chart: 1 error occurred:
        * Deployment.apps "rplaceholder-ff55435345834a3fe224936776c2aa15f6ed5358" is invalid: spec.template.spec.containers[0].image: Required value
```

Because the chart need additional parameters, we need to provide one `values.yaml`. This can be done through the `--values` command.

```
$ cat /tmp/values.yaml
epinio:
  tlsIssuer: "x"
  ingress: "x"
  appName: "placeholder"
  replicaCount: "1"
  stageID: "999"
  imageURL: "ubuntu"
  username: "user"
  routes: []
  env: []
  configurations: []
  start: []
$ ./suse-cli-challenge install chart epinio-application --values /tmp/values.yaml
2023/06/15 16:25:13 creating 2 resource(s)
```
