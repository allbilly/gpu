all: kernel

kernel: spmv.c spmv.cl
	gcc spmv.c -g -lOpenCL -o kernel
clean:
	rm kernel