// static u64
// amd64_rdtsc() {
//     u64 low, high;
//     __asm__ volatile ("rdtsc" : "=a" (low), "=d" (high));
//     return (high << 32) | low;
// }

// static u64
// cpu_clock() {
//     return amd64_rdtsc();
// }

#define OS_LINUX_ERROR_CODE_NOT_EXIST 2

#define OS_LINUX_AMD64_SYSCALL_READ 0

static sint
os_linux_amd64_syscall_read(uint fd, void* buf, uint size)
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
        : "0"(OS_LINUX_AMD64_SYSCALL_READ), 
        //  RDI      RSI       RDX
			"D"(fd), "S"(buf), "d"(size)

		// two registers are clobbered after system call
        : "rcx", "r11", 
			"memory"
    );
    return ret;
}

#define OS_LINUX_AMD64_SYSCALL_WRITE 1

static sint
os_linux_amd64_syscall_write(uint fd, const void* buf, uint size)
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


#define OS_LINUX_AMD64_SYSCALL_OPEN 2

static sint
os_linux_amd64_syscall_open(const u8* path, u32 flags, u32 mode)
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
        : "0"(OS_LINUX_AMD64_SYSCALL_OPEN), 
        //  RDI      RSI       RDX
			"D"(path), "S"(flags), "d"(mode)

		// two registers are clobbered after system call
        : "rcx", "r11", 
			"memory"
    );
    return ret;
}

#define OS_LINUX_AMD64_SYSCALL_CLOSE 3

static sint
os_linux_amd64_syscall_close(uint fd)
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
        : "0"(OS_LINUX_AMD64_SYSCALL_CLOSE), 
        //  RDI
			"D"(fd)

		// two registers are clobbered after system call
        : "rcx", "r11", 
			"memory"
    );
    return ret;
}

typedef struct {
  s64 sec;
  s64 nano;
} LinuxTimespec;

typedef struct {
    // ID of device containing file
    u64 device;

    // Inode number
    u64 inode;

    // Number of hard links
    u64 num_links;

    // File type and mode
    u32 mode;

    // User ID of owner
    u32 user_id;

    // Group ID of owner
    u32 group_id;

    // This field is placed here purely for padding. It does not
    // carry any useful information
    u32 padding;

    // Device ID (if special file)
    u64 r_device;

    // Total size in bytes
    u64 size;

    // Block size for filesystem I/O
    u64 block_size;

    // Number of 512B blocks allocated
    u64 num_blocks;

    // Time of last access
    LinuxTimespec access_time;

    // Time of last modification
    LinuxTimespec mod_time;

    // Time of last status change
    LinuxTimespec status_change_time;

    // Additional reserved fields for future compatibility
    u64 reserved[3];
} LinuxFileStat;

#define OS_LINUX_AMD64_SYSCALL_STAT 4

static sint
os_linux_amd64_syscall_stat(const u8* path, LinuxFileStat* stat)
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
        : "0"(OS_LINUX_AMD64_SYSCALL_STAT), 
        //  RDI      RSI
			"D"(path), "S"(stat)

		// two registers are clobbered after system call
        : "rcx", "r11", 
			"memory"
    );
    return ret;
}

#define OS_LINUX_AMD64_SYSCALL_MMAP 9

#define OS_LINUX_MEMORY_MAP_PROT_READ  0x1
#define OS_LINUX_MEMORY_MAP_PROT_WRITE 0x2

#define OS_LINUX_MEMORY_MAP_SHARED    0x01
#define OS_LINUX_MEMORY_MAP_PRIVATE   0x02
#define OS_LINUX_MEMORY_MAP_ANONYMOUS 0x20

static sint
os_linux_amd64_syscall_mmap(void* ptr, uint len, uint prot, uint flags, uint fd, uint offset) {
    register sint  rax __asm__ ("rax") = OS_LINUX_AMD64_SYSCALL_MMAP;
    register void* rdi __asm__ ("rdi") = ptr;
    register uint  rsi __asm__ ("rsi") = len;
    register uint  rdx __asm__ ("rdx") = prot;
    register uint  r10 __asm__ ("r10") = flags;
    register uint  r8  __asm__ ("r8")  = fd;
    register uint  r9  __asm__ ("r9")  = offset;
    __asm__ __volatile__ (
        "syscall"
        : "+r" (rax)
        : "r" (rdi), "r" (rsi), "r" (rdx), "r" (r10), "r" (r8), "r" (r9)
        : "rcx", "r11", "memory"
    );
    return rax;
}

#define OS_LINUX_AMD64_SYSCALL_MUNMAP 11

static sint
os_linux_amd64_syscall_munmap(void* ptr, uint len) {
    register sint  rax __asm__ ("rax") = OS_LINUX_AMD64_SYSCALL_MUNMAP;
    register void* rdi __asm__ ("rdi") = ptr;
    register uint  rsi __asm__ ("rsi") = len;
    __asm__ __volatile__ (
        "syscall"
        : "+r" (rax)
        : "r" (rdi), "r" (rsi)
        : "rcx", "r11", "memory"
    );
    return rax;
}

#define OS_LINUX_AMD64_SYSCALL_IOCTL 16

static sint
os_linux_amd64_syscall_ioctl(uint fd, uint op, void* ptr) {
    register sint  rax __asm__ ("rax") = OS_LINUX_AMD64_SYSCALL_IOCTL;
    register uint  rdi __asm__ ("rdi") = fd;
    register uint  rsi __asm__ ("rsi") = op;
    register void* rdx __asm__ ("rdx") = ptr;
    __asm__ __volatile__ (
        "syscall"
        : "+r" (rax)
        : "r" (rdi), "r" (rsi), "r" (rdx)
        : "rcx", "r11", "memory"
    );
    return rax;
}

#define OS_LINUX_AMD64_SYSCALL_EXIT 60

static _Noreturn void
os_linux_amd64_syscall_exit(uint c) {
    sint ret;
    __asm__ volatile
    (
        "syscall"
		
		// outputs
		// RAX
        : "=a" (ret)
		
		// inputs
		// RAX
        : "0"(OS_LINUX_AMD64_SYSCALL_EXIT), 
        //  RDI     
			"D"(c)

		// clobbers
		// two registers are clobbered after system call
        : "rcx", "r11"
    );
	__builtin_unreachable();
}
