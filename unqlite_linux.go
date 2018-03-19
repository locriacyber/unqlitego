// +build linux,cgo

package unqlitego

// #cgo LDFLAGS: -lm -lpthread
// #cgo CFLAGS: -DUNQLITE_ENABLE_THREADS=1 -Wno-unused-but-set-variable -DJX9_ENABLE_MATH_FUNC -DUNQLITE_ENABLE_JX9_HASH_IO
import "C"
