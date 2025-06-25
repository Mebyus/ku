package butler

import "strings"

type formatBuffer struct {
	buf strings.Builder
}

func (g *formatBuffer) nl() {
	g.putb('\n')
}

func (g *formatBuffer) space() {
	g.putb(' ')
}

func (g *formatBuffer) indent() {
	g.putb('\t')
}

func (g *formatBuffer) puts(s string) {
	_, _ = g.buf.WriteString(s)
}

func (g *formatBuffer) putb(b byte) {
	_ = g.buf.WriteByte(b)
}

func (g *formatBuffer) String() string {
	return g.buf.String()
}
