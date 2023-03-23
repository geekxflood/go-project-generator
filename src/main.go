package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "project-init",
	Short: "Project Init is a Golang project initializer",
	Run: func(cmd *cobra.Command, args []string) {
		projectName, _ := cmd.Flags().GetString("project")
		createProject(projectName)
	},
}

func main() {
	rootCmd.PersistentFlags().StringP("project", "p", "", "Project name")
	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createProject(projectPath string) {
	if projectPath == "" {
		log.Fatal("Project path must be provided with --project or -p flag")
	}

	// Create project directories
	os.MkdirAll(fmt.Sprintf("%s/cmd", projectPath), 0755)
	os.MkdirAll(fmt.Sprintf("%s/pkg", projectPath), 0755)

	// Create .gitignore file
	gitignore := []byte("*.log\n*.swp\n*~\n*.out\n*.exe\n*.test\nvendor/\n")
	os.WriteFile(fmt.Sprintf("%s/.gitignore", projectPath), gitignore, 0644)

	// Create README.md file
	readme := []byte("# " + filepath.Base(projectPath))
	os.WriteFile(fmt.Sprintf("%s/README.md", projectPath), readme, 0644)

	// Create Makefile
	makefile := []byte("IMAGE_NAME = $(shell folder_name=$$(basename \"$$(pwd)\"); camel_case_name=$$(echo \"$$folder_name\" | awk -F' ' '{for(i=1;i<=NF;++i)printf \"%s\",toupper(substr($$i,1,1))tolower(substr($$i,2))}'); echo \"$$camel_case_name\")\n\n.PHONY: all build clean install build_image\n\nall: build clean\n\nbuild:\n\tmkdir -p binary\n\tgo build -a  \\\n\t\t-gcflags=all=\"-l -B\" \\\n\t\t-ldflags=\"-w -s\" \\\n\t\t-o binary/$(IMAGE_NAME) \\\n\t\t./...\n\ninstall: build\n\tsudo cp binary/$(IMAGE_NAME) /usr/local/bin\n\nclean:\n\trm -rf binary\n\nbuild_image: build\n\tdocker build --build-arg BINARY_NAME=$(IMAGE_NAME) -t $(IMAGE_NAME)-image .\n")
	os.WriteFile(fmt.Sprintf("%s/Makefile", projectPath), makefile, 0644)

	// Create Dockerfile
	dockerfile := []byte("FROM golang:1.20-alpine AS builder\n\nARG ARCH=amd64\nARG BINARY_NAME\n\nENV GOROOT /usr/local/go\nENV GOPATH /go\nENV PATH $GOPATH/bin:$GOROOT/bin:$PATH\nENV GO_VERSION 1.20\nENV GO111MODULE on\nENV CGO_ENABLED=0\n\nWORKDIR /go/src/\nCOPY . .\nRUN apk update && apk add make git\nRUN go get ./...\nRUN mkdir /go/src/build\nRUN go build -a -gcflags=all=\"-l -B\" -ldflags=\"-w -s\" -o build/$${BINARY_NAME} ./...\n\nFROM alpine:3.17\n\nARG BINARY_NAME\n\nCOPY --from=builder /go/src/build/$${BINARY_NAME} /usr/local/bin/$${BINARY_NAME}\nCMD [\"/usr/local/bin/$${BINARY_NAME}\"]\nEXPOSE 9740\n")
	os.WriteFile(fmt.Sprintf("%s/Dockerfile", projectPath), dockerfile, 0644)

	// Create .dockerignore file
	dockerignore := []byte(".env\n.gitlab-ci.yml\n.dockerignore\n.gcloudignore\n.gitignore\n.github/\n.gitlab/\n.git/\n*.md\nbuild/\ndoc/\ndashboard/\nvendor/\nLICENSE\n")
	os.WriteFile(fmt.Sprintf("%s/.dockerignore", projectPath), dockerignore, 0644)

	// Create src/main.go folder and file
	os.MkdirAll(fmt.Sprintf("%s/src", projectPath), 0755)
	main := []byte("package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tfmt.Println(\"Hello World\")\n}\n")
	os.WriteFile(fmt.Sprintf("%s/src/main.go", projectPath), main, 0644)

	// Run go mod init
	cmd := exec.Command("go", "mod", "init", filepath.Base(projectPath))
	cmd.Dir = projectPath
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Project %s has been initialized\n", projectPath)
}
