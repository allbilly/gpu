/**
 * spmv.c
 * 
 * most code converted to c from https://github.com/vetter/shoc/blob/master/src/opencl/level1/spmv/Spmv.cpp
 *  code simplified not all kernels are being used
 */

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <sys/time.h>
#include <math.h>
#include <assert.h>
#include <CL/cl.h>


//define the error threshold for the results "not matching"
#define PERCENT_DIFF_ERROR_THRESHOLD 0.05
static const int BLOCK_SIZE =100;

#define MAX_SOURCE_SIZE (0x100000)

// alignment factor in terms of number of floats, used to enforce
// memory coalescing
static const int PAD_FACTOR = 1;


#define SMALL_FLOAT_VAL 0.00000001f

char str_temp[1024];

cl_platform_id platform_id;
cl_device_id device_id;
cl_device_id* devices;
cl_uint num_devices;
cl_uint num_platforms;
cl_int errcode;
cl_context clGPUContext;
cl_kernel clKernel1;
cl_kernel clKernel2;
cl_command_queue clCommandQue;
cl_program clProgram;


// Device data structures
cl_mem d_val, d_vec, d_out; // floating point
cl_mem d_cols, d_rowLengths, d_rowDelimiters; // integer


FILE *fp;
char *source_str;
size_t source_size;

uint64_t getTimeInNSecs()
{
	struct timespec time;
	clock_gettime(CLOCK_MONOTONIC, &time);
	uint64_t timeInNSec = time.tv_sec * 1e9 + time.tv_nsec;
	return timeInNSec;
}

// ****************************************************************************
// Function: spmvCpu
//
// Purpose:
//   Runs sparse matrix vector multiplication on the CPU
//
// Arguements:
//   val: array holding the non-zero values for the matrix
//   cols: array of column indices for each element of A
//   rowDelimiters: array of size dim+1 holding indices to rows of A;
//                  last element is the index one past the last
//                  element of A
//   vec: dense vector of size dim to be used for multiplication
//   dim: number of rows/columns in the matrix
//   out: input - buffer of size dim
//        output - result from the spmv calculation
//
// Programmer: Lukasz Wesolowski
// Creation: June 23, 2010
// Returns:
//   nothing directly
//   out indirectly through a pointer
// ****************************************************************************
# define floatType float

void spmvCpu(const floatType *val, const int *cols, const int *rowDelimiters,
	     const floatType *vec, int dim, floatType *out)
{

    for (int i=0; i<dim; i++)
    {
        floatType t = 0;
        for (int j=rowDelimiters[i]; j<rowDelimiters[i+1]; j++)
        {
            int col = cols[j];
            t += val[j] * vec[col];
        }
        out[i] = t;
    }

}

// ****************************************************************************
// Function: verifyResults
//
// Purpose:
//   Verifies correctness of GPU results by comparing to CPU results
//
// Arguments:
//   cpuResults: array holding the CPU result vector
//   gpuResults: array hodling the GPU result vector
//   size: number of elements per vector
//   pass: optional iteration number
//
// Programmer: Lukasz Wesolowski
// Creation: June 23, 2010
// Returns:
//   nothing
//   prints "Passed" if the vectors agree within a relative error of
//   MAX_RELATIVE_ERROR and "FAILED" if they are different
// ****************************************************************************

int verifyResults(const floatType *cpuResults, const floatType *gpuResults,
                   const int size, const int pass)
{
	printf("verifying \n");
	printf("%f \n", fabs(cpuResults[0] - gpuResults[0]) / cpuResults[0]);
    int passed = 1;
    for (int i=0; i<size; i++)
    {
        if (fabs(cpuResults[i] - gpuResults[i]) / cpuResults[i]
            > PERCENT_DIFF_ERROR_THRESHOLD)
        {
#ifdef VERBOSE_OUTPUT
           cout << "Mismatch at i: "<< i << " ref: " << cpuResults[i] <<
                " dev: " << gpuResults[i] << endl;
#endif
            passed = 0;
        }
    }/*	
    if (pass != -1)
    {
        cout << "Test ";
    }
    if (passed)
    {
        cout << "Passed" << endl;
    }
    else
    {
        cout << "Failed" << endl;
    } */
    return passed;
}
// ****************************

// ****************************************************************************
// Function: fill
//
// Purpose:
//   Simple routine to initialize input array
//
// Arguments:
//   A: pointer to the array to initialize
//   n: number of elements in the array
//   maxi: specifies range of random values
//
// Programmer: Lukasz Wesolowski
// Creation: June 21, 2010
// Returns:  nothing
//
// ****************************************************************************

void fill(float *A, const int n, const float maxi)
{
    for (int j = 0; j < n; j++)
    {
        A[j] = ((floatType) maxi * (rand() / (RAND_MAX + 1.0f)));
		
    }
}

// ****************************************************************************
// Function initRandomMatrix
//
// Purpose:
//   Assigns random positions to a given number of elements in a square
//   matrix, A.  The function encodes these positions in compressed sparse
//   row format.
//
// Arguments:
//   cols:          array for column indexes of elements (size should be = n)
//   rowDelimiters: array of size dim+1 holding indices to rows of A;
//                  last element is the index one past the last element of A
//   n:             number of nonzero elements in A
//   dim:           number of rows/columns in A
//
// Programmer: Kyle Spafford
// Creation: July 28, 2010
// Returns: nothing
//
// ****************************************************************************

void initRandomMatrix(int *cols, int *rowDelimiters, const int n, const int dim)
{
    int nnzAssigned = 0;

    // Figure out the probability that a nonzero should be assigned to a given
    // spot in the matrix
    double prob = (double)n / ((double)dim * (double)dim);

    // Seed random number generator
    srand48(8675309L);

    // Randomly decide whether entry i,j gets a value, but ensure n values
    // are assigned
    int fillRemaining = 0;
    for (int i = 0; i < dim; i++)
    {
        rowDelimiters[i] = nnzAssigned;
        for (int j = 0; j < dim; j++)
        {
            int numEntriesLeft = (dim * dim) - ((i * dim) + j);
            int needToAssign   = n - nnzAssigned;
            if (numEntriesLeft <= needToAssign) {
                fillRemaining = 1;
            }
            if ((nnzAssigned < n && drand48() <= prob) || fillRemaining)
            {
                // Assign (i,j) a value
                cols[nnzAssigned] = j;
                nnzAssigned++;
            }
        }
    }
    // Observe the convention to put the number of non zeroes at the end of the
    // row delimiters array
    rowDelimiters[dim] = n;
    assert(nnzAssigned == n);
}

// ****************************************************************************
// Function: convertToPadded
//
// Purpose: pads a CSR matrix with zeros so that each line of values
//          for the matrix is aligned to PAD_FACTOR*4 bytes
//
// Arguments:
//   A: array holding the non-zero values for the matrix
//   cols: array of column indices of the sparse matrix
//   dim: number of rows/columns in the matrix
//   rowDelimiters: array holding indices in A to rows of the sparse matrix
//   newA_ptr: input - pointer to an uninitialized pointer
//             output - pointer to padded A
//   newcols_ptr: input - pointer to an uninitialized pointer
//                output - pointer to padded cols
//   newIndices: input - buffer of size dim + 1
//               output - array holding indices in newA to rows of the
//                        sparse matrix
//   newSize: input - pointer to uninitialized int
//            output - pointer to the size of A
//
// Programmer: Lukasz Wesolowski
// Creation: July 8, 2010
// Returns:
//   nothing directly
//   allocates and returns *newA_ptr and *newcols_ptr indirectly
//   returns newIndices and newSize indirectly through pointers
// ******************************

void convertToPadded(floatType *A, int *cols, int dim, int *rowDelimiters,
                     floatType **newA_ptr, int **newcols_ptr, int *newIndices,
                     int *newSize)
{

    // determine total padded size and new row indices
    int paddedSize = 0;
    int rowSize;

    for (int i=0; i<dim; i++)
    {
        newIndices[i] = paddedSize;
        rowSize = rowDelimiters[i+1] - rowDelimiters[i];
        if (rowSize % PAD_FACTOR != 0)
        {
            rowSize += PAD_FACTOR - rowSize % PAD_FACTOR;
        }
        paddedSize += rowSize;
    }
    *newSize = paddedSize;
    newIndices[dim] = paddedSize;

    //*newA_ptr = pmsAllocHostBuffer<floatType>( paddedSize ); //just create a new array

	*newA_ptr = (floatType *) malloc( paddedSize * sizeof(paddedSize));

    //*newcols_ptr = pmsAllocHostBuffer<int>( paddedSize ); 

	*newcols_ptr = (int *)malloc(paddedSize * sizeof(int));

    floatType *newA = *newA_ptr;
    int *newcols = *newcols_ptr;

    memset(newA, 0, paddedSize * sizeof(floatType));
    memset(newcols, 0, paddedSize*sizeof(int));

    // fill newA and newcols
    for (int i=0; i<dim; i++)
    {
        for (int j=rowDelimiters[i], k=newIndices[i]; j<rowDelimiters[i+1];
             j++, k++)
        {
            newA[k] = A[j];
            newcols[k] = cols[j];
        }
    }

}

void read_cl_file()
{
	// Load the kernel source code into the array source_str
	fp = fopen("spmv.cl", "r");
	if (!fp)
	{
		fprintf(stderr, "Failed to load kernel.\n");
		exit(1);
	}
	source_str = (char *)malloc(MAX_SOURCE_SIZE);
	source_size = fread(source_str, 1, MAX_SOURCE_SIZE, fp);
	fclose(fp);
}

void init_array(float *A, float *p, float *r)
{
}

void cl_initialization()
{
	// Get platform and device information
	errcode = clGetPlatformIDs(1, &platform_id, &num_platforms);
	if (errcode == CL_SUCCESS)
		printf("number of platforms is %d\n", num_platforms);
	else
		printf("Error getting platform IDs\n");

	errcode = clGetPlatformInfo(platform_id, CL_PLATFORM_NAME, sizeof(str_temp), str_temp, NULL);
	if (errcode == CL_SUCCESS)
		printf("platform name is %s\n", str_temp);
	else
		printf("Error getting platform name\n");

	errcode = clGetPlatformInfo(platform_id, CL_PLATFORM_VERSION, sizeof(str_temp), str_temp, NULL);
	if (errcode == CL_SUCCESS)
		printf("platform version is %s\n", str_temp);
	else
		printf("Error getting platform version\n");

	errcode = clGetDeviceIDs(platform_id, CL_DEVICE_TYPE_GPU, 1, &device_id, &num_devices);
	if (errcode == CL_SUCCESS)
		printf("number of devices is %d\n", num_devices);
	else
		printf("Error getting device IDs\n");

	errcode = clGetDeviceInfo(device_id, CL_DEVICE_NAME, sizeof(str_temp), str_temp, NULL);
	if (errcode == CL_SUCCESS)
		printf("device name is %s\n", str_temp);
	else
		printf("Error getting device name\n");

	// Create an OpenCL context
	clGPUContext = clCreateContext(NULL, 1, &device_id, NULL, NULL, &errcode);
	if (errcode != CL_SUCCESS)
		printf("Error in creating context\n");

	//Create a command-queue
	clCommandQue = clCreateCommandQueue(clGPUContext, device_id, 0, &errcode);
	if (errcode != CL_SUCCESS)
		printf("Error in creating command queue\n");
}

void cl_load_prog()
{
	// Create a program from the kernel source
	clProgram = clCreateProgramWithSource(clGPUContext, 1, (const char **)&source_str, (const size_t *)&source_size, &errcode);

	if (errcode != CL_SUCCESS)
	{

		printf("Error in creating program\n");
	}

	// Build the program
	errcode = clBuildProgram(clProgram, 1, &device_id, NULL, NULL, NULL);
	//if(errcode != CL_SUCCESS) printf("Error in building program\n");
	if (errcode != CL_SUCCESS)
	{
		char *buff_erro;
		cl_int errcode;
		size_t build_log_len;
		errcode = clGetProgramBuildInfo(clProgram, device_id, CL_PROGRAM_BUILD_LOG, 0, NULL, &build_log_len);
		if (errcode)
		{
			printf("clGetProgramBuildInfo failed at line %d\n", __LINE__);
			exit(-1);
		}

		buff_erro = malloc(build_log_len);
		if (!buff_erro)
		{
			printf("malloc failed at line %d\n", __LINE__);
			exit(-2);
		}

		errcode = clGetProgramBuildInfo(clProgram, device_id, CL_PROGRAM_BUILD_LOG, build_log_len, buff_erro, NULL);
		if (errcode)
		{
			printf("clGetProgramBuildInfo failed at line %d\n", __LINE__);
			exit(-3);
		}

		fprintf(stderr, "Build log: \n%s\n", buff_erro); //Be careful with  the fprint
		free(buff_erro);
		fprintf(stderr, "clBuildProgram failed\n");
		exit(EXIT_FAILURE);
	}
	// Create the 1st OpenCL kernel
	clKernel1 = clCreateKernel(clProgram, "spmv_csr_scalar_kernel", &errcode);
	if (errcode != CL_SUCCESS)
		printf("Error in creating kernel\n");

	//Create the 2nd OpenCL kernel
	clKernel2 = clCreateKernel(clProgram, "spmv_csr_vector_kernel", &errcode);
	if(errcode != CL_SUCCESS) 
		printf("Error in creating kernel\n");

	clFinish(clCommandQue);
}

void cl_clean_up()
{
	// Clean up
	errcode = clFlush(clCommandQue);
	errcode = clFinish(clCommandQue);
	errcode = clReleaseKernel(clKernel1);
	//errcode = clReleaseKernel(clKernel2);
	errcode = clReleaseProgram(clProgram);
	errcode = clReleaseMemObject(d_val);
	errcode = clReleaseMemObject(d_out);
	errcode = clReleaseMemObject(d_vec);
	errcode = clReleaseMemObject(d_cols);
	errcode = clReleaseMemObject(d_rowDelimiters);
	errcode = clReleaseMemObject(d_rowLengths);

	errcode = clReleaseCommandQueue(clCommandQue);
	errcode = clReleaseContext(clGPUContext);
	if (errcode != CL_SUCCESS)
		printf("Error in cleanup\n");
}

#define floatType float



int main(int argc, char **argv)
{

   // Host data structures
    // array of values in the sparse matrix
    //floatType *h_val;
	floatType *h_valPad;
    // array of column indices for each value in h_val
    //int *h_cols, *h_colsPad;
    // array of indices to the start of each row in h_val/valPad
    //int *h_rowDelimiters, *h_rowDelimitersPad;
    // Dense vector of values
    //floatType *h_vec;
    // Output vector
    //floatType *h_out;
    // Reference solution computed by cpu
    //floatType *refOut;

    int nItems;            // number of non-zero elements in the matrix
    int nItemsPadded;
    int numRows;           // number of rows in the matrix


	// generate random matrix
	 // If we're not opening a file, the dimension of the matrix
    // has been passed in as an argument
	int nRows = 100;

    numRows = nRows;
    nItems = numRows * numRows / 100; // 1% of entries will be non-zero
    //float maxval = op.getOptionFloat("maxval");

	printf("%d \n", nItems);
	fflush(stdout); 
	float maxval = 100.0;

	floatType h_val[nItems];
	int h_cols[nItems];
    int h_rowDelimiters[nRows + 1];

	// fill data
    fill(h_val, nItems, maxval);
    initRandomMatrix(h_cols, h_rowDelimiters, nItems, numRows);


    // Set up remaining host data
	floatType h_vec[numRows];
	floatType refOut[numRows];
	int h_rowDelimitersPad[numRows + 1];

    fill(h_vec, numRows, maxval);

	// Set up the padded data structures
    int paddedSize = numRows + (PAD_FACTOR - numRows % PAD_FACTOR);
    int * h_colsPad;
    floatType h_out[paddedSize];

    convertToPadded(h_val, h_cols, numRows, h_rowDelimiters, &h_valPad,
            &h_colsPad, h_rowDelimitersPad, &nItemsPadded);

    // Compute reference solution
    spmvCpu(h_val, h_cols, h_rowDelimiters, h_vec, numRows, refOut);

	// read cl and initialize
	read_cl_file();
	cl_initialization();

	int numNonZeroes = nItems;
	int cmSize = 1;

	#define clFloatType cl_float;

	//. create buffer for all items
	/// Allocate device memory
    d_val = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, numNonZeroes * sizeof(cl_float), NULL, &errcode);
	d_cols = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, numNonZeroes * sizeof(cl_int), NULL, &errcode);
    d_vec = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, numRows * sizeof(cl_float), NULL, &errcode);
    d_out = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, numRows * sizeof(cl_float), NULL, &errcode);
    d_rowDelimiters = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, (numRows+1) * sizeof(cl_int), NULL, &errcode);
	d_rowLengths = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE, cmSize * sizeof(int), NULL, &errcode);


	// Transfer data to device
	errcode = clEnqueueWriteBuffer(clCommandQue, d_val, CL_TRUE, 0, numNonZeroes * sizeof(floatType), h_val, 0, NULL, NULL);
	if (errcode != CL_SUCCESS)
		printf("Error in writing buffers\n");
	errcode = clEnqueueWriteBuffer(clCommandQue, d_cols, CL_TRUE, 0, numNonZeroes *sizeof(int), h_cols, 0, NULL, NULL);
	if (errcode != CL_SUCCESS)
		printf("Error in writing buffers\n");
	errcode = clEnqueueWriteBuffer(clCommandQue, d_vec, CL_TRUE, 0, numRows * sizeof(floatType), h_vec, 0, NULL, NULL);
	if (errcode != CL_SUCCESS)
		printf("Error in writing buffers\n");
	errcode = clEnqueueWriteBuffer(clCommandQue, d_rowDelimiters, CL_TRUE, 0, (numRows+1) * sizeof(int), h_rowDelimiters, 0, NULL, NULL);
	if (errcode != CL_SUCCESS)
		printf("Error in writing buffers\n");
	errcode = clFinish(clCommandQue);
	if (errcode != CL_SUCCESS)
		printf("Error in writing buffers\n");

	printf("Laoding program \n");
	cl_load_prog();
	
	uint64_t t_start, t_end, duration;


	// Set up CSR Kernels
	
	//clKernel1  = clCreateKernel(prog, "spmv_csr_scalar_kernel", &errcode);
	//CL_CHECK_ERROR(errcode);
	errcode = clSetKernelArg(clKernel1, 0, sizeof(cl_mem), (void*) &d_val);
	//CL_CHECK_ERROR(errcode);
	if (errcode != CL_SUCCESS)
		printf("Error in seting arguments\n");
	errcode = clSetKernelArg(clKernel1, 1, sizeof(cl_mem), (void*) &d_vec);
	//CL_CHECK_ERROR(errcode);
	if (errcode != CL_SUCCESS)
		printf("Error in seting arguments\n");
	errcode = clSetKernelArg(clKernel1, 2, sizeof(cl_mem), (void*) &d_cols);
	//CL_CHECK_ERROR(errcode);
	if (errcode != CL_SUCCESS)
		printf("Error in seting arguments\n");
	errcode = clSetKernelArg(clKernel1, 3, sizeof(cl_mem), (void*) &d_rowDelimiters);
	//CL_CHECK_ERROR(errcode);
	errcode = clSetKernelArg(clKernel1, 4, sizeof(cl_int), (void*) &numRows);
	//CL_CHECK_ERROR(errcode);
	errcode = clSetKernelArg(clKernel1, 5, sizeof(cl_mem), (void*) &d_out);
    //CL_CHECK_ERROR(errcode);

	if (errcode != CL_SUCCESS) {
		printf("Error in seting arguments\n");
	} else { printf (" Success");};
	
	t_start = getTimeInNSecs();
	// Execute the 1st OpenCL kernel
	//errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel1, 1, NULL, &global, &local, 0, NULL, NULL);

	const size_t scalarGlobalWSize = numRows;
	size_t localWorkSize = BLOCK_SIZE;

	printf("Setting args \n");
	printf("%d \n",CL_DEVICE_MAX_WORK_GROUP_SIZE);
	errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel1, 1, NULL,
			&scalarGlobalWSize, &localWorkSize, 0, NULL,
			NULL);
	if (errcode != CL_SUCCESS) printf("Error in launching kernel %d\n", errcode);

	// Transfer data back to host
    printf("Transfering data back to host \n");      
    errcode = clEnqueueReadBuffer(clCommandQue, d_out, 1, 0, numRows *sizeof(floatType), h_out, 0, NULL, NULL);

	if (errcode != CL_SUCCESS) printf("Error in transfering back results %d \n", errcode);

	int vfy;
	vfy = verifyResults(refOut, h_out, numRows, 1);

	if (vfy == 1) printf("Passed \n");

	t_end = getTimeInNSecs();
	clFinish(clCommandQue);
	duration = t_end - t_start;

	//global = (size_t)(((float)NY) / ((float)local)+1) * local;

	//verify(A, r, p, s, s_outputFromGpu, q, q_outputFromGpu, NX, NY);

	cl_clean_up();

	free(d_val);
	free(d_out);
	free(d_vec);
	free(d_cols);
	free(d_rowDelimiters);
	free(d_rowLengths);


	return 0;
}
