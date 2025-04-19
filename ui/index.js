window.addEventListener('DOMContentLoaded', (ev) => {
    init();
});

// sample data
const certsJson = `
[
    {
        "host": "1.2.3.4",
        "port": 443,
        "updateTime": "2025-1-10",
        "daysLeft": 30
    },
    {
        "host": "google.com",
        "port": 123,
        "updateTime": "2025-10-10",
        "daysLeft": 3
    },
    {
        "host": "baidu.com",
        "port": 443,
        "updateTime": "2025-20-10",
        "daysLeft": 7 
    }
]
`;

function cert2tr(cert) {
    const trE = document.createElement("tr");
    if (cert.daysLeft <= 3) {
        trE.className = "warning";
    } else if (cert.daysLeft <= 7) {
        trE.className = "caution";
    }
    trE.innerHTML = `<td>${cert.host}</td>
                    <td>${cert.port}</td>
                    <td>${cert.updateTime}</td>
                    <td>${cert.daysLeft}</td>`;

    return trE;
}

function certs2table(certs){
    const tableE = document.querySelector("#cert-table");
    for(let i=0; i<certs.length; i++){
        tableE.appendChild(cert2tr(certs[i]));
    }
    return tableE;
}

function certs2csv(certs){

}

function init_ui(data) {
    const certs = JSON.parse(certsJson);
    certs2table(certs);
}

async function fetch_data() {
    const url = "./fetch";
    try {
        const response = await fetch(url);
        const data = await response.json();
        console.log(data);
        return data;
    } catch (error) {
        console.error("Fetch error:", error);
        return null;
    }
}

function init() {
    const data = fetch_data();
    if(data == null){
        return;
    }

    init_ui(data);
}