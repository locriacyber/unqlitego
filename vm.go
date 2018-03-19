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

func (vm *VM) Unqlite_vm_extract_variable(variable_name string) *Unqlite_value {
	/*This function must be used with extra causion since it might return
	a variable from the type of *C.unqlite_value ,be sure to free this pointer
	In case of no such variable of out-of-memory issue NULL is returned

	In case where the VM have not been executed the return value will be C.int(0)
	In case of unqlite is compiled with threads support and the vm.vm instance have been releases
	by a different thread 0 will be returned

	For summary:
	-----------
	*) 0 or NULL = Bad
	*) *C.unqlite_value = Good */
	c_variable_name := C.CString(variable_name)
	defer C.free(unsafe.Pointer(c_variable_name))
	unqlite_value := C.unqlite_vm_extract_variable(vm.vm, c_variable_name)
	return &Unqlite_value{unqlite_value}
}

func (vm *VM) VM_extract_output() string {
	//extract the output from a vm after execution
	var len C.int
	var buff *C.char
	buff = C.extract_vm_output(vm.vm, &len)
	q := C.GoStringN(buff, len)
	fmt.Printf("%s", q)
	return ""
}

func (vm *VM) VM_execute() int {
	res := C.unqlite_vm_exec(vm.vm)
	return int(res)
}

func (vm *VM) Extract_variable_as_int(variable_name string) (int, error) {
	/*Extract a variable from the VM after if have been executed
	If something went wrong return nil
	*/
	var unqlite_value *Unqlite_value
	unqlite_value = vm.Unqlite_vm_extract_variable(variable_name)
	if !unqlite_value_ok(unqlite_value) {
		return 0, nil
	}

	res := int(C.unqlite_value_to_int(unqlite_value.unqlite_value))
	return res, GlobaLError("OK")
}

func (vm *VM) Extract_variable_as_string(variable_name string) (string, error) {
	/*Extract a variable from the VM after if have been executed
	If something went wrong return nil
	*/
	var unqlite_value *Unqlite_value
	unqlite_value = vm.Unqlite_vm_extract_variable(variable_name)
	if !unqlite_value_ok(unqlite_value) {
		return "", nil
	}
	var plen C.int
	c_res := C.extract_variable_as_string(unqlite_value.unqlite_value, &plen)
	res := C.GoStringN(c_res, plen)

	return res, GlobaLError("OK")
}

func (vm *VM) Extract_variable_as_bool(variable_name string) (bool, error) {
	/*Extract a variable from the VM after if have been executed
	If something went wrong return nil
	*/
	var unqlite_value *Unqlite_value
	unqlite_value = vm.Unqlite_vm_extract_variable(variable_name)
	if !unqlite_value_ok(unqlite_value) {
		return false, nil
	}

	res := int(C.unqlite_value_to_bool(unqlite_value.unqlite_value))
	return res != 0, GlobaLError("OK")
}

func (vm *VM) Extract_variable_as_int64(variable_name string) (int64, error) {
	/*Extract a variable from the VM after if have been executed
	If something went wrong return nil
	*/
	var unqlite_value *Unqlite_value
	unqlite_value = vm.Unqlite_vm_extract_variable(variable_name)
	if !unqlite_value_ok(unqlite_value) {
		return 0, nil
	}

	res := int64(C.unqlite_value_to_int64(unqlite_value.unqlite_value))
	return res, GlobaLError("OK")
}

func (vm *VM) Extract_variable_as_double(variable_name string) (float64, error) {
	/*Extract a variable from the VM after if have been executed
	If something went wrong return nil
	*/
	var unqlite_value *Unqlite_value
	unqlite_value = vm.Unqlite_vm_extract_variable(variable_name)
	if !unqlite_value_ok(unqlite_value) {
		return 0.0, nil
	}

	res := float64(C.unqlite_value_to_double(unqlite_value.unqlite_value))
	return res, GlobaLError("OK")
}
