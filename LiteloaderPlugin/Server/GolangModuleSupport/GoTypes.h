#pragma once

#include <stddef.h> /* for ptrdiff_t below */
#include <vector>
#include <string>
/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
#ifndef _WIN64
static_assert(false,"This Header is only working in 64bit Windows")
#endif

static_assert (sizeof(void*) == 64 / 8, "check for 64 bit pointer matching GoInt");

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef size_t GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;

typedef void* GoMap;
typedef void* GoChan;
typedef struct { void* t; void* v; } GoInterface;

//typedef struct { const char* p; ptrdiff_t n; } _GoString_;
struct GoString {
	const char* Data;
	size_t Length;
	// std::string Conventor
	operator std::string() const {
		return std::string(Data, Length);
	}
	std::string ToString() const {
		return std::string(Data, Length);
	}

	// from std::string
	GoString(const std::string& str) {
		Data = str.c_str();
		Length = str.length();
	}
	
	//destructor
	~GoString() {
		Data = nullptr;
		Length = 0;
	}
};

//typedef struct { void* data; GoInt len; GoInt cap; } GoSlice;
template <typename T>
struct GoSlice {
	T* Data;
	GoInt Length;
	GoInt Capacity;
	// std::vector Conventor
	operator std::vector<T>() const {
		return std::vector<T>(Data, Data + Length);
	}
	std::vector<T> ToVector() const {
		return std::vector<T>(Data, Data + Length);
	}

	// from std::vector
	GoSlice(const std::vector<T>& vec) {
		Data = vec.data();
		Length = vec.size();
		Capacity = vec.capacity();
	}
	
	//destructor
	~GoSlice() {
		Data = nullptr;
		Length = 0;
		Capacity = 0;
	}
};
