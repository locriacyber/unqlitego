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
	// Database Pointer
	db *Database

	// Cursor Handle
	handle *C.unqlite_kv_cursor
}

// Cursor creates and initializes a new UnQLite database cursor.
func (db *Database) Cursor() (*Cursor, error) {
	var err error

	c := &Cursor{db: db}
	res := C.unqlite_kv_cursor_init(db.conn, &c.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	runtime.SetFinalizer(c, (*Cursor).Close)

	return nil, err
}

// Close closes the open cursor.
func (cr *Cursor) Close() error {
	var err error

	if cr.db.conn != nil && cr.handle != nil {
		res := C.unqlite_kv_cursor_release(cr.db.conn, cr.handle)
		if res != C.UNQLITE_OK {
			err = UnQLiteError(res)
		}
		cr.handle = nil
	}

	return err
}

// Seek will search the cursor for an exact match. If the record exists the cursor is left pointing to it.
// Otherwise it is left pointing to EOF and UNQLITE_NOTFOUND is returned.
func (cr *Cursor) Seek(key []byte) error {
	var err error
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_cursor_seek(cr.handle, k, C.int(len(key)), C.UNQLITE_CURSOR_MATCH_EXACT)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return err
}

// SeekLE will search the cursor and left pointing to the largest key in the database that is smaller than (pKey/nKey).
// If the database contains no keys smaller than (pKey/nKey), the cursor is left at EOF.
// This option have sense only if the underlying key/value storage subsystem support range search (i.e: B+Tree, R+Tree, etc.).
// Otherwise this option is ignored and an exact match is performed.
func (cr *Cursor) SeekLE(key []byte) error {
	var err error
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_cursor_seek(cr.handle, k, C.int(len(key)), C.UNQLITE_CURSOR_MATCH_LE)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return err
}

// SeekGE will search the cursor and left pointing to the smallest key in the database that is larger than (pKey/nKey).
// If the database contains no keys larger than (pKey/nKey), the cursor is left at EOF.
// This option have sense only if the underlying key/value storage subsystem support range search (i.e: B+Tree, R+Tree, etc.).
// Otherwise this option is ignored and an exact match is performed.
func (cr *Cursor) SeekGE(key []byte) error {
	var err error
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_cursor_seek(cr.handle, k, C.int(len(key)), C.UNQLITE_CURSOR_MATCH_GE)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return err
}

// First returns the first entry of the cursor.
func (cr *Cursor) First() error {
	var err error

	res := C.unqlite_kv_cursor_first_entry(cr.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return err
}

// Last returns the last entry of the cursor.
func (cr *Cursor) Last() error {
	var err error

	res := C.unqlite_kv_cursor_last_entry(cr.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return err
}

// IsValid returns a boolean indicating if the cursor is valid.
func (cr *Cursor) IsValid() bool {
	return C.unqlite_kv_cursor_valid_entry(cr.handle) == 1
}

// Next moves the cursor to the next entry.
func (cr *Cursor) Next() error {
	var err error

	res := C.unqlite_kv_cursor_next_entry(cr.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return err
}

// Prev moves the cursor the the previous entry.
func (cr *Cursor) Prev() error {
	var err error

	res := C.unqlite_kv_cursor_prev_entry(cr.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return err
}

// Delete current cursor entry.
func (cr *Cursor) Delete() error {
	var err error

	res := C.unqlite_kv_cursor_delete_entry(cr.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return err
}

// Reset the cursor.
func (cr *Cursor) Reset() error {
	var err error

	res := C.unqlite_kv_cursor_reset(cr.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return err
}

// Key returns the key at the current cursor location.
func (cr *Cursor) Key() (key []byte, err error) {
	var n C.int

	res := C.unqlite_kv_cursor_key(cr.handle, nil, &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
		return
	}

	key = make([]byte, int(n))
	res = C.unqlite_kv_cursor_key(cr.handle, unsafe.Pointer(&key[0]), &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}

// Value returns the value at the current cursor position.
func (cr *Cursor) Value() (value []byte, err error) {
	var n C.unqlite_int64

	res := C.unqlite_kv_cursor_data(cr.handle, nil, &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
		return
	}

	value = make([]byte, int(n))
	res = C.unqlite_kv_cursor_data(cr.handle, unsafe.Pointer(&value[0]), &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}

	return
}
