package test

import (
	"reflect"
	"shopping-cart/server/domain"
	"testing"
)

var euro2decimals = &domain.Currency{
	Symbol:   "€",
	Decimals: 2,
}

var euro1decimal = &domain.Currency{
	Symbol:   "€",
	Decimals: 2,
}

func TestProduct_GetDiscountPriceIfAny(t *testing.T) {
	type fields struct {
		Code     string
		Name     string
		Price    *domain.Money
		Discount domain.DiscountType
	}
	type args struct {
		units int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "2x1 when odd",
			fields: fields{
				Code: "ACODE",
				Name: "2x1 when odd",
				Price: &domain.Money{
					Amount:   4,
					Currency: nil,
				},
				Discount: domain.BuyTwoGetOne,
			},
			args: args{3},
			want: int64(8),
		},
		{
			name: "2x1 when even",
			fields: fields{
				Code: "ACODE",
				Name: "2x1 when even",
				Price: &domain.Money{
					Amount:   4,
					Currency: nil,
				},
				Discount: domain.BuyTwoGetOne,
			},
			args: args{2},
			want: int64(4),
		},
		{
			name: "25% discount when more than 2",
			fields: fields{
				Code: "ACODE",
				Name: "25% discount when more than 2",
				Price: &domain.Money{
					Amount:   4,
					Currency: euro1decimal,
				},
				Discount: domain.Reduce25Percent,
			}, args: args{3}, want: int64(9)},
		{
			name: "No 25% discount when less than 3",
			fields: fields{
				Code: "ACODE",
				Name: "No 25% discount when less than 3",
				Price: &domain.Money{
					Amount:   4,
					Currency: euro1decimal,
				},
				Discount: domain.Reduce25Percent,
			},
			args: args{2},
			want: int64(8),
		},
		{
			name: "No 25% discount when more than 2 near zero  with 2 decimal",
			fields: fields{
				Code: "ACODE",
				Name: "No 25% discount when more than 2 near zero  with 2 decimal",
				Price: &domain.Money{
					Amount:   1,
					Currency: euro2decimals,
				},
				Discount: domain.Reduce25Percent,
			},
			args: args{3},
			want: int64(2),
		},
		{
			name: "No 25% discount when less than 3 near zero with 1 decimal",
			fields: fields{
				Code: "ACODE",
				Name: "No 25% discount when less than 3 near zero  with 1 decimal",
				Price: &domain.Money{
					Amount:   1,
					Currency: euro1decimal,
				},
				Discount: domain.Reduce25Percent,
			},
			args: args{2},
			want: int64(2),
		},
		{
			name: "25% discount when price near 0",
			fields: fields{
				Code: "ACODE",
				Name: "25% discount when price near 0",
				Price: &domain.Money{
					Amount:   1,
					Currency: euro1decimal,
				},
				Discount: domain.Reduce25Percent,
			},
			args: args{3},
			want: int64(2),
		},
		{
			name: "2x1 when price near 0",
			fields: fields{
				Code: "ACODE",
				Name: "2x1 when price near 0",
				Price: &domain.Money{
					Amount:   1,
					Currency: nil,
				},
				Discount: domain.BuyTwoGetOne,
			},
			args: args{2},
			want: int64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &domain.Product{
				Code:     tt.fields.Code,
				Name:     tt.fields.Name,
				Price:    tt.fields.Price,
				Discount: tt.fields.Discount,
			}
			if got := p.GetDiscountPriceIfAny(tt.args.units); got != tt.want {
				t.Errorf("GetDiscountPriceIfAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProducts(t *testing.T) {
	tests := []struct {
		name string
		want map[string]*domain.Product
	}{
		{name: "Check product \"database\"", want: map[string]*domain.Product{
			"PEN": {
				Code:     "PEN",
				Name:     "Lana Pen",
				Price:    &domain.Money{Amount: 500, Currency: euro2decimals},
				Discount: domain.BuyTwoGetOne,
			},
			"TSHIRT": {
				Code:     "TSHIRT",
				Name:     "Lana T-Shirt",
				Price:    &domain.Money{Amount: 2000, Currency: euro2decimals},
				Discount: domain.Reduce25Percent,
			},
			"MUG": {
				Code:     "MUG",
				Name:     "Lana Coffee Mug",
				Price:    &domain.Money{Amount: 750, Currency: euro2decimals},
				Discount: domain.None,
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := domain.GetProducts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProducts() = %v, want %v", got, tt.want)
			}
		})
	}
}
