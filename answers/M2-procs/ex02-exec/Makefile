NAME = test
SRCS = main.c
OBJS = $(SRCS:.c=.o)

%.o:%.c
	gcc $(CFLAGS) -c $^ -o $@

all:$(NAME)

$(NAME):$(OBJS)
	gcc $(CFLAGS) $^ -o $@

fclean: clean
	rm -f $(NAME)

clean:
	rm -f $(OBJS)

re: fclean all

.PHONY: fclean all re clean