__kernel void microbench(__global const int *in ,int count) {
    
     int a;
     for(int i = 0; i < count; i++)
    {	     
       __asm__ volatile(
               "flat_load_dword %0, %1\n\t" 
               "s_waitcnt vmcnt(0)" 
               : "=v"(a) 
               : "v"(in));
    }
}
