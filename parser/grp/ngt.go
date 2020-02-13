// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type Ngt struct {
	Header *Ngt_Header
	Stuff uint32
	Delim []byte
	Lines []*Ngt_Line
	_io *kaitai.Stream
	_root *Ngt
	_parent interface{}
}

func (this *Ngt) Read(io *kaitai.Stream, parent interface{}, root *Ngt) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1 := new(Ngt_Header)
	err = tmp1.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Header = tmp1
	tmp2, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Stuff = uint32(tmp2)
	tmp3, err := this._io.ReadBytes(int(1))
	if err != nil {
		return err
	}
	tmp3 = tmp3
	this.Delim = tmp3
	this.Lines = make([]*Ngt_Line, (this.Header.Length - 1))
	for i := range this.Lines {
		tmp4 := new(Ngt_Line)
		err = tmp4.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Lines[i] = tmp4
	}
	return err
}
type Ngt_Header struct {
	Length uint32
	_io *kaitai.Stream
	_root *Ngt
	_parent interface{}
}

func (this *Ngt_Header) Read(io *kaitai.Stream, parent interface{}, root *Ngt) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp5, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Length = uint32(tmp5)
	return err
}
type Ngt_Line struct {
	Delim []byte
	LineHeader *Ngt_Header
	Items []*Ngt_Item
	_io *kaitai.Stream
	_root *Ngt
	_parent *Ngt
}

func (this *Ngt_Line) Read(io *kaitai.Stream, parent *Ngt, root *Ngt) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp6, err := this._io.ReadBytes(int(1))
	if err != nil {
		return err
	}
	tmp6 = tmp6
	this.Delim = tmp6
	tmp7 := new(Ngt_Header)
	err = tmp7.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LineHeader = tmp7
	this.Items = make([]*Ngt_Item, this.LineHeader.Length)
	for i := range this.Items {
		tmp8 := new(Ngt_Item)
		err = tmp8.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Items[i] = tmp8
	}
	return err
}
type Ngt_Item struct {
	Id uint32
	Distance float32
	_io *kaitai.Stream
	_root *Ngt
	_parent *Ngt_Line
}

func (this *Ngt_Item) Read(io *kaitai.Stream, parent *Ngt_Line, root *Ngt) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp9, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Id = uint32(tmp9)
	tmp10, err := this._io.ReadF4le()
	if err != nil {
		return err
	}
	this.Distance = float32(tmp10)
	return err
}
