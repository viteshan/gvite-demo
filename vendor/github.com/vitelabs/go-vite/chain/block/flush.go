package chain_block

import (
	"bytes"
	"fmt"
	"github.com/vitelabs/go-vite/chain/utils"
	"github.com/vitelabs/go-vite/common/types"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type BufWriter struct {
	Buffer *bytes.Buffer
	Err    error
}

func NewBufWriter() *BufWriter {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	return &BufWriter{
		Buffer: buf,
	}
}

func (w *BufWriter) Write(data []byte) error {
	w.Buffer.Write(data)
	return nil
}
func (w *BufWriter) WriteError(err error) {
	w.Err = err
}
func (w *BufWriter) Close() error {
	return nil
}

func (w *BufWriter) Release() {
	bufPool.Put(w.Buffer)
}

func (bDB *BlockDB) Id() types.Hash {
	return bDB.id
}

// lock write
func (bDB *BlockDB) Prepare() {
	// set bDB.flushTargetLocation
	bDB.flushStartLocation = bDB.fm.NextFlushStartLocation()
	// set bDB.flushTargetLocation
	bDB.flushTargetLocation = bDB.fm.LatestLocation()

	// set bDB.flushBuf
	bufWriter := NewBufWriter()

	bDB.fm.ReadRange(bDB.flushStartLocation, bDB.flushTargetLocation, bufWriter)

	if bufWriter.Err != nil {
		panic(fmt.Sprintf("BlockDB prepare failed when flush, start location is %+v, target location is %+v. Error: %s",
			bDB.flushStartLocation, bDB.flushTargetLocation, bufWriter.Err))
	}

	bDB.flushBuf = bufWriter

	// set next flush start location
	bDB.fm.SetNextFlushStartLocation(bDB.flushTargetLocation)
}

// Assume locking writing (lock the insertion and deletion of account block and snapshot block)
// nextFlushStartLocation is the start location of the next flush.
// bDB.flushStartLocation is
func (bDB *BlockDB) CancelPrepare() {
	nextFlushStartLocation := bDB.fm.NextFlushStartLocation()

	if nextFlushStartLocation.Compare(bDB.flushStartLocation) > 0 {
		bDB.fm.SetNextFlushStartLocation(bDB.flushStartLocation)
	}

	bDB.flushStartLocation = nil
	bDB.flushTargetLocation = nil

	bDB.flushBuf.Release()
	bDB.flushBuf = nil
}

func (bDB *BlockDB) RedoLog() ([]byte, error) {
	var redoLog []byte
	data := bDB.flushBuf.Buffer.Bytes()
	redoLog = make([]byte, 0, 24+len(data))

	redoLog = append(redoLog, chain_utils.SerializeLocation(bDB.flushStartLocation)...)
	redoLog = append(redoLog, chain_utils.SerializeLocation(bDB.flushTargetLocation)...)
	redoLog = append(redoLog, data...)

	return redoLog, nil

}

func (bDB *BlockDB) Commit() error {
	return bDB.fm.Flush(bDB.flushStartLocation, bDB.flushTargetLocation, bDB.flushBuf.Buffer.Bytes())
}

// lock write
func (bDB *BlockDB) AfterCommit() {
	bDB.flushStartLocation = nil
	bDB.flushTargetLocation = nil

	bDB.flushBuf.Release()
	bDB.flushBuf = nil

}

func (bDB *BlockDB) BeforeRecover(redoLog []byte) {
	flushStartLocation := chain_utils.DeserializeLocation(redoLog[:12])

	if err := bDB.fm.DeleteTo(flushStartLocation); err != nil {
		panic(err)
	}

	if _, err := bDB.fm.Write(redoLog[24:]); err != nil {
		panic(err)
	}
}

func (bDB *BlockDB) AfterRecover() {}

func (bDB *BlockDB) PatchRedoLog(redoLog []byte) error {
	flushStartLocation := chain_utils.DeserializeLocation(redoLog[:12])
	flushTargetLocation := chain_utils.DeserializeLocation(redoLog[12:24])

	return bDB.fm.Flush(flushStartLocation, flushTargetLocation, redoLog[24:])
}
