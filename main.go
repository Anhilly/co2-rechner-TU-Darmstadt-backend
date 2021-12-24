package main

import (
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/server"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"log"
)

func main() {
	err := database.ConnectDatabase()
	if err != nil {
		panic(err)
	}

	server.StartServer()
}

// Funktion um Umfragen in die Datenbank zu schreiben
func testumfrage() {
	umfrage := structs.InsertUmfrage{
		Mitarbeiteranzahl: 10,
		Jahr:              2020,
		NutzerEmail:       "anton@tobi.com",
		Gebaeude: []structs.UmfrageGebaeude{
			{GebaeudeNr: 1101, Nutzflaeche: 1000}, // Waerme
			{GebaeudeNr: 1220, Nutzflaeche: 1000}, // Strom, Kaelte
		},
		ITGeraete: []structs.UmfrageITGeraete{
			{IDITGeraete: 4, Anzahl: 5},  // Beamer
			{IDITGeraete: 6, Anzahl: 5},  // Server
			{IDITGeraete: 7, Anzahl: 5},  // Multifunktionsdrucker
			{IDITGeraete: 8, Anzahl: 5},  // Patronen Multifunktionsdrucker
			{IDITGeraete: 9, Anzahl: 5},  // Laser-/Tintenstrahldrucker
			{IDITGeraete: 10, Anzahl: 5}, // Patronen Laser-/Tintenstrahldrucker
		},
	}

	id, err := database.UmfrageInsert(umfrage)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(id)

	mitarbeiterumfrage1 := structs.InsertMitarbeiterUmfrage{
		IDUmfrage: id,
		Pendelweg: []structs.UmfragePendelweg{
			{IDPendelweg: 1, Strecke: 100, Personenanzahl: 1},  // Fahrrad
			{IDPendelweg: 2, Strecke: 100, Personenanzahl: 1},  // E-Fahrrad
			{IDPendelweg: 3, Strecke: 100, Personenanzahl: 1},  // Motorisiertes Zweirad
			{IDPendelweg: 4, Strecke: 100, Personenanzahl: 1},  // PKW (Diesel)
			{IDPendelweg: 5, Strecke: 100, Personenanzahl: 1},  // PKW (Benzin)
			{IDPendelweg: 6, Strecke: 100, Personenanzahl: 1},  // Bus
			{IDPendelweg: 7, Strecke: 100, Personenanzahl: 1},  // Bahn
			{IDPendelweg: 8, Strecke: 100, Personenanzahl: 1},  // U-Bahn
			{IDPendelweg: 9, Strecke: 100, Personenanzahl: 1},  // Straßenbahn
			{IDPendelweg: 10, Strecke: 100, Personenanzahl: 1}, // Mix oeffentlichen Verkehrsmittel mit U-Bahn
			{IDPendelweg: 11, Strecke: 100, Personenanzahl: 1}, // Mix oeffentlichen Verkehrsmittel ohne U-Bahn
		},
		TageImBuero: 2,
		Dienstreise: []structs.UmfrageDienstreise{
			{IDDienstreise: structs.IDDienstreiseBahn, Strecke: 1000},
			{IDDienstreise: structs.IDDienstreiseAuto, Strecke: 1000, Tankart: "Benzin"},
			{IDDienstreise: structs.IDDienstreiseAuto, Strecke: 1000, Tankart: "Diesel"},
			{IDDienstreise: structs.IDDienstreiseFlugzeug, Strecke: 1000, Streckentyp: "Langstrecke"},
			{IDDienstreise: structs.IDDienstreiseFlugzeug, Strecke: 1000, Streckentyp: "Kurzstrecke"},
		},
		ITGeraete: []structs.UmfrageITGeraete{
			{IDITGeraete: 1, Anzahl: 5}, // Notebooks
			{IDITGeraete: 2, Anzahl: 5}, // Desktop PC
			{IDITGeraete: 3, Anzahl: 5}, // Bildschirm
			{IDITGeraete: 5, Anzahl: 5}, // Mobiltelefon
		},
	}

	id, err = database.MitarbeiterUmfrageInsert(mitarbeiterumfrage1)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(id)
}
