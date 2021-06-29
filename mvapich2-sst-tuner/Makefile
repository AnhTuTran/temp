CC=mpicc
CFLAGS=-fPIC -g -Wall

all:
	$(CC) $(CFLAGS) -c mv2_coll_interceptor.c -o mv2_coll_interceptor.o
	ar rcs libmv2_coll_interceptor.a mv2_coll_interceptor.o
	$(CC) -shared -o libmv2_coll_interceptor.so mv2_coll_interceptor.o

clean:
	rm -f *.o libmv2_coll_interceptor.*

