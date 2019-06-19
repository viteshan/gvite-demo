package reward

import (
	"fmt"
	"math/big"
	"math/rand"
	"os/user"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vitelabs/go-vite/rpcapi/api"

	"github.com/vitelabs/go-vite/client"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/wallet"
)

func init() {
	current, _ := user.Current()
	home := current.HomeDir
	WalletDir = path.Join(home, ".gvite/wallet")
}

var WalletDir string

var RawUrl = ""

var selfAddr = types.HexToAddressPanic("")
var sbpFixAddr = types.HexToAddressPanic("")
var passwd = ""
var mnemonic = ""

func TestCalReward(t *testing.T) {
	cli, err := client.NewRpcClient(RawUrl)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	startIdx := uint64(0)
	endIdx := uint64(26)
	//todo reward address
	votes, err := CalReward(startIdx, endIdx, "", sbpFixAddr, cli)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(votes)

	rewardTotal, err := MergeReward(votes)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("total reward", big.NewInt(0).Div(rewardTotal, OnePercentVite).String())

	details, err := CalRewardDropDetails(votes)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(details)

	{ // print to csv
		var logs = ""
		for i := startIdx; i <= endIdx; i++ {
			details := details[i]
			for kk, vv := range details {
				amount := big.NewInt(0).Set(vv.Amount)
				amountStr := amount.String()
				logs += fmt.Sprintf("%d,%s,%f,%s\n", i, kk, float64(amount.Div(vv.Amount, OnePercentVite).Int64())/100.0, amountStr)
			}
		}
		fmt.Println(logs)
	}

	mergedDetails, err := MergeRewardDrop(details)

	sendTotal := big.NewInt(0)
	for k, v := range mergedDetails {
		amount := big.NewInt(0).Set(v.Amount)
		t.Log(k, amount.Div(amount, OnePercentVite))
		sendTotal.Add(sendTotal, amount)
	}

	//t.Log("all send", sendTotal.Div(sendTotal, oneVite))
	t.Log("all send", sendTotal)

	var finalLog []*SBPRewardDropDetails

	for _, v := range mergedDetails {
		amount := big.NewInt(0).Set(v.Amount)
		if amount.Cmp(OnePercentVite) < 0 {
			t.Log("ignore amount", v.ToAddr, v.Amount)
			continue
		}
		finalLog = append(finalLog, &SBPRewardDropDetails{ToAddr: v.ToAddr, Amount: amount})
	}

	for k, v := range finalLog {
		t.Log(k, v.ToAddr, v.Amount.String())
	}

	return

	txId := fmt.Sprintf("all send: %d", rand.Intn(10000000))

	w := wallet.New(&wallet.Config{
		DataDir:        WalletDir,
		MaxSearchIndex: 100000,
	})
	w.Start()

	em, err := w.RecoverEntropyStoreFromMnemonic(mnemonic, passwd)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	em.Unlock(passwd)

	all, err := TransferRewardDrop(selfAddr, finalLog, txId, cli, func(addr types.Address, data []byte) (signedData, pubkey []byte, err error) {
		time.Sleep(time.Second)
		return em.SignData(addr, data)
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, v := range all {
		t.Log(selfAddr, v.Height, v.Hash)
	}
}

func TestTransferRewardDrop(t *testing.T) {
	// "address, amount"
	logs := []string{}

	var finalLog []*SBPRewardDropDetails
	for _, v := range logs {
		tmps := strings.Split(v, ",")

		addr, err := types.HexToAddress(tmps[0])
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		amount, flag := big.NewInt(0).SetString(tmps[1], 10)
		if !flag {
			t.Error("string to bigInt fail.", tmps[1])
			t.FailNow()
		}
		if amount.Cmp(OnePercentVite) < 0 {
			t.Log("ignore amount", addr, amount)
			continue
		}
		finalLog = append(finalLog, &SBPRewardDropDetails{ToAddr: addr, Amount: amount})
	}

	sendTotal := big.NewInt(0)
	for _, v := range finalLog {
		sendTotal.Add(sendTotal, v.Amount)
	}

	for k, v := range finalLog {
		t.Log(k, v.ToAddr, v.Amount.String())
	}
	t.Log("all send", sendTotal.Div(sendTotal, OnePercentVite))

	w := wallet.New(&wallet.Config{
		DataDir:        WalletDir,
		MaxSearchIndex: 100000,
	})
	w.Start()

	em, err := w.RecoverEntropyStoreFromMnemonic(mnemonic, passwd)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	em.Unlock(passwd)

	cli, err := client.NewRpcClient(RawUrl)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	txId := fmt.Sprintf("%d", rand.Intn(10000000))
	t.Log("txId", txId)
	all, err := TransferRewardDrop(selfAddr, finalLog, txId, cli, func(addr types.Address, data []byte) (signedData, pubkey []byte, err error) {
		time.Sleep(time.Millisecond * 1000)
		return em.SignData(addr, data)
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for k, v := range all {
		t.Log(selfAddr, k, v.Height, v.Hash)
	}
}

func TestCalRewardTestNet(t *testing.T) {
	startIdx := uint64(172)
	endIdx := uint64(194)
	votes, err := calRewardForTestNet(startIdx, endIdx, "", sbpFixAddr, rewards_172_194, votesDetails_172_194)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	rewardTotal, err := MergeReward(votes)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("total reward", big.NewInt(0).Div(rewardTotal, OnePercentVite).String())

	details, err := CalRewardDropDetails(votes)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(details)

	{ // print to csv
		var logs = ""
		for i := startIdx; i <= endIdx; i++ {
			details := details[i]
			for kk, vv := range details {
				amount := big.NewInt(0).Set(vv.Amount)
				amountStr := amount.String()
				logs += fmt.Sprintf("%d,%s,%f,%s\n", i, kk, float64(amount.Div(vv.Amount, OnePercentVite).Int64())/100.0, amountStr)
			}
		}
		fmt.Println(logs)
	}

	mergedDetails, err := MergeRewardDrop(details)

	sendTotal := big.NewInt(0)
	for k, v := range mergedDetails {
		amount := big.NewInt(0).Set(v.Amount)
		t.Log(k, amount.Div(amount, OnePercentVite))
		sendTotal.Add(sendTotal, amount)
	}
	t.Log(sendTotal)

	var finalLog []*SBPRewardDropDetails

	for _, v := range mergedDetails {
		amount := big.NewInt(0).Set(v.Amount)
		if amount.Cmp(OnePercentVite) < 0 {
			t.Log("ignore amount", v.ToAddr, v.Amount)
			continue
		}
		finalLog = append(finalLog, &SBPRewardDropDetails{ToAddr: v.ToAddr, Amount: amount})
	}

	for k, v := range finalLog {
		t.Log(k, v.ToAddr, v.Amount.String())
	}

	//return

	cli, err := client.NewRpcClient(RawUrl)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	txId := fmt.Sprintf("%d", rand.Intn(10000000))

	w := wallet.New(&wallet.Config{
		DataDir:        WalletDir,
		MaxSearchIndex: 100000,
	})
	w.Start()

	em, err := w.RecoverEntropyStoreFromMnemonic(mnemonic, passwd)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	em.Unlock(passwd)

	all, err := TransferRewardDrop(selfAddr, finalLog, txId, cli, func(addr types.Address, data []byte) (signedData, pubkey []byte, err error) {
		time.Sleep(time.Second)
		return em.SignData(addr, data)
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, v := range all {
		t.Log(selfAddr, v.Height, v.Hash)
	}
}

func calRewardForTestNet(startIdx uint64, endIdx uint64, sbp string, sbpRewardAddr types.Address, rewards []string, votesDetails []string) (map[uint64]*SBPRewardVote, error) {

	rewardMap := make(map[uint64]*api.Reward)

	for _, v := range rewards {
		tmps := strings.Split(v, ",")
		if len(tmps) != 2 {
			panic(v)
		}
		t0, err := strconv.ParseUint(tmps[0], 10, 64)
		if err != nil {
			panic(err)
		}
		f1, err := strconv.ParseFloat(tmps[1], 64)
		if err != nil {
			panic(err)
		}

		rewardMap[t0] = &api.Reward{TotalReward: floatOneViteToBigInt(f1).String()}
	}
	votesDetailMap := make(map[uint64]map[types.Address]*big.Int)

	for _, v := range votesDetails {
		tmps := strings.Split(v, ",")
		if len(tmps) != 3 {
			panic(v)
		}

		t0, err := strconv.ParseUint(tmps[0], 10, 64)
		if err != nil {
			panic(err)
		}
		t1 := types.HexToAddressPanic(tmps[1])

		f2, err := strconv.ParseFloat(tmps[2], 64)
		if err != nil {
			panic(err)
		}
		detail, ok := votesDetailMap[t0]
		if !ok {
			detail = make(map[types.Address]*big.Int)
			votesDetailMap[t0] = detail
		}
		_, ok = detail[t1]
		if ok {
			panic("t1 dup")
		}
		detail[t1] = floatOneViteToBigInt(f2)
	}

	result := make(map[uint64]*SBPRewardVote)
	for i := startIdx; i <= endIdx; i++ {
		detailTmp, ok := votesDetailMap[i]
		if !ok {
			panic(fmt.Sprintf("%d not exist", i))
		}
		rewardTmp, ok := rewardMap[i]
		if !ok {
			panic(fmt.Sprintf("%d not exist", i))
		}

		tmp := &SBPRewardVote{}
		result[i] = tmp

		tmp.SbpReward = rewardTmp

		totalBalance := big.NewInt(0)
		for _, v := range detailTmp {
			totalBalance.Add(totalBalance, v)
		}
		tmp.VoteTotal, tmp.VoteDetails = reCalRewardForTestNet(sbpRewardAddr, totalBalance, detailTmp)
	}
	return result, nil
}

func reCalRewardForTestNet(sbpRewardAddr types.Address, totalBalance *big.Int, voteDetails map[types.Address]*big.Int) (*big.Int, map[types.Address]*big.Int) {
	sbpFixVote := big.NewInt(500000)
	sbpFixVote.Mul(sbpFixVote, big.NewInt(1152))
	sbpFixVote.Mul(sbpFixVote, OneVite)

	reTotalBalance := big.NewInt(0)
	reTotalBalance.Add(reTotalBalance, totalBalance)
	reTotalBalance.Add(reTotalBalance, sbpFixVote)

	reSbp, ok := voteDetails[sbpRewardAddr]
	if !ok {
		voteDetails[sbpRewardAddr] = sbpFixVote
	} else {
		reSbp.Add(reSbp, sbpFixVote)
	}
	return reTotalBalance, voteDetails
}

func floatOneViteToBigInt(f float64) *big.Int {
	bf := big.NewFloat(f)
	bf.Mul(bf, OneViteFloat)
	if bf.IsInt() {
		result, _ := bf.Int(nil)
		return result
	} else {
		panic("not bigInt")
	}
}

func TestWallet(t *testing.T) {
	w := wallet.New(&wallet.Config{
		DataDir:        WalletDir,
		MaxSearchIndex: 100000,
	})
	w.Start()

	em, err := w.RecoverEntropyStoreFromMnemonic(mnemonic, passwd)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	em.Unlock(passwd)

	t.Log(em.GetPrimaryAddr().Hex())
}
