NAME = test 
SRCS = main.c ft_strcmp.c 
OBJS = $(SRCS:.c=.o)

all:$(NAME)

$(NAME): $(OBJS)
	gcc $(CFLAGS) $^ -o $@ 


%.o: %.c 
	gcc $(CFLAGS) -c $< -o $@

clean:
	rm -f $(OBJS)

fclean: clean
	rm -f $(NAME)

re: fclean all

.PHONY: re all fclean clean
