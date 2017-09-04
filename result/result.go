package result

type Item struct {
	Message		int
	Roundtrip	int64
}

type Result struct {
	Items		[]Item
}
