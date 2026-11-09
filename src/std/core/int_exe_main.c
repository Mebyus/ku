#define ARGS_BUFFER_SIZE 64

int main(int argc, char** argv) {
    str args_buf[ARGS_BUFFER_SIZE];

    uint n = cast(uint, argc);
    if (n > ARGS_BUFFER_SIZE) {
        n = ARGS_BUFFER_SIZE;
    }

	for (uint i = 0; i < n; i += 1) {
		u8* ptr = cast(u8*, argv[i]);
		if (ptr == nil) {
            n = i;
			break;
		}
		args_buf[i] = unsafe_make_c_string_from_ptr(ptr);
	}

    s_str args = make_span_str(args_buf, n);
    ku_run(args);
}
