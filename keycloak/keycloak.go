package keycloak

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/config"
	"github.com/Nerzal/gocloak/v13"
	"github.com/pkg/errors"
)

var KeycloakClient *gocloak.GoCloak

func SetupKeycloakClient(mode string) error {
	if mode == "prod" {
		KeycloakClient = gocloak.NewClient(config.ProdKeycloakUrl)
	} else if mode == "dev" {
		KeycloakClient = gocloak.NewClient(config.DevKeycloakUrl)
	} else {
		return errors.New("MODE not set")
	}

	return nil
}
