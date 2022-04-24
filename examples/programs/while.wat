(module
  (import "console" "log" (func $writeln_i32 (param $value i32)))
  (import "console" "log" (func $writeln_f64 (param $value f64)))
  (global $i (mut i32) (i32.const 0))
  (global $count (mut i32) (i32.const 0))
  (func (export "main")
    (global.set $count (i32.const 1))
    (global.set $i (i32.const 0))
    (if (i32.lt_s (global.get $i) (i32.const 5))
      (then
        (loop
          (call $writeln_i32 (global.get $i))
          (call $writeln_i32 (global.get $count))
          (global.set $i (i32.add (global.get $i) (i32.const 1)))
          (global.set $count (i32.mul (global.get $count) (i32.const 2)))
          (br_if 0 (i32.lt_s (global.get $i) (i32.const 5)))
        )
      )
    )
    (call $writeln_i32 (global.get $i))
    (call $writeln_i32 (global.get $count))
  )
)