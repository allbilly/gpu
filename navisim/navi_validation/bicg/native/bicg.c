/**
 * bicg.c: This file is part of the PolyBench/GPU 1.0 test suite.
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

#define SMALL_FLOAT_VAL 0.00000001f

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
cl_mem r_mem_obj;
cl_mem p_mem_obj;
cl_mem q_mem_obj;
cl_mem s_mem_obj;

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

void verify(float* A, float* r, float* p, float* s, float* s_outputFromGpu, float* q, float* q_outputFromGpu, int NX, int NY)
{
	int fail;
	fail = 0;

	int i,j;
	
  	for (i = 0; i < NY; i++)
	{
		s[i] = 0.0;
	}

	for (i = 0; i < NX; i++)
	{
		q[i] = 0.0;
		for (j = 0; j < NY; j++)
	  	{
	    		s[j] = s[j] + r[i] * A[i*NY + j];
	    		q[i] = q[i] + A[i*NY + j] * p[j];
	  	}
	}

	// Compare s with s_cuda
	for (i=0; i<NX; i++)
	{
		if (percentDiff(q[i], q_outputFromGpu[i]) > PERCENT_DIFF_ERROR_THRESHOLD)
		{
			fail++;
		}
	}

	for (i=0; i<NY; i++)
	{
		if (percentDiff(s[i], s_outputFromGpu[i]) > PERCENT_DIFF_ERROR_THRESHOLD)
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
	fp = fopen("bicg.cl", "r");
	if (!fp) {
		fprintf(stderr, "Failed to load kernel.\n");
		exit(1);
	}
	source_str = (char*)malloc(MAX_SOURCE_SIZE);
	source_size = fread( source_str, 1, MAX_SOURCE_SIZE, fp);
	fclose( fp );
}


void init_array(float *A, float *p, float *r)
{
	
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
		
	// Create the 1st OpenCL kernel
	clKernel1 = clCreateKernel(clProgram, "bicgKernel1", &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating kernel\n");

	// Create the 2nd OpenCL kernel
	clKernel2 = clCreateKernel(clProgram, "bicgKernel2", &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating kernel\n");

	clFinish(clCommandQue);
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
	errcode = clReleaseMemObject(p_mem_obj);
	errcode = clReleaseMemObject(q_mem_obj);
	errcode = clReleaseMemObject(r_mem_obj);
	errcode = clReleaseMemObject(s_mem_obj);
	errcode = clReleaseCommandQueue(clCommandQue);
	errcode = clReleaseContext(clGPUContext);
	if(errcode != CL_SUCCESS) printf("Error in cleanup\n");
}




int main(int argc, char** argv) 
{
	if(argc != 2){
		printf("./kernel NX=NY\n");
		exit(-1);
	}
	
	int NX;
	int NY;
	NX = atoi(argv[1]);
	NY = atoi(argv[1]);
	float* A;
	float* r;
	float* s;
	float* p;
	float* q;
	float* s_outputFromGpu;
	float* q_outputFromGpu;
 	
	A = (float*)malloc(NX*NY*sizeof(float));
	r = (float*)malloc(NX*sizeof(float));
	s = (float*)malloc(NY*sizeof(float));
	p = (float*)malloc(NY*sizeof(float));
	q = (float*)malloc(NX*sizeof(float));
	s_outputFromGpu = (float*)malloc(NY*sizeof(float));
	q_outputFromGpu = (float*)malloc(NX*sizeof(float));
		

  	for (int i = 0; i < NX; i++)
	{
    		r[i] = i * M_PI;

    		for (int j = 0; j < NY; j++)
		{
      			A[i*NY + j] = ((float) i*j) / NX;
		}
 	}
	
	for (int i = 0; i < NY; i++)
	{
    		p[i] = i * M_PI;
	}
	read_cl_file();
	cl_initialization();
	
	a_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(float) * NX * NY, NULL, &errcode);
	r_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(float) * NX, NULL, &errcode);
	s_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(float) * NX, NULL, &errcode);
	p_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(float) * NX, NULL, &errcode);
	q_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(float) * NX, NULL, &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating buffers\n");
	
	errcode = clEnqueueWriteBuffer(clCommandQue, a_mem_obj, CL_TRUE, 0, sizeof(float) * NX * NY, A, 0, NULL, NULL);
	errcode = clEnqueueWriteBuffer(clCommandQue, r_mem_obj, CL_TRUE, 0, sizeof(float) * NX, r, 0, NULL, NULL);
	errcode = clEnqueueWriteBuffer(clCommandQue, s_mem_obj, CL_TRUE, 0, sizeof(float) * NX, s, 0, NULL, NULL);
	errcode = clEnqueueWriteBuffer(clCommandQue, p_mem_obj, CL_TRUE, 0, sizeof(float) * NX, p, 0, NULL, NULL);
	errcode = clEnqueueWriteBuffer(clCommandQue, q_mem_obj, CL_TRUE, 0, sizeof(float) * NX, q, 0, NULL, NULL);
	if(errcode != CL_SUCCESS)printf("Error in writing buffers\n");	
	
	cl_load_prog();

	uint64_t t_start, t_end, duration;

	size_t local, global;
	local=256;
	global = (size_t)(((float)NX) / ((float)local)+1) * local;	
	
	// Set the arguments of the kernel
	errcode =  clSetKernelArg(clKernel1, 0, sizeof(cl_mem), (void *)&a_mem_obj);
	errcode |= clSetKernelArg(clKernel1, 1, sizeof(cl_mem), (void *)&p_mem_obj);
	errcode |= clSetKernelArg(clKernel1, 2, sizeof(cl_mem), (void *)&q_mem_obj);
	errcode |= clSetKernelArg(clKernel1, 3, sizeof(int), &NX);
        errcode |= clSetKernelArg(clKernel1, 4, sizeof(int), &NY);
	if(errcode != CL_SUCCESS) printf("Error in seting arguments\n");

	t_start = getTimeInNSecs();
	// Execute the 1st OpenCL kernel
	errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel1, 1, NULL, &global, &local, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in launching kernel\n");
	t_end = getTimeInNSecs();
	clFinish(clCommandQue);
	duration = t_end - t_start;
	
	global = (size_t)(((float)NY) / ((float)local)+1) * local;

	// Set the arguments of the kernel
	errcode =  clSetKernelArg(clKernel2, 0, sizeof(cl_mem), (void *)&a_mem_obj);
	errcode |= clSetKernelArg(clKernel2, 1, sizeof(cl_mem), (void *)&r_mem_obj);
	errcode |= clSetKernelArg(clKernel2, 2, sizeof(cl_mem), (void *)&s_mem_obj);
	errcode |= clSetKernelArg(clKernel2, 3, sizeof(int), &NX);
        errcode |= clSetKernelArg(clKernel2, 4, sizeof(int), &NY);
	if(errcode != CL_SUCCESS) printf("Error in seting arguments\n");

	t_start = getTimeInNSecs();	
	// Execute the 2nd OpenCL kernel
	errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel2, 1, NULL, &global, &local, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in launching kernel\n");
	t_end = getTimeInNSecs();
	clFinish(clCommandQue);

	duration += t_end - t_start;
	fprintf(stdout, "GPU Runtime: %f\n", duration/1e9);

	errcode = clEnqueueReadBuffer(clCommandQue, s_mem_obj, CL_TRUE, 0, NY*sizeof(float), s_outputFromGpu, 0, NULL, NULL);
	errcode = clEnqueueReadBuffer(clCommandQue, q_mem_obj, CL_TRUE, 0, NX*sizeof(float), q_outputFromGpu, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in reading GPU mem\n");  

	verify(A, r, p, s, s_outputFromGpu, q, q_outputFromGpu, NX, NY);

	cl_clean_up();
	
	free(A);
	free(r);
	free(s);
	free(p);
	free(q);
	free(s_outputFromGpu);
	free(q_outputFromGpu);
	
    	return 0;
}

