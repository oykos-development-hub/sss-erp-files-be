package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// File struct
type File struct {
	ID          int       `db:"id,omitempty"`
	ParentID    *int      `db:"parent_id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	Size        int       `db:"size"`
	Type        string    `db:"type"`
	CreatedAt   time.Time `db:"created_at,omitempty"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Table returns the table name
func (t *File) Table() string {
	return "files"
}

// GetAll gets all records from the database, using upper
func (t *File) GetAll(condition *up.Cond) ([]*File, error) {
	collection := upper.Collection(t.Table())
	var all []*File
	var res up.Result

	if condition != nil {
		res = collection.Find(*condition)
	} else {
		res = collection.Find()
	}

	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, err
}

// Get gets one record from the database, by id, using upper
func (t *File) Get(id int) (*File, error) {
	var one File
	collection := upper.Collection(t.Table())

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&one)
	if err != nil {
		return nil, err
	}
	return &one, nil
}

// Update updates a record in the database, using upper
func (t *File) Update(m File) error {
	m.UpdatedAt = time.Now()
	collection := upper.Collection(t.Table())
	res := collection.Find(m.ID)
	err := res.Update(&m)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by id, using upper
func (t *File) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a model into the database, using upper
func (t *File) Insert(m File) (int, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	collection := upper.Collection(t.Table())
	res, err := collection.Insert(m)
	if err != nil {
		return 0, err
	}

	id := getInsertId(res.ID())

	return id, nil
}
