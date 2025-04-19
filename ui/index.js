window.addEventListener('DOMContentLoaded', (ev) => {
    init();
    const refreshE = document.querySelector("#refresh");
    refreshE.addEventListener("click", (ev)=>{
        update();
    });
});

// sample data
const certsJson = `
[
    {
        "addr": "google.ca:443",
        "dns": "",
        "updateTime": "2025-1-10",
        "daysLeft": 30,
        "updateStatus": "ok",
    },
    {
        "addr": "1.2.3.4:443",
        "dns": "myexample.com",
        "updateTime": "2025-20-10",
        "daysLeft": 7,
        "updateStatus": "failed to connection."
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
    trE.innerHTML = `<td>${cert.addr}</td>
                    <td>${cert.dns}</td>
                    <td>${cert.updateTime}</td>
                    <td>${cert.daysLeft}</td>
                    <td>${cert.updateStatus}</td>`;

    return trE;
}

function certs2table(certs){
    const tableE = document.querySelector("#cert-table");
    tbodyE = tableE.children[1];
    tbodyE.innerHTML = "";
    for(let i=0; i<certs.length; i++){
        tbodyE.appendChild(cert2tr(certs[i]));
    }
    return tableE;
}

function certs2csv(certs){

}

function update_ui(data) {
    certs2table(data);
}

async function fetch_data() {
    const url = "./fetch";
    try {
        const response = await fetch(url);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error("Fetch error:", error);
        return null;
    }
}

async function update(){
    const data = await fetch_data();
    console.log(data);
    if(data == null){
        return;
    }
    update_ui(data);
}

async function init() {
    update();
}