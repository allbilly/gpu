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


//define the error threshold for the results "not matching"
#define PERCENT_DIFF_ERROR_THRESHOLD 0.05

#define MAX_SOURCE_SIZE (0x100000)


#ifndef M_PI
#define M_PI 3.14159
#endif

#define SMALL_FLOAT_VAL 0.00000001
char str_temp[1024];


cl_platform_id platform_id;
cl_device_id device_id;   
cl_uint num_devices;
cl_uint num_platforms;
cl_int errcode;
cl_context clGPUContext;
cl_kernel clKernel1;
cl_kernel clKernel2;
cl_command_queue clCommandQue;
cl_program clProgram;
cl_mem a_mem_obj;
cl_mem x_mem_obj;
cl_mem y_mem_obj;
cl_mem tmp_mem_obj;
FILE *fp;
char *source_str;
size_t source_size;

uint64_t getTimeInNSecs(){
    struct timespec time;
    clock_gettime(CLOCK_MONOTONIC, &time);
    uint64_t timeInNSec = time.tv_sec * 1e9 + time.tv_nsec;
    return timeInNSec;
}

float absVal(float a)
{
	if(a < 0)
	{
		return (a * -1);
	}
   	else
	{ 
		return a;
	}
}



float percentDiff(double val1, double val2)
{
	if ((absVal(val1) < 0.01) && (absVal(val2) < 0.01))
	{
		return 0.0f;
	}

	else
	{
    		return 100.0f * (absVal(absVal(val1 - val2) / absVal(val1 + SMALL_FLOAT_VAL)));
	}
}

void verify(float* A, float* x, float* y, float* tmp, float *y_outputFromGpu, int NX, int NY)
{
	for (int i= 0; i < NY; i++)
	{
    		y[i] = 0;
	}
  
	for (int i = 0; i < NX; i++)
 	{
		tmp[i] = 0;

		for (int j = 0; j < NY; j++)
		{
			tmp[i] = tmp[i] + A[i*NY + j] * x[j];
		}
			
		for (int j = 0; j < NY; j++)
		{
			y[j] = y[j] + A[i*NY + j] * tmp[i];
		}
    	}
	int i, fail;
	fail = 0;

	for (i=0; i<NY; i++)
	{
		if (percentDiff(y[i], y_outputFromGpu[i]) > PERCENT_DIFF_ERROR_THRESHOLD)
		{
			fail++;
		}	
	}
	
	// print results
	printf("Non-Matching CPU-GPU Outputs Beyond Error Threshold of %4.2f Percent: %d\n", PERCENT_DIFF_ERROR_THRESHOLD, fail);

}


void read_cl_file()
{
	// Load the kernel source code into the array source_str
	fp = fopen("atax.cl", "r");
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


void cl_mem_init(float* A, float* x, float* y, float* tmp)
{
	
}


void cl_load_prog()
{
	// Create a program from the kernel source
	clProgram = clCreateProgramWithSource(clGPUContext, 1, (const char **)&source_str, (const size_t *)&source_size, &errcode);

	if(errcode != CL_SUCCESS) printf("Error in creating program\n");

	// Build the program
	errcode = clBuildProgram(clProgram, 1, &device_id, NULL, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in building program\n");
		
	// Create the 1st OpenCL kernel
	clKernel1 = clCreateKernel(clProgram, "atax_kernel1", &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating kernel\n");

	// Create the 2nd OpenCL kernel
	clKernel2 = clCreateKernel(clProgram, "atax_kernel2", &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating kernel\n");
	clFinish(clCommandQue);
}


void cl_launch_kernel()
{
	
}


void cl_clean_up()
{
	// Clean up
	errcode = clFlush(clCommandQue);
	errcode = clFinish(clCommandQue);
	errcode = clReleaseKernel(clKernel1);
	errcode = clReleaseKernel(clKernel2);
	errcode = clReleaseProgram(clProgram);
	errcode = clReleaseMemObject(a_mem_obj);
	errcode = clReleaseMemObject(x_mem_obj);
	errcode = clReleaseMemObject(y_mem_obj);
	errcode = clReleaseMemObject(tmp_mem_obj);
	errcode = clReleaseCommandQueue(clCommandQue);
	errcode = clReleaseContext(clGPUContext);
	if(errcode != CL_SUCCESS) printf("Error in cleanup\n");
}


void atax_cpu()
{
	int i,j;
	
	
}


int main(int argc, char** argv) 
{
	if(argc != 2){
		printf("./kernel NX=NY\n");
		exit(-1);
	}
	int NX,NY;

	NX = atoi(argv[1]);
	NY = atoi(argv[1]);

	float* A;
	float* x;
	float* y;
	float* y_outputFromGpu;
	float* tmp;

	A = (float*)malloc(NX*NY*sizeof(float));
	x = (float*)malloc(NY*sizeof(float));
	y = (float*)malloc(NY*sizeof(float));
	y_outputFromGpu = (float*)malloc(NY*sizeof(float));
	tmp = (float*)malloc(NX*sizeof(float));
	
	for (int i = 0; i < NX; i++)
	{
		x[i] = i * M_PI;
		for (int j = 0; j < NY; j++)
		{
			A[i*NY + j] = ((float) i*(j)) / NX;
		}
	}
	read_cl_file();
	cl_initialization();
	
	a_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(float) * NX * NY, NULL, &errcode);
	x_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(float) * NY, NULL, &errcode);
	y_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(float) * NY, NULL, &errcode);
	tmp_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(float) * NX, NULL, &errcode);		
	if(errcode != CL_SUCCESS) printf("Error in creating buffers\n");
	
	errcode = clEnqueueWriteBuffer(clCommandQue, a_mem_obj, CL_TRUE, 0, sizeof(float) * NX * NY, A, 0, NULL, NULL);
	errcode = clEnqueueWriteBuffer(clCommandQue, x_mem_obj, CL_TRUE, 0, sizeof(float) * NY, x, 0, NULL, NULL);
	errcode = clEnqueueWriteBuffer(clCommandQue, y_mem_obj, CL_TRUE, 0, sizeof(float) * NY, y, 0, NULL, NULL);
	errcode = clEnqueueWriteBuffer(clCommandQue, tmp_mem_obj, CL_TRUE, 0, sizeof(float) * NX, tmp, 0, NULL, NULL);
	if(errcode != CL_SUCCESS)printf("Error in writing buffers\n");
	cl_load_prog();

	uint64_t t_start, t_end, duration;

	size_t local, global;
	local=256;
	global = (size_t)(((float)NX) / ((float)local)+1) * local;
	
	// Set the arguments of the kernel
	errcode =  clSetKernelArg(clKernel1, 0, sizeof(cl_mem), (void *)&a_mem_obj);
	errcode |= clSetKernelArg(clKernel1, 1, sizeof(cl_mem), (void *)&x_mem_obj);
	errcode |= clSetKernelArg(clKernel1, 2, sizeof(cl_mem), (void *)&tmp_mem_obj);
	errcode |= clSetKernelArg(clKernel1, 3, sizeof(int), (void *)&NX);
	errcode |= clSetKernelArg(clKernel1, 4, sizeof(int), (void *)&NY);
	if(errcode != CL_SUCCESS) printf("Error in seting arguments\n");

	t_start = getTimeInNSecs();
	// Execute the OpenCL kernel
	errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel1, 1, NULL, &global, &local, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in launching kernel\n");
	clEnqueueBarrier(clCommandQue);
	t_end = getTimeInNSecs();	
	duration = t_end - t_start;
	
	global = (size_t)(((float)NY) / ((float)local)+1) * local;

	// Set the arguments of the kernel
	errcode =  clSetKernelArg(clKernel2, 0, sizeof(cl_mem), (void *)&a_mem_obj);
	errcode |= clSetKernelArg(clKernel2, 1, sizeof(cl_mem), (void *)&y_mem_obj);
	errcode |= clSetKernelArg(clKernel2, 2, sizeof(cl_mem), (void *)&tmp_mem_obj);
	errcode |= clSetKernelArg(clKernel2, 3, sizeof(int), (void *)&NX);
	errcode |= clSetKernelArg(clKernel2, 4, sizeof(int), (void *)&NY);
	if(errcode != CL_SUCCESS) printf("Error in seting arguments\n");

	t_start = getTimeInNSecs();	
	errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel2, 1, NULL, &global, &local, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in launching kernel\n");
	clFinish(clCommandQue);
	t_end = getTimeInNSecs();
	duration += t_end - t_start;
	fprintf(stdout, "GPU Runtime: %f\n", duration/1e9);

	errcode = clEnqueueReadBuffer(clCommandQue, y_mem_obj, CL_TRUE, 0, NY*sizeof(float), y_outputFromGpu, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in reading GPU mem\n");

	verify(A, x, y, tmp, y_outputFromGpu, NX, NY);
	cl_clean_up();
	
	free(A);
	free(x);
	free(y);
	free(y_outputFromGpu);
	free(tmp);
	
    return 0;
}

