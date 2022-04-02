#include "pch.h"
#include <EventAPI.h>
#include <LoggerAPI.h>
#include <MC/Level.hpp>
#include <MC/BlockInstance.hpp>
#include <MC/Block.hpp>
#include <MC/BlockSource.hpp>
#include <MC/Actor.hpp>
#include <MC/Player.hpp>
#include <MC/ItemStack.hpp>
#include <LLAPI.h>

#include "GolangModuleSupport/MemoryModule.h"

#include "GolangModuleSupport/GoFuncDef.h"

#include "resource.h"

// LiteLoader Style Logger
Logger logger("LLWS");
// subLogger ref to logger splited by different level
Logger::OutputStream* subLogger[] = { &logger.info,&logger.warn,&logger.debug,&logger.error };

HMODULE dllModule;

extern "C" void* (*getProcAddrWrapper)(void* hModule, LPCSTR lpProcName);
extern "C" void* subModuleHandler;


extern "C" void* (*getProcAddrWrapper)(void* hModule, LPCSTR lpProcName);

std::mutex DllLock;
extern "C" void LoadMemDLL() {
	if (std::filesystem::exists(L"llws_golang.dll")) {
		subModuleHandler = LoadLibrary(L"llws_golang.dll");
		getProcAddrWrapper = (void* (*)(void*, LPCSTR))GetProcAddress;
		logger.info("Loaded GolangModule From Filesystem {} {}", (__int64)subModuleHandler, (__int64)getProcAddrWrapper);
	}
	else {

		HRSRC thisRes = ::FindResource(dllModule, MAKEINTRESOURCE(Websocket), L"DLL");
		unsigned int ResSize = ::SizeofResource(dllModule, thisRes);
		HGLOBAL ResData = ::LoadResource(dllModule, thisRes);
		void* ResDataRef = ::LockResource(ResData);
		subModuleHandler = MemoryLoadLibrary(ResDataRef, ResSize);
		getProcAddrWrapper = (void* (*)(void*, LPCSTR))MemoryGetProcAddress;
		logger.info("Loaded GolangModule From Resoueces {} {}", (__int64)subModuleHandler, (__int64)getProcAddrWrapper);
	}
}

void GoLoggerWrapper(char* str, size_t size,char logLevel) {
	std::string cpplog(str, size);
	*subLogger[logLevel] << cpplog << logger.endl;
}

struct GoString* GoRuncmdWrapper(char* cmd, long long len) {
	std::string cppcmd(cmd, len);
	auto result = Global<Level>->runcmdEx(cppcmd);
	return new GoString(result.second);
}

void GoMSVCFreeWrapper(void* p) {
	free(p);
}

void PluginInit(){
	LoadMemDLL();
	Init(GoLoggerWrapper);
	SetRuncmdFunc(GoRuncmdWrapper);
	SetMSVCFreeFunc(GoMSVCFreeWrapper);

	Event::ServerStartedEvent::subscribe([](Event::ServerStartedEvent ev)->bool {
		StartServer();
		return true;
	});
	
}
