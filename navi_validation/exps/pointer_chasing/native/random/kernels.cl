__kernel void microbench(__global ulong *memory, __global ulong* offset, int count){
    int global_id = get_global_id(0);

    __global ulong *p = memory;
    p += offset[global_id];

    ulong q;
    //printf("%p ",memory);
    for(int i=0; i<count; i++){
        //printf("%p ",p); 
        q = (ulong)memory + *p; 
        p = (__global ulong*)q;
    }
    *p=1;
}
