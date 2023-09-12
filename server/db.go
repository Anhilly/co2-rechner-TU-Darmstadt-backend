package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"io/ioutil"
	"log"
	"net/http"
)

// PostAddFaktor fuegt einen neuen CO2-Faktor fuer die Energieversorgung eines bestimmten Jahres in die DB ein,
// sofern der Nutzer authentifizierter Admin ist und sendet eine Response mit null zurueck.
func PostAddFaktor(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.AddCO2Faktor{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostAddFaktor")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.EnergieversorgungAddFaktor(data)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
		} else {
			err := database.RemoveDump(ordner)
			if err != nil {
				log.Println(err)
			}
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.RemoveDump(ordner)
	if err != nil {
		log.Println(err)
	}

	sendResponse(res, true, nil, http.StatusOK)
}

// PostAddZaehlerdaten fuegt Zaehlerdaten fuer einen bestimmten Zaehler in die DB ein,
// sofern der Nutzer authentifizierter Admin ist und sendet eine Response mit null zurueck.
func PostAddZaehlerdaten(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.AddZaehlerdaten{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostAddZaehlerdaten")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.ZaehlerAddZaehlerdaten(data)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
		} else {
			err := database.RemoveDump(ordner)
			if err != nil {
				log.Println(err)
			}
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.RemoveDump(ordner)
	if err != nil {
		log.Println(err)
	}

	sendResponse(res, true, nil, http.StatusOK)
}

// PostAddZaehlerdaten fuegt Zaehlerdaten fuer einen Liste an Zaehler in die DB ein,
// sofern der Nutzer authentifizierter Admin ist und sendet eine Response mit null zurueck.
func PostAddZaehlerdatenCSV(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.AddZaehlerdatenCSV{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	combined_error := "Folgende Fehler sind aufgetreten:"
	error_encountered := false

	for i := 0; i < len(data.PKEnergie); i++ { // rufe für jeden uebergebenen Wert die Hinzufuegefunktion einzeln auf
		ordner, err := database.CreateDump("PostAddZaehlerdatenCSV")
		if err != nil {
			errorResponse(res, err, http.StatusInternalServerError)
			return
		}

		log.Printf("PK: %d, Type: %d, Jahr: %d, Wert: %f\n", data.PKEnergie[i], data.IDEnergieversorgung[i], data.Jahr, data.Wert[i])

		eachValue := structs.AddZaehlerdaten{
			PKEnergie:           data.PKEnergie[i],
			IDEnergieversorgung: data.IDEnergieversorgung[i],
			Jahr:                data.Jahr,
			Wert:                data.Wert[i],
		}

		err = database.ZaehlerAddZaehlerdaten(eachValue)

		if err != nil {
			error_encountered = true
			combined_error += fmt.Sprintf("\n\t-Zähler %d: %s", data.PKEnergie[i], err.Error())

			err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
			if err2 != nil {
				log.Println(err2)
				errorResponse(res, err2, http.StatusInternalServerError)
				return
			}
		}

		err = database.RemoveDump(ordner)
		if err != nil {
			log.Println(err)
		}
	}

	if error_encountered {
		errorResponse(res, errors.New(combined_error), http.StatusInternalServerError)
		return
	}

	sendResponse(res, true, nil, http.StatusOK)
}

// PostAddStandardZaehlerdaten fuegt alle Zaehlern ohne Zaehlerwert für das gegebene Jahr den Zaehlerwert 0.0 ein
func PostAddStandardZaehlerdaten(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.AddStandardZaehlerdaten{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostAddStandardZaehlerdaten")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.ZaehlerAddStandardZaehlerdaten(data)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
		} else {
			err := database.RemoveDump(ordner)
			if err != nil {
				log.Println(err)
			}
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.RemoveDump(ordner)
	if err != nil {
		log.Println(err)
	}

	sendResponse(res, true, nil, http.StatusOK)
}

// PostAddVersorger fügt den gegebenem Gebaeude den gegebenen Versorger hinzu, solange kein Versorger fuer
// das gegebene Jahr vorhanden ist
func PostAddVersorger(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.AddVersorger{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostAddVersorger")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.GebaeudeAddVersorger(data)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
		} else {
			err := database.RemoveDump(ordner)
			if err != nil {
				log.Println(err)
			}
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.RemoveDump(ordner)
	if err != nil {
		log.Println(err)
	}

	sendResponse(res, true, nil, http.StatusOK)
}

// PostAddStandardVersorger fuegt allen Gebaeuden ohne einen Versorger fuer das gegebene Jahr den
// Standard-Versorger 1 hinzu
func PostAddStandardVersorger(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.AddStandardVersorger{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostAddStandardVersorger")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.GebaeudeAddStandardVersorger(data)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
		} else {
			err := database.RemoveDump(ordner)
			if err != nil {
				log.Println(err)
			}
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.RemoveDump(ordner)
	if err != nil {
		log.Println(err)
	}

	sendResponse(res, true, nil, http.StatusOK)
}

// PostInsertZaehler fuegt einen neuen Zaehler in die DB ein, sofern der Nutzer authentifizierter Admin ist
// und sendet eine Response mit null zurueck.
func PostInsertZaehler(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.InsertZaehler{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostInsertZaehler")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.ZaehlerInsert(data)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
		} else {
			err := database.RemoveDump(ordner)
			if err != nil {
				log.Println(err)
			}
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.RemoveDump(ordner)
	if err != nil {
		log.Println(err)
	}

	sendResponse(res, true, nil, http.StatusOK)
}

// PostInsertGebaeude fuegt ein neues Gebeaude in die DB ein, sofern der Nutzer authentifizierter Admin ist
// und sendet eine Response mit null zurueck.
func PostInsertGebaeude(res http.ResponseWriter, req *http.Request) {
	s, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	data := structs.InsertGebaeude{}
	err = json.Unmarshal(s, &data)
	if err != nil {
		errorResponse(res, err, http.StatusBadRequest)
		return
	}

	// Datenverarbeitung
	ordner, err := database.CreateDump("PostInsertGebaeude")
	if err != nil {
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.GebaeudeInsert(data)
	if err != nil {
		err2 := database.RestoreDump(ordner) // im Fehlerfall wird vorheriger Zustand wiederhergestellt
		if err2 != nil {
			log.Println(err2)
		} else {
			err := database.RemoveDump(ordner)
			if err != nil {
				log.Println(err)
			}
		}
		errorResponse(res, err, http.StatusInternalServerError)
		return
	}

	err = database.RemoveDump(ordner)
	if err != nil {
		log.Println(err)
	}

	sendResponse(res, true, nil, http.StatusOK)
}
