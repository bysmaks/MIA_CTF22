#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <time.h>

#define UPPER 99999999
#define LOWER 10000000

#define flag "CTF{1_g3t_s3cr3t_c0d3_3z}"

int main(int argc, char ** argv)
{
	srand(time(NULL));

	unsigned int secret = (rand() % (UPPER - LOWER + 1)) + LOWER;

	printf("Hi. Tell me the secret code\n> ");

	unsigned int input;

	if (scanf("%u", &input) != 1)
	{
		fprintf(stderr, "Input failed\n");
		return -1;
	}

	if (input == secret)
	{
		printf("Correct code! Take it %s\n'", flag);
	}
	else
	{
		printf("Incorrect code. Try harder\n");
	}

	return 0;
}
