package fakes

import (
	"math"
	"time"

	"github.com/primandproper/platform/identifiers"
	"github.com/primandproper/platform/numbers"

	fake "github.com/brianvoe/gofakeit/v7"
)

func init() {
	if err := fake.Seed(time.Now().UnixNano()); err != nil {
		panic(err)
	}
}

const (
	exampleQuantity = 3
)

// BuildFakeID builds a fake ID.
func BuildFakeID() string {
	return identifiers.New()
}

// BuildFakeTime builds a fake time.
func BuildFakeTime() time.Time {
	return fake.Date().Add(0).Truncate(time.Second).UTC()
}

func buildFakeNumber() float64 {
	return math.Round(float64((fake.Number(101, math.MaxInt8-1) * 100) / 100))
}

// buildUniqueString builds a fake string.
func buildUniqueString() string {
	return fake.LoremIpsumSentence(7)
}

// buildFakePassword builds a fake password.
func buildFakePassword() string {
	return fake.Password(true, true, true, true, false, 32)
}

func BuildFakeFloat32RangeWithOptionalMax() numbers.MinRange[float32] {
	minimum := float32(buildFakeNumber())
	return numbers.MinRange[float32]{
		Min: minimum,
		Max: new(float32(buildFakeNumber()) + minimum),
	}
}

func BuildFakeOptionalFloat32Range() numbers.OpenRange[float32] {
	minimum := float32(buildFakeNumber())
	return numbers.OpenRange[float32]{
		Min: new(minimum),
		Max: new(float32(buildFakeNumber()) + minimum),
	}
}

func BuildFakeOptionalUint32Range() numbers.OpenRange[uint32] {
	minimum := uint32(buildFakeNumber())
	return numbers.OpenRange[uint32]{
		Min: new(minimum),
		Max: new(uint32(buildFakeNumber()) + minimum),
	}
}

func BuildFakeUint16RangeWithOptionalMax() numbers.MinRange[uint16] {
	minimum := uint16(buildFakeNumber())
	return numbers.MinRange[uint16]{
		Min: minimum,
		Max: new(uint16(buildFakeNumber()) + minimum),
	}
}

func BuildFakeUint32RangeWithOptionalMax() numbers.MinRange[uint32] {
	minimum := uint32(buildFakeNumber())
	return numbers.MinRange[uint32]{
		Min: minimum,
		Max: new(uint32(buildFakeNumber()) + minimum),
	}
}

func BuildFakeUint32RangeWithOptionalMaxUpdateRequestInput() numbers.OpenRangeUpdateRequestInput[uint32] {
	minimum := uint32(buildFakeNumber())
	return numbers.OpenRangeUpdateRequestInput[uint32]{
		Min: &minimum,
		Max: new(uint32(buildFakeNumber()) + minimum),
	}
}
