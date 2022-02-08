package server

import (
	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SendeBestaetigungsMail(id primitive.ObjectID, username string) error {
	mailjetClient := mailjet.NewMailjetClient(publicAPIKey, privateAPIKey)

	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "anton.hillmann@stud.tu-darmstadt.de",
				Name:  "CO2-Rechner",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: username,
					Name:  "",
				},
			},
			TemplateID:       3532127,
			TemplateLanguage: true,
			Subject:          "",
			Variables: map[string]interface{}{
				"confirmation_link": "https://localhost:8080/#/user/" + id.Hex(),
			},
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		return err
	}

	return nil
}

func SendePasswortVergessenMail(username string, passwort string) error {

	mailjetClient := mailjet.NewMailjetClient(publicAPIKey, privateAPIKey)
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "anton.hillmann@stud.tu-darmstadt.de",
				Name:  "CO2-Rechner",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: username,
					Name:  "",
				},
			},
			TemplateID:       3535931,
			TemplateLanguage: true,
			Subject:          "CO2-Rechner Passwort Zurücksetzen",
			Variables: map[string]interface{}{
				"passwort": passwort,
			},
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mailjetClient.SendMailV31(&messages)

	if err != nil {
		return err
	}

	return nil
}
