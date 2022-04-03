package server

/*
#include "c_wrappers.h"
#cgo LDFLAGS: -Wl,--allow-multiple-definition -L.. -lfakedll
*/
import "C"
import (
	"BDSWebsocket/logger"
	"unsafe"
)

//export Init
// Init is a wrapper that init Golang Logger
func Init() {
	logger.SetOutput(logger.Writer_wrapper{
		WriteFunc: func(p []byte) {
			str := C.CString(string(p))
			C.LoggerWrapper(str, C.ulonglong(len(p)), logger.LInfo)
			C.free(unsafe.Pointer(str))
		},
	})
	logger.Warn.SetOutput(logger.Writer_wrapper{
		WriteFunc: func(p []byte) {
			str := C.CString(string(p))
			C.LoggerWrapper(str, C.ulonglong(len(p)), logger.LWarn)
			C.free(unsafe.Pointer(str))
		},
	})
	logger.Debug.SetOutput(logger.Writer_wrapper{
		WriteFunc: func(p []byte) {
			str := C.CString(string(p))
			C.LoggerWrapper(str, C.ulonglong(len(p)), logger.LDebug)
			C.free(unsafe.Pointer(str))
		},
	})
	logger.Error.SetOutput(logger.Writer_wrapper{
		WriteFunc: func(p []byte) {
			str := C.CString(string(p))
			C.LoggerWrapper(str, C.ulonglong(len(p)), logger.LError)
			C.free(unsafe.Pointer(str))
		},
	})
}

// CallRuncmdFunc is a wrapper that call the callback function that execute the Minecraft command
func CallRuncmdFunc(cmd string) string {
	str := C.CString(cmd)
	result := C.RuncmdWrapper(str, C.ulonglong(len(cmd)))
	C.free(unsafe.Pointer(str))
	resultStr := C.GoBytes(unsafe.Pointer(result.Data), C.int(result.Length))
	CallMSVCFreeFunc(unsafe.Pointer(result))
	return string(resultStr)
}

// CallMSVCFreeFunc is a wrapper that call the callback function that free the memory allocated by the MSVC
func CallMSVCFreeFunc(block unsafe.Pointer) {
	C.MSVCFreeWrapper(block)
}

// CallBroadcastMessageWrapper is a wrapper that call the callback function that broadcast the message to all Players
func CallBroadcastMessageWrapper(message string, textType int) {
	str := C.CString(message)
	C.BroadcastMessageWrapper(str, C.ulonglong(len(message)), C.int(textType))
	C.free(unsafe.Pointer(str))
}
