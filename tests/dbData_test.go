package tests

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
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

	for i := 1000; i <= 6000; i++ { // gehe alle Gebaeude durch
		gebaeude, err := database.GebaeudeFind(int32(i))
		if err != nil {
			continue
		}

		// wenn Gebaeude vorhanden, pruefe alle Zaehlerreferenzen
		for _, referenz := range gebaeude.WaermeRef {
			zaehler, err := database.ZaehlerFind(referenz, structs.IDEnergieversorgungWaerme)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", referenz)
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range zaehler.GebaeudeRef {
				if gegenReferenz == gebaeude.Nr {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", referenz)
				fmt.Println(zaehler)
			}
		}

		for _, referenz := range gebaeude.StromRef {
			zaehler, err := database.ZaehlerFind(referenz, structs.IDEnergieversorgungStrom)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", referenz)
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range zaehler.GebaeudeRef {
				if gegenReferenz == gebaeude.Nr {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", referenz)
				fmt.Println(zaehler)
			}
		}

		for _, referenz := range gebaeude.KaelteRef {
			zaehler, err := database.ZaehlerFind(referenz, structs.IDEnergieversorgungKaelte)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", referenz)
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range zaehler.GebaeudeRef {
				if gegenReferenz == gebaeude.Nr {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %d, ", gebaeude.Nr)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", referenz)
				fmt.Println(zaehler)
			}
		}
	}
}

func TestGebaeuderefKonsistenz(t *testing.T) {
	is := is.NewRelaxed(t)

	for i := 1000; i <= 12000; i++ { // gehe alle Gebaeude durch
		zaehler, err := database.ZaehlerFind(int32(i), structs.IDEnergieversorgungWaerme)
		if err != nil {
			continue
		}

		for _, referenz := range zaehler.GebaeudeRef {
			gebaeude, err := database.GebaeudeFind(referenz)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %d, ", referenz)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", i)
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range gebaeude.WaermeRef {
				if gegenReferenz == zaehler.PKEnergie {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %d, ", referenz)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", zaehler.PKEnergie)
				fmt.Println(zaehler)
			}
		}
	}

	for i := 1000; i <= 12000; i++ { // gehe alle Gebaeude durch
		zaehler, err := database.ZaehlerFind(int32(i), structs.IDEnergieversorgungStrom)
		if err != nil {
			continue
		}

		for _, referenz := range zaehler.GebaeudeRef {
			gebaeude, err := database.GebaeudeFind(referenz)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %d, ", referenz)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", i)
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range gebaeude.StromRef {
				if gegenReferenz == zaehler.PKEnergie {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %d, ", referenz)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", zaehler.PKEnergie)
				fmt.Println(zaehler)
			}
		}
	}

	for i := 1000; i <= 12000; i++ { // gehe alle Gebaeude durch
		zaehler, err := database.ZaehlerFind(int32(i), structs.IDEnergieversorgungKaelte)
		if err != nil {
			continue
		}

		for _, referenz := range zaehler.GebaeudeRef {
			gebaeude, err := database.GebaeudeFind(referenz)
			is.NoErr(err)

			if err != nil { // Ausgabe, falls ein Fehler gefunden wurde
				fmt.Println("Referenz nicht gefunden:")
				fmt.Printf("Gebaeude-Nummer: %d, ", referenz)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", i)
				fmt.Println(zaehler)
			}

			found := false
			for _, gegenReferenz := range gebaeude.KaelteRef {
				if gegenReferenz == zaehler.PKEnergie {
					found = true
					break
				}
			}
			is.Equal(found, true)
			if !found {
				fmt.Println("Fehlende Rueckreferenz:")
				fmt.Printf("Gebaeude-Nummer: %d, ", referenz)
				fmt.Println(gebaeude)
				fmt.Printf("Zaehler-Nummer: %d, ", zaehler.PKEnergie)
				fmt.Println(zaehler)
			}
		}
	}
}
