(module
  (type (;0;) (func (param i32)))
  (type (;1;) (func (param i32 i32)))
  (type (;2;) (func (result i32)))
  (type (;3;) (func (param i32) (result i32)))
  (import "env" "print_s" (func (;0;) (type 0)))
  (import "env" "TrigEvent" (func (;1;) (type 1)))
  (func (;2;) (type 2) (result i32)
    i32.const 16384
    call 0
    i32.const 16403)
  (func (;3;) (type 0) (param i32)
    i32.const 16404
    call 0
    local.get 0
    call 0
    i32.const 16447
    i32.const 16427
    call 1)
  (func (;4;) (type 3) (param i32) (result i32)
    i32.const 16464
    call 0
    i32.const 1)
  (func (;5;) (type 3) (param i32) (result i32)
    i32.const 16494
    call 0
    i32.const 16447
    i32.const 16427
    call 1
    i32.const 16403)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 1)
  (global (;0;) (mut i32) (i32.const 16384))
  (global (;1;) i32 (i32.const 16527))
  (global (;2;) i32 (i32.const 16527))
  (export "memory" (memory 0))
  (export "__heap_base" (global 1))
  (export "__data_end" (global 2))
  (export "init" (func 2))
  (export "testFunc" (func 3))
  (export "testFuncWithInt" (func 4))
  (export "testFuncWithString" (func 5))
  (data (;0;) (i32.const 16384) "TestContract::init\00")
  (data (;1;) (i32.const 16403) "\00")
  (data (;2;) (i32.const 16404) "TestContract::testFunc\00")
  (data (;3;) (i32.const 16427) "{\22Param\22:\22testStr\22}\00")
  (data (;4;) (i32.const 16447) "testFunc(string)\00")
  (data (;5;) (i32.const 16464) "TestContract::testFuncWithInt\00")
  (data (;6;) (i32.const 16494) "TestContract::testFuncWithString\00"))
