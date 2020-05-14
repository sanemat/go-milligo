(module
    ;; Import the required proc_exit WASI function
    (import "wasi_unstable" "proc_exit" (func $proc_exit (param i32)))

    (memory 1)
    (export "memory" (memory 0))
    (func $__ (result i32)
        (local $tmp i32)
        i32.const 5
        set_local $tmp
        i32.const 20
        return
        set_local $tmp
        i32.const 6
        set_local $tmp
        get_local $tmp
    )

    (func $main (export "_start")
        (call $proc_exit
            (call $__)
        )
    )
)
