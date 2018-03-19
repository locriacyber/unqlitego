package unqlitego

// #include <unqlite.h>
// #include <wrappers.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
)

// UnQLiteError ... standard error for this module

type GlobaLError string

func (s GlobaLError) Error() string {
	return string(s)
}

type UnQLiteError int

func (e UnQLiteError) Error() string {
	s := errString[e]
	if s == "" {
		return fmt.Sprintf("errno %d", int(e))
	}
	return s
}

var errString = map[UnQLiteError]string{
	C.UNQLITE_LOCKERR:        "Locking protocol error",
	C.UNQLITE_READ_ONLY:      "Read only Key/Value storage engine",
	C.UNQLITE_CANTOPEN:       "Unable to open the database file",
	C.UNQLITE_FULL:           "Full database",
	C.UNQLITE_VM_ERR:         "Virtual machine error",
	C.UNQLITE_COMPILE_ERR:    "Compilation error",
	C.UNQLITE_DONE:           "Operation done", // Not an error.
	C.UNQLITE_CORRUPT:        "Corrupt pointer",
	C.UNQLITE_NOOP:           "No such method",
	C.UNQLITE_PERM:           "Permission error",
	C.UNQLITE_EOF:            "End Of Input",
	C.UNQLITE_NOTIMPLEMENTED: "Method not implemented by the underlying Key/Value storage engine",
	C.UNQLITE_BUSY:           "The database file is locked",
	C.UNQLITE_UNKNOWN:        "Unknown configuration option",
	C.UNQLITE_EXISTS:         "Record exists",
	C.UNQLITE_ABORT:          "Another thread have released this instance",
	C.UNQLITE_INVALID:        "Invalid parameter",
	C.UNQLITE_LIMIT:          "Database limit reached",
	C.UNQLITE_NOTFOUND:       "No such record",
	C.UNQLITE_LOCKED:         "Forbidden Operation",
	C.UNQLITE_EMPTY:          "Empty record",
	C.UNQLITE_IOERR:          "IO error",
	C.UNQLITE_NOMEM:          "Out of memory",
}

type Unqlite_value struct {
	unqlite_value *C.unqlite_value
}

func init() {
	C.unqlite_lib_init()
	if !IsThreadSafe() {
		panic("unqlite library was not compiled for thread-safe option UNQLITE_ENABLE_THREADS=1")
	}
}

func unqlite_value_ok(unqlite_value *Unqlite_value) bool {
	switch unqlite_value.unqlite_value {
	case nil:
		return false //User Data is wrong
	default:
		return true
	}
}

// Shutdown the UnQlite Library.
func Shutdown() (err error) {
	res := C.unqlite_lib_shutdown()
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

/* TODO: implement

// Database Engine Handle
int unqlite_config(unqlite *pDb,int nOp,...);

// Key/Value (KV) Store Interfaces
int unqlite_kv_fetch_callback(unqlite *pDb,const void *pKey,
	                    int nKeyLen,int (*xConsumer)(const void *,unsigned int,void *),void *pUserData);
int unqlite_kv_config(unqlite *pDb,int iOp,...);

//  Cursor Iterator Interfaces
int unqlite_kv_cursor_key_callback(unqlite_kv_cursor *pCursor,int (*xConsumer)(const void *,unsigned int,void *),void *pUserData);
int unqlite_kv_cursor_data_callback(unqlite_kv_cursor *pCursor,int (*xConsumer)(const void *,unsigned int,void *),void *pUserData);

// Utility interfaces
int unqlite_util_load_mmaped_file(const char *zFile,void **ppMap,unqlite_int64 *pFileSize);
int unqlite_util_release_mmaped_file(void *pMap,unqlite_int64 iFileSize);

// Global Library Management Interfaces
int unqlite_lib_config(int nConfigOp,...);
*/
