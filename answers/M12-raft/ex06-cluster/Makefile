PATH := $(HOME)/.local/go/bin:$(PATH)
GO   ?= go
NAME  = cluster

all: $(NAME)
$(NAME):
	$(GO) build -o $(NAME) .
clean:
	rm -f $(NAME)
fclean: clean
re: fclean all
.PHONY: all clean fclean re