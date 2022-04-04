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

namespace SubmoduleFuncBinder {

	extern "C" void* (*getProcAddrWrapper)(void* hModule, LPCSTR lpProcName);
	extern "C" void* subModuleHandler;


	extern "C" void* (*getProcAddrWrapper)(void* hModule, LPCSTR lpProcName);

	std::map<std::string, __int64> ModuleList;
	std::map<__int64, std::map<std::string, __int64>> FakeDllExports;

	HCUSTOMMODULE CustomLoadLibrary(LPCSTR filename, void* userdata){
		//logger.warn("LoadLibrary {}", filename);
		auto iter = ModuleList.find(std::string(filename));
		if (iter != ModuleList.end()) {
			//logger.warn("CustomBind Module {} -> {}", filename, it->second);
			return (void*)iter->second;
		}
		HMODULE result;
		UNREFERENCED_PARAMETER(userdata);
		result = LoadLibraryA(filename);
		if (result == NULL) {
			return NULL;
		}

		return (HCUSTOMMODULE)result;
	}

	FARPROC CustomGetProcAddress(HCUSTOMMODULE module, LPCSTR name, void* userdata){
		//logger.warn("GetProcAddress {} {}", (__int64)module, name);
		auto iter = FakeDllExports.find((__int64)module);
		if (iter != FakeDllExports.end()) {
			//logger.warn("CustomBind FakeModule {}", (__int64)module);
			auto const& mExports = iter->second;
			auto expIter = mExports.find(std::string(name));
			if (expIter != mExports.end()) {
				//logger.warn("CustomBind FakeFunc {} -> {}", std::string(name), (long long)(expIter->second));
				return (FARPROC)expIter->second;
			}
			else {
				logger.warn("CustomBind FakeFunc {} -> NULL", std::string(name));
				return NULL;
			}
		}

		UNREFERENCED_PARAMETER(userdata);
		return (FARPROC)GetProcAddress((HMODULE)module, name);
	}
	extern "C" void LoadMemDLL() {
		if (std::filesystem::exists(L"llws_golang.dll")) {
			std::ifstream dllFile("llws_golang.dll", std::ifstream::binary);
			if (dllFile) {
				dllFile.seekg(0, dllFile.end);
				size_t length = dllFile.tellg();
				dllFile.seekg(0, dllFile.beg);

				char* buffer = new char[length];
				dllFile.read(buffer, length);

				if (!dllFile) {
					logger.error("[Debug] Failed to read llws_golang.dll from file system");
					dllFile.close();
				}
				dllFile.close();

				subModuleHandler = MemoryLoadLibraryEx(buffer, length, MemoryDefaultAlloc, MemoryDefaultFree, CustomLoadLibrary, CustomGetProcAddress, MemoryDefaultFreeLibrary, NULL);
				getProcAddrWrapper = (void* (*)(void*, LPCSTR))MemoryGetProcAddress;
				logger.info("Loaded GolangModule From Filesystem {} {}", (__int64)subModuleHandler, (__int64)getProcAddrWrapper);
			}
			else {
				logger.error("[Debug] Failed to load llws_golang.dll from file system");
				dllFile.close();
			}
		}
		else {
			HRSRC thisRes = ::FindResource(dllModule, MAKEINTRESOURCE(Websocket), L"DLL");
			unsigned int resSize = ::SizeofResource(dllModule, thisRes);
			void* resData = ::LockResource(::LoadResource(dllModule, thisRes));
			subModuleHandler = MemoryLoadLibraryEx(resData, resSize, MemoryDefaultAlloc, MemoryDefaultFree, CustomLoadLibrary, CustomGetProcAddress, MemoryDefaultFreeLibrary, NULL);
			int err = GetLastError();
			std::cout << err << std::endl;
			getProcAddrWrapper = (void* (*)(void*, LPCSTR))MemoryGetProcAddress;
			logger.info("Loaded GolangModule From Resoueces {} {}", (__int64)subModuleHandler, (__int64)getProcAddrWrapper);
		}
	}
}

namespace GolangWrappers {
	void LoggerWrapper(char* str, size_t size, char logLevel) {
		std::string cpplog(str, size);
		*subLogger[logLevel] << cpplog << logger.endl;
	}

	struct GoString* RuncmdWrapper(char* cmd, size_t len) {
		std::string cppcmd(cmd, len);
		auto result = Global<Level>->runcmdEx(cppcmd);
		return new GoString(result.second);
	}

	void BroadcastMessageWrapper(char* msg, size_t len, int textType) {
		auto players = Global<Level>->getAllPlayers();
		auto msgStr = std::string(msg, len);
		for (auto& player : players) {
			player->sendText(msgStr, (TextType)textType);
		}
	}		
	
	void MSVCFreeWrapper(void* p) {
		free(p);
	}

}

void PluginInit(){
	SubmoduleFuncBinder::ModuleList.insert({ "llws_cpp_host.dll", (__int64)0x1001 });
	auto& FakeDll_1001 = SubmoduleFuncBinder::FakeDllExports[0x1001];
	FakeDll_1001.insert({ "LoggerWrapper", (__int64)GolangWrappers::LoggerWrapper });
	FakeDll_1001.insert({ "RuncmdWrapper",(__int64)GolangWrappers::RuncmdWrapper });
	FakeDll_1001.insert({ "MSVCFreeWrapper",(__int64)GolangWrappers::MSVCFreeWrapper });
	FakeDll_1001.insert({ "BroadcastMessageWrapper",(__int64)GolangWrappers::BroadcastMessageWrapper });
	SubmoduleFuncBinder::LoadMemDLL();
	Init();

	Event::ServerStartedEvent::subscribe([](Event::ServerStartedEvent ev)->bool {
		StartServer();
		return true;
		}
	);
	Event::PlayerChatEvent::subscribe([](Event::PlayerChatEvent ev)->bool {
		ChatEventBroadcast(ev.mPlayer->getRealName(), ev.mMessage);
		return true;
		}
	);
	Event::PlayerJoinEvent::subscribe([](Event::PlayerJoinEvent ev)->bool {
		auto& mPlayer = ev.mPlayer;
		auto& mPos = mPlayer->getPos();
		GoSlice<float> goPos(vector<float>({ mPos.x, mPos.y, mPos.z }));
		JoinEventBroadcast(mPlayer->getRealName(), mPlayer->getXuid(), mPlayer->getUuid(), mPlayer->getIP(), goPos, mPlayer->getDimensionId());
		return true;
		}
	);
	Event::PlayerLeftEvent::subscribe([](Event::PlayerLeftEvent ev)->bool {
		auto& mPlayer = ev.mPlayer;
		auto& mPos = mPlayer->getPos();
		GoSlice<float> goPos(vector<float>({ mPos.x, mPos.y, mPos.z }));
		LeftEventBroadcast(mPlayer->getRealName(), mPlayer->getXuid(), mPlayer->getUuid(), goPos, mPlayer->getDimensionId());
		return true;
		}
	);
}
