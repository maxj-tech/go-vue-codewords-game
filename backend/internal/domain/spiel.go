package domain

type Spiel struct {
	Karten         []Karte
	Teams          []Team
	StartendesTeam Teamfarbe
	AktuellerZug   Team
	Spielstand     map[Teamfarbe]int // todo why could not get to work: map[Team]int ?
}

func NeuesSpiel(startTeam Teamfarbe) *Spiel {
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

	return &Spiel{
		Karten:         karten,
		Teams:          []Team{{Farbe: TeamRot}, {Farbe: TeamBlau}},
		StartendesTeam: startTeam,
		AktuellerZug:   Team{Farbe: startTeam},
		Spielstand:     map[Teamfarbe]int{TeamRot: 0, TeamBlau: 0},
	}
}
