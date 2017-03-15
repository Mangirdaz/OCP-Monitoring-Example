test:
	go test -v `go list ./... | egrep -v /vendor/`

proto:
	cd config/ && protoc -I=./ --go_out=./ proto.proto

build-api: 
	cd backend-api/ && CGO_ENABLED=0 GOOS=linux go build -o api .

build-api-mon: 
	cd backend-api-monitor/ && CGO_ENABLED=0 GOOS=linux go build -o api .	

build-api-ext: 
	cd backend-external/ && CGO_ENABLED=0 GOOS=linux go build -o api .

build-ui: 
	cd frontend/ui/ && npm install && npm run build
	cd frontend && rm -rf statik/statik.go && ${GOPATH}/bin/statik -src=./ui/dist
	cd frontend && CGO_ENABLED=0 GOOS=linux go build -o fe .	

run-api: build-api
	./backend-api/api

run-api-ext: export API_PORT=8001
run-api-ext: build-api-ext 
	./backend-external/api		
	
run-api-mon: build-api-mon
	./backend-api-monitor/api	

run-consul:
	docker run -p 8400:8400 -p 8500:8500 -p 8600:53/udp -h node1 progrium/consul -server -bootstrap -ui-dir /ui

run-ui:
	cd frontend/ui/ && npm run dev 
	
docker-build-ui:
	cd ./frontend/ && docker build -t docker.io/mangirdas/ocp-example-fe:v0.5 -f Dockerfile .
	cd ./frontend/ && docker build -t mangirdas/ocp-example-fe:v0.5 -f Dockerfile .
docker-build-api:
	cd ./backend-api/ && docker build -t docker.io/mangirdas/ocp-example-api:v0.5 -f Dockerfile .
	cd ./backend-api/ && docker build -t mangirdas/ocp-example-api:v0.5 -f Dockerfile .
docker-build-api-ext:
	cd ./backend-external/ && docker build -t docker.io/mangirdas/ocp-example-api-ext:v0.5 -f Dockerfile .
	cd ./backend-external/ && docker build -t mangirdas/ocp-example-api-ext:v0.5 -f Dockerfile .
docker-build-api-mon:
	cd ./backend-api-monitor/ && docker build -t docker.io/mangirdas/ocp-example-api-mon:v0.5 -f Dockerfile .
	cd ./backend-api-monitor/ && docker build -t mangirdas/ocp-example-api-mon:v0.5 -f Dockerfile .
docker-push-ui:
	docker push docker.io/mangirdas/ocp-example-fe:v0.5
docker-push-api:
	docker push docker.io/mangirdas/ocp-example-api:v0.5
docker-push-api-ext:
	docker push docker.io/mangirdas/ocp-example-api-ext:v0.5
docker-push-api-mon:
	docker push docker.io/mangirdas/ocp-example-api-mon:v0.5

build: build-api build-ui build-api-ext build-api-mon docker-build-ui docker-build-api docker-build-api-ext docker-build-api-mon docker-push-ui docker-push-api docker-push-api-ext docker-push-api-mon
	