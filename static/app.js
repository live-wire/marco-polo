const width = window.innerWidth;

const height = window.innerHeight;

const projection = d3.geoMercator().fitSize([width, height], geoJSON);

const geoGenerator = d3.geoPath(projection);


const svg = d3.select("body")
    .append("svg")
    .attr("width", width)
    .attr("height", height)
    .attr("viewBox", [150, 0, width / 1.25, height / 1.25 ])
    .style("background-color", "#02101b")

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

const traffic = svg.append("g")
    .attr("fill", "red")
    .attr("stroke", "none")
    .attr("cursor", "pointer");


setInterval(() => {
    fetchLiveTraffic().then(d=>{
        const data = d["default"]; // TODO: hardcoded default here (Use the selected service 

        traffic.selectAll("path")
            .data(data)
            .join("path")
            .attr("d", geoGenerator)
            .selectAll("title")
            .data((d, i) => data.slice(i, i+1))
            .join("title")
            .text(d => `${d.properties.ip}, ${d.properties.city}, ${d.properties.country}`)

    })
}, 3000);



