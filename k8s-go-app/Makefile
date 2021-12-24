#write here your username
USERNAME := fe4p3b
APP_NAME := k8s-go-app
VERSION := 1.0.0

#write here path for your project
PROJECT := /home/almat/go/src/github.com/Fe4p3b/go-backend-2/k8s-go-app
GIT_COMMIT := $(shell git rev-parse HEAD)


all: run

run:
	go install -ldflags="-X '$(PROJECT)/version.Version=$(VERSION)' \
	-X '$(PROJECT)/version.Commit=$(GIT_COMMIT)'" && $(APP_NAME)

build_container:
	docker build -t docker.io/$(USERNAME)/$(APP_NAME):$(VERSION) --build-arg GIT_COMMIT=$(GIT_COMMIT) --build-arg VERSION=$(VERSION) --build-arg PROJECT=$(PROJECT) .

push_container:
	docker push  docker.io/$(USERNAME)/$(APP_NAME):$(VERSION)

apply_deploy:
	kubectl apply -f k8s/deployment.yaml

apply_service:
	kubectl apply -f k8s/service.yaml

apply_ingress:
	kubectl apply -f k8s/ingress.yaml

