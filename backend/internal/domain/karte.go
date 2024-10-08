package domain

import "math/rand/v2"

type Kartenfarbe int

const (
	KarteRot     = Kartenfarbe(Rot)
	KarteBlau    = Kartenfarbe(Blau)
	KarteBeige   = Kartenfarbe(Beige)
	KarteSchwarz = Kartenfarbe(Schwarz)
)

type Karte struct {
	Begriff string
	Getippt bool
	Farbe   Kartenfarbe
}

func erzeugeUndMischeFarben(roteKarten int, blaueKarten int) []Kartenfarbe {
	var farben []Kartenfarbe
	farben = append(farben, wiederholeFarbe(KarteRot, roteKarten)...)
	farben = append(farben, wiederholeFarbe(KarteBlau, blaueKarten)...)

	farben = append(farben, wiederholeFarbe(KarteBeige, 7)...)
	farben = append(farben, KarteSchwarz)

	rand.Shuffle(len(farben), func(i, j int) { farben[i], farben[j] = farben[j], farben[i] })

	return farben
}

func wiederholeFarbe(farbe Kartenfarbe, anzahl int) []Kartenfarbe {
	result := make([]Kartenfarbe, anzahl)
	for i := range result {
		result[i] = farbe
	}
	return result
}

func erzeugeKarten(begriffe []Begriff, farben []Kartenfarbe) []Karte {
	var karten []Karte
	for i, begriff := range begriffe {
		karten = append(karten, Karte{Begriff: string(begriff), Farbe: farben[i]})
	}
	return karten
}
