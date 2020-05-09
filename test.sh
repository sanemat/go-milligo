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
assert 21 '5+20-4'
assert 41 " 12 + 34 - 5 "
assert 47 '5+6*7'
assert 5 '3+6/3'
assert 15 '5*(9-6)'
assert 4 '(3+5)/2'

echo OK
