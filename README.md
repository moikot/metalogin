# Metalogin

Metalogin simplifies access to your lovely bare-metal Kubernetes cluster. 
 
It receives necessary information from a bare-metal 
Kubernetes cluster via SSH and creates a context in your local ~/.kube/config.

Literally, after executing this:

```bash
ssh [user]@[cluster-IP] "cat ~/.kube/config" \
  | docker run -i --rm -v ~/.kube/:/kube moikot/metalogin -c /kube/config
```

You should be able to execute `kubectl get nodes`. No installation, no fiddling 
with certificates, contexts or users. This command does require Docker though. 

You can also build and run it locally if you have a Golang environment. 
In such case you need to execute

```bash
go get github.com/moikot/metalogin
ssh [user]@[cluster-IP] "cat ~/.kube/config" | ~/go/bin/metalogin -c ~/.kube/config

```

### What it actually does
1. First of all, it receives `config` file from your Kubernetes API node and 
deserializes it.
2. It tries to find a cluster record there which has name `kubernetes` this 
record corresponds to the bare-metal Kubernetes cluster.
3. It uses `server` field and assuming that it has a correct URI format, it 
tries to extract the server host name. Usually it's the IP address you used in
the SSH call.
4. Creates a cluster record in the local configuration with name 
`kubernetes-[host_name]` where `host_name` is the host name extracted on 
the previous step. All the other fields like `certificate-authority-data` 
and `server` are copied from source.
5. It tries to find a user record with name `kubernetes-admin` and if succeeds 
it creates a user record in the local configuration with name 
`kubernetes-admin-[host_name]` and then copies content of `client-certificate-data` 
and `client-key-data` fields from the found one.
6. Creates a context with name `kubernetes-admin-[host_name]@kubernetes-[host_name]` 
using previously created cluster and user. 
7. Sets the created context as the current one.
