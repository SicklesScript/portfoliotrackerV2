document.addEventListener('DOMContentLoaded', async () => {
    const token = localStorage.getItem("token");
    const tableBody = document.getElementById("dashboard-table-body");
    
    // Summary Elements
    const elTotal = document.getElementById("summary-total");
    const elRoe = document.getElementById("summary-roe");
    const elRoic = document.getElementById("summary-roic");

    if (!token) {
        window.location.href = "/login.html";
        return;
    }

    try {
        // Fetch user's holdings from DB
        const response = await fetch("/api/dashboard", {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            console.error("Failed to load dashboard");
            return;
        }

        const holdings = await response.json();
        tableBody.innerHTML = "";

        if (holdings.length === 0) {
            tableBody.innerHTML = "<tr><td colspan='7'>No holdings found. Add a transaction!</td></tr>";
            elTotal.textContent = "$0.00";
            return;
        }

        let grandTotalInvested = 0;

        const rowPromises = holdings.map(async (item) => {
            const row = document.createElement("tr");

            // Strings to floats)
            const shares = parseFloat(item.TotalShares);
            const avgPrice = parseFloat(item.AvgCostPerShare);
            const totalVal = parseFloat(item.TotalInvested);
            
            grandTotalInvested += totalVal;

            // Unique IDs for cells to update later
            const roeId = `roe-${item.Ticker}`;
            const roicId = `roic-${item.Ticker}`;
            const crId = `cr-${item.Ticker}`;

            // Initial Table Row
            row.innerHTML = `
                <td><strong>${item.Ticker}</strong></td>
                <td>${shares.toFixed(2)}</td>
                <td>$${avgPrice.toFixed(2)}</td>
                <td>$${totalVal.toFixed(2)}</td>
                <td id="${roeId}">Loading...</td>
                <td id="${roicId}">Loading...</td>
                <td id="${crId}">Loading...</td>
            `;
            tableBody.appendChild(row);

            // Fetch specific metrics for this stock
            let data = { roe: 0, roic: 0 };
            
            try {
                const res = await fetch(`/api/metrics?ticker=${item.Ticker}`);
                if (res.ok) {
                    const metrics = await res.json();
                    
                    // Store for weighting calculation
                    data.roe = metrics.returnOnEquity || 0;
                    data.roic = metrics.returnOnInvestedCapital || 0;
                    const cr = metrics.currentRatio || 0;

                    // Update Table Cells
                    document.getElementById(roeId).textContent = (data.roe * 100).toFixed(2) + "%";
                    document.getElementById(roicId).textContent = (data.roic * 100).toFixed(2) + "%";
                    document.getElementById(crId).textContent = cr.toFixed(2);
                }
            } catch (err) {
                console.error(`Metrics error for ${item.Ticker}`, err);
                document.getElementById(roeId).textContent = "-";
                document.getElementById(roicId).textContent = "-";
                document.getElementById(crId).textContent = "-";
            }

            // Return value and metrics for the summary calculation
            return {
                value: totalVal,
                roe: data.roe,
                roic: data.roic
            };
        });

        // Wait for all API calls to finish
        const results = await Promise.all(rowPromises);

        // Calculate Weighted Averages
        let weightedRoeSum = 0;
        let weightedRoicSum = 0;

        results.forEach(stock => {
            if (grandTotalInvested > 0) {
                const weight = stock.value / grandTotalInvested;
                weightedRoeSum += (weight * stock.roe);
                weightedRoicSum += (weight * stock.roic);
            }
        });

        // Update summary headers
        elTotal.textContent = "$" + grandTotalInvested.toLocaleString(undefined, {
            minimumFractionDigits: 2, 
            maximumFractionDigits: 2
        });
        elRoe.textContent = (weightedRoeSum * 100).toFixed(2) + "%";
        elRoic.textContent = (weightedRoicSum * 100).toFixed(2) + "%";

    } catch (err) {
        console.error("Dashboard Error:", err);
    }
});