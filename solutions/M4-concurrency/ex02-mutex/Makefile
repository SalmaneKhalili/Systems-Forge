CC     ?= cc
NAME    = test
SRCS    = main.c
OBJS    = $(SRCS:.c=.o)

all: $(NAME)

$(NAME): $(OBJS)
	$(CC) $(CFLAGS) -pthread -o $@ $(OBJS)

%.o: %.c
	$(CC) $(CFLAGS) -pthread -c $< -o $@

clean:
	rm -f $(OBJS)

fclean: clean
	rm -f $(NAME)

re: fclean all

.PHONY: all clean fclean re