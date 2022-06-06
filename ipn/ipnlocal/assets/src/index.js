let data 

fetch('/json')
    .then(response => response.json())
    .then(data => {
        // console.log(data)
        newGrid(data.Peers)
    });

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

newGrid = d => {
    new gridjs.Grid({
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
                name: 'IPs'
            }, 
            {
                id: 'OS',
                name: 'OS'
            },
            {
                id: 'ActAgo',
                name: 'Last Seen'
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
        search: true,
        data: d
    }).render(document.getElementById("wrapper"));
}

// newGrid([])
