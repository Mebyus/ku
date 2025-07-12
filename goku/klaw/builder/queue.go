package builder

type QueueItem struct {
	// Path to unit inside the queue.
	Path string

	// Place where unit with this path is imported.
	//
	// There may be more than one import place for a specific unit inside the
	// whole program. This field tracks the first one we encounter during unit
	// walk phase.
	// Pin srcmap.Pin
}

// UnitQueue keeps track which unit were already visited and
// which unit should be visited next during unit discovery phase.
type UnitQueue struct {
	// List of paths waiting in queue.
	backlog []QueueItem

	// List of all collected units.
	units []*Unit

	// Set which contains paths of already visited units.
	visited map[string]struct{}
}

func NewUnitQueue() *UnitQueue {
	return &UnitQueue{
		visited: make(map[string]struct{}),
	}
}

// Add tries to add item to backlog. If a given path was
// already visited then this call will be no-op.
func (q *UnitQueue) Add(item QueueItem) {
	_, ok := q.visited[item.Path]
	if ok {
		// given path is already known to walker
		return
	}

	q.visited[item.Path] = struct{}{}
	q.backlog = append(q.backlog, item)
}

// Next get next item from backlog and write it to a given pointer.
// Returns true if operation was successfull and false when there are no items left.
func (q *UnitQueue) Next(item *QueueItem) bool {
	if len(q.backlog) == 0 {
		return false
	}

	last := len(q.backlog) - 1
	*item = q.backlog[last]

	// shrink slice, but keep its capacity
	q.backlog = q.backlog[:last]
	return true
}

func (q *UnitQueue) AddUnit(unit *Unit) {
	q.units = append(q.units, unit)

	for _, p := range unit.Imports {
		q.Add(QueueItem{
			Path: p,
		})
	}
}
