all: kernel

kernel: relu.c relu.cl
	gcc relu.c -lOpenCL -o kernel
clean:
	rm kernel
