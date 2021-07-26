build:
	go build
	cd frontend; npm install; npm run build

test:
	go fmt $(go list ./...)
	go vet $(go list ./...)
	go test ./...

