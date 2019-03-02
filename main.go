package main

import (
	"flag"
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"github.com/pkg/errors"
	"bufio"
	"io"
	"net/url"
)

type Config struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind string
	CurrentContext string `yaml:"current-context"`
	Clusters []*ClusterEntry
	Contexts []*ContextEntry
	Users []*UserEntry
}

type UserEntry struct {
	User *User
	Name string
}

type User struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData string `yaml:"client-key-data"`
}

type ContextEntry struct {
	Context *Context
	Name string
}

type Context struct {
	Cluster string
	User string
}

type ClusterEntry struct {
	Cluster *Cluster
	Name string
}

type Cluster struct {
	CAData string `yaml:"certificate-authority-data"`
	Server string
}

func main() {
	var kube string
	flag.StringVar(&kube, "k", "", "Kubernetes config folder.")
	flag.Parse()

	if kube == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("The program is intended to work with pipes.")
		fmt.Println("Usage: ssh user@my-server \"cat ~/.kube/config\" | metalogin")
		os.Exit(1)
	}

	srcConfigBytes, err := readFromInput()
	if err != nil {
		errors.Wrapf(err, "failed to read configuration from the input stream")
		fmt.Print(err)
		os.Exit(1)
	}

	srcConfig := Config{}
	err = yaml.Unmarshal(srcConfigBytes, &srcConfig)
	if err != nil {
		errors.Wrapf(err, "failed to unmarshal configuration from the input stream")
		fmt.Print(err)
		os.Exit(1)
	}

	kubeconfig := filepath.Join(kube, "config")
	dstConfigBytes, err := ioutil.ReadFile(kubeconfig)
	if err != nil {
		errors.Wrapf(err, "failed to read %s", kubeconfig)
		fmt.Print(err)
		os.Exit(1)
	}

	dstConfig := Config{}
	err = yaml.Unmarshal(dstConfigBytes, &dstConfig)
	if err != nil {
		errors.Wrapf(err, "failed to unmarshal %s", kubeconfig)
		fmt.Print(err)
		os.Exit(1)
	}

// Update cluster
	srcCluster := getClusterEntry(srcConfig, "kubernetes")
	if srcCluster == nil {
		fmt.Printf("cluster %s is not defined", "kubernetes")
		os.Exit(1)
	}

	server, err := url.Parse(srcCluster.Cluster.Server)
	if err != nil {
		errors.Wrapf(err, "server name `%s` is not a valid URL", srcCluster.Cluster.Server)
		fmt.Print(err)
		os.Exit(1)
	}

	clusterName := "kubernetes-" + server.Hostname()

	dstCluster := getClusterEntry(dstConfig, clusterName)
	if dstCluster == nil {
		dstCluster = &ClusterEntry{ Name: clusterName}
		dstConfig.Clusters = append(dstConfig.Clusters, dstCluster)
	}
	dstCluster.Cluster = srcCluster.Cluster

// Update user
	srcUser := getUserEntry(srcConfig, "kubernetes-admin")

	userName := "kubernetes-admin-" + server.Hostname()

	dstUser := getUserEntry(dstConfig, userName)
	if dstUser == nil {
		dstUser = &UserEntry{ Name: userName}
		dstConfig.Users = append(dstConfig.Users, dstUser)
	}
	dstUser.User = srcUser.User

// Update context
	contextName := userName + "@" + clusterName

	dstContext := getContextEntry(dstConfig, contextName)
	if dstContext == nil {
		dstContext = &ContextEntry{ Name: contextName }
		dstConfig.Contexts = append(dstConfig.Contexts, dstContext)
	}

	dstContext.Context.User = userName
	dstContext.Context.Cluster = clusterName

// Set current context
	dstConfig.CurrentContext = contextName

	b, err := yaml.Marshal(dstConfig)
	if err != nil {
		errors.Wrapf(err, "failed to marshal the resulted configuration")
		fmt.Print(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(kubeconfig, b, 0644)
	if err != nil {
		errors.Wrapf(err, "failed to write the resulted configuration")
		fmt.Print(err)
		os.Exit(1)
	}
}

func getClusterEntry(config Config, clusterName string) *ClusterEntry {
	for _, cluster := range config.Clusters {
		if cluster.Name == clusterName {
			return cluster
		}
	}
	return nil
}

func getUserEntry(config Config, userName string) *UserEntry {
	for _, user := range config.Users {
		if user.Name == userName {
			return user
		}
	}
	return nil
}

func getContextEntry(config Config, contextName string) *ContextEntry {
	for _, context := range config.Contexts {
		if context.Name == contextName {
			return context
		}
	}
	return nil
}

func readFromInput() ([]byte, error) {
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