const char *kernels_ocl =
"__kernel void microbench(__global const int *in ,int count) {\n"
"    \n"
"     int a;\n"
"     for(int i = 0; i < count; i++)\n"
"    {	     \n"
"       __asm__ volatile(\n"
"               \"flat_load_dword %0, %1\\n\\t\" \n"
"               \"s_waitcnt vmcnt(0)\" \n"
"               : \"=v\"(a) \n"
"               : \"v\"(in));\n"
"    }\n"
"}\n"
;
