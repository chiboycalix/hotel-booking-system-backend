package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/models"
	"gopkg.in/gomail.v2"
)

type forgetPassword struct {
	ID        string `json:"id" bson:"_id"`
	Email     string `json:"email" bson:"email"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
}

func SendMailService(user models.User, templatePath string) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Println(err, "err")
		log.Fatal("error parsing template")
		return err
	}
	t.Execute(&body, forgetPassword{ID: user.ID, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName})
	m := gomail.NewMessage()
	m.SetHeader("From", common.EMAIL())
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Forget Password")
	m.SetBody("text/html", body.String())
	port, err := strconv.Atoi(os.Getenv("EMAIL_SERVICE_PORT"))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	d := gomail.NewDialer(os.Getenv("EMAIL_SERVICE_HOST"), port, common.EMAIL(), common.EmailPassword())
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	fmt.Println("email sent successfully")
	return nil
}
