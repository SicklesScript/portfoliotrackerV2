const form = document.getElementById('login-form');

// Listen for submit event
form.addEventListener('submit', async (e) => {
    e.preventDefault();

    // Extract data from form
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    const data = {
        email: email,
        password: password
    };

    try {
        // Pass data to loginAPI
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            // Convert data to json friendly format
            body: JSON.stringify(data)
        });

        // Store token and user_email in header
        if (response.ok) {
            const info = await response.json();

            localStorage.setItem('token', info.token);

            localStorage.setItem('user_email', info.email);
            
            // Redirect user to home screen *for now*
            window.location.href = '/'
        } else {
            alert("Login failed!")
        }
    } catch (err) {
        console.error("Error:", err)
    }
});