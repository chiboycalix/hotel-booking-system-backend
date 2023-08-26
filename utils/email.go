package utils

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"

	brevo "github.com/getbrevo/brevo-go/lib"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/models"
)

type forgetPassword struct {
	ID          string `json:"id"          bson:"_id"`
	Email       string `json:"email"       bson:"email"`
	FirstName   string `json:"firstName"   bson:"firstName"`
	LastName    string `json:"lastName"    bson:"lastName"`
	FrontendUrl string `json:"frontendUrl" bson:"frontendUrl"`
}

func SendMailService(user models.User, templatePath string, subject string) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal("error parsing template")
		return err
	}
	t.Execute(
		&body,
		forgetPassword{
			ID:          user.ID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			FrontendUrl: common.FrontendUrl(),
		},
	)

	var ctx context.Context
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", common.BrevoAPIKey())
	br := brevo.NewAPIClient(cfg)
	_, _, err = br.TransactionalEmailsApi.SendTransacEmail(ctx, brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  "Hotel Booking System",
			Email: common.SenderEmail(),
		},
		To: []brevo.SendSmtpEmailTo{
			{Name: "chi", Email: user.Email},
		},
		HtmlContent: body.String(),
		Subject:     subject,
	})
	if err != nil {
		fmt.Println("Error ", err)
		return err
	}
	fmt.Println("Email sent successfully")
	return nil
}
