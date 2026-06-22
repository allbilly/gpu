all: kernel

kernel: matrixmultiplication.c matrixmultiplication.cl
	gcc matrixmultiplication.c -lOpenCL -o kernel
clean:
	rm kernel