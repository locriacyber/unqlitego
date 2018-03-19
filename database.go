package unqlitego

// #include <unqlite.h>
// #include <wrappers.h>
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// Database ...
type Database struct {
	handle *C.unqlite
}

// NewDatabase ...
func NewDatabase(filename string) (db *Database, err error) {
	// TODO: Enforce call to check library lock and call forced call to init.
	db = &Database{}
	name := C.CString(filename)
	defer C.free(unsafe.Pointer(name))
	res := C.unqlite_open(&db.handle, name, C.UNQLITE_OPEN_CREATE)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	if db.handle != nil {
		runtime.SetFinalizer(db, (*Database).Close)
	}
	return
}

func (db *Database) Unqlite_compile(jx9_script string, vm *VM) (error, string) {
	res := C.unqlite_compile(db.handle, C.CString(jx9_script), C.int(len(jx9_script)), &vm.vm)
	if res != C.UNQLITE_OK {
		if res == C.UNQLITE_COMPILE_ERR {
			err := UnQLiteError(res)
			error_log := new(C.char)
			err_msg := C.extract_unqlite_log_error(db.handle, error_log)
			g_err_msg := C.GoString(err_msg)
			//C.free(unsafe.Pointer(err_msg))
			return err, g_err_msg
		}
	}
	return nil, ""
}

// Close ...
func (db *Database) Close() (err error) {
	if db.handle != nil {
		res := C.unqlite_close(db.handle)
		if res != C.UNQLITE_OK {
			err = UnQLiteError(res)
		}
		db.handle = nil
	}
	return
}

// Store ...
func (db *Database) Store(key, value []byte) (err error) {
	var k, v unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	if len(value) > 0 {
		v = unsafe.Pointer(&value[0])
	}

	res := C.unqlite_kv_store(db.handle,
		k, C.int(len(key)),
		v, C.unqlite_int64(len(value)))
	if res == C.UNQLITE_OK {
		return nil
	}
	return UnQLiteError(res)
}

// Append ...
func (db *Database) Append(key, value []byte) (err error) {
	var k, v unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	if len(value) > 0 {
		v = unsafe.Pointer(&value[0])
	}

	res := C.unqlite_kv_append(db.handle,
		k, C.int(len(key)),
		v, C.unqlite_int64(len(value)))
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Fetch ...
func (db *Database) Fetch(key []byte) (value []byte, err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	var n C.unqlite_int64
	res := C.unqlite_kv_fetch(db.handle, k, C.int(len(key)), nil, &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
		return
	}
	value = make([]byte, int(n))
	res = C.unqlite_kv_fetch(db.handle, k, C.int(len(key)), unsafe.Pointer(&value[0]), &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Delete ...
func (db *Database) Delete(key []byte) (err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_delete(db.handle, k, C.int(len(key)))
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Begin ...
func (db *Database) Begin() (err error) {
	res := C.unqlite_begin(db.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Commit ...
func (db *Database) Commit() (err error) {
	res := C.unqlite_commit(db.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Rollback ...
func (db *Database) Rollback() (err error) {
	res := C.unqlite_rollback(db.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}
