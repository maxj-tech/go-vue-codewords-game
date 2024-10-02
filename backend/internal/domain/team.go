package domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type Teamfarbe int

const (
	TeamRot  = Teamfarbe(Rot)
	TeamBlau = Teamfarbe(Blau)
)

func (c Teamfarbe) String() string {
	switch c {
	case TeamRot:
		return "TeamRot(0)"
	case TeamBlau:
		return "TeamBlau(1)"
	default:
		return fmt.Sprintf("Unbekannte Teamfarbe: %d", c)
	}
}

type Team struct {
	Farbe     Teamfarbe
	Ermittler map[uuid.UUID]Spieler
	Chef      Spieler
}

func NewTeam(farbe Teamfarbe, chef Spieler, ermittler []Spieler) (Team, error) {

	if farbe != TeamRot && farbe != TeamBlau {
		return Team{}, errors.New("Ungültige Teamfarbe.")
	}

	if chef.Rolle != Chef {
		return Team{}, errors.New("Chef muss auch Chef-Rolle haben")
	}

	if len(ermittler) == 0 {
		return Team{}, errors.New("Es muss mind. ein Ermittler da sein")
	}

	ermittlerSet := make(map[uuid.UUID]Spieler)
	for _, e := range ermittler {
		if e.Rolle != Ermittler {
			return Team{}, errors.New("Alle Ermittler müssen die Ermittler-Rolle haben.")
		}
		if _, exists := ermittlerSet[e.id]; exists {
			return Team{}, errors.New("Kein Ermittler darf doppelt vorkommen.")
		}
		ermittlerSet[e.id] = e
	}

	return Team{
		Farbe:     farbe,
		Chef:      chef,
		Ermittler: ermittlerSet,
	}, nil
}

func (t Team) String() string {
	return fmt.Sprintf("Team: %v, Chef: %s, Spieler: %v", t.Farbe, t.Chef, t.Ermittler)
}
