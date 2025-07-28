package genc

func (g *Gen) MainTestDriver(tests []string) {
	g.nl()
	if len(tests) == 0 {
		g.puts("static uint")
		g.nl()
		g.puts("ku_main() {}")
		g.nl()
		return
	}

	g.puts("#define TEST_OUTPUT_BUFFER_SIZE 1 << 16")
	g.nl()
	g.nl()

	g.puts("static u8 test_output_buffer[TEST_OUTPUT_BUFFER_SIZE];")
	g.nl()
	g.nl()

	g.puts("static uint")
	g.nl()
	g.puts("ku_main() {")
	g.nl()
	g.inc()
	g.level += 1

	g.indent()
	g.puts("uint fail_count = 0;")
	g.nl()

	g.indent()
	g.puts("FormatCapBuffer buf;")
	g.nl()

	g.indent()
	g.puts("init_fmt_cap_buffer(&buf, make_span_u8(test_output_buffer, TEST_OUTPUT_BUFFER_SIZE));")
	g.nl()

	g.indent()
	g.puts("TestContext t;")
	g.nl()
	g.nl()

	for _, t := range tests {
		name := getTestFunName(t)

		g.indent()
		g.puts("test_reset(&t, ")
		g.str(t)
		g.puts(");")
		g.nl()

		g.indent()
		g.puts(name)
		g.puts("(&t);")
		g.nl()

		g.indent()
		g.puts("if (t.failed) {")
		g.nl()
		g.inc()

		g.indent()
		g.puts("fmt_cap_buffer_put_test(&buf, &t);")
		g.nl()

		g.indent()
		g.puts("fail_count += 1;")
		g.nl()

		g.dec()
		g.indent()
		g.puts("}")
		g.nl()

		g.nl()
	}

	g.nl()

	g.indent()
	g.puts("stdout_print(fmt_cap_buffer_take(&buf));")
	g.nl()

	g.indent()
	g.puts("if (fail_count != 0) {")
	g.nl()
	g.inc()

	g.indent()
	g.puts("os_exit(1);")
	g.nl()

	g.dec()
	g.indent()
	g.puts("}")
	g.nl()

	g.nl()
	g.indent()
	g.puts("return 0;")
	g.nl()

	g.level -= 1
	g.dec()
	g.puts("}")
	g.nl()
}
