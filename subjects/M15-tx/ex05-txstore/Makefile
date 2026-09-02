PATH := $(HOME)/.local/go/bin:$(PATH)
GO   ?= go

all: txstore
txstore:
	$(GO) build -o txstore .
clean:
	rm -f txstore
fclean: clean
re: fclean all
.PHONY: all clean fclean re
