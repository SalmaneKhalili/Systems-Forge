PATH := $(HOME)/.local/go/bin:$(PATH)
GO   ?= go

all: node
node:
	$(GO) build -o node .
clean:
	rm -f node
fclean: clean
re: fclean all
.PHONY: all clean fclean re