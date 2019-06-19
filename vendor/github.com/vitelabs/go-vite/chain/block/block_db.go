package chain_block

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/snappy"
	"github.com/pkg/errors"
	"github.com/vitelabs/go-vite/chain/file_manager"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/crypto"
	"github.com/vitelabs/go-vite/interfaces"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/log15"
	"io"
	"path"
	"sync"
)

type BlockDB struct {
	fm *chain_file_manager.FileManager

	snappyWriteBuffer []byte
	wg                sync.WaitGroup

	fileSize int64
	id       types.Hash

	flushStartLocation  *chain_file_manager.Location
	flushTargetLocation *chain_file_manager.Location
	flushBuf            *BufWriter

	log log15.Logger
}

func NewBlockDB(chainDir string) (*BlockDB, error) {
	id, _ := types.BytesToHash(crypto.Hash256([]byte("blockDb")))

	fileSize := int64(10 * 1024 * 1024) // 10M
	fm, err := chain_file_manager.NewFileManager(path.Join(chainDir, "blocks"), fileSize, 10)
	if err != nil {
		return nil, err
	}

	return &BlockDB{
		fm:                fm,
		fileSize:          fileSize,
		snappyWriteBuffer: make([]byte, fileSize),
		id:                id,
		log:               log15.New("module", "blockDB"),
	}, nil
}

func (bDB *BlockDB) FileSize() int64 {
	return bDB.fileSize
}

func (bDB *BlockDB) Close() error {
	if err := bDB.fm.Close(); err != nil {
		return errors.New(fmt.Sprintf("bDB.fm.Close failed, error is %s", err))
	}

	bDB.fm = nil
	return nil
}

func (bDB *BlockDB) Write(ss *ledger.SnapshotChunk) ([]*chain_file_manager.Location, *chain_file_manager.Location, error) {

	accountBlocksLocation := make([]*chain_file_manager.Location, 0, len(ss.AccountBlocks))

	for _, accountBlock := range ss.AccountBlocks {
		buf, err := accountBlock.Serialize()
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("ss.AccountBlocks.Serialize failed, error is %s, accountBlock is %+v", err.Error(), accountBlock))
		}

		if location, err := bDB.fm.Write(makeWriteBytes(bDB.snappyWriteBuffer, BlockTypeAccountBlock, buf)); err != nil {
			return nil, nil, errors.New(fmt.Sprintf("bDB.fm.Write failed, error is %s, accountBlock is %+v", err.Error(), accountBlock))
		} else {
			accountBlocksLocation = append(accountBlocksLocation, location)
		}
	}

	buf, err := ss.SnapshotBlock.Serialize()
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("ss.SnapshotBlock.Serialize failed, error is %s, snapshotBlock is %+v", err.Error(), ss.SnapshotBlock))
	}

	snapshotBlockLocation, err := bDB.fm.Write(makeWriteBytes(bDB.snappyWriteBuffer, BlockTypeSnapshotBlock, buf))

	//bDB.log.Info(fmt.Sprintf("sb %s %d %d", ss.SnapshotBlock.Hash, ss.SnapshotBlock.Height, data), "method", "Write")

	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("bDB.fm.Write failed, error is %s, snapshotBlock is %+v", err.Error(), ss.SnapshotBlock))
	}
	return accountBlocksLocation, snapshotBlockLocation, nil
}

func (bDB *BlockDB) Read(location *chain_file_manager.Location) ([]byte, error) {
	buf, _, err := bDB.fm.Read(location)
	if err != nil {
		return nil, err
	}
	if len(buf) <= 0 {
		return nil, nil
	}

	sBuf, err := snappy.Decode(nil, buf[1:])
	if err != nil {
		return nil, err
	}
	return sBuf, nil
}

func (bDB *BlockDB) ReadRaw(startLocation *chain_file_manager.Location, buf []byte) (*chain_file_manager.Location, int, error) {
	return bDB.fm.ReadRaw(startLocation, buf)
}

func (bDB *BlockDB) ReadRange(startLocation *chain_file_manager.Location, endLocation *chain_file_manager.Location) ([]*ledger.SnapshotChunk, error) {
	bfp := newBlockFileParser()

	endLocation = bDB.maxLocation(endLocation)

	bDB.wg.Add(1)
	go func() {
		defer bDB.wg.Done()
		bDB.fm.ReadRange(startLocation, endLocation, bfp)
		if endLocation != nil {
			buf, _, err := bDB.fm.Read(endLocation)

			if len(buf) >= 0 {
				bufSizeBytes := make([]byte, 4)
				binary.BigEndian.PutUint32(bufSizeBytes, uint32(len(buf)))
				bfp.Write(bufSizeBytes)
				bfp.Write(buf)
			}

			if err != nil && err != io.EOF {
				bfp.WriteError(err)
				return
			}

		}
		bfp.Close()
	}()

	var segList []*ledger.SnapshotChunk
	var seg *ledger.SnapshotChunk

	var snappyReadBuffer = make([]byte, 0, 8*1024) // 8kb
	iterator := bfp.Iterator()

	for buf := range iterator {
		if seg == nil {
			seg = &ledger.SnapshotChunk{}
		}

		sBuf, err := snappy.Decode(snappyReadBuffer, buf.Buffer)
		if err != nil {
			return nil, err
		}

		if buf.BlockType == BlockTypeSnapshotBlock {

			sb := &ledger.SnapshotBlock{}
			if err := sb.Deserialize(sBuf); err != nil {
				return nil, err
			}
			seg.SnapshotBlock = sb
			segList = append(segList, seg)
			seg = nil
		} else if buf.BlockType == BlockTypeAccountBlock {
			ab := &ledger.AccountBlock{}
			if err := ab.Deserialize(sBuf); err != nil {
				return nil, err
			}
			seg.AccountBlocks = append(seg.AccountBlocks, ab)
		}
	}

	if err := bfp.Error(); err != nil {
		return nil, err
	}

	if seg != nil {
		segList = append(segList, seg)
	}

	return segList, nil
}

func (bDB *BlockDB) GetNextLocation(location *chain_file_manager.Location) (*chain_file_manager.Location, error) {
	nextLocation, err := bDB.fm.GetNextLocation(location)
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
		return nil, nil
	}
	return nextLocation, nil
}

func (bDB *BlockDB) PrepareRollback(location *chain_file_manager.Location) ([]*ledger.SnapshotChunk, error) {
	bfp := newBlockFileParser()

	bDB.wg.Add(1)
	go func() {
		defer bDB.wg.Done()
		bDB.fm.ReadRange(location, bDB.fm.LatestLocation(), bfp)
		bfp.Close()
	}()

	var segList []*ledger.SnapshotChunk
	var seg *ledger.SnapshotChunk
	var snappyReadBuffer = make([]byte, 0, 4*1024) // 4KB

	iterator := bfp.Iterator()

	for buf := range iterator {
		if seg == nil {
			seg = &ledger.SnapshotChunk{}
		}

		sBuf, err := snappy.Decode(snappyReadBuffer, buf.Buffer)
		if err != nil {
			return nil, err
		}

		if buf.BlockType == BlockTypeSnapshotBlock {

			sb := &ledger.SnapshotBlock{}
			if err := sb.Deserialize(sBuf); err != nil {
				return nil, err
			}
			seg.SnapshotBlock = sb
			segList = append(segList, seg)
			seg = nil
		} else if buf.BlockType == BlockTypeAccountBlock {

			ab := &ledger.AccountBlock{}
			if err := ab.Deserialize(sBuf); err != nil {
				return nil, err
			}
			seg.AccountBlocks = append(seg.AccountBlocks, ab)
		}
	}

	if seg != nil {
		segList = append(segList, seg)
	}

	if err := bfp.Error(); err != nil {
		return nil, err
	}

	return segList, nil

}

func (bDB *BlockDB) Rollback(location *chain_file_manager.Location) error {
	return bDB.fm.DeleteTo(location)
}

func (bDB *BlockDB) SetLog(h log15.Handler) {
	bDB.log.SetHandler(h)
	bDB.fm.SetLog(h)
}

func (bDB *BlockDB) GetStatus() []interfaces.DBStatus {
	return bDB.fm.GetCacheStatusList()
}

func (bDB *BlockDB) maxLocation(location *chain_file_manager.Location) *chain_file_manager.Location {
	latestLocation := bDB.fm.LatestLocation()

	if location == nil || (latestLocation != nil && location.Compare(latestLocation) > 0) {
		return latestLocation
	}
	return location
}
