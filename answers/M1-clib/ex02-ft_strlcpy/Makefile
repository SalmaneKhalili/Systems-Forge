NAME = test 

all: $(NAME)

$(NAME): main.c ft_strlcpy.c
	gcc $(CFLAGS) main.c ft_strlcpy.c -o $@

fclean:
	rm -rf $(NAME)

re: fclean all

.PHONY: fclean all re
