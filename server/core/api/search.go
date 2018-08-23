package api

type SimpleTraitAutocompleteRequest struct {
	Prefix string `json:"prefix" binding:"required"`
	Size   int    `json:"size" binding:"required"`
}
