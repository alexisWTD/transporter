package transporter

import (
	"testing"
)

var (
	fakesourceCN = NewNode("source1", "source", map[string]interface{}{"value": "rockettes"})
	fileNode     = NewNode("localfile", "file", map[string]interface{}{"uri": "file:///tmp/crap"})
)

var (
	testEmptyApiConfig = Api{
		Uri:             "",
		MetricsInterval: 100,
	}
)

func TestPipelineString(t *testing.T) {
	nodeRegistry["source"] = NewTestSourceImpl

	data := []struct {
		in           *Node
		terminalNode *Node
		out          string
	}{
		{
			&fakesourceCN,
			nil,
			" - Pipeline\n  - Source:                      source          no namespace set               no uri set\n",
		},
		{
			&fakesourceCN,
			&fileNode,
			" - Pipeline\n  - Source:                      source          no namespace set               no uri set\n  - Sink:   localfile            file            no namespace set               file:///tmp/crap\n",
		},
	}

	for _, v := range data {
		if v.terminalNode != nil {
			v.in.Attach(v.terminalNode)
		}
		p, err := NewPipeline(v.in, testEmptyApiConfig)
		if err != nil {
			t.Errorf("can't create pipeline, got %s", err.Error())
			t.FailNow()
		}
		if p.String() != v.out {
			t.Errorf("\nexp %s\ngot %s", v.out, p.String())
		}
	}
}
