PATH := $(HOME)/.local/go/bin:$(PATH)
GO   ?= go

all: watch
watch:
	$(GO) build -o watch .
clean:
	rm -f watch
fclean: clean
re: fclean all
.PHONY: all clean fclean re