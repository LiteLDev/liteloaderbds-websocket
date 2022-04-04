#pragma once

#include "GoTypes.h"

#ifdef __cplusplus
#define EXTERN_C_GOFunc extern "C"
#else
#define EXTERN_C_GOFunc extern
#endif

EXTERN_C_GOFunc void Init();
EXTERN_C_GOFunc void StartServer();
EXTERN_C_GOFunc void ChatEventBroadcast(GoString playerName, GoString message);
EXTERN_C_GOFunc void JoinEventBroadcast(GoString playerName, GoString XUID, GoString UUID, GoString ipAddress, GoSlice<float> position, GoInt dimensionId);
EXTERN_C_GOFunc void LeftEventBroadcast(GoString playerName, GoString XUID, GoString UUID, GoSlice<float> position, GoInt dimensionId);

#undef EXTERN_C_GOFunc