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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
