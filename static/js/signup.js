const form = document.getElementById("signup-form");

form.addEventListener("submit", async(e) => {
    e.preventDefault();

    const userEmail = document.getElementById("email")
    const userPass = document.getElementById("password")
    const userPassConfirmed = document.getElementById("passwordConfirmed")

    const data = {
        email: userEmail,
        password: userPass,
        passwordConfirmed: userPassConfirmed,
    };

    try {
        const response = await fetch("/api/signup", {
            method: "POST",
            headers: {
                "Cotent-Type": "application/json",
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            window.location.href("/login")
        } else {
            alert("Failed to create account")
        }
    }
    catch (err) {
        console.error("Error:", err)
    }
})