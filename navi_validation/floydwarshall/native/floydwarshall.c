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
cl_mem path_mem_obj, distance_mem_obj;
FILE *fp;
char *source_str;
size_t source_size;


uint64_t getTimeInNSecs(){
    struct timespec time;
    clock_gettime(CLOCK_MONOTONIC, &time);
    uint64_t timeInNSec = time.tv_sec * 1e9 + time.tv_nsec;
    return timeInNSec;
}

void read_cl_file()
{
	// Load the kernel source code into the array source_str
	fp = fopen("floydwarshall.cl", "r");
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
	clKernel = clCreateKernel(clProgram, "floydWarshallPass", &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating kernel\n");

}

void cl_clean_up()
{
	// Clean up
	errcode = clFlush(clCommandQue);
	errcode = clFinish(clCommandQue);
	errcode = clReleaseKernel(clKernel);
	errcode = clReleaseProgram(clProgram);
	errcode = clReleaseMemObject(path_mem_obj);
	errcode = clReleaseMemObject(distance_mem_obj);
	errcode = clReleaseCommandQueue(clCommandQue);
	errcode = clReleaseContext(clGPUContext);
	if(errcode != CL_SUCCESS) printf("Error in cleanup\n");
}

int main(int argc, char **argv) 
{
	if(argc != 2){
		printf("./kernel The number of nodes in the graph\n");
		exit(-1);
	}
	int N;
	int numIterations = 0;
	N = atoi(argv[1]);
	int* hOutputPathMatrix;
	int* hOutputPathDistanceMatrix;
	int* hVerificationPathMatrix;
	int* hVerificationPathDistanceMatrix;

	hOutputPathMatrix = (int*)malloc(N*N*sizeof(int));
	hOutputPathDistanceMatrix = (int*)malloc(N*N*sizeof(int));
	hVerificationPathMatrix = (int*)malloc(N*N*sizeof(int));
	hVerificationPathDistanceMatrix = (int*)malloc(N*N*sizeof(int));
	
	if(numIterations == 0 || numIterations > N){
		numIterations = N;
	}
	
	for (int i = 0; i < N; i++){
		for(int j = 0; j < i; j++){
			int temp = rand()%10;
			hOutputPathDistanceMatrix[i*N+j] = temp;
			hOutputPathDistanceMatrix[j*N+i] = temp;
		}
	}
	for (int i = 0; i < N; i++){
		int iXWidth = i*N;
		hOutputPathDistanceMatrix[iXWidth+i] = 0;
	}
	for (int i = 0; i < N; i++){
		for(int j = 0; j < i; j++){
			hOutputPathMatrix[i*N+j] = i;
			hOutputPathMatrix[j*N+i] = j;
		}
		hOutputPathMatrix[i*N+i] = i;
	}

	for(int i=0; i < N*N; i++){
		hVerificationPathDistanceMatrix[i] = hOutputPathDistanceMatrix[i];
		hVerificationPathMatrix[i] = hOutputPathMatrix[i];
	}
	read_cl_file();
	cl_initialization();
	path_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_ONLY, sizeof(int) * N*N, NULL, &errcode);
	distance_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_ONLY, sizeof(int) * N*N, NULL, &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating buffers\n");
	
	errcode = clEnqueueWriteBuffer(clCommandQue, path_mem_obj, CL_TRUE, 0, sizeof(int) * N*N, hOutputPathMatrix, 0, NULL, NULL);
	errcode = clEnqueueWriteBuffer(clCommandQue, distance_mem_obj, CL_TRUE, 0, sizeof(int) * N*N, hOutputPathDistanceMatrix, 0, NULL, NULL);
	if(errcode != CL_SUCCESS)printf("Error in writing buffers\n");

	cl_load_prog();

	uint64_t t_start, t_end, duration=0;
	int blockSize = 8;
	if (N%blockSize != 0){
		N = (N/blockSize + 1) * blockSize;
	}
	size_t local[2], global[2];
	local[0] = blockSize;
	local[1] = blockSize;
	global[0] = N;
	global[1] = N;

	for(int k = 0; k < numIterations; k++){
		int pass = k;
		// Set the arguments of the kernel
		errcode =  clSetKernelArg(clKernel, 0, sizeof(cl_mem), (void *)&distance_mem_obj);
		errcode =  clSetKernelArg(clKernel, 1, sizeof(cl_mem), (void *)&path_mem_obj);
		errcode |= clSetKernelArg(clKernel, 2, sizeof(int), (void *)&N);
		errcode |= clSetKernelArg(clKernel, 3, sizeof(int), (void *)&pass);
		if(errcode != CL_SUCCESS) printf("Error in seting arguments\n");
		t_start = getTimeInNSecs();
		errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel, 2, NULL, global, local, 0, NULL, NULL);
		if(errcode != CL_SUCCESS) printf("Error in launching kernel\n");
		clFinish(clCommandQue);
		t_end = getTimeInNSecs();
		duration += t_end - t_start;
	}	
	fprintf(stdout, "GPU Runtime: %f\n", duration/1e9);

	errcode = clEnqueueReadBuffer(clCommandQue, path_mem_obj, CL_TRUE, 0, N*N*sizeof(int), hOutputPathMatrix, 0, NULL, NULL);
	errcode = clEnqueueReadBuffer(clCommandQue, distance_mem_obj, CL_TRUE, 0, N*N*sizeof(int), hOutputPathDistanceMatrix, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in reading GPU mem\n");

	int distanceYtoX, distanceYtoK, distanceKtoX, indirectDistance;
	int width = N, yXwidth;

	for (int k = 0; k < numIterations; k++) {
		for (int y = 0; y < N; y++) {
			yXwidth = y * N;
			for (int x = 0; x < N; x++) {
				distanceYtoX = hVerificationPathDistanceMatrix[yXwidth+x];
				distanceYtoK = hVerificationPathDistanceMatrix[yXwidth+k];
				distanceKtoX = hVerificationPathDistanceMatrix[k*width+x];

				indirectDistance = distanceYtoK + distanceKtoX;

				if (indirectDistance < distanceYtoX) {
					hVerificationPathDistanceMatrix[yXwidth+x] = indirectDistance;
					hVerificationPathMatrix[yXwidth+x] = k;
				}
			}
		}
	}

	bool mismatch;
	for (int i = 0; i < N; i++) {
		for (int j = 0; j < N; j++) {
			if (hOutputPathMatrix[i*N+j] != hVerificationPathMatrix[i*N+j]) {
				mismatch = true;
				printf("Mismatch at row %d col %d, expected %d got %d\n", i, j,
					hVerificationPathMatrix[i*N+j],
					hOutputPathMatrix[i*N+j]);
			}
			if (hOutputPathDistanceMatrix[i*N+j] != hVerificationPathDistanceMatrix[i*N+j]) {
				mismatch = true;
				printf("Mismatch at row %d col %d, expected %d got %d\n", i, j,
					hVerificationPathDistanceMatrix[i*N+j],
					hOutputPathDistanceMatrix[i*N+j]);
			}
		}
	}
	if(!mismatch){
		printf("Passed!\n");	
	}
	cl_clean_up();
	
	free(hOutputPathMatrix);
	free(hOutputPathDistanceMatrix);
	free(hVerificationPathMatrix);
	free(hVerificationPathDistanceMatrix);
	
    return 0;
}

