# Make variables
BIN_OUTPUT_FOLDER := bin
BIN_MAIN_FILE := dist

# Commands
default: 
	@go run cmd/xqtr/main.go

build: 
	@go build -o $(BIN_OUTPUT_FOLDER)/$(BIN_MAIN_FILE) -v cmd/xqtr/*.go

tests: 
	@go test ./... -cover ./... -v