static uint __attribute__ ((noinline))
save_stack_trace(void** buf, uint size) {
    void** frame_pointer = __builtin_frame_address(0);
    uint i = 0;
    while (i < size) {
        void** next_frame_pointer = *frame_pointer; 
        void*  pc = frame_pointer[1];

        if (next_frame_pointer <= frame_pointer) {
            return i;
        }

        buf[i] = pc;
        frame_pointer = next_frame_pointer;
        i += 1;
    }
    return i;
}
