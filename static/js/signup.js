const form = document.getElementById("signup-form");

form.addEventListener("submit", async(e) => {
    e.preventDefault();

    const userEmail = document.getElementById("email").value
    const userPass = document.getElementById("password").value
    const userPassConfirmed = document.getElementById("passwordConfirmed").value

    console.log(userPass);
    console.log(userPassConfirmed);

    if (userPass != userPassConfirmed) {
        alert("Passwords do not match");
        return;
    };

    const data = {
        email: userEmail,
        password: userPass,
        passwordConfirmed: userPassConfirmed,
    };

    try {
        const response = await fetch("/api/signup", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data)
        });

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