package types

type Header struct {
	cols []string
}

type Row struct {
	cols []string
}

type Table struct {
	header Header
	rows   []Row
}
