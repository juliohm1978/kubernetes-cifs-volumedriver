TAGNAME = juliohm/kubernetes-cifs-volumedriver-installer
VERSION = 2.0-beta

build:
	go build

test:
	go test

docker: build test
	docker build -t $(TAGNAME):$(VERSION) .

push: docker
	docker push $(TAGNAME):$(VERSION)

install:
	kubectl apply -f install.yaml

delete:
	kubectl delete -f install.yaml

clean:
	rm -fr kubernetes-cifs-volumedriver