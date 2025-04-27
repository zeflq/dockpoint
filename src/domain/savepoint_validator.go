package domain

type SavepointValidator interface {
	Validate(savepoints []Savepoint) error
}
