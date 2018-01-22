package data_layer

type Entity struct {
	Id int
}

type IEntity interface {
	Save() error
	Reload() error
}

func (e Entity) Save() error {

	// try saving the entity into our database
	return nil
}
