CMD_NAME := gpta

compile:
	go build -o $(CMD_NAME) cmd/main.go

clean:
	rm $(CMD_NAME)