(module
    ;; Imports.
    (import "console" "log" (func $log_i32 (param $0 i32)))

    ;; Globals.
    ;; (global $g (import "js" "global") (mut i32))
    (global $g (mut i32) (i32.const 42))
    (global $pi f64 (f64.const 3.141592653589793))
    (global $e f64 (f64.const 2.718281828459045))

    ;; Main.
    (func (export "main") (param $a i32) (param $b i32)
        (if (i32.gt_s (local.get $a) (local.get $b))
            (call $log_i32 (local.get $a))
        )

        (if (i32.eq (i32.const 1) (i32.const 0))
            (block
                (i32.const 0)
                call $log_i32
            )
            (block
                (i32.const 1)
                call $log_i32
            )
        )

        (loop
            (block
                (call $log_i32 (i32.const 123))
                (br_if 0 (i32.eq (i32.const 0) (i32.const 0)))
            )
        )

        (call $log_i32 (i32.trunc_f64_s (f64.const 3.141592653589793)))

        (call $log_i32 (i32.trunc_f64_s (f64.convert_i32_s (i32.const 13))))

        (call $log_i32 (i32.add (local.get $a) (i32.add (local.get $a) (local.get $b))))
    )
)
