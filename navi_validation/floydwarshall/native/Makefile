all: kernel

kernel: floydwarshall.c floydwarshall.cl
	gcc floydwarshall.c -lOpenCL -o kernel
clean:
	rm kernel
