package domain

import "github.com/google/uuid"

type Rolle int

const (
	Chef Rolle = iota
	Ermittler
)

type Spieler struct {
	Id    uuid.UUID
	Name  string
	Rolle Rolle
}

func NewSpieler(name string, rolle Rolle) Spieler {
	return Spieler{
		Id:    uuid.New(),
		Name:  name,
		Rolle: rolle,
	}
}
