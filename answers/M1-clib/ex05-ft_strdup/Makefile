NAME = test 
SRCS = main.c ft_strdup.c 
OBJS = $(SRCS:.c=.o)

%.o: %.c 
	gcc $(CFLAGS) -c $< -o $@

all: $(NAME)

$(NAME):$(OBJS)
	gcc $(CFLAGS) $^ -o $@

clean:
	rm -f *.o 
fclean: clean
	rm -f $(NAME)

re: fclean all 

.PHONY: re fclean clean all

