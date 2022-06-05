let data 

fetch('/json')
    .then(response => response.json())
    .then(data => {
        console.log(data)
        // d = data.Peers.map(x=>{
        //     console.log([x.HostName, x.OS])
        //     return [
        //         x.HostName, 
        //         x.IPv4,
        //         x.OS
        //     ]
        // })
        newGrid(data.Peers)
    });

newGrid = d => {
    new gridjs.Grid({
        columns: [
            {
                id: 'HostName',
                name: 'Hostname'
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
        data: d
        // data: [
        //   ["John", "john@example.com", "(353) 01 222 3333"],
        //   ["Mark", "mark@gmail.com", "(01) 22 888 4444"],
        //   ["Eoin", "eoin@gmail.com", "0097 22 654 00033"],
        //   ["Sarah", "sarahcdd@gmail.com", "+322 876 1233"],
        //   ["Afshin", "afshin@mail.com", "(353) 22 87 8356"]
        // ]
    }).render(document.getElementById("wrapper"));
}

// newGrid([])
