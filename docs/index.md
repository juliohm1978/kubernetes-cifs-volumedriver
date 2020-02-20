# Kubernetes CIFS Volume Driver

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

> NOTE: Starting at v2.0, the driver has been fully reimplemented using [Go](https://golang.org/). As a full-fledged programming language, it provides a more robust solution and better error handling.
>
> Because Go can handle Json objects natively, the `jq` dependency is no longer necessary. The driver still relies on the `mount.cifs` binary, however, which is used to issue mount commands in the host OS. Aside from a different code base, all features should work the same as expected.
>
> The last implementation using Bash was v0.6. You can visit the tag to review the documentation for that release.
> <https://github.com/juliohm1978/kubernetes-cifs-volumedriver/tree/v0.6>

## Pre-requisites

On your Kubernetes nodes, simply install the dependency:

* `cifs-utils` because the host itself will do the mounting

```bash
sudo apt-get install -y cifs-utils
```

For CentOS:

```bash
yum -y install cifs-utils
```

## Manual Installation

Flexvolumes are very straight forward. The driver needs to be copied into a special volume plugin directory of your Kubernetes cluster.

For manuall installation, you will need to compile the code to create the binary executable. If you have Go installed, it should be easy as `make`.

```bash
go get github.com/juliohm1978/kubernetes-cifs-volumedriver
cd $GOPATH/src/github.com/juliohm1978/kubernetes-cifs-volumedriver

make

## if you want to be sure, run the test suite
make test

## Alternatively, you can use docker build to create the binary inside the docker image.
make docker
```

That should create the binary `kubernetes-cifs-volumedriver` that you can copy to your Kubernetes nodes.

```bash
## as root in all kubernetes nodes
mkdir -p /usr/libexec/kubernetes/kubelet-plugins/volume/exec/juliohm~cifs
cp -vr kubernetes-cifs-volumedriver /usr/libexec/kubernetes/kubelet-plugins/volume/exec/juliohm~cifs/cifs
chmod +x /usr/libexec/kubernetes/kubelet-plugins/volume/exec/juliohm~cifs/*
```

This procedure should be simple enough for testing purposes, so feel free to automate this in any way. Once the binary is copied and marked as executable, Kubelet should automatically pick it up and it should be working.

## DaemonSet Installation

When dealing with a large cluster, manually copying the driver to all hosts becomes inhuman. As proposed in [Flexvolume's documentation](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/storage/flexvolume-deployment.md#recommended-driver-deployment-method), the recommended driver deployment method is to have a DaemonSet install the driver cluster-wide automatically.

A Docker image [juliohm/kubernetes-cifs-volumedriver-intaller](https://hub.docker.com/r/juliohm/kubernetes-cifs-volumedriver-installer/) is available for this purpose, which can be deployed into a Kubernetes cluster using the `install.yaml` from this repository. The image is built `FROM busybox`, so the it's essentially very small.

The installer image allows you to install without the need to compile the project. Deploying the volume driver should be as easy as `make install`:

```bash
git clone https://github.com/juliohm1978/kubernetes-cifs-volumedriver.git
cd kubernetes-cifs-volumedriver

make install
```

The `install` target uses kubectl to create a privileged DaemonSet with pods that mount the host directory `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/` for installation. Check the output from the deployed containers to make sure it did not produce any errors. Crashing pods mean something went wrong.

> *NOTE*: This deployment does NOT install host dependencies, which still needs to be done manually on all hosts. See previous chapter *Pre-requisites*.

If you need to tweak or customize the installation, you can modify the `install.yaml` directly.

Installing is a one time job. So, once you have verified that it's completed, the DaemonSet can be safely removed.

```bash
make delete
```

## The Volume Plugin Directory

As of today with Kubernetes v1.16, the kubelet's default directory for volume plugins is `/usr/libexec/kubernetes/kubelet-plugins/volume/exec/`. This could be different if your installation changed this directory using the `--volume-plugin-dir` parameter.

A known example of this change is the installation provided by [Kubespray](https://github.com/kubernetes-incubator/kubespray), which at version v2.4.0 uses `/var/lib/kubelet/volume-plugins`.

Please, review the [kubelet command line parameters](https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/) (namely `--volume-plugin-dir`) and make sure it matches the directory where the driver will be installed.

You can modify `install.yaml` and change the field `spec.template.spec.volumes.hostPath.path` to the path used by your Kubernetes installation.

## Customizing the Vendor/Driver name

By default, the driver installation path is `$KUBELET_PLUGIN_DIRECTORY/juliohm~cifs/cifs`.

For some installations, you may need to change the vendor+driver name. Starting at v2.0, you can customize the vendor name/directory for your installation by tweaking `install.yaml` and defining `VENDOR` and `DRIVER` environment variables.

```yaml

## snippet ##

      containers:
        - image: juliohm/kubernetes-cifs-volumedriver-installer:2.0
          env:
            - name: VENDOR
              value: mycompany
            - name: DRIVER
              value: mycifs

## snippet ##

```

The example above will install the driver in the path `$KUBELET_PLUGIN_DIRECTORY/mycompany~mycifs/mycifs`. For the most part, changig the `VENDOR` variable should be enough to make your installation unique to your needs.

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

If you are using versions older than v0.5, you can still workaround by including these values in the `spec.flexVolume.options.opts` field of the PersistentVolume.

```yaml
## PV spec
spec:
  flexVolume:
    driver: juliohm/cifs
    options:
      opts: domain=Foo,uid=33,gid=33
```

## Troubleshooting and Known Issues

Because the driver is fundamentally a wrapper to `mount.cifs`, understanding what's happening at runtime can be challenging.

Remember to install the dependency: `cifs-utils`. Whatever your host OS should be, the driver attempts to issue `mount -t cifs...` to mount the volume. It should be installed on every node of the cluster.

Pay attention to the secret's `type` field, which **MUST** match the volume driver name. Otherwise, the secret values will not be passed to the driver.

If your Pod is stuck trying to mount a volume, check the events from the Kubernetes API.

```shell
kubectl describe pod POD_NAME
```

Events at the bottom of the Pod'd description will show errors collected by the driver while trying to mount the volume. So, you might see somethingn like this:

```text
Events:
  Type     Reason       Age                From               Message
  ----     ------       ----               ----               -------
  Normal   Scheduled    34s                default-scheduler  Successfully assigned default/nginx-deployment-58fc77b8db-hf47f to minikube
  Warning  FailedMount  18s (x6 over 34s)  kubelet, minikube  MountVolume.SetUp failed for volume "test-volume" : Couldn't get secret default/my-secret err: secrets "my-secret" not found
  Warning  FailedMount  2s                 kubelet, minikube  MountVolume.SetUp failed for volume "test-volume" : mount command failed, status: Failure, reason: Error: exit status 32
```

... where `Error: exit status 32` is the output of the mount command.

### Enabling Logs

Starting at v2.0, the Go implementation provides a basic logging mechanism. This could help you even further to understand why your volume fails to mount.

The driver attempts to write log messages to `/var/log/kubernetes-cifs-volumedriver.log`, but **only if that file already exists on disk and is writable**. Because log messages show all arguments issued to the `mount` command, password and secrets can be exposed. For that reason, logging is disabled by default. To enable, simply create the log file and wait for messages to come in.

```shell
## On the Kubernetes node where the Pod is scheduled
touch /var/log/kubernetes-cifs-volumedriver.log
tail -F /var/log/kubernetes-cifs-volumedriver.log
```

Once you no longer need to keep log messages, you can disable it by simply removing the log file.

```shell
rm /var/log/kubernetes-cifs-volumedriver.log
```

Log messages will look similar to this:

```text
# tail -F kubernetes-cifs-volumedriver.log
2019/11/17 04:46:05 [mount -t cifs -o,uid=333,gid=333,rw,username=pass123,password=user123,domain=Foo //10.0.0.114/publico /var/lib/kubelet/pods/7c92dd1d-5303-479e-8ff2-16713ae655c9/volumes/juliohm~cifs/test-volume]
2019/11/17 04:46:05 {"Status":"Failure","Message":"Error: exit status 32","Capabilities":{"Attach":false}}
```

Note that the complete `mount` command is on display, along with the response given to the Kubernetes API. That should give you a clear idea of what the driver is trying to do, and possibly some insight into the root cause of the problem.

### Kubelet Logs

While diagnosing issues, you might also want to checkout the output of the `kubelet` daemon. It rus on every node, in the host OS and is responsible for creating and destroying pods/containers.

The location for its log file can vary, depending on how you provisioned your cluster. For Ubuntu, `kubelet` is usually installed as system service. In that case, you can use `journalctl`.

```shell
journalctl -f -u kubelet
```

The output of `kubelet` may also give you clues and relevant error messages.
