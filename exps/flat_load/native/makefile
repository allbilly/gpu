all: kernel

kernel: kernel.c kernels.h
	gcc kernel.c -o kernel -lOpenCL

kernels.h: kernels.cl
	./stringify_ocl.sh kernels.cl
clean:
	rm *.h kernel
