PATH := $(HOME)/.local/go/bin:$(PATH)
GO   ?= go

all: replica
replica:
	$(GO) build -o replica .
clean:
	rm -f replica
fclean: clean
re: fclean all
.PHONY: all clean fclean re