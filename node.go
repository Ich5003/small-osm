package osm

import (
	"encoding/xml"
	"sort"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/paulmach/go.osm/internal/osmpb"
)

// NodeID corresponds the primary key of a node.
// The node id + version uniquely identify a node.
type NodeID int64

// Node is an osm point and allows for marshalling to/from osm xml.
type Node struct {
	XMLName     xml.Name    `xml:"node"`
	ID          NodeID      `xml:"id,attr"`
	Lat         float64     `xml:"lat,attr"`
	Lon         float64     `xml:"lon,attr"`
	User        string      `xml:"user,attr"`
	UserID      UserID      `xml:"uid,attr"`
	Visible     bool        `xml:"visible,attr"`
	Version     int         `xml:"version,attr"`
	ChangesetID ChangesetID `xml:"changeset,attr"`
	Timestamp   time.Time   `xml:"timestamp,attr"`
	Tags        Tags        `xml:"tag"`

	// Committed, is the estimated time this object was committed
	// and made visible in the central OSM database.
	Committed *time.Time `xml:"commited,attr,omitempty"`
}

// Nodes is a set of nodes with helper functions on top.
type Nodes []*Node

// Marshal encodes the nodes using protocol buffers.
func (ns Nodes) Marshal() ([]byte, error) {
	if len(ns) == 0 {
		return nil, nil
	}

	ss := &stringSet{}
	encoded := marshalNodes(ns, ss, true)
	encoded.Strings = ss.Strings()

	return proto.Marshal(encoded)
}

// UnmarshalNodes will unmarshal the data into a list of nodes.
func UnmarshalNodes(data []byte) (Nodes, error) {
	if len(data) == 0 {
		return nil, nil
	}

	pbf := &osmpb.DenseNodes{}
	err := proto.Unmarshal(data, pbf)
	if err != nil {
		return nil, err
	}

	return unmarshalNodes(pbf, pbf.GetStrings(), nil)
}

type nodesSort Nodes

// SortByIDVersion will sort the set of nodes first by id and then version
// in ascending order.
func (ns Nodes) SortByIDVersion() {
	sort.Sort(nodesSort(ns))
}

func (ns nodesSort) Len() int      { return len(ns) }
func (ns nodesSort) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }
func (ns nodesSort) Less(i, j int) bool {
	if ns[i].ID == ns[j].ID {
		return ns[i].Version < ns[j].Version
	}

	return ns[i].ID < ns[j].ID
}
