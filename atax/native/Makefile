all: kernel

kernel: atax.c atax.cl
	gcc atax.c -lOpenCL -o kernel
clean:
	rm kernel
