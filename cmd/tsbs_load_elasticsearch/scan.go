package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/timescale/tsbs/pkg/data"
	"github.com/timescale/tsbs/pkg/data/usecases/common"
	"github.com/timescale/tsbs/pkg/targets"
)

var newLine = []byte("\n")

type fileDataSource struct {
	scanner *bufio.Scanner
}

func (d *fileDataSource) NextItem() data.LoadedPoint {
	ok := d.scanner.Scan()
	if !ok && d.scanner.Err() == nil { // nothing scanned & no error = EOF
		return data.LoadedPoint{}
	} else if !ok {
		fatal("scan error: %v", d.scanner.Err())
		return data.LoadedPoint{}
	}
	return data.NewLoadedPoint(d.scanner.Bytes())
}

func (d *fileDataSource) Headers() *common.GeneratedDataHeaders { return nil }

type batch struct {
	buf     *bytes.Buffer
	rows    uint
	metrics uint64
}

func (b *batch) Len() uint {
	return b.rows
}

func (b *batch) Append(item data.LoadedPoint) {
	that := item.Data.([]byte)
	b.rows++
	//b.metrics += uint64(len(strings.Split(args[1], ",")))
	b.buf.Write([]byte(fmt.Sprintf(`{ "index" : { "_index" : "%s"} }`, loader.DatabaseName())))
	b.buf.Write(newLine)
	b.buf.Write(that)
	b.buf.Write(newLine)
}

type factory struct{}

func (f *factory) New() targets.Batch {
	return &batch{buf: bufPool.Get().(*bytes.Buffer)}
}
