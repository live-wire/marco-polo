/**
 * Sample traffic data
 */

// Just flush/ returns traffic for all services
// TODO: decide which service you want the traffic for using a dropdown maybe?
const url = "/flush/default" 

// Default traffic object should be GeoJSON parseable 
var traffic = {"type": "FeatureCollection", "features": []}

fetchLiveTraffic = () => {
    fetch("/flush").then(d=>d.json()).then(d=>{
        traffic["features"] = d["default"]; // TODO: hardcoded default here (Use the selected service name here)
        console.log(traffic);
    })
}
fetchLiveTraffic()
setInterval(fetchLiveTraffic, 3000)
