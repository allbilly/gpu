/**
 * atax.c: This file is part of the PolyBench/GPU 1.0 test suite.
 *
 *
 * Contact: Scott Grauer-Gray <sgrauerg@gmail.com>
 * Louis-Noel Pouchet <pouchet@cse.ohio-state.edu>
 * Web address: http://www.cse.ohio-state.edu/~pouchet/software/polybench/GPU
 */

#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <sys/time.h>
#include <math.h>

#include <CL/cl.h>

#define MAX_SOURCE_SIZE (0x100000)

char str_temp[1024];

cl_platform_id platform_id;
cl_device_id device_id;   
cl_uint num_devices;
cl_uint num_platforms;
cl_int errcode;
cl_context clGPUContext;
cl_kernel clKernel;
cl_command_queue clCommandQue;
cl_program clProgram;
cl_mem input_mem_obj;
FILE *fp;
char *source_str;
size_t source_size;

uint64_t getTimeInNSecs(){
    struct timespec time;
    clock_gettime(CLOCK_MONOTONIC, &time);
    uint64_t timeInNSec = time.tv_sec * 1e9 + time.tv_nsec;
    return timeInNSec;
}

void verify(int *output, int LEN)
{
	for(int i = 0; i < LEN -1; i++) {
		if (output[i] > output[i+1]){
			printf("Error: array[%d] > array[%d]: %d %d\n", i, i+1,
					output[i], output[i+1]);
			return;
		}
	}
	printf("Passed!\n");
}


void read_cl_file()
{
	// Load the kernel source code into the array source_str
	fp = fopen("bs.cl", "r");
	if (!fp) {
		fprintf(stderr, "Failed to load kernel.\n");
		exit(1);
	}
	source_str = (char*)malloc(MAX_SOURCE_SIZE);
	source_size = fread( source_str, 1, MAX_SOURCE_SIZE, fp);
	fclose( fp );
}

void cl_initialization()
{	
	// Get platform and device information
	errcode = clGetPlatformIDs(1, &platform_id, &num_platforms);
	if(errcode == CL_SUCCESS) printf("number of platforms is %d\n",num_platforms);
	else printf("Error getting platform IDs\n");

	errcode = clGetPlatformInfo(platform_id,CL_PLATFORM_NAME, sizeof(str_temp), str_temp,NULL);
	if(errcode == CL_SUCCESS) printf("platform name is %s\n",str_temp);
	else printf("Error getting platform name\n");

	errcode = clGetPlatformInfo(platform_id, CL_PLATFORM_VERSION, sizeof(str_temp), str_temp,NULL);
	if(errcode == CL_SUCCESS) printf("platform version is %s\n",str_temp);
	else printf("Error getting platform version\n");

	errcode = clGetDeviceIDs( platform_id, CL_DEVICE_TYPE_GPU, 1, &device_id, &num_devices);
	if(errcode == CL_SUCCESS) printf("number of devices is %d\n", num_devices);
	else printf("Error getting device IDs\n");

	errcode = clGetDeviceInfo(device_id,CL_DEVICE_NAME, sizeof(str_temp), str_temp,NULL);
	if(errcode == CL_SUCCESS) printf("device name is %s\n",str_temp);
	else printf("Error getting device name\n");
	
	// Create an OpenCL context
	clGPUContext = clCreateContext( NULL, 1, &device_id, NULL, NULL, &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating context\n");
 
	//Create a command-queue
	clCommandQue = clCreateCommandQueue(clGPUContext, device_id, 0, &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating command queue\n");
}


void cl_load_prog()
{
	// Create a program from the kernel source
	clProgram = clCreateProgramWithSource(clGPUContext, 1, (const char **)&source_str, (const size_t *)&source_size, &errcode);

	if(errcode != CL_SUCCESS) printf("Error in creating program\n");

	// Build the program
	errcode = clBuildProgram(clProgram, 1, &device_id, NULL, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in building program\n");
		
	// Create the OpenCL kernel
	clKernel = clCreateKernel(clProgram, "BitonicSort", &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating kernel\n");

}

void cl_clean_up()
{
	// Clean up
	errcode = clFlush(clCommandQue);
	errcode = clFinish(clCommandQue);
	errcode = clReleaseKernel(clKernel);
	errcode = clReleaseProgram(clProgram);
	errcode = clReleaseMemObject(input_mem_obj);
	errcode = clReleaseCommandQueue(clCommandQue);
	errcode = clReleaseContext(clGPUContext);
	if(errcode != CL_SUCCESS) printf("Error in cleanup\n");
}

int main(int argc, char **argv) 
{
	if(argc != 2){
		printf("./kernel length\n");
		exit(-1);
	}
	int LEN;
	LEN = atoi(argv[1]);
	
	int *input = (int*)malloc(LEN*sizeof(int));
	int *output = (int*)malloc(LEN*sizeof(int));

	// perPassIn = (int*)malloc(LEN*sizeof(int));
	// perPassOut = (int*)malloc(LEN*sizeof(int));
	
	for (int i = 0; i < LEN; i++)
	{
		input[i] = rand();
	}
	read_cl_file();
	cl_initialization();
	
	input_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(int) * LEN, NULL, &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating buffers\n");
	
	errcode = clEnqueueWriteBuffer(clCommandQue, input_mem_obj, CL_TRUE, 0, sizeof(int) * LEN, input, 0, NULL, NULL);
	if(errcode != CL_SUCCESS)printf("Error in writing buffers\n");

	int num_stages = 0;
	for(int temp = LEN; temp > 1; temp >>= 1){
		num_stages++;
	}
	int direction = 1;
	
	cl_load_prog();

	size_t local, global;
	local = 64;
	global = LEN/2;
	uint64_t t_start,t_end;
	uint64_t duration = 0;
	
	for(int stage = 0; stage < num_stages; stage++){
		for(int passOfStage = 0; passOfStage < stage+1; passOfStage++){
			// Set the arguments of the kernel
			errcode =  clSetKernelArg(clKernel, 0, sizeof(cl_mem), (void *)&input_mem_obj);
			errcode |= clSetKernelArg(clKernel, 1, sizeof(int), (void *)&stage);
			errcode |= clSetKernelArg(clKernel, 2, sizeof(int), (void *)&passOfStage);
			errcode |= clSetKernelArg(clKernel, 3, sizeof(int), (void *)&direction);
			if(errcode != CL_SUCCESS) printf("Error in seting arguments\n");

			t_start = getTimeInNSecs();
			errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel, 1, NULL, &global, &local, 0, NULL, NULL);
			if(errcode != CL_SUCCESS) printf("Error in launching kernel\n");
			clFinish(clCommandQue);	
			t_end = getTimeInNSecs();
			duration += t_end - t_start;
		}
	}
	
	fprintf(stdout, "GPU Runtime: %f\n", duration/1e9);
	
	errcode = clEnqueueReadBuffer(clCommandQue, input_mem_obj, CL_TRUE, 0, LEN*sizeof(int), output, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in reading GPU mem\n");

	verify(output,LEN);
	cl_clean_up();
	
	free(input);
	free(output);
	
	
    return 0;
}

