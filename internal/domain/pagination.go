package domain

type Pagination struct {
	Page         uint64
	CountPerPage uint64
}

type Page struct {
	Pages uint64
	Total uint64
}
