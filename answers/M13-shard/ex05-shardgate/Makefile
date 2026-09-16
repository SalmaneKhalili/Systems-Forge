PATH := $(HOME)/.local/go/bin:$(PATH)
GO   ?= go

all: gate
gate:
	$(GO) build -o gate .
clean:
	rm -f gate
fclean: clean
re: fclean all
.PHONY: all clean fclean re