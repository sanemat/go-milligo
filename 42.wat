(module
    ;; Import the required proc_exit WASI function
    (import "wasi_unstable" "proc_exit" (func $proc_exit (param i32)))

    (memory 1)
    (export "memory" (memory 0))

    (func $main (export "_start")
        (call $proc_exit
            (i32.const 42) ;; exit code
        )
    )
)
