const form = document.getElementById('login-form');

form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const userEmail = document.getElementById('email').value;
    const userPassword = document.getElementById('password').value;

    const data = {
        email: userEmail,
        password: userPassword
    };

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            const info = await response.json();

            localStorage.setItem('token', info.token);

            localStorage.setItem('user_email', info.email);

            window.location.href = '/'
        } else {
            alert("Login failed!")
        }
    } catch (err) {
        console.error("Error:", err)
    }
});