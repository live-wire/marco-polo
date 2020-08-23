const width = window.innerWidth;

const height = window.innerHeight;

const projection = d3.geoMercator().fitSize([width, height], geoJSON);

const geoGenerator = d3.geoPath(projection);


const svg = d3.create("svg")
    .attr("width", width)
    .attr("height", height)
    .attr("viewbox", [0, 0, width, height])
    .attr("style", "background-color: #02101b");


svg.append("path")
    .datum(geoJSON)
    .attr("d", geoGenerator)
    .attr("fill", "#2a2b28")
    .attr("stroke", "none");


svg.append("path")
    .datum(traffic)
    .attr("d", geoGenerator)
    .attr("fill", "red")
    .attr("stroke", "none");

    

document.body.append(svg.node())
