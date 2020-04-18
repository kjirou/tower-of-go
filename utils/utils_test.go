package utils

import "testing"

func TestMatrixPosition_GetX(t *testing.T) {
	type fields struct {
		X int
		Y int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "returns X value",
			fields: fields{X: 2, Y: 1},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matrixPosition := &MatrixPosition{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if got := matrixPosition.GetX(); got != tt.want {
				t.Errorf("MatrixPosition.GetX() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrixPosition_GetY(t *testing.T) {
	type fields struct {
		X int
		Y int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "returns Y value",
			fields: fields{X: 1, Y: 2},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matrixPosition := &MatrixPosition{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if got := matrixPosition.GetY(); got != tt.want {
				t.Errorf("MatrixPosition.GetY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrixPosition_Validate(t *testing.T) {
	type fields struct {
		X int
		Y int
	}
	type args struct {
		rowLength    int
		columnLength int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "returns false when Y is less than 0",
			fields: fields{X: 0, Y: -1},
			args: args{rowLength: 1, columnLength: 1},
			want: false,
		},
		{
			name: "returns false when X is less than 0",
			fields: fields{X: -1, Y: 0},
			args: args{rowLength: 1, columnLength: 1},
			want: false,
		},
		{
			name: "returns false when Y is greater than or equal to rowLength",
			fields: fields{X: 0, Y: 1},
			args: args{rowLength: 1, columnLength: 1},
			want: false,
		},
		{
			name: "returns false when X is greater than or equal to columnLength",
			fields: fields{X: 1, Y: 0},
			args: args{rowLength: 1, columnLength: 1},
			want: false,
		},
		{
			name: "returns true if Y is greater than 0 and less than rowLength and X is greater than 0 and less than columnLength",
			fields: fields{X: 0, Y: 0},
			args: args{rowLength: 1, columnLength: 1},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matrixPosition := &MatrixPosition{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if got := matrixPosition.Validate(tt.args.rowLength, tt.args.columnLength); got != tt.want {
				t.Errorf("MatrixPosition.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
