document.addEventListener('DOMContentLoaded', async () => {
    const token = localStorage.getItem("token");
    const selectElement = document.getElementById("portfolio_id");

    try {
        const response = await fetch("/api/getPortfolio", {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        if (response.ok) {
            const portfolios = await response.json();
            
            selectElement.innerHTML = '<option value="" disabled selected>Select a portfolio</option>';

            // Loop through portfolios and make an <option> for each
            portfolios.forEach(portfolio => {
                const option = document.createElement("option");
                option.value = portfolio.ID;    
                option.textContent = portfolio.Name; 
                selectElement.appendChild(option);
            });
        } else {
            console.error("Failed to fetch portfolios");
            selectElement.innerHTML = '<option value="" disabled>Error loading portfolios</option>';
        }
    } catch (err) {
        console.error("Error:", err);
    }
});

const form = document.getElementById("transaction-form");

form.addEventListener("submit", async(e) => {
    e.preventDefault();

    // Gather values from the DOM
    const portfolioId = document.getElementById("portfolio_id").value;
    const stockName = document.getElementById("stock_name").value;
    const ticker = document.getElementById("ticker").value;
    // Parse numbers to ensure they aren't sent as strings 
    const shares = document.getElementById("shares").value;
    const pricePerShare = document.getElementById("price_per_share").value;

    const token = localStorage.getItem("token");

    const data = {
        portfolio_id: portfolioId,
        stock_name: stockName,
        ticker: ticker,
        shares: shares, 
        price_per_share: pricePerShare
    };

    try {
        const response = await fetch("/api/createTransaction", {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        if (response.ok) {
            alert("Transaction successfully added!");
            // Redirect to dashboard or clear form
            window.location.href = "/"; 
        } else {
            alert("Unable to add transaction");
            console.log(await response.text()); 
        }
    } catch (err) {
        console.error("Error:", err);
    }
});