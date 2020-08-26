run:
	go run ./cmd/controller

test:
	go vet ./...
	ginkgo -r

install:
	kapp deploy -c -a scm -f ./config/crd

uninstall:
	kapp delete -a scm

image:
	pack build kontinue/scm-controller

gen:
	controller-gen \
		object \
		paths=./api/v1
	controller-gen \
		crd \
		paths=./api/v1/
