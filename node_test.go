package osm

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestNode(t *testing.T) {
	data := []byte(`<node id="123" changeset="456" timestamp="2014-04-10T00:43:05Z" version="1" visible="true" user="user" uid="1357" lat="50.7107023" lon="6.0043943"/>`)

	n := Node{}
	err := xml.Unmarshal(data, &n)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if v := n.ID; v != 123 {
		t.Errorf("incorrect id, got %v", v)
	}

	if v := n.Version; v != 1 {
		t.Errorf("incorrect version, got %v", v)
	}

	if v := n.Lat; v != 50.7107023 {
		t.Errorf("incorrect lat, got %v", v)
	}

	if v := n.Lon; v != 6.0043943 {
		t.Errorf("incorrect lon, got %v", v)
	}
}

func TestNode_MarshalJSON(t *testing.T) {
	n := Node{
		ID: 123,
	}

	data, err := json.Marshal(n)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if !bytes.Equal(data, []byte(`{"type":"node","id":123,"lat":0,"lon":0,"visible":false,"timestamp":"0001-01-01T00:00:00Z"}`)) {
		t.Errorf("incorrect json: %v", string(data))
	}
}

func TestNode_MarshalXML(t *testing.T) {
	n := Node{
		ID: 123,
	}

	data, err := xml.Marshal(n)
	if err != nil {
		t.Fatalf("xml marshal error: %v", err)
	}

	expected := `<node id="123" lat="0" lon="0" user="" uid="0" visible="false" version="0" changeset="0" timestamp="0001-01-01T00:00:00Z"></node>`
	if !bytes.Equal(data, []byte(expected)) {
		t.Errorf("incorrect marshal, got: %s", string(data))
	}
}

func TestUnmarshalNodes(t *testing.T) {
	ns := Nodes{
		{ID: 123},
		{ID: 321},
	}

	data, err := ns.Marshal()
	if err != nil {
		t.Fatalf("nodes marshal error: %v", err)
	}

	ns2, err := UnmarshalNodes(data)
	if err != nil {
		t.Fatalf("nodes unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(ns, ns2) {
		t.Errorf("nodes not equal")
		t.Logf("%+v", ns)
		t.Logf("%+v", ns2)
	}

	// empty nodes
	ns = Nodes{}

	data, err = ns.Marshal()
	if err != nil {
		t.Fatalf("nodes marshal error: %v", err)
	}

	if l := len(data); l != 0 {
		t.Errorf("length of node data should be 0, got %v", l)
	}

	ns2, err = UnmarshalNodes(data)
	if err != nil {
		t.Fatalf("nodes unmarshal error: %v", err)
	}

	if ns2 != nil {
		t.Errorf("should return nil Nodes for empty data, got %v", ns2)
	}
}

func TestNodes_ids(t *testing.T) {
	ns := Nodes{
		{ID: 1, Version: 3},
		{ID: 2, Version: 4},
	}

	eids := ElementIDs{NodeID(1).ElementID(3), NodeID(2).ElementID(4)}
	if ids := ns.ElementIDs(); !reflect.DeepEqual(ids, eids) {
		t.Errorf("incorrect element ids: %v", ids)
	}

	fids := FeatureIDs{NodeID(1).FeatureID(), NodeID(2).FeatureID()}
	if ids := ns.FeatureIDs(); !reflect.DeepEqual(ids, fids) {
		t.Errorf("incorrect feature ids: %v", ids)
	}

	nids := []NodeID{1, 2}
	if ids := ns.IDs(); !reflect.DeepEqual(ids, nids) {
		t.Errorf("incorrect node ids: %v", nids)
	}
}

func TestNodes_SortByIDVersion(t *testing.T) {
	ns := Nodes{
		{ID: 7, Version: 3},
		{ID: 2, Version: 4},
		{ID: 5, Version: 2},
		{ID: 5, Version: 3},
		{ID: 5, Version: 4},
		{ID: 3, Version: 4},
		{ID: 4, Version: 4},
		{ID: 9, Version: 4},
	}

	ns.SortByIDVersion()

	eids := ElementIDs{
		NodeID(2).ElementID(4),
		NodeID(3).ElementID(4),
		NodeID(4).ElementID(4),
		NodeID(5).ElementID(2),
		NodeID(5).ElementID(3),
		NodeID(5).ElementID(4),
		NodeID(7).ElementID(3),
		NodeID(9).ElementID(4),
	}

	if ids := ns.ElementIDs(); !reflect.DeepEqual(ids, eids) {
		t.Errorf("incorrect sort: %v", eids)
	}
}
