package main

import (
	"net/url"

	"github.com/pkg/errors"
)

const (
	ClusterName = "kubernetes"
	UserName    = "kubernetes-admin"
)

func Merge(srcConfig, dstConfig *Config) error {

	// Merge cluster
	srcCluster := srcConfig.GetClusterEntry(ClusterName)
	if srcCluster == nil {
		return errors.Errorf("cluster `%s` is not defined in source configuration", ClusterName)
	}

	server, err := url.Parse(srcCluster.Cluster.Server)
	if err != nil {
		return errors.Wrapf(err, "server name `%s` is not a valid URL", srcCluster.Cluster.Server)
	}

	clusterName := ClusterName + "-" + server.Hostname()

	dstCluster := dstConfig.GetClusterEntry(clusterName)
	if dstCluster == nil {
		dstCluster = &ClusterEntry{Name: clusterName}
		dstConfig.Clusters = append(dstConfig.Clusters, dstCluster)
	}
	dstCluster.Cluster = srcCluster.Cluster

	// Merge user
	srcUser := srcConfig.GetUserEntry(UserName)
	if srcUser == nil {
		return errors.Errorf("user `%s` is not defined in source configuration", UserName)
	}

	userName := UserName + "-" + server.Hostname()

	dstUser := dstConfig.GetUserEntry(userName)
	if dstUser == nil {
		dstUser = &UserEntry{Name: userName}
		dstConfig.Users = append(dstConfig.Users, dstUser)
	}
	dstUser.User = srcUser.User

	// Update context
	contextName := userName + "@" + clusterName

	dstContext := dstConfig.GetContextEntry(contextName)
	if dstContext == nil {
		dstContext = &ContextEntry{Name: contextName}
		dstConfig.Contexts = append(dstConfig.Contexts, dstContext)
	}

	dstContext.Context = &Context{
		User:    userName,
		Cluster: clusterName,
	}

	// Set current context
	dstConfig.CurrentContext = contextName

	return nil
}
