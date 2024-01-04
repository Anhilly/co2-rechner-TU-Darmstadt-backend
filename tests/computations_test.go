package tests

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestComputations(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase("dev")
	is.NoErr(err)
	defer func() {
		err := database.DisconnectDatabase()
		is.NoErr(err)
	}()

	t.Run("TestBerechneITGeraete", TestBerechneITGeraete)
	t.Run("TestBerechnePendelweg", TestBerechnePendelweg)
	t.Run("TestBerechneDienstreisen", TestBerechneDienstreisen)
}

func TestBerechneITGeraete(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("BerechneITGeraete: Slice = nil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var itGeraeteDaten []structs.UmfrageITGeraete = nil

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
		is.Equal(emissionenAufgeteilt, map[int32]float64{})
	})

	t.Run("BerechneITGeraete: leere Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.UmfrageITGeraete{}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
		is.Equal(emissionenAufgeteilt, map[int32]float64{})
	})

	t.Run("BerechneITGeraete: einelementige Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.UmfrageITGeraete{{IDITGeraete: 1, Anzahl: 1}}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)                  // Normalfall wirft keine Errors
		is.Equal(emissionen, 147000.0) // erwartetes Ergebnis: 147000.0 (CO2FaktorJahr fuer Notebooks (ID = 1))
		is.Equal(emissionenAufgeteilt, map[int32]float64{
			1: 147000.0,
		})
	})

	t.Run("BerechneITGeraete: komplexe Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.UmfrageITGeraete{
			{IDITGeraete: 1, Anzahl: 5},
			{IDITGeraete: 4, Anzahl: 7},
			{IDITGeraete: 6, Anzahl: 4},
			{IDITGeraete: 3, Anzahl: 3},
			{IDITGeraete: 7, Anzahl: 1},
			{IDITGeraete: 2, Anzahl: 4},
			{IDITGeraete: 4, Anzahl: 3},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)                   // Normalfall wirft keine Errors
		is.Equal(emissionen, 1594235.0) // Erwarteter Wert: 1594235.0
		is.Equal(emissionenAufgeteilt, map[int32]float64{
			1: 735000.0,
			2: 237400.0,
			3: 220500.0,
			7: 74625.0,
			4: 89710.0,
			6: 237000.0,
		})
	})

	t.Run("BerechneITGeraete: Anzahl 0 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.UmfrageITGeraete{
			{IDITGeraete: 1, Anzahl: 0},
			{IDITGeraete: 4, Anzahl: 0},
			{IDITGeraete: 6, Anzahl: 0},
			{IDITGeraete: 3, Anzahl: 0},
			{IDITGeraete: 7, Anzahl: 0},
			{IDITGeraete: 2, Anzahl: 0},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // Erwarteter Wert: 0.0 (kein Anzahl = keine Emissionen)
		is.Equal(emissionenAufgeteilt, map[int32]float64{})
	})

	t.Run("BerechneITGeraete: Toner", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.UmfrageITGeraete{
			{IDITGeraete: 8, Anzahl: 1},
			{IDITGeraete: 10, Anzahl: 1},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(emissionen, 27000.0) // Erwarteter Wert: 27000.0
		is.Equal(emissionenAufgeteilt, map[int32]float64{
			8:  13500.0,
			10: 13500.0,
		})
	})

	// Fehlertests
	t.Run("BerechneITGeraete: ID = 100 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.UmfrageITGeraete{
			{IDITGeraete: 100, Anzahl: 5},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(emissionen, 0.0)           // Fehlerfall liefert 0.0
		is.Equal(emissionenAufgeteilt, nil)
	})

	t.Run("BerechneITGeraete: negative Anzahl eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		itGeraeteDaten := []structs.UmfrageITGeraete{
			{IDITGeraete: 1, Anzahl: -5},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneITGeraete(itGeraeteDaten)

		is.Equal(err, structs.ErrAnzahlNegativ) // Funktion wirft ErrAnzahlNegativ
		is.Equal(emissionen, 0.0)               // Fehlerfall liefert 0.0
		is.Equal(emissionenAufgeteilt, nil)
	})

	// Fehler ErrStrEinheitUnbekannt momentan nicht abpruefbar. Benoetigt falschen Datensatz in Datenbank
}

func TestBerechnePendelweg(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("BerechnePendelweg: Slice = nil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var pendelwegDaten []structs.UmfragePendelweg = nil
		var tageImBuero int32 = 1
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
		is.Equal(emissionenAufgeteilt, map[int32]float64{})
	})

	t.Run("BerechnePendelweg: leerer Slice", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{}
		var tageImBuero int32 = 1

		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
		is.Equal(emissionenAufgeteilt, map[int32]float64{})
	})

	t.Run("BerechnePendelweg: tageImBuero = 0 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{
			{IDPendelweg: 1, Strecke: 1, Personenanzahl: 1},
		}
		var tageImBuero int32 = 0
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (kein Pendelweg = keine Emissionen)
		is.Equal(emissionenAufgeteilt, map[int32]float64{})
	})

	t.Run("BerechnePendelweg: leerer Slice, tageImBuero = 0 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{}
		var tageImBuero int32 = 0
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
		is.Equal(emissionenAufgeteilt, map[int32]float64{})
	})

	t.Run("BerechnePendelweg: einfache Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{
			{IDPendelweg: 1, Strecke: 10, Personenanzahl: 1},
		}
		var tageImBuero int32 = 1
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(emissionen, 3680.0) // erwartetes Ergebnis: 3680.0
		is.Equal(emissionenAufgeteilt, map[int32]float64{
			1: 3680.0,
		})
	})

	t.Run("BerechnePendelweg: Fussgaenger", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{
			{IDPendelweg: 12, Strecke: 20, Personenanzahl: 1},
		}
		var tageImBuero int32 = 3
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.NoErr(err)             //Normalfall wirft keinen Error
		is.Equal(emissionen, 0.0) // erwartet Ergebnis 0, da Fussgaenger
		is.Equal(emissionenAufgeteilt, map[int32]float64{
			12: 0.0,
		})
	})

	t.Run("BerechnePendelweg: Eingabe mit Weg == 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{
			{IDPendelweg: 2, Strecke: 0, Personenanzahl: 1},
			{IDPendelweg: 1, Strecke: 10, Personenanzahl: 1},
		}
		var tageImBuero int32 = 1
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(emissionen, 3680.0) // erwartetes Ergebnis: 3680.0
		is.Equal(emissionenAufgeteilt, map[int32]float64{
			1: 3680.0,
		})
	})

	t.Run("BerechnePendelweg: Rundung auf 2 Nachkommastellen", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{
			{IDPendelweg: 1, Strecke: 10, Personenanzahl: 3},
		}
		var tageImBuero int32 = 1
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(emissionen, 1226.67) // erwartetes Ergebnis: 1226.67
		is.Equal(emissionenAufgeteilt, map[int32]float64{
			1: 1226.67,
		})
	})

	t.Run("BerechnePendelweg: komplexe Berechnung", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{
			{IDPendelweg: 1, Strecke: 10, Personenanzahl: 1},
			{IDPendelweg: 4, Strecke: 15, Personenanzahl: 1},
			{IDPendelweg: 10, Strecke: 60, Personenanzahl: 1},
			{IDPendelweg: 7, Strecke: 3, Personenanzahl: 1},
			{IDPendelweg: 5, Strecke: 5, Personenanzahl: 1},
			{IDPendelweg: 10, Strecke: 30, Personenanzahl: 1},
		}
		var tageImBuero int32 = 3
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.NoErr(err)                   // Normalfall wirft keine Errors
		is.Equal(emissionen, 2181504.0) // erwartetes Ergebnis: 2181504.0
		is.Equal(emissionenAufgeteilt, map[int32]float64{
			1:  11040,
			4:  1010160,
			5:  358800,
			7:  6624,
			10: 794880,
		})
	})

	// Errortests
	t.Run("BerechnePendelweg: ID = 100 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{
			{IDPendelweg: 100, Strecke: 1, Personenanzahl: 1},
		}
		var tageImBuero int32 = 1
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(emissionen, 0.0)           // Fehlerfall liefert 0.0
		is.Equal(emissionenAufgeteilt, nil) // Fehlerfall liefert nil
	})

	t.Run("BerechnePendelweg: Personenanzahl < 1 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{
			{IDPendelweg: 1, Strecke: 1, Personenanzahl: 0},
		}
		var tageImBuero int32 = 1
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.Equal(err, structs.ErrPersonenzahlZuKlein) // Funktion wirft ErrPersonenzahlZuKlein
		is.Equal(emissionen, 0.0)                     // Fehlerfall liefert 0.0
		is.Equal(emissionenAufgeteilt, nil)           // Fehlerfall liefert nil
	})

	t.Run("BerechnePendelweg: negative Strecke eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		pendelwegDaten := []structs.UmfragePendelweg{
			{IDPendelweg: 1, Strecke: -100, Personenanzahl: 1},
		}
		var tageImBuero int32 = 1
		var allePendelwege = []structs.AllePendelwege{
			{
				Pendelwege:  pendelwegDaten,
				TageImBuero: tageImBuero,
			},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechnePendelweg(allePendelwege)

		is.Equal(err, structs.ErrStreckeNegativ) // Funktion wirft ErrStreckeNegativ
		is.Equal(emissionen, 0.0)                // Fehlerfall liefert 0.0
		is.Equal(emissionenAufgeteilt, nil)      // Fehlerfall liefert nil
	})

	// Fehler ErrStrEinheitUnbekannt momentan nicht abpruefbar. Benoetigt falschen Datensatz in Datenbank
}

func TestBerechneDienstreisen(t *testing.T) { //nolint:funlen
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("BerechneDienstreisen: Slice = nil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var dienstreisenDaten []structs.UmfrageDienstreise = nil

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
		is.Equal(emissionenAufgeteilt, map[string]float64{})
	})

	t.Run("BerechneDienstreisen: leerer Slice", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (bei leerer Eingabe keine Emissionen)
		is.Equal(emissionenAufgeteilt, map[string]float64{})
	})

	t.Run("BerechneDienstreisen: einfache Eingabe Bahn", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 1, Strecke: 100},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(emissionen, 1600.0) // erwartetes Ergebnis: 1600.0
		is.Equal(emissionenAufgeteilt, map[string]float64{
			"1": 1600.0,
		})
	})

	t.Run("BerechneDienstreisen: Bahn; Felder Tankart, Streckentyp egal", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 1, Strecke: 10, Streckentyp: "unbekannt", Tankart: "unbekannt"},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)               // Normalfall wirft keine Errors
		is.Equal(emissionen, 160.0) // erwartetes Ergebnis: 160.0
		is.Equal(emissionenAufgeteilt, map[string]float64{
			"1": 160.0,
		})
	})

	t.Run("BerechneDienstreisen: einfache Eingabe Auto", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 2, Strecke: 100, Tankart: "Diesel"},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                 // Normalfall wirft keine Errors
		is.Equal(emissionen, 48800.0) // erwartetes Ergebnis: 48800.0
		is.Equal(emissionenAufgeteilt, map[string]float64{
			"2-Diesel": 48800.0,
		})
	})

	t.Run("BerechneDienstreisen: Auto; Feld Streckentyp egal", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 2, Strecke: 10, Tankart: "Benzin", Streckentyp: "unbekannt"},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(emissionen, 5200.0) // erwartetes Ergebnis: 5200.0
		is.Equal(emissionenAufgeteilt, map[string]float64{
			"2-Benzin": 5200.0,
		})
	})

	t.Run("BerechneDienstreisen: einfache Eingabe Flugzeug", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 3, Strecke: 100, Streckentyp: "Kurzstrecke"},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                  // Normalfall wirft keine Errors
		is.Equal(emissionen, 177600.0) // erwartetes Ergebnis: 177600.0
		is.Equal(emissionenAufgeteilt, map[string]float64{
			"3-Kurzstrecke": 177600.0,
		})
	})

	t.Run("BerechneDienstreisen: Flugzeug; Feld Tankart egal", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 3, Strecke: 10, Streckentyp: "Langstrecke", Tankart: "unbekannt"},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                // Normalfall wirft keine Errors
		is.Equal(emissionen, 8360.0) // erwartetes Ergebnis: 8360.0
		is.Equal(emissionenAufgeteilt, map[string]float64{
			"3-Langstrecke": 8360.0,
		})
	})

	t.Run("BerechneDienstreisen: komplexe Eingabe", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 2, Strecke: 1200, Tankart: "Diesel"},
			{IDDienstreise: 1, Strecke: 150},
			{IDDienstreise: 3, Strecke: 750, Streckentyp: "Langstrecke"},
			{IDDienstreise: 3, Strecke: 1000, Streckentyp: "Kurzstrecke"},
			{IDDienstreise: 2, Strecke: 45, Tankart: "Benzin"},
			{IDDienstreise: 1, Strecke: 1},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)                   // Normalfall wirft keine Errors
		is.Equal(emissionen, 3014416.0) // erwartetes Ergebnis: 3014416.0
		is.Equal(emissionenAufgeteilt, map[string]float64{
			"1":             2416.0,
			"2-Diesel":      585600.0,
			"2-Benzin":      23400.0,
			"3-Kurzstrecke": 1776000.0,
			"3-Langstrecke": 627000.0,
		})
	})

	t.Run("BerechneDienstreisen: Strecke = 0 eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 1, Strecke: 0},
			{IDDienstreise: 2, Strecke: 0, Tankart: "Benzin"},
			{IDDienstreise: 2, Strecke: 0, Tankart: "Diesel"},
			{IDDienstreise: 3, Strecke: 0, Streckentyp: "Kurzstrecke"},
			{IDDienstreise: 3, Strecke: 0, Streckentyp: "Langstrecke"},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.NoErr(err)             // Normalfall wirft keine Errors
		is.Equal(emissionen, 0.0) // erwartetes Ergebnis: 0.0 (keine Strecke = keine Emissionen)
		is.Equal(emissionenAufgeteilt, map[string]float64{})
	})

	// Errortests
	t.Run("BerechneDienstreisen: ID = 0 nicht vorhanden", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 0, Strecke: 100},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(emissionen, 0.0)           // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
		is.Equal(emissionenAufgeteilt, nil)
	})

	t.Run("BerechneDienstreisen: Unbekannte DienstreisenID", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: -1, Strecke: 100},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.Equal(err, mongo.ErrNoDocuments) // Datenbank wirft ErrNoDocuments
		is.Equal(emissionen, 0.0)           // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
		is.Equal(emissionenAufgeteilt, nil)
	})

	t.Run("BerechneDienstreisen: unbekannte Tankart", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 2, Strecke: 100, Tankart: "nicht definiert"},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.Equal(err, structs.ErrTankartUnbekannt) // Funktion wirft ErrTankartUnbekannt
		is.Equal(emissionen, 0.0)                  // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
		is.Equal(emissionenAufgeteilt, nil)
	})

	t.Run("BerechneDienstreisen: unbekannter Streckentyp", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 3, Strecke: 100, Streckentyp: "nicht definiert"},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.Equal(err, structs.ErrStreckentypUnbekannt) // Funktion wirft ErrStreckentypUnbekannt
		is.Equal(emissionen, 0.0)                      // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
		is.Equal(emissionenAufgeteilt, nil)
	})

	t.Run("BerechneDienstreisen: negative Strecke eingegeben", func(t *testing.T) {
		is := is.NewRelaxed(t)

		dienstreisenDaten := []structs.UmfrageDienstreise{
			{IDDienstreise: 1, Strecke: -100},
		}

		emissionen, emissionenAufgeteilt, err := co2computation.BerechneDienstreisen(dienstreisenDaten)

		is.Equal(err, structs.ErrStreckeNegativ) // Funktion wirft ErrStreckeNegativ
		is.Equal(emissionen, 0.0)                // bei Fehlern wird 0.0 als Ergebnis zurückgegeben
		is.Equal(emissionenAufgeteilt, nil)
	})

	// Fehler ErrStrEinheitUnbekannt momentan nicht abpruefbar. Benoetigt falschen Datensatz in Datenbank
	// Fehler ErrGebaeudeSpezialfall momentan nicht abpruefbar. Benoetigt Datensatz mit spezialfall == 0 in Datenbank
}
