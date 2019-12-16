
build:
	GOOS=linux GOARCH=amd64 go build -o wso2am-controller ./cmd/controller/

run:
	./wso2am-controller -kubeconfig=$(HOME)/.kube/config

all: build run
















