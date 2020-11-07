package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
)

// ArtifactEnvelope contains artifact content (as JSON) and relative original file path
type ArtifactEnvelope struct {
	Path    string      `json:"path"`
	Content interface{} `json:"content"`
}

var config *Config
var confPath string
var repositoryPath string
var silent bool

func init() {
	flag.BoolVar(&silent, "silent", false, "Disables logging")
	flag.StringVar(&repositoryPath, "repository", ".", "Path to a repository to check")
	flag.StringVar(&confPath, "conf", "conf.yml", "Path to the configuration file")
	flag.Parse()

	if silent {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetFlags(log.Ltime | log.Lmicroseconds)
	}

	cfg, err := NewConfig(confPath)
	if err != nil {
		panic(err)
	}
	config = cfg
}

func main() {
	log.Printf("Looking for artifacts in %v", repositoryPath)

	// Set current dir to the repository root
	err := os.Chdir(repositoryPath)
	if err != nil {
		panic(err)
	}

	input := make(map[string][]ArtifactEnvelope)
	for _, artifact := range config.Artifacts {
		var files []ArtifactEnvelope

		filePaths := getFilePaths(&artifact)

		for _, filePath := range filePaths {
			content := ConvertToJSON(artifact.Type, filePath)
			files = append(files, ArtifactEnvelope{
				Path:    filePath,
				Content: content,
			})
		}
		input[artifact.Name] = files
	}

	// Log output data
	log.Printf(describeResult(input))

	// Print ready input json to the stdout
	result, err := json.Marshal(input)
	fmt.Print(string(result))
}

func describeResult(result map[string][]ArtifactEnvelope) string {
	var builder strings.Builder
	builder.WriteString("Gathered artifacts:\n")
	for key, value := range result {
		fmt.Fprintf(&builder, "  %s:\n", key)
		for _, artifact := range value {
			fmt.Fprintf(&builder, "    - %s\n", artifact.Path)
		}
	}
	return builder.String()
}

func getFilePaths(artifact *Artifact) []string {

	matches, err := doublestar.Glob(artifact.SearchPattern)
	if err != nil {
		panic(err)
	}

	return matches
}
