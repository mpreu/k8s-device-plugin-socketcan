format:
	go fmt ./pkg/...
	go vet ./pkg/...

build: format build-socketcan

build-docker:
	docker build -t mpreu/k8s-device-plugin-socketcan:latest -f ./build/package/Dockerfile .

build-socketcan:
	cd cmd/socketcan && go fmt && go build

deploy-ds:
	cd deployments && kubectl apply -f socketcan-ds.yml
	@echo "Run following command to get ds pod"
	@echo 'POD=$$(kubectl -n kube-system get pods -o go-template='\''{{range .items}}{{if .metadata.labels.name}}{{if eq .metadata.labels.name "k8s-device-plugin-socketcan"}}{{.metadata.name}}{{end}}{{end}}{{end}}'\'')'

remove-ds:
	cd deployments && kubectl delete -f socketcan-ds.yml

install:
	cd cmd/socketcan && go install

.PHONY: format build build-docker build-socketcan deploy-ds install