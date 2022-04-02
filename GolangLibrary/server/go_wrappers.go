package server

/*
#include "c_wrappers.h"
#cgo LDFLAGS: -Wl,--allow-multiple-definition
*/
import "C"
import (
	"BDSWebsocket/logger"
	"unsafe"
)

//export Init
// Init is a wrapper that init Golang Logger
func Init(loggerFuncPtr unsafe.Pointer) {
	C.set_logger(loggerFuncPtr)
	logger.SetOutput(logger.Writer_wrapper{
		WriteFunc: func(p []byte) {
			str := C.CString(string(p))
			C.call_logger(str, C.longlong(len(p)), logger.LInfo)
			C.free(unsafe.Pointer(str))
		},
	})
	logger.Warn.SetOutput(logger.Writer_wrapper{
		WriteFunc: func(p []byte) {
			str := C.CString(string(p))
			C.call_logger(str, C.longlong(len(p)), logger.LWarn)
			C.free(unsafe.Pointer(str))
		},
	})
	logger.Debug.SetOutput(logger.Writer_wrapper{
		WriteFunc: func(p []byte) {
			str := C.CString(string(p))
			C.call_logger(str, C.longlong(len(p)), logger.LDebug)
			C.free(unsafe.Pointer(str))
		},
	})
	logger.Error.SetOutput(logger.Writer_wrapper{
		WriteFunc: func(p []byte) {
			str := C.CString(string(p))
			C.call_logger(str, C.longlong(len(p)), logger.LError)
			C.free(unsafe.Pointer(str))
		},
	})
}

//export SetRuncmdFunc
// SetRuncmdFunc is a wrapper that set the callback function for runcmd
func SetRuncmdFunc(runcmdFuncPtr unsafe.Pointer) {
	C.set_runcmd(runcmdFuncPtr)
}

// CallRuncmdFunc is a wrapper that call the callback function that execute the Minecraft command
func CallRuncmdFunc(cmd string) string {
	str := C.CString(cmd)
	result := C.call_runcmd(str, C.longlong(len(cmd)))
	C.free(unsafe.Pointer(str))
	resultStr := C.GoBytes(unsafe.Pointer(result.Data), C.int(result.Length))
	CallMSVCFreeFunc(unsafe.Pointer(result))
	return string(resultStr)
}

//export SetMSVCFreeFunc
// SetMSVCFreeFunc is a wrapper that set the callback function for NSVVC free function
// we need it to free the memory allocated by the MSVC due to incompatible Memory Allocation
func SetMSVCFreeFunc(runcmdFuncPtr unsafe.Pointer) {
	C.set_msvc_free(runcmdFuncPtr)
}

// CallMSVCFreeFunc is a wrapper that call the callback function that free the memory allocated by the MSVC
func CallMSVCFreeFunc(block unsafe.Pointer) {
	C.call_msvc_free(block)
}
