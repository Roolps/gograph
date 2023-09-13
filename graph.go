package gograph

import (
	"errors"
	"fmt"
)

// Available graph types
type GraphType int

const (
	GraphTypeLine    GraphType = 1
	GraphTypeBar     GraphType = 2
	GraphTypeScatter GraphType = 3
)

type Graph struct {
	// Type of graph, eg line, bar, scatter
	GraphType GraphType `json:"graph_type"`

	// Graph x and y Axis
	XAxis *graphAxis `json:"x_axis"`
	YAxis *graphAxis `json:"y_axis"`

	// Dimensions in px
	Width  int64 `json:"width"`
	Height int64 `json:"height"`

	Settings *graphSettings `json:"settings"`
}

type graphAxis struct {
	// axis label (if enabled in settings)
	Label string `json:"label"`

	// min and max
	Min int64 `json:"min"`
	Max int64 `json:"max"`

	// increment for axis
	Increment int64 `json:"increment"`
}

type graphSettings struct {
	Labels bool `json:"labels,omitempty"`
}

// New gograph!
func (g *Graph) New() (string, error) {
	switch g.GraphType {
	case 1:
	case 2:
	case 3:
	default:
		return "", errors.New(fmt.Sprint("unknown graph type %v", g.GraphType))
	}
	return "", nil
}
