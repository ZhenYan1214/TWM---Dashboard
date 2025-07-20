# 環境變數配置說明

## 快速開始

1. 複製範例配置文件：
   ```bash
   cp .env.sample .env
   ```

2. 根據您的需求修改 `.env` 文件中的配置值

3. 重新啟動開發服務器：
   ```bash
   npm run dev
   ```

## 配置參數說明

### API 配置

| 變數名 | 預設值 | 說明 |
|--------|--------|------|
| `VITE_ETHERSCAN_API_KEY` | `RKFRN7F2EGCMFWPMVTX5RE4ZIQ7E1RHJXI` | Etherscan API 密鑰，用於查詢以太坊鏈上數據 |
| `VITE_RPC_URL` | QuickNode URL | 以太坊節點 RPC 端點 |

### 區塊鏈參數

| 變數名 | 預設值 | 說明 |
|--------|--------|------|
| `VITE_SECONDS_PER_BLOCK` | `12` | 每個區塊的平均時間（秒），以太坊主網約為12秒 |

### 請求批次處理配置

| 變數名 | 預設值 | 說明 |
|--------|--------|------|
| `VITE_BATCH_SIZE` | `8` | 一般操作的批次大小，每批次最多處理的請求數量 |
| `VITE_BATCH_DELAY` | `1000` | 一般操作的批次延遲（毫秒），批次間的等待時間 |
| `VITE_CHART_BATCH_SIZE` | `6` | 圖表數據的批次大小，每批次最多處理的區塊數量 |
| `VITE_CHART_BATCH_DELAY` | `1000` | 圖表數據的批次延遲（毫秒），批次間的等待時間 |

### 外部 API 配置

| 變數名 | 預設值 | 說明 |
|--------|--------|------|
| `VITE_LIDO_API_BASE_URL` | `https://eth-api.lido.fi` | Lido Protocol API 基礎 URL，用於獲取 stETH APR 數據 |

## 效能調優建議

### 網路條件良好時
```env
VITE_BATCH_SIZE=12
VITE_BATCH_DELAY=500
VITE_CHART_BATCH_SIZE=8
VITE_CHART_BATCH_DELAY=800
```

### 網路條件一般時（預設）
```env
VITE_BATCH_SIZE=8
VITE_BATCH_DELAY=1000
VITE_CHART_BATCH_SIZE=6
VITE_CHART_BATCH_DELAY=1000
```

### 網路條件較差時
```env
VITE_BATCH_SIZE=4
VITE_BATCH_DELAY=1500
VITE_CHART_BATCH_SIZE=3
VITE_CHART_BATCH_DELAY=2000
```

## 注意事項

1. **VITE_ 前綴**: 所有在客戶端使用的環境變數都必須以 `VITE_` 開頭
2. **API 限制**: 請注意 Etherscan 和 RPC 提供商的 API 限制
3. **開發模式**: 在開發模式下，控制台會顯示當前配置資訊
4. **重新啟動**: 修改 `.env` 文件後需要重新啟動開發服務器

## 疑難排解

### 如果遇到 API 限制錯誤
- 增加 `BATCH_DELAY` 和 `CHART_BATCH_DELAY` 的值
- 減少 `BATCH_SIZE` 和 `CHART_BATCH_SIZE` 的值

### 如果載入速度太慢
- 減少延遲時間（但要注意 API 限制）
- 增加批次大小（在 API 限制允許的範圍內） 