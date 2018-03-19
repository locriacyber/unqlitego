// +build linux,cgo

package unqlitego

// #cgo linux LDFLAGS: -lm -lpthread
// #cgo linux CFLAGS: -DUNQLITE_ENABLE_THREADS -DJX9_ENABLE_MATH_FUNC -DUNQLITE_ENABLE_JX9_HASH_IO -Wno-unused-but-set-variable
import "C"
