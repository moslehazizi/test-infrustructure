package response

type Paginated[T any] struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Data    T   `json:"data"`
}
