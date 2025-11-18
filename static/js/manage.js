const form = document.getElementById("create-portfolio-form");

form.addEventListener("submit", async(e) => {
    e.preventDefault();

    const portfolioName = document.getElementById("portfolioName").value;
    const token = localStorage.getItem("token");

    const data = {
        portfolioName: portfolioName
    };

    try {
        const response = await fetch("/api/createPortfolio", {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            alert("Portfolio successfully created!");
            window.location.href = "/";
        } else {
            alert("Unable to create portfolio");
        }
    } catch (err) {
        console.error("Error:", err)
    }
});
