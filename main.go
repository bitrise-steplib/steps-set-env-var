package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/log"
)

// ConfigsModel ...
type ConfigsModel struct {
	Value           string
	DestinationKeys string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		Value:           os.Getenv("value"),
		DestinationKeys: os.Getenv("destination_keys"),
	}
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")
	log.Printf("- Value: %s", configs.Value)
	log.Printf("- DestinationKeys: %s", configs.DestinationKeys)
}

func (configs ConfigsModel) validate() error {
	if configs.DestinationKeys == "" {
		return errors.New("no DestinationKeys parameter specified")
	}
	return nil
}

func main() {
	configs := createConfigsModelFromEnvs()

	fmt.Println()
	configs.print()
	fmt.Println()

	if err := configs.validate(); err != nil {
		log.Errorf("Issue with input: %s", err)
		os.Exit(1)
	}

	if configs.Value == "" {
		log.Warnf("The value of environment variable (%s) is empty", configs.Value)
		fmt.Println()
	}

	log.Infof("Exported environment variables:")
	for _, dest := range strings.Split(configs.DestinationKeys, " ") {
		cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", dest, "--value", configs.Value).CombinedOutput()
		if err != nil {
			fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
			os.Exit(1)
		}
		fmt.Printf("- %s: %s\n", dest, configs.Value)
	}
}
