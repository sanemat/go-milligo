(module
    ;; Import the required proc_exit WASI function
    (import "wasi_unstable" "proc_exit" (func $proc_exit (param i32)))

    (memory 1)
    (export "memory" (memory 0))

    (func $main (export "_start")
        (call $proc_exit
            i32.const 5
            i32.const 20
            i32.add
            i32.const 4
            i32.sub
        )
    )
)
