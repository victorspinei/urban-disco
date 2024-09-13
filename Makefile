BINARY_NAME=main.out

build:
	go build -o ./bin/${BINARY_NAME} ./main.go

run:
	go build -o ./bin/${BINARY_NAME} ./main.go
	./bin/${BINARY_NAME}
clean: 
	rm tmp/*
	rm bin/*