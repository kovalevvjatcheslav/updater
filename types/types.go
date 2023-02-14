package types

type Row struct {
	Cols []string
}

type Table struct {
	Header Row
	Rows   []Row
}
