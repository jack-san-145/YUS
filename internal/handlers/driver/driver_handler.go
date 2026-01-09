package driver

import (
	"log"
	"net/http"
	"strconv"

	"yus/internal/handlers/common/response"
	"yus/internal/services"
)

func (h *DriverHandler) SendOTPHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var otp_status = make(map[string]bool)
	err := r.ParseForm()
	if err != nil {
		log.Println("error while parsing the form - ", err)
		return
	}
	driver_id := r.FormValue("driver_id")
	email := r.FormValue("email")

	driver_id_int, err := strconv.Atoi(driver_id)
	if err != nil {
		log.Println("error while converting the driver_id string to driver_id_int - ", err)
	}
	if exists, _ := h.Store.DB.DriverExists(ctx, driver_id_int); !exists {
		response.WriteJSON(w, r, map[string]string{"status": "no driver found"})
		return
	}
	if services.ValidateEmail(email) {
		otp := services.GenerateOtp()
		is_email_sent := services.SendEmailTo(email, otp)
		if is_email_sent {
			h.Store.InMemoryDB.SetOtp(ctx, email, otp) //set otp to redis if otp sent to email successfully
		}

		otp_status["otp_sent"] = is_email_sent
	} else {
		otp_status["otp_sent"] = false
	}

	response.WriteJSON(w, r, otp_status)
}

func (h *DriverHandler) VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var pass_status = make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		log.Println("error while parsing form")
		return
	}
	driver_id_string := r.FormValue("driver_id")
	driver_id_int, err := strconv.Atoi(driver_id_string)
	if err != nil {
		log.Println("error while converting the driver_id string to driver_id_int - ", err)
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	given_otp := r.FormValue("otp")
	log.Println("verify otp for - ", email, given_otp)

	if services.ValidateEmail(email) && services.ValidatePassword(password) {

		otp, _ := h.Store.InMemoryDB.GetOtp(ctx, email)
		if given_otp == otp {
			ok, _ := h.Store.DB.SetDriverPassword(ctx, driver_id_int, email, password)
			if ok {
				pass_status["status"] = "success"
			} else {
				pass_status["status"] = "failed"
			}

		} else {
			pass_status["status"] = "failed"
		}

	} else {
		pass_status["status"] = "failed"
	}
	response.WriteJSON(w, r, pass_status)
}

func (h *DriverHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var login_status = make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		log.Println("error while parsing driver login form")
		login_status["login_status"] = "invalid"
	} else {
		driver_id := r.FormValue("driver_id")

		driver_id_int, err := strconv.Atoi(driver_id)
		if err != nil {
			log.Println("error while converting the driver_id string to driver_id_int - ", err)
		}
		password := r.FormValue("password")
		if valid, _ := h.Store.DB.ValidateDriver(ctx, driver_id_int, password); valid {
			login_status["login_status"] = "valid"
			session_id, err := h.Store.InMemoryDB.CreateDriverSession(ctx, driver_id_int)
			if err != nil {
				login_status["login_status"] = "invalid"
				return
			}
			login_status["session_id"] = session_id
		} else {
			login_status["login_status"] = "invalid"
		}
	}
	response.WriteJSON(w, r, login_status)
}

func (h *DriverHandler) GetAllocatedBusHandler(w http.ResponseWriter, r *http.Request) {
	//yus.kwscloud.in/yus/get-allotted-bus

	ctx := r.Context()

	driver_id := r.Context().Value("DRIVER_ID").(int)

	log.Println("driver_id- ", driver_id)
	alloted_bus, _ := h.Store.DB.GetAllottedBusForDriver(ctx, driver_id)
	if alloted_bus.BusID != 0 && alloted_bus.RouteId != 0 {
		response.WriteJSON(w, r, alloted_bus)
	} else {
		response.WriteJSON(w, r, map[string]string{"status": "no bus allotted"})
	}
}
