// dllmain.cpp : 定义 DLL 应用程序的入口点。
#include "pch.h"
#include <LLAPI.h>
#pragma comment(lib, "../SDK/Lib/bedrock_server_api.lib")
#pragma comment(lib, "../SDK/Lib/bedrock_server_var.lib")
#pragma comment(lib, "../SDK/Lib/LiteLoader.lib")
#pragma comment(lib, "GolangModuleSupport/GoDll.lib")
// using Builtin SymDBHelper instead of LL provided one
// for using convenience GolangDLL Support
//#pragma comment(lib, "../SDK/Lib/SymDBHelper.lib")
extern HMODULE dllModule;
BOOL APIENTRY DllMain( HMODULE hModule,
                       DWORD  ul_reason_for_call,
                       LPVOID lpReserved
                     )
{
    dllModule = hModule;
    switch (ul_reason_for_call)
    {
    case DLL_PROCESS_ATTACH:
        LL::registerPlugin("Websocket", "using websocket to manage this server", LL::Version(2, 0, 0), {
                { "Author", "LiteLDev" },
                { "Github", "https://github.com/LiteLDev/Websocket" }
            }
        );
        break;
    case DLL_THREAD_ATTACH:
    case DLL_THREAD_DETACH:
    case DLL_PROCESS_DETACH:
        break;
    }
    return TRUE;
}

void PluginInit();

extern "C" {
    // Do something after all the plugins loaded
    _declspec(dllexport) void onPostInit() {
        std::ios::sync_with_stdio(false);
        PluginInit();
    }
}
