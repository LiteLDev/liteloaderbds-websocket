#pragma once

#include "GoTypes.h"

#ifdef __cplusplus
#define EXTERN_C_GOFunc extern "C"
#else
#define EXTERN_C_GOFunc extern
#endif

EXTERN_C_GOFunc void Init();
EXTERN_C_GOFunc void StartServer();

#undef EXTERN_C_GOFunc