all: build

NAME = converter

build:
	@go build -o $(NAME) cmd/main.go

fclean:
	@rm $(NAME)