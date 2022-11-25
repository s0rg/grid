package grid

type pqueue []path

func (q pqueue) Len() int           { return len(q) }
func (q pqueue) Less(i, j int) bool { return q[i].Cost < q[j].Cost }
func (q pqueue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *pqueue) Push(x any) {
	*q = append(*q, x.(path))
}

func (q *pqueue) Pop() (x any) {
	t := *q
	n := t.Len()

	x = t[n-1]
	*q = t[0 : n-1]

	return x
}
