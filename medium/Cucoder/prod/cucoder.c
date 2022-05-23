#include <stdio.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <string.h>
#include <errno.h>
#include <unistd.h>

#define FLAG_FILE "flag.txt"
#define OUT_FILE  "out.txt"

#define BUF_SIZE 32

static int get_flag(char * flag, size_t size)
{
	int ret = 0;

	int fd = open(FLAG_FILE, O_RDONLY);
	if (-1 == fd)
	{
		fprintf(stderr, "open '%s' failed: %s", FLAG_FILE, strerror(errno));
		return -1;
	}

	ssize_t rc = read(fd, flag, size);
	if ((0 == rc) || (rc > size))
	{
		fprintf(stderr, "read failed: %s", strerror(errno));
		ret = -1;
	}

	close(fd);

	return ret;
}

static void encrypt_flag(char * flag, size_t len)
{
	for (int i = 0; i < len; i++)
	{
		if ((flag[i] | 0xFE) == 0xFE)
		{
			flag[i] = flag[i] - 1;
		}
		else
		{
			flag[i] = flag[i] + 1;
		}
	}
}

static int save_flag(char * flag, size_t size)
{
	int ret = 0;

	int fd = open(OUT_FILE, O_WRONLY | O_CREAT, 0666);
	if (-1 == fd)
	{
		fprintf(stderr, "open '%s' failed: %s", OUT_FILE, strerror(errno));
		return -1;
	}

	ssize_t rc = write(fd, flag, size);
	if ((0 == rc) || (rc > size))
	{
		fprintf(stderr, "write failed: %s", strerror(errno));
		ret = -1;
	}

	close(fd);

	return ret;
}

int main(int argc, char ** argv)
{
	char flag[BUF_SIZE] = { 0 };

	if (get_flag(flag, sizeof(flag)))
	{
		fprintf(stderr, "get_flag failed\n");
		return -1;
	}

	//printf("%s\n", flag);

	encrypt_flag(flag, strlen(flag));

	//printf("%s\n", flag);

	if (save_flag(flag, sizeof(flag)))
	{
		fprintf(stderr, "save_flag failed\n");
		return -1;
	}

	return 0;
}
