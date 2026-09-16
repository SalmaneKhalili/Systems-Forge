NAME = hello

all: $(NAME)

$(NAME): main.c
	$(CC) $(CFLAGS) -o $@ $<

fclean:
	rm -rf $(NAME)

re: fclean all

.PHONY: all fclean re