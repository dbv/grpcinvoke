package grpcinvoke

var (
	EmptyInnerMessage = InnerMessage{}
)

type InnerMessage struct {
}

func (me *InnerMessage) String() string {
	return ""
}

func (me *InnerMessage) IsEmpty() bool {
	return true
}
