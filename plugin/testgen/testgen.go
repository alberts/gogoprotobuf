// Copyright (c) 2013, Vastech SA (PTY) LTD. All rights reserved.
// http://code.google.com/p/gogoprotobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

/*
The testgen plugin generates Test and Benchmark functions for each message.

Tests are enabled using the following extensions:

  - testgen
  - testgen_all

Benchmarks are enabled using the following extensions:

  - benchgen
  - benchgen_all

Let us look at:

  code.google.com/p/gogoprotobuf/test/example/example.proto

Btw all the output can be seen at:

  code.google.com/p/gogoprotobuf/test/example/*

The following message:

  option (gogoproto.testgen_all) = true;
  option (gogoproto.benchgen_all) = true;

  message A {
	optional string Description = 1 [(gogoproto.nullable) = false];
	optional int64 Number = 2 [(gogoproto.nullable) = false];
	optional bytes Id = 3 [(gogoproto.customtype) = "code.google.com/p/gogoprotobuf/test/custom.Uuid", (gogoproto.nullable) = false];
  }

given to the testgen plugin, will generate the following test code:

	func TestAProto(t *testing.T) {
		popr := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
		p := NewPopulatedA(popr, false)
		data, err := code_google_com_p_gogoprotobuf_proto.Marshal(p)
		if err != nil {
			panic(err)
		}
		msg := &A{}
		if err := code_google_com_p_gogoprotobuf_proto.Unmarshal(data, msg); err != nil {
			panic(err)
		}
		for i := range data {
			data[i] = byte(popr.Intn(256))
		}
		if err := p.VerboseEqual(msg); err != nil {
			t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)
		}
		if !p.Equal(msg) {
			t.Fatalf("%#v !Proto %#v", msg, p)
		}
	}

	func BenchmarkAProtoMarshal(b *testing.B) {
		popr := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
		total := 0
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			p := NewPopulatedA(popr, true)
			b.StartTimer()
			data, err := code_google_com_p_gogoprotobuf_proto.Marshal(p)
			if err != nil {
				panic(err)
			}
			b.StopTimer()
			total += len(data)
		}
		b.SetBytes(int64(total / b.N))
	}

	func BenchmarkAProtoUnmarshal(b *testing.B) {
		popr := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
		total := 0
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			p := NewPopulatedA(popr, true)
			data, err := code_google_com_p_gogoprotobuf_proto.Marshal(p)
			if err != nil {
				panic(err)
			}
			msg := &A{}
			total += len(data)
			b.StartTimer()
			if err := code_google_com_p_gogoprotobuf_proto.Unmarshal(data, msg); err != nil {
				panic(err)
			}
		}
		b.SetBytes(int64(total / b.N))
	}


	func TestAJSON(t *testing1.T) {
		popr := math_rand1.New(math_rand1.NewSource(time1.Now().UnixNano()))
		p := NewPopulatedA(popr, true)
		jsondata, err := encoding_json.Marshal(p)
		if err != nil {
			panic(err)
		}
		msg := &A{}
		err = encoding_json.Unmarshal(jsondata, msg)
		if err != nil {
			panic(err)
		}
		if err := p.VerboseEqual(msg); err != nil {
			t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)
		}
		if !p.Equal(msg) {
			t.Fatalf("%#v !Json Equal %#v", msg, p)
		}
	}

	func TestAProtoText(t *testing2.T) {
		popr := math_rand2.New(math_rand2.NewSource(time2.Now().UnixNano()))
		p := NewPopulatedA(popr, true)
		data := code_google_com_p_gogoprotobuf_proto1.MarshalTextString(p)
		msg := &A{}
		if err := code_google_com_p_gogoprotobuf_proto1.UnmarshalText(data, msg); err != nil {
			panic(err)
		}
		if err := p.VerboseEqual(msg); err != nil {
			t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)
		}
		if !p.Equal(msg) {
			t.Fatalf("%#v !Proto %#v", msg, p)
		}
	}

	func TestAProtoCompactText(t *testing2.T) {
		popr := math_rand2.New(math_rand2.NewSource(time2.Now().UnixNano()))
		p := NewPopulatedA(popr, true)
		data := code_google_com_p_gogoprotobuf_proto1.CompactTextString(p)
		msg := &A{}
		if err := code_google_com_p_gogoprotobuf_proto1.UnmarshalText(data, msg); err != nil {
			panic(err)
		}
		if err := p.VerboseEqual(msg); err != nil {
			t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)
		}
		if !p.Equal(msg) {
			t.Fatalf("%#v !Proto %#v", msg, p)
		}
	}

Other registered tests are also generated.
Tests are registered to this test plugin by calling the following function.

  func RegisterTestPlugin(newFunc NewTestPlugin)

where NewTestPlugin is:

  type NewTestPlugin func(g *generator.Generator) TestPlugin

and TestPlugin is an interface:

  type TestPlugin interface {
	Generate(imports generator.PluginImports, file *generator.FileDescriptor) (used bool)
  }

Plugins that use this interface include:

  - populate
  - gostring
  - equal
  - union
  - and more

Please look at these plugins as examples of how to create your own.
A good idea is to let each plugin generate its own tests.

*/
package testgen

import (
	"code.google.com/p/gogoprotobuf/gogoproto"
	"code.google.com/p/gogoprotobuf/protoc-gen-gogo/generator"
)

type TestPlugin interface {
	Generate(imports generator.PluginImports, file *generator.FileDescriptor) (used bool)
}

type NewTestPlugin func(g *generator.Generator) TestPlugin

var testplugins = make([]NewTestPlugin, 0)

func RegisterTestPlugin(newFunc NewTestPlugin) {
	testplugins = append(testplugins, newFunc)
}

type plugin struct {
	*generator.Generator
	generator.PluginImports
	tests []TestPlugin
}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return "testgen"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
	p.tests = make([]TestPlugin, 0, len(testplugins))
	for i := range testplugins {
		p.tests = append(p.tests, testplugins[i](g))
	}
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	atLeastOne := false
	for i := range p.tests {
		used := p.tests[i].Generate(p.PluginImports, file)
		if used {
			atLeastOne = true
		}
	}
	if atLeastOne {
		p.P(`//These tests are generated by code.google.com/p/gogoprotobuf/plugin/testgen`)
	}
}

type testProto struct {
	*generator.Generator
}

func newProto(g *generator.Generator) TestPlugin {
	return &testProto{g}
}

func (p *testProto) Generate(imports generator.PluginImports, file *generator.FileDescriptor) bool {
	used := false
	testingPkg := imports.NewImport("testing")
	randPkg := imports.NewImport("math/rand")
	timePkg := imports.NewImport("time")
	protoPkg := imports.NewImport("code.google.com/p/gogoprotobuf/proto")
	for _, message := range file.Messages() {
		ccTypeName := generator.CamelCaseSlice(message.TypeName())
		if gogoproto.HasTestGen(file.FileDescriptorProto, message.DescriptorProto) {
			used = true

			p.P(`func Test`, ccTypeName, `Proto(t *`, testingPkg.Use(), `.T) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
			p.P(`p := NewPopulated`, ccTypeName, `(popr, false)`)
			p.P(`data, err := `, protoPkg.Use(), `.Marshal(p)`)
			p.P(`if err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			p.P(`msg := &`, ccTypeName, `{}`)
			p.P(`if err := `, protoPkg.Use(), `.Unmarshal(data, msg); err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			p.P(`for i := range data {`)
			p.In()
			p.P(`data[i] = byte(popr.Intn(256))`)
			p.Out()
			p.P(`}`)
			if gogoproto.HasVerboseEqual(file.FileDescriptorProto, message.DescriptorProto) {
				p.P(`if err := p.VerboseEqual(msg); err != nil {`)
				p.In()
				p.P(`t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)`)
				p.Out()
				p.P(`}`)
			}
			p.P(`if !p.Equal(msg) {`)
			p.In()
			p.P(`t.Fatalf("%#v !Proto %#v", msg, p)`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
			p.P()
		}

		if gogoproto.HasTestGen(file.FileDescriptorProto, message.DescriptorProto) {
			if gogoproto.IsMarshaler(file.FileDescriptorProto, message.DescriptorProto) || gogoproto.IsUnsafeMarshaler(file.FileDescriptorProto, message.DescriptorProto) {
				p.P(`func Test`, ccTypeName, `MarshalTo(t *`, testingPkg.Use(), `.T) {`)
				p.In()
				p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
				p.P(`p := NewPopulated`, ccTypeName, `(popr, false)`)
				p.P(`size := p.Size()`)
				p.P(`data := make([]byte, size)`)
				p.P(`for i := range data {`)
				p.In()
				p.P(`data[i] = byte(popr.Intn(256))`)
				p.Out()
				p.P(`}`)
				p.P(`_, err := p.MarshalTo(data)`)
				p.P(`if err != nil {`)
				p.In()
				p.P(`panic(err)`)
				p.Out()
				p.P(`}`)
				p.P(`msg := &`, ccTypeName, `{}`)
				p.P(`if err := `, protoPkg.Use(), `.Unmarshal(data, msg); err != nil {`)
				p.In()
				p.P(`panic(err)`)
				p.Out()
				p.P(`}`)
				p.P(`for i := range data {`)
				p.In()
				p.P(`data[i] = byte(popr.Intn(256))`)
				p.Out()
				p.P(`}`)
				if gogoproto.HasVerboseEqual(file.FileDescriptorProto, message.DescriptorProto) {
					p.P(`if err := p.VerboseEqual(msg); err != nil {`)
					p.In()
					p.P(`t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)`)
					p.Out()
					p.P(`}`)
				}
				p.P(`if !p.Equal(msg) {`)
				p.In()
				p.P(`t.Fatalf("%#v !Proto %#v", msg, p)`)
				p.Out()
				p.P(`}`)
				p.Out()
				p.P(`}`)
				p.P()
			}
		}

		if gogoproto.HasBenchGen(file.FileDescriptorProto, message.DescriptorProto) {
			used = true
			p.P(`func Benchmark`, ccTypeName, `ProtoMarshal(b *`, testingPkg.Use(), `.B) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(0))`)
			p.P(`var populated []*`, ccTypeName)
			p.P(`for i := 0; i < 10900; i++ {`)
			p.In()
			p.P(`p := NewPopulated`, ccTypeName, `(popr, false)`)
			p.P(`populated = append(populated, p)`)
			p.Out()
			p.P(`}`)
			p.P(`total := 0`)
			p.P(`b.ResetTimer()`)
			p.P(`for i := 0; i < b.N; i++ {`)
			p.In()
			p.P(`p := populated[i % len(populated)]`)
			p.P(`data, err := `, protoPkg.Use(), `.Marshal(p)`)
			p.P(`if err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			p.P(`total += len(data)`)
			p.Out()
			p.P(`}`)
			p.P(`b.SetBytes(int64(total / b.N))`)
			p.Out()
			p.P(`}`)
			p.P()

			p.P(`func Benchmark`, ccTypeName, `ProtoUnmarshal(b *`, testingPkg.Use(), `.B) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(0))`)
			p.P(`var populated []*`, ccTypeName)
			p.P(`for i := 0; i < 10009; i++ {`)
			p.In()
			p.P(`p := NewPopulated`, ccTypeName, `(popr, false)`)
			p.P(`populated = append(populated, p)`)
			p.Out()
			p.P(`}`)
			p.P(`total := 0`)
			p.P(`b.ResetTimer()`)
			p.P(`for i := 0; i < b.N; i++ {`)
			p.In()
			p.P(`p := populated[i % len(populated)]`)
			p.P(`data, err := `, protoPkg.Use(), `.Marshal(p)`)
			p.P(`if err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			p.P(`msg := &`, ccTypeName, `{}`)
			p.P(`total += len(data)`)
			p.P(`b.StartTimer()`)
			p.P(`if err := `, protoPkg.Use(), `.Unmarshal(data, msg); err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
			p.P(`b.SetBytes(int64(total / b.N))`)
			p.Out()
			p.P(`}`)
			p.P()
		}
	}
	return used
}

type testJson struct {
	*generator.Generator
}

func newJson(g *generator.Generator) TestPlugin {
	return &testJson{g}
}

func (p *testJson) Generate(imports generator.PluginImports, file *generator.FileDescriptor) bool {
	used := false
	testingPkg := imports.NewImport("testing")
	randPkg := imports.NewImport("math/rand")
	timePkg := imports.NewImport("time")
	jsonPkg := imports.NewImport("encoding/json")
	for _, message := range file.Messages() {
		ccTypeName := generator.CamelCaseSlice(message.TypeName())
		if gogoproto.HasTestGen(file.FileDescriptorProto, message.DescriptorProto) {
			used = true

			p.P(`func Test`, ccTypeName, `JSON(t *`, testingPkg.Use(), `.T) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
			p.P(`p := NewPopulated`, ccTypeName, `(popr, true)`)
			p.P(`jsondata, err := `, jsonPkg.Use(), `.Marshal(p)`)
			p.P(`if err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			p.P(`msg := &`, ccTypeName, `{}`)
			p.P(`err = `, jsonPkg.Use(), `.Unmarshal(jsondata, msg)`)
			p.P(`if err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			if gogoproto.HasVerboseEqual(file.FileDescriptorProto, message.DescriptorProto) {
				p.P(`if err := p.VerboseEqual(msg); err != nil {`)
				p.In()
				p.P(`t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)`)
				p.Out()
				p.P(`}`)
			}
			p.P(`if !p.Equal(msg) {`)
			p.In()
			p.P(`t.Fatalf("%#v !Json Equal %#v", msg, p)`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)

		}

		/*if gogoproto.HasBenchGen(file.FileDescriptorProto, message.DescriptorProto) {
			used = true
			p.P(`func Benchmark`, ccTypeName, `JSONMarshal(b *`, testingPkg.Use(), `.B) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
			p.P(`total := 0`)
			p.P(`b.ResetTimer()`)
			p.P(`b.StopTimer()`)
			p.P(`for i := 0; i < b.N; i++ {`)
			p.In()
			p.P(`p := NewPopulated`, ccTypeName, `(popr, true)`)
			p.P(`b.StartTimer()`)
			p.P(`data, err := `, jsonPkg.Use(), `.Marshal(p)`)
			p.P(`if err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			p.P(`b.StopTimer()`)
			p.P(`total += len(data)`)
			p.Out()
			p.P(`}`)
			p.P(`b.SetBytes(int64(total / b.N))`)
			p.Out()
			p.P(`}`)
			p.P()

			p.P(`func Benchmark`, ccTypeName, `JSONUnmarshal(b *`, testingPkg.Use(), `.B) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
			p.P(`total := 0`)
			p.P(`b.ResetTimer()`)
			p.P(`for i := 0; i < b.N; i++ {`)
			p.In()
			p.P(`b.StopTimer()`)
			p.P(`p := NewPopulated`, ccTypeName, `(popr, true)`)
			p.P(`data, err := `, jsonPkg.Use(), `.Marshal(p)`)
			p.P(`if err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			p.P(`msg := &`, ccTypeName, `{}`)
			p.P(`total += len(data)`)
			p.P(`b.StartTimer()`)
			p.P(`if err := `, jsonPkg.Use(), `.Unmarshal(data, msg); err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
			p.P(`b.SetBytes(int64(total / b.N))`)
			p.Out()
			p.P(`}`)
			p.P()
		}*/
	}
	return used
}

type testText struct {
	*generator.Generator
}

func newText(g *generator.Generator) TestPlugin {
	return &testText{g}
}

func (p *testText) Generate(imports generator.PluginImports, file *generator.FileDescriptor) bool {
	used := false
	testingPkg := imports.NewImport("testing")
	randPkg := imports.NewImport("math/rand")
	timePkg := imports.NewImport("time")
	protoPkg := imports.NewImport("code.google.com/p/gogoprotobuf/proto")
	//fmtPkg := imports.NewImport("fmt")
	for _, message := range file.Messages() {
		ccTypeName := generator.CamelCaseSlice(message.TypeName())
		if gogoproto.HasTestGen(file.FileDescriptorProto, message.DescriptorProto) {
			used = true

			p.P(`func Test`, ccTypeName, `ProtoText(t *`, testingPkg.Use(), `.T) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
			p.P(`p := NewPopulated`, ccTypeName, `(popr, true)`)
			p.P(`data := `, protoPkg.Use(), `.MarshalTextString(p)`)
			p.P(`msg := &`, ccTypeName, `{}`)
			p.P(`if err := `, protoPkg.Use(), `.UnmarshalText(data, msg); err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			if gogoproto.HasVerboseEqual(file.FileDescriptorProto, message.DescriptorProto) {
				p.P(`if err := p.VerboseEqual(msg); err != nil {`)
				p.In()
				p.P(`t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)`)
				p.Out()
				p.P(`}`)
			}
			p.P(`if !p.Equal(msg) {`)
			p.In()
			p.P(`t.Fatalf("%#v !Proto %#v", msg, p)`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
			p.P()

			p.P(`func Test`, ccTypeName, `ProtoCompactText(t *`, testingPkg.Use(), `.T) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
			p.P(`p := NewPopulated`, ccTypeName, `(popr, true)`)
			p.P(`data := `, protoPkg.Use(), `.CompactTextString(p)`)
			p.P(`msg := &`, ccTypeName, `{}`)
			p.P(`if err := `, protoPkg.Use(), `.UnmarshalText(data, msg); err != nil {`)
			p.In()
			p.P(`panic(err)`)
			p.Out()
			p.P(`}`)
			if gogoproto.HasVerboseEqual(file.FileDescriptorProto, message.DescriptorProto) {
				p.P(`if err := p.VerboseEqual(msg); err != nil {`)
				p.In()
				p.P(`t.Fatalf("%#v !VerboseProto %#v, since %v", msg, p, err)`)
				p.Out()
				p.P(`}`)
			}
			p.P(`if !p.Equal(msg) {`)
			p.In()
			p.P(`t.Fatalf("%#v !Proto %#v", msg, p)`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
			p.P()

		}

		/*if gogoproto.HasBenchGen(file.FileDescriptorProto, message.DescriptorProto) {
			used = true

			p.P(`func Benchmark`, ccTypeName, `ProtoTextMarshal(b *`, testingPkg.Use(), `.B) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
			p.P(`total := 0`)
			p.P(`b.ResetTimer()`)
			p.P(`b.StopTimer()`)
			p.P(`for i := 0; i < b.N; i++ {`)
			p.In()
			p.P(`p := NewPopulated`, ccTypeName, `(popr, true)`)
			p.P(`b.StartTimer()`)
			p.P(`data := `, protoPkg.Use(), `.MarshalTextString(p)`)
			p.P(`b.StopTimer()`)
			p.P(`total += len(data)`)
			p.Out()
			p.P(`}`)
			p.P(`b.SetBytes(int64(total / b.N))`)
			p.Out()
			p.P(`}`)
			p.P()

			p.P(`func Benchmark`, ccTypeName, `ProtoCompactTextMarshal(b *`, testingPkg.Use(), `.B) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
			p.P(`total := 0`)
			p.P(`b.ResetTimer()`)
			p.P(`b.StopTimer()`)
			p.P(`for i := 0; i < b.N; i++ {`)
			p.In()
			p.P(`p := NewPopulated`, ccTypeName, `(popr, true)`)
			p.P(`b.StartTimer()`)
			p.P(`data := `, protoPkg.Use(), `.CompactTextString(p)`)
			p.P(`b.StopTimer()`)
			p.P(`total += len(data)`)
			p.Out()
			p.P(`}`)
			p.P(`b.SetBytes(int64(total / b.N))`)
			p.Out()
			p.P(`}`)
			p.P()

			p.P(`func Benchmark`, ccTypeName, `ProtoTextUnmarshal(b *`, testingPkg.Use(), `.B) {`)
			p.In()
			p.P(`popr := `, randPkg.Use(), `.New(`, randPkg.Use(), `.NewSource(`, timePkg.Use(), `.Now().UnixNano()))`)
			p.P(`total := 0`)
			p.P(`b.ResetTimer()`)
			p.P(`for i := 0; i < b.N; i++ {`)
			p.In()
			p.P(`b.StopTimer()`)
			p.P(`p := NewPopulated`, ccTypeName, `(popr, true)`)
			p.P(`data := `, protoPkg.Use(), `.MarshalTextString(p)`)
			p.P(`msg := &`, ccTypeName, `{}`)
			p.P(`total += len(data)`)
			p.P(`msg.Reset()`)
			p.P(`b.StartTimer()`)
			p.P(`if err := `, protoPkg.Use(), `.UnmarshalText(data, msg); err != nil {`)
			p.In()
			p.P(`panic(`, fmtPkg.Use(), `.Sprintf("%v given %v", err, data))`)
			p.Out()
			p.P(`}`)
			p.Out()
			p.P(`}`)
			p.P(`b.SetBytes(int64(total / b.N))`)
			p.Out()
			p.P(`}`)
			p.P()
		}*/
	}
	return used
}

func init() {
	RegisterTestPlugin(newProto)
	RegisterTestPlugin(newJson)
	RegisterTestPlugin(newText)
}
