package tests

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"io"
	"testing"
)

func TestComputations(t *testing.T) {
	database.ConnectDatabase()
	defer database.DisconnectDatabase()

	t.Run("TestBerechneITGeraete", TestBerechneITGeraete)
	t.Run("TestBerechnePendelweg", TestBerechnePendelweg)
	t.Run("TestBerechneDienstreisen", TestBerechneDienstreisen)
}

func TestBerechneITGeraete(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("BerechneITGeraete: Slice = nil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var itGeraeteDaten []structs.ITGeraeteAnzahl = nil

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechneITGeraete: leere Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechneITGeraete: einelementige Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{{IDITGeraete: 1, Anzahl: 1}}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)                  // Normalfall wirft keine Errors
		is.Equal(emissionen, 147000.0) // erwartetes Ergebnis: 147000.0 (CO2FaktorJahr fuer Notebooks (ID = 1))
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

		is.NoErr(err)                   // Normalfall wirft keine Errors
		is.Equal(emissionen, 1594235.0) // Erwarteter Wert: 1594235.0
	})

	t.Run("BerechneITGeraete: Anzahl 0 eingegeben", func(t *testing.T) {
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

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // Erwarteter Wert: 0.0 (kein Anzahl = keine Emissionen)
	})

	t.Run("BerechneITGeraete: Toner", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{
			{IDITGeraete: 8, Anzahl: 1},
			{IDITGeraete: 10, Anzahl: 1},
		}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(emissionen, 27000.0) // Erwarteter Wert: 27000.0
	})

	// Fehlertests
	t.Run("BerechneITGeraete: ID = 100 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{
			{IDITGeraete: 100, Anzahl: 5},
		}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.Equal(err, io.EOF)     // Datenbank wirft EOF
		is.Equal(emissionen, 0.0) // Fehlerfall liefert 0.0
	})

	t.Run("BerechneITGeraete: negative Anzahl eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.ITGeraeteAnzahl{
			{IDITGeraete: 1, Anzahl: -5},
		}

		emissionen, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.Equal(err, co2computation.ErrAnzahlNegativ) // Funktion wirft ErrAnzahlNegativ
		is.Equal(emissionen, 0.0)                      // Fehlerfall liefert 0.0
	})

	// Fehler ErrStrEinheitUnbekannt momentan nicht abpruefbar. Benoetigt falschen Datensatz in Datenbank
}

func TestBerechnePendelweg(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("BerechnePendelweg: Slice = nil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pendelwegDaten []structs.PendelwegElement = nil
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechnePendelweg: leerer Slice", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{}
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechnePendelweg: tageImBuero = 0 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{
			{IDPendelweg: 1, Strecke: 1, Personenanzahl: 1},
		}
		var tageImBuero int32 = 0

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (kein Pendelweg = keine Emissionen)
	})

	t.Run("BerechnePendelweg: leerer Slice, tageImBuero = 0 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{}
		var tageImBuero int32 = 0

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechnePendelweg: einfache Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{
			{IDPendelweg: 1, Strecke: 10, Personenanzahl: 1},
		}
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.NoErr(err)                // Normalfall wirft keine Errors
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

		is.NoErr(err)                   // Normalfall wirft keine Errors
		is.Equal(emissionen, 2181504.0) // erwartetes Ergebnis: 2181504.0
	})

	// Errortests
	t.Run("BerechnePendelweg: ID = 100 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{
			{IDPendelweg: 100, Strecke: 1, Personenanzahl: 1},
		}
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.Equal(err, io.EOF)     // Datenbank wirft EOF Error
		is.Equal(emissionen, 0.0) // Fehlerfall liefert 0.0
	})

	t.Run("BerechnePendelweg: Personenanzahl < 1 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{
			{IDPendelweg: 1, Strecke: 1, Personenanzahl: 0},
		}
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.Equal(err, co2computation.ErrPersonenzahlZuKlein) // Funktion wirft ErrPersonenzahlZuKlein
		is.Equal(emissionen, 0.0)                            // Fehlerfall liefert 0.0
	})

	t.Run("BerechnePendelweg: negative Strecke eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.PendelwegElement{
			{IDPendelweg: 1, Strecke: -100, Personenanzahl: 1},
		}
		var tageImBuero int32 = 1

		emissionen, err := co2computation.BerechnePendelweg(pendelwegDaten, tageImBuero)

		is.Equal(err, co2computation.ErrStreckeNegativ) // Funktion wirft ErrStreckeNegativ
		is.Equal(emissionen, 0.0)                       // Fehlerfall liefert 0.0
	})

	// Fehler ErrStrEinheitUnbekannt momentan nicht abpruefbar. Benoetigt falschen Datensatz in Datenbank
}

func TestBerechneDienstreisen(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("BerechneDienstreisen: Slice = nil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var dienstreisenDaten []structs.DienstreiseElement = nil

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechneDienstreisen: leerer Slice", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
	})

	t.Run("BerechneDienstreisen: einfache Eingabe Bahn", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 1, Strecke: 100},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(emissionen, 1600.0) // erwartetes Ergebnis: 1600.0
	})

	t.Run("BerechneDienstreisen: Bahn; Felder Tankart, Streckentyp egal", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 1, Strecke: 10, Streckentyp: "unbekannt", Tankart: "unbekannt"},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)               // Normalfall wirft keine Errors
		is.Equal(emissionen, 160.0) // erwartetes Ergebnis: 160.0
	})

	t.Run("BerechneDienstreisen: einfache Eingabe Auto", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 2, Strecke: 100, Tankart: "Diesel"},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(emissionen, 48800.0) // erwartetes Ergebnis: 48800.0
	})

	t.Run("BerechneDienstreisen: Auto; Feld Streckentyp egal", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 2, Strecke: 10, Tankart: "Benzin", Streckentyp: "unbekannt"},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(emissionen, 5200.0) // erwartetes Ergebnis: 5200.0
	})

	t.Run("BerechneDienstreisen: einfache Eingabe Flugzeug", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 3, Strecke: 100, Streckentyp: "Kurzstrecke"},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                  // Normalfall wirft keine Errors
		is.Equal(emissionen, 177600.0) // erwartetes Ergebnis: 177600.0
	})

	t.Run("BerechneDienstreisen: Flugzeug; Feld Tankart egal", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 3, Strecke: 10, Streckentyp: "Langstrecke", Tankart: "unbekannt"},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(emissionen, 8360.0) // erwartetes Ergebnis: 8360.0
	})

	t.Run("BerechneDienstreisen: komplexe Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 2, Strecke: 1200, Tankart: "Diesel"},
			{IDDienstreise: 1, Strecke: 150},
			{IDDienstreise: 3, Strecke: 750, Streckentyp: "Langstrecke"},
			{IDDienstreise: 3, Strecke: 1000, Streckentyp: "Kurzstrecke"},
			{IDDienstreise: 2, Strecke: 45, Tankart: "Benzin"},
			{IDDienstreise: 1, Strecke: 1},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                   // Normalfall wirft keine Errors
		is.Equal(emissionen, 3014416.0) // erwartetes Ergebnis: 3014416.0
	})

	t.Run("BerechneDienstreisen: Strecke = 0 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 1, Strecke: 0},
			{IDDienstreise: 2, Strecke: 0, Tankart: "Benzin"},
			{IDDienstreise: 2, Strecke: 0, Tankart: "Diesel"},
			{IDDienstreise: 3, Strecke: 0, Streckentyp: "Kurzstrecke"},
			{IDDienstreise: 3, Strecke: 0, Streckentyp: "Langstrecke"},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (keine Strecke = keine Emissionen)
	})

	// Errortests
	t.Run("BerechneDienstreisen: ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 0, Strecke: 100},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.Equal(err, io.EOF)     // Datenbank wirft EOF Error
		is.Equal(emissionen, 0.0) // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
	})

	t.Run("BerechneDienstreisen: unbekannte Tankart", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 2, Strecke: 100, Tankart: "nicht definiert"},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.Equal(err, co2computation.ErrTankartUnbekannt) // Funktione wirft ErrTankartUnbekannt
		is.Equal(emissionen, 0.0)                         // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
	})

	t.Run("BerechneDienstreisen: unbekannter Streckentyp", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 3, Strecke: 100, Streckentyp: "nicht definiert"},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.Equal(err, co2computation.ErrStreckentypUnbekannt) // Funktione wirft ErrStreckentypUnbekannt
		is.Equal(emissionen, 0.0)                             // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
	})

	t.Run("BerechneDienstreisen: negative Strecke eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.DienstreiseElement{
			{IDDienstreise: 1, Strecke: -100},
		}

		emissionen, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.Equal(err, co2computation.ErrStreckeNegativ) // Funktione wirft ErrStreckeNegativ)
		is.Equal(emissionen, 0.0)                       // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
	})

	// Fehler ErrStrEinheitUnbekannt momentan nicht abpruefbar. Benoetigt falschen Datensatz in Datenbank
}
