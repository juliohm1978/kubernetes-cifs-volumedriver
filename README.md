# Kubernetes CIFS Volume Driver

[![nodesource/node](https://dockeri.co/image/juliohm/kubernetes-cifs-volumedriver-installer)](https://registry.hub.docker.com/u/juliohm/kubernetes-cifs-volumedriver-installer/)

## Important note

As of September 2023, personally, I have had no time to maintain this repo. It hasn't had any updates for some time now, and my current work schedule and priorities does not allow me the time for it.

Support for flexVolumes is not deprecated and will soon be removed from Kubernete's main releases.
https://github.com/juliohm1978/kubernetes-cifs-volumedriver/issues/36

I have no plans to implement this using the new CSI spec. Feel free to search for the best alternatives. I would like to thank the entire community for embracing this initial implementation. It's been a pleasure!

## Supported versions

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
* 1.17.x
* 1.18.x

## Documentation moved go Github Pages

Head over to <https://k8scifsvol.morimoto.net.br> for a complete installation guide and examples.

> WARNING: The documentation for this project is no longer hosted as subdomain of juliohm.com.br.
> Afer Nov/7th/2021, the domain has changed to <https://k8scifsvol.morimoto.net.br>
