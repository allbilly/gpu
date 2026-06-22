all: kernel

kernel: kernel.c kernels.h
	gcc kernel.c -o kernel -L/opt/rocm/opencl/lib -lOpenCL

kernels.h: kernels.cl
	./stringify_ocl.sh kernels.cl
clean:
	rm *.h kernel _temp*
