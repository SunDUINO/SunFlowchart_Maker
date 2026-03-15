package model

// SetNextID allows the history package to restore the ID counter.
func (d *Diagram) SetNextID(id int) {
	d.nextID = id
}
