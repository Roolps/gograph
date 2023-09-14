package gograph

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"text/template"
	"time"
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

	Settings *graphSettings `json:"settings"`
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

type graphSettings struct {
	Labels bool `json:"labels,omitempty"`
}

// Create new go graph!
func (g *Graph) New(exportPath string) (string, error) {
	switch g.GraphType {
	case 1:
		// ok... g, err :=
		g, err := g.newAxis(g.YAxis.Key, float64(g.Height))
		if err == nil {
			t := template.Must(template.New("").Parse(`
			<svg viewBox="0 0 {{.Width}} {{.Height}}" width="{{.Width}}" height="{{.Height}}" fill="red" preserveAspectRatio="none" xmlns="http://www.w3.org/2000/svg">       
			
			{{range $l := .YAxis.Labels}}<text>{{$l}}</text>{{end}}
			
			</svg>`))
			f, err := os.Create(exportPath + "/" + fmt.Sprintf("%x", time.Now().Unix()) + ".svg")
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
func (g *Graph) newAxis(key string, dimension float64) (*Graph, error) {
	keys := []float64{}

	if g == nil {
		return g, errors.New("declared graph is nill")
	}
	if g.DataSet == nil {
		return g, errors.New("graph dataset is empty")
	}

	// Grab all values in dataset and order them
	for _, d := range *g.DataSet {
		keys = append(keys, d[key])
	}
	sort.Float64s(keys)

	// Get the increment factor (1:10 = 1, 10:100 = 10, 100:1000 = 100) etc
	log.Println(keys[len(keys)-1])
	max := math.Ceil(math.Log10(keys[len(keys)-1]*1.2)) - 1

	g.YAxis.Max = roundToMult(keys[len(keys)-1]*1.2, math.Pow(10, max))
	g.YAxis.Increment = g.YAxis.Max / (dimension / 40)

	g.YAxis.Labels = *increment(g.YAxis.Increment, int(dimension/40))
	log.Println(g.YAxis.Labels)
	return g, nil
}

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
