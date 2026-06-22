#include <stdio.h>
#include <stdlib.h>
#include <CL/cl.h>
#include "kernels.h"
#include <time.h>
#define COUNT 10000000


uint64_t getTimeInNSecs(){
    struct timespec time;
    clock_gettime(CLOCK_MONOTONIC, &time);
    uint64_t timeInSec = time.tv_sec * 1e9 + time.tv_nsec;
    return timeInSec;
}

int main(int argc, char **argv){

    cl_int err;
    uint64_t start, end, kernel_time;
    int len, size;
    unsigned int count;
    size_t grid_size;
    cl_mem d_memory;
    cl_mem d_offset;

    if (argc != 2){
        printf("./kernel memory_size\n");
        exit(-1);
    }


    size = atoi(argv[1]);
    //count = atoi(argv[2]);
    count = COUNT;
    printf("memory_size %d, count %u\n", size, count);

    grid_size = 1;
    //workgroup_size = 32;

    len = size / sizeof(void*);

    
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
    //d_memory = clCreateBuffer(context,  CL_MEM_READ_WRITE,  sizeof(void*)*len, NULL, &err);
    d_memory = clCreateBuffer(context,  CL_MEM_READ_WRITE,  sizeof(uint64_t)*len, NULL, &err);
    if(err != CL_SUCCESS){
        printf("create buffer d_memory failed:err_%d\n",err);
    }

    d_offset = clCreateBuffer(context,  CL_MEM_READ_WRITE,  sizeof(uint64_t)*64, NULL, &err);
    if(err != CL_SUCCESS){
        printf("create buffer d_memory failed:err_%d\n",err);
    }


    uint64_t *memory = malloc(len * sizeof(uint64_t));

    for(int i=0; i<len; i++){
        if(i == len - 1){
            memory[i] = 0;
        }else{
            //memory[i] = (void*)(&(d_memory[i]) + 1);
            memory[i] = (i+1)*8;
        }
        //printf("%lx ",(uint64_t)memory[i]);
    }
    //printf("\n");
    err = clEnqueueWriteBuffer(queue, d_memory, CL_TRUE, 0, sizeof(void*)*len, memory, 0, NULL, NULL);
    if(err != CL_SUCCESS){
        printf("copy memory to d_memory failed:err_%d\n",err);
    }
    free(memory);

    ulong zero = 0;
    err = clEnqueueFillBuffer(queue, d_offset, &zero, sizeof(uint64_t), 0, sizeof(uint64_t)*64, 0, NULL, NULL);
    if(err != CL_SUCCESS){
        printf("fill buffer failed:err_%d\n",err);
    }

    err = clSetKernelArg(kernel,0,sizeof(cl_mem), &d_memory);
    err |= clSetKernelArg(kernel,1,sizeof(cl_mem), &d_offset);
    err |= clSetKernelArg(kernel,2,sizeof(unsigned int), &count);
    if(err != CL_SUCCESS){
        printf("Setting kernel arguments failed:err_%d\n",err);
    }
    start = getTimeInNSecs();
    err = clEnqueueNDRangeKernel(queue, kernel, 1, NULL, &grid_size, NULL, 0, NULL, NULL);
    if(err != CL_SUCCESS){
        printf("kernel enqueue failed:err_%d\n",err);
    }
    clFinish(queue);
    // void** back = malloc(sizeof(void*)*len);
    // err = clEnqueueReadBuffer(queue, d_memory, CL_TRUE, 0, sizeof(void*)*len, back, 0, NULL, NULL);
    // if(err != CL_SUCCESS){
    //     printf("copy d_memory to back failed:err_%d\n",err);
    // }
    // for(int i = 0; i < len; i++){
    //     printf("%lx ",(uint64_t)back[i]);
    // }
    // printf("\n");
    end = getTimeInNSecs();
    kernel_time = end - start;
    printf("kernel time: %ld\n",kernel_time);

    return 0;
}
    

