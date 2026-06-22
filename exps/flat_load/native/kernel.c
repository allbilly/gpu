#include <stdio.h>
#include <stdlib.h>
#include <CL/cl.h>
#include "kernels.h"
#include <time.h>

uint64_t getTimeInNSecs(){
    struct timespec time;
    clock_gettime(CLOCK_MONOTONIC, &time);
    uint64_t timeInSec = time.tv_sec * 1e9 + time.tv_nsec;
    return timeInSec;
}

int main(int argc, char **argv){

    cl_int err;
    uint64_t start, end, kernel_time;
    int loop_count;
    size_t grid_size, workgroup_size;

    cl_mem d_in;

    if (argc != 2){
        printf("./kernel loop_count num_wf\n");
        exit(-1);
    }
    loop_count = atoi(argv[1]);
    //num_wf = atoi(argv[2]);
    printf("loop_count %d, num_wf %d\n",loop_count);

    grid_size = 1;
    
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

    cl_program program = clCreateProgramWithSource(context, 1, &kernels_ocl, NULL, &err);
    if(err != CL_SUCCESS){
        printf("createt program failed:err_%d\n",err);
    }

    err = clBuildProgram(program, 1, &device, "-save-temps", NULL, NULL);
    if(err != CL_SUCCESS){
        size_t len;
        char buffer[2048];

        printf("kernel build failed:err_%d\n",err);
        clGetProgramBuildInfo(program, device, CL_PROGRAM_BUILD_LOG, sizeof(buffer), buffer, &len);
        printf("%s\n", buffer);
    }

    cl_kernel kernel = clCreateKernel(program, "microbench", &err);
    if(err != CL_SUCCESS){
        printf("create kernel failed:err_%d\n",err);
    }
    d_in = clCreateBuffer(context,  CL_MEM_READ_ONLY,  sizeof(int)*loop_count, NULL, &err);
    if(err != CL_SUCCESS){
        printf("create buffer d_in failed:err_%d\n",err);
    }

    err = clSetKernelArg(kernel,0,sizeof(cl_mem), &d_in);
    err |= clSetKernelArg(kernel,1,sizeof(unsigned int), &loop_count);
    if(err != CL_SUCCESS){
        printf("Setting kernel arguments failed:err_%d\n",err);
    }
    start = getTimeInNSecs();
    err = clEnqueueNDRangeKernel(queue, kernel, 1, NULL, &grid_size, NULL, 0, NULL, NULL);
    if(err != CL_SUCCESS){
        printf("kernel enqueue failed:err_%d\n",err);
    }
    clFinish(queue);

    end = getTimeInNSecs();
    kernel_time = end - start;
    printf("kernel time: %ld\n",kernel_time);

    return 0;

}

