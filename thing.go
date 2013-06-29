package reddit

import ()

type Thing struct {
	Id   string
	Name string
	Kind string
}

func (t *Thing) Report() error {
	return nil
}

func (t *Thing) Hide() error {
	return nil
}

func (t *Thing) Unhide() error {
	return nil
}

func (t *Thing) Info() error {
	return nil
}

func (t *Thing) MarkNsfw() error {
	return nil
}

func (t *Thing) UnmarkNsfw() error {
	return nil
}
