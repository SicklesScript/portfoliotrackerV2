const form = document.getElementById("signup-form");

// Listens for submit event 
form.addEventListener("submit", async(e) => {
    e.preventDefault();

    // Extracts form values 
    const email = document.getElementById("email").value
    const pass = document.getElementById("password").value
    const passConfirmed = document.getElementById("passwordConfirmed").value

    // If both password entries do not match, alert and return
    if (pass != passConfirmed) {
        alert("Passwords do not match");
        return;
    };

    const data = {
        email: email,
        password: pass,
        passwordConfirmed: passConfirmed,
    };

    try {
        // Pass data to signupAPI
        const response = await fetch("/api/signup", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            // Convert data to json friendly format
            body: JSON.stringify(data)
        });

        // Redirect user to login upon successful signup
        if (response.ok) {
            window.location.href = "/login"
        } else {
            alert("Failed to create account")
        }
    }
    catch (err) {
        console.error("Error:", err)
    }
})