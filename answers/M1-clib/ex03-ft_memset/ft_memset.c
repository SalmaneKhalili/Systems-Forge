#include <stddef.h>

void	*ft_memset(void *s, int c, size_t n)
{
	char	*cpy;
	size_t		i;

	if (n == 0)
		return (s);
	cpy = (char *)s;
	i = 0;
	while (i < n)
	{
		*cpy = (unsigned char)c;
		cpy++;
		i++;
	}
	return (s);
}
