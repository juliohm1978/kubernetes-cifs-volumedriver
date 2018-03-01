# Kubernetes CIFS Volume Driver

A simple volume driver based on [Kubernetes' Flexvolume](https://github.com/kubernetes/community/edit/master/contributors/devel/flexvolume.md) that allows Kubernetes hosts to mount CIFS volumes (samba shares) into pods and containers.

It has been tested under Kubernetes 1.8.x and 1.9.x.

## Pre-requisites

There are few requirements for this to work. On your Kubernetes nodes, simply install a couple of dependencies: `cifs-utils` because the host itself will do the mounting and `jq` to parse json coming from the k8s api.

```bash
sudo apt-get install -y jq cifs-utils
```

For CentOS:

```bash
yum -y install jq cifs-utils
```

## Installation

This should not be a problem for different hosts, since it's very straight forward. The `juliohm1978~cifs` directory simply needs to be copied to `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/` and the script `cifs` needs permission to be executed.

The `install.sh` script included does that as an example:

```bash
## as root
cp -vr juliohm1978~cifs /usr/libexec/kubernetes/kubelet-plugins/volume/exec/
chmod +x /usr/libexec/kubernetes/kubelet-plugins/volume/exec/juliohm1978~cifs/*
```

Feel free to automate your installation in any way, shape or form. Once the script is copied and marked as executable, Kubelet should automatically pick it up and it should be working.

## Example of PersistentVolume

The following is an example of PersistentVolume that uses the volume driver.

```yaml
kind: PersistentVolume
metadata:
  name: mycifspv
spec:
  capacity:
    storage: 1Gi
  flexVolume:
    driver: juliohm1978/cifs
    options:
      opts: sec=ntlm,uid=1000
      server: my-cifs-host
      share: /MySharedDirectory
    secretRef:
      name: my-secret
```

Credentials are passed using a Secret, which can be declared as follows.

```yaml
apiVersion: v1
data:
  password: ###
  username: ###
kind: Secret
metadata:
  name: my-secret
type: juliohm1978/cifs
```

*NOTE*: Pay attention to the secret's `type` field, which MUST match the volume driver name. Otherwise the secret values will not be passed to the mount script.
