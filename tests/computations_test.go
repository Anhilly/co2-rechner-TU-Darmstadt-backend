package tests

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"io"
	"testing"
)

func TestCompuations(t *testing.T) {
	database.ConnectDatabase()
	defer database.DisconnectDatabase()

	t.Run("TestBerechneITGeraete", TestBerechneITGeraete)
	t.Run("TestBerechnePendelweg", TestBerechnePendelweg)
}

func TestTester(t *testing.T) {
	database.ConnectDatabase()
	defer database.DisconnectDatabase()

	t.Run("TestBerechnePendelweg", TestBerechnePendelweg)
}

func TestBerechneITGeraete(t *testing.T) {
	is := is.NewRelaxed(t)

	// normale Berechnungen
	t.Run("BerechneITGeraete: leere Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var itGeraeteDaten []structs.ITGeraeteAnzahl = nil

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)             // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // bei leerer Eingabe gibt es keine Emissionen
	})

	t.Run("BerechneITGeraete: leere Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)             // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // bei leerer Eingabe gibt es keine Emissionen
	})

	t.Run("BerechneITGeraete: einelementige Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{{IDITGeraete: 1, Anzahl: 1}}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)                  // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 147000.0) // CO2KatorJahr (588000) fuer Notebooks (ID = 1)
	})

	t.Run("BerechneITGeraete: komplexe Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{
			{IDITGeraete: 1, Anzahl: 5},
			{IDITGeraete: 4, Anzahl: 10},
			{IDITGeraete: 6, Anzahl: 4},
			{IDITGeraete: 3, Anzahl: 3},
			{IDITGeraete: 7, Anzahl: 1},
			{IDITGeraete: 2, Anzahl: 4},
		}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)                   // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 1594235.0) // Erwarteter Wert: 1594235
	})

	t.Run("BerechneITGeraete: Anzahl 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{
			{IDITGeraete: 1, Anzahl: 0},
			{IDITGeraete: 4, Anzahl: 0},
			{IDITGeraete: 6, Anzahl: 0},
			{IDITGeraete: 3, Anzahl: 0},
			{IDITGeraete: 7, Anzahl: 0},
			{IDITGeraete: 2, Anzahl: 0},
		}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)             // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // Erwarteter Wert: 0.0
	})

	t.Run("BerechneITGeraete: Toner", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{
			{IDITGeraete: 8, Anzahl: 1},
			{IDITGeraete: 10, Anzahl: 1},
		}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)                 // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 27000.0) // Erwarteter Wert: 27000.0
	})

	// Fehlertests
	t.Run("BerechneITGeraete: ID nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{
			{IDITGeraete: 100, Anzahl: 5},
		}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.Equal(err, io.EOF)     // EOF von der Dantenbank erwartet
		is.Equal(emissionen, 0.0) // Erwarteter Wert: 0.0
	})
}

func TestBerechnePendelweg(t *testing.T) {
	is := is.NewRelaxed(t)

	// normale Berechnungen
	t.Run("BerechnePendelweg: Slice = nil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pendelwegDaten []structs.PendelwegElement = nil
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)             // soll ErrPersonenzahlZuKlein werfen
		is.Equal(emissionen, 0.0) // bei Fehler ist Ergebnis 0.0
	})

	t.Run("BerechnePendelweg: leerer Slice", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{}
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)             // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // bei leerer Eingabe gibt es keine Emissionen
	})

	t.Run("BerechnePendelweg: tageImBuero = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{
			{IDPendelweg: 1, Strecke: 1, Personenanzahl: 1},
		}
		var tageImBuero int32 = 0

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)             // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // bei leerer Eingabe gibt es keine Emissionen
	})

	t.Run("BerechnePendelweg: leerer Slice, tageImBuero = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{}
		var tageImBuero int32 = 0

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)             // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // bei leerer Eingabe gibt es keine Emissionen
	})

	t.Run("BerechnePendelweg: einfache Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{
			{IDPendelweg: 1, Strecke: 10, Personenanzahl: 1},
		}
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)                // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 3680.0) // erwartetes Ergebnis: 3680.0
	})

	t.Run("BerechnePendelweg: komplexe Berechnung", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{
			{IDPendelweg: 1, Strecke: 10, Personenanzahl: 1},
			{IDPendelweg: 4, Strecke: 15, Personenanzahl: 1},
			{IDPendelweg: 10, Strecke: 90, Personenanzahl: 1},
			{IDPendelweg: 7, Strecke: 3, Personenanzahl: 1},
			{IDPendelweg: 5, Strecke: 5, Personenanzahl: 1},
		}
		var tageImBuero int32 = 3

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)                   // bei normaler Berechnung sollte kein Error geworfen werden
		is.Equal(emissionen, 2181504.0) // erwartetes Ergebnis: 2181504.0
	})

	// Errortests
	t.Run("BerechnePendelweg: ID = 100", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{
			{IDPendelweg: 100, Strecke: 1, Personenanzahl: 1},
		}
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.Equal(err, io.EOF)     // Datenbank gibt EOF Error zurueck, wenn ID unbekannt
		is.Equal(emissionen, 0.0) // bei Fehler ist Ergebnis 0.0
	})

	t.Run("BerechnePendelweg: Personenanzahl < 1", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{
			{IDPendelweg: 1, Strecke: 1, Personenanzahl: 0},
		}
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.Equal(err, co2computation.ErrPersonenzahlZuKlein) // soll ErrPersonenzahlZuKlein werfen
		is.Equal(emissionen, 0.0)                            // bei Fehler ist Ergebnis 0.0
	})
}