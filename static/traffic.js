/**
 * Sample traffic data
 */

// Just flush/ returns traffic for all services
// TODO: decide which service you want the traffic for using a dropdown maybe?
const url = "/flush/default" 

// Default traffic object should be GeoJSON parseable 

const fetchLiveTraffic = () => fetch(url).then(d=>d.json());
