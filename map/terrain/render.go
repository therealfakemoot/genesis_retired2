package genesis

import (
	"html/template"
	"io"

	"github.com/spf13/viper"
	Q "github.com/therealfakemoot/go-quantize"

	l "github.com/therealfakemoot/genesis/log"
)

var topoMap = `
<!DOCTYPE html>
<svg width="{{ $.Width }}" height="{{ $.Height }}" stroke="#fff" stroke-width="0.5"></svg>
<script src="https://d3js.org/d3.v4.min.js"></script>
<script src="https://d3js.org/d3-hsv.v0.1.min.js"></script>
<script src="https://d3js.org/d3-contour.v1.min.js"></script>
<script>

var svg = d3.select("svg"),
width = +svg.attr("width"),
height = +svg.attr("height");

var i0 = d3.interpolateHsvLong(d3.hsv(120, 1, 0.65), d3.hsv(60, 1, 0.90)),
i1 = d3.interpolateHsvLong(d3.hsv(60, 1, 0.90), d3.hsv(0, 0, 0.95)),
interpolateTerrain = function(t) { return t < 0.5 ? i0(t * 2) : i1((t - 0.5) * 2); },
color = d3.scaleSequential(interpolateTerrain).domain([{{ $.Domain.Min }}, {{ $.Domain.Max }}]);

d3.json("terrain.json", function(error, terrain) {
	if (error) throw error;

	svg.selectAll("path")
	.dev(true)
	.data(d3.contours()
	.size([terrain.width, terrain.height])
	.thresholds(d3.range({{ $.Domain.Max }}, {{ $.Domain.Max }}, 5))
	(terrain.values))
	.enter().append("path")
	.attr("d", d3.geoPath(d3.geoIdentity().scale(width / terrain.width)))
	.attr("fill", function(d) { return color(d.value); });
});

</script>
`

// RenderTopoHTML emits
func RenderTopoHTML(w io.Writer) {
	t, err := template.New("terrain").Parse(topoMap)

	if err != nil {
		l.Term.WithError(err).Error("Failed to parse topoMap template.")
	}

	d := Q.Domain{
		Max:  viper.GetFloat64("Terrain.Domain.Max"),
		Min:  viper.GetFloat64("Terrain.Domain.Min"),
		Step: viper.GetFloat64("Terrain.Domain.Step"),
	}

	v := struct {
		Width  float64
		Height float64
		Domain Q.Domain
	}{
		Width:  float64(viper.GetInt("mapX")),
		Height: float64(viper.GetInt("mapY")),
		Domain: d,
	}

	t.Execute(w, v)
}
