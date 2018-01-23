package data_layer

type Entity struct {
	Id int
}

type IEntity interface {
	Save() error
	Reload() error
}

func CreateEntity() Entity {
	return Entity{}
}

func (e Entity) Save() error {

	// try saving the entity into our database
	// FIXME: actually save into database
	return nil
}
