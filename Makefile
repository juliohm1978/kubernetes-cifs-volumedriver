TAGNAME = juliohm/kubernetes-cifs-volumedriver-installer
VERSION = 0.5

build: Dockerfile
	docker build -t $(TAGNAME):$(VERSION) .

push:
	docker push $(TAGNAME):$(VERSION)
