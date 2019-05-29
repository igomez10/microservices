package main

import "testing"

func TestGetBitcoinPrice(t *testing.T) {

	_, err := getBitcoinPrice()
	if err != nil {
		t.Errorf("error getting bitcoint price %s", err)
	}

}
