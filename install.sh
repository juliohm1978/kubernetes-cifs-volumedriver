#!/bin/bash -e

cp -vr juliohm~cifs /usr/libexec/kubernetes/kubelet-plugins/volume/exec/
chmod +x /usr/libexec/kubernetes/kubelet-plugins/volume/exec/*
