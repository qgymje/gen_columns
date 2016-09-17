package models

type _UserColumn struct {
	ID   string
	Name string
}

var UserColumns _UserColumn

func init() {
	UserColumns.ID = "id"
	UserColumns.Name = "name"

}
