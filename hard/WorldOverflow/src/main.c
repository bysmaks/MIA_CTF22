#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>

#define FLAG_FILE "flag.txt"

#define WORLD_BANNER \
"o               .        ___---___                    .      \n" \
"       .              .--\\        --.     .     .         .  \n" \
"                    ./.;_.\\     __/~ \\.                      \n" \
"                   /;  / `-'  __\\    . \\                     \n" \
" .        .       / ,--'     / .   .;   \\        |           \n" \
"                 | .|       /       __   |      -O-       .  \n" \
"                |__/    __ |  . ;   \\ | . |      |           \n" \
"                |      /  \\_    . ;| \\___|                  \n" \
"   .    o       |      \\  .~\\___,--'     |           .      \n" \
"                 |     | . ; ~~~~\\_    __|                   \n" \
"    |             \\    \\   .  .  ; \\  /_/   .                \n" \
"   -O-        .    \\   /         . |  ~/                  .  \n" \
"    |    .          ~\\ \\   .      /  /~          o          '\n" \
"  .                   ~--___ ; ___--~                        \n" \
"                 .          ---         .                    \n"

int get_flag(char * flag, size_t size)
{
	int ret = 0;

	int fd = open(FLAG_FILE, O_RDONLY);
	if (-1 == fd)
	{
		fprintf(stderr, "open '%s' failed: %s\n", FLAG_FILE, strerror(errno));
		return -1;
	}

	ssize_t rc = read(fd, flag, size);
	if ((0 == rc) || (rc > size))
	{
		fprintf(stderr, "read failed: %s\n", strerror(errno));
		ret = -1;
	}

	close(fd);

	return ret;
}

int main()
{
	setvbuf(stdin, NULL, _IONBF, 0);
	setvbuf(stdout, NULL, _IONBF, 0);
	setvbuf(stderr, NULL, _IONBF, 0);

	char flag[32] = { 0 };
	char buf[64] = { 0 };

	if (get_flag(flag, sizeof(flag)))
	{
		fprintf(stderr, "get_flag failed\n");
		return -1;
	}

	printf("%s", WORLD_BANNER);
	printf("What is your main wish?\n> ");

	read(0, buf, 256);

	printf("\nYour wish will be fulfilled!\n");
	printf("%s", buf);

	return 0;
}
