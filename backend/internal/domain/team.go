package domain

import (
	"errors"
	"github.com/google/uuid"
)

type Teamfarbe int

const (
	TeamRot  = Teamfarbe(Rot)
	TeamBlau = Teamfarbe(Blau)
)

type Team struct {
	Farbe     Teamfarbe
	Ermittler map[uuid.UUID]Spieler
	Chef      Spieler
}

func NewTeam(farbe Teamfarbe, chef Spieler, ermittler []Spieler) (Team, error) {

	if chef.Rolle != Chef {
		return Team{}, errors.New("chef must have role Chef")
	}

	if len(ermittler) == 0 {
		return Team{}, errors.New("at least one ermittler must be set")
	}

	ermittlerSet := make(map[uuid.UUID]Spieler)
	for _, e := range ermittler {
		if e.Rolle != Ermittler {
			return Team{}, errors.New("all ermittler must have role Ermittler")
		}
		if _, exists := ermittlerSet[e.id]; exists {
			return Team{}, errors.New("duplicate ermittler found")
		}
		ermittlerSet[e.id] = e
	}

	return Team{
		Farbe:     farbe,
		Chef:      chef,
		Ermittler: ermittlerSet,
	}, nil
}
