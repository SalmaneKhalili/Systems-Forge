CC ?= cc
CFLAGS += -Wall -Wextra -Werror

all: micro
micro: micro.c
	$(CC) $(CFLAGS) -o micro micro.c
clean:
	rm -f micro
fclean: clean
re: fclean all
.PHONY: all clean fclean re