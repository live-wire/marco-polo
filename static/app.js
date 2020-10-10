const width = window.innerWidth;

const height = window.innerHeight;

const projection = d3.geoMercator().fitSize([width, height], geoJSON);

const geoGenerator = d3.geoPath(projection);



const svg = d3.create("svg")
    .attr("id", "svgd3")
    .attr("width", width)
    .attr("height", height)
    .attr("viewBox", [150, 0, width / 1.25, height / 1.25 ])
    .attr("style", "background-color: #02101b")

const countries = svg.append("g")
    .attr("fill", "#2a2b28")
    .attr("stroke", "none")
    .attr("cursor", "pointer")
    .selectAll("path")
    .data(geoJSON.features)
    .join("path")
    .attr("d", geoGenerator);

countries.append("title")
    .text(d => d.properties["NAME"])


const trafficPath = svg.append("path")
    .attr("fill", "red")
    .attr("stroke", "none");


setInterval(() => {
    const traffic = {"type": "FeatureCollection", "features": []}
    
    fetchLiveTraffic().then(d=>{
        traffic["features"] = d["default"]; // TODO: hardcoded default here (Use the selected service 
    }).then(() => {
        trafficPath.datum(traffic)
        .attr("d", geoGenerator)   
    })
}, 3000)    

    
document.body.append(svg.node());

