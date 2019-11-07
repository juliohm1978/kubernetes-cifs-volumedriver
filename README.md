# Kubernetes CIFS Volume Driver

[![nodesource/node](http://dockeri.co/image/juliohm/kubernetes-cifs-volumedriver-installer)](https://registry.hub.docker.com/u/juliohm/kubernetes-cifs-volumedriver-installer/)

A simple volume driver based on [Kubernetes' Flexvolume](https://github.com/kubernetes/community/blob/master/contributors/devel/flexvolume.md) that allows Kubernetes hosts to mount CIFS volumes (samba shares) into pods and containers.

It has been tested under Kubernetes versions:

* 1.8.x
* 1.9.x
* 1.10.x
* 1.11.x
* 1.12.x
* 1.13.x
* 1.14.x
* 1.15.x
* 1.16.x

## Pre-requisites

On your Kubernetes nodes, simply install a couple of dependencies:

* `cifs-utils` because the host itself will do the mounting
* `jq` to parse json coming from the k8s api

```bash
sudo apt-get install -y jq cifs-utils
```

For CentOS:

```bash
yum -y install jq cifs-utils
```

## Manual Installation

Flexvolumes are very straight forward. The `juliohm~cifs` directory simply needs to be copied into the volume plugin directory of your Kubernetes cluster.

Below is an example:

```bash
## as root in all kubernetes nodes
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

> *NOTE*: This deployment does NOT install host dependencies, which still needs to be done manually on all hosts. See previous chapter *Pre-requisites*.

Once you have verified that installation was completed, the DaemonSet can be safely removed.

```bash
kubectl delete -f install.yaml
```

## The Volume Plugin Directory

As of today with Kubernetes v1.15, the kubelet's default directory for volume plugins is `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/`. This could be different if your installation changed this directory using the `--volume-plugin-dir` parameter.

A known example of this change is the installation provided by [Kubespray](https://github.com/kubernetes-incubator/kubespray), which at version v2.4.0 uses `/var/lib/kubelet/volume-plugins`.

Please, review the [kubelet command line parameters](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/) (namely `--volume-plugin-dir`) and make sure it matches the directory where the driver will be installed.

You can modify `install.yaml` and change the field `spec.template.spec.volumes.hostPath.path` to the path used by your Kubernetes installation.

## Example of PersistentVolume

The following is an example of PersistentVolume that uses the volume driver.

```yaml
apiVersion: v1
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
  accessModes:
    - ReadWriteMany
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

## Using `securityContext` to inform uid/gid parameters

Starting at version 0.5, the driver will also accept values coming from the Pod's `securityContext`.

For example, consider the following Deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
        volumeMounts:
          - name: test
            mountPath: /dados
      securityContext:
        runAsUser: 33
        runAsGroup: 33
        fsGroup: 33
      volumes:
        - name: test
          persistentVolumeClaim:
            claimName: test-claim
```

... which defines a `securityContext`.

```yaml
      securityContext:
        runAsUser: 33
        runAsGroup: 33
        fsGroup: 33
```

The value of `fsGroup` is passed to the volume driver, but previous versions would ignore that. It is now used to construct `uid` and `gid` parameters for the mount command.

If you are using versions older than 0.5, you can still workaround by including these values in the `spec.flexVolume.options.opts` field of the PersistentVolume.

```yaml
## PV spec
spec:
  flexVolume:
    driver: juliohm/cifs
    options:
      opts: domain=Foo,uid=33,gid=33
```

## Notes Failures and Known Issues

For most issues reported until now, the root cause was not related to the driver itself. Understanding what's happening at runtime can be challenging.

Remember to install the dependencies: `jq` and `cifs-utils`. These should be installed on every node of the cluster.

Pay attention to the secret's `type` field, which **MUST** match the volume driver name. Otherwise, the secret values will not be passed to the mount script.

Watching the kubelet's logs on the nodes where the pod is scheduled may help you diagnose some problems. More often than not, the driver fails because the `mount` command fails with a non-zero exit code.

Take note of the field `spec.flexVolume.options.opts` used in your PV and try to manually mount the volume on the same node where the pod is scheduled using the same options and credentials. Given the PV yaml in the example above, the driver would issue a command line similar to this:

```
mount -t cifs -o sec=ntlm,uid=1000,username=***,password=*** //my-cifs-host/MySharedDirectory /mnt/temp/dir
```

If that fails, it will likely give you an insight into the root cause of the problem.
