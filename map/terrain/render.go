package genesis

// import "text/template"

var topoMap = `
<!DOCTYPE html>
<svg width="960" height="673" stroke="#fff" stroke-width="0.5"></svg>
<script src="https://d3js.org/d3.v4.min.js"></script>
<script src="https://d3js.org/d3-hsv.v0.1.min.js"></script>
<script src="https://d3js.org/d3-contour.v1.min.js"></script>
<script>

var svg = d3.select("svg"),
width = +svg.attr("{{ .mapWidth }}"),
height = +svg.attr("{{ .mapHeigth }}");

var i0 = d3.interpolateHsvLong(d3.hsv(120, 1, 0.65), d3.hsv(60, 1, 0.90)),
i1 = d3.interpolateHsvLong(d3.hsv(60, 1, 0.90), d3.hsv(0, 0, 0.95)),
interpolateTerrain = function(t) { return t < 0.5 ? i0(t * 2) : i1((t - 0.5) * 2); },
color = d3.scaleSequential(interpolateTerrain).domain([90, 190]);

d3.json("volcano.json", function(error, volcano) {
	if (error) throw error;

	svg.selectAll("path")
	.data(d3.contours()
	.size([volcano.width, volcano.height])
	.thresholds(d3.range(90, 195, 5))
	(volcano.values))
	.enter().append("path")
	.attr("d", d3.geoPath(d3.geoIdentity().scale(width / volcano.width)))
	.attr("fill", function(d) { return color(d.value); });
});

</script>
`

type TopoParams struct {
	Width  int
	Height int
}

func RenderTopoHtml(width, height int) {}
