package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSuitecrmUpInitialValue(t *testing.T) {
	initialiseMetrics()

	if v := testutil.ToFloat64(suitecrmUp); v != 1 {
		t.Errorf("waiting suitecrm_up = 1, but got %v", v)
	}
}
