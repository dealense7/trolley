package country

type Id int64

const (
	Spain Id = 1
	Italy Id = 2
)

type Model struct {
	Id     Id     `db:"id"`
	Code   string `db:"code"`
	Name   string `db:"name"`
	Status bool   `db:"status"`
}
