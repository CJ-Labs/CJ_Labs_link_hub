package http

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	tonCenterBaseURL = "https://toncenter.com/api/v3"
	apiKey           = "b3f03f4307143359fddb56fbd6e82ac7bf7eaa968582a1b5a768dc319c68d02e"
)

// TestGetTransactionsByMasterchainBlock 测试获取主链区块交易
func TestGetTransactionsByMasterchainBlock(t *testing.T) {
	// 创建 HTTP 客户端配置
	cfg := Config{
		BaseURL: tonCenterBaseURL,
		Timeout: 30 * time.Second,
		Headers: map[string]string{
			"X-API-Key": apiKey,
		},
		RetryCount: 3,
		RetryWait:  time.Second,
	}

	// 创建 HTTP 客户端
	client := NewClient(cfg)

	// 创建上下文
	ctx := context.Background()

	// 定义查询参数
	queryParams := map[string]string{
		"seqno":  "45310979",
		"limit":  "256",
		"offset": "0",
	}

	// 定义响应变量
	var result Response

	// 发送请求
	resp, err := client.Get(ctx, "/transactionsByMasterchainBlock", queryParams, &result)

	// 验证请求是否成功
	require.NoError(t, err, "请求应该成功")
	require.NotNil(t, resp, "响应不应为空")
	require.Equal(t, 200, resp.StatusCode(), "状态码应为 200")

	// 验证响应内容
	require.NotEmpty(t, result.Transactions, "交易列表不应为空")

	// 验证第一笔交易的详细信息
	firstTx := result.Transactions[0]
	t.Run("verify_transaction_details", func(t *testing.T) {
		assert.Equal(t, "-1:5555555555555555555555555555555555555555555555555555555555555555", firstTx.Account)
		assert.Equal(t, "54371914000003", firstTx.Lt)
		assert.Equal(t, 45310979, firstTx.McBlockSeqno)

		// 验证交易描述
		assert.Equal(t, "tick_tock", firstTx.Description.Type)
		assert.False(t, firstTx.Description.Aborted)
		assert.True(t, firstTx.Description.IsTock)

		// 验证计算阶段
		assert.True(t, firstTx.Description.ComputePh.Success)
		assert.Equal(t, "2486", firstTx.Description.ComputePh.GasFees)
		assert.Equal(t, "70000000", firstTx.Description.ComputePh.GasLimit)
		assert.Equal(t, 0, firstTx.Description.ComputePh.ExitCode)

		// 验证账户状态
		assert.Equal(t, "1740941620298", firstTx.AccountStateBefore.Balance)
		assert.Equal(t, "active", firstTx.AccountStateBefore.AccountStatus)
		assert.Equal(t, "1740941620298", firstTx.AccountStateAfter.Balance)
		assert.Equal(t, "active", firstTx.AccountStateAfter.AccountStatus)
	})

	// 验证地址簿
	t.Run("verify_address_book", func(t *testing.T) {
		address := "-1:3333333333333333333333333333333333333333333333333333333333333333"
		addressInfo, exists := result.AddressBook[address]
		assert.True(t, exists, "地址应存在于地址簿中")
		assert.Equal(t, "Ef8zMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzM0vF", addressInfo.UserFriendly)
		assert.Nil(t, addressInfo.Domain)
	})

	// 打印响应信息用于调试
	t.Logf("Response status code: %d", resp.StatusCode())
	t.Logf("First transaction hash: %s", firstTx.Hash)
	t.Logf("Total transactions: %d", len(result.Transactions))
}

func TestGetEthereumTransaction(t *testing.T) {
	// 创建 HTTP 客户端配置
	cfg := Config{
		BaseURL: "https://eth-mainnet.g.alchemy.com/v2/dSe_ey3M3YqwXJtyDnkPFpdNlUtQpafS",
		Timeout: 30 * time.Second,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		RetryCount: 3,
		RetryWait:  time.Second,
	}

	// 创建 HTTP 客户端
	client := NewClient(cfg)

	// 创建上下文
	ctx := context.Background()

	// 构造请求体
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionByHash",
		"params":  []string{"0x9af2d56d7232d6b368a0b060216f0edde355698326eff06064fd7c7765612133"},
		"id":      1,
	}

	// 定义响应变量
	var result JsonRpcResponse

	// 发送 POST 请求
	resp, err := client.Post(ctx, "", requestBody, &result)

	// 基本验证
	require.NoError(t, err, "请求应该成功")
	require.NotNil(t, resp, "响应不应为空")
	require.Equal(t, 200, resp.StatusCode(), "状态码应为 200")

	// 验证响应内容
	t.Run("verify_response_metadata", func(t *testing.T) {
		assert.Equal(t, "2.0", result.JsonRpc)
		assert.Equal(t, 1, result.Id)
	})

	// 验证交易详情
	t.Run("verify_transaction_details", func(t *testing.T) {
		tx := result.Result
		assert.Equal(t, "0x9af2d56d7232d6b368a0b060216f0edde355698326eff06064fd7c7765612133", tx.Hash)
		assert.Equal(t, "0xd3276", tx.Nonce)
		assert.Equal(t, "0x7870dcd64575e362d989b3ff0c40c10b70d926ab0a463e2cff8e72ce5a184696", tx.BlockHash)
		assert.Equal(t, "0x14f9ed4", tx.BlockNumber)
		assert.Equal(t, "0x0", tx.TransactionIndex)
		assert.Equal(t, "0x93793bd1f3e35a0efd098c30e486a860a0ef7551", tx.From)
		assert.Equal(t, "0x68d3a973e7272eb388022a5c6518d9b2a2e66fbf", tx.To)
	})

	// 验证交易参数
	t.Run("verify_transaction_parameters", func(t *testing.T) {
		tx := result.Result
		assert.Equal(t, "0x14f9ed4", tx.Value)
		assert.Equal(t, "0xdfe9fbf73", tx.GasPrice)
		assert.Equal(t, "0x3b1b2", tx.Gas)
		assert.Equal(t, "0x2", tx.Type)
		assert.Equal(t, "0x1", tx.ChainId)
	})

	// 验证签名数据
	t.Run("verify_signature_data", func(t *testing.T) {
		tx := result.Result
		assert.Equal(t, "0x3805d4243d44ce753df2db962a31ffba334681441ab61e227b1aa6c8787375c0", tx.R)
		assert.Equal(t, "0x1910b396d892bbe877aaaedf5a5ff06d1e6cde79d1c0f0eb0cd073323d783889", tx.S)
		assert.Equal(t, "0x1", tx.V)
		assert.Equal(t, "0x1", tx.YParity)
	})

	// 打印调试信息
	t.Logf("Response status code: %d", resp.StatusCode())
	t.Logf("Transaction hash: %s", result.Result.Hash)
	t.Logf("Block number: %s", result.Result.BlockNumber)
}

// Transaction 交易信息结构
type Transaction struct {
	Account      string `json:"account"`
	Hash         string `json:"hash"`
	Lt           string `json:"lt"`
	Now          int64  `json:"now"`
	McBlockSeqno int    `json:"mc_block_seqno"`
	TraceID      string `json:"trace_id"`
	Description  struct {
		Type      string `json:"type"`
		Aborted   bool   `json:"aborted"`
		Destroyed bool   `json:"destroyed"`
		IsTock    bool   `json:"is_tock"`
		StoragePh struct {
			StorageFeesCollected string `json:"storage_fees_collected"`
			StatusChange         string `json:"status_change"`
		} `json:"storage_ph"`
		ComputePh struct {
			Success  bool   `json:"success"`
			GasFees  string `json:"gas_fees"`
			GasUsed  string `json:"gas_used"`
			GasLimit string `json:"gas_limit"`
			ExitCode int    `json:"exit_code"`
			VmSteps  int    `json:"vm_steps"`
		} `json:"compute_ph"`
	} `json:"description"`
	AccountStateBefore struct {
		Hash          string `json:"hash"`
		Balance       string `json:"balance"`
		AccountStatus string `json:"account_status"`
		DataHash      string `json:"data_hash"`
		CodeHash      string `json:"code_hash"`
	} `json:"account_state_before"`
	AccountStateAfter struct {
		Hash          string `json:"hash"`
		Balance       string `json:"balance"`
		AccountStatus string `json:"account_status"`
		DataHash      string `json:"data_hash"`
		CodeHash      string `json:"code_hash"`
	} `json:"account_state_after"`
}

type AddressBook struct {
	UserFriendly string  `json:"user_friendly"`
	Domain       *string `json:"domain"`
}

type Response struct {
	Transactions []Transaction          `json:"transactions"`
	AddressBook  map[string]AddressBook `json:"address_book"`
}

// EthTransaction 以太坊交易结构
type EthTransaction struct {
	Hash                 string   `json:"hash"`
	Nonce                string   `json:"nonce"`
	BlockHash            string   `json:"blockHash"`
	BlockNumber          string   `json:"blockNumber"`
	TransactionIndex     string   `json:"transactionIndex"`
	From                 string   `json:"from"`
	To                   string   `json:"to"`
	Value                string   `json:"value"`
	GasPrice             string   `json:"gasPrice"`
	Gas                  string   `json:"gas"`
	MaxFeePerGas         string   `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string   `json:"maxPriorityFeePerGas"`
	Input                string   `json:"input"`
	R                    string   `json:"r"`
	S                    string   `json:"s"`
	V                    string   `json:"v"`
	YParity              string   `json:"yParity"`
	ChainId              string   `json:"chainId"`
	AccessList           []string `json:"accessList"`
	Type                 string   `json:"type"`
}

// JsonRpcResponse JSON-RPC 响应结构
type JsonRpcResponse struct {
	JsonRpc string         `json:"jsonrpc"`
	Id      int            `json:"id"`
	Result  EthTransaction `json:"result"`
}
