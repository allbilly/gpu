all: kernel

kernel: bicg.c bicg.cl
	gcc bicg.c -lOpenCL -o kernel
clean:
	rm kernel
