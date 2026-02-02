package repo

import (
	"sync"
	"time"

	"example.com/notes-api/internal/core"
)

type NoteRepoMem struct {
	mu    sync.Mutex
	notes map[int64]*core.Note
	next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
	return &NoteRepoMem{
		notes: make(map[int64]*core.Note),
	}
}

func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.next++
	n.ID = r.next
	n.CreatedAt = time.Now()

	r.notes[n.ID] = &n
	return n.ID, nil
}

func (r *NoteRepoMem) List() []*core.Note {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]*core.Note, 0, len(r.notes))
	for _, n := range r.notes {
		result = append(result, n)
	}
	return result
}

func (r *NoteRepoMem) Get(id int64) (*core.Note, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notes[id]
	return n, ok
}

func (r *NoteRepoMem) Update(id int64, input core.NoteUpdate, updatedAt time.Time) (*core.Note, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, ok := r.notes[id]
	if !ok {
		return nil, false
	}

	if input.Title != nil {
		n.Title = *input.Title
	}
	if input.Content != nil {
		n.Content = *input.Content
	}

	n.UpdatedAt = &updatedAt
	return n, true
}

func (r *NoteRepoMem) Delete(id int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.notes[id]; !ok {
		return false
	}

	delete(r.notes, id)
	return true
}
