package unqlitego

// #include <unqlite.h>
// #include <wrappers.h>
// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"
)

// Cursor represents an UnQLite database cursor.
type Cursor struct {
	parent *Database
	handle *C.unqlite_kv_cursor
}

// Cursor creates and initializes a new UnQLite database cursor.
func (db *Database) Cursor() (*Cursor, error) {
	c := &Cursor{parent: db}
	res := C.unqlite_kv_cursor_init(db.handle, &cursor.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	runtime.SetFinalizer(cursor, (*Cursor).Close)

	return
}

// Close closes the open cursor.
func (curs *Cursor) Close() (err error) {
	if curs.parent.handle != nil && curs.handle != nil {
		res := C.unqlite_kv_cursor_release(curs.parent.handle, curs.handle)
		if res != C.UNQLITE_OK {
			err = UnQLiteError(res)
		}
		curs.handle = nil
	}

	return
}

// Seek will search the cursor for an exact match. If the record exists the cursor is left pointing to it.
// Otherwise it is left pointing to EOF and UNQLITE_NOTFOUND is returned.
func (curs *Cursor) Seek(key []byte) (err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_cursor_seek(curs.handle, k, C.int(len(key)), C.UNQLITE_CURSOR_MATCH_EXACT)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// SeekLE will search the cursor and left pointing to the largest key in the database that is smaller than (pKey/nKey).
// If the database contains no keys smaller than (pKey/nKey), the cursor is left at EOF.
// This option have sense only if the underlying key/value storage subsystem support range search (i.e: B+Tree, R+Tree, etc.).
// Otherwise this option is ignored and an exact match is performed.
func (curs *Cursor) SeekLE(key []byte) (err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_cursor_seek(curs.handle, k, C.int(len(key)), C.UNQLITE_CURSOR_MATCH_LE)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// SeekGE will search the cursor and left pointing to the smallest key in the database that is larger than (pKey/nKey).
// If the database contains no keys larger than (pKey/nKey), the cursor is left at EOF.
// This option have sense only if the underlying key/value storage subsystem support range search (i.e: B+Tree, R+Tree, etc.).
// Otherwise this option is ignored and an exact match is performed.
func (curs *Cursor) SeekGE(key []byte) (err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_cursor_seek(curs.handle, k, C.int(len(key)), C.UNQLITE_CURSOR_MATCH_GE)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// First returns the first entry of the cursor.
func (curs *Cursor) First() (err error) {
	res := C.unqlite_kv_cursor_first_entry(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Last returns the last entry of the cursor.
func (curs *Cursor) Last() (err error) {
	res := C.unqlite_kv_cursor_last_entry(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// IsValid returns a boolean indicating if the cursor is valid.
func (curs *Cursor) IsValid() (ok bool) {
	return C.unqlite_kv_cursor_valid_entry(curs.handle) == 1
}

// Next moves the cursor to the next entry.
func (curs *Cursor) Next() (err error) {
	res := C.unqlite_kv_cursor_next_entry(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Prev moves the cursor the the previous entry.
func (curs *Cursor) Prev() (err error) {
	res := C.unqlite_kv_cursor_prev_entry(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Delete current cursor entry.
func (curs *Cursor) Delete() (err error) {
	res := C.unqlite_kv_cursor_delete_entry(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Reset the cursor.
func (curs *Cursor) Reset() (err error) {
	res := C.unqlite_kv_cursor_reset(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Key returns the key at the current cursor location.
func (curs *Cursor) Key() (key []byte, err error) {
	var n C.int
	res := C.unqlite_kv_cursor_key(curs.handle, nil, &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
		return
	}

	key = make([]byte, int(n))
	res = C.unqlite_kv_cursor_key(curs.handle, unsafe.Pointer(&key[0]), &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Value returns the value at the current cursor position.
func (curs *Cursor) Value() (value []byte, err error) {
	var n C.unqlite_int64
	res := C.unqlite_kv_cursor_data(curs.handle, nil, &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
		return
	}

	value = make([]byte, int(n))
	res = C.unqlite_kv_cursor_data(curs.handle, unsafe.Pointer(&value[0]), &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}
