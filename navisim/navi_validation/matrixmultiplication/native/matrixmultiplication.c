#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <sys/time.h>
#include <time.h>
#include <stdbool.h>

#include <CL/cl.h>
#include "../../common/polybenchUtilFuncts.h"
#define MAX_SOURCE_SIZE (0x100000)

/* Problem size. */
#define X 64
#define Y 64
#define Z 64


/* Thread block dimensions */
#define DIM_LOCAL_WORK_GROUP_X 64


/* Can switch DATA_TYPE between float and double */
typedef float DATA_TYPE;

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
cl_mem a_mem_obj;
cl_mem b_mem_obj;
cl_mem c_mem_obj;
FILE *fp;
char *source_str;
size_t source_size;

typedef struct Matrix{
	DATA_TYPE *data;
	int height;
	int width;
} Matrix;

Matrix newMatrix(int height, int width){
	Matrix matrix;
	matrix.height = height;
	matrix.width = width;
	matrix.data = malloc(height*width*sizeof(DATA_TYPE));
	return matrix;
}

void verify(Matrix *a, Matrix *b, Matrix *gc)
{
	bool mismatch;
	Matrix c = newMatrix(X,Z);
	for(int i=0; i < Z; i++){
		for(int j=0; j < X; j++){
			int indexC = j*Z + i;
			DATA_TYPE sum=0;
			for(int k=0; k < Y; k++){
				int indexA = j*Y + k;
				int indexB = k*Z + i;
				sum += a->data[indexA] * b->data[indexB];
			}
			c.data[indexC] = sum;
		}
	}
	for(int i = 0; i < X; i++){
		for(int j= 0; j < Z; j++){
			int index = i*Z + j;
			if (absVal(gc->data[index] - c.data[index]) > 1e-3){
				mismatch = true;
				//if (index){
					printf("mismatch at (%d,%d), expect %f, output %f\n", i,j,
						c.data[index], gc->data[index]);
					return;
				//}
				
			}
		}
	}
	if(!mismatch){
		printf("Passed!\n");
	}
}


void read_cl_file()
{
	// Load the kernel source code into the array source_str
	fp = fopen("matrixmultiplication.cl", "r");
	if (!fp) {
		fprintf(stderr, "Failed to load kernel.\n");
		exit(1);
	}
	source_str = (char*)malloc(MAX_SOURCE_SIZE);
	source_size = fread( source_str, 1, MAX_SOURCE_SIZE, fp);
	fclose( fp );
}


void init_array(Matrix *a, Matrix *b)
{
	for (int i = 0; i < X; i++){
		for (int j = 0; j < Y; j++){
			a->data[j*X+i] = (float)rand()/(float)RAND_MAX; 
		}
	}
	for (int i = 0; i < Y; i++){
		for (int j = 0; j < Z; j++){
			b->data[j*Y+i] = (float)rand()/(float)RAND_MAX; 
		}
	}
	
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


void cl_mem_init(Matrix *a, Matrix *b)
{
	a_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_ONLY, sizeof(DATA_TYPE) * X * Y, NULL, &errcode);
	b_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_ONLY, sizeof(DATA_TYPE) * Y * Z, NULL, &errcode);
	c_mem_obj = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, sizeof(DATA_TYPE) * X * Z, NULL, &errcode);
		
	if(errcode != CL_SUCCESS) printf("Error in creating buffers\n");
	
	errcode = clEnqueueWriteBuffer(clCommandQue, a_mem_obj, CL_TRUE, 0, sizeof(DATA_TYPE) * X * Y, a->data, 0, NULL, NULL);
	errcode = clEnqueueWriteBuffer(clCommandQue, b_mem_obj, CL_TRUE, 0, sizeof(DATA_TYPE) * Y * Z, b->data, 0, NULL, NULL);
	if(errcode != CL_SUCCESS)printf("Error in writing buffers\n");
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
	clKernel = clCreateKernel(clProgram, "mmmKernel_local", &errcode);
	if(errcode != CL_SUCCESS) printf("Error in creating kernel\n");

}


void cl_launch_kernel()
{
	double t_start, t_end;

	size_t localWorkSize[2], globalWorkSize[2];
	localWorkSize[0] = 8;
	localWorkSize[1] = 8;
	globalWorkSize[0] = Z / 4;
	globalWorkSize[1] = X / 4;
	int width = Y;

	t_start = rtclock();

	// Set the arguments of the kernel
	errcode =  clSetKernelArg(clKernel, 0, sizeof(cl_mem), (void *)&a_mem_obj);
	errcode |= clSetKernelArg(clKernel, 1, sizeof(cl_mem), (void *)&b_mem_obj);
	errcode |= clSetKernelArg(clKernel, 2, sizeof(cl_mem), (void *)&c_mem_obj);
	errcode |= clSetKernelArg(clKernel, 3, sizeof(int), (void *)&width);
	errcode |= clSetKernelArg(clKernel, 4, 32*32*4, NULL);
	if(errcode != CL_SUCCESS) printf("Error in seting arguments\n");
	errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel, 1, NULL, globalWorkSize, localWorkSize, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in launching kernel\n");
	clFinish(clCommandQue);

	t_end = rtclock();
	fprintf(stdout, "GPU Runtime: %0.6lfs\n", t_end - t_start);
}


void cl_clean_up()
{
	// Clean up
	errcode = clFlush(clCommandQue);
	errcode = clFinish(clCommandQue);
	errcode = clReleaseKernel(clKernel);
	errcode = clReleaseProgram(clProgram);
	errcode = clReleaseMemObject(a_mem_obj);
	errcode = clReleaseMemObject(b_mem_obj);
	errcode = clReleaseMemObject(c_mem_obj);
	errcode = clReleaseCommandQueue(clCommandQue);
	errcode = clReleaseContext(clGPUContext);
	if(errcode != CL_SUCCESS) printf("Error in cleanup\n");
}

int main(void) 
{
	double t_start, t_end;

	
	Matrix a = newMatrix(X,Y);
	Matrix b = newMatrix(Y,Z);
	Matrix c = newMatrix(X,Z);
	
	init_array(&a,&b);
	read_cl_file();
	cl_initialization();
	cl_mem_init(&a,&b);
	cl_load_prog();

	cl_launch_kernel();

	errcode = clEnqueueReadBuffer(clCommandQue, c_mem_obj, CL_TRUE, 0, X*Z*sizeof(DATA_TYPE), c.data, 0, NULL, NULL);
	if(errcode != CL_SUCCESS) printf("Error in reading GPU mem\n");

	verify(&a, &b, &c);
	cl_clean_up();
	
    return 0;
}