CC = cc
CFLAGS = -Wall -Wextra -Werror
NAME = test
SRCS = main.c
RM = rm -f
OBJS = $(SRCS:.c=.o)

%.o:%.c
	$(CC) $(CFLAGS) -c $< -o $@

all: $(NAME)

$(NAME):$(OBJS)
	$(CC) $(CFLAGS) $^ -o $@

clean:
	$(RM) $(OBJS)

fclean: clean
	$(RM) $(NAME)

re: fclean all

.PHONY: re fclean clean all