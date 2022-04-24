(module
  (import "console" "log" (func $writeln_i32 (param $value i32)))
  (import "console" "log" (func $writeln_f64 (param $value f64)))
  (global $x (mut i32) (i32.const 0))
  (func (export "main")
    (global.set $x (i32.const 0))
    (if (i32.lt_s (i32.add (global.get $x) (i32.const 1)) (i32.const 5))
      (then
        (global.set $x (i32.const 100))
      )
    )
  )
)