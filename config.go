package main

import (
	"gopkg.in/yaml.v2"
)

type Config struct {
	APIVersion     string `yaml:"apiVersion"`
	Kind           string
	CurrentContext string `yaml:"current-context"`
	Clusters       []*ClusterEntry
	Contexts       []*ContextEntry
	Users          []*UserEntry
}

type UserEntry struct {
	User *User
	Name string
}

type User struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}

type ContextEntry struct {
	Context *Context
	Name    string
}

type Context struct {
	Cluster string
	User    string
}

type ClusterEntry struct {
	Cluster *Cluster
	Name    string
}

type Cluster struct {
	CAData string `yaml:"certificate-authority-data"`
	Server string
}

func Unmarshal(data []byte) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) GetClusterEntry(clusterName string) *ClusterEntry {
	for _, cluster := range c.Clusters {
		if cluster.Name == clusterName {
			return cluster
		}
	}
	return nil
}

func (c *Config) GetUserEntry(userName string) *UserEntry {
	for _, user := range c.Users {
		if user.Name == userName {
			return user
		}
	}
	return nil
}

func (c *Config) GetContextEntry(contextName string) *ContextEntry {
	for _, context := range c.Contexts {
		if context.Name == contextName {
			return context
		}
	}
	return nil
}
