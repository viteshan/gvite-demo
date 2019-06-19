package chain

import (
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/vm_db"
	"sync"
)

type Listener interface{}

const (
	prepareInsertAbsEvent = byte(1)
	insertAbsEvent        = byte(2)

	prepareInsertSbsEvent = byte(3)
	InsertSbsEvent        = byte(4)

	prepareDeleteAbsEvent = byte(5)
	DeleteAbsEvent        = byte(6)

	prepareDeleteSbsEvent = byte(7)
	deleteSbsEvent        = byte(8)
)

type eventManager struct {
	listenerList []EventListener

	maxHandlerId uint32
	mu           sync.Mutex
}

func newEventManager() *eventManager {
	return &eventManager{
		maxHandlerId: 0,
		listenerList: make([]EventListener, 0),
	}
}

func (em *eventManager) TriggerInsertAbs(eventType byte, vmAccountBlocks []*vm_db.VmAccountBlock) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	if len(em.listenerList) <= 0 {
		return nil
	}

	switch eventType {
	case prepareInsertAbsEvent:
		for _, listener := range em.listenerList {
			if err := listener.PrepareInsertAccountBlocks(vmAccountBlocks); err != nil {
				return err
			}
		}
	case insertAbsEvent:
		for _, listener := range em.listenerList {
			listener.InsertAccountBlocks(vmAccountBlocks)
		}

	}
	return nil
}

func (em *eventManager) TriggerDeleteAbs(eventType byte, accountBlocks []*ledger.AccountBlock) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	if len(em.listenerList) <= 0 {
		return nil
	}

	switch eventType {
	case prepareDeleteAbsEvent:
		for _, listener := range em.listenerList {
			if err := listener.PrepareDeleteAccountBlocks(accountBlocks); err != nil {
				return err
			}
		}
	case DeleteAbsEvent:
		for _, listener := range em.listenerList {
			listener.DeleteAccountBlocks(accountBlocks)
		}
	}
	return nil
}

func (em *eventManager) TriggerInsertSbs(eventType byte, chunks []*ledger.SnapshotChunk) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	if len(em.listenerList) <= 0 {
		return nil
	}

	switch eventType {

	case prepareInsertSbsEvent:
		for _, listener := range em.listenerList {
			if err := listener.PrepareInsertSnapshotBlocks(chunks); err != nil {
				return err
			}
		}
	case InsertSbsEvent:
		for _, listener := range em.listenerList {
			listener.InsertSnapshotBlocks(chunks)
		}
	}
	return nil
}

func (em *eventManager) TriggerDeleteSbs(eventType byte, chunks []*ledger.SnapshotChunk) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	if len(em.listenerList) <= 0 {
		return nil
	}
	switch eventType {
	case prepareDeleteSbsEvent:
		for _, listener := range em.listenerList {
			if err := listener.PrepareDeleteSnapshotBlocks(chunks); err != nil {
				return err
			}
		}
	case deleteSbsEvent:
		for _, listener := range em.listenerList {
			listener.DeleteSnapshotBlocks(chunks)
		}
	}
	return nil
}

func (em *eventManager) Register(listener EventListener) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.listenerList = append(em.listenerList, listener)
}
func (em *eventManager) UnRegister(listener EventListener) {
	em.mu.Lock()
	defer em.mu.Unlock()

	for index, listener := range em.listenerList {
		if listener == listener {
			em.listenerList = append(em.listenerList[:index], em.listenerList[index+1:]...)
			break
		}
	}
}

func (c *chain) Register(listener EventListener) {

	c.em.Register(listener)

}

func (c *chain) UnRegister(listener EventListener) {
	c.em.UnRegister(listener)
}
