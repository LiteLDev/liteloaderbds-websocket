#pragma once

#include "GoTypes.h"

#ifdef __cplusplus
#define EXTERN_C_GOFunc extern "C"
#else
#define EXTERN_C_GOFunc extern
#endif

EXTERN_C_GOFunc void Init(void(*)(char*, size_t, char));
EXTERN_C_GOFunc void StartServer();
EXTERN_C_GOFunc void SetRuncmdFunc(struct GoString* (*)(char*, long long));
EXTERN_C_GOFunc void SetMSVCFreeFunc(void (*)(void*));

#undef EXTERN_C_GOFunc