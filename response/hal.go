package response

// Link represents a hypermedia link
type Link struct {
	Href string `json:"href"`
}

// HAL represents a hypermedia (HATEOAS) style response
type HAL[T any] struct {
	Links map[string]Link `json:"_links"`
	Data  T               `json:"data"`
}

// NewHAL creates a new HAL response with links and data
func NewHAL[T any](data T, links map[string]Link) HAL[T] {
	return HAL[T]{
		Links: links,
		Data:  data,
	}
}

// NewHALWithSelfLink creates a new HAL response with a self link
func NewHALWithSelfLink[T any](data T, selfHref string) HAL[T] {
	return HAL[T]{
		Links: map[string]Link{
			"self": {Href: selfHref},
		},
		Data: data,
	}
}
