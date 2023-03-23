# Golang Project Initializer

This Golang project initializer creates a new project with the necessary folder structure, a `.gitignore` file, a `Makefile`, a `Dockerfile`, and a `.dockerignore` file. It also initializes a Go module for the project.

## Usage

1. Install the Golang project initializer:

```sh
go get ./...
```

2. Build the Golang project initializer:

```sh
make build
```

or directly install
    
```sh
make install
```

3. Run the Golang project initializer, providing the desired project name using the --project or -p flag:

```sh
ggp --project <project_name>
```

```sh
ggp --p <project_name>
```

This will create a new project in the specified directory with the following structure:

```md
<project_name>/
├── .gitignore
├── README.md
├── Makefile
├── Dockerfile
├── .dockerignore
├── src/
│   └── main.go
└── go.mod
```

- The Makefile includes targets for building, cleaning, installing, and building a Docker image for the project.
- The Dockerfile is a multi-stage build file that uses the golang:1.20-alpine and alpine:3.17 images.
