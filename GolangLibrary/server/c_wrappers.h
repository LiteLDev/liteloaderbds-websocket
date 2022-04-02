#include <stdlib.h>

struct cgoString {
	const char* Data;
	unsigned long long Length;
};

void* logger_func_ptr = 0;
void call_logger(char* buf,long long size,char logLevel){
	((void(*)(char*,long long,char))logger_func_ptr)(buf,size,logLevel);
	return;
}
void set_logger(void* fn){
	logger_func_ptr = fn;
	return;
}

void* msvc_free_func_ptr = 0;
void call_msvc_free(void* block){
    return ((void(*)(void*))msvc_free_func_ptr)(block);
}
void set_msvc_free(void* fn){
    msvc_free_func_ptr = fn;
    return;
}


void* runcmd_func_ptr = 0;
struct cgoString* call_runcmd(char* cmd, long long size){
    return ((struct cgoString*(*)(char*,long long))runcmd_func_ptr)(cmd,size);
}
void set_runcmd(void* fn){
    runcmd_func_ptr = fn;
    return;
}
