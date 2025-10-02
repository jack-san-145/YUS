package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"strconv"
)

func SendEmailTo(email string, otp string) bool {

	from := "yellohbus@gmail.com"
	password := os.Getenv("YUS_EMAIL_PASSWORD")
	fmt.Println("password send mail - ", password)
	msg := []byte("Subject: Your Otp for YUS - " + otp + "\r\n\r\nYour One time OTP is " + otp + "\r\n\r\n- YUS(Yelloh Bus) Team")
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{email}, msg)
	if err != nil {
		fmt.Println("otp not sent - ", err)
		// ch <- false
		return false
	}
	fmt.Println("Email sent ")
	// ch <- true
	return true
}

func GenerateOtp() string {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		fmt.Println("error while generating OTP - ", err)
		return ""
	}
	otp := int(n.Int64()) + 100000
	fmt.Println("generated otp - ", otp)
	otp_string := strconv.Itoa(otp)
	return otp_string
}
