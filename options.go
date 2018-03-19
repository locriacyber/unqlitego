package unqlitego

// #include <unqlite.h>
// #include <wrappers.h>
// #include <stdlib.h>
import "C"

// IsThreadSafe returns a boolean identifiying if the UnQlite library is
// compiled Thread Safe.
func IsThreadSafe() bool {
	return C.unqlite_lib_is_threadsafe() == 1
}

// Version returns the version string.
func Version() string {
	return C.GoString(C.unqlite_lib_version())
}

// Signature returns the Signature string.
func Signature() string {
	return C.GoString(C.unqlite_lib_signature())
}

// Ident returns the UnQlite identification string.
func Ident() string {
	return C.GoString(C.unqlite_lib_ident())
}

// Copyright returns the Copyright string.
func Copyright() string {
	return C.GoString(C.unqlite_lib_copyright())
}
