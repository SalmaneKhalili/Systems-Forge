PATH := $(HOME)/.local/go/bin:$(PATH)
GO   ?= go

all: switch svc
switch:
	$(GO) build -o switch .
svc:
	$(GO) build -o svc ./backend
clean:
	rm -f switch svc
fclean: clean
re: fclean all
.PHONY: all clean fclean re
