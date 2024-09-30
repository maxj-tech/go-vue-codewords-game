package domain

import "math/rand/v2"

type Begriff string

// todo make configurable
var begriffsVorrat = []Begriff{
	"Akt", "Affe", "Auge", "Axt", "Angriff", "Ausweichen", "Anker", "Appetit", "App", "Ampel", "Ampere", "Anhalter", "Aal",
	"Balkan", "Baum", "Blatt", "Blau", "Blume", "Blumenkohl", "Blut", "Blutdruck", "Blutegel", "Braun", "Brause", "Bronze",
	"Blutgruppe",
	"China", "Chaos", "Chor", "Cyan", "Calcium", "Creme", "Creme brulee", "Creme fraiche", "Creme de la creme", "Chlor",
	"Division", "Druck", "David", "Dach", "Dauer", "Duft", "Dunkel", "Dank", "Derby", "Delfin", "Diode", "Dose", "Draht",
	"Ei", "Elefant", "Eis", "Eisbär", "Eisberg", "Eiche", "Eidechse", "Eierbecher", "Ebenholz", "Erde", "Ehe", "Ehre",
	"Fleck", "Fliege",
	"Giraffe", "Gold",
	"Hahn", "Hahnenschrei", "Hamburger", "Hund",
	"Irrlicht", "Irland",
	"Japan", "Jäger",
	"Karte", "Kartoffel", "Kartoffelsalat", "Kartoffelbrei", "Katze", "Kerze", "Käse", "Käsebrot", "Käsekuchen", "Käsespätzle", "König",
	"Königskrone", "Königreich", "Kran", "Kranich", "Kranz", "Krake", "Kreide", "Kreis", "Kreisfahrt", "Kreuz",
	"Kreuzheben", "Kreuzritter", "Kreuzung", "Kreuzworträtsel", "Kreisverkehr", "Kreislauf", "Rezension",
	"Luxemburg", "Luchs", "Lachs", "Lachen", "Lack", "Leckstein", "Leder", "Leber", "Luder", "Löwe", "Laudatio", "Licht",
	"Marker", "Maus", "Melone", "Meter", "Moskau", "Mutter",
	"Nacht", "Natur", "Nase", "Nebel", "Nebel", "Nebenwirkung", "Neuheit", "Nichts", "Nautik", "Nordpol", "Nordsee",
	"Olymp", "Optik", "Ort", "Ostern", "Ostsee", "Ozean", "Ozelot", "Obst", "Oberfläche", "Original", "Ordnung", "Orange",
	"Python", "Puck", "Petri Heil", "Pfanne", "Pfau", "Pfauenauge", "Papier", "Politur", "Pinsel", "Polizei", "Pflanze",
	"Quadrat", "Quadratur", "Quelle", "Qualle", "Quark", "Quarz", "Querflöte", "Queren", "Qualität", "Qualifikation",
	"Rezept", "Rezension", "Rezession", "Riese", "Riesenrad", "Roulette", "Rute", "Ruhe", "Ruhestand", "Rache", "Rasen",
	"Schiff", "Schlange", "Schlangenbiss", "Schlangenlinie", "Schlüssel", "Schloss", "Schuh", "Schuhkarton", "Schuhcreme",
	"Sohle", "Suche", "Suchmaschine", "Steuer", "Skelett", "Skala", "Strauß", "Strand", "Stift", "Stuhl", "Sex", "Sonne",
	"Toast", "Tradition", "Tisch", "Tischdecke", "Tischtennis", "Tischler", "Tischlerei", "Tabledance", "Thermometer",
	"Uhr", "Uhu", "Urahn", "Uran", "Ureinwohner", "Ulme", "Ufer", "Übermut", "Überfluss", "Umsatz", "Umschlag", "Umwelt",
	"Vater", "Verzeichnis", "Verehrer", "Verwandtschaft", "Vogel", "Vogelhaus", "Verein", "Vereinigung", "Vereinbarung",
	"Wand", "Wasser", "Weide", "Wald", "Wiese", "Wurzel", "Wermut", "Wolke", "Welle", "Wichtig", "Wicht", "Warze", "Wurst",
	"Zyklop", "Zyklon", "Zylinder", "Zeiger", "Zacken", "Zahn", "Zeit", "Zaun", "Zauber", "Zenith", "Zentrum", "Zebra",
}

func waehleRandomBegriffe(n int) []Begriff {
	randomBegriffe := make([]Begriff, len(begriffsVorrat))
	copy(randomBegriffe, begriffsVorrat)
	rand.Shuffle(len(randomBegriffe), func(i, j int) {
		randomBegriffe[i], randomBegriffe[j] = randomBegriffe[j], randomBegriffe[i]
	})
	return randomBegriffe[:n]
}
