package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:   "Milk",
		Price:  10,
		SKU: "abc-abc-abs",
	}
	err := p.ProductValidator()
	if err != nil {
		t.Fatal(err)
	}
}