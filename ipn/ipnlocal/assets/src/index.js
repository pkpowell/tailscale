let data 

fetch('/json')
    .then(response => response.json())
    .then(data => {
        console.log(data)
        d = data.Peers.map(x=>{
            console.log([x.HostName, x.OS])
            return [
                x.HostName, 
                x.IPv4,
                x.OS
            ]
        })
        newGrid(d)
    });

newGrid = d => {
    new gridjs.Grid({
        columns: ["Machine", "IP", "OS"],
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
