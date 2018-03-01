# Kubernetes CIFS Volume Driver

A simple volume driver based on [Kubernetes' Flexvolume](https://github.com/kubernetes/community/blob/master/contributors/devel/flexvolume.md) that allows Kubernetes hosts to mount CIFS volumes (samba shares) into pods and containers.

It has been tested under Kubernetes 1.8.x and 1.9.x.

## Pre-requisites

On your Kubernetes nodes, simply install a couple of dependencies: `cifs-utils` because the host itself will do the mounting and `jq` to parse json coming from the k8s api.

```bash
sudo apt-get install -y jq cifs-utils
```

For CentOS:

```bash
yum -y install jq cifs-utils
```

## Manual Installation

Flexvolumes are very straight forward. The `juliohm~cifs` directory simply needs to be copied to `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/` and the script `cifs` needs permission to be executed.

Below is an example:

```bash
## as root
cp -vr juliohm~cifs /usr/libexec/kubernetes/kubelet-plugins/volume/exec/
chmod +x /usr/libexec/kubernetes/kubelet-plugins/volume/exec/juliohm~cifs/*
```

This procedure should be simple enough for testing purposes, so feel free to automate this in any way, shape or form. Once the script is copied and marked as executable, Kubelet should automatically pick it up and it should be working.

When dealing with a large cluster, manually copying the driver to all hosts becomes inhuman. For that and most cases in general, the DaemonSet installation should make things easier.

## DaemonSet Installation

As proposed in [Flexvolume's documentation](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/storage/flexvolume-deployment.md#recommended-driver-deployment-method), the recommended driver deployment method is to have a DaemonSet install the driver cluster-wide automatically.

A Docker image [juliohm/kubernetes-cifs-volumedriver-intaller](https://hub.docker.com/r/juliohm/kubernetes-cifs-volumedriver-installer/) is available for this purpose, which can be deployed into a Kubernetes cluster using the `install.yaml` from this repository. The image is built `FROM busybox`, so the it's essentially very small and slightly over 1MB.

Deploying the volume driver should be as easy as:

```bash
kubectl apply -f install.yaml
```

This creates a privileged DaemonSet with pods that mount the host directory `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/` internally as `/flexmnt` for installation. Check the output from the deployed containers to make sure it did not produce any errors. Crashing pods mean something went wrong.

*NOTE*: This deployment does NOT install host dependencies, which still needs to be done manually on all hosts. See previous chapter *Pre-requisites*.

Once you have verified that installation was completed, the DaemonSet can be safely removed.

```bash
kubectl delete -f install.yaml
```

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
    driver: juliohm/cifs
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
type: juliohm/cifs
```

*NOTE*: Pay attention to the secret's `type` field, which MUST match the volume driver name. Otherwise the secret values will not be passed to the mount script.
