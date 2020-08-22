/*
 * Copyright 2020 National Library of Norway.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package index

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/nlnwa/gowarc/warcrecord"
)

type CdxWriter interface {
	Init() error
	Close()
	Write(wr warcrecord.WarcRecord, fileName string, offset int64, rle int) error
}

type CdxLegacy struct {
}
type CdxJ struct {
	jsonMarshaler *jsonpb.Marshaler
}
type CdxPb struct {
	jsonMarshaler *jsonpb.Marshaler
}
type CdxDb struct {
	dbDir string
	db    *Db
}

func NewCdxDb(dbDir string) *CdxDb {
	return &CdxDb{dbDir: dbDir}
}

func (c *CdxDb) Init() (err error) {
	c.db, err = NewIndexDb(c.dbDir)
	if err != nil {
		return err
	}
	return nil
}

func (c *CdxDb) Close() {
	c.db.Flush()
	c.db.Close()
}

func (c *CdxDb) Write(wr warcrecord.WarcRecord, fileName string, offset int64, rle int) error {
	return c.db.Add(wr, fileName, offset, rle)
}

func (c *CdxLegacy) Init() (err error) {
	return nil
}

func (c *CdxLegacy) Close() {
}

func (c *CdxLegacy) Write(wr warcrecord.WarcRecord, fileName string, offset int64, rle int) error {
	return nil
}

func (c *CdxJ) Init() (err error) {
	c.jsonMarshaler = &jsonpb.Marshaler{}
	return nil
}

func (c *CdxJ) Close() {
}

func (c *CdxJ) Write(wr warcrecord.WarcRecord, fileName string, offset int64, rle int) error {
	if wr.Type() == warcrecord.RESPONSE {
		rec := NewCdxRecord(wr, fileName, offset, rle)
		cdxj, err := c.jsonMarshaler.MarshalToString(rec)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s %s %s\n", rec.Ssu, rec.Sts, rec.Srt, cdxj)
	}
	return nil
}

func (c *CdxPb) Init() (err error) {
	return nil
}

func (c *CdxPb) Close() {
}

func (c *CdxPb) Write(wr warcrecord.WarcRecord, fileName string, offset int64, rle int) error {
	if wr.Type() == warcrecord.RESPONSE {
		rec := NewCdxRecord(wr, fileName, offset, rle)
		cdxpb, err := proto.Marshal(rec)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s %s %s\n", rec.Ssu, rec.Sts, rec.Srt, cdxpb)
	}
	return nil
}
