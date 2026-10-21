// Package region declares enumerations of text.
package region

// Region is a region of a cloud, whose text is its value.
//
//dokimi:enum
type Region string

// The regions of a cloud.
const (
	US Region = "us-east"
	EU Region = "eu-west"
	AP Region = "ap-south"
)

// Zone is a zone of a region, whose zero value is its last variant.
//
//dokimi:enum
type Zone string

// The zones of a region.
const (
	North Zone = "north"
	South Zone = "south"
	Unset Zone = ""
)

// Lane is a lane of a road. Its directive takes the text of each variant from
// the trailing comment of the variant, so a text is not the value.
//
//dokimi:enum text=comment
type Lane string

// The lanes of a road.
const (
	Left  Lane = "l" // left
	Right Lane = "r" // right
)
