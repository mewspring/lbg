package main

import (
	"fmt"

	"github.com/llir/l/ir/types"
	"github.com/mewmew/lbg/internal/syntax"
)

// llType translates the given Go type to an equivalent LLVM IR type.
func (c *compiler) llType(t syntax.Expr) types.Type {
	switch t := t.(type) {
	case *syntax.Name:
		return c.llNamedType(t)
	case *syntax.FuncType:
		return c.llFuncType(t)
	default:
		panic(fmt.Errorf("support for %T not yet implemented", t))
	}
}

// TODO: add *compiler context for llNamedType, to resolve type from qualified
// identifiers.

func createUniverseTypes() map[string]types.Type {
	// bool
	boolType := types.NewInt(1)
	boolType.Alias = "bool"

	// byte
	byteType := types.NewInt(8)
	byteType.Alias = "byte"

	// complex64
	complex64Type := types.NewStruct(types.Float, types.Float)
	complex64Type.Alias = "complex64"

	// complex128
	complex128Type := types.NewStruct(types.Double, types.Double)
	complex128Type.Alias = "complex128"

	// TODO: add error interface.

	// float32
	float32Type := &types.FloatType{
		Alias: "float32",
		Kind:  types.FloatKindFloat,
	}

	// float64
	float64Type := &types.FloatType{
		Alias: "float64",
		Kind:  types.FloatKindDouble,
	}

	// int
	// TODO: set size based on target architecture?
	intType := types.NewInt(32)
	intType.Alias = "int"

	// int8
	int8Type := types.NewInt(8)
	int8Type.Alias = "int8"

	// int16
	int16Type := types.NewInt(16)
	int16Type.Alias = "int16"

	// int32
	int32Type := types.NewInt(32)
	int32Type.Alias = "int32"

	// int64
	int64Type := types.NewInt(64)
	int64Type.Alias = "int64"

	// rune
	runeType := types.NewInt(32)
	runeType.Alias = "rune"

	// string
	data := types.NewPointer(types.NewArray(0, types.I8))
	stringType := types.NewStruct(data, intType)
	stringType.Alias = "string"

	// uint8
	uint8Type := types.NewInt(8)
	uint8Type.Alias = "uint8"

	// uint16
	uint16Type := types.NewInt(16)
	uint16Type.Alias = "uint16"

	// uint32
	uint32Type := types.NewInt(32)
	uint32Type.Alias = "uint32"

	// uint64
	uint64Type := types.NewInt(64)
	uint64Type.Alias = "uint64"

	universe := map[string]types.Type{
		"bool":       boolType,
		"byte":       byteType,
		"complex64":  complex64Type,
		"complex128": complex128Type,
		"float32":    float32Type,
		"float64":    float64Type,
		"int":        intType,
		"int8":       int8Type,
		"int16":      int16Type,
		"int32":      int32Type,
		"int64":      int64Type,
		"rune":       runeType,
		"string":     stringType,
		"uint8":      uint8Type,
		"uint16":     uint16Type,
		"uint32":     uint32Type,
		"uint64":     uint64Type,
	}
	return universe
}

// llNamedType translates the given named Go type to an equivalent LLVM IR type.
func (c *compiler) llNamedType(t *syntax.Name) types.Type {
	//universe := createUniverseTypes()
	// TODO: resolve identifiers from scope, where the universe scope would
	// contain the predeclared types
	//if typ, ok := c.scope.findType(t.Value); ok {
	//	return typ
	//}
	panic(fmt.Errorf("support for named type %q not yet implemented", t.Value))
}

// llFuncType translates the given Go function type to an equivalent LLVM IR
// type.
func (c *compiler) llFuncType(t *syntax.FuncType) *types.FuncType {
	var retType types.Type
	switch len(t.ResultList) {
	case 0:
		retType = types.Void
	case 1:
		retType = c.llType(t.ResultList[0].Type)
	default:
		var resultTypes []types.Type
		for _, result := range t.ResultList {
			resultType := c.llType(result.Type)
			resultTypes = append(resultTypes, resultType)
		}
		retType = types.NewStruct(resultTypes...)
	}
	var paramTypes []types.Type
	variadic := false
	for i, param := range t.ParamList {
		if _, ok := param.Type.(*syntax.DotsType); ok {
			if i != len(t.ParamList)-1 {
				// TODO: report error through *compiler instead of panic.
				panic(fmt.Errorf("invalid use of variadic parameter type for parameter %q (param %d of %d); must be the last parameter", param.Name.Value, i+1, len(t.ParamList)))
			}
			variadic = true
			continue
		}
		paramType := c.llType(param.Type)
		paramTypes = append(paramTypes, paramType)
	}
	typ := types.NewFunc(retType, paramTypes...)
	typ.Variadic = variadic
	return typ
}
