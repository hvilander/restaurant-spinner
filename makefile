install :
	@go install github.com/a-h/templ/cmd/templ@latest

execute :
	@./bin/main

	
build : 
	@templ generate templates 
	@go build -o bin/main main.go

run : build execute
