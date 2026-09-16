NAME = test
SRCS = main.c
OBJS = $(SRCS:.c=.o)

%.o: %.c
	gcc $(CFLAGS) -c $^ -o $@

all: $(NAME)

$(NAME): $(OBJS)
	gcc $(CFLAGS) $^ -o $@

clean:
	rm -f $(OBJS)
fclean: clean
	rm -f $(NAME)
re: fclean all

.PHONY: re all fclean clean
