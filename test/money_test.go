package test

import (
	"reflect"
	"shopping-cart/server/domain"
	"testing"
)

func TestMoney_String(t *testing.T) {
	euro := &domain.Currency{
		Symbol:   "$",
		Decimals: 2,
	}

	type fields struct {
		Amount   int64
		Currency *domain.Currency
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Zero", fields: fields{0, euro}, want: "0.00$"},
		{name: "Only decimals", fields: fields{22, euro}, want: "0.22$"},
		{name: "Integer part and no decimals", fields: fields{4575800, euro}, want: "45758.00$"},
		{name: "Integer part and decimals", fields: fields{542213, euro}, want: "5422.13$"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &domain.Money{
				Amount:   tt.fields.Amount,
				Currency: tt.fields.Currency,
			}
			if got := m.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
