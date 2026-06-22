/* Compiler options:
-c -emit-llvm -target amdgcn-amd-amdhsa -x cl -O3 -cl-kernel-arg-info -D__OPENCL_VERSION__=200 -D__IMAGE_SUPPORT__=1 -Xclang -cl-ext=+cl_khr_fp64,+cl_khr_global_int32_base_atomics,+cl_khr_global_int32_extended_atomics,+cl_khr_local_int32_base_atomics,+cl_khr_local_int32_extended_atomics,+cl_khr_int64_base_atomics,+cl_khr_int64_extended_atomics,+cl_khr_3d_image_writes,+cl_khr_byte_addressable_store,+cl_khr_fp16,+cl_khr_gl_sharing,+cl_khr_gl_depth_images,+cl_amd_device_attribute_query,+cl_amd_media_ops,+cl_amd_media_ops2,+cl_khr_image2d_from_buffer,+cl_khr_subgroups,+cl_khr_gl_event,+cl_khr_mipmap_image,+cl_khr_mipmap_image_writes,+cl_amd_copy_buffer_p2p -mllvm -amdgpu-prelink  -include opencl-c.h 
*/

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
