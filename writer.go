package yed

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

func newWriter(w io.Writer) *writer {
	return &writer{w: w}
}

type writer struct {
	w      io.Writer
	err    error
	header bool
	footer bool
}

func (w *writer) writeString(s string) error {
	if w.err != nil {
		return w.err
	}
	_, err := w.w.Write([]byte(s))
	if err != nil {
		w.err = err
	}
	return err
}

func (w *writer) printf(format string, args ...interface{}) error {
	if w.err != nil {
		return w.err
	}
	_, err := fmt.Fprintf(w.w, format, args...)
	if err != nil {
		w.err = err
	}
	return err
}

func (w *writer) writeHeader() error {
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

func (w *writer) WriteGraph(g *Graph) error {
	if err := w.writeHeader(); err != nil {
		return err
	}
	id := g.ID()
	if id == "" {
		id = "G"
	}
	_ = w.printf(`
  <graph edgedefault="directed" id="%s">
    <data key="d0"/>`, id)
	for _, n := range g.sub {
		if err := w.writeNode(n); err != nil {
			return err
		}
	}
	return w.writeString(`
  </graph>`)
}

func (w *writer) writeNode(n *Node) error {
	if n.sub != nil {
		_ = w.printf("\n\t<node id=\"%s\" yfiles.foldertype=\"group\">", n.ID())
		if err := w.writeDescription(n.Description); err != nil {
			return err
		}
		_ = w.printf(`
      <data key="d6">
        <y:ProxyAutoBoundsNode>
          <y:Realizers active="0">
            <y:GroupNode>
              <y:Geometry height="181.87067499999998" width="294.52628859375017" x="186.07787140624984" y="135.68900499999998"/>
              <y:Fill color="#F5F5F5" transparent="false"/>
              <y:BorderStyle color="#000000" type="dashed" width="1.0"/>
              <y:NodeLabel alignment="right" autoSizePolicy="node_width" backgroundColor="#EBEBEB" borderDistance="0.0" fontFamily="Dialog" fontSize="15" fontStyle="plain" hasLineColor="false" height="21.4609375" horizontalTextPosition="center" iconTextGap="4" modelName="internal" modelPosition="t" textColor="#000000" verticalTextPosition="bottom" visible="true" width="294.52628859375017" x="0.0" xml:space="preserve" y="0.0">%s</y:NodeLabel>
              <y:Shape type="roundrectangle"/>
              <y:State closed="false" closedHeight="50.0" closedWidth="50.0" innerGraphDisplayEnabled="false"/>
              <y:NodeBounds considerNodeLabelSize="true"/>
              <y:Insets bottom="15" bottomF="15.0" left="15" leftF="15.0" right="15" rightF="15.0" top="15" topF="15.0"/>
              <y:BorderInsets bottom="5" bottomF="5.1200000000000045" left="7" leftF="7.15380859375" right="7" rightF="7.040000000000134" top="0" topF="0.0"/>
            </y:GroupNode>
            <y:GroupNode>
              <y:Geometry height="50.0" width="50.0" x="0.0" y="60.0"/>
              <y:Fill color="#F5F5F5" transparent="false"/>
              <y:BorderStyle color="#000000" type="dashed" width="1.0"/>
              <y:NodeLabel alignment="right" autoSizePolicy="node_width" backgroundColor="#EBEBEB" borderDistance="0.0" fontFamily="Dialog" fontSize="15" fontStyle="plain" hasLineColor="false" height="21.4609375" horizontalTextPosition="center" iconTextGap="4" modelName="internal" modelPosition="t" textColor="#000000" verticalTextPosition="bottom" visible="true" width="64.3076171875" x="-7.15380859375" xml:space="preserve" y="0.0">%s</y:NodeLabel>
              <y:Shape type="roundrectangle"/>
              <y:State closed="true" closedHeight="50.0" closedWidth="50.0" innerGraphDisplayEnabled="false"/>
              <y:Insets bottom="5" bottomF="5.0" left="5" leftF="5.0" right="5" rightF="5.0" top="5" topF="5.0"/>
              <y:BorderInsets bottom="0" bottomF="0.0" left="0" leftF="0.0" right="0" rightF="0.0" top="0" topF="0.0"/>
            </y:GroupNode>
          </y:Realizers>
        </y:ProxyAutoBoundsNode>
      </data>`, escape(n.Label), escape(n.Label))
		if err := w.WriteGraph(n.sub); err != nil {
			return err
		}
		return w.writeString("\n\t</node>")

	}
	_ = w.printf("\n\t<node id=\"%s\">", n.ID())
	if err := w.writeDescription(n.Description); err != nil {
		return err
	}
	if err := w.writeNodeGraphics(n.Label, n.Style); err != nil {
		return err
	}
	return w.writeString("\n\t</node>")
}

func (w *writer) writeDescription(desc string) error {
	if desc == "" {
		return w.writeString("\n\t\t<data key=\"d5\"/>")
	}
	return w.printf("\n\t\t<data key=\"d5\" xml:space=\"preserve\">%s</data>", escape(desc))
}

func (w *writer) writeNodeGraphics(label string, s *NodeStyle) error {
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

func (w *writer) WriteEdge(e *Edge) error {
	_ = w.printf("\n\t<edge id=\"%s\" source=\"%s\" target=\"%s\">",
		e.ID(), e.Source().ID(), e.Target().ID())
	if err := w.writeDescription(e.Description); err != nil {
		return err
	}
	if err := w.writeEdgeGraphics(e.Label, e.Style); err != nil {
		return err
	}
	_ = w.writeString("\n\t</edge>")
	return w.err
}

func (w *writer) writeEdgeGraphics(label string, s *EdgeStyle) error {
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

func (w *writer) Close() error {
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
`

const footer = `
  <data key="d7">
    <y:Resources/>
  </data>
</graphml>
`
