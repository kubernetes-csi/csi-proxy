package internal

type FooRequest struct {
	SubMessage *FooRequestSubMessage
	Blob       []byte
}

type FooRequestSubMessage struct {
	Blah          int32
	SubSubMessage []*FooRequestSubSubMessage
}

type FooRequestSubSubMessage struct {
	Bools []bool
}

type FooResponse struct {
	Response32 int32
	SubMessage *FooRequestSubMessage
}
