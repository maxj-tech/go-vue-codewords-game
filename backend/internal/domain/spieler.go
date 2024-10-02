package domain

import (
	"errors"
	"fmt"
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
		return Spieler{}, errors.New("Ungültiger Name")
	}

	if rolle != Chef && rolle != Ermittler {
		return Spieler{}, errors.New("Ungültige Rolle")
	}

	return Spieler{
		id:    uuid.New(),
		Name:  name,
		Rolle: rolle,
	}, nil
}

func (s Spieler) ID() uuid.UUID {
	return s.id
}

func (s Spieler) String() string {
	return fmt.Sprintf("%s (ID: %s, Rolle: %d)", s.Name, s.id, s.Rolle)
}
