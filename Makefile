TAGNAME=juliohm/kubernetes-cifs-volumedriver-installer
VERSION=2.3-beta
DOCKER_CLI_EXPERIMENTAL=enabled
PLATFORMS=linux/amd64,linux/386,linux/arm,linux/arm64,linux/ppc64le

build:
	go build -a -installsuffix cgo

test:
	go test

docker:
	docker buildx build \
		-t $(TAGNAME):$(VERSION) \
		--progress plain \
		--platform=$(PLATFORMS) \
		.

push:
	docker buildx build \
		--push \
		-t $(TAGNAME):$(VERSION) \
		--progress plain \
		--platform=$(PLATFORMS) \
		.

install:
	kubectl apply -f install.yaml

delete:
	kubectl delete -f install.yaml

clean:
	rm -fr kubernetes-cifs-volumedriver
