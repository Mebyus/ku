#define LIBC_CLOCK_MONOTONIC 1

s32 // linkname
clock_gettime(uint clock_id, TimeDur* t);

void // linkname
free(void* ptr);

static TimeDur
clock_mono() {
    TimeDur dur;
    clock_gettime(LIBC_CLOCK_MONOTONIC, &dur);
    return dur;
}
