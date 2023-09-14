package gograph

import (
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"text/template"
)

// Available graph types
type GraphType int

const (
	GraphTypeBar     GraphType = 1
	GraphTypeLine    GraphType = 2
	GraphTypeScatter GraphType = 3
)

type Graph struct {
	// Type of graph, eg line, bar, scatter
	GraphType GraphType `json:"graph_type"`

	// Graph x and y Axis
	XAxis *GraphAxis `json:"x_axis"`
	YAxis *GraphAxis `json:"y_axis"`

	// Dimensions in px
	Width  int64 `json:"width"`
	Height int64 `json:"height"`

	// Dataset for the graph
	DataSet *[]map[string]float64

	Settings *GraphSettings `json:"settings"`
}

type GraphAxis struct {
	// Key in map of values
	Key string `json:"key"`

	// axis label (if enabled in settings)
	Label string `json:"label"`

	// min and max
	// Min float64 `json:"min"`
	Max float64 `json:"max"`

	// increment for axis
	Increment float64 `json:"increment"`

	// Units to display on the axis
	Unit string `json:"unit"`

	Labels []float64
}

type GraphSettings struct {
	Labels bool `json:"labels,omitempty"`

	Fill       string `json:"fill,omitempty"`
	Background string `json:"background,omitempty"`
	Radius     int64  `json:"radius,omitempty"`
}

// Create new go graph!
func (g *Graph) New(exportPath string) (string, error) {
	switch g.GraphType {
	case 1:
		// ok... g, err :=
		g, err := g.newAxis(g.YAxis.Key, g.XAxis.Key, float64(g.Height))
		if err == nil {
			// w/h * h = w
			t := template.Must(template.New("").Funcs(template.FuncMap{
				"calcY": func(max float64, val float64) float64 {
					return max - val
				},
				"calcX": func(i int, increment float64) float64 {
					return float64(i) * increment
				},
				"labelY": func(height int64) float64 {
					return 40 - (float64(height) / 2)
				},
			},
			).Parse(`
			<div id="graph-container" style="position:absolute;top:50%;left:50%;transform:translate(-50%,-50%)">
				<div style="display:flex;flex-direction:row;">
					<div class="y-labels" style="position:relative;display:flex;height:{{.Height}}px;justify-content:space-between;flex-direction:column-reverse;width:30px;align-items:flex-end;margin-right:5px;">
						{{if not (eq .Settings nil)}}
						{{if .Settings.Labels}}
						<p class="axis-label" style="position:absolute;width:{{.Height}}px;top:50%;right:{{labelY .Height}}px;transform:translateY(-50%) rotate(-90deg);text-align:center;margin:0;">{{.YAxis.Label}}</p>
						{{end}}
						{{end}}
						{{range $l := .YAxis.Labels}}
						<p style="margin: 0;">{{$l}}</p>
						{{end}}
					</div>
					<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {{.XAxis.Max}} {{.YAxis.Max}}" width="{{.Width}}" height="{{.Height}}" preserveAspectRatio="none" {{if not (eq .Settings nil)}}fill="{{.Settings.Fill}}" style="background-color:{{.Settings.Background}}"{{end}}>       
						{{$xAxis := .XAxis}}
						{{$yAxis := .YAxis}}
						{{$settings := .Settings}}
						{{range $i, $d := .DataSet}}
						<rect width="{{$xAxis.Increment}}" height="{{index $d $yAxis.Key}}" x="{{calcX $i $xAxis.Increment}}" y="{{calcY $yAxis.Max (index $d $yAxis.Key)}}" {{if not (eq $settings nil)}}rx="{{$settings.Radius}}"{{end}}/>
						{{end}}
					</svg>				
				</div>
				<div class="x-labels" style="position:relative;display:flex;flex-direction:row;width:{{.Width}}px;justify-content:space-between;margin-left:30px;margin-top:5px;">
					{{if not (eq .Settings nil)}}
					{{if .Settings.Labels}}
					<p class="axis-label" style="position:absolute;width:{{.Width}}px;left:50%;transform:translateX(-50%);text-align:center;bottom:-25px;margin:0;">{{.XAxis.Label}}</p>
					{{end}}
					{{end}}
					{{range $i, $d := .DataSet}}
					<p style="width:calc(100% / {{len $xAxis.Labels}});text-align:center;margin:0;">{{index $d $xAxis.Key}}</p>
					{{end}}
				</div>
			</div>
			`))
			f, err := os.Create(exportPath + "/testing.html")
			if err != nil {
				return "", errors.New("create file: " + err.Error())
			}
			err = t.Execute(f, g)
			if err != nil {
				return "", errors.New("execute: " + err.Error())
			}
			f.Close()
		}
	case 2:
	case 3:
	default:
		return "", errors.New(fmt.Sprint("unknown graph type %v", g.GraphType))
	}
	return exportPath, nil
}

// Generate friendly axis from dataset
func (g *Graph) newAxis(Ykey string, XKey string, dimension float64) (*Graph, error) {
	yKeys := []float64{}
	xKeys := []float64{}

	if g == nil {
		return g, errors.New("declared graph is nill")
	}
	if g.DataSet == nil {
		return g, errors.New("graph dataset is empty")
	}

	// Grab all values in dataset and order them
	for _, d := range *g.DataSet {
		yKeys = append(yKeys, d[Ykey])

		xKeys = append(xKeys, d[XKey])
		g.XAxis.Labels = append(g.XAxis.Labels, d[XKey])
	}
	sort.Float64s(yKeys)
	sort.Float64s(xKeys)

	// Get the increment factor (1:10 = 1, 10:100 = 10, 100:1000 = 100) etc
	max := math.Ceil(math.Log10(yKeys[len(yKeys)-1]*1.15)) - 1

	// Set max value around 20% above the top value rounded to the closest multiple of the increment
	g.YAxis.Max = roundToMult(yKeys[len(yKeys)-1]*1.15, math.Pow(10, max))
	g.YAxis.Increment = g.YAxis.Max / (dimension / 40)

	// Generate the labels for the Y axis following the increment
	g.YAxis.Labels = *increment(g.YAxis.Increment, int(dimension/40))

	g.XAxis.Max = float64(g.Width) / float64(g.Height) * g.YAxis.Max
	g.XAxis.Increment = g.XAxis.Max / float64(len(xKeys))

	return g, nil
}

// ------------------
// UTILITY FUNCTIONS
// ------------------

func roundToMult(number float64, multiple float64) float64 {
	return math.Ceil(number/multiple) * multiple
}

func increment(increment float64, max int) *[]float64 {
	labels := []float64{}
	for i := 0; i <= max; i++ {
		labels = append(labels, float64(i)*increment)
	}
	return &labels
}

func reverse[T any](s []T) []T {
	f := 0
	l := len(s) - 1
	for f < l {
		s[f], s[l] = s[l], s[f]
		f++
		l--
	}
	return s
}
