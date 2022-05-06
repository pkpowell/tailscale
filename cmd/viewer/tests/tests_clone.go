// Copyright (c) 2022 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by tailscale.com/cmd/cloner; DO NOT EDIT.

package tests

import (
	"inet.af/netaddr"
)

// Clone makes a deep copy of StructWithPtrs.
// The result aliases no memory with the original.
func (src *StructWithPtrs) Clone() *StructWithPtrs {
	if src == nil {
		return nil
	}
	dst := new(StructWithPtrs)
	*dst = *src
	if dst.Value != nil {
		dst.Value = new(StructWithoutPtrs)
		*dst.Value = *src.Value
	}
	if dst.Int != nil {
		dst.Int = new(int)
		*dst.Int = *src.Int
	}
	return dst
}

// A compilation failure here means this code must be regenerated, with the command at the top of this file.
var _StructWithPtrsCloneNeedsRegeneration = StructWithPtrs(struct {
	Value        *StructWithoutPtrs
	Int          *int
	NoCloneValue *StructWithoutPtrs
}{})

// Clone makes a deep copy of StructWithoutPtrs.
// The result aliases no memory with the original.
func (src *StructWithoutPtrs) Clone() *StructWithoutPtrs {
	if src == nil {
		return nil
	}
	dst := new(StructWithoutPtrs)
	*dst = *src
	return dst
}

// A compilation failure here means this code must be regenerated, with the command at the top of this file.
var _StructWithoutPtrsCloneNeedsRegeneration = StructWithoutPtrs(struct {
	Int int
	Pfx netaddr.IPPrefix
}{})

// Clone makes a deep copy of Map.
// The result aliases no memory with the original.
func (src *Map) Clone() *Map {
	if src == nil {
		return nil
	}
	dst := new(Map)
	*dst = *src
	if dst.M != nil {
		dst.M = map[string]int{}
		for k, v := range src.M {
			dst.M[k] = v
		}
	}
	return dst
}

// A compilation failure here means this code must be regenerated, with the command at the top of this file.
var _MapCloneNeedsRegeneration = Map(struct {
	M map[string]int
}{})

// Clone makes a deep copy of StructWithSlices.
// The result aliases no memory with the original.
func (src *StructWithSlices) Clone() *StructWithSlices {
	if src == nil {
		return nil
	}
	dst := new(StructWithSlices)
	*dst = *src
	dst.Values = append(src.Values[:0:0], src.Values...)
	dst.ValuePointers = make([]*StructWithoutPtrs, len(src.ValuePointers))
	for i := range dst.ValuePointers {
		dst.ValuePointers[i] = src.ValuePointers[i].Clone()
	}
	dst.StructPointers = make([]*StructWithPtrs, len(src.StructPointers))
	for i := range dst.StructPointers {
		dst.StructPointers[i] = src.StructPointers[i].Clone()
	}
	dst.Structs = make([]StructWithPtrs, len(src.Structs))
	for i := range dst.Structs {
		dst.Structs[i] = *src.Structs[i].Clone()
	}
	dst.Ints = make([]*int, len(src.Ints))
	for i := range dst.Ints {
		x := *src.Ints[i]
		dst.Ints[i] = &x
	}
	dst.Slice = append(src.Slice[:0:0], src.Slice...)
	dst.Prefixes = append(src.Prefixes[:0:0], src.Prefixes...)
	dst.Data = append(src.Data[:0:0], src.Data...)
	return dst
}

// A compilation failure here means this code must be regenerated, with the command at the top of this file.
var _StructWithSlicesCloneNeedsRegeneration = StructWithSlices(struct {
	Values         []StructWithoutPtrs
	ValuePointers  []*StructWithoutPtrs
	StructPointers []*StructWithPtrs
	Structs        []StructWithPtrs
	Ints           []*int
	Slice          []string
	Prefixes       []netaddr.IPPrefix
	Data           []byte
}{})
