package api

// Used for autocomplete requests for simple_traits, roles and organizations
type AutocompleteRequest struct {
	Prefix string `json:"prefix" binding:"required"`
	Size   int    `json:"size" binding:"required"`
}
