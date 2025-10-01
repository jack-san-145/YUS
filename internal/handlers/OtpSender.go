package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tet/internals/services"
	"tet/internals/storage/postgres"
	"tet/internals/storage/redis"
	"time"
)

func generateOtp() string {
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
