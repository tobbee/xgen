// Copyright 2020 The xgen Authors. All rights reserved. Use of this source
// code is governed by a BSD-style license that can be found in the LICENSE
// file.
//
// Package xgen written in pure Go providing a set of functions that allow you
// to parse XSD (XML schema files). This library needs Go version 1.10 or
// later.

package xgen

import "encoding/xml"

func (opt *Options) OnGroup(ele xml.StartElement, protoTree []interface{}) (err error) {
	group := Group{}
	for _, attr := range ele.Attr {
		if attr.Name.Local == "name" {
			group.Name = attr.Value
		}
		if attr.Name.Local == "ref" {
			group.Name = attr.Value
			group.Ref, err = opt.GetValueType(attr.Value, protoTree)
			if err != nil {
				return
			}
		}
		if attr.Name.Local == "maxOccurs" {
			if attr.Value != "0" {
				group.Plural = true
			}
		}
	}
	if opt.ComplexType.Len() == 0 {
		if opt.InGroup == 0 {
			opt.InGroup++
			opt.CurrentEle = opt.InElement
			opt.Group = &group
			return
		}
		if opt.InGroup > 0 {
			opt.InGroup++
			opt.Group.Groups = append(opt.Group.Groups, group)
			return
		}

	}
	if opt.ComplexType.Len() > 0 {
		opt.ComplexType.Peek().(*ComplexType).Groups = append(opt.ComplexType.Peek().(*ComplexType).Groups, group)
		return
	}
	return
}
