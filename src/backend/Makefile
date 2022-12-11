run:
	docker-compose down --volumes
	docker build -f ./docker/app/Dockerfile . --tag app
	docker-compose up --build

build:
	go build -o app.out -v ./cmd/server/main.go

make-mocks:
	go generate ./...

generate-api:
	go install github.com/swaggo/swag/cmd/swag@v1.6.5
	swag init -g ./cmd/server/main.go -o docs
	
clean:
	rm -rf *.out *.exe *.html *.csv
