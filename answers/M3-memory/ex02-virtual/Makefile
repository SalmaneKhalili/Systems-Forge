NAME = test
SRCS = main.c
OBJS = $(SRCS:.c=.o)

%.o: %.c
	gcc $(CFLAGS) -c $< -o $@

all: $(OBJS)
	gcc $(CFLAGS) $(OBJS) -o $(NAME)

clean:
	rm $(OBJS)

fclean: clean
	rm $(NAME)

re: fclean all 

.PHONY: re fclean all clean
