package tests

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestData(t *testing.T) {
	is := is.NewRelaxed(t)

	err := database.ConnectDatabase("dev")
	is.NoErr(err)
	defer func() {
		err := database.DisconnectDatabase()
		is.NoErr(err)
	}()

	t.Run("TestZaehlerrefKonsistenz", TestZaehlerrefKonsistenz)
	t.Run("TestGebaeuderefKonsistenz", TestGebaeuderefKonsistenz)
}

func TestZaehlerrefKonsistenz(t *testing.T) {
	is := is.NewRelaxed(t)

	alle_gebaeude, err := database.GebaeudeAlleNr()
	is.NoErr(err)

	for _, i := range alle_gebaeude { // gehe alle Gebaeude durch
		gebaeude, err := database.GebaeudeFind(int32(i))
		is.NoErr(err)

		// wenn Gebaeude vorhanden, pruefe alle Zaehlerreferenzen
		for _, referenz := range gebaeude.WaermeRef {
			zaehler, err := database.ZaehlerFindOID(referenz, structs.IDEnergieversorgungWaerme)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", referenz.Hex())
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range zaehler.GebaeudeRef {
				if gegenReferenz == gebaeude.GebaeudeID {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", referenz.Hex())
				fmt.Println(zaehler)
			}
		}

		for _, referenz := range gebaeude.StromRef {
			zaehler, err := database.ZaehlerFindOID(referenz, structs.IDEnergieversorgungStrom)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", referenz.Hex())
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range zaehler.GebaeudeRef {
				if gegenReferenz == gebaeude.GebaeudeID {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", referenz.Hex())
				fmt.Println(zaehler)
			}
		}

		for _, referenz := range gebaeude.KaelteRef {
			zaehler, err := database.ZaehlerFindOID(referenz, structs.IDEnergieversorgungKaelte)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", referenz.Hex())
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range zaehler.GebaeudeRef {
				if gegenReferenz == gebaeude.GebaeudeID {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", referenz.Hex())
				fmt.Println(zaehler)
			}
		}
	}
}

func TestGebaeuderefKonsistenz(t *testing.T) {
	is := is.NewRelaxed(t)

	alle_zaehler, err := database.ZaehlerAlleZaehlerUndDaten()
	is.NoErr(err)

	var alleZaehlerOIDs []primitive.ObjectID
	for _, zaehler := range alle_zaehler { // gehe alle Zaehler durch
		alleZaehlerOIDs = append(alleZaehlerOIDs, zaehler.ZaehlerID)
	}

	for _, zaehlerOID := range alleZaehlerOIDs { // gehe alle Zaehler durch
		zaehler, err := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungWaerme)
		if err != nil {
			continue
		}

		for _, referenz := range zaehler.GebaeudeRef {
			gebaeude, err := database.GebaeudeFindOID(referenz)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %s, ", referenz.Hex())
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", zaehlerOID.Hex())
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range gebaeude.WaermeRef {
				if gegenReferenz == zaehler.ZaehlerID {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %s, ", referenz.Hex())
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", zaehlerOID.Hex())
				fmt.Println(zaehler)
			}
		}
	}

	for _, zaehlerOID := range alleZaehlerOIDs { // gehe alle Gebaeude durch
		zaehler, err := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungStrom)
		if err != nil {
			continue
		}

		for _, referenz := range zaehler.GebaeudeRef {
			gebaeude, err := database.GebaeudeFindOID(referenz)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %s, ", referenz.Hex())
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", zaehlerOID.Hex())
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range gebaeude.StromRef {
				if gegenReferenz == zaehler.ZaehlerID {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %s, ", referenz.Hex())
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", zaehlerOID.Hex())
				fmt.Println(zaehler)
			}
		}
	}

	for _, zaehlerOID := range alleZaehlerOIDs { // gehe alle Gebaeude durch
		zaehler, err := database.ZaehlerFindOID(zaehlerOID, structs.IDEnergieversorgungKaelte)
		if err != nil {
			continue
		}

		for _, referenz := range zaehler.GebaeudeRef {
			gebaeude, err := database.GebaeudeFindOID(referenz)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %s, ", referenz.Hex())
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", zaehlerOID.Hex())
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range gebaeude.KaelteRef {
				if gegenReferenz == zaehler.ZaehlerID {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %s, ", referenz.Hex())
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %s, ", zaehlerOID.Hex())
				fmt.Println(zaehler)
			}
		}
	}
}
