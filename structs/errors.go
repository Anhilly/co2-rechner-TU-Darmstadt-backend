package structs

import "errors"

var ( // Fehler von Add und Insert Funktionen
	// Fehler durch Nutzereingabe
	ErrZaehlerVorhanden = errors.New("Es ist schon ein Zaehler mit dem PK vorhanden")

	// Fehler durch Nutzereingabe
	ErrFehlendeGebaeuderef = errors.New("Neuer Zaehler hat keine Referenzen auf Gebaeude")

	// Fehler durch Nutzereingabe
	ErrJahrVorhanden = errors.New("Ein Wert ist fuer das angegebene Jahr schon vorhanden")

	// Fehler durch Nutzereingabe
	ErrGebaeudeVorhanden = errors.New("Ein Gebaeude mit der angegeben Nummer existiert schon in der Datenbank")

	// Fehler durch Nutzereingabe
	ErrIDEnergieversorgungNichtVorhanden = errors.New("Die angegebene IDEnergieversorgung ist nicht vorhanden")

	// Fehler bei Erstellung der ObjectID
	ErrObjectIDNichtKonvertierbar = errors.New("ObjectID Konvertierung fehlerhaft")
)

var ( // Fehler, die bei Berechnungen auftreten
	// Fehler durch Nutzereingabe
	ErrJahrNichtVorhanden = errors.New("getEnergieCO2Faktor: Kein CO2 Faktor für angegebenes Jahr vorhanden")

	// Fehler durch Nutzereingabe
	ErrFlaecheNegativ = errors.New("gebaeudeNormalfall: Nutzflaeche ist negativ")

	// Fehler durch fehlende Behandlung eines Gebaeudespezialfalls im Code
	ErrGebaeudeSpezialfall = errors.New("BerechneEnergieverbrauch: Spezialfall fuer Gebaeude nicht abgedeckt")

	// Fehler durch fehlende Behandlung eines Zaehlerspezialfalls im Code
	ErrZaehlerSpezialfall = errors.New("gebaeudeNormalfall: Spezialfall fuer Zaehler nicht abgedeckt")

	// Fehler durch falsche Daten in Datenbank
	ErrStrGebaeuderefFehlt = "%s: Zaehler %d hat keine Referenzen auf Gebaeude"

	// Fehler durch fehlende Werte in Datenbank
	ErrStrVerbrauchFehlt = "%s: Kein Verbrauch für das Jahr %d, Zaehler: %d"

	// Fehler durch nicht behandelte Einheit oder Fehler in der Datenbank
	ErrStrEinheitUnbekannt = "%s: Einheit %s unbekannt"

	// Fehler durch Nutzereingabe
	ErrPersonenzahlZuKlein = errors.New("BerechnePendelweg: Personenzahl ist kleiner als 1")

	// Fehler durch Nutzereingabe (oder Wert fehlt in Datenbank)
	ErrTankartUnbekannt = errors.New("BerechneDienstreisen: Tankart nicht vorhanden")

	// Fehler durch Nutzereingabe (oder Wert fehlt in Datenbank)
	ErrStreckentypUnbekannt = errors.New("BerechneDienstreisen: Streckentyp nicht vorhanden")

	// Fehler durch Nutzereingabe
	ErrStreckeNegativ = errors.New("BerechneDienstreise / BerechnePendelweg: Strecke ist negativ")

	// Fehler durch Nutzereingabe
	ErrAnzahlNegativ = errors.New("BerechneITGeraete: Anzahl an IT-Geraeten ist negativ")

	// Fehler durch fehlende Implementierung einer Berechnung
	ErrBerechnungUnbekannt = errors.New("BerechneDienstreisen: Keine Berechnung fuer angegeben ID vorhanden")
)

var ( // Fehler die bei der Authentifizierung auftreten
	// Nutzer will Account mit bestehender Email registrieren
	ErrInsertExistingAccount = errors.New("Account mit dieser Email existiert bereits")

	// Fehler das für Nutzer kein Sessiontoken registriert ist
	ErrNutzerHatKeinenSessiontoken = errors.New("Email hat keinen Sessiontoken registriert")

	// Abgelaufener Sessiontoken
	ErrAbgelaufenerSessiontoken = errors.New("Der Sessiontoken ist abgelaufen")

	// Falsches Passwort beim Anmelden
	ErrFalschesPasswortError = errors.New("Die Kombination aus Passwort und Email stimmt nicht überein")

	// Falsche Email beim Anmelden
	ErrNichtExistenteEmail = errors.New("Es existiert kein Konto mit dieser Email")

	// Falscher Sessiontoken beim Authentifizieren
	ErrFalscherSessiontoken = errors.New("Falscher Sessiontoken für diesen Benutzer")
)