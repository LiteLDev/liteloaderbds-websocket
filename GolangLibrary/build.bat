dlltool -l libfakedll.a -d .\server\c_wrappers.def

go build -buildmode=c-shared -trimpath -o llws_golang.dll
go build -trimpath -o llws_client.exe