package structs

import "errors"

var ( // Fehler von Add, Delete und Insert Funktionen
	// Fehler durch Nutzereingabe
	ErrZaehlerVorhanden = errors.New("Es ist schon ein Zaehler mit dem DPName vorhanden")

	// Fehler durch Nutzereingabe
	ErrFehlendeGebaeuderef = errors.New("Neuer Zaehler hat keine Referenzen auf Gebaeude")

	// Fehler durch Nutzereingabe
	ErrJahrVorhanden = errors.New("Ein Wert ist fuer das angegebene Jahr schon vorhanden")

	// Fehler durch Nutzereingabe
	ErrJahrUndMonatVorhanden = errors.New("Ein Wert fuer den angegebene Monat im Jahr ist schon vorhanden")

	// Fehler durch Nutzereingabe
	ErrGebaeudeVorhanden = errors.New("Ein Gebaeude mit der angegeben Nummer existiert schon in der Datenbank")

	// Fehler durch Nutzereingabe
	ErrGebaeudeNichtVorhanden = errors.New("Ein Gebaeude mit der angegeben Nummer konnte nicht gefunden werden")

	// Fehler durch Nutzereingabe
	ErrIDEnergieversorgungNichtVorhanden = errors.New("Die angegebene IDEnergieversorgung ist nicht vorhanden")

	// Fehler bei Erstellung der ObjectID
	ErrObjectIDNichtKonvertierbar = errors.New("ObjectID Konvertierung fehlerhaft")

	// Fehler beim loeschen des Eintrags wo ObjectID nicht gefunden wurde
	ErrObjectIDNichtGefunden = errors.New("ObjectID nicht gefunden")

	// Fehler beim loeschen eines Nutzers wo Username nicht gefunden wurde
	ErrUsernameLoeschenFehlgeschlagen = errors.New("Username konnte nicht geloescht werden")

	// Fehler beim loeschen des Eintrags wo ObjectID nicht gefunden wurde
	ErrUmfrageVollstaendig = errors.New("Umfrage ist bereits von allen Mitarbeitenden ausgefüllt.")
)

var ( // Fehler von Find Funktionen
	// Fehler beim Abrufen von mehreren Dokumenten
	ErrDokumenteNichtGefunden = errors.New("Es konnten nicht alle angefragten Dokumente gefunden werden!")

	// Fehler beim Finden einer Umfrage zu MitarbeiterUmfragen
	ErrMitarbeiterUmfrageMehrfachAssoziiert = errors.New("Die gegebene MitarbeiterUmfrage ist in mehreren Umfragen referenziert.")
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
	ErrStrGebaeuderefFehlt = "%s: Zaehler %s hat keine Referenzen auf Gebaeude"

	// Fehler durch fehlende Werte in Datenbank
	ErrStrVerbrauchFehlt = "%s: Kein Verbrauch für das Jahr %d, Zaehler: %s"

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

	// Fehler durch fehlende Angabe des Versorgers für Jahr
	ErrStrKeinVersorger = "gebaeudeNormalfall: Gebaeude %d hat keinen Versorger für Jahr %d eingetragen"

	// Fehler durch fehlenden Energiefaktor für bestimmtes Jahr und Energieform
	ErrStrKeinFaktorFuerVertrag = "Für Jahr %d ist kein Faktor für Energie %d und Vertrag %d"

	// Fehler
	ErrVertragNichtVorhanden = errors.New("Angegeben ID fuer den Vertrag ist nicht vorhanden")
)

var ( // Fehler die bei der Authentifizierung auftreten
	// Nutzer will Account mit bestehender Username registrieren
	ErrInsertExistingAccount = errors.New("Account mit diesem Nutzernamen existiert bereits")

	// Nutzer nicht berechtigt oder kein Admin
	ErrNutzerHatKeineBerechtigung = errors.New("Der Nutzer hat nicht die passende Berechtigung")
)
