package frontendsrc

import (
	"testing"

	"gitlab.com/verygoodsoftwarenotvirus/naff/models/testprojects"

	"github.com/stretchr/testify/assert"
)

func Test_buildSomethingFrontendModels(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		proj := testprojects.BuildTodoApp()
		typ := proj.LastDataType()

		expected := `import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class Item {
  id: number;
  name: string;
  details: string;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    name = "";
    details = "";
    this.createdOn = 0;
  }

static areEqual = function(
  i1: Item,
  i2: Item,
): boolean {
    return (
      i1.id === i2.id &&
      i1.name === i2.name &&
      i1.details === i2.details &&
      i1.archivedOn === i2.archivedOn
    );
  }
}

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<ValidIngredient> ({
  name: Factory.Sync.each(() =>  faker.random.word()),
  details: Factory.Sync.each(() =>  faker.random.word()),
  ...defaultFactories,
});
`
		actual := buildSomethingFrontendModels(typ)()

		assert.Equal(t, expected, actual, "expected and actual do not match")
	})

	T.Run("gamut", func(t *testing.T) {
		proj := testprojects.BuildEveryTypeApp()
		typ := proj.LastDataType()

		expected := `import * as Factory from "factory.ts";
import faker from "faker";
import {defaultFactories} from "@/models/fakes";

export class EveryType {
  id: number;
  string: string;
  pointerToString?: string;
  bool: boolean;
  pointerToBool?: boolean;
  int: number;
  pointerToInt?: number;
  int8: number;
  pointerToInt8?: number;
  int16: number;
  pointerToInt16?: number;
  int32: number;
  pointerToInt32?: number;
  int64: number;
  pointerToInt64?: number;
  uint: number;
  pointerToUint?: number;
  uint8: number;
  pointerToUint8?: number;
  uint16: number;
  pointerToUint16?: number;
  uint32: number;
  pointerToUint32?: number;
  uint64: number;
  pointerToUint64?: number;
  float32: number;
  pointerToFloat32?: number;
  float64: number;
  pointerToFloat64?: number;
  createdOn: number;
  updatedOn?: number;
  archivedOn?: number;

  constructor() {
    this.id = 0;
    string = "";
    pointerToString = "";
    bool = false;
    pointerToBool = false;
    int = 0;
    pointerToInt = 0;
    int8 = 0;
    pointerToInt8 = 0;
    int16 = 0;
    pointerToInt16 = 0;
    int32 = 0;
    pointerToInt32 = 0;
    int64 = 0;
    pointerToInt64 = 0;
    uint = 0;
    pointerToUint = 0;
    uint8 = 0;
    pointerToUint8 = 0;
    uint16 = 0;
    pointerToUint16 = 0;
    uint32 = 0;
    pointerToUint32 = 0;
    uint64 = 0;
    pointerToUint64 = 0;
    float32 = 0;
    pointerToFloat32 = 0;
    float64 = 0;
    pointerToFloat64 = 0;
    this.createdOn = 0;
  }

static areEqual = function(
  et1: EveryType,
  et2: EveryType,
): boolean {
    return (
      et1.id === et2.id &&
      et1.string === et2.string &&
      et1.pointerToString === et2.pointerToString &&
      et1.bool === et2.bool &&
      et1.pointerToBool === et2.pointerToBool &&
      et1.int === et2.int &&
      et1.pointerToInt === et2.pointerToInt &&
      et1.int8 === et2.int8 &&
      et1.pointerToInt8 === et2.pointerToInt8 &&
      et1.int16 === et2.int16 &&
      et1.pointerToInt16 === et2.pointerToInt16 &&
      et1.int32 === et2.int32 &&
      et1.pointerToInt32 === et2.pointerToInt32 &&
      et1.int64 === et2.int64 &&
      et1.pointerToInt64 === et2.pointerToInt64 &&
      et1.uint === et2.uint &&
      et1.pointerToUint === et2.pointerToUint &&
      et1.uint8 === et2.uint8 &&
      et1.pointerToUint8 === et2.pointerToUint8 &&
      et1.uint16 === et2.uint16 &&
      et1.pointerToUint16 === et2.pointerToUint16 &&
      et1.uint32 === et2.uint32 &&
      et1.pointerToUint32 === et2.pointerToUint32 &&
      et1.uint64 === et2.uint64 &&
      et1.pointerToUint64 === et2.pointerToUint64 &&
      et1.float32 === et2.float32 &&
      et1.pointerToFloat32 === et2.pointerToFloat32 &&
      et1.float64 === et2.float64 &&
      et1.pointerToFloat64 === et2.pointerToFloat64 &&
      et1.archivedOn === et2.archivedOn
    );
  }
}

export const fakeValidIngredientFactory = Factory.Sync.makeFactory<ValidIngredient> ({
  string: Factory.Sync.each(() =>  faker.random.word()),
  pointerToString: Factory.Sync.each(() =>  faker.random.word()),
  bool: Factory.Sync.each(() =>  faker.random.boolean()),
  pointerToBool: Factory.Sync.each(() =>  faker.random.boolean()),
  int: Factory.Sync.each(() =>  faker.random.number()),
  pointerToInt: Factory.Sync.each(() =>  faker.random.number()),
  int8: Factory.Sync.each(() =>  faker.random.number()),
  pointerToInt8: Factory.Sync.each(() =>  faker.random.number()),
  int16: Factory.Sync.each(() =>  faker.random.number()),
  pointerToInt16: Factory.Sync.each(() =>  faker.random.number()),
  int32: Factory.Sync.each(() =>  faker.random.number()),
  pointerToInt32: Factory.Sync.each(() =>  faker.random.number()),
  int64: Factory.Sync.each(() =>  faker.random.number()),
  pointerToInt64: Factory.Sync.each(() =>  faker.random.number()),
  uint: Factory.Sync.each(() =>  faker.random.number()),
  pointerToUint: Factory.Sync.each(() =>  faker.random.number()),
  uint8: Factory.Sync.each(() =>  faker.random.number()),
  pointerToUint8: Factory.Sync.each(() =>  faker.random.number()),
  uint16: Factory.Sync.each(() =>  faker.random.number()),
  pointerToUint16: Factory.Sync.each(() =>  faker.random.number()),
  uint32: Factory.Sync.each(() =>  faker.random.number()),
  pointerToUint32: Factory.Sync.each(() =>  faker.random.number()),
  uint64: Factory.Sync.each(() =>  faker.random.number()),
  pointerToUint64: Factory.Sync.each(() =>  faker.random.number()),
  float32: Factory.Sync.each(() =>  faker.random.number()),
  pointerToFloat32: Factory.Sync.each(() =>  faker.random.number()),
  float64: Factory.Sync.each(() =>  faker.random.number()),
  pointerToFloat64: Factory.Sync.each(() =>  faker.random.number()),
  ...defaultFactories,
});
`
		actual := buildSomethingFrontendModels(typ)()

		assert.Equal(t, expected, actual, "expected and actual do not match")
	})
}

var allTypes = []string{
	"string",
	"*string",
	"bool",
	"*bool",
	"int",
	"*int",
	"int8",
	"*int8",
	"int16",
	"*int16",
	"int32",
	"*int32",
	"int64",
	"*int64",
	"uint",
	"*uint",
	"uint8",
	"*uint8",
	"uint16",
	"*uint16",
	"uint32",
	"*uint32",
	"uint64",
	"*uint64",
	"float32",
	"*float32",
	"float64",
	"*float64",
}

func Test_typeToDefaultJSValue(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		for _, typ := range allTypes {
			assert.NotEmpty(t, typeToDefaultJSValue(typ), "should have an entry for %s", typ)
		}
	})
}

func Test_typeToFakerValue(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		for _, typ := range allTypes {
			assert.NotEmpty(t, typeToFakerValue(typ), "should have an entry for %s", typ)
		}
	})
}

func Test_typeToJSType(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		for _, typ := range allTypes {
			assert.NotEmpty(t, typeToJSType(typ), "should have an entry for %s", typ)
		}
	})
}
