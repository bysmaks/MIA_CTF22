global _start

section .text

_print:
	mov eax, 4
	mov ebx, 1
	lea ecx, [esp+8]
	mov edx, dword [esp+4]
	int 0x80

	ret

_start:
	push 0x0
	push 0x0a216f
	push 0x6c6c6548
	push 9
	call _print

	xor eax, eax
	push eax
	mov eax, 3
	xor ebx, ebx
	lea ecx, [esp]
	mov edx, 4
	int 0x80

	mov esi, [esp]

	cmp esi, 0x0a465443
	je _equal

	push 0x00
	push 0x0a706f4e
	push 5
	call _print

	jmp _exit

_equal:
	push 0x0a7d396e
	push 0x316e7233
	push 0x6c5f7a33
	push 0x5f6d3534
	push 0x7b465443
	push 20
	call _print

_exit:
	; exit
	mov eax, 1
	xor ebx, ebx
	int 0x80
