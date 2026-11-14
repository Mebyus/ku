#define OS_LINUX_AMD64_SYSCALL_WRITE 1

static sint
os_linux_amd64_syscall_write(uint fd, void* buf, uint size)
{
    sint ret;
    __asm__ volatile
    (
        "syscall"
		
        // outputs
		// RAX
        : "=a" (ret)

        // inputs
		// RAX
        : "0"(OS_LINUX_AMD64_SYSCALL_WRITE), 
        //  RDI      RSI       RDX
			"D"(fd), "S"(buf), "d"(size)

		// two registers are clobbered after system call
        : "rcx", "r11", 
			"memory"
    );
    return ret;
}

#define OS_LINUX_STDOUT 1

static void
print(str s) {
	if (s.num == 0) {
		return;
	}

    uint i = 0;
	while (i < s.num) {
		sint r = os_linux_amd64_syscall_write(OS_LINUX_STDOUT, s.ptr + i, s.num - i);
        if (r < 0) {
            // stop writing after any error
            return;
        }
        i += cast(uint, r);
	}
}
