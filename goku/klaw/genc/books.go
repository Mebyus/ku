package genc

func (g *Gen) NameBooks() {
	g.typesBook()
	g.errorsBook()
}

func (g *Gen) errorsBook() {
	records := g.State.ErrorRecords()

	g.nl()
	g.puts("/*/errors")
	g.nl()

	for _, r := range records {
		g.putn(r.Id)
		g.space()
		g.puts(r.Name)
		g.nl()
	}

	g.puts("*/")
	g.nl()
}

func (g *Gen) typesBook() {

}
