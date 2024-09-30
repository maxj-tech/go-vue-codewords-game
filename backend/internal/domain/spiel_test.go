package domain

import "testing"

func TestNeuesSpielHat25Karten(t *testing.T) {
	spiel := NeuesSpiel(TeamRot)

	if len(spiel.Karten) != 25 {
		t.Fatalf("expected 25 cards, got %d", len(spiel.Karten))
	}
}

func TestNeuesSpielKeineDoppeltenBegriffe(t *testing.T) {
	spiel := NeuesSpiel(TeamBlau)

	begriffMap := make(map[string]bool)
	for _, karte := range spiel.Karten {
		if begriffMap[karte.Begriff] {
			t.Fatalf("duplicate card found: %s", karte.Begriff)
		}
		begriffMap[karte.Begriff] = true
	}
}

func TestNeuesSpielFarbenVerteilung(t *testing.T) {
	tests := []struct {
		startendesTeam Teamfarbe
		expectedRote   int
		expectedBlaue  int
	}{
		{TeamRot, 9, 8},
		{TeamBlau, 8, 9},
	}

	for _, tt := range tests {
		spiel := NeuesSpiel(tt.startendesTeam)

		countColors := map[Kartenfarbe]int{
			KarteRot:     0,
			KarteBlau:    0,
			KarteBeige:   0,
			KarteSchwarz: 0,
		}

		for _, karte := range spiel.Karten {
			countColors[karte.Farbe]++
		}

		if countColors[KarteRot] != tt.expectedRote {
			t.Errorf("expected %d red cards for starting team %v, got %d", tt.expectedRote, tt.startendesTeam, countColors[KarteRot])
		}

		if countColors[KarteBlau] != tt.expectedBlaue {
			t.Errorf("expected %d blue cards for starting team %v, got %d", tt.expectedBlaue, tt.startendesTeam, countColors[KarteBlau])
		}

		if countColors[KarteBeige] != 7 {
			t.Errorf("expected 7 beige cards for starting team %v, got %d", tt.startendesTeam, countColors[KarteBeige])
		}

		if countColors[KarteSchwarz] != 1 {
			t.Errorf("expected 1 black card for starting team %v, got %d", tt.startendesTeam, countColors[KarteSchwarz])
		}

	}
}

func TestNeuesSpielAktuellerZugFuerStartendesTeam(t *testing.T) {
	var startendesTeam Teamfarbe = TeamBlau
	spiel := NeuesSpiel(startendesTeam)
	if spiel.AktuellerZug.Farbe != startendesTeam {
		t.Errorf("expected current turn to be with team %v but got team %v", startendesTeam, spiel.AktuellerZug.Farbe)
	}
}

func TestNeuesSpielStartendesTeamRichtigGesetzt(t *testing.T) {
	var startendesTeam Teamfarbe = TeamRot
	spiel := NeuesSpiel(startendesTeam)

	if spiel.StartendesTeam != startendesTeam {
		t.Errorf("expected starting team to be %v but got %v", startendesTeam, spiel.StartendesTeam)
	}
}

func TestNeuesSpielInitialerSpielstand(t *testing.T) {
	spiel := NeuesSpiel(TeamRot)

	if spiel.Spielstand[TeamRot] != 0 || spiel.Spielstand[TeamBlau] != 0 {
		t.Errorf("expected initial scores to be zero but got Rot:%d Blau:%d", spiel.Spielstand[TeamRot], spiel.Spielstand[TeamBlau])
	}
}

func TestNeuesSpielTeamsInitialisiert(t *testing.T) {
	spiel := NeuesSpiel(TeamRot)

	if len(spiel.Teams) != 2 || spiel.Teams[0].Farbe != TeamRot || spiel.Teams[1].Farbe != TeamBlau {
		t.Errorf("teams not initialized correctly")
	}
}
