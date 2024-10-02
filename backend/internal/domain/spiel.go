package domain

import "fmt"

type Spiel struct {
	Karten            []Karte
	Teams             []Team
	AlsNaechstesAmZug Teamfarbe
	Spielstand        map[Teamfarbe]int
}

func NewSpiel(startTeam Teamfarbe) Spiel {
	gewaehlteBegriffe := waehleRandomBegriffe(25)
	var roteKarten, blaueKarten int
	if startTeam == TeamRot {
		roteKarten = 9
		blaueKarten = 8
	} else {
		roteKarten = 8
		blaueKarten = 9
	}
	gewaehlteFarben := erzeugeUndMischeFarben(roteKarten, blaueKarten)

	karten := erzeugeKarten(gewaehlteBegriffe, gewaehlteFarben)

	return Spiel{
		Karten:            karten,
		AlsNaechstesAmZug: startTeam,
		Spielstand:        map[Teamfarbe]int{TeamRot: 0, TeamBlau: 0},
	}
}

func (s *Spiel) SetTeams(team1, team2 Team) error {
	if team1.Farbe == team2.Farbe {
		return fmt.Errorf("Die Teams dürfen nicht die gleiche Farbe haben")
	}
	if team1.Chef.ID() == team2.Chef.ID() {
		return fmt.Errorf("Die Chefs der Teams dürfen nicht die gleichen sein")
	}
	for id := range team1.Ermittler {
		if _, exists := team2.Ermittler[id]; exists {
			return fmt.Errorf("Kein Ermittler darf in beiden Teams sein.")
		}
	}
	s.Teams = []Team{team1, team2}
	return nil
}
