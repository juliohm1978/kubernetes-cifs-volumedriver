package main

import (
	"strings"
	"testing"
)

// Running without arguments should return an error message.
func TestMainWithoutArgs(t *testing.T) {
	args := []string{}
	msg := driverMain(args)
	if msg.Status != retStatFailure {
		t.Error()
	}
	if msg.Message != retMsgInsufficientArgs {
		t.Error()
	}

	args = []string{"/path/to/binary"}
	msg = driverMain(args)
	if msg.Status != retStatFailure {
		t.Error()
	}
	if msg.Message != retMsgInsufficientArgs {
		t.Error()
	}
}

// There is nothing to init at this point. Should return a sucess message.
func TestInit(t *testing.T) {
	args := []string{"/path/to/binary", "init"}
	msg := driverMain(args)
	if msg.Status != retStatSuccess {
		t.Error()
	}
	if msg.Capabilities.Attach {
		t.Error()
	}
}

func TestUnsupportedOperation(t *testing.T) {
	args := []string{"/path/to/binary", "i_do_not_exist"}
	msg := driverMain(args)
	if msg.Status != retStatNotSupported {
		t.Error()
	}
	if !strings.HasPrefix(msg.Message, retMsgUnsupportedOperation) {
		t.Error()
	}
}

func TestUnmountCmd(t *testing.T) {

	args := []string{"/path/to/binary", "unmount", "/mnt/point"}
	mountCmd := createUmountCmd(args)
	if mountCmd == nil {
		t.Error()
	}

	if mountCmd.Args[0] != "umount" {
		t.Error()
	}
	if mountCmd.Args[1] != "/mnt/point" {
		t.Error()
	}
}

// Teste with all possible arguments
func TestMountCmdComplete(t *testing.T) {

	jsonArgs := `{
	  "kubernetes.io/mounterArgs.FsGroup": "33",
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/readwrite": "rw",
	  "kubernetes.io/serviceAccount.name": "default",
	  "kubernetes.io/secret/domain": "ZG9tYWluMTIz",
	  "kubernetes.io/secret/username": "dXNlcjEyMw==",
	  "kubernetes.io/secret/password": "cGFzczEyMw==",
	  "opts": "domain=Foo",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	mountCmd := createMountCmd(args)
	if mountCmd == nil {
		t.Error()
	}

	if mountCmd.Args[0] != "mount" {
		t.Error()
	}
	if mountCmd.Args[1] != "-t" {
		t.Error()
	}
	if mountCmd.Args[2] != "cifs" {
		t.Error()
	}
	if mountCmd.Args[3] != "-o,uid=33,gid=33,rw,domain=domain123,username=user123,password=pass123,domain=Foo" {
		t.Error()
	}
	if mountCmd.Args[4] != "//fooserver123/test" {
		t.Error()
	}
	if mountCmd.Args[5] != "/mnt/point" {
		t.Error()
	}
}

// Simplest test, without any of:
// * fsGroup
// * Opts
// * Credentials
func TestMountCmdSimplest(t *testing.T) {

	jsonArgs := `{
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/serviceAccount.name": "default",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	mountCmd := createMountCmd(args)
	if mountCmd == nil {
		t.Error()
	}

	if mountCmd.Args[0] != "mount" {
		t.Error()
	}
	if mountCmd.Args[1] != "-t" {
		t.Error()
	}
	if mountCmd.Args[2] != "cifs" {
		t.Error()
	}
	if mountCmd.Args[3] != "//fooserver123/test" {
		t.Error()
	}
	if mountCmd.Args[4] != "/mnt/point" {
		t.Error()
	}
}

func TestMountCmdWithoutCredentials(t *testing.T) {

	jsonArgs := `{
	  "kubernetes.io/mounterArgs.FsGroup": "33",
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/readwrite": "rw",
	  "kubernetes.io/serviceAccount.name": "default",
	  "opts": "domain=Foo",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	mountCmd := createMountCmd(args)
	if mountCmd == nil {
		t.Error()
	}

	if mountCmd.Args[0] != "mount" {
		t.Error()
	}
	if mountCmd.Args[1] != "-t" {
		t.Error()
	}
	if mountCmd.Args[2] != "cifs" {
		t.Error()
	}
	if mountCmd.Args[3] != "-o,uid=33,gid=33,rw,domain=Foo" {
		t.Error()
	}
	if mountCmd.Args[4] != "//fooserver123/test" {
		t.Error()
	}
	if mountCmd.Args[5] != "/mnt/point" {
		t.Error()
	}
}

func TestMountCmdFsGroupLegacy(t *testing.T) {

	jsonArgs := `{
	  "kubernetes.io/fsGroup": "33",
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/readwrite": "rw",
	  "kubernetes.io/serviceAccount.name": "default",
	  "kubernetes.io/secret/domain": "ZG9tYWluMTIz",
	  "kubernetes.io/secret/username": "dXNlcjEyMw==",
	  "kubernetes.io/secret/password": "cGFzczEyMw==",
	  "opts": "domain=Foo",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	mountCmd := createMountCmd(args)
	if mountCmd == nil {
		t.Error()
	}

	if mountCmd.Args[0] != "mount" {
		t.Error()
	}
	if mountCmd.Args[1] != "-t" {
		t.Error()
	}
	if mountCmd.Args[2] != "cifs" {
		t.Error()
	}
	if mountCmd.Args[3] != "-o,uid=33,gid=33,rw,domain=domain123,username=user123,password=pass123,domain=Foo" {
		t.Error()
	}
	if mountCmd.Args[4] != "//fooserver123/test" {
		t.Error()
	}
	if mountCmd.Args[5] != "/mnt/point" {
		t.Error()
	}
}

func TestMountCmdWithoutFsGroup(t *testing.T) {

	jsonArgs := `{
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/readwrite": "rw",
	  "kubernetes.io/serviceAccount.name": "default",
	  "kubernetes.io/secret/domain": "ZG9tYWluMTIz",
	  "kubernetes.io/secret/username": "dXNlcjEyMw==",
	  "kubernetes.io/secret/password": "cGFzczEyMw==",
	  "opts": "domain=Foo",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	mountCmd := createMountCmd(args)
	if mountCmd == nil {
		t.Error()
	}

	if mountCmd.Args[0] != "mount" {
		t.Error()
	}
	if mountCmd.Args[1] != "-t" {
		t.Error()
	}
	if mountCmd.Args[2] != "cifs" {
		t.Error()
	}
	if mountCmd.Args[3] != "-o,rw,domain=domain123,username=user123,password=pass123,domain=Foo" {
		t.Error()
	}
	if mountCmd.Args[4] != "//fooserver123/test" {
		t.Error()
	}
	if mountCmd.Args[5] != "/mnt/point" {
		t.Error()
	}
}

func TestMountCmdInvalidCredendialDomain(t *testing.T) {

	// recover from panic, which is a good sign here
	defer func() {
		recover()
	}()

	jsonArgs := `{
	  "kubernetes.io/mounterArgs.FsGroup": "33",
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/readwrite": "rw",
	  "kubernetes.io/serviceAccount.name": "default",
	  "kubernetes.io/secret/domain": "INVALID_BASE64",
	  "kubernetes.io/secret/username": "dXNlcjEyMw==",
	  "kubernetes.io/secret/password": "cGFzczEyMw==",
	  "opts": "domain=Foo",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	createMountCmd(args)
	t.Error("Invalid base64 did not cause panic")
}

func TestMountCmdInvalidCredendialUser(t *testing.T) {

	// recover from panic, which is a good sign here
	defer func() {
		recover()
	}()

	jsonArgs := `{
	  "kubernetes.io/mounterArgs.FsGroup": "33",
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/readwrite": "rw",
	  "kubernetes.io/serviceAccount.name": "default",
	  "kubernetes.io/secret/domain": "ZG9tYWluMTIz",
	  "kubernetes.io/secret/username": "INVALID_BASE64",
	  "kubernetes.io/secret/password": "cGFzczEyMw==",
	  "opts": "domain=Foo",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	createMountCmd(args)
	t.Error("Invalid base64 did not cause panic")
}

func TestMountCmdInvalidCredendialPassword(t *testing.T) {

	// recover from panic, which is a good sign here
	defer func() {
		recover()
	}()

	jsonArgs := `{
	  "kubernetes.io/mounterArgs.FsGroup": "33",
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/readwrite": "rw",
	  "kubernetes.io/serviceAccount.name": "default",
	  "kubernetes.io/secret/domain": "ZG9tYWluMTIz",
	  "kubernetes.io/secret/username": "dXNlcjEyMw==",
	  "kubernetes.io/secret/password": "INVALID_BASE64",
	  "opts": "domain=Foo",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	createMountCmd(args)
	t.Error("Invalid base64 did not cause panic")
}

func TestMountCmdWithoutReadWrite(t *testing.T) {

	jsonArgs := `{
	  "kubernetes.io/mounterArgs.FsGroup": "33",
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/serviceAccount.name": "default",
	  "kubernetes.io/secret/domain": "ZG9tYWluMTIz",
	  "kubernetes.io/secret/username": "dXNlcjEyMw==",
	  "kubernetes.io/secret/password": "cGFzczEyMw==",
	  "opts": "domain=Foo",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	mountCmd := createMountCmd(args)
	if mountCmd == nil {
		t.Error()
	}

	if mountCmd.Args[0] != "mount" {
		t.Error()
	}
	if mountCmd.Args[1] != "-t" {
		t.Error()
	}
	if mountCmd.Args[2] != "cifs" {
		t.Error()
	}
	if mountCmd.Args[3] != "-o,uid=33,gid=33,domain=domain123,username=user123,password=pass123,domain=Foo" {
		t.Error()
	}
	if mountCmd.Args[4] != "//fooserver123/test" {
		t.Error()
	}
	if mountCmd.Args[5] != "/mnt/point" {
		t.Error()
	}
}

func TestMountCmdNoCredentialsAndNoOpts(t *testing.T) {

	jsonArgs := `{
	  "kubernetes.io/mounterArgs.FsGroup": "33",
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/readwrite": "rw",
	  "kubernetes.io/serviceAccount.name": "default",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	mountCmd := createMountCmd(args)
	if mountCmd == nil {
		t.Error()
	}

	if mountCmd.Args[0] != "mount" {
		t.Error()
	}
	if mountCmd.Args[1] != "-t" {
		t.Error()
	}
	if mountCmd.Args[2] != "cifs" {
		t.Error()
	}
	if mountCmd.Args[3] != "-o,uid=33,gid=33,rw" {
		t.Error()
	}
	if mountCmd.Args[4] != "//fooserver123/test" {
		t.Error()
	}
	if mountCmd.Args[5] != "/mnt/point" {
		t.Error()
	}
}

func TestMountCmdNoReadWrite(t *testing.T) {

	jsonArgs := `{
	  "kubernetes.io/mounterArgs.FsGroup": "33",
	  "kubernetes.io/fsType": "",
	  "kubernetes.io/pod.name": "nginx-deployment-549ddfb5fc-rnqk8",
	  "kubernetes.io/pod.namespace": "default",
	  "kubernetes.io/pod.uid": "bb6b2e46-c80d-4c86-920c-8e08736fa211",
	  "kubernetes.io/pvOrVolumeName": "test-volume",
	  "kubernetes.io/serviceAccount.name": "default",
	  "kubernetes.io/secret/domain": "ZG9tYWluMTIz",
	  "kubernetes.io/secret/username": "dXNlcjEyMw==",
	  "kubernetes.io/secret/password": "cGFzczEyMw==",
	  "opts": "domain=Foo",
	  "server": "fooserver123",
	  "share": "/test"
	}`

	args := []string{"/path/to/binary", "mount", "/mnt/point", jsonArgs}
	mountCmd := createMountCmd(args)
	if mountCmd == nil {
		t.Error()
	}

	if mountCmd.Args[0] != "mount" {
		t.Error()
	}
	if mountCmd.Args[1] != "-t" {
		t.Error()
	}
	if mountCmd.Args[2] != "cifs" {
		t.Error()
	}
	if mountCmd.Args[3] != "-o,uid=33,gid=33,domain=domain123,username=user123,password=pass123,domain=Foo" {
		t.Error()
	}
	if mountCmd.Args[4] != "//fooserver123/test" {
		t.Error()
	}
	if mountCmd.Args[5] != "/mnt/point" {
		t.Error()
	}
}
