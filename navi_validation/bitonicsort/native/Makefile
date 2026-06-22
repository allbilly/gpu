all: kernel

kernel: bs.c bs.cl
	gcc bs.c -lOpenCL -o kernel
clean:
	rm kernel
