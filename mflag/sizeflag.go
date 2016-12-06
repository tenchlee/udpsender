package mflag

import (
	"github.com/dustin/go-humanize"
	"flag"
)

type Bytes uint64

func (b *Bytes) String() string {
	return humanize.IBytes(uint64(*b))
}

func (d *Bytes) Set(s string) error {
	v, err := humanize.ParseBytes(s)
	*d = Bytes(v)
	return err
}

func FlagBytesVar(dest *Bytes, name string, value Bytes, usage string)  {
	*dest = value
	flag.Var(dest, name, usage)
}

func FlagBytes(name string, value Bytes, usage string) *Bytes {
	d := new(Bytes)
	FlagBytesVar(d, name, value, usage)
	return d
}

