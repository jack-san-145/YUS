// =======================
// Combined JavaScript for Registration and OTP Verification
// =======================

// ✅ Validation Functions (match Go backend logic)
function validatePassword(password) {
    if (password.length < 8) return false;
    const hasLower = /[a-z]/.test(password);
    const hasUpper = /[A-Z]/.test(password);
    const hasNumber = /[0-9]/.test(password);
    const hasSpecial = /[@$!%*?&]/.test(password);
    return hasLower && hasUpper && hasNumber && hasSpecial;
}

function validateCollegeEmail(email) {
    const emailPattern = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/;
    if (!emailPattern.test(email)) return false;
    const domain = email.split("@")[1];
    return domain === "kamarajengg.edu.in";
}

function validateName(name) {
    const namePattern = /^[A-Za-z ]{2,50}$/;
    return namePattern.test(name);
}

// =======================
// Registration Page
// =======================
if (document.getElementById('registerbox')) {
    document.addEventListener('DOMContentLoaded', function () {
        const emailInput = document.getElementById("email");
        const unameInput = document.getElementById("username");
        const passwordInput = document.getElementById("password");
        const rbtn = document.querySelector(".regisbtn");
        const evalid = document.getElementById("evalid");
        const nvalid = document.getElementById("nvalid");
        const pvalid = document.getElementById("pvalid");

        rbtn.addEventListener("click", async function (event) {
            event.preventDefault();
            let isValid = true;

            // Clear previous messages
            evalid.textContent = "";
            nvalid.textContent = "";
            pvalid.textContent = "";

            // Email validation
            if (!validateCollegeEmail(emailInput.value.trim())) {
                evalid.textContent = "Invalid college email (must end with @kamarajengg.edu.in)";
                isValid = false;
            }

            // Name validation
            if (!validateName(unameInput.value.trim())) {
                nvalid.textContent = "Name must contain only letters and spaces (2–50 chars)";
                isValid = false;
            }

            // Password validation
            if (!validatePassword(passwordInput.value.trim())) {
                pvalid.textContent = "Password must be ≥8 chars with upper, lower, number & special (@$!%*?&)";
                isValid = false;
            }

            if (!isValid) return;

            // ✅ Prepare form data for backend (x-www-form-urlencoded)
            const formData = new URLSearchParams();
            formData.append("email", emailInput.value.trim());
            formData.append("name", unameInput.value.trim());
            formData.append("password", passwordInput.value.trim());

            console.log("Sending OTP request:", formData.toString());

            try {
                const response = await fetch("https://yus.kwscloud.in/yus/send-otp-admin", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/x-www-form-urlencoded"
                    },
                    body: formData
                });

                const data = await response.json();
                console.log("OTP Send Response:", data);

                if (data.otp_sent === true) {
                    sessionStorage.setItem('adminEmail', emailInput.value.trim());
                    sessionStorage.setItem('adminName', unameInput.value.trim());
                    sessionStorage.setItem('adminPassword', passwordInput.value.trim());

                    alert("OTP sent successfully! Please check your email.");
                    window.location.href = "/yus/serve-otp-verify-page";
                } else if (data.otp_sent === "Admin already exists") {
                    alert("Admin already exists with this email.");
                } else if (data.otp_sent === false) {
                    alert("Failed to send OTP. Please try again.");
                } else {
                    alert("Unexpected response from server.");
                }
            } catch (error) {
                console.error("Error sending OTP:", error);
                alert("Network error. Please check your connection and try again.");
            }
        });
    });
}

// =======================
// OTP Verification Page
// =======================
if (document.getElementById('otpbox')) {
    document.addEventListener('DOMContentLoaded', function () {
        const otpInputs = document.querySelectorAll('.otp-input');
        const form = document.getElementById('otpForm');
        const resendLink = document.getElementById('resend-link');
        const timerElement = document.getElementById('timer');

        let countdown = 180;
        let timerInterval;

        startTimer();

        otpInputs.forEach((input, index) => {
            input.addEventListener('input', function () {
                if (this.value) {
                    this.classList.add('filled');
                    this.classList.add('pulse');
                    setTimeout(() => this.classList.remove('pulse'), 500);
                } else {
                    this.classList.remove('filled');
                }

                if (this.value && index < otpInputs.length - 1) {
                    otpInputs[index + 1].focus();
                }
            });

            input.addEventListener('keydown', function (e) {
                if (e.key === 'Backspace' && !this.value && index > 0) {
                    otpInputs[index - 1].focus();
                }
            });

            input.addEventListener('keypress', function (e) {
                if (!/[0-9]/.test(e.key)) e.preventDefault();
            });
        });

        // ✅ OTP submission (x-www-form-urlencoded)
        form.addEventListener('submit', async function (e) {
            e.preventDefault();
            let otp = '';
            otpInputs.forEach(input => otp += input.value);

            if (otp.length !== 6) {
                alert('Please enter a valid 6-digit OTP');
                return;
            }

            const email = sessionStorage.getItem('adminEmail');
            const name = sessionStorage.getItem('adminName');
            const password = sessionStorage.getItem('adminPassword');

            if (!email || !name || !password) {
                alert('Session expired. Please restart the registration process.');
                window.location.href = "registerform.html";
                return;
            }

            // ✅ Prepare form data
            const formData = new URLSearchParams();
            formData.append("email", email);
            formData.append("name", name);
            formData.append("password", password);
            formData.append("otp", otp);

            try {
                const response = await fetch("https://yus.kwscloud.in/yus/verify-otp-admin", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/x-www-form-urlencoded"
                    },
                    body: formData
                });

                const data = await response.json();
                console.log("OTP Verification Response:", data);

                if (data.status === "successfully added admin") {
                    alert("Registration successful! Welcome to YUS.");
                    sessionStorage.clear();
                    window.location.href = "/yus/serve-index-page";
                } else if (data.status === "invalid otp") {
                    alert("Invalid OTP. Please try again.");
                    otpInputs.forEach(input => { input.value = ''; input.classList.remove('filled'); });
                    otpInputs[0].focus();
                } else if (data.status === "Admin already exists") {
                    alert("Admin already exists with this email.");
                    window.location.href = "registerform.html";
                } else {
                    alert("Unexpected response: " + data.status);
                }
            } catch (error) {
                console.error("Error verifying OTP:", error);
                alert("Network error during verification. Please try again.");
            }
        });

        // ✅ Resend OTP (x-www-form-urlencoded)
        resendLink.addEventListener('click', async function (e) {
            e.preventDefault();
            if (countdown !== 0) return;

            const email = sessionStorage.getItem('adminEmail');
            const name = sessionStorage.getItem('adminName');
            const password = sessionStorage.getItem('adminPassword');

            if (!email || !name || !password) {
                alert('Session expired. Please restart the registration process.');
                return;
            }

            const formData = new URLSearchParams();
            formData.append("email", email);
            formData.append("name", name);
            formData.append("password", password);

            try {
                const response = await fetch("https://yus.kwscloud.in/yus/send-otp-admin", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/x-www-form-urlencoded"
                    },
                    body: formData
                });

                const data = await response.json();
                console.log("Resend OTP Response:", data);

                if (data.otp_sent === true) {
                    countdown = 180;
                    startTimer();
                    alert('OTP has been resent to your email!');
                } else {
                    alert('Failed to resend OTP. Please try again.');
                }
            } catch (error) {
                console.error("Error resending OTP:", error);
                alert("Network error while resending OTP.");
            }
        });

        // Timer
        function startTimer() {
            if (timerInterval) clearInterval(timerInterval);
            updateTimerDisplay();
            timerInterval = setInterval(() => {
                countdown--;
                updateTimerDisplay();
                if (countdown <= 0) clearInterval(timerInterval);
            }, 1000);
        }

        function updateTimerDisplay() {
            if (countdown > 0) {
                const minutes = Math.floor(countdown / 60);
                const seconds = countdown % 60;
                timerElement.textContent = ` (${minutes}:${seconds < 10 ? '0' : ''}${seconds})`;
                resendLink.style.pointerEvents = 'none';
                resendLink.style.opacity = '0.5';
            } else {
                timerElement.textContent = '';
                resendLink.style.pointerEvents = 'auto';
                resendLink.style.opacity = '1';
            }
        }

        otpInputs[0].focus();
    });
}