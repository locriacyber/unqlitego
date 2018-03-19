package unqlitego

// #include <unqlite.h>
// #include <wrappers.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"sync"
)

var (
	lib *Library
)

// init will initialize the library.
// This is to replace to previous implementation of init for the UnQLite library.
// In order to be fully ThreadSafe the library itself requires a global Mutex.
//
// Calls to 'C.unqlite_lib_init()' are ThreadSafe and any subsequents call will result in NOOP.
// Calls to 'C.unqlite_lib_shutdown()' are NOT ThreadSafe and require a library Mutex.
// Any subsequent calls to 'C.unqlite_lib_shutdown()' will result in NOOP.
// Calls to 'unqlite_lib_config' are not ThreadSafe and require a library Mutex.
func init() {
	lib = &Library{
		mu:  new(sync.Mutex),
		lck: new(sync.Mutex),
	}

	if !lib.IsThreadSafe() {
		panic("UnQLite was not compiled ThreadSafe, please include 'UNQLITE_ENABLE_THREADS' directive in compilation.")
	}
}

// Library represents the UnQLite Library Control.
type Library struct {
	// ThreadSafe Mutex
	mu *sync.Mutex

	// Library Lock
	lck *sync.Mutex
}

// Info returns the UnQLite Library.
func Info() *Library {
	return lib
}

// Init will initialize the library. After call init, the library must be shutdown first before any
// re-configuration can be done. Subsequent calls result in NOOP.
// This is autocalled when a database is opened.
func (l *Library) Init() {
	// Initalize ThreadSafe Library Init.
	l.mu.Lock()
	defer l.mu.Unlock()

	// Initalize Library Lock
	// This lock will remain until Shutdown() is called.
	// Library lock is required for calls to C.unqlite_lib_config.
	// C.unqlite_lib_config is not ThreadSafe and is only allowed to be called
	// when no call has yet been made to C.unqlite_lib_init().
	// Therefor set global library lock, and release on call to Shutdown().
	defer l.lck.Lock()

	C.unqlite_lib_init()
}

// IsThreadSafe returns a boolean identifiying if the UnQLite library is
// compiled Thread Safe.
func (l *Library) IsThreadSafe() bool {
	return C.unqlite_lib_is_threadsafe() == 1
}

// Version returns the version string.
func (l *Library) Version() string {
	return C.GoString(C.unqlite_lib_version())
}

// Signature returns the Signature string.
func (l *Library) Signature() string {
	return C.GoString(C.unqlite_lib_signature())
}

// Ident returns the UnQLite identification string.
func (l *Library) Ident() string {
	return C.GoString(C.unqlite_lib_ident())
}

// Copyright returns the Copyright string.
func (l *Library) Copyright() string {
	return C.GoString(C.unqlite_lib_copyright())
}

// Shutdown the UnQLite Library.
func (l *Library) Shutdown() error {
	var err error

	// Enforce ThreadSafe operation.
	l.mu.Lock()
	defer l.mu.Unlock()

	// Release Library lock
	defer l.lck.Unlock()

	res := C.unqlite_lib_shutdown()
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return err
}

// UnQLiteError is returned on both UnQLite native errors as well as UnQLite Go errors.
// Native errors are within the range of '<= 0' (Equal and lesser then Zero) while the
// errors from UnQLite Go are > 0 (Greater then Zero).
//
// 	( UnQLiteError <= 0 ) Native Errors
//	( UnQLiteError >  0 ) UnQLite Go Errors.
type UnQLiteError int

// Error returns the string representation of the UnQLiteError.
func (e UnQLiteError) Error() string {
	s := errString[e]
	if s == "" {
		return fmt.Sprintf("err: %d", int(e))
	}

	return s
}

func (e UnQLiteError) String() string {
	return e.Error()
}

var errString = map[UnQLiteError]string{
	C.UNQLITE_OK:             "OK",
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

// unQLiteValue represents an UnQLite Value object which is returned from the unlying UnQLite Database.
type unQLiteValue struct {
	// Pointer to underlying UnQLite Value.
	v *C.unqlite_value
}

// Nil returns a boolean indicating if the unQLiteValue is nil.
func (uv *unQLiteValue) Nil() bool {
	switch uv.v {
	case nil:
		// User Data is wrong
		return true
	default:
		return false
	}
}
