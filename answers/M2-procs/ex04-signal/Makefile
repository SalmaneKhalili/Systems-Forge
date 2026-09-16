NAME = test
SRCS = main.c
OBJS = $(SRCS:.c=.o)

CC = cc
RM = rm -f
# CFLAGS = -Wall -Wextra -Werror  this is commented out for now since the grader establishes its own flags

all: $(NAME)

$(NAME): $(OBJS)
	$(CC) $(CFLAGS) $^ -o $@

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	$(RM) $(OBJS)

fclean: clean
	$(RM) $(NAME)

re: fclean all

.PHONY: fclean re clean all