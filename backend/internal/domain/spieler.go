package domain

import (
	"errors"
	"github.com/google/uuid"
	"strings"
)

type Rolle int

const (
	Chef Rolle = iota
	Ermittler
)

type Spieler struct {
	id    uuid.UUID
	Name  string
	Rolle Rolle
}

func NewSpieler(name string, rolle Rolle) (Spieler, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Spieler{}, errors.New("name cannot be empty or whitespace")
	}

	return Spieler{
		id:    uuid.New(),
		Name:  name,
		Rolle: rolle,
	}, nil
}
