package unqlitego

// #include <unqlite.h>
// #include <wrappers.h>
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// Database represents a UnQLite Database.
type Database struct {
	conn *C.unqlite
}

// NewDatabase creates and initalizes a new UnQLite database connection.
func NewDatabase(filename string) (db *Database, err error) {
	// TODO: Enforce call to check library lock and call forced call to init.
	db = &Database{}

	name := C.CString(filename)
	defer C.free(unsafe.Pointer(name))

	// Call Library Init
	// This will lock the library against modifications.
	Info().Init()

	res := C.unqlite_open(&db.conn, name, C.UNQLITE_OPEN_CREATE)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	if db.conn != nil {
		runtime.SetFinalizer(db, (*Database).Close)
	}

	return
}

// Compile a JX9 Script into a Virtual Machine.
func (db *Database) Compile(jx9 string, vm *VM) (string, error) {

	res := C.unqlite_compile(db.conn, C.CString(jx9), C.int(len(jx9)), &vm.vm)
	if res != C.UNQLITE_OK {
		if res == C.UNQLITE_COMPILE_ERR {
			err := UnQLiteError(res)

			// Error Log
			elog := new(C.char)

			// Error Message
			emsg := C.extract_unqlite_log_error(db.conn, elog)

			// Global Error Message
			gmsg := C.GoString(emsg)

			return gmsg, err
		}
	}

	return "", nil
}

// Close will close the database connection.
func (db *Database) Close() (err error) {
	if db.conn != nil {
		res := C.unqlite_close(db.conn)
		if res != C.UNQLITE_OK {
			err = UnQLiteError(res)
		}
		db.conn = nil
	}

	return
}

// Store will store a new Key/Value pair in the database.
func (db *Database) Store(key, value []byte) (err error) {
	var k, v unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	if len(value) > 0 {
		v = unsafe.Pointer(&value[0])
	}

	res := C.unqlite_kv_store(db.conn,
		k, C.int(len(key)),
		v, C.unqlite_int64(len(value)))
	if res == C.UNQLITE_OK {
		return nil
	}

	return UnQLiteError(res)
}

// Append will write a new record into the database.
// If the record does not exists it will be created, else the new
// data is appended to the end of the old data.
func (db *Database) Append(key, value []byte) (err error) {
	var k, v unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	if len(value) > 0 {
		v = unsafe.Pointer(&value[0])
	}

	res := C.unqlite_kv_append(db.conn,
		k, C.int(len(key)),
		v, C.unqlite_int64(len(value)))

	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Fetch a record from the database.
func (db *Database) Fetch(key []byte) (value []byte, err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	var n C.unqlite_int64
	res := C.unqlite_kv_fetch(db.conn, k, C.int(len(key)), nil, &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
		return
	}

	value = make([]byte, int(n))
	res = C.unqlite_kv_fetch(db.conn, k, C.int(len(key)), unsafe.Pointer(&value[0]), &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Delete a record from the database.
func (db *Database) Delete(key []byte) (err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_delete(db.conn, k, C.int(len(key)))
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Begin will start a transaction.
func (db *Database) Begin() (err error) {
	res := C.unqlite_begin(db.conn)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Commit will commit a transaction to the database.
func (db *Database) Commit() (err error) {
	res := C.unqlite_commit(db.conn)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Rollback a transaction.
func (db *Database) Rollback() (err error) {
	res := C.unqlite_rollback(db.conn)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}
