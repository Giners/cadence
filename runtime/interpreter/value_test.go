/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package interpreter_test

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"testing"

	"github.com/onflow/atree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence/runtime/common"
	. "github.com/onflow/cadence/runtime/interpreter"
	"github.com/onflow/cadence/runtime/sema"
	checkerUtils "github.com/onflow/cadence/runtime/tests/checker"
	"github.com/onflow/cadence/runtime/tests/utils"
)

func newTestCompositeValue(storage Storage, owner common.Address) *CompositeValue {
	return NewCompositeValue(
		storage,
		utils.TestLocation,
		"Test",
		common.CompositeKindStructure,
		NewStringValueOrderedMap(),
		owner,
	)
}

func TestOwnerNewArray(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}

	value := newTestCompositeValue(storage, oldOwner)

	assert.Equal(t, oldOwner, value.GetOwner())

	array := NewArrayValue(
		VariableSizedStaticType{
			Type: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		value,
	)

	value = array.GetIndex(ReturnEmptyLocationRange, 0).(*CompositeValue)

	assert.Equal(t, common.Address{}, array.GetOwner())
	assert.Equal(t, common.Address{}, value.GetOwner())
}

func TestOwnerArrayDeepCopy(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	value := newTestCompositeValue(storage, oldOwner)

	array := NewArrayValue(
		VariableSizedStaticType{
			Type: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		value,
	)

	arrayCopy, err := array.DeepCopy(storage, atree.Address(newOwner))
	require.NoError(t, err)
	array = arrayCopy.(*ArrayValue)

	value = array.GetIndex(ReturnEmptyLocationRange, 0).(*CompositeValue)

	assert.Equal(t, newOwner, array.GetOwner())
	assert.Equal(t, newOwner, value.GetOwner())
}

func TestOwnerArrayElement(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	value := newTestCompositeValue(storage, oldOwner)

	array := NewArrayValueWithAddress(
		VariableSizedStaticType{
			Type: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		newOwner,
		value,
	)

	value = array.GetIndex(ReturnEmptyLocationRange, 0).(*CompositeValue)

	assert.Equal(t, newOwner, array.GetOwner())
	assert.Equal(t, newOwner, value.GetOwner())
}

func TestOwnerArraySetIndex(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	value1 := newTestCompositeValue(storage, oldOwner)
	value2 := newTestCompositeValue(storage, oldOwner)

	array := NewArrayValueWithAddress(
		VariableSizedStaticType{
			Type: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		newOwner,
		value1,
	)

	value1 = array.GetIndex(ReturnEmptyLocationRange, 0).(*CompositeValue)

	assert.Equal(t, newOwner, array.GetOwner())
	assert.Equal(t, newOwner, value1.GetOwner())
	assert.Equal(t, oldOwner, value2.GetOwner())

	array.SetIndex(ReturnEmptyLocationRange, 0, value2)

	value2 = array.GetIndex(ReturnEmptyLocationRange, 0).(*CompositeValue)

	assert.Equal(t, newOwner, array.GetOwner())
	assert.Equal(t, newOwner, value1.GetOwner())
	assert.Equal(t, newOwner, value2.GetOwner())
}

func TestOwnerArrayAppend(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	value := newTestCompositeValue(storage, oldOwner)

	array := NewArrayValueWithAddress(
		VariableSizedStaticType{
			Type: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		newOwner,
	)

	assert.Equal(t, newOwner, array.GetOwner())
	assert.Equal(t, oldOwner, value.GetOwner())

	array.Append(value)

	value = array.GetIndex(ReturnEmptyLocationRange, 0).(*CompositeValue)

	assert.Equal(t, newOwner, array.GetOwner())
	assert.Equal(t, newOwner, value.GetOwner())
}

func TestOwnerArrayInsert(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	value := newTestCompositeValue(storage, oldOwner)

	array := NewArrayValueWithAddress(
		VariableSizedStaticType{
			Type: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		newOwner,
	)

	assert.Equal(t, newOwner, array.GetOwner())
	assert.Equal(t, oldOwner, value.GetOwner())

	array.Insert(ReturnEmptyLocationRange, 0, value)

	value = array.GetIndex(ReturnEmptyLocationRange, 0).(*CompositeValue)

	assert.Equal(t, newOwner, array.GetOwner())
	assert.Equal(t, newOwner, value.GetOwner())
}

func TestOwnerArrayRemove(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	owner := common.Address{0x1}

	value := newTestCompositeValue(storage, owner)

	array := NewArrayValueWithAddress(
		VariableSizedStaticType{
			Type: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		owner,
		value,
	)

	assert.Equal(t, owner, array.GetOwner())
	assert.Equal(t, owner, value.GetOwner())

	value = array.Remove(ReturnEmptyLocationRange, 0).(*CompositeValue)

	assert.Equal(t, owner, array.GetOwner())
	assert.Equal(t, common.Address{}, value.GetOwner())
}

func TestOwnerNewDictionary(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}

	keyValue := NewStringValue("test")
	value := newTestCompositeValue(storage, oldOwner)

	assert.Equal(t, oldOwner, value.GetOwner())

	dictionary := NewDictionaryValue(
		DictionaryStaticType{
			KeyType:   PrimitiveStaticTypeString,
			ValueType: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		keyValue, value,
	)

	// NOTE: keyValue is string, has no owner

	queriedValue, _, _ := dictionary.GetKey(keyValue)
	value = queriedValue.(*CompositeValue)

	assert.Equal(t, common.Address{}, dictionary.GetOwner())
	assert.Equal(t, common.Address{}, value.GetOwner())
}

func TestOwnerDictionary(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	keyValue := NewStringValue("test")
	value := newTestCompositeValue(storage, oldOwner)

	dictionary := NewDictionaryValueWithAddress(
		DictionaryStaticType{
			KeyType:   PrimitiveStaticTypeString,
			ValueType: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		newOwner,
		keyValue, value,
	)

	// NOTE: keyValue is string, has no owner

	queriedValue, _, _ := dictionary.GetKey(keyValue)
	value = queriedValue.(*CompositeValue)

	assert.Equal(t, newOwner, dictionary.GetOwner())
	assert.Equal(t, newOwner, value.GetOwner())
}

func TestOwnerDictionaryCopy(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	keyValue := NewStringValue("test")
	value := newTestCompositeValue(storage, oldOwner)

	dictionary := NewDictionaryValueWithAddress(
		DictionaryStaticType{
			KeyType:   PrimitiveStaticTypeString,
			ValueType: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		newOwner,
		keyValue, value,
	)

	copyResult, err := dictionary.DeepCopy(storage, atree.Address{})
	require.NoError(t, err)

	dictionaryCopy := copyResult.(*DictionaryValue)

	queriedValue, _, _ := dictionaryCopy.GetKey(keyValue)
	value = queriedValue.(*CompositeValue)

	assert.Equal(t, common.Address{}, dictionaryCopy.GetOwner())
	assert.Equal(t, common.Address{}, value.GetOwner())
}

func TestOwnerDictionarySetSome(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	keyValue := NewStringValue("test")
	value := newTestCompositeValue(storage, oldOwner)

	dictionary := NewDictionaryValueWithAddress(
		DictionaryStaticType{
			KeyType:   PrimitiveStaticTypeString,
			ValueType: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		newOwner,
	)

	assert.Equal(t, newOwner, dictionary.GetOwner())
	assert.Equal(t, oldOwner, value.GetOwner())

	inter, err := NewInterpreter(
		nil,
		utils.TestLocation,
		WithStorage(storage),
	)
	require.NoError(t, err)

	dictionary.Set(
		inter,
		ReturnEmptyLocationRange,
		keyValue,
		NewSomeValueNonCopying(value),
	)

	queriedValue, _, _ := dictionary.GetKey(keyValue)
	value = queriedValue.(*CompositeValue)

	assert.Equal(t, newOwner, dictionary.GetOwner())
	assert.Equal(t, newOwner, value.GetOwner())
}

func TestOwnerDictionaryInsertNonExisting(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	keyValue := NewStringValue("test")
	value := newTestCompositeValue(storage, oldOwner)

	dictionary := NewDictionaryValueWithAddress(
		DictionaryStaticType{
			KeyType:   PrimitiveStaticTypeString,
			ValueType: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		newOwner,
	)

	assert.Equal(t, newOwner, dictionary.GetOwner())
	assert.Equal(t, oldOwner, value.GetOwner())

	inter, err := NewInterpreter(
		nil,
		utils.TestLocation,
		WithStorage(storage),
	)
	require.NoError(t, err)

	existingValue := dictionary.Insert(
		inter.Storage,
		ReturnEmptyLocationRange,
		keyValue,
		value,
	)
	assert.Equal(t, NilValue{}, existingValue)

	queriedValue, _, _ := dictionary.GetKey(keyValue)
	value = queriedValue.(*CompositeValue)

	assert.Equal(t, newOwner, dictionary.GetOwner())
	assert.Equal(t, newOwner, value.GetOwner())
}

func TestOwnerDictionaryRemove(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	keyValue := NewStringValue("test")
	value1 := newTestCompositeValue(storage, oldOwner)
	value2 := newTestCompositeValue(storage, oldOwner)

	dictionary := NewDictionaryValueWithAddress(
		DictionaryStaticType{
			KeyType:   PrimitiveStaticTypeString,
			ValueType: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		newOwner,
		keyValue, value1,
	)

	assert.Equal(t, newOwner, dictionary.GetOwner())
	assert.Equal(t, oldOwner, value1.GetOwner())
	assert.Equal(t, oldOwner, value2.GetOwner())

	inter, err := NewInterpreter(
		nil,
		utils.TestLocation,
		WithStorage(storage),
	)
	require.NoError(t, err)

	existingValue := dictionary.Insert(
		inter.Storage,
		ReturnEmptyLocationRange,
		keyValue,
		value2,
	)
	require.IsType(t, &SomeValue{}, existingValue)
	value1 = existingValue.(*SomeValue).Value.(*CompositeValue)

	queriedValue, _, _ := dictionary.GetKey(keyValue)
	value2 = queriedValue.(*CompositeValue)

	assert.Equal(t, newOwner, dictionary.GetOwner())
	assert.Equal(t, common.Address{}, value1.GetOwner())
	assert.Equal(t, newOwner, value2.GetOwner())
}

func TestOwnerDictionaryInsertExisting(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	keyValue := NewStringValue("test")
	value := newTestCompositeValue(storage, oldOwner)

	dictionary := NewDictionaryValueWithAddress(
		DictionaryStaticType{
			KeyType:   PrimitiveStaticTypeString,
			ValueType: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		newOwner,
		keyValue, value,
	)

	assert.Equal(t, newOwner, dictionary.GetOwner())
	assert.Equal(t, oldOwner, value.GetOwner())

	inter, err := NewInterpreter(
		nil,
		utils.TestLocation,
		WithStorage(storage),
	)
	require.NoError(t, err)

	existingValue := dictionary.Remove(
		inter.Storage,
		ReturnEmptyLocationRange,
		keyValue,
	)
	require.IsType(t, &SomeValue{}, existingValue)
	value = existingValue.(*SomeValue).Value.(*CompositeValue)

	assert.Equal(t, newOwner, dictionary.GetOwner())
	assert.Equal(t, common.Address{}, value.GetOwner())
}

func TestOwnerNewComposite(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}

	composite := newTestCompositeValue(storage, oldOwner)

	assert.Equal(t, oldOwner, composite.GetOwner())
}

func TestOwnerCompositeSet(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}
	newOwner := common.Address{0x2}

	value := newTestCompositeValue(storage, oldOwner)
	composite := newTestCompositeValue(storage, newOwner)

	assert.Equal(t, oldOwner, value.GetOwner())
	assert.Equal(t, newOwner, composite.GetOwner())

	inter, err := NewInterpreter(
		nil,
		utils.TestLocation,
		WithStorage(storage),
	)
	require.NoError(t, err)

	const fieldName = "test"

	composite.SetMember(inter, ReturnEmptyLocationRange, fieldName, value)

	value = composite.GetMember(inter, ReturnEmptyLocationRange, fieldName).(*CompositeValue)

	assert.Equal(t, newOwner, composite.GetOwner())
	assert.Equal(t, newOwner, value.GetOwner())
}

func TestOwnerCompositeCopy(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	oldOwner := common.Address{0x1}

	value := newTestCompositeValue(storage, oldOwner)
	composite := newTestCompositeValue(storage, oldOwner)

	inter, err := NewInterpreter(
		nil,
		utils.TestLocation,
		WithStorage(storage),
	)
	require.NoError(t, err)

	const fieldName = "test"

	composite.SetMember(
		inter,
		ReturnEmptyLocationRange,
		fieldName,
		value,
	)

	copyResult, err := composite.DeepCopy(storage, atree.Address{})
	require.NoError(t, err)
	composite = copyResult.(*CompositeValue)

	value = composite.GetMember(inter, ReturnEmptyLocationRange, fieldName).(*CompositeValue)

	assert.Equal(t, common.Address{}, composite.GetOwner())
	assert.Equal(t, common.Address{}, value.GetOwner())
}

func TestStringer(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	type testCase struct {
		value    Value
		expected string
	}

	stringerTests := map[string]testCase{
		"UInt": {
			value:    NewUIntValueFromUint64(10),
			expected: "10",
		},
		"UInt8": {
			value:    UInt8Value(8),
			expected: "8",
		},
		"UInt16": {
			value:    UInt16Value(16),
			expected: "16",
		},
		"UInt32": {
			value:    UInt32Value(32),
			expected: "32",
		},
		"UInt64": {
			value:    UInt64Value(64),
			expected: "64",
		},
		"UInt128": {
			value:    NewUInt128ValueFromUint64(128),
			expected: "128",
		},
		"UInt256": {
			value:    NewUInt256ValueFromUint64(256),
			expected: "256",
		},
		"Int8": {
			value:    Int8Value(-8),
			expected: "-8",
		},
		"Int16": {
			value:    Int16Value(-16),
			expected: "-16",
		},
		"Int32": {
			value:    Int32Value(-32),
			expected: "-32",
		},
		"Int64": {
			value:    Int64Value(-64),
			expected: "-64",
		},
		"Int128": {
			value:    NewInt128ValueFromInt64(-128),
			expected: "-128",
		},
		"Int256": {
			value:    NewInt256ValueFromInt64(-256),
			expected: "-256",
		},
		"Word8": {
			value:    Word8Value(8),
			expected: "8",
		},
		"Word16": {
			value:    Word16Value(16),
			expected: "16",
		},
		"Word32": {
			value:    Word32Value(32),
			expected: "32",
		},
		"Word64": {
			value:    Word64Value(64),
			expected: "64",
		},
		"UFix64": {
			value:    NewUFix64ValueWithInteger(64),
			expected: "64.00000000",
		},
		"Fix64": {
			value:    NewFix64ValueWithInteger(-32),
			expected: "-32.00000000",
		},
		"Void": {
			value:    VoidValue{},
			expected: "()",
		},
		"true": {
			value:    BoolValue(true),
			expected: "true",
		},
		"false": {
			value:    BoolValue(false),
			expected: "false",
		},
		"some": {
			value:    NewSomeValueNonCopying(BoolValue(true)),
			expected: "true",
		},
		"nil": {
			value:    NilValue{},
			expected: "nil",
		},
		"String": {
			value:    NewStringValue("Flow ridah!"),
			expected: "\"Flow ridah!\"",
		},
		"Array": {
			value: NewArrayValue(
				VariableSizedStaticType{
					Type: PrimitiveStaticTypeAnyStruct,
				},
				storage,
				NewIntValueFromInt64(10),
				NewStringValue("TEST"),
			),
			expected: "[10, \"TEST\"]",
		},
		"Address": {
			value:    NewAddressValue(common.Address{0, 0, 0, 0, 0, 0, 0, 1}),
			expected: "0x1",
		},
		"composite": {
			value: func() Value {
				members := NewStringValueOrderedMap()
				members.Set("y", NewStringValue("bar"))

				return NewCompositeValue(
					storage,
					utils.TestLocation,
					"Foo",
					common.CompositeKindResource,
					members,
					common.Address{},
				)
			}(),
			expected: "S.test.Foo(y: \"bar\")",
		},
		"composite with custom stringer": {
			value: func() Value {
				members := NewStringValueOrderedMap()
				members.Set("y", NewStringValue("bar"))

				compositeValue := NewCompositeValue(
					storage,
					utils.TestLocation,
					"Foo",
					common.CompositeKindResource,
					members,
					common.Address{},
				)

				compositeValue.Stringer = func(_ SeenReferences) string {
					return "y --> bar"
				}

				return compositeValue
			}(),
			expected: "y --> bar",
		},
		"Link": {
			value: LinkValue{
				TargetPath: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "foo",
				},
				Type: PrimitiveStaticTypeInt,
			},
			expected: "Link<Int>(/storage/foo)",
		},
		"Path": {
			value: PathValue{
				Domain:     common.PathDomainStorage,
				Identifier: "foo",
			},
			expected: "/storage/foo",
		},
		"Type": {
			value:    TypeValue{Type: PrimitiveStaticTypeInt},
			expected: "Type<Int>()",
		},
		"Capability with borrow type": {
			value: &CapabilityValue{
				Path: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "foo",
				},
				Address:    NewAddressValueFromBytes([]byte{1, 2, 3, 4, 5}),
				BorrowType: PrimitiveStaticTypeInt,
			},
			expected: "Capability<Int>(address: 0x102030405, path: /storage/foo)",
		},
		"Capability without borrow type": {
			value: &CapabilityValue{
				Path: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "foo",
				},
				Address: NewAddressValueFromBytes([]byte{1, 2, 3, 4, 5}),
			},
			expected: "Capability(address: 0x102030405, path: /storage/foo)",
		},
		"Dictionary": {
			value: NewDictionaryValue(
				DictionaryStaticType{
					KeyType:   PrimitiveStaticTypeString,
					ValueType: PrimitiveStaticTypeUInt8,
				},
				storage,
				NewStringValue("a"), UInt8Value(42),
				NewStringValue("b"), UInt8Value(99),
			),
			expected: `{"a": 42, "b": 99}`,
		},
		"Recursive ephemeral reference (array)": {
			value: func() Value {
				array := NewArrayValue(
					VariableSizedStaticType{
						Type: PrimitiveStaticTypeAnyStruct,
					},
					storage,
				)
				arrayRef := &EphemeralReferenceValue{Value: array}
				array.Insert(ReturnEmptyLocationRange, 0, arrayRef)
				return array
			}(),
			expected: `[[...]]`,
		},
	}

	test := func(name string, testCase testCase) {

		t.Run(name, func(t *testing.T) {

			t.Parallel()

			assert.Equal(t,
				testCase.expected,
				testCase.value.String(),
			)
		})
	}

	for name, testCase := range stringerTests {
		test(name, testCase)
	}
}

func TestVisitor(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	var intVisits, stringVisits int

	visitor := EmptyVisitor{
		IntValueVisitor: func(interpreter *Interpreter, value IntValue) {
			intVisits++
		},
		StringValueVisitor: func(interpreter *Interpreter, value *StringValue) {
			stringVisits++
		},
	}

	var value Value
	value = NewIntValueFromInt64(42)
	value = NewSomeValueNonCopying(value)
	value = NewArrayValue(
		VariableSizedStaticType{
			Type: PrimitiveStaticTypeAnyStruct,
		},
		storage,
		value,
	)
	value = NewDictionaryValue(
		DictionaryStaticType{
			KeyType:   PrimitiveStaticTypeString,
			ValueType: PrimitiveStaticTypeAny,
		},
		storage,
		NewStringValue("42"), value,
	)
	members := NewStringValueOrderedMap()
	members.Set("foo", value)
	value = NewCompositeValue(
		storage,
		utils.TestLocation,
		"Foo",
		common.CompositeKindStructure,
		members,
		common.Address{},
	)

	value.Accept(nil, visitor)

	require.Equal(t, 1, intVisits)
	require.Equal(t, 1, stringVisits)
}

func TestKeyString(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	type testCase struct {
		value    HasKeyString
		expected string
	}

	stringerTests := map[string]testCase{
		"UInt": {
			value:    NewUIntValueFromUint64(10),
			expected: "10",
		},
		"UInt8": {
			value:    UInt8Value(8),
			expected: "8",
		},
		"UInt16": {
			value:    UInt16Value(16),
			expected: "16",
		},
		"UInt32": {
			value:    UInt32Value(32),
			expected: "32",
		},
		"UInt64": {
			value:    UInt64Value(64),
			expected: "64",
		},
		"UInt128": {
			value:    NewUInt128ValueFromUint64(128),
			expected: "128",
		},
		"UInt256": {
			value:    NewUInt256ValueFromUint64(256),
			expected: "256",
		},
		"Int8": {
			value:    Int8Value(-8),
			expected: "-8",
		},
		"Int16": {
			value:    Int16Value(-16),
			expected: "-16",
		},
		"Int32": {
			value:    Int32Value(-32),
			expected: "-32",
		},
		"Int64": {
			value:    Int64Value(-64),
			expected: "-64",
		},
		"Int128": {
			value:    NewInt128ValueFromInt64(-128),
			expected: "-128",
		},
		"Int256": {
			value:    NewInt256ValueFromInt64(-256),
			expected: "-256",
		},
		"Word8": {
			value:    Word8Value(8),
			expected: "8",
		},
		"Word16": {
			value:    Word16Value(16),
			expected: "16",
		},
		"Word32": {
			value:    Word32Value(32),
			expected: "32",
		},
		"Word64": {
			value:    Word64Value(64),
			expected: "64",
		},
		"UFix64": {
			value:    NewUFix64ValueWithInteger(64),
			expected: "64.00000000",
		},
		"Fix64": {
			value:    NewFix64ValueWithInteger(-32),
			expected: "-32.00000000",
		},
		"true": {
			value:    BoolValue(true),
			expected: "true",
		},
		"false": {
			value:    BoolValue(false),
			expected: "false",
		},
		"String": {
			value:    NewStringValue("Flow ridah!"),
			expected: "Flow ridah!",
		},
		"Address": {
			value:    NewAddressValue(common.Address{0, 0, 0, 0, 0, 0, 0, 1}),
			expected: "0x1",
		},
		"enum": {
			value: func() HasKeyString {
				members := NewStringValueOrderedMap()
				members.Set("rawValue", UInt8Value(42))
				return NewCompositeValue(
					storage,
					utils.TestLocation,
					"Foo",
					common.CompositeKindEnum,
					members,
					common.Address{},
				)
			}(),
			expected: "42",
		},
		"Path": {
			value: PathValue{
				Domain:     common.PathDomainStorage,
				Identifier: "foo",
			},
			// NOTE: this is an unfortunate mistake,
			// the KeyString function should have been using Domain.Identifier()
			expected: "/PathDomainStorage/foo",
		},
	}

	test := func(name string, testCase testCase) {

		t.Run(name, func(t *testing.T) {

			t.Parallel()

			assert.Equal(t,
				testCase.expected,
				testCase.value.KeyString(),
			)
		})
	}

	for name, testCase := range stringerTests {
		test(name, testCase)
	}
}

func TestBlockValue(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	block := BlockValue{
		Height: 4,
		View:   5,
		ID: NewArrayValue(
			ByteArrayStaticType,
			storage,
		),
		Timestamp: 5.0,
	}

	// static type test
	var actualTs = block.Timestamp
	const expectedTs UFix64Value = 5.0
	assert.Equal(t, expectedTs, actualTs)
}

func TestEphemeralReferenceTypeConformance(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	// Obtain a self referencing (cyclic) ephemeral reference value.

	code := `
        pub fun getEphemeralRef(): &Foo {
            var foo = Foo()
            var fooRef = &foo as &Foo

            // Create the cyclic reference
            fooRef.bar = fooRef

            return fooRef
        }

        pub struct Foo {

            pub(set) var bar: &Foo?

            init() {
                self.bar = nil
            }
        }`

	checker, err := checkerUtils.ParseAndCheckWithOptions(t,
		code,
		checkerUtils.ParseAndCheckOptions{},
	)

	require.NoError(t, err)

	inter, err := NewInterpreter(
		ProgramFromChecker(checker),
		checker.Location,
		WithStorage(storage),
	)

	require.NoError(t, err)

	err = inter.Interpret()
	require.NoError(t, err)

	value, err := inter.Invoke("getEphemeralRef")
	require.NoError(t, err)
	require.IsType(t, &EphemeralReferenceValue{}, value)

	dynamicType := value.DynamicType(inter, SeenReferences{})

	// Check the dynamic type conformance on a cyclic value.
	conforms := value.ConformsToDynamicType(inter, dynamicType, TypeConformanceResults{})
	assert.True(t, conforms)

	// Check against a non-conforming type
	conforms = value.ConformsToDynamicType(inter, EphemeralReferenceDynamicType{}, TypeConformanceResults{})
	assert.False(t, conforms)
}

func TestCapabilityValue_Equal(t *testing.T) {

	t.Parallel()

	t.Run("equal, borrow type", func(t *testing.T) {

		t.Parallel()

		require.True(t,
			(&CapabilityValue{
				Address: AddressValue{0x1},
				Path: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "test",
				},
				BorrowType: PrimitiveStaticTypeInt,
			}).Equal(
				&CapabilityValue{
					Address: AddressValue{0x1},
					Path: PathValue{
						Domain:     common.PathDomainStorage,
						Identifier: "test",
					},
					BorrowType: PrimitiveStaticTypeInt,
				},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("equal, no borrow type", func(t *testing.T) {

		t.Parallel()

		require.True(t,
			(&CapabilityValue{
				Address: AddressValue{0x1},
				Path: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "test",
				},
			}).Equal(
				&CapabilityValue{
					Address: AddressValue{0x1},
					Path: PathValue{
						Domain:     common.PathDomainStorage,
						Identifier: "test",
					},
				},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different paths", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			(&CapabilityValue{
				Address: AddressValue{0x1},
				Path: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "test1",
				},
				BorrowType: PrimitiveStaticTypeInt,
			}).Equal(
				&CapabilityValue{
					Address: AddressValue{0x1},
					Path: PathValue{
						Domain:     common.PathDomainStorage,
						Identifier: "test2",
					},
					BorrowType: PrimitiveStaticTypeInt,
				},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different addresses", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			(&CapabilityValue{
				Address: AddressValue{0x1},
				Path: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "test",
				},
				BorrowType: PrimitiveStaticTypeInt,
			}).Equal(
				&CapabilityValue{
					Address: AddressValue{0x2},
					Path: PathValue{
						Domain:     common.PathDomainStorage,
						Identifier: "test",
					},
					BorrowType: PrimitiveStaticTypeInt,
				},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different borrow types", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			(&CapabilityValue{
				Address: AddressValue{0x1},
				Path: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "test",
				},
				BorrowType: PrimitiveStaticTypeInt,
			}).Equal(
				&CapabilityValue{
					Address: AddressValue{0x1},
					Path: PathValue{
						Domain:     common.PathDomainStorage,
						Identifier: "test",
					},
					BorrowType: PrimitiveStaticTypeString,
				},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			(&CapabilityValue{
				Address: AddressValue{0x1},
				Path: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "test",
				},
				BorrowType: PrimitiveStaticTypeInt,
			}).Equal(
				NewStringValue("test"),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestAddressValue_Equal(t *testing.T) {

	t.Parallel()

	t.Run("equal", func(t *testing.T) {

		t.Parallel()

		require.True(t,
			AddressValue{0x1}.Equal(
				AddressValue{0x1},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			AddressValue{0x1}.Equal(
				AddressValue{0x2},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			AddressValue{0x1}.Equal(
				UInt8Value(1),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestBoolValue_Equal(t *testing.T) {

	t.Parallel()

	t.Run("equal true", func(t *testing.T) {

		t.Parallel()

		require.True(t,
			BoolValue(true).Equal(
				BoolValue(true),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("equal false", func(t *testing.T) {

		t.Parallel()

		require.True(t,
			BoolValue(false).Equal(
				BoolValue(false),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			BoolValue(true).Equal(
				BoolValue(false),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			BoolValue(true).Equal(
				UInt8Value(1),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestStringValue_Equal(t *testing.T) {

	t.Parallel()

	t.Run("equal", func(t *testing.T) {

		t.Parallel()

		require.True(t,
			NewStringValue("test").Equal(
				NewStringValue("test"),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			NewStringValue("test").Equal(
				NewStringValue("foo"),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			NewStringValue("1").Equal(
				UInt8Value(1),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestNilValue_Equal(t *testing.T) {

	t.Parallel()

	t.Run("equal", func(t *testing.T) {

		t.Parallel()

		require.True(t,
			NilValue{}.Equal(
				NilValue{},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			NilValue{}.Equal(
				UInt8Value(0),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestSomeValue_Equal(t *testing.T) {

	t.Parallel()

	t.Run("equal", func(t *testing.T) {

		t.Parallel()

		require.True(t,
			NewSomeValueNonCopying(NewStringValue("test")).Equal(
				NewSomeValueNonCopying(NewStringValue("test")),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			NewSomeValueNonCopying(NewStringValue("test")).Equal(
				NewSomeValueNonCopying(NewStringValue("foo")),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			NewSomeValueNonCopying(NewStringValue("1")).Equal(
				UInt8Value(1),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestTypeValue_Equal(t *testing.T) {

	t.Parallel()

	t.Run("equal", func(t *testing.T) {

		t.Parallel()

		require.True(t,
			TypeValue{
				Type: PrimitiveStaticTypeString,
			}.Equal(
				TypeValue{
					Type: PrimitiveStaticTypeString,
				},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			TypeValue{
				Type: PrimitiveStaticTypeString,
			}.Equal(
				TypeValue{
					Type: PrimitiveStaticTypeInt,
				},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			TypeValue{
				Type: PrimitiveStaticTypeString,
			}.Equal(
				NewStringValue("String"),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestPathValue_Equal(t *testing.T) {

	t.Parallel()

	for _, domain := range common.AllPathDomains {

		t.Run(fmt.Sprintf("equal, %s", domain), func(t *testing.T) {

			require.True(t,
				PathValue{
					Domain:     domain,
					Identifier: "test",
				}.Equal(
					PathValue{
						Domain:     domain,
						Identifier: "test",
					},
					ReturnEmptyLocationRange,
				),
			)
		})
	}

	for _, domain := range common.AllPathDomains {
		for _, otherDomain := range common.AllPathDomains {

			if domain == otherDomain {
				continue
			}

			t.Run(fmt.Sprintf("different domains %s %s", domain, otherDomain), func(t *testing.T) {

				require.False(t,
					PathValue{
						Domain:     domain,
						Identifier: "test",
					}.Equal(
						PathValue{
							Domain:     otherDomain,
							Identifier: "test",
						},
						ReturnEmptyLocationRange,
					),
				)
			})
		}
	}

	for _, domain := range common.AllPathDomains {

		t.Run(fmt.Sprintf("different identifiers, %s", domain), func(t *testing.T) {

			require.False(t,
				PathValue{
					Domain:     domain,
					Identifier: "test1",
				}.Equal(
					PathValue{
						Domain:     domain,
						Identifier: "test2",
					},
					ReturnEmptyLocationRange,
				),
			)
		})
	}

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			PathValue{
				Domain:     common.PathDomainStorage,
				Identifier: "test",
			}.Equal(
				NewStringValue("/storage/test"),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestLinkValue_Equal(t *testing.T) {

	t.Parallel()

	t.Run("equal, borrow type", func(t *testing.T) {

		t.Parallel()

		require.True(t,
			LinkValue{
				TargetPath: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "test",
				},
				Type: PrimitiveStaticTypeInt,
			}.Equal(
				LinkValue{
					TargetPath: PathValue{
						Domain:     common.PathDomainStorage,
						Identifier: "test",
					},
					Type: PrimitiveStaticTypeInt,
				},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different paths", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			LinkValue{
				TargetPath: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "test1",
				},
				Type: PrimitiveStaticTypeInt,
			}.Equal(
				LinkValue{
					TargetPath: PathValue{
						Domain:     common.PathDomainStorage,
						Identifier: "test2",
					},
					Type: PrimitiveStaticTypeInt,
				},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different types", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			LinkValue{
				TargetPath: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "test",
				},
				Type: PrimitiveStaticTypeInt,
			}.Equal(
				LinkValue{
					TargetPath: PathValue{
						Domain:     common.PathDomainStorage,
						Identifier: "test",
					},
					Type: PrimitiveStaticTypeString,
				},
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		require.False(t,
			LinkValue{
				TargetPath: PathValue{
					Domain:     common.PathDomainStorage,
					Identifier: "test",
				},
				Type: PrimitiveStaticTypeInt,
			}.Equal(
				NewStringValue("test"),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestArrayValue_Equal(t *testing.T) {

	t.Parallel()

	uint8ArrayStaticType := VariableSizedStaticType{
		Type: PrimitiveStaticTypeUInt8,
	}

	t.Run("equal", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.True(t,
			NewArrayValue(
				uint8ArrayStaticType,
				storage,
				UInt8Value(1),
				UInt8Value(2),
			).Equal(
				NewArrayValue(
					uint8ArrayStaticType,
					storage,
					UInt8Value(1),
					UInt8Value(2),
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different elements", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewArrayValue(
				uint8ArrayStaticType,
				storage,
				UInt8Value(1),
				UInt8Value(2),
			).Equal(
				NewArrayValue(
					uint8ArrayStaticType,
					storage,
					UInt8Value(2),
					UInt8Value(3),
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("more elements", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewArrayValue(
				uint8ArrayStaticType,
				storage,
				UInt8Value(1),
			).Equal(
				NewArrayValue(
					uint8ArrayStaticType,
					storage,
					UInt8Value(1),
					UInt8Value(2),
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("fewer elements", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewArrayValue(
				uint8ArrayStaticType,
				storage,
				UInt8Value(1),
				UInt8Value(2),
			).Equal(
				NewArrayValue(
					uint8ArrayStaticType,
					storage,
					UInt8Value(1),
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different types", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		uint16ArrayStaticType := VariableSizedStaticType{
			Type: PrimitiveStaticTypeUInt16,
		}

		require.False(t,
			NewArrayValue(
				uint8ArrayStaticType,
				storage,
			).Equal(
				NewArrayValue(
					uint16ArrayStaticType,
					storage,
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("no type, type", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewArrayValue(
				nil,
				storage,
			).Equal(
				NewArrayValue(
					uint8ArrayStaticType,
					storage,
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("type, no type", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewArrayValue(
				uint8ArrayStaticType,
				storage,
			).Equal(
				NewArrayValue(
					nil,
					storage,
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("no types", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.True(t,
			NewArrayValue(
				nil,
				storage,
			).Equal(
				NewArrayValue(
					nil,
					storage,
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewArrayValue(
				uint8ArrayStaticType,
				storage,
				UInt8Value(1),
			).Equal(
				UInt8Value(1),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestDictionaryValue_Equal(t *testing.T) {

	t.Parallel()

	byteStringDictionaryType := DictionaryStaticType{
		KeyType:   PrimitiveStaticTypeUInt8,
		ValueType: PrimitiveStaticTypeString,
	}

	t.Run("equal", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.True(t,
			NewDictionaryValue(
				byteStringDictionaryType,
				storage,
				UInt8Value(1),
				NewStringValue("1"),
				UInt8Value(2),
				NewStringValue("2"),
			).Equal(
				NewDictionaryValue(
					byteStringDictionaryType,
					storage,
					UInt8Value(1),
					NewStringValue("1"),
					UInt8Value(2),
					NewStringValue("2"),
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different keys", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewDictionaryValue(
				byteStringDictionaryType,
				storage,
				UInt8Value(1),
				NewStringValue("1"),
				UInt8Value(2),
				NewStringValue("2"),
			).Equal(
				NewDictionaryValue(
					byteStringDictionaryType,
					storage,
					UInt8Value(2),
					NewStringValue("1"),
					UInt8Value(3),
					NewStringValue("2"),
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different values", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewDictionaryValue(
				byteStringDictionaryType,
				storage,
				UInt8Value(1),
				NewStringValue("1"),
				UInt8Value(2),
				NewStringValue("2"),
			).Equal(
				NewDictionaryValue(
					byteStringDictionaryType,
					storage,
					UInt8Value(1),
					NewStringValue("2"),
					UInt8Value(2),
					NewStringValue("3"),
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("more elements", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewDictionaryValue(
				byteStringDictionaryType,
				storage,
				UInt8Value(1),
				NewStringValue("1"),
			).Equal(
				NewDictionaryValue(
					byteStringDictionaryType,
					storage,
					UInt8Value(1),
					NewStringValue("1"),
					UInt8Value(2),
					NewStringValue("2"),
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("fewer elements", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewDictionaryValue(
				byteStringDictionaryType,
				storage,
				UInt8Value(1),
				NewStringValue("1"),
				UInt8Value(2),
				NewStringValue("2"),
			).Equal(
				NewDictionaryValue(
					byteStringDictionaryType,
					storage,
					UInt8Value(1),
					NewStringValue("1"),
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different types", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		stringByteDictionaryStaticType := DictionaryStaticType{
			KeyType:   PrimitiveStaticTypeString,
			ValueType: PrimitiveStaticTypeUInt8,
		}

		require.False(t,
			NewDictionaryValue(
				byteStringDictionaryType,
				storage,
			).Equal(
				NewDictionaryValue(
					stringByteDictionaryStaticType,
					storage,
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different kind", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		require.False(t,
			NewDictionaryValue(
				byteStringDictionaryType,
				storage,
				UInt8Value(1),
				NewStringValue("1"),
				UInt8Value(2),
				NewStringValue("2"),
			).Equal(
				NewArrayValue(
					ByteArrayStaticType,
					storage,
					UInt8Value(1),
					UInt8Value(2),
				),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestCompositeValue_Equal(t *testing.T) {

	t.Parallel()

	t.Run("equal", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		fields1 := NewStringValueOrderedMap()
		fields1.Set("a", NewStringValue("a"))

		fields2 := NewStringValueOrderedMap()
		fields2.Set("a", NewStringValue("a"))

		require.True(t,
			NewCompositeValue(
				storage,
				utils.TestLocation,
				"X",
				common.CompositeKindStructure,
				fields1,
				common.Address{},
			).Equal(
				NewCompositeValue(
					storage,
					utils.TestLocation,
					"X",
					common.CompositeKindStructure,
					fields2,
					common.Address{},
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different location", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		fields1 := NewStringValueOrderedMap()
		fields1.Set("a", NewStringValue("a"))

		fields2 := NewStringValueOrderedMap()
		fields2.Set("a", NewStringValue("a"))

		require.False(t,
			NewCompositeValue(
				storage,
				common.IdentifierLocation("A"),
				"X",
				common.CompositeKindStructure,
				fields1,
				common.Address{},
			).Equal(
				NewCompositeValue(
					storage,
					common.IdentifierLocation("B"),
					"X",
					common.CompositeKindStructure,
					fields2,
					common.Address{},
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different identifier", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		fields1 := NewStringValueOrderedMap()
		fields1.Set("a", NewStringValue("a"))

		fields2 := NewStringValueOrderedMap()
		fields2.Set("a", NewStringValue("a"))

		require.False(t,
			NewCompositeValue(
				storage,
				common.IdentifierLocation("A"),
				"X",
				common.CompositeKindStructure,
				fields1,
				common.Address{},
			).Equal(
				NewCompositeValue(
					storage,
					common.IdentifierLocation("A"),
					"Y",
					common.CompositeKindStructure,
					fields2,
					common.Address{},
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different fields", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		fields1 := NewStringValueOrderedMap()
		fields1.Set("a", NewStringValue("a"))

		fields2 := NewStringValueOrderedMap()
		fields2.Set("a", NewStringValue("b"))

		require.False(t,
			NewCompositeValue(
				storage,
				common.IdentifierLocation("A"),
				"X",
				common.CompositeKindStructure,
				fields1,
				common.Address{},
			).Equal(
				NewCompositeValue(
					storage,
					common.IdentifierLocation("A"),
					"X",
					common.CompositeKindStructure,
					fields2,
					common.Address{},
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("more fields", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		fields1 := NewStringValueOrderedMap()
		fields1.Set("a", NewStringValue("a"))

		fields2 := NewStringValueOrderedMap()
		fields2.Set("a", NewStringValue("a"))
		fields2.Set("b", NewStringValue("b"))

		require.False(t,
			NewCompositeValue(
				storage,
				common.IdentifierLocation("A"),
				"X",
				common.CompositeKindStructure,
				fields1,
				common.Address{},
			).Equal(
				NewCompositeValue(
					storage,
					common.IdentifierLocation("A"),
					"X",
					common.CompositeKindStructure,
					fields2,
					common.Address{},
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("fewer fields", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		fields1 := NewStringValueOrderedMap()
		fields1.Set("a", NewStringValue("a"))
		fields1.Set("b", NewStringValue("b"))

		fields2 := NewStringValueOrderedMap()
		fields2.Set("a", NewStringValue("a"))

		require.False(t,
			NewCompositeValue(
				storage,
				common.IdentifierLocation("A"),
				"X",
				common.CompositeKindStructure,
				fields1,
				common.Address{},
			).Equal(
				NewCompositeValue(
					storage,
					common.IdentifierLocation("A"),
					"X",
					common.CompositeKindStructure,
					fields2,
					common.Address{},
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different composite kind", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		fields1 := NewStringValueOrderedMap()
		fields1.Set("a", NewStringValue("a"))

		fields2 := NewStringValueOrderedMap()
		fields2.Set("a", NewStringValue("a"))

		require.False(t,
			NewCompositeValue(
				storage,
				common.IdentifierLocation("A"),
				"X",
				common.CompositeKindStructure,
				fields1,
				common.Address{},
			).Equal(
				NewCompositeValue(
					storage,
					common.IdentifierLocation("A"),
					"X",
					common.CompositeKindResource,
					fields2,
					common.Address{},
				),
				ReturnEmptyLocationRange,
			),
		)
	})

	t.Run("different composite kind", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		fields1 := NewStringValueOrderedMap()
		fields1.Set("a", NewStringValue("a"))

		require.False(t,
			NewCompositeValue(
				storage,
				common.IdentifierLocation("A"),
				"X",
				common.CompositeKindStructure,
				fields1,
				common.Address{},
			).Equal(
				NewStringValue("test"),
				ReturnEmptyLocationRange,
			),
		)
	})
}

func TestNumberValue_Equal(t *testing.T) {

	t.Parallel()

	testValues := map[string]EquatableValue{
		"UInt":    NewUIntValueFromUint64(10),
		"UInt8":   UInt8Value(8),
		"UInt16":  UInt16Value(16),
		"UInt32":  UInt32Value(32),
		"UInt64":  UInt64Value(64),
		"UInt128": NewUInt128ValueFromUint64(128),
		"UInt256": NewUInt256ValueFromUint64(256),
		"Int8":    Int8Value(-8),
		"Int16":   Int16Value(-16),
		"Int32":   Int32Value(-32),
		"Int64":   Int64Value(-64),
		"Int128":  NewInt128ValueFromInt64(-128),
		"Int256":  NewInt256ValueFromInt64(-256),
		"Word8":   Word8Value(8),
		"Word16":  Word16Value(16),
		"Word32":  Word32Value(32),
		"Word64":  Word64Value(64),
		"UFix64":  NewUFix64ValueWithInteger(64),
		"Fix64":   NewFix64ValueWithInteger(-32),
	}

	for name, value := range testValues {

		t.Run(fmt.Sprintf("equal, %s", name), func(t *testing.T) {

			require.True(t,
				value.Equal(
					value,
					ReturnEmptyLocationRange,
				),
			)
		})
	}

	for name, value := range testValues {
		for otherName, otherValue := range testValues {

			if name == otherName {
				continue
			}

			t.Run(fmt.Sprintf("unequal, %s %s", name, otherName), func(t *testing.T) {

				require.False(t,
					value.Equal(
						otherValue,
						ReturnEmptyLocationRange,
					),
				)
			})
		}
	}

	for name, value := range testValues {

		t.Run(fmt.Sprintf("different kind, %s", name), func(t *testing.T) {

			t.Parallel()

			require.False(t,
				value.Equal(
					AddressValue{0x1},
					ReturnEmptyLocationRange,
				),
			)
		})
	}
}

func TestPublicKeyValue(t *testing.T) {

	t.Parallel()

	t.Run("Stringer output includes public key value", func(t *testing.T) {

		t.Parallel()

		storage := NewInMemoryStorage()

		publicKey := NewArrayValue(
			VariableSizedStaticType{
				Type: PrimitiveStaticTypeInt,
			},
			storage,
			NewIntValueFromInt64(1),
			NewIntValueFromInt64(7),
			NewIntValueFromInt64(3),
		)

		publicKeyString := "[1, 7, 3]"

		sigAlgo := func() *CompositeValue {
			fields := NewStringValueOrderedMap()
			fields.Set(sema.EnumRawValueFieldName, UInt8Value(sema.SignatureAlgorithmECDSA_secp256k1.RawValue()))

			return &CompositeValue{
				QualifiedIdentifier: sema.SignatureAlgorithmType.QualifiedIdentifier(),
				Kind:                sema.SignatureAlgorithmType.Kind,
				Fields:              fields,
			}
		}

		interpreter, err := NewInterpreter(
			nil,
			utils.TestLocation,
			WithStorage(storage),
			WithPublicKeyValidationHandler(
				func(_ *Interpreter, _ *CompositeValue) BoolValue {
					return true
				},
			),
		)
		require.NoError(t, err)

		key := NewPublicKeyValue(
			interpreter,
			publicKey,
			sigAlgo(),
			interpreter.PublicKeyValidationHandler,
		)

		require.Contains(t,
			key.String(),
			publicKeyString,
		)
	})
}

func TestNonStorable(t *testing.T) {

	t.Parallel()

	storage := NewInMemoryStorage()

	code := `
      pub struct Foo {

          let bar: &Int?

          init() {
              self.bar = &1 as &Int
          }
      }

      fun foo(): &Int? {
          return Foo().bar
      }
    `

	checker, err := checkerUtils.ParseAndCheckWithOptions(t,
		code,
		checkerUtils.ParseAndCheckOptions{},
	)

	require.NoError(t, err)

	inter, err := NewInterpreter(
		ProgramFromChecker(checker),
		checker.Location,
		WithStorage(storage),
	)

	require.NoError(t, err)

	err = inter.Interpret()
	require.NoError(t, err)

	_, err = inter.Invoke("foo")
	require.NoError(t, err)

}

func TestHashable(t *testing.T) {

	// Assert that all Value and DynamicType implementations are hashable

	pkgs, err := packages.Load(
		&packages.Config{
			// https://github.com/golang/go/issues/45218
			Mode: packages.NeedImports | packages.NeedTypes,
		},
		"github.com/onflow/cadence/runtime/interpreter",
	)
	require.NoError(t, err)

	pkg := pkgs[0]
	scope := pkg.Types.Scope()

	test := func(interfaceName string) {

		t.Run(interfaceName, func(t *testing.T) {

			interfaceType, ok := scope.Lookup(interfaceName).Type().Underlying().(*types.Interface)
			require.True(t, ok)

			for _, name := range scope.Names() {
				object := scope.Lookup(name)
				_, ok := object.(*types.TypeName)
				if !ok {
					continue
				}

				implementationType := object.Type()
				if !types.Implements(implementationType, interfaceType) {
					continue
				}

				err := checkHashable(implementationType)
				if !assert.NoError(t,
					err,
					"%s implementation is not hashable: %s",
					interfaceType.String(),
					implementationType,
				) {
					continue
				}
			}
		})
	}

	test("Value")
	test("DynamicType")
}

func checkHashable(ty types.Type) error {

	// TODO: extend the notion of unhashable types,
	//  see https://github.com/golang/go/blob/a22e3172200d4bdd0afcbbe6564dbb67fea4b03a/src/runtime/alg.go#L144

	switch ty := ty.(type) {
	case *types.Basic:
		switch ty.Kind() {
		case types.Bool,
			types.Int,
			types.Int8,
			types.Int16,
			types.Int32,
			types.Int64,
			types.Uint,
			types.Uint8,
			types.Uint16,
			types.Uint32,
			types.Uint64,
			types.Float32,
			types.Float64,
			types.String:
			return nil
		}
	case *types.Pointer,
		*types.Array,
		*types.Interface:
		return nil

	case *types.Struct:
		numFields := ty.NumFields()
		for i := 0; i < numFields; i++ {
			field := ty.Field(i)
			fieldTy := field.Type()
			err := checkHashable(fieldTy)
			if err != nil {
				return fmt.Errorf(
					"struct type has unhashable field %s: %w",
					field.Name(),
					err,
				)
			}
		}
		return nil

	case *types.Named:
		return checkHashable(ty.Underlying())
	}

	return fmt.Errorf(
		"type %s is potentially not hashable",
		ty.String(),
	)
}
