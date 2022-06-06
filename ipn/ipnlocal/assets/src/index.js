import { Grid, html } from "/assets/src/gridjs.js"
let data 

const compReg = (a, b) => {
    a = a.toLowerCase()
    b = b.toLowerCase()
    if (a > b) {
        return 1;
    } else if (b > a) {
        return -1;
    } else {
        return 0;
    }
}

const newGrid = d => {
    const options = { weekday: 'short', year: 'numeric', month: 'short', day: 'numeric' }
    new Grid({
        columns: [
            {
                id: 'HostName',
                name: 'Hostname',
                sort: {
                    compare: compReg
                }
            }, 
            {
                id: 'IPs',
                name: 'IPs',
                // formatter: cell => `${cell.join(", ")}`
                formatter: html(cell => cell.forEach(c => `<div>${c}</div>`))
            }, 
            {
                id: 'OS',
                name: 'OS'
            },
            {
                id: 'LastSeen',
                name: 'Last Seen',
                formatter: cell => new Date(cell).toLocaleDateString("en-US", options)
            }, 
            {
                id: 'DNSName',
                name: 'DNS'
            }, 
            {
                id: 'RX',
                name: 'RX'
            },
            {
                id: 'TX',
                name: 'TX'
            },
        ],
        sort: true,
        // search: true,
        data: d
    }).render(document.getElementById("wrapper"));
}

fetch('/json')
    .then(response => response.json())
    .then(data => {
        console.log(data)
        newGrid(data.Peers)
    })