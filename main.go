package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "", "Kubernetes config file.")
	flag.Parse()

	if configFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := Run(configFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func Run(configFile string) error {
	srcConfig, err := ReadSourceConfig()
	if err != nil {
		return errors.Wrapf(err, "failed to read source configuration")
	}

	dstConfig, err := ReadDestinationConfig(configFile)
	if err != nil {
		return errors.Wrapf(err, "failed to read destination configuration")
	}

	err = Merge(srcConfig, dstConfig)
	if err != nil {
		return errors.Wrapf(err, "failed to merge source configuration into destination one")
	}

	resConfigBytes, err := yaml.Marshal(dstConfig)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal result configuration")
	}

	var perm os.FileMode
	info, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			perm = 0600 // User read/write only
		} else {
			return errors.Wrapf(err, "failed to get information about `%s`", configFile)
		}
	} else {
		perm = info.Mode().Perm()
	}

	err = ioutil.WriteFile(configFile, resConfigBytes, perm)
	if err != nil {
		return errors.Wrapf(err, "failed to write result configuration")
	}

	return nil
}

func ReadSourceConfig() (*Config, error) {
	configBytes, err := ReadFromInput()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read configuration from the input stream")
	}

	if len(configBytes) == 0 {
		return nil, errors.New("no configuration in the input stream")
	}

	conf, err := Unmarshal(configBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal configuration from the input stream")
	}

	return conf, nil
}

func ReadFromInput() ([]byte, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		msg := "The program is intended to work with pipes.\n"
		msg += "Usage: ssh user@my-server \"cat ~/.kube/config\" | metalogin -c ~/.kube/config"
		return nil, errors.New(msg)
	}

	reader := bufio.NewReader(os.Stdin)
	var output []byte

	for {
		input, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		output = append(output, input)
	}

	return output, nil
}

func ReadDestinationConfig(configFile string) (*Config, error) {

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return NewConfig(), nil
	}

	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read `%s`", configFile)
	}

	conf, err := Unmarshal(configBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal `%s`", configFile)
	}

	return conf, nil
}
