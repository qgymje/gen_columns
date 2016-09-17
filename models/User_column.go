package models

type UserColumn struct {
	ID   string
	Name string
}

var UserColumns UserColumn

func init() {
	UserColumns.ID = "id"
	UserColumns.Name = "name"

}
