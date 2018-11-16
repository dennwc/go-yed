package yed

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

type ID int

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

type Writer struct {
	w      io.Writer
	err    error
	header bool
	footer bool

	lastNode ID
	lastEdge ID
}

func (w *Writer) writeString(s string) error {
	if w.err != nil {
		return w.err
	}
	_, err := w.w.Write([]byte(s))
	if err != nil {
		w.err = err
	}
	return err
}

func (w *Writer) printf(format string, args ...interface{}) error {
	if w.err != nil {
		return w.err
	}
	_, err := fmt.Fprintf(w.w, format, args...)
	if err != nil {
		w.err = err
	}
	return err
}

func (w *Writer) writeHeader() error {
	if w.footer {
		return errors.New("yed: writer is closed")
	} else if w.header {
		return w.err
	}
	if err := w.writeString(header); err != nil {
		return err
	}
	w.header = true
	return nil
}

func escape(s string) string {
	buf := bytes.NewBuffer(nil)
	xml.EscapeText(buf, []byte(s))
	return buf.String()
}

type Node struct {
	Label       string
	Description string

	Style *NodeStyle
}

func (w *Writer) WriteNode(n Node) (ID, error) {
	id := w.lastNode
	w.lastNode++
	if err := w.WriteNodeWithID(id, n); err != nil {
		return 0, err
	}
	return id, nil
}

func (w *Writer) WriteNodeWithID(id ID, n Node) error {
	if err := w.writeHeader(); err != nil {
		return err
	}
	_ = w.printf("\n\t<node id=\"n%d\">", id)
	if err := w.writeDescription(n.Description); err != nil {
		return err
	}
	if err := w.writeNodeGraphics(n.Label, n.Style); err != nil {
		return err
	}
	_ = w.writeString("\n\t</node>")
	return w.err
}

func (w *Writer) writeDescription(desc string) error {
	if desc == "" {
		return w.writeString("\n\t\t<data key=\"d5\"/>")
	}
	return w.printf("\n\t\t<data key=\"d5\" xml:space=\"preserve\">%s</data>", escape(desc))
}

func (w *Writer) writeNodeGraphics(label string, s *NodeStyle) error {
	if s == nil {
		s = &defaultNodeStyle
	} else {
		s.setDefaults()
	}
	bs := s.Border
	_ = w.printf(`
      <data key="d6">
        <y:ShapeNode>
          <y:Geometry height="%f" x="0.0" y="0.0"/>
          <y:Fill color="%s" transparent="false"/>
          <y:BorderStyle color="%s" raised="false" type="line" width="%f"/>
`, s.Height, string(s.Color), string(bs.Color), bs.Width)
	if label != "" {
		ls := s.Label
		_ = w.printf(`
          <y:NodeLabel alignment="center" autoSizePolicy="content" fontFamily="Dialog" fontSize="%d" fontStyle="plain" hasBackgroundColor="false" hasLineColor="false" height="17.96875" horizontalTextPosition="center" iconTextGap="4" modelName="custom" textColor="%s" verticalTextPosition="bottom" visible="true" x="5.0" xml:space="preserve" y="6.015625">%s<y:LabelModel><y:SmartNodeLabelModel distance="4.0"/></y:LabelModel><y:ModelParameter><y:SmartNodeLabelModelParameter labelRatioX="0.0" labelRatioY="0.0" nodeRatioX="0.0" nodeRatioY="0.0" offsetX="0.0" offsetY="0.0" upX="0.0" upY="-1.0"/></y:ModelParameter></y:NodeLabel>`,
			ls.FontSize, string(ls.Color), escape(label))
	}
	_ = w.printf(`
          <y:Shape type="%s"/>
        </y:ShapeNode>
      </data>`, string(s.Shape))
	return w.err
}

type Edge struct {
	Label       string
	Description string

	Style *EdgeStyle
}

func (w *Writer) WriteEdge(from, to ID, e *Edge) error {
	if err := w.writeHeader(); err != nil {
		return err
	}
	if e == nil {
		e = &Edge{}
	}
	id := w.lastEdge
	w.lastEdge++

	_ = w.printf("\n\t<edge id=\"e%d\" source=\"n%d\" target=\"n%d\">", id, from, to)
	if err := w.writeDescription(e.Description); err != nil {
		return err
	}
	if err := w.writeEdgeGraphics(e.Label, e.Style); err != nil {
		return err
	}
	_ = w.writeString("\n\t</edge>")
	return w.err
}

func (w *Writer) writeEdgeGraphics(label string, s *EdgeStyle) error {
	if s == nil {
		s = &defaultEdgeStyle
	} else {
		s.setDefaults()
	}
	bs := s.Line
	_ = w.printf(`
      <data key="d10">
        <y:PolyLineEdge>
          <y:Path sx="0.0" sy="0.0" tx="0.0" ty="0.0"/>
          <y:LineStyle color="%s" type="line" width="%f"/>
          <y:Arrows source="%s" target="%s"/>
`, string(bs.Color), bs.Width, string(s.Source), string(s.Target))
	if label != "" {
		ls := s.Label
		_ = w.printf(`
          <y:EdgeLabel alignment="center" configuration="AutoFlippingLabel" distance="2.0" fontFamily="Dialog" fontSize="%d" fontStyle="plain" hasBackgroundColor="false" hasLineColor="false" height="17.96875" horizontalTextPosition="center" iconTextGap="4" modelName="three_center" modelPosition="center" preferredPlacement="anywhere" ratio="0.5" textColor="%s" verticalTextPosition="bottom" visible="true" x="-2.14501953125" xml:space="preserve" y="11.50439453125">%s<y:PreferredPlacementDescriptor angle="0.0" angleOffsetOnRightSide="0" angleReference="absolute" angleRotationOnRightSide="co" distance="-1.0" frozen="true" placement="anywhere" side="anywhere" sideReference="relative_to_edge_flow"/></y:EdgeLabel>`,
			ls.FontSize, string(ls.Color), escape(label))
	}
	_ = w.writeString(`
          <y:BendStyle smoothed="false"/>
        </y:PolyLineEdge>
      </data>`)
	return w.err
}

func (w *Writer) Close() error {
	if w.footer {
		return nil
	}
	if err := w.writeString(footer); err != nil {
		return err
	}
	w.footer = true
	return nil
}

const header = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<graphml xmlns="http://graphml.graphdrawing.org/xmlns" xmlns:java="http://www.yworks.com/xml/yfiles-common/1.0/java" xmlns:sys="http://www.yworks.com/xml/yfiles-common/markup/primitives/2.0" xmlns:x="http://www.yworks.com/xml/yfiles-common/markup/2.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:y="http://www.yworks.com/xml/graphml" xmlns:yed="http://www.yworks.com/xml/yed/3" xsi:schemaLocation="http://graphml.graphdrawing.org/xmlns http://www.yworks.com/xml/schema/graphml/1.1/ygraphml.xsd">
  <!--Created by github.com/dennwc/go-yed -->
  <key attr.name="Description" attr.type="string" for="graph" id="d0"/>
  <key for="port" id="d1" yfiles.type="portgraphics"/>
  <key for="port" id="d2" yfiles.type="portgeometry"/>
  <key for="port" id="d3" yfiles.type="portuserdata"/>
  <key attr.name="url" attr.type="string" for="node" id="d4"/>
  <key attr.name="description" attr.type="string" for="node" id="d5"/>
  <key for="node" id="d6" yfiles.type="nodegraphics"/>
  <key for="graphml" id="d7" yfiles.type="resources"/>
  <key attr.name="url" attr.type="string" for="edge" id="d8"/>
  <key attr.name="description" attr.type="string" for="edge" id="d9"/>
  <key for="edge" id="d10" yfiles.type="edgegraphics"/>
  <graph edgedefault="directed" id="G">
    <data key="d0"/>
`

const footer = `
  </graph>
  <data key="d7">
    <y:Resources/>
  </data>
</graphml>
`
