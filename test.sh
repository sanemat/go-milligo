#!/bin/bash
assert() {
  expected="$1"
  input="$2"

  ./milligo "$input" > tmp.wat
  wasmtime tmp.wat
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

assert 42 42
assert 0 0
assert 9 '10-1'
assert 11 '10+1'

echo OK
