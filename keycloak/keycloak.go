package keycloak

import (
	"context"
	"errors"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/config"
	"github.com/Nerzal/gocloak/v13"
)

var (
	KeycloakClient *gocloak.GoCloak
	realm          string
)

// SetupKeycloakClient setzt den KeycloakClient
func SetupKeycloakClient(mode string) error {
	if mode == "prod" {
		KeycloakClient = gocloak.NewClient(config.ProdKeycloakUrl)
		realm = config.ProdKeycloakRealm
	} else if mode == "dev" {
		KeycloakClient = gocloak.NewClient(config.DevKeycloakUrl)
		realm = config.DevKeycloakRealm
	} else if mode == "test" {
		KeycloakClient = gocloak.NewClient(config.TestKeycloakUrl)
		realm = config.TestKeycloakRealm
	} else {
		return errors.New("MODE not set")
	}

	return nil
}

// GetEmailFromToken gibt die E-Mail aus dem Token zurueck.
func GetEmailFromToken(token string, ctx context.Context) (string, error) {
	userInfo, err := KeycloakClient.GetUserInfo(ctx, token, realm)
	if err != nil {
		return "", err
	}

	var email string
	if userInfo.Email != nil {
		email = *userInfo.Email
	} else {
		return "", errors.New("Cannot retrieve email from token")
	}

	return email, nil
}

// GetUsernameFromToken gibt die Nutzernamen aus dem Token zurueck.
func GetUsernameFromToken(token string, ctx context.Context) (string, error) {
	userInfo, err := KeycloakClient.GetUserInfo(ctx, token, realm)
	if err != nil {
		return "", err
	}

	var nutzername string
	if userInfo.PreferredUsername != nil {
		nutzername = *userInfo.PreferredUsername
	} else {
		return "", errors.New("Cannot retrieve username from token")
	}

	return nutzername, nil
}
