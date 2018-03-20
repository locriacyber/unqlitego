package unqlitego

// #include <unqlite.h>
// #include <wrappers.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"unsafe"
)

// VM Represents an UnQLite/Jx9 Virtual Machine.
type VM struct {
	vm *C.unqlite_vm
}

// NewVM creates and intializes a new UnQlite/Jx9 Virtual Machine.
func NewVM() (vm *VM) {
	vm = &VM{}
	return
}

// Free releases the Virtual Machine.
func (vm VM) Free() {
	C.free(unsafe.Pointer(vm.vm))
}

// Execute Virtual Machine.
func (vm *VM) Execute() int {
	res := C.unqlite_vm_exec(vm.vm)

	return int(res)
}

// Result will return the output of the Virtual Machine after execution.
func (vm *VM) Result() string {
	var len C.int
	var buff *C.char

	buff = C.extract_vm_output(vm.vm, &len)
	q := C.GoStringN(buff, len)

	return fmt.Sprintf("%s", q)
}

/*
	This function must be used with extra causion since it might return
	a variable from the type of *C.unqlite_value ,be sure to free this pointer
	In case of no such variable of out-of-memory issue NULL is returned

	In case where the VM have not been executed the return value will be C.int(0)
	In case of unqlite is compiled with threads support and the vm.vm instance have been releases
	by a different thread 0 will be returned

	For summary:
	-----------
	*) 0 or NULL = Bad
	*) *C.unqlite_value = Good
*/
func (vm *VM) extract(v string) *unQLiteValue {
	uval := &unQLiteValue{}
	cvar := C.CString(v)

	// Mark cvar for Clean-up
	defer C.free(unsafe.Pointer(cvar))

	// Extract Value from VM
	uval.v = C.unqlite_vm_extract_variable(vm.vm, cvar)

	return uval
}

// ExtractInt will extract the result from the Virtual Machine as int.
func (vm *VM) ExtractInt(v string) (int, error) {
	/*
		Extract a variable from the VM after if have been executed
		If something went wrong return nil
	*/
	uval := vm.extract(v)
	if uval.Nil() {
		return 0, nil
	}

	res := int(C.unqlite_value_to_int(uval.v))

	return res, UnQLiteError(0)
}

// ExtractString will extract the result from the Virtual Machine as string.
func (vm *VM) ExtractString(v string) (string, error) {
	/*
		Extract a variable from the VM after if have been executed
		If something went wrong return nil
	*/
	uval := vm.extract(v)
	if uval.Nil() {
		return "", nil
	}

	var plen C.int
	cvar := C.extract_variable_as_string(uval.v, &plen)
	res := C.GoStringN(cvar, plen)

	return res, UnQLiteError(0)
}

// ExtractBool will extract the result from the Virtual Machine as bool.
func (vm *VM) ExtractBool(v string) (bool, error) {
	/*
		Extract a variable from the VM after if have been executed
		If something went wrong return nil
	*/
	uval := vm.extract(v)
	if uval.Nil() {
		return false, nil
	}

	res := int(C.unqlite_value_to_bool(uval.v))

	return res != 0, UnQLiteError(0)
}

// ExtractInt64 will extract the result from the Virtual Machine as int64.
func (vm *VM) ExtractInt64(v string) (int64, error) {
	/*
		Extract a variable from the VM after if have been executed
		If something went wrong return nil
	*/
	uval := vm.extract(v)
	if uval.Nil() {
		return 0, nil
	}

	res := int64(C.unqlite_value_to_int64(uval.v))

	return res, UnQLiteError(0)
}

// ExtractFloat64 will extract the result from the Virtual Machine as float64.
func (vm *VM) ExtractFloat64(v string) (float64, error) {
	/*
		Extract a variable from the VM after if have been executed
		If something went wrong return nil
	*/
	uval := vm.extract(v)
	if uval.Nil() {
		return 0.0, nil
	}

	res := float64(C.unqlite_value_to_double(uval.v))

	return res, UnQLiteError(0)
}
