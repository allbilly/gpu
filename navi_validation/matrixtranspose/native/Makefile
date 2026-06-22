all: kernel

kernel: matrixtranspose.c matrixtranspose.cl
	gcc matrixtranspose.c -lOpenCL -o kernel
clean:
	rm kernel
