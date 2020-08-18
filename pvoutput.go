package pvoutput

// PVEncodable is an interface that API objects need to implement
type PVEncodable interface {
	Encode() (string, error)
}
