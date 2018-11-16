package yed

var (
	defaultNodeStyle = NodeStyle{
		Color:  "#FFCC00",
		Shape:  RoundedRectangle,
		Height: 30.0,
		Border: &defaultBorderStyle,
		Label:  &defaultLabelStyle,
	}
	defaultEdgeStyle = EdgeStyle{
		Source: NoArrow,
		Target: StdArrow,
		Line:   &defaultLineStyle,
		Label:  &defaultLabelStyle,
	}

	defaultBorderStyle = BorderStyle{
		Color: Black,
		Width: 1.0,
	}
	defaultLineStyle = LineStyle{
		Color: Black,
		Width: 1.0,
	}
	defaultLabelStyle = LabelStyle{
		FontSize: 12,
		Color:    Black,
	}
)

type Color string

const (
	Black = Color("#000000")
	White = Color("#FFFFFF")
)

type Shape string

const (
	RoundedRectangle = Shape("roundrectangle")
	Diamond          = Shape("diamond")
)

type NodeStyle struct {
	Color  Color
	Shape  Shape
	Height float64
	Border *BorderStyle
	Label  *LabelStyle
}

func (s *NodeStyle) setDefaults() {
	if s.Color == "" {
		s.Color = defaultNodeStyle.Color
	}
	if s.Shape == "" {
		s.Shape = defaultNodeStyle.Shape
	}
	if s.Height == 0 {
		s.Height = defaultNodeStyle.Height
	}
	if s.Border == nil {
		s.Border = &defaultBorderStyle
	} else {
		s.Border.setDefaults()
	}
	if s.Label == nil {
		s.Label = &defaultLabelStyle
	} else {
		s.Label.setDefaults()
	}
}

type Arrow string

const (
	NoArrow  = Arrow("none")
	StdArrow = Arrow("standard")
)

type EdgeStyle struct {
	Source Arrow
	Target Arrow
	Line   *LineStyle
	Label  *LabelStyle
}

func (s *EdgeStyle) setDefaults() {
	if s.Source == "" {
		s.Source = defaultEdgeStyle.Source
	}
	if s.Target == "" {
		s.Target = defaultEdgeStyle.Target
	}
	if s.Line == nil {
		s.Line = &defaultLineStyle
	} else {
		s.Line.setDefaults()
	}
	if s.Label == nil {
		s.Label = &defaultLabelStyle
	} else {
		s.Label.setDefaults()
	}
}

type BorderStyle = LineStyle

type LineStyle struct {
	Color Color
	Width float64
}

func (s *LineStyle) setDefaults() {
	if s.Color == "" {
		s.Color = defaultBorderStyle.Color
	}
	if s.Width == 0 {
		s.Width = defaultBorderStyle.Width
	}
}

type LabelStyle struct {
	FontSize int
	Color    Color
}

func (s *LabelStyle) setDefaults() {
	if s.FontSize == 0 {
		s.FontSize = defaultLabelStyle.FontSize
	}
	if s.Color == "" {
		s.Color = defaultLabelStyle.Color
	}
}
