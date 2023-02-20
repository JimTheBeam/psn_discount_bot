package parser

import "testing"

func Test_parsePriceCurrency(t *testing.T) {
	type args struct {
		priceStr string
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 string
	}{
		{
			name: "1",
			args: args{
				priceStr: "1.700,00 TL",
			},
			want:  1700,
			want1: "TL",
		},
		{
			name: "2",
			args: args{
				priceStr: "139,80 TL",
			},
			want:  139.80,
			want1: "TL",
		},
		{
			name: "3",
			args: args{
				priceStr: "56,87 TL",
			},
			want:  56.87,
			want1: "TL",
		},
		{
			name: "4",
			args: args{
				priceStr: "700,00 TL",
			},
			want:  700.00,
			want1: "TL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parsePriceCurrency(tt.args.priceStr)
			if got != tt.want {
				t.Errorf("parsePriceCurrency() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parsePriceCurrency() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
