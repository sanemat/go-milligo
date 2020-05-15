(module
    ;; Import the required proc_exit WASI function
    (import "wasi_unstable" "proc_exit" (func $proc_exit (param i32)))

    (memory 1)
    (export "memory" (memory 0))

    (func $main (export "_start")
        (local $tmp i32)
        (call $proc_exit
            i32.const 5
            i32.const 20
            i32.add
            set_local $tmp
            i32.const 5
            set_local $tmp
            get_local $tmp
        )
    )
)
