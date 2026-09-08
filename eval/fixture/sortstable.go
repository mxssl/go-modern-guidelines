package legacy

import "sort"

// Task is a unit of queued work. Tasks with equal priority must keep the
// order in which they were enqueued.
type Task struct {
	ID       int
	Priority int
}

// OrderByPriority returns tasks sorted by descending priority. Tasks with the
// same priority stay in enqueue order.
func OrderByPriority(tasks []Task) []Task {
	out := make([]Task, len(tasks))
	copy(out, tasks)
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Priority > out[j].Priority
	})
	return out
}
