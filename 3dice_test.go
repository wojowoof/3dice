package main

import (
	"testing"

	"wojones.com/src/dicegame"
)

func Test_rollcheck(t *testing.T) {
	type args struct {
		dg   *dicegame.DiceGame
		argv []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rollcheck(tt.args.dg, tt.args.argv)
			if (err != nil) != tt.wantErr {
				t.Errorf("rollcheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("rollcheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
