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
#include <stdbool.h>

#include <CL/cl.h>


#define MAX_SOURCE_SIZE (0x100000)

/* Thread block dimensions */
#define DIM_LOCAL_WORK_GROUP_X 64

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
cl_mem output_mem_obj;
FILE *fp;
char *source_str;
size_t source_size;

uint64_t getTimeInNSecs(){
    struct timespec time;
    clock_gettime(CLOCK_MONOTONIC, &time);
    uint64_t timeInNSec = time.tv_sec * 1e9 + time.tv_nsec;
    return timeInNSec;
}

void verify(int *input, int *output, int width, int height)
{
	bool mismatch;
	int *output_cpu=malloc(width*height*sizeof(int));
	for(int i=0; i < height; i++){
        	for(int j=0; j < width; j++){
        		output_cpu[j*width+i] = input[i*width+j];
        	}
	}
	for(int i = 0; i < height; i++) {
		for(int j=0; j<width; j++){
			int actual = output[i*width+j];
			int expected = output_cpu[i*width+j];
			if (actual != expected){
				mismatch = true;
				printf("mismatch at (%d, %d), expected %d, but get %d\n", i,j, expected, actual);
			}
		}
	}
	free(output_cpu);
	if(!mismatch){
		printf("Passed!\n");
	}
}

void read_cl_file()
{
	// Load the kernel source code into the array source_str
	fp = fopen("matrixtranspose.cl", "r");
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
	clKernel = clCreateKernel(clProgram, "matrixTranspose", &errcode);
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
	errcode = clReleaseMemObject(output_mem_obj);
	errcode = clReleaseCommandQueue(clCommandQue);
	errcode = clReleaseContext(clGPUContext);
	if(errcode != CL_SUCCESS) printf("Error in cleanup\n");
}

int main(int argc, char** argv) 
{
	if(argc != 2){
		printf("./kernel size of square matrix\n");
		exit(-1);
	}
	int width, height,numData;
	int elemsPerThread1Dim = 4;
	int blockSize = 4;
	width = atoi(argv[1]);
	height = atoi(argv[1]);
	numData = width*height;
	int* input;
	int* output;

	input = (int*)malloc(numData*sizeof(int));
	output = (int*)malloc(numData*sizeof(int));
	
	for (int i = 0; i < numData; i++)
	{
		input[i] = i;
	}
	read_cl_file();
	cl_initialization();

	input_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_ONLY, sizeof(int) * numData, NULL, &errcode);
	output_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(int) * numData, NULL, &errcode);		
	if(errcode != CL_SUCCESS) printf("Error in creating buffers\n");
	
	errcode = clEnqueueWriteBuffer(clCommandQue, input_mem_obj, CL_TRUE, 0, sizeof(int) * numData, input, 0, NULL, NULL);
	if(errcode != CL_SUCCESS)printf("Error in writing buffers\n");

	cl_load_prog();
	int neededLocalMemory = blockSize * blockSize * elemsPerThread1Dim * elemsPerThread1Dim * sizeof(int);
	int wiWidth = width/elemsPerThread1Dim;
	int wiHeight = height/elemsPerThread1Dim;
	int numWGWidth = wiWidth / blockSize;

	uint64_t t_start, t_end;

	size_t local[2] = {blockSize, blockSize};
	size_t global[2] = {wiWidth, wiHeight};;

	// Set the arguments of the kernel
	errcode =  clSetKernelArg(clKernel, 0, sizeof(cl_mem), (void *)&output_mem_obj);
	errcode |= clSetKernelArg(clKernel, 1, sizeof(cl_mem), (void *)&input_mem_obj);
	errcode |= clSetKernelArg(clKernel, 2, (size_t)neededLocalMemory, NULL);
	errcode |=  clSetKernelArg(clKernel, 3, sizeof(int), (void *)&wiWidth);
	errcode |= clSetKernelArg(clKernel, 4, sizeof(int), (void *)&wiHeight);
	errcode |= clSetKernelArg(clKernel, 5, sizeof(int), (void *)&numWGWidth);
	uint32_t group_x_offset = 0;
	errcode |= clSetKernelArg(clKernel,6, sizeof(uint32_t),(void *)&group_x_offset);
	uint32_t group_y_offset = 0;
	errcode |= clSetKernelArg(clKernel,7, sizeof(uint32_t),(void *)&group_y_offset);
	if(errcode != CL_SUCCESS) printf("Error in seting arguments\n");

	t_start = getTimeInNSecs();
	errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel, 2, NULL, global, local, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in launching kernel\n");
	
	clFinish(clCommandQue);
	t_end = getTimeInNSecs();
	
	fprintf(stdout, "GPU Runtime: %f\n", (t_end - t_start)/1e9);

	errcode = clEnqueueReadBuffer(clCommandQue, output_mem_obj, CL_TRUE, 0, numData*sizeof(int), output, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in reading GPU mem\n");

	verify(input, output, width, height);
	cl_clean_up();
	
	free(input);
	free(output);
	
    return 0;
}

