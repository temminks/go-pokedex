package catching

import (
	"math"
	"testing"
)

func TestGetProbability(t *testing.T) {
	probability635 := GetProbability(635)
	tolerance := 1e-12
	if math.Abs(probability635-0.15) > tolerance {
		t.Errorf("probability for 635 should be 0.15 but was %.18f", probability635)
	}
	probability256 := GetProbability(256)
	if math.Abs(probability256-0.2) > tolerance {
		t.Errorf("probability for 256 should be 0.20 but was %.18f", probability256)
	}
	probability0 := GetProbability(0)
	if math.Abs(probability0-0.7) > tolerance {
		t.Errorf("probability for 0 should be 0.70 but was %.18f", probability0)
	}
	probability255 := GetProbability(255)
	if math.Abs(probability255-0.2) > tolerance {
		t.Errorf("probability for 255 should be 0.2 but was %.18f", probability255)
	}
}
