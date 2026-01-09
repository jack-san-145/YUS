package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"strconv"
)

// func SendEmailTo(email string, otp string) bool {

// 	from := "yellohbus@gmail.com"
// 	password := os.Getenv("YUS_EMAIL_PASSWORD")
// 	fmt.Println("password send mail - ", password)
// 	msg := []byte("Subject: Your Otp for YUS - " + otp + "\r\n\r\nYour One time OTP is " + otp + "\r\n\r\n- YUS(Yelloh Bus) Team")
// 	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
// 	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{email}, msg)
// 	if err != nil {
// 		fmt.Println("otp not sent - ", err)
// 		// ch <- false
// 		return false
// 	}
// 	fmt.Println("Email sent ")
// 	// ch <- true
// 	return true
// }

func SendEmailTo(email string, otp string) bool {

	from := "yusofficialteam@gmail.com"
	password := os.Getenv("YUS_OFFICIAL_TEAM_EMAIL_PASSWORD")

	// HTML email template
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<body style="font-family: Arial, sans-serif; padding: 20px; background-color: #f5f5f5;">
			<div style="max-width: 500px; margin: auto; background: white; padding: 25px; border-radius: 10px;">
				<h2 style="color: #1A73E8;">🔐 YUS Verification Code</h2>
				<p>Your one-time verification code is:</p>
				<h1 style="color:#1A73E8; letter-spacing: 5px; margin: 20px 0;">%s</h1>
				<p>This code is valid for <b>3 minutes</b>. Do not share it with anyone.</p>
				<br>
				<p>Thank you,<br><b>YUS (Yelloh Bus) Team</b></p>
			</div>
		</body>
		</html>
	`, otp)

	// Plain text fallback
	plainBody := "Your Driver Verification Code: " + otp + "\nValid for 3 minutes.\nDo not share this code with anyone."

	// Construct MIME message
	msg := []byte(
		"Subject: YUS Verification Code : " + otp + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: multipart/alternative; boundary=YUSBOUNDARY\r\n\r\n" +
			"--YUSBOUNDARY\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			plainBody + "\r\n\r\n" +
			"--YUSBOUNDARY\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
			htmlBody + "\r\n\r\n" +
			"--YUSBOUNDARY--",
	)

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{email}, msg)
	if err != nil {
		fmt.Println("otp not sent - ", err)
		return false
	}

	fmt.Println("Email sent to:", email)
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
