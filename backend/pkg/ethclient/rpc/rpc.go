package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethclient2 "github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/ethclient"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type rpcTransaction struct {
	TxExtraInfo []txExtraInfo `json:"transactions"`
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

type rpcBlock struct {
	Hash         common.Hash          `json:"hash"`
	Transactions []*types.Transaction `json:"transactions"`
	UncleHashes  []common.Hash        `json:"uncles"`
}

type Client struct {
	*ethclient.Client
	rpcClient *rpc.Client
}

func NewEthClient(rpcClient *rpc.Client) ethclient2.Client {
	return &Client{
		rpcClient: rpcClient,
		Client:    ethclient.NewClient(rpcClient),
	}
}

func (c *Client) GetLogs(ctx context.Context, fromBlock *big.Int, toBlock *big.Int) ([]*types.Log, error) {
	getLogsReq := struct {
		FromBlock string `json:"fromBlock"`
		ToBlock   string `json:"toBlock"`
	}{
		FromBlock: hexutil.EncodeBig(fromBlock),
		ToBlock:   hexutil.EncodeBig(toBlock),
	}
	var logs []*types.Log
	if err := c.rpcClient.CallContext(ctx, &logs, "eth_getLogs", &getLogsReq); err != nil {
		return nil, err
	}
	return logs, nil
}

func (c *Client) GetLatestBlockNumber(ctx context.Context) (*big.Int, error) {
	var raw json.RawMessage
	err := c.rpcClient.CallContext(ctx, &raw, "eth_blockNumber")
	if err != nil {
		return nil, err
	} else if len(raw) != 0 {
		var res *big.Int
		res, err = hexutil.DecodeBig(strings.Trim(string(raw), "\""))
		if err != nil {
			return nil, err
		} else {
			return res, nil
		}
	} else {
		return nil, fmt.Errorf("failed to process response: %s", string(raw))
	}
}

func (c *Client) CustomBlockByNumber(ctx context.Context, args ...interface{}) (*types.Block, error) {
	var raw json.RawMessage
	err := c.rpcClient.CallContext(ctx, &raw, "eth_getBlockByNumber", args...)
	if err != nil {
		return nil, fmt.Errorf("CustomBlockByNumber/CallContext: %w", err)
	} else if len(raw) == 0 {
		return nil, ethereum.NotFound
	}
	// Decode header and transactions.
	var head *types.Header
	var body rpcBlock
	var txInfo rpcTransaction
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, fmt.Errorf("CustomBlockByNumber/Unmarshal/head: %w", err)
	}
	if head == nil {
		return nil, fmt.Errorf("CustomBlockByNumber: null raw: %s", string(raw))
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, fmt.Errorf("CustomBlockByNumber/Unmarshal/body: %w", err)
	}
	if err := json.Unmarshal(raw, &txInfo); err != nil {
		return nil, fmt.Errorf("CustomBlockByNumber/Unmarshal/txInfo: %w", err)
	}

	if len(body.Transactions) != len(txInfo.TxExtraInfo) {
		return nil, fmt.Errorf("CustomBlockByNumber: incorrect transaction processing")
	}
	// Load uncles because they are not included in the block response.
	var uncles []*types.Header
	// Fill the sender cache of transactions in the block.
	txs := make([]*types.Transaction, len(body.Transactions))
	for i, tx := range body.Transactions {
		if txInfo.TxExtraInfo[i].From != nil {
			setSenderFromServer(tx, *txInfo.TxExtraInfo[i].From, body.Hash)
		}
		txs[i] = tx
	}

	return types.NewBlockWithHeader(head).WithBody(types.Body{
		Transactions: txs,
		Uncles:       uncles,
	}), nil
}

// senderFromServer is a types.Signer that remembers the sender address returned by the RPC
// server. It is stored in the transaction's sender address cache to avoid an additional
// request in TransactionSender.
type senderFromServer struct {
	addr      common.Address
	blockhash common.Hash
}

var errNotCached = errors.New("sender not cached")

func setSenderFromServer(tx *types.Transaction, addr common.Address, block common.Hash) {
	// Use types.Sender for side-effect to store our signer into the cache.
	types.Sender(&senderFromServer{addr, block}, tx)
}

func (s *senderFromServer) Equal(other types.Signer) bool {
	os, ok := other.(*senderFromServer)
	return ok && os.blockhash == s.blockhash
}

func (s *senderFromServer) Sender(tx *types.Transaction) (common.Address, error) {
	if s.addr == (common.Address{}) {
		return common.Address{}, errNotCached
	}
	return s.addr, nil
}

func (s *senderFromServer) ChainID() *big.Int {
	panic("can't sign with senderFromServer")
}
func (s *senderFromServer) Hash(tx *types.Transaction) common.Hash {
	panic("can't sign with senderFromServer")
}
func (s *senderFromServer) SignatureValues(tx *types.Transaction, sig []byte) (R, S, V *big.Int, err error) {
	panic("can't sign with senderFromServer")
}
