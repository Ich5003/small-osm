osm [![CI](https://github.com/ich5003/small-osm/workflows/CI/badge.svg)](https://github.com/ich5003/small-osm/actions?query=workflow%3ACI+event%3Apush) [![Go Report Card](https://goreportcard.com/badge/github.com/ich5003/small-osm)](https://goreportcard.com/report/github.com/ich5003/small-osm) [![Go Reference](https://pkg.go.dev/badge/github.com/ich5003/small-osm.svg)](https://pkg.go.dev/github.com/ich5003/small-osm)
=====

This package is a general purpose library for reading, writing and working
with [OpenStreetMap](https://osm.org) data in Go (golang). It has the ability to
read [OSM XML](https://wiki.openstreetmap.org/wiki/OSM_XML) and
[PBF](https://wiki.openstreetmap.org/wiki/PBF_Format) data formats available at
[planet.osm.org](https://planet.osm.org/) or via the
[v0.6 API](https://wiki.openstreetmap.org/wiki/API_v0.6).


Made available by the package are the following types:

* Node
* Way
* Relation
* Changeset
* Note
* User

And the following “container” types:

* OSM - container returned via API
* Change - used by the replication API
* Diff - corresponds to [Overpass Augmented diffs](https://wiki.openstreetmap.org/wiki/Overpass_API/Augmented_Diffs)

## List of sub-package utilities

* [`annotate`](annotate) - adds lon/lat, version, changeset and orientation data to way and relation members
* [`osmapi`](osmapi) - supports all the v0.6 read/data endpoints
* [`osmgeojson`](osmgeojson) - OSM to GeoJSON conversion compatible with [osmtogeojson](https://github.com/tyrasd/osmtogeojson)
* [`osmpbf`](osmpbf) - stream processing of `*.osm.pbf` files
* [`osmxml`](osmxml) - stream processing of `*.osm` xml files
* [`replication`](replication) - fetch replication state and change files

## Concepts

This package refers to the core OSM data types as **Objects**. The Node, Way,
Relation, Changeset, Note and User types implement the `osm.Object` interface
and can be referenced using the `osm.ObjectID` type. As a result it is possible
to have a slice of `[]osm.Object` that contains nodes, changesets and users.

Individual versions of the core OSM Map Data types are referred to as **Elements**
and the set of versions for a give Node, Way or Relation is referred to as a
**Feature**. For example, an `osm.ElementID` could refer to "Node with id 10 and
version 3" and the `osm.FeatureID` would refer to "all versions of node with id 10."
Put another way, features represent a road and how it's changed over time and an
element is a specific version of that feature.

A number of helper methods are provided for dealing with features and elements.
The idea is to make it easy to work with a Way and its member nodes, for example.

## Scanning large data files

For small data it is possible to use the `encoding/xml` package in the
Go standard library to marshal/unmarshal the data. This is typically done using the
`osm.OSM` or `osm.Change` "container" structs.

For large data the package defines the `Scanner` interface implemented in both the [osmxml](osmxml)
and [osmpbf](osmpbf) sub-packages.

```go
type osm.Scanner interface {
	Scan() bool
	Object() osm.Object
	Err() error
	Close() error
}
```

This interface is designed to mimic the [bufio.Scanner](https://golang.org/pkg/bufio/#Scanner)
interface found in the Go standard library.

Example usage:

```go
f, err := os.Open("./delaware-latest.osm.pbf")
if err != nil {
	panic(err)
}
defer f.Close()

scanner := osmpbf.New(context.Background(), f, 3)
defer scanner.Close()

for scanner.Scan() {
	o := scanner.Object()
	// do something
}

scanErr := scanner.Err()
if scanErr != nil {
	panic(scanErr)
}
```

**Note:** Scanners are **not** safe for parallel use. One should feed the
objects into a channel and have workers read from that.
