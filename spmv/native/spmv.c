/**
 * spmv.c
 *
 * most code converted to c from
 * https://github.com/vetter/shoc/blob/master/src/opencl/level1/spmv/Spmv.cpp
 *  code simplified not all kernels are being used
 */

#include <CL/cl.h>
#include <assert.h>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/time.h>
#include <time.h>

// define the error threshold for the results "not matching"
#define PERCENT_DIFF_ERROR_THRESHOLD 0.05
static const int BLOCK_SIZE = 128;

#define MAX_SOURCE_SIZE (0x100000)

// alignment factor in terms of number of floats, used to enforce
// memory coalescing
static const int PAD_FACTOR = 10;

#define SMALL_FLOAT_VAL 0.00000001f

char str_temp[1024];

cl_platform_id platform_id;
cl_device_id device_id;
cl_device_id *devices;
cl_uint num_devices;
cl_uint num_platforms;
cl_int errcode;
cl_context clGPUContext;
cl_kernel clKernel1;
cl_command_queue clCommandQue;
cl_program clProgram;

// Device data structures
cl_mem d_val, d_vec, d_out;     // floating point
cl_mem d_cols, d_rowDelimiters; // integer

FILE *fp;
char *source_str;
size_t source_size;

uint64_t getTimeInNSecs() {
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
#define floatType float

void spmvCpu(const floatType *val, const int *cols, const int *rowDelimiters,
             const floatType *vec, int dim, floatType *out) {

    for (int i = 0; i < dim; i++) {
        floatType t = 0;
        for (int j = rowDelimiters[i]; j < rowDelimiters[i + 1]; j++) {
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
                  const int size) {
    int passed = 1;
    for (int i = 0; i < size; i++) {
        if (fabs(cpuResults[i] - gpuResults[i]) / cpuResults[i] >
            PERCENT_DIFF_ERROR_THRESHOLD) {
            passed = 0;
        }
    }
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

void fill(float *A, const int n, const float maxi) {
    for (int j = 0; j < n; j++) {
        A[j] = ((floatType)maxi * (rand() / (RAND_MAX + 1.0f)));
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

void initRandomMatrix(int *cols, int *rowDelimiters, const int n,
                      const uint64_t dim) {
    int nnzAssigned = 0;

    // Figure out the probability that a nonzero should be assigned to a given
    // spot in the matrix
    double prob = (double)n / ((double)dim * (double)dim);

    // Seed random number generator
    srand48(8675309L);

    // Randomly decide whether entry i,j gets a value, but ensure n values
    // are assigned
    int fillRemaining = 0;
    for (int i = 0; i < dim; i++) {
        rowDelimiters[i] = nnzAssigned;
        for (int j = 0; j < dim; j++) {
            uint64_t total = dim * dim;
            uint64_t cur = i * dim + j;
            uint64_t numEntriesLeft = total - cur;
            if (numEntriesLeft < 0) {
                printf("long is %ld bytes\n", sizeof(long));
                printf("error here\n");
            }
            int needToAssign = n - nnzAssigned;
            if (numEntriesLeft <= needToAssign) {
                fillRemaining = 1;
            }
            if ((nnzAssigned < n && drand48() <= prob) || fillRemaining) {
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

void read_cl_file() {
    // Load the kernel source code into the array source_str
    fp = fopen("spmv.cl", "r");
    if (!fp) {
        fprintf(stderr, "Failed to load kernel.\n");
        exit(1);
    }
    source_str = (char *)malloc(MAX_SOURCE_SIZE);
    source_size = fread(source_str, 1, MAX_SOURCE_SIZE, fp);
    fclose(fp);
}

void init_array(float *A, float *p, float *r) {}

void cl_initialization() {
    // Get platform and device information
    errcode = clGetPlatformIDs(1, &platform_id, &num_platforms);
    if (errcode == CL_SUCCESS)
        printf("number of platforms is %d\n", num_platforms);
    else
        printf("Error getting platform IDs\n");

    errcode = clGetPlatformInfo(platform_id, CL_PLATFORM_NAME, sizeof(str_temp),
                                str_temp, NULL);
    if (errcode == CL_SUCCESS)
        printf("platform name is %s\n", str_temp);
    else
        printf("Error getting platform name\n");

    errcode = clGetPlatformInfo(platform_id, CL_PLATFORM_VERSION,
                                sizeof(str_temp), str_temp, NULL);
    if (errcode == CL_SUCCESS)
        printf("platform version is %s\n", str_temp);
    else
        printf("Error getting platform version\n");

    errcode = clGetDeviceIDs(platform_id, CL_DEVICE_TYPE_GPU, 1, &device_id,
                             &num_devices);
    if (errcode == CL_SUCCESS)
        printf("number of devices is %d\n", num_devices);
    else
        printf("Error getting device IDs\n");

    errcode = clGetDeviceInfo(device_id, CL_DEVICE_NAME, sizeof(str_temp),
                              str_temp, NULL);
    if (errcode == CL_SUCCESS)
        printf("device name is %s\n", str_temp);
    else
        printf("Error getting device name\n");

    // Create an OpenCL context
    clGPUContext = clCreateContext(NULL, 1, &device_id, NULL, NULL, &errcode);
    if (errcode != CL_SUCCESS)
        printf("Error in creating context\n");

    // Create a command-queue
    clCommandQue = clCreateCommandQueue(clGPUContext, device_id,
                                        CL_QUEUE_PROFILING_ENABLE, &errcode);
    if (errcode != CL_SUCCESS)
        printf("Error in creating command queue\n");
}

void cl_load_prog() {
    // Create a program from the kernel source
    clProgram =
        clCreateProgramWithSource(clGPUContext, 1, (const char **)&source_str,
                                  (const size_t *)&source_size, &errcode);

    if (errcode != CL_SUCCESS) {

        printf("Error in creating program\n");
    }

    // Build the program
    errcode = clBuildProgram(clProgram, 1, &device_id, NULL, NULL, NULL);
    // if(errcode != CL_SUCCESS) printf("Error in building program\n");
    if (errcode != CL_SUCCESS) {
        char *buff_erro;
        cl_int errcode;
        size_t build_log_len;
        errcode =
            clGetProgramBuildInfo(clProgram, device_id, CL_PROGRAM_BUILD_LOG, 0,
                                  NULL, &build_log_len);
        if (errcode) {
            printf("clGetProgramBuildInfo failed at line %d\n", __LINE__);
            exit(-1);
        }

        buff_erro = malloc(build_log_len);
        if (!buff_erro) {
            printf("malloc failed at line %d\n", __LINE__);
            exit(-2);
        }

        errcode =
            clGetProgramBuildInfo(clProgram, device_id, CL_PROGRAM_BUILD_LOG,
                                  build_log_len, buff_erro, NULL);
        if (errcode) {
            printf("clGetProgramBuildInfo failed at line %d\n", __LINE__);
            exit(-3);
        }

        fprintf(stderr, "Build log: \n%s\n",
                buff_erro); // Be careful with  the fprint
        free(buff_erro);
        fprintf(stderr, "clBuildProgram failed\n");
        exit(EXIT_FAILURE);
    }
    // Create the 1st OpenCL kernel
    clKernel1 = clCreateKernel(clProgram, "spmv_csr_scalar_kernel", &errcode);
    if (errcode != CL_SUCCESS)
        printf("Error in creating kernel\n");

    clFinish(clCommandQue);
}

void cl_clean_up() {
    // Clean up
    errcode = clFlush(clCommandQue);
    errcode = clFinish(clCommandQue);
    errcode = clReleaseKernel(clKernel1);
    errcode = clReleaseProgram(clProgram);
    errcode = clReleaseMemObject(d_val);
    errcode = clReleaseMemObject(d_out);
    errcode = clReleaseMemObject(d_vec);
    errcode = clReleaseMemObject(d_cols);
    errcode = clReleaseMemObject(d_rowDelimiters);

    errcode = clReleaseCommandQueue(clCommandQue);
    errcode = clReleaseContext(clGPUContext);
    if (errcode != CL_SUCCESS)
        printf("Error in cleanup\n");
}

#define floatType float

int main(int argc, char **argv) {

    uint64_t nItems;  // number of non-zero elements in the matrix
    uint64_t numRows; // number of rows in the matrix
    float sparsity;

    if (argc != 3) {
        printf("./kernel dim sparsity\n");
        exit(-1);
    }
    numRows = atoi(argv[1]);
    sparsity = atof(argv[2]);

    nItems = numRows * numRows;
    printf("numNonZeroes is %ld\n", nItems);
    nItems = nItems * sparsity;
    float maxval = 10.0;
    printf("numNonZeroes is %ld\n", nItems);

    floatType *h_val =
        malloc(sizeof(floatType) * nItems); // array of values in sparse matrix
    int *h_cols =
        malloc(sizeof(int) *
               nItems); // // array of column indices for each value in h_val
    int *h_rowDelimiters =
        malloc(sizeof(int) * (numRows + 1)); // array of indices to the start of
                                             // each row in h_val/

    // fill data
    fill(h_val, nItems, maxval);
    initRandomMatrix(h_cols, h_rowDelimiters, nItems, numRows);

    // Set up remaining host data
    floatType *h_vec =
        malloc(sizeof(floatType) * numRows); // Dense vector of values
    floatType *refOut = malloc(sizeof(floatType) *
                               numRows); // Reference solution computed by cpu

    fill(h_vec, numRows, maxval);

    floatType *h_out = malloc(sizeof(floatType) * numRows); // host output

    // Compute reference solution
    spmvCpu(h_val, h_cols, h_rowDelimiters, h_vec, numRows, refOut);

    // read cl and initialize
    read_cl_file();
    cl_initialization();

    int numNonZeroes = nItems;

#define clFloatType cl_float;

    //. create buffer for all items
    /// Allocate device memory

    d_val = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE,
                           numNonZeroes * sizeof(cl_float), NULL, &errcode);
    if (errcode != CL_SUCCESS)
        printf("Error in creating buffers: %d\n", errcode);
    d_cols = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE,
                            numNonZeroes * sizeof(cl_int), NULL, &errcode);
    d_vec = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE,
                           numRows * sizeof(cl_float), NULL, &errcode);
    d_out = clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE,
                           numRows * sizeof(cl_float), NULL, &errcode);
    d_rowDelimiters =
        clCreateBuffer(clGPUContext, CL_MEM_READ_WRITE,
                       (numRows + 1) * sizeof(cl_int), NULL, &errcode);

    // Transfer data to device
    errcode = clEnqueueWriteBuffer(clCommandQue, d_val, CL_TRUE, 0,
                                   numNonZeroes * sizeof(floatType), h_val, 0,
                                   NULL, NULL);
    if (errcode != CL_SUCCESS)
        printf("Error in writing buffers: %d\n", errcode);
    errcode =
        clEnqueueWriteBuffer(clCommandQue, d_cols, CL_TRUE, 0,
                             numNonZeroes * sizeof(int), h_cols, 0, NULL, NULL);
    if (errcode != CL_SUCCESS)
        printf("Error in writing buffers: %d\n", errcode);
    errcode =
        clEnqueueWriteBuffer(clCommandQue, d_vec, CL_TRUE, 0,
                             numRows * sizeof(floatType), h_vec, 0, NULL, NULL);
    if (errcode != CL_SUCCESS)
        printf("Error in writing buffers: %d\n", errcode);
    errcode = clEnqueueWriteBuffer(clCommandQue, d_rowDelimiters, CL_TRUE, 0,
                                   (numRows + 1) * sizeof(int), h_rowDelimiters,
                                   0, NULL, NULL);
    if (errcode != CL_SUCCESS)
        printf("Error in writing buffers: %d\n", errcode);
    errcode = clFinish(clCommandQue);
    if (errcode != CL_SUCCESS)
        printf("Error in writing buffers\n");

    cl_load_prog();

    // Set up CSR Kernels

    // set up scalar kernel
    errcode = clSetKernelArg(clKernel1, 0, sizeof(cl_mem), (void *)&d_val);
    if (errcode != CL_SUCCESS)
        printf("Error in seting arguments\n");
    errcode = clSetKernelArg(clKernel1, 1, sizeof(cl_mem), (void *)&d_vec);
    if (errcode != CL_SUCCESS)
        printf("Error in seting arguments\n");
    errcode = clSetKernelArg(clKernel1, 2, sizeof(cl_mem), (void *)&d_cols);
    if (errcode != CL_SUCCESS)
        printf("Error in seting arguments\n");
    errcode =
        clSetKernelArg(clKernel1, 3, sizeof(cl_mem), (void *)&d_rowDelimiters);
    if (errcode != CL_SUCCESS)
        printf("Error in seting arguments\n");
    errcode = clSetKernelArg(clKernel1, 4, sizeof(cl_int), (void *)&numRows);
    if (errcode != CL_SUCCESS)
        printf("Error in seting arguments\n");
    errcode = clSetKernelArg(clKernel1, 5, sizeof(cl_mem), (void *)&d_out);
    if (errcode != CL_SUCCESS)
        printf("Error in seting arguments\n");

    const size_t scalarGlobalWSize = numRows;
    size_t localWorkSize = BLOCK_SIZE;
    cl_event event;
    cl_ulong t_start, t_end, duration = 0;

    errcode = clEnqueueNDRangeKernel(clCommandQue, clKernel1, 1, NULL,
                                     &scalarGlobalWSize, &localWorkSize, 0,
                                     NULL, &event);

    if (errcode != CL_SUCCESS)
        printf("Error in launching kernel %d\n", errcode);
    clFinish(clCommandQue);
    clGetEventProfilingInfo(event, CL_PROFILING_COMMAND_START, sizeof(cl_ulong),
                            &t_start, NULL);
    clGetEventProfilingInfo(event, CL_PROFILING_COMMAND_END, sizeof(cl_ulong),
                            &t_end, NULL);
    duration = t_end - t_start;

    errcode =
        clEnqueueReadBuffer(clCommandQue, d_out, 1, 0,
                            numRows * sizeof(floatType), h_out, 0, NULL, NULL);

    if (errcode != CL_SUCCESS)
        printf("Error in transfering back results %d \n", errcode);

    int vfy;

    vfy = verifyResults(refOut, h_out, numRows);

    if (vfy == 1)
        printf("Passed \n");

    printf("GPU Runtime: %.10f\n", (double)duration / 1e9);

    cl_clean_up();

    return 0;
}
