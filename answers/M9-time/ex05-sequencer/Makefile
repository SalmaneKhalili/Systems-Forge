PATH := $(HOME)/.local/go/bin:$(PATH)
GO   ?= go

all: srv
srv:
	$(GO) build -o srv .
clean:
	rm -f srv
fclean: clean
re: fclean all
.PHONY: all clean fclean re