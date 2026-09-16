#include <pthread.h>
#include <stdio.h>

static pthread_rwlock_t rw = PTHREAD_RWLOCK_INITIALIZER;
static pthread_mutex_t ctl = PTHREAD_MUTEX_INITIALIZER;
static long a, b; /* invariant: b == 2 * a */
static int version;
static long reads, consistent;

#define NREADERS 4
#define NWRITERS 2
#define RREADS 50
#define RWRITES 10

static void *reader(void *arg)
{
	int r;
	(void)arg;
	for (r = 0; r < RREADS; r++) {
		int ok;
		pthread_rwlock_rdlock(&rw);
		ok = (b == 2 * a);
		pthread_rwlock_unlock(&rw);
		pthread_mutex_lock(&ctl);
		reads++;
		if (ok)
			consistent++;
		pthread_mutex_unlock(&ctl);
	}
	return NULL;
}

static void *writer(void *arg)
{
	int w;
	(void)arg;
	for (w = 0; w < RWRITES; w++) {
		pthread_rwlock_wrlock(&rw);
		a++;
		b = 2 * a;
		version++;
		pthread_rwlock_unlock(&rw);
	}
	return NULL;
}

int main(void)
{
	pthread_t rd[NREADERS], wr[NWRITERS];
	int i;

	for (i = 0; i < NREADERS; i++)
		if (pthread_create(&rd[i], NULL, reader, NULL) != 0) {
			fprintf(stderr, "create reader failed\n");
			return 1;
		}
	for (i = 0; i < NWRITERS; i++)
		if (pthread_create(&wr[i], NULL, writer, NULL) != 0) {
			fprintf(stderr, "create writer failed\n");
			return 1;
		}
	for (i = 0; i < NREADERS; i++)
		pthread_join(rd[i], NULL);
	for (i = 0; i < NWRITERS; i++)
		pthread_join(wr[i], NULL);

	if (reads != NREADERS * RREADS) {
		printf("read rounds %ld (want %d)\n", reads, NREADERS * RREADS);
		return 1;
	}
	printf("read rounds 200\n");
	printf("write rounds %d\n", NWRITERS * RWRITES);
	if (consistent != NREADERS * RREADS) {
		printf("consistent %ld (want %d)\n", consistent, NREADERS * RREADS);
		return 1;
	}
	printf("consistent 200\n");
	if (version != NWRITERS * RWRITES) {
		printf("version %d (want %d)\n", version, NWRITERS * RWRITES);
		return 1;
	}
	printf("version 20\n");
	printf("ALL PASS\n");
	return 0;
}