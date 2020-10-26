package seller

import (
	"fmt"
	"go.uber.org/zap"
	"net/smtp"
)

func NewEmailProvider() *EmailProvider {
	return &EmailProvider{}
}

type EmailProvider struct {
}

func (ep *EmailProvider) StockChanged(oldStock int, newStock int, seller Seller, product string) {
	fmt.Println(oldStock, newStock, product)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// logging product stock change
	logger.Info("SMS Warning sent to",
		// Structured context as strongly typed Field values.
		zap.String("Seller", seller.UUID),
		zap.String("Phone:", seller.Phone),
		zap.String("Product", product),
		zap.String("Event", "Product stock changed"),
	)
	//SMS Warning sent to {seller_UUID} (Phone: {seller_Phone}): {Product_name} Product stock changed
	err:=emailSender(seller.Name,seller.Email, product)
	if err != nil{
		logger.Error("An error happened while sending email to seller",zap.Error(err))
	}
}

func emailSender(name,email, product string) error {
	// Gmail account just for sending test emails
	from := "go.challenge.gfg@gmail.com"
	password := "thisIsAgoChallenge"
	// Receiver email address.
	to := []string{"japalaciosmo@gmail.com"}
	// Message.
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + "Stock change" + "!\n"
	msg := []byte(subject + mime + "\n" + " Dear Mr./Ms./Dr. "+name+" Product: "+product + "Product stock changed")
	// Authentication.
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	// Sending email.
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
