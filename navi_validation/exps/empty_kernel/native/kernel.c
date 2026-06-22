#include <time.h>
#include <stdio.h>
#include <stdlib.h>
#include <CL/cl.h>

const char *kernelSource = "__kernel void empty()\n "
"{\n"
"}\n";

int main(int argc, char **argv){
    cl_int err;
    clock_t start, end;
    double kernel_time;
    int num_wg, num_wf_per_wg;
    size_t grid_size, workgroup_size;
    if (argc != 3){
        printf("./main num_wg num_wf_per_wg\n");
        exit(-1);
    }
    num_wg = atoi(argv[1]);
    num_wf_per_wg = atoi(argv[2]);
    printf("num_wg %d, num_wf_per_wg %d\n",num_wg,num_wf_per_wg);

    grid_size = 32*num_wf_per_wg*num_wg;
    workgroup_size = 32*num_wf_per_wg;
    
    cl_platform_id platform;
    clGetPlatformIDs(1,&platform, NULL);

    cl_device_id device;
    clGetDeviceIDs(platform, CL_DEVICE_TYPE_GPU, 1, &device, NULL);

    cl_context context = clCreateContext(NULL, 1, &device, NULL, NULL, &err);
    if(err != CL_SUCCESS){
        printf("create context failed:err_%d\n",err);
    }

    cl_command_queue queue = clCreateCommandQueue(context, device, 0, &err);
    if(err != CL_SUCCESS){
        printf("createt queue failed:err_%d\n",err);
    }

    cl_program program = clCreateProgramWithSource(context, 1, &kernelSource, NULL, &err);
    if(err != CL_SUCCESS){
        printf("createt program failed:err_%d\n",err);
    }

    err = clBuildProgram(program, 1, &device, NULL, NULL, NULL);
    if(err != CL_SUCCESS){
        printf("kernel build failed:err_%d\n",err);
    }

    cl_kernel kernel = clCreateKernel(program, "empty", &err);
    if(err != CL_SUCCESS){
        printf("create kernel failed:err_%d\n",err);
    }

    start = clock();
    err = clEnqueueNDRangeKernel(queue, kernel, 1, NULL, &grid_size, &workgroup_size, 0, NULL, NULL);
    if(err != CL_SUCCESS){
        printf("kernel enqueue failed:err_%d\n",err);
    }
    clFinish(queue);

    end = clock();
    kernel_time = ((double)(end - start))/CLOCKS_PER_SEC;
    printf("kernel time: %g\n",kernel_time);

    return 0;

}
