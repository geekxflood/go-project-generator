package main

import (
	"fmt"
	"io/ioutil"
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
	ioutil.WriteFile(fmt.Sprintf("%s/.gitignore", projectPath), gitignore, 0644)

	// Create Makefile
	makefile := []byte(`.PHONY: build

build:
	go build -o bin/ ./cmd/...
`)
	ioutil.WriteFile(fmt.Sprintf("%s/Makefile", projectPath), makefile, 0644)

	// Run go mod init
	cmd := exec.Command("go", "mod", "init", filepath.Base(projectPath))
	cmd.Dir = projectPath
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Project %s has been initialized\n", projectPath)
}
