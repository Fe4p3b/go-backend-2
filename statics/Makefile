#write here your username
USERNAME := fe4p3b
APP_NAME := statics
VERSION := 1.0.0

#write here path for your project
PROJECT := /home/almat/go/src/github.com/Fe4p3b/go-backend-2/statics
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
	kubectl -n geekbrains apply -f k8s/deployment.yaml

apply_service:
	kubectl -n geekbrains apply -f k8s/service.yaml

apply_ingress:
	kubectl -n geekbrains apply -f k8s/ingress.yaml


apply_rolling_update:
	kubectl -n geekbrains apply -f k8s/rolling-update/manifest.yaml

apply_fixed_deployment:
	kubectl -n geekbrains apply -f k8s/fixed-deployment/manifest.yaml

apply_blue-green:
	kubectl -n geekbrains apply -f k8s/blue-green/manifest-v1.yaml; \
	sleep 3; \
	kubectl -n geekbrains apply -f k8s/blue-green/manifest-v2.yaml; \
	kubectl -n geekbrains patch service statics -p '{"spec":{"selector":{"version":"v2.0.0"}}}';\
	kubectl -n geekbrains delete deploy statics-v1

apply_canary:
	kubectl -n geekbrains apply -f k8s/canary/manifest-v1.yaml; \
	sleep 3; \
	kubectl -n geekbrains apply -f k8s/canary/manifest-v2.yaml; \
	kubectl -n geekbrains scale --replicas=9 deploy statics-v1; \
	sleep 3; \
	kubectl -n geekbrains scale --replicas=10 deploy statics-v2; \
	kubectl -n geekbrains delete deploy statics-v1