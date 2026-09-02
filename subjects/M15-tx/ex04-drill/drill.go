package main

// Drill is a deterministic chaos drill over a small cluster. Nodes are
// up or down; each op touches one node and succeeds only if that node is up
// "at the current tick". Time advances as explicit Tick steps — never the
// wall clock.
type Drill struct {
	up   map[string]bool
	tick int
}

func NewDrill(up []string) *Drill {
	d := &Drill{up: make(map[string]bool)}
	for _, n := range up {
		d.up[n] = true
	}
	return d
}

func (d *Drill) Tick() {
	d.tick++
}

func (d *Drill) Up(node string) {
	d.up[node] = true
}

func (d *Drill) Down(node string) {
	d.up[node] = false
}

// Op reports "ok" if the node is up at the current tick, else "fail".
func (d *Drill) Op(node string) string {
	if d.up[node] {
		return "ok"
	}
	return "fail"
}
