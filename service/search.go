package service

// Search is an interface of a search object.
type Search interface {
	MarshalProto(permission interface{}) error
}
