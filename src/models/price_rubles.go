package models

import (
	"encoding/json"
)

// PriceRubles is stored in whole rubles. Fractional input values are truncated toward zero.
//
// @swaggertype number
// @minimum 0
type PriceRubles int

func (p *PriceRubles) UnmarshalJSON(data []byte) error {
	var n float64
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*p = PriceRubles(int(n))
	return nil
}
