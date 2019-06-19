package chain_plugins

import (
	"errors"
	"fmt"
	"github.com/vitelabs/go-vite/chain/db"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/vm_db"
	"os"
	"path"
	"sync"
	"sync/atomic"
)

const roundSize = uint64(10)

const (
	stop  = 0
	start = 1
)

type Plugins struct {
	dataDir string

	log     log15.Logger
	chain   Chain
	store   *chain_db.Store
	plugins map[string]Plugin

	writeStatus uint32
	mu          sync.RWMutex
}

func NewPlugins(chainDir string, chain Chain) (*Plugins, error) {
	var err error

	dataDir := path.Join(chainDir, "plugins")

	store, err := chain_db.NewStore(dataDir, "plugins")
	if err != nil {
		return nil, err
	}

	plugins := map[string]Plugin{
		"filterToken": newFilterToken(store, chain),
		"onRoadInfo":  newOnRoadInfo(store, chain),
	}

	return &Plugins{
		dataDir:     dataDir,
		chain:       chain,
		store:       store,
		plugins:     plugins,
		writeStatus: start,
		log:         log15.New("module", "chain_plugins"),
	}, nil
}

func (p *Plugins) StopWrite() {
	if !atomic.CompareAndSwapUint32(&p.writeStatus, start, stop) {
		return
	}

	p.mu.Lock()
}

func (p *Plugins) StartWrite() {
	if !atomic.CompareAndSwapUint32(&p.writeStatus, stop, start) {
		return
	}

	p.mu.Unlock()
}

func (p *Plugins) RebuildData() error {
	p.StopWrite()
	defer p.StartWrite()

	p.log.Info("Start rebuild plugin data")

	if err := p.store.Close(); err != nil {
		return err
	}

	// remove data
	os.RemoveAll(p.dataDir)

	// set new store
	store, err := chain_db.NewStore(p.dataDir, "plugins")
	if err != nil {
		return err
	}

	p.store = store

	for _, plugin := range p.plugins {
		plugin.SetStore(store)
	}

	// replace flusher store
	flusher := p.chain.Flusher()
	flusher.ReplaceStore(p.store.Id(), store)

	// get latest snapshot block
	latestSnapshot := p.chain.GetLatestSnapshotBlock()
	if latestSnapshot == nil {
		return errors.New("GetLatestSnapshotBlock fail")
	}

	p.log.Info(fmt.Sprintf("latestSnapshot[%v %v]", latestSnapshot.Hash, latestSnapshot.Height), "method", "RebuildData")

	// build data
	h := uint64(0)

	for h < latestSnapshot.Height {
		targetH := h + roundSize
		if targetH > latestSnapshot.Height {
			targetH = latestSnapshot.Height
		}

		chunks, err := p.chain.GetSubLedger(h, targetH)
		if err != nil {
			return err
		}

		p.log.Info(fmt.Sprintf("rebuild %d - %d", h+1, targetH), "method", "RebuildData")

		for _, chunk := range chunks {

			if chunk.SnapshotBlock != nil &&
				chunk.SnapshotBlock.Height == h {
				continue
			}
			// write ab
			for _, ab := range chunk.AccountBlocks {

				batch := p.store.NewBatch()

				for _, plugin := range p.plugins {
					if err := plugin.InsertAccountBlock(batch, ab); err != nil {
						return err
					}
				}
				p.store.WriteAccountBlock(batch, ab)
			}

			// write sb
			batch := p.store.NewBatch()

			for _, plugin := range p.plugins {
				if err := plugin.InsertSnapshotBlock(batch, chunk.SnapshotBlock, chunk.AccountBlocks); err != nil {
					pErr := errors.New(fmt.Sprintf("InsertSnapshotBlock fail, err:%v, sb[%v, %v,len=%v] ", err, chunk.SnapshotBlock.Height, chunk.SnapshotBlock.Hash, len(chunk.AccountBlocks)))
					p.log.Error(pErr.Error(), "method", "RebuildData")
					return pErr
				}
			}

			p.store.WriteSnapshot(batch, chunk.AccountBlocks)

		}

		// flush to disk
		flusher.Flush()

		h = targetH
	}

	// success
	p.log.Info("Succeed rebuild plugin data")
	return nil
}

func (p *Plugins) Close() error {
	if err := p.store.Close(); err != nil {
		return err
	}
	return nil
}

func (p *Plugins) Store() *chain_db.Store {
	return p.store
}

func (p *Plugins) GetPlugin(name string) Plugin {
	return p.plugins[name]
}

func (p *Plugins) RemovePlugin(name string) {
	delete(p.plugins, name)
}

func (p *Plugins) PrepareInsertAccountBlocks(vmBlocks []*vm_db.VmAccountBlock) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// for recover
	for _, vmBlock := range vmBlocks {
		batch := p.store.NewBatch()

		for _, plugin := range p.plugins {
			if err := plugin.InsertAccountBlock(batch, vmBlock.AccountBlock); err != nil {
				return err
			}
		}
		p.store.WriteAccountBlock(batch, vmBlock.AccountBlock)
	}

	return nil
}

func (p *Plugins) PrepareInsertSnapshotBlocks(chunks []*ledger.SnapshotChunk) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, chunk := range chunks {
		batch := p.store.NewBatch()

		for _, plugin := range p.plugins {

			if err := plugin.InsertSnapshotBlock(batch, chunk.SnapshotBlock, chunk.AccountBlocks); err != nil {
				return err
			}
		}
		p.store.WriteSnapshot(batch, chunk.AccountBlocks)

	}

	return nil
}

func (p *Plugins) PrepareDeleteAccountBlocks(blocks []*ledger.AccountBlock) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	batch := p.store.NewBatch()

	for _, plugin := range p.plugins {
		if err := plugin.DeleteAccountBlocks(batch, blocks); err != nil {
			return err
		}
	}
	p.store.RollbackAccountBlocks(batch, blocks)

	return nil
}

func (p *Plugins) PrepareDeleteSnapshotBlocks(chunks []*ledger.SnapshotChunk) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	batch := p.store.NewBatch()

	for _, plugin := range p.plugins {

		if err := plugin.DeleteSnapshotBlocks(batch, chunks); err != nil {
			return err
		}

	}
	p.store.RollbackSnapshot(batch)

	return nil
}

func (p *Plugins) DeleteSnapshotBlocks(chunks []*ledger.SnapshotChunk) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	allUnconfirmedBlocks := p.chain.GetAllUnconfirmedBlocks()

	rollbackBatch := p.store.NewBatch()

	for _, plugin := range p.plugins {
		if err := plugin.RemoveNewUnconfirmed(rollbackBatch, allUnconfirmedBlocks); err != nil {
			return err
		}
	}

	p.store.RollbackSnapshot(rollbackBatch)

	for _, plugin := range p.plugins {
		batch := p.store.NewBatch()
		for _, unconfirmedBlock := range allUnconfirmedBlocks {
			if err := plugin.InsertAccountBlock(rollbackBatch, unconfirmedBlock); err != nil {
				return err
			}
			p.store.WriteAccountBlock(batch, unconfirmedBlock)
		}

	}
	return nil
}

func (p *Plugins) InsertAccountBlocks(blocks []*vm_db.VmAccountBlock) error {
	return nil
}
func (p *Plugins) InsertSnapshotBlocks(chunks []*ledger.SnapshotChunk) error {
	return nil
}
func (p *Plugins) DeleteAccountBlocks(blocks []*ledger.AccountBlock) error {
	return nil
}

func (p *Plugins) checkAndRecover() (*chain_db.Store, error) {
	return nil, nil
}
