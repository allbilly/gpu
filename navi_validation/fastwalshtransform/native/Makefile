all: kernel

kernel: fastwalshtransform.c fastwalshtransform.cl
	gcc fastwalshtransform.c -lOpenCL -o kernel
clean:
	rm kernel
