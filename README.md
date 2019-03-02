# Metalogin

[![Build Status](https://travis-ci.org/moikot/metalogin.svg?branch=master)](https://travis-ci.org/moikot/metalogin)
[![Go Report Card](https://goreportcard.com/badge/github.com/moikot/metalogin)](https://goreportcard.com/report/github.com/moikot/metalogin)
[![Coverage Status](https://coveralls.io/repos/github/moikot/metalogin/badge.svg?branch=master)](https://coveralls.io/github/moikot/metalogin?branch=master)

Metalogin simplifies access to your lovely bare-metal Kubernetes cluster.

It receives necessary information from Kubernetes API node via SSH and creates
a context in your local ~/.kube/config.

Literally, after executing this:

```bash
ssh [user]@[cluster-IP] "cat ~/.kube/config" \
  | docker run -i --rm -v ~/.kube/:/kube moikot/metalogin -c /kube/config
```

You should be able to execute `kubectl get nodes` on your local machine.
No installation, no fiddling with certificates, contexts or users.
This command does require Docker though.

You can also build and run it locally if you have a Golang environment.
In such case you need to run the following commands:

```bash
go get github.com/moikot/metalogin
ssh [user]@[cluster-IP] "cat ~/.kube/config" | ~/go/bin/metalogin -c ~/.kube/config

```

### What it actually does
1. First of all, it receives `config` file from your Kubernetes API node and
deserializes it.
2. It tries to find a cluster record in it with name `kubernetes`. This
record corresponds to the bare-metal Kubernetes cluster.
3. It uses `server` field and assuming that it has a correct URI format, it
tries to extract the server host name. Usually it's the IP address you used in
the SSH call.
4. It creates a cluster record in the local configuration with name
`kubernetes-[host_name]` where `host_name` is the host name extracted on
the previous step. All the other fields like `certificate-authority-data`
and `server` are copied from the source record.
5. It tries to find a user record with name `kubernetes-admin` and when it succeeds
it creates a user record in the local configuration with name
`kubernetes-admin-[host_name]` and then copies content of `client-certificate-data`
and `client-key-data` fields from the source.
6. It creates a context with name `kubernetes-admin-[host_name]@kubernetes-[host_name]`
using previously created cluster and user.
7. Finally, it sets the created context as the current one.
