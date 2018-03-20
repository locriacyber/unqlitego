// +build darwin,cgo

package unqlitego

// #cgo darwin LDFLAGS: -lm -lpthread
// #cgo darwin CFLAGS: -DUNQLITE_ENABLE_THREADS -DJX9_ENABLE_MATH_FUNC -DUNQLITE_ENABLE_JX9_HASH_IO
import "C"
