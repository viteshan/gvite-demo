package main_test

import (
	"io/ioutil"
	"testing"

	"github.com/vitelabs/go-vite/wallet"
)

// go test -run TestWallet_RandomAddr -v
func TestWallet_RandomAddr(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "wallet")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tmpDir)
	manager := wallet.New(&wallet.Config{
		DataDir: tmpDir,
	})
	manager.Start()
	mnemonic, storeManager, err := manager.NewMnemonicAndEntropyStore("123456")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(mnemonic)
	storeManager.Unlock("123456")
	randomAddressList, err := storeManager.ListAddress(0, 1)
	if err != nil {
		t.Fatal(err)
	}

	for _, random := range randomAddressList {
		t.Log(random.Hex())
	}
}
