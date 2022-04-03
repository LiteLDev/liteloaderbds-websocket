#include <stdlib.h>
#include <stdint.h>

struct cgoString {
	const char* Data;
	unsigned long long Length;
};

extern void LoggerWrapper(char* str, size_t size,char logLevel);

extern struct cgoString* RuncmdWrapper(char* cmd, size_t len);

extern void MSVCFreeWrapper(void* p);

extern void BroadcastMessageWrapper(char* msg, size_t len, int textType);
