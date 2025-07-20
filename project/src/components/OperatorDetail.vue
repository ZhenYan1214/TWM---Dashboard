<template>
  <div class="operator-detail-dashboard">
    
    <!-- 載入進度畫面 -->
    <div v-if="isPageLoading" class="loading-overlay">
      <div class="loading-container">
        <div class="loading-header">
          <div class="loading-icon">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2L15.09 8.26L22 9L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9L8.91 8.26L12 2Z"/>
            </svg>
          </div>
          <h2 class="loading-title">載入操作者資料</h2>
          <p class="loading-subtitle">正在準備 Operator #{{ operatorId }} 的詳細資訊</p>
          
          <!-- 取消按鈕 -->
          <button @click="cancelLoading" class="cancel-loading-btn" title="取消載入並返回 Obol 列表">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
            返回列表
          </button>
        </div>

        <!-- 進度條 -->
        <div class="progress-section">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: loadingProgress + '%' }"></div>
          </div>
          <div class="progress-info">
            <div class="progress-text">{{ Math.round(loadingProgress) }}%</div>
            <div class="progress-hint">點擊下方「返回列表」可隨時返回 Obol 頁面</div>
          </div>
        </div>

        <!-- 載入步驟 -->
        <div class="loading-steps">
          <div v-for="(step, index) in loadingSteps" 
               :key="index" 
               :class="['loading-step', { 
                 'completed': step.completed, 
                 'current': currentLoadingStep === step.name && !step.completed 
               }]">
            <div class="step-indicator">
              <svg v-if="step.completed" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 6L9 17L4 12"/>
              </svg>
              <div v-else-if="currentLoadingStep === step.name" class="step-loading">
                <div class="loading-spinner"></div>
              </div>
              <div v-else class="step-pending"></div>
            </div>
            <div class="step-content">
              <span class="step-text">{{ step.name }}</span>
              <!-- 圖表初始化進度顯示 -->
              <div v-if="step.name === '初始化所有圖表' && chartsInitializing" class="charts-progress">
                <div class="charts-progress-bar">
                  <div class="charts-progress-fill" :style="{ width: chartsProgress + '%' }"></div>
                </div>
                <span class="charts-progress-text">{{ Math.round(chartsProgress) }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 實際頁面內容 -->
    <div v-else class="page-content">
    <!-- Top Overview Section -->
    <section class="overview-section">
      <div class="operator-header-card">
        <div class="operator-header-content">
          <div class="operator-icon">
            <span class="operator-number">#{{ operatorId }}</span>
          </div>
          <div class="operator-info">
            <h2 class="operator-name">{{ operatorName }}</h2>
            <div class="operator-subtitle">
              <span class="operator-id">操作者編號: #{{ operatorId }}</span>
              <div class="status-badge" :class="operatorInfo.status">
                <span class="status-dot"></span>
                <span class="status-text">{{ operatorInfo.statusText }}</span>
              </div>
            </div>
          </div>
          <button @click="goBack" class="back-button">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            返回
          </button>
        </div>
      </div>
    </section>

    <!-- Statistics Overview Cards -->
    <section class="stats-overview-section">
      <div v-for="card in overviewCards" :key="card.label" class="overview-card" @mouseenter="addCardHover" @mouseleave="removeCardHover">
        <div class="card-header">
          <div class="card-icon" :class="card.iconClass">
            <svg :width="24" :height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path :d="card.iconPath"/>
            </svg>
          </div>
          <div class="status-badge" :class="card.status">
            <span class="status-text">{{ card.statusText }}</span>
            <div class="status-indicator"></div>
          </div>
        </div>
        
        <div class="card-content">
          <div class="card-label">{{ card.label }}</div>
          <div class="main-amount">
            {{ card.amount }} <span class="unit">{{ card.unit }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Contract Address Section -->
    <section class="contract-address-section">
      <div class="contract-address-card">
        <div class="address-header">
          <div class="address-icon" @click="openEtherscan(operatorInfo.rewardAddress)" title="點擊在 Etherscan 中查看" v-if="operatorInfo.rewardAddress">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
            </svg>
          </div>
          <div class="address-content">
            <h3 class="address-title">合約地址</h3>
            <p class="address-description contract-address" v-if="operatorInfo.rewardAddress" 
               @click="openEtherscan(operatorInfo.rewardAddress)" 
               title="點擊在 Etherscan 中查看">
              {{ operatorInfo.rewardAddress }}
            </p>
            <p class="address-description" v-else>暫無地址</p>
          </div>
          <div class="etherscan-badge" @click="openEtherscan(operatorInfo.rewardAddress)" v-if="operatorInfo.rewardAddress">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
              <polyline points="15,3 21,3 21,9"/>
              <line x1="10" y1="14" x2="21" y2="3"/>
            </svg>
            查看
          </div>
        </div>
      </div>
    </section>

    <!-- Divider -->
    <div class="divider"></div>

    <!-- Split Wallet Address Section -->
    <section class="split-wallet-section">
      <div class="split-wallet-card">
        <div class="address-header">
          <div class="address-icon" @click="openEtherscan(splitWalletAddress)" title="點擊在 Etherscan 中查看" v-if="splitWalletAddress">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12V7H5a2 2 0 0 1 0-4h14v4"/>
              <path d="M3 5v14a2 2 0 0 0 2 2h16v-5"/>
              <path d="M18 12a2 2 0 0 0 0 4h4v-4h-4z"/>
            </svg>
          </div>
          <div class="address-content">
            <h3 class="address-title">Split Wallet 地址</h3>
            <p class="address-description" v-if="splitWalletLoading">正在載入中...</p>
            <p class="address-description" v-else-if="splitWalletError">載入失敗</p>
            <p class="address-description split-wallet-address" v-else-if="splitWalletAddress" 
               @click="openEtherscan(splitWalletAddress)" 
               title="點擊在 Etherscan 中查看">
              {{ splitWalletAddress }}
            </p>
            <p class="address-description" v-else>暫無地址</p>
          </div>
          <div class="etherscan-badge" @click="openEtherscan(splitWalletAddress)" v-if="splitWalletAddress">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
              <polyline points="15,3 21,3 21,9"/>
              <line x1="10" y1="14" x2="21" y2="3"/>
            </svg>
            查看
          </div>
          <div class="loading-badge" v-else-if="splitWalletLoading">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 2l-2 2m-7.61 2.61L18 12l-2.39 7.39L12 18l-7.39 2.39L6 18l-6-6 6-1.39L8.61 6.61"/>
            </svg>
            載入中
          </div>
        </div>
      </div>
    </section>

    <!-- Reward Share Section -->
    <section class="reward-share-section">
      <div class="reward-share-card">
        <div class="section-header">
          <div class="section-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2L2 7l10 5 10-5-10-5z"/>
              <path d="M2 17l10 5 10-5"/>
              <path d="M2 12l10 5 10-5"/>
            </svg>
          </div>
          <div class="section-content">
            <h3 class="section-title">Cluster 分潤配置</h3>
            <p class="section-description" v-if="rewardShareLoading">正在載入分潤資料...</p>
            <p class="section-description" v-else-if="rewardShareError">{{ rewardShareError }}</p>
            <p class="section-description" v-else-if="rewardShareData">
              共 {{ rewardShareData.rewardAddress?.length || 0 }} 個獎勵地址
            </p>
            <p class="section-description" v-else>等待 Split Wallet 地址載入</p>
          </div>
          <div class="loading-indicator" v-if="rewardShareLoading">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 2l-2 2m-7.61 2.61L18 12l-2.39 7.39L12 18l-7.39 2.39L6 18l-6-6 6-1.39L8.61 6.61"/>
            </svg>
            載入中
          </div>
        </div>

        <!-- Reward Share List -->
        <div v-if="rewardShareData && rewardShareData.rewardAddress" class="reward-share-list">
          <div class="list-header">
            <span class="header-title">獎勵地址分潤配置</span>
            <span class="header-count">{{ rewardShareData.rewardAddress.length }} 個地址</span>
          </div>
          
          <div class="share-items">
            <div v-for="(address, index) in rewardShareData.rewardAddress" 
                 :key="index" 
                 class="share-item">
                             <div class="share-address-info">
                 <div class="address-label">地址 {{ index + 1 }}</div>
                 <div class="address-value" @click="openEtherscan(address)" title="點擊在 Etherscan 中查看">
                   <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                     <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                     <polyline points="15,3 21,3 21,9"/>
                     <line x1="10" y1="14" x2="21" y2="3"/>
                   </svg>
                   {{ address }}
                 </div>
               </div>
              <div class="share-percentage">
                <div class="percentage-label">分潤比例</div>
                <div class="percentage-value">
                  {{ calculateSharePercentage(rewardShareData.rewardShare[index]) }}
                </div>
                
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="!rewardShareLoading && !rewardShareError" class="empty-state">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
            <path d="M12 2L2 7l10 5 10-5-10-5z"/>
            <path d="M2 17l10 5 10-5"/>
            <path d="M2 12l10 5 10-5"/>
          </svg>
          <p class="empty-message">尚未載入分潤資訊</p>
          <p class="empty-description">等待 Split Wallet 地址載入後自動獲取</p>
        </div>
      </div>
    </section>

    <!-- wstETH Token 收益統計 Section -->
    <section class="wsteth-rewards-section">
      <div class="wsteth-rewards-card">
        <div class="section-header">
          <div class="section-icon wsteth-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="M12 6v6l4 2"/>
            </svg>
          </div>
          <div class="section-content">
            <h3 class="section-title">wstETH 收益統計</h3>
            <p class="section-description" v-if="wstETHLoading">正在載入 wstETH 收益數據...</p>
            <p class="section-description" v-else-if="wstETHError">{{ wstETHError }}</p>
            <p class="section-description" v-else-if="wstETHSummary">
              共 {{ incomingTransactionsCount }} 筆收入記錄
            </p>
            <p class="section-description" v-else>等待 Split Wallet 地址載入</p>
          </div>
          <div class="refresh-button" @click="refreshWstETHData" v-if="!wstETHLoading && splitWalletAddress" title="重新載入數據">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M23 4v6h-6M1 20v-6h6"/>
              <path d="M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15"/>
            </svg>
            重新載入
          </div>
          <div class="loading-indicator" v-if="wstETHLoading">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 2l-2 2m-7.61 2.61L18 12l-2.39 7.39L12 18l-7.39 2.39L6 18l-6-6 6-1.39L8.61 6.61"/>
            </svg>
            載入中
          </div>
        </div>

        <!-- wstETH 統計卡片 -->
        <div v-if="wstETHSummary && !wstETHLoading" class="wsteth-stats-cards">
          <div class="stats-card total-received">
            <div class="card-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2L15.09 8.26L22 9L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9L8.91 8.26L12 2Z"/>
              </svg>
            </div>
            <div class="card-content">
              <div class="card-label">總收到 wstETH</div>
              <div class="card-value">{{ formatWstETHAmount(wstETHSummary.totalReceived) }}</div>
              <div class="card-unit">wstETH</div>
            </div>
          </div>

          <div class="stats-card transaction-count">
            <div class="card-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
                <polyline points="10,9 9,9 8,9"/>
              </svg>
            </div>
            <div class="card-content">
              <div class="card-label">收入記錄筆數</div>
              <div class="card-value">{{ incomingTransactionsCount }}</div>
              <div class="card-unit">筆</div>
            </div>
          </div>

          <!-- 收益預估卡片 -->
          <div class="stats-card estimated-earnings">
            <div class="card-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="1" x2="12" y2="23"/>
                <path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/>
              </svg>
            </div>
            <div class="card-content">
              <div class="card-label">目前預估收益</div>
              <div class="card-value" :class="{ 
                'estimated': predictedWstETH === null || predictionError,
                'error': predictionError 
              }">
                {{ formatPredictionResult() }}
              </div>
              <div class="card-unit">wstETH</div>
            </div>
            <div class="estimation-note" v-if="predictionError">
              無法預估(沒有分潤紀錄)
            </div>
            <div class="estimation-note" v-else-if="predictedWstETH === null">
              {{ getOperatorType() }} {{ lidoAPR ? '(APR: ' + (lidoAPR * 100).toFixed(2) + '%)' : '' }}
            </div>
            <div class="estimation-note" v-else>
              基於 {{ getOperatorType() }} (APR: {{ (lidoAPR * 100).toFixed(2) }}%)
            </div>
          </div>
        </div>

        <!-- 交易記錄列表 -->
        <div v-if="incomingTransactions && incomingTransactions.length > 0" class="wsteth-transactions">
          <div class="transactions-header">
            <h4 class="transactions-title">收入記錄</h4>
            <span class="transactions-count">{{ incomingTransactions.length }} 筆收入</span>
          </div>
          
          <div class="transactions-list">
            <div v-for="(tx, index) in (showAllTransactions ? incomingTransactions : incomingTransactions.slice(0, 10))" 
                 :key="tx.hash" 
                 class="transaction-item">
              <div class="tx-info">
                <div class="tx-type incoming">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 5v14M5 12l7 7 7-7"/>
                  </svg>
                  收入
                </div>
                <div class="tx-date">{{ new Date(tx.timeStamp * 1000).toLocaleString('zh-TW') }}</div>
              </div>
              
              <div class="tx-amount incoming">
                <span class="amount-value">
                  +{{ formatWstETHAmount(tx.value) }}
                </span>
                <span class="amount-unit">wstETH</span>
              </div>
              
              <div class="tx-hash" @click="openEtherscanTx(tx.hash)" title="在 Etherscan 中查看交易">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                  <polyline points="15,3 21,3 21,9"/>
                  <line x1="10" y1="14" x2="21" y2="3"/>
                </svg>
                {{ tx.hash.substring(0, 10) }}...
              </div>
            </div>
          </div>
          
          <div v-if="incomingTransactions.length > 10" class="show-more-transactions">
            <button class="show-more-btn" @click="showAllTransactions = !showAllTransactions">
              {{ showAllTransactions ? '收起' : `查看全部 ${incomingTransactions.length} 筆收入記錄` }}
            </button>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="!wstETHLoading && !wstETHError" class="empty-state">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
            <circle cx="12" cy="12" r="10"/>
            <path d="M12 6v6l4 2"/>
          </svg>
          <p class="empty-message">尚未載入 wstETH 收入記錄</p>
          <p class="empty-description">等待 Split Wallet 地址載入後自動獲取收入資訊</p>
        </div>
      </div>
    </section>

    <!-- Historical Trends Section - 切換式圖表 -->
    <section class="historical-trends-section">
      <div class="trends-card">
        <div class="chart-header">
          <div class="header-left">
            <h3 class="chart-title">Validator 狀態趨勢</h3>

          </div>
          
          <!-- 時間維度選擇器 -->
          <div class="period-selector">
            <button v-for="period in availablePeriods" 
                    :key="period.value" 
                    @click="switchPeriod(period.value)"
                    :class="['period-btn', { 
                      active: selectedPeriod === period.value,
                      loading: charts[period.value].loading
                    }]"
                    :disabled="charts[period.value].loading">
              <span v-if="charts[period.value].loading" class="btn-loading-indicator"></span>
              {{ period.label }}
            </button>
          </div>
        </div>

        <!-- 單一圖表容器 -->
        <div class="chart-container">
          <!-- Loading overlay -->
          <div v-if="charts[selectedPeriod].loading" class="chart-loading">
            <div class="loading-spinner-large"></div>
            <div class="loading-text">
              <p>載入{{ selectedPeriodText }}數據中...</p>
              <small>請稍候...</small>
            </div>
          </div>
          
          <!-- Chart canvas - 所有畫布都渲染但只顯示當前選中的 -->
          <canvas v-for="period in availablePeriods"
                  :key="period.value"
                  :ref="`chartCanvas_${period.value}`" 
                  v-show="selectedPeriod === period.value && !charts[period.value].loading && !charts[period.value].error && charts[period.value].data"
                  class="chart-canvas"></canvas>
          
          <!-- Error state -->
          <div v-if="!charts[selectedPeriod].loading && charts[selectedPeriod].error" class="chart-error">
            <div class="error-icon">⚠️</div>
            <p>{{ charts[selectedPeriod].error }}</p>
            <div class="error-actions">
              <button @click="loadChartData(selectedPeriod)" class="retry-btn">重試載入</button>
              <button @click="recreateChart(selectedPeriod)" class="debug-btn">重新創建圖表</button>
            </div>
          </div>
          
          <!-- Empty state -->
          <div v-if="!charts[selectedPeriod].loading && !charts[selectedPeriod].error && !charts[selectedPeriod].data" class="chart-empty">
            <div class="empty-icon">📊</div>
            <p>尚未載入{{ selectedPeriodText }}數據</p>
            <button @click="loadChartData(selectedPeriod)" class="retry-btn">載入數據</button>
          </div>
        </div>
      </div>
    </section>

    </div> <!-- page-content 結束 -->
  </div>
</template>

<script>
import { ether_obol } from '../utils/obol.js'
import { 
  Chart,
  LineController,
  LineElement, 
  PointElement,
  LinearScale,
  CategoryScale,
  Title,
  Tooltip,
  Legend
} from 'chart.js'

// 只註冊必要的 Chart.js 組件
Chart.register(
  LineController,
  LineElement, 
  PointElement,
  LinearScale,
  CategoryScale,
  Title,
  Tooltip,
  Legend
)

export default {
  name: 'OperatorDetail',
  props: {
    operatorData: {
      type: Object,
      required: true
    },
    operatorType: {
      type: String,
      default: 'Obol' // 預設為 Obol，可以是 'Obol' 或 'SSV'
    }
  },
  data() {
    return {
      // 從props初始化數據
      operatorId: '',
      operatorInfo: {
        active: false,
        name: '',
        rewardAddress: '',
        totalVettedValidators: 0,
        totalExitedValidators: 0,
        totalAddedValidators: 0,
        totalDepositedValidators: 0,
        status: 'info',
        statusText: '未知'
      },
      // Split Wallet 相關數據
      splitWalletAddress: null,
      splitWalletLoading: false,
      splitWalletError: null,
      // Reward Share 相關數據
      rewardShareData: null,
      rewardShareLoading: false,
      rewardShareError: null,
      // wstETH Token 相關數據
      wstETHSummary: null,
      wstETHTransactions: null,
      wstETHLoading: false,
      wstETHError: null,
      showAllTransactions: false,
      // 收益預測相關數據
      lidoAPR: null,
      predictedWstETH: null,
      predictionError: null,
      // Chart 相關數據 - 重構為多圖表實例
      charts: {
        '1m': { data: null, loading: false, error: null, instance: null },
        '1y': { data: null, loading: false, error: null, instance: null },
        'all': { data: null, loading: false, error: null, instance: null }
      },
      chartsInitializing: false,
      chartsProgress: 0,
      selectedPeriod: '1m', // 當前選中的時間維度
      // 精簡時間維度選擇
      availablePeriods: [
        { value: '1m', label: '1個月' },
        { value: '1y', label: '1年' },
        { value: 'all', label: '全部' }
      ],
      // 頁面載入進度
      isPageLoading: true,
      loadingProgress: 0,
      loadingCancelled: false, // 載入取消標誌
      loadingSteps: [
        { name: '初始化操作者資訊', completed: false },
        { name: '載入 Split Wallet 地址', completed: false },
        { name: '載入分潤配置資料', completed: false },
        { name: '載入 wstETH 收益數據', completed: false },
        { name: '初始化所有圖表', completed: false }
      ],
      currentLoadingStep: ''
    }
  },
  computed: {
    operatorName() {
      return this.operatorInfo.name || `Operator #${this.operatorId}`
    },

    selectedPeriodText() {
      const period = this.availablePeriods.find(p => p.value === this.selectedPeriod)
      return period ? period.label : '未知'
    },

    // 過濾出收入交易記錄
    incomingTransactions() {
      if (!this.wstETHTransactions || !Array.isArray(this.wstETHTransactions)) {
        return []
      }
      return this.wstETHTransactions.filter(tx => tx.isIncoming === true)
    },

    // 計算收入交易筆數
    incomingTransactionsCount() {
      return this.incomingTransactions.length
    },

    overviewCards() {
      return [
        {
          label: '已退出驗證器',
          amount: this.operatorInfo.totalExitedValidators,
          unit: '個',
          status: 'info',
          statusText: '統計',
          iconClass: 'exited',
          iconPath: 'M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9'
        },
        {
          label: '總添加驗證器',
          amount: this.operatorInfo.totalAddedValidators,
          unit: '個',
          status: 'healthy',
          statusText: '正常',
          iconClass: 'added',
          iconPath: 'M12 2L15.09 8.26L22 9L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9L8.91 8.26L12 2Z'
        },
        {
          label: '已存入驗證器',
          amount: this.operatorInfo.totalDepositedValidators,
          unit: '個',
          status: 'healthy',
          statusText: '正常',
          iconClass: 'deposited',
          iconPath: 'M20 6L9 17L4 12'
        }
      ]
    }
  },
  watch: {
    operatorData: {
      immediate: true,
      handler(newData) {
        if (newData) {
          this.startLoadingSequence(newData)
        }
      }
    }
  },
  methods: {
    initializeData(data) {
      // 從props初始化數據
      if (data && data.operatorData) {
        this.operatorId = data.operatorId || '0'
        this.operatorInfo = {
          active: data.operatorData.active || false,
          name: data.operatorData.name || '',
          rewardAddress: data.operatorData.rewardAddress || '',
          totalVettedValidators: data.operatorData.totalVettedValidators || 0,
          totalExitedValidators: data.operatorData.totalExitedValidators || 0,
          totalAddedValidators: data.operatorData.totalAddedValidators || 0,
          totalDepositedValidators: data.operatorData.totalDepositedValidators || 0,
          status: data.operatorData.status || 'info',
          statusText: data.operatorData.statusText || '未知'
        }
      } else {
        console.warn('No operator data provided, using default values')
      }
    },
    goBack() {
      // 發送事件到父組件返回儀表板
      this.$emit('go-back')
    },

    addCardHover(event) {
      event.currentTarget.style.transform = 'translateY(-2px)'
      event.currentTarget.style.boxShadow = '0 8px 25px rgba(0, 0, 0, 0.15)'
    },
    removeCardHover(event) {
      event.currentTarget.style.transform = 'translateY(0)'
      event.currentTarget.style.boxShadow = '0 2px 8px rgba(0, 0, 0, 0.1)'
    },
    
    // Split Wallet 相關方法（用於手動刷新）
    async fetchSplitWalletData() {
      await this.fetchSplitWalletDataWithProgress()
      
      // 自動載入 Reward Share 資料
      if (this.splitWalletAddress) {
        this.fetchRewardShareData(this.splitWalletAddress)
      }
    },
    
    // Reward Share 相關方法（用於手動刷新）
    async fetchRewardShareData(splitWalletAddress) {
      await this.fetchRewardShareDataWithProgress(splitWalletAddress)
    },
    
    // 計算分潤百分比
    calculateSharePercentage(share) {
      if (!this.rewardShareData || !this.rewardShareData.rewardShare) return '0%'
      
      // 假設總share是所有share的總和
      const totalShare = this.rewardShareData.rewardShare.reduce((sum, s) => sum + Number(s), 0)
      if (totalShare === 0) return '0%'
      
      const percentage = (Number(share) / totalShare * 100).toFixed(2)
      return `${percentage}%`
    },
    
    // 跳轉到 Etherscan
    openEtherscan(address) {
      if (!address) return
      
      const etherscanUrl = `https://etherscan.io/address/${address}`
      window.open(etherscanUrl, '_blank')
    },

    // 跳轉到 Etherscan 交易頁面
    openEtherscanTx(txHash) {
      if (!txHash) return
      
      const etherscanUrl = `https://etherscan.io/tx/${txHash}`
      window.open(etherscanUrl, '_blank')
    },

    // 載入進度相關方法
    async startLoadingSequence(data) {
      this.isPageLoading = true
      this.loadingProgress = 0
      this.loadingCancelled = false // 重置取消標誌
      this.resetLoadingSteps()

      try {
        // 步驟 1: 初始化操作者資訊
        await this.executeLoadingStep('初始化操作者資訊', async () => {
          this.initializeData(data)
          await this.delay(500) // 模擬載入時間
        })

        // 檢查是否已取消
        if (this.loadingCancelled) return

        // 步驟 2: 載入 Split Wallet 地址
        await this.executeLoadingStep('載入 Split Wallet 地址', async () => {
          await this.fetchSplitWalletDataWithProgress()
        })

        // 檢查是否已取消
        if (this.loadingCancelled) return

        // 步驟 3: 載入分潤配置資料
        await this.executeLoadingStep('載入分潤配置資料', async () => {
          if (this.splitWalletAddress) {
            await this.fetchRewardShareDataWithProgress(this.splitWalletAddress)
          } else {
            await this.delay(300) // 如果沒有地址，短暫延遲
          }
        })

        // 檢查是否已取消
        if (this.loadingCancelled) return

        // 步驟 4: 載入 wstETH 收益數據
        await this.executeLoadingStep('載入 wstETH 收益數據', async () => {
          if (this.splitWalletAddress) {
            await this.fetchWstETHDataWithProgress(this.splitWalletAddress)
          } else {
            await this.delay(300) // 如果沒有地址，短暫延遲
          }
        })

        // 檢查是否已取消
        if (this.loadingCancelled) return

        // 步驟 5: 初始化所有圖表
        await this.executeLoadingStep('初始化所有圖表', async () => {
          await this.initializeAllCharts()
        })

        // 檢查是否已取消
        if (this.loadingCancelled) return

        // 完成載入
        this.loadingProgress = 100
        await this.delay(500) // 顯示100%一會兒
        
        // 最後檢查是否已取消
        if (!this.loadingCancelled) {
          this.isPageLoading = false
        }

      } catch (error) {
        console.error('載入序列失敗:', error)
        
        // 如果不是因為取消而失敗，顯示頁面
        if (!this.loadingCancelled) {
          this.isPageLoading = false
        }
      }
    },

    resetLoadingSteps() {
      this.loadingSteps.forEach(step => {
        step.completed = false
      })
      this.currentLoadingStep = ''
    },

    async executeLoadingStep(stepName, asyncFunction) {
      // 檢查是否已取消載入
      if (this.loadingCancelled) {
        console.log('⏹️ 載入已取消，跳過步驟:', stepName, '- 用戶選擇返回 Obol 列表')
        return
      }
      
      this.currentLoadingStep = stepName
      
      try {
        await asyncFunction()
        
        // 再次檢查是否已取消（異步操作完成後）
        if (this.loadingCancelled) {
          console.log('⏹️ 載入已取消，停止處理步驟結果:', stepName, '- 返回 Obol 列表')
          return
        }
        
        // 標記步驟完成
        const step = this.loadingSteps.find(s => s.name === stepName)
        if (step) {
          step.completed = true
        }
        
        // 更新進度
        const completedSteps = this.loadingSteps.filter(s => s.completed).length
        this.loadingProgress = (completedSteps / this.loadingSteps.length) * 100
        
        this.currentLoadingStep = ''
        await this.delay(200) // 步驟間短暫停頓
        
      } catch (error) {
        console.error(`執行載入步驟失敗: ${stepName}`, error)
        
        // 檢查是否為取消導致的錯誤
        if (this.loadingCancelled) {
          console.log('⏹️ 載入已取消，不處理錯誤:', stepName, '- 用戶返回 Obol 列表')
          return
        }
        
        // 即使失敗也標記為完成，繼續下一步
        const step = this.loadingSteps.find(s => s.name === stepName)
        if (step) {
          step.completed = true
        }
        this.currentLoadingStep = ''
      }
    },

    // 取消載入
    cancelLoading() {
      console.log('🚫 用戶取消載入操作，返回 Obol 列表')
      
      // 設置取消標誌
      this.loadingCancelled = true
      
      // 立即隱藏載入畫面
      this.isPageLoading = false
      
      // 重置載入狀態
      this.resetLoadingState()
      
      // 發送事件返回 Obol 畫面
      this.$emit('go-back')
    },

    // 重置載入狀態
    resetLoadingState() {
      this.loadingProgress = 0
      this.currentLoadingStep = ''
      this.loadingCancelled = false
      this.resetLoadingSteps()
      
      // 重置所有載入狀態
      this.splitWalletLoading = false
      this.rewardShareLoading = false
      this.wstETHLoading = false
      this.chartsInitializing = false
      
      // 清理圖表載入狀態
      Object.keys(this.charts).forEach(period => {
        this.charts[period].loading = false
      })
    },

    delay(ms) {
      return new Promise((resolve, reject) => {
        const timer = setTimeout(() => {
          if (this.loadingCancelled) {
            reject(new Error('Loading cancelled'))
          } else {
            resolve()
          }
        }, ms)
        
        // 如果載入已經被取消，立即清理定時器
        if (this.loadingCancelled) {
          clearTimeout(timer)
          reject(new Error('Loading cancelled'))
        }
      })
    },

    async fetchSplitWalletDataWithProgress() {
      // 檢查是否已取消載入
      if (this.loadingCancelled) {
        console.log('⏹️ 載入已取消，跳過 Split Wallet 載入 - 返回 Obol 列表')
        return
      }

      if (!this.operatorInfo.rewardAddress) {
        this.splitWalletError = '沒有獎勵地址，無法查詢 Split Wallet 資料'
        return
      }

      this.splitWalletLoading = true
      this.splitWalletError = null

      try {
        const address = await ether_obol.getObolOperatorSplitWallets(this.operatorInfo.rewardAddress)
        
        // 檢查請求完成後是否已取消
        if (this.loadingCancelled) {
          console.log('⏹️ 載入已取消，忽略 Split Wallet 結果 - 返回 Obol 列表')
          return
        }
        
        this.splitWalletAddress = address
      } catch (error) {
        if (!this.loadingCancelled) {
          console.error('Error fetching split wallet address:', error)
          this.splitWalletError = `載入失敗: ${error.message || '未知錯誤'}`
        }
      } finally {
        if (!this.loadingCancelled) {
          this.splitWalletLoading = false
        }
      }
    },
    
    async fetchRewardShareDataWithProgress(splitWalletAddress) {
      // 檢查是否已取消載入
      if (this.loadingCancelled) {
        console.log('⏹️ 載入已取消，跳過分潤資料載入 - 返回 Obol 列表')
        return
      }

      this.rewardShareLoading = true
      this.rewardShareError = null
      
      try {
        const data = await ether_obol.getObolOperatorRewardshare(splitWalletAddress)
        
        // 檢查請求完成後是否已取消
        if (this.loadingCancelled) {
          console.log('⏹️ 載入已取消，忽略分潤資料結果 - 返回 Obol 列表')
          return
        }
        
        this.rewardShareData = data
      } catch (error) {
        if (!this.loadingCancelled) {
          console.error('Error fetching reward share data:', error)
          this.rewardShareError = `載入失敗: ${error.message || '未知錯誤'}`
        }
      } finally {
        if (!this.loadingCancelled) {
          this.rewardShareLoading = false
        }
      }
    },

    // wstETH 數據載入方法
    async fetchWstETHDataWithProgress(splitWalletAddress) {
      // 檢查是否已取消載入
      if (this.loadingCancelled) {
        console.log('⏹️ 載入已取消，跳過 wstETH 數據載入 - 返回 Obol 列表')
        return
      }

      this.wstETHLoading = true
      this.wstETHError = null
      this.predictionError = null
      
      try {
        console.log('🚀 開始載入 wstETH 數據:', splitWalletAddress)
        
        // 同時載入摘要數據和 Lido APR
        const [summaryData, lidoAPR] = await Promise.all([
          ether_obol.getObolOperatorWstETHSummary(splitWalletAddress),
          ether_obol.getLidoProtocolAPR().catch(error => {
            console.warn('Lido APR 載入失敗:', error)
            return null
          })
        ])
        
        // 檢查請求完成後是否已取消
        if (this.loadingCancelled) {
          console.log('⏹️ 載入已取消，忽略 wstETH 數據結果 - 返回 Obol 列表')
          return
        }
        
        this.wstETHSummary = summaryData
        this.wstETHTransactions = summaryData.transactions || []
        this.lidoAPR = lidoAPR
        
        // 計算收益預測
        try {
          await this.calculateWstETHPrediction(summaryData.transactions || [])
        } catch (predictionError) {
          if (!this.loadingCancelled) {
            console.warn('收益預測計算失敗:', predictionError)
            this.predictionError = predictionError.message
          }
        }
        
        if (!this.loadingCancelled) {
          console.log('✅ wstETH 數據載入成功')
        }
      } catch (error) {
        if (!this.loadingCancelled) {
          console.error('❌ wstETH 數據載入失敗:', error)
          this.wstETHError = `載入失敗: ${error.message || '未知錯誤'}`
        }
      } finally {
        if (!this.loadingCancelled) {
          this.wstETHLoading = false
        }
      }
    },

    // 計算 wstETH 收益預測
    async calculateWstETHPrediction(transactions) {
      try {
        console.log('📊 開始計算收益預測')
        
        // 檢查必要參數
        if (!this.lidoAPR) {
          throw new Error('Lido APR 尚未載入')
        }
        
        // 獲取活躍驗證器數量
        const activeValidators = this.operatorInfo.totalDepositedValidators || 0
        
        // 判斷是否可計算（是否有交易記錄）
        const incomingTxs = transactions.filter(tx => tx.isIncoming === true)
        const isComputable = incomingTxs.length > 0
        
        // 獲取最新一筆交易的時間戳
        let latestTimestamp = null
        if (isComputable && incomingTxs.length > 0) {
          // 交易已按時間排序，取第一筆（最新的）
          latestTimestamp = incomingTxs[0].timeStamp
        }
        
        console.log('📋 預測參數:', {
          lidoAPR: this.lidoAPR,
          activeValidators: activeValidators,
          operatorType: this.operatorType,
          isComputable: isComputable,
          latestTimestamp: latestTimestamp,
          incomingTransactionsCount: incomingTxs.length
        })
        
        // 調用預測函數
        const prediction = await ether_obol.getPridictionWstETH(
          this.lidoAPR,
          activeValidators,
          this.operatorType,
          isComputable,
          latestTimestamp
        )
        
        this.predictedWstETH = prediction
        
        if (prediction !== null && prediction !== undefined) {
          console.log('✅ 收益預測計算成功:', prediction)
        } else {
          console.log('ℹ️ 收益預測結果為空（可能因為不可計算）')
        }
        
      } catch (error) {
        console.error('❌ 收益預測計算失敗:', error)
        this.predictionError = error.message
        throw error
      }
    },

    // 手動重新載入 wstETH 數據
    async refreshWstETHData() {
      if (this.splitWalletAddress) {
        console.log('🔄 手動重新載入 wstETH 數據和預測')
        await this.fetchWstETHDataWithProgress(this.splitWalletAddress)
      }
    },

    // 格式化 wstETH 數量
    formatWstETHAmount(amount) {
      if (!amount || amount === 0) return '0'
      if (amount < 0.001) return amount.toFixed(8)
      if (amount < 1) return amount.toFixed(6)
      return amount.toFixed(4)
    },

    // 格式化預測結果
    formatPredictionResult() {
      if (this.predictionError) {
        return '計算失敗'
      }
      
      if (this.predictedWstETH === null || this.predictedWstETH === undefined) {
        return '無法計算'
      }
      
      return this.formatWstETHAmount(this.predictedWstETH)
    },

    // 獲取操作者類型
    getOperatorType() {
      return this.operatorType || 'Obol'
    },

    // 計算收益率（預留邏輯）
    calculateYield() {
      // 這裡預留給用戶實現收益計算邏輯
      console.log('💡 收益計算邏輯預留給用戶實現')
      return {
        daily: null,
        monthly: null,
        yearly: null,
        apy: null
      }
    },
    
    async loadChartDataWithProgress() {
      // 使用統一的載入邏輯
      await this.loadChartData()
    },

    // ========== 多圖表實例方法 ==========
    
    // 初始化所有圖表
    async initializeAllCharts() {
      if (!this.operatorId) {
        console.warn('沒有操作者 ID，跳過圖表初始化')
        return
      }

      console.log('🚀 開始初始化所有圖表')
      this.chartsInitializing = true
      this.chartsProgress = 0

      try {
        let completedCount = 0

        // 並行載入所有圖表數據
        const promises = this.availablePeriods.map(async (period) => {
          try {
            console.log(`📊 載入圖表: ${period.label}`)
            this.charts[period.value].loading = true
            this.charts[period.value].error = null
            
            const data = await ether_obol.getObolOperatorHistoryValidators(this.operatorId, period.value)
            this.charts[period.value].data = data
            this.charts[period.value].loading = false
            
            console.log(`✅ 圖表載入成功: ${period.label}`)
          } catch (error) {
            console.error(`❌ 圖表載入失敗: ${period.label}`, error)
            this.charts[period.value].error = `載入失敗: ${error.message}`
            this.charts[period.value].loading = false
          }
          
          completedCount++
          this.chartsProgress = (completedCount / this.availablePeriods.length) * 100
        })

        await Promise.all(promises)
        
        // 等待 DOM 更新後渲染圖表
        this.$nextTick(() => {
          this.renderAllCharts()
        })

        console.log('🎉 所有圖表初始化完成')
      } catch (error) {
        console.error('圖表初始化失敗:', error)
      } finally {
        this.chartsInitializing = false
        this.chartsProgress = 100
      }
    },

    // 載入單個圖表數據
    async loadChartData(period) {
      if (!this.operatorId) {
        this.charts[period].error = '沒有操作者 ID'
        return
      }

      console.log(`📊 載入圖表數據: ${period}`)
      this.charts[period].loading = true
      this.charts[period].error = null

      try {
        const data = await ether_obol.getObolOperatorHistoryValidators(this.operatorId, period)
        this.charts[period].data = data
        this.charts[period].loading = false
        
        console.log(`✅ 圖表載入成功: ${period}`)
        
        // 渲染圖表
        this.$nextTick(() => {
          this.renderChart(period)
        })
        
      } catch (error) {
        console.error(`❌ 圖表載入失敗: ${period}`, error)
        this.charts[period].error = `載入失敗: ${error.message}`
        this.charts[period].loading = false
      }
    },

    // 渲染所有圖表
    renderAllCharts() {
      console.log('🎯 開始渲染所有圖表')
      this.availablePeriods.forEach(period => {
        if (this.charts[period.value].data && !this.charts[period.value].loading) {
          this.renderChart(period.value)
        }
      })
    },

    // 渲染單個圖表
    renderChart(period) {
      console.log(`🎯 開始渲染圖表: ${period}`)
      
      // 基本檢查
      if (!Chart) {
        this.charts[period].error = 'Chart.js 未載入'
        return
      }

      const chartData = this.charts[period].data
      if (!chartData || !chartData.length) {
        console.log(`⚠️ 沒有數據，跳過渲染: ${period}`)
        return
      }

      const canvas = this.$refs[`chartCanvas_${period}`]
      if (!canvas || !canvas[0]) {
        console.log(`⚠️ Canvas 不存在，延遲渲染: ${period}`)
        setTimeout(() => this.renderChart(period), 100)
        return
      }

      try {
        // 清理舊圖表
        this.destroyChart(period)

        // 確保 Canvas 可見
        const canvasElement = canvas[0]
        canvasElement.style.display = 'block'
        canvasElement.style.width = '100%'
        canvasElement.style.height = '100%'

        // 準備數據
        const { labels, datasets } = this.prepareChartData(chartData)
        
        // 創建新圖表
        const ctx = canvasElement.getContext('2d')
        this.charts[period].instance = new Chart(ctx, {
          type: 'line',
          data: { labels, datasets },
          options: this.getChartOptions()
        })

        console.log(`🎉 圖表渲染成功: ${period}`)

      } catch (error) {
        console.error(`❌ 圖表渲染失敗: ${period}`, error)
        this.charts[period].error = `渲染失敗: ${error.message}`
      }
    },

    // 數據準備方法
    prepareChartData(chartData) {
      const labels = chartData.map(item => {
        const date = new Date(item.timestamp)
        return date.toLocaleDateString('zh-TW', { month: 'numeric', day: 'numeric' })
      })

      const totalAddedData = chartData.map(item => 
        Math.max(0, Number(item.data?.totalAddedValidators) || 0)
      )
      const totalDepositedData = chartData.map(item => 
        Math.max(0, Number(item.data?.totalDepositedValidators) || 0)
      )
      const totalExitedData = chartData.map(item => 
        Math.max(0, Number(item.data?.totalExitedValidators) || 0)
      )

      const datasets = [
        {
          label: '總添加驗證器',
          data: totalAddedData,
          borderColor: '#3B82F6',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          borderWidth: 2,
          fill: false
        },
        {
          label: '已啟動驗證器', 
          data: totalDepositedData,
          borderColor: '#10B981',
          backgroundColor: 'rgba(16, 185, 129, 0.1)',
          borderWidth: 2,
          fill: false
        },
        {
          label: '已退出驗證器',
          data: totalExitedData,
          borderColor: '#9CA3AF',
          backgroundColor: 'rgba(156, 163, 175, 0.1)',
          borderWidth: 2,
          fill: false
        }
      ]

      return { labels, datasets }
    },
    
    // 5. 圖表配置
    getChartOptions() {
      return {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: true, position: 'top' },
          tooltip: { mode: 'index', intersect: false }
        },
        scales: {
          x: { display: true, grid: { display: false } },
          y: { display: true, beginAtZero: true }
        },
        interaction: { mode: 'index', intersect: false }
      }
    },

    // 清理圖表實例
    destroyChart(period) {
      if (period && this.charts[period] && this.charts[period].instance) {
        try {
          this.charts[period].instance.destroy()
          console.log(`🧹 清理圖表實例: ${period}`)
        } catch (error) {
          console.warn(`清理圖表時出錯: ${period}`, error)
        }
        this.charts[period].instance = null
      }
    },

    // 清理所有圖表實例
    destroyAllCharts() {
      console.log('🧹 清理所有圖表實例')
      this.availablePeriods.forEach(period => {
        this.destroyChart(period.value)
      })
    },

    // 獲取圖表狀態摘要
    getChartsStatusSummary() {
      let loaded = 0, loading = 0, error = 0
      
      this.availablePeriods.forEach(period => {
        const chart = this.charts[period.value]
        if (chart.loading) loading++
        else if (chart.error) error++
        else if (chart.data) loaded++
      })
      
      return { loaded, loading, error, total: this.availablePeriods.length }
    },

    // 切換時間維度
    switchPeriod(period) {
      if (this.selectedPeriod === period) {
        return // 避免重複切換
      }
      
      console.log(`🔄 切換時間維度: ${this.selectedPeriod} -> ${period}`)
      this.selectedPeriod = period
      
      // 如果該時間維度的數據尚未載入，則載入數據
      if (!this.charts[period].data && !this.charts[period].loading) {
        this.loadChartData(period)
      }
    },

    // 重新創建特定圖表
    recreateChart(period) {
      console.log(`🔄 重新創建圖表: ${period}`)
      this.destroyChart(period)
      if (this.charts[period].data) {
        this.$nextTick(() => {
          this.renderChart(period)
        })
      }
    }
  },

  mounted() {
    console.log('📱 OperatorDetail 組件已掛載')
  },

  beforeUnmount() {
    console.log('🧹 OperatorDetail 組件即將卸載，清理所有資源並返回 Obol 列表')
    
    // 取消任何正在進行的載入
    this.loadingCancelled = true
    
    // 清理圖表實例
    this.destroyAllCharts()
    
    // 重置載入狀態
    this.resetLoadingState()
  }
}
</script>

<style scoped>
/* Main Container */
.operator-detail-dashboard {
  min-height: 100vh;
  background: var(--bg-primary);
  padding: 0;
  position: relative;
}

/* 載入進度畫面 */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, var(--bg-primary) 0%, var(--bg-secondary) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.loading-container {
  background: var(--bg-card);
  border: 1px solid var(--border-color, rgba(0, 0, 0, 0.1));
  border-radius: 20px;
  padding: 48px 40px;
  box-shadow: var(--shadow-lg, 0 20px 40px rgba(0, 0, 0, 0.1));
  max-width: 480px;
  width: 90%;
  text-align: center;
}

.loading-header {
  margin-bottom: 32px;
}

.loading-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  color: var(--brand-primary);
  animation: float 3s ease-in-out infinite;
}

.loading-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.loading-subtitle {
  font-size: 16px;
  color: var(--text-secondary);
  margin: 0;
}

/* 進度條 */
.progress-section {
  margin-bottom: 32px;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: var(--progress-bg, rgba(59, 130, 246, 0.1));
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 12px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--brand-primary), var(--brand-secondary));
  border-radius: 4px;
  transition: width 0.5s ease;
  position: relative;
}

.progress-fill::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.3), transparent);
  animation: shimmer 2s infinite;
}

.progress-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.progress-text {
  font-size: 14px;
  font-weight: 600;
  color: var(--brand-primary);
}

.progress-hint {
  font-size: 12px;
  color: var(--text-muted);
  opacity: 0.8;
  font-weight: 400;
}

/* 載入步驟 */
.loading-steps {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.loading-step {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.loading-step.current {
  background: var(--step-current-bg, rgba(59, 130, 246, 0.05));
  border: 1px solid var(--step-current-border, rgba(59, 130, 246, 0.2));
}

.loading-step.completed {
  background: var(--step-completed-bg, rgba(16, 185, 129, 0.05));
}

.step-indicator {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.step-indicator svg {
  color: var(--success);
}

.step-loading {
  width: 16px;
  height: 16px;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(59, 130, 246, 0.3);
  border-top: 2px solid var(--brand-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.step-pending {
  width: 8px;
  height: 8px;
  background: rgba(156, 163, 175, 0.3);
  border-radius: 50%;
}

.step-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.step-text {
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 500;
}

.loading-step.current .step-text {
  color: var(--brand-primary);
  font-weight: 600;
}

.loading-step.completed .step-text {
  color: var(--success);
}

/* 圖表初始化進度條樣式 */
.charts-progress {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
}

.charts-progress-bar {
  flex: 1;
  height: 4px;
  background: var(--progress-bg, rgba(59, 130, 246, 0.1));
  border-radius: 2px;
  overflow: hidden;
}

.charts-progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--brand-primary), var(--brand-secondary));
  border-radius: 2px;
  transition: width 0.3s ease;
  position: relative;
}

.charts-progress-fill::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.4), transparent);
  animation: shimmer 1.5s infinite;
}

.charts-progress-text {
  font-size: 11px;
  font-weight: 600;
  color: var(--brand-primary);
  min-width: 32px;
  text-align: right;
}

/* 動畫 */
@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 取消載入按鈕 */
.cancel-loading-btn {
  margin-top: 20px;
  padding: 10px 20px;
  background: var(--cancel-btn-bg, rgba(239, 68, 68, 0.1));
  border: 1px solid var(--cancel-btn-border, rgba(239, 68, 68, 0.3));
  color: var(--cancel-btn-color, #ef4444);
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
  margin-left: auto;
  margin-right: auto;
}

.cancel-loading-btn:hover {
  background: var(--cancel-btn-hover-bg, rgba(239, 68, 68, 0.2));
  border-color: var(--cancel-btn-hover-border, rgba(239, 68, 68, 0.5));
  color: var(--cancel-btn-hover-color, #dc2626);
  transform: translateY(-1px);
  box-shadow: var(--cancel-btn-hover-shadow, 0 4px 12px rgba(239, 68, 68, 0.2));
}

.cancel-loading-btn:active {
  transform: translateY(0);
  box-shadow: var(--cancel-btn-active-shadow, 0 2px 6px rgba(239, 68, 68, 0.15));
}

.cancel-loading-btn svg {
  transition: transform 0.2s ease;
}

.cancel-loading-btn:hover svg {
  transform: rotate(90deg);
}

/* 深色模式支援 */
@media (prefers-color-scheme: dark) {
  .loading-overlay {
    --bg-primary: #0f172a;
    --bg-secondary: #1e293b;
    --bg-card: #1e293b;
    --text-primary: #f1f5f9;
    --text-secondary: #cbd5e1;
    --text-muted: #94a3b8;
    --border-color: rgba(255, 255, 255, 0.1);
    --shadow-lg: 0 20px 40px rgba(0, 0, 0, 0.4);
    
    /* 進度條深色模式 */
    --progress-bg: rgba(59, 130, 246, 0.2);
    
    /* 載入步驟深色模式 */
    --step-current-bg: rgba(59, 130, 246, 0.15);
    --step-current-border: rgba(59, 130, 246, 0.4);
    --step-completed-bg: rgba(16, 185, 129, 0.15);
    
    /* 取消按鈕深色模式 */
    --cancel-btn-bg: rgba(239, 68, 68, 0.15);
    --cancel-btn-border: rgba(239, 68, 68, 0.4);
    --cancel-btn-color: #fca5a5;
    --cancel-btn-hover-bg: rgba(239, 68, 68, 0.25);
    --cancel-btn-hover-border: rgba(239, 68, 68, 0.6);
    --cancel-btn-hover-color: #fecaca;
    --cancel-btn-hover-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
    --cancel-btn-active-shadow: 0 2px 6px rgba(239, 68, 68, 0.25);
  }
  
  .loading-container {
    background: var(--bg-card);
    border-color: var(--border-color);
    box-shadow: var(--shadow-lg);
  }
  
  .loading-title {
    color: var(--text-primary);
  }
  
  .loading-subtitle {
    color: var(--text-secondary);
  }
  
  .step-text {
    color: var(--text-primary);
  }
  
  .progress-hint {
    color: var(--text-muted);
  }
  
  .step-pending {
    background: var(--text-muted);
  }
  
  .charts-progress-text {
    color: var(--brand-primary);
  }
  
  .loading-icon {
    color: var(--brand-primary);
  }
}

/* Top Overview Section */
.overview-section {
  padding: 32px 24px;
}

.operator-header-card {
  background: var(--bg-card);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  padding: 24px;
  margin-bottom: 32px;
}

.operator-header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
}

.operator-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--brand-primary), var(--brand-secondary));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.operator-number {
  font-size: 24px;
  font-weight: 700;
  color: white;
}

.operator-info {
  flex: 1;
  margin-left: 20px;
}

.operator-name {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.operator-subtitle {
  display: flex;
  align-items: center;
  gap: 16px;
}

.operator-id {
  font-size: 16px;
  color: var(--text-secondary);
  font-weight: 500;
}

.back-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: transparent;
  border: 1px solid var(--brand-primary);
  color: var(--brand-primary);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
  width: fit-content;
}

.back-button:hover {
  background: var(--brand-primary);
  color: white;
}

/* Stats Overview Section */
.stats-overview-section {
  display: flex;
  gap: 24px;
  margin-bottom: 32px;
  padding: 0 24px;
}

.overview-card {
  flex: 1;
  background: var(--bg-card);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  padding: 28px 24px;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  min-height: 160px;
}

.overview-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.overview-card:hover::before {
  opacity: 1;
}

.overview-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.overview-card .card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.overview-card .card-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
}

.overview-card .card-icon.active {
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
}

.overview-card .card-icon.exited {
  background: rgba(156, 163, 175, 0.1);
  color: var(--text-muted);
}

.overview-card .card-icon.added {
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
}

.overview-card .card-icon.deposited {
  background: rgba(99, 102, 241, 0.1);
  color: var(--brand-secondary);
}

.overview-card .status-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
}

.overview-card .status-badge.healthy {
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
}

.overview-card .status-badge.warning {
  background: rgba(245, 158, 11, 0.1);
  color: #F59E0B;
}

.overview-card .status-badge.error {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
}

.overview-card .status-badge.info {
  background: rgba(99, 102, 241, 0.1);
  color: var(--brand-secondary);
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.status-badge.healthy .status-indicator {
  background: var(--success);
}

.status-badge.warning .status-indicator {
  background: #F59E0B;
}

.status-badge.error .status-indicator {
  background: var(--danger);
}

.status-badge.info .status-indicator {
  background: var(--brand-secondary);
}

.overview-card .card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.overview-card .card-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 12px;
}

.overview-card .main-amount {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.1;
}

.overview-card .unit {
  font-size: 18px;
  font-weight: 500;
  color: var(--text-muted);
  margin-left: 0.8rem;
  opacity: 0.7;
}

/* Contract Address Section */
.contract-address-section {
  margin-bottom: 32px;
  padding: 0 24px;
}

.contract-address-card {
  background: var(--bg-card);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  padding: 24px;
  transition: all 0.3s ease;
}

.contract-address-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.address-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.address-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(99, 102, 241, 0.1);
  color: var(--brand-secondary);
  cursor: pointer;
  transition: all 0.3s ease;
}

.address-icon:hover {
  background: rgba(99, 102, 241, 0.2);
  transform: scale(1.05);
}

.address-content {
  flex: 1;
}

.address-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.address-description {
  color: var(--text-secondary);
  font-size: 14px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  word-break: break-all;
}

.etherscan-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.etherscan-badge:hover {
  background: rgba(59, 130, 246, 0.2);
  transform: translateY(-1px);
}

.contract-address, .split-wallet-address {
  cursor: pointer;
  transition: all 0.3s ease;
  word-break: break-all;
  line-height: 1.4;
}

.contract-address:hover, .split-wallet-address:hover {
  color: var(--brand-primary);
  text-decoration: underline;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.status-badge.healthy .status-dot {
  background: var(--success);
}

.status-badge.warning .status-dot {
  background: #F59E0B;
}

.status-badge.error .status-dot {
  background: var(--danger);
}

.status-badge.info .status-dot {
  background: var(--brand-secondary);
}

/* Divider */
.divider {
  width: 100%;
  height: 1px;
  background: linear-gradient(90deg, 
    transparent 0%, 
    rgba(59, 130, 246, 0.2) 20%, 
    rgba(99, 102, 241, 0.3) 50%, 
    rgba(59, 130, 246, 0.2) 80%, 
    transparent 100%
  );
  margin: 32px 0;
}

/* Split Wallet Section */
.split-wallet-section {
  margin-bottom: 32px;
  padding: 0 24px;
}

.split-wallet-card {
  background: var(--bg-card);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  padding: 24px;
  transition: all 0.3s ease;
}

.split-wallet-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.loading-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  background: rgba(156, 163, 175, 0.1);
  color: var(--text-muted);
  border: 1px solid rgba(156, 163, 175, 0.2);
}

.loading-badge svg {
  animation: spin 1s linear infinite;
}

/* Split Wallet specific icon styling */
.split-wallet-card .address-icon {
  background: rgba(156, 163, 175, 0.1);
  color: var(--text-muted);
}

.split-wallet-card .address-icon:hover {
  background: rgba(156, 163, 175, 0.2);
  transform: scale(1.05);
}

/* Reward Share Section */
.reward-share-section {
  margin-bottom: 32px;
  padding: 0 24px;
}

.reward-share-card {
  background: var(--bg-card);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  padding: 24px;
  transition: all 0.3s ease;
}

.reward-share-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.section-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.section-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(245, 158, 11, 0.1);
  color: #F59E0B;
}

.section-content {
  flex: 1;
}

.section-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.section-description {
  color: var(--text-secondary);
  font-size: 14px;
  margin: 0;
}

.loading-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  background: rgba(156, 163, 175, 0.1);
  color: var(--text-muted);
  border: 1px solid rgba(156, 163, 175, 0.2);
}

.loading-indicator svg {
  animation: spin 1s linear infinite;
}

.reward-share-list {
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  padding-top: 20px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.header-count {
  font-size: 14px;
  color: var(--text-secondary);
  background: rgba(245, 158, 11, 0.1);
  padding: 4px 8px;
  border-radius: 4px;
}

.share-items {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.share-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: rgba(0, 0, 0, 0.02);
  border-radius: 8px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
}

.share-item:hover {
  background: rgba(245, 158, 11, 0.02);
  border-color: rgba(245, 158, 11, 0.2);
}

.share-address-info {
  flex: 1;
}

.address-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 6px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.address-value {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  cursor: pointer;
  transition: all 0.3s ease;
  word-break: break-all;
  line-height: 1.4;
}

.address-value:hover {
  color: var(--brand-primary);
  text-decoration: underline;
}

.share-percentage {
  text-align: right;
  min-width: 120px;
}

.percentage-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 4px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.percentage-value {
  font-size: 18px;
  font-weight: 700;
  color: #F59E0B;
  margin-bottom: 2px;
}

.share-raw-value {
  font-size: 11px;
  color: var(--text-muted);
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.empty-state {
  text-align: center;
  padding: 48px 24px;
  color: var(--text-muted);
}

.empty-state svg {
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-message {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.empty-description {
  font-size: 14px;
  margin: 0;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* wstETH Rewards Section */
.wsteth-rewards-section {
  margin-bottom: 32px;
  padding: 0 24px;
}

.wsteth-rewards-card {
  background: var(--bg-card);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  padding: 24px;
  transition: all 0.3s ease;
}

.wsteth-rewards-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.wsteth-icon {
  background: rgba(245, 158, 11, 0.1);
  color: #F59E0B;
}

.refresh-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.refresh-button:hover {
  background: rgba(59, 130, 246, 0.2);
  transform: translateY(-1px);
}

/* wstETH 統計卡片 */
.wsteth-stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  margin: 24px 0;
  padding-top: 20px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

.stats-card {
  background: rgba(0, 0, 0, 0.02);
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.stats-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, transparent, currentColor, transparent);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.stats-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stats-card:hover::before {
  opacity: 0.6;
}

.stats-card.total-received {
  color: var(--success);
}

.stats-card.transaction-count {
  color: var(--brand-primary);
}

.stats-card.estimated-earnings {
  color: var(--brand-secondary);
  position: relative;
}

.stats-card .card-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.05);
  color: inherit;
  flex-shrink: 0;
}

.stats-card .card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stats-card .card-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.stats-card .card-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.stats-card .card-value.negative {
  color: var(--danger);
}

.stats-card .card-value.estimated {
  color: var(--text-muted);
  font-style: italic;
}

.stats-card .card-value.error {
  color: var(--danger);
  font-style: italic;
}

.stats-card .card-unit {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-muted);
  opacity: 0.8;
}

.estimation-note {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 10px;
  color: var(--text-muted);
  background: rgba(0, 0, 0, 0.05);
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 500;
}

/* 交易記錄列表 */
.wsteth-transactions {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

.transactions-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.transactions-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.transactions-count {
  font-size: 14px;
  color: var(--text-secondary);
  background: rgba(245, 158, 11, 0.1);
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.transactions-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.transaction-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: rgba(0, 0, 0, 0.02);
  border: 1px solid rgba(0, 0, 0, 0.04);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.transaction-item:hover {
  background: rgba(245, 158, 11, 0.02);
  border-color: rgba(245, 158, 11, 0.1);
  transform: translateX(2px);
}

.tx-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tx-type {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
}

.tx-type.incoming {
  color: var(--success);
}

.tx-type.outgoing {
  color: #F59E0B;
}

.tx-date {
  font-size: 12px;
  color: var(--text-muted);
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.tx-amount {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  min-width: 120px;
}

.amount-value {
  font-size: 16px;
  font-weight: 700;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.tx-amount.incoming .amount-value {
  color: var(--success);
}

.tx-amount.outgoing .amount-value {
  color: #F59E0B;
}

.amount-unit {
  font-size: 12px;
  color: var(--text-muted);
  font-weight: 500;
}

.tx-hash {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-secondary);
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  cursor: pointer;
  transition: all 0.3s ease;
  padding: 4px 8px;
  border-radius: 4px;
  min-width: 120px;
  justify-content: flex-end;
}

.tx-hash:hover {
  color: var(--brand-primary);
  background: rgba(59, 130, 246, 0.1);
  text-decoration: underline;
}

.show-more-transactions {
  text-align: center;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.04);
}

.show-more-btn {
  padding: 8px 16px;
  background: transparent;
  color: var(--brand-primary);
  border: 1px solid var(--brand-primary);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.show-more-btn:hover {
  background: var(--brand-primary);
  color: white;
}

/* Historical Trends Section - 切換式圖表 */
.historical-trends-section {
  margin-bottom: 32px;
  padding: 0 24px;
}

.trends-card {
  background: var(--bg-card);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.trends-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

/* 圖表標題和控制區域 */
.chart-header {
  padding: 20px 24px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  flex-wrap: wrap;
}

.header-left {
  flex: 1;
  min-width: 200px;
}

.chart-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.chart-status {
  font-size: 14px;
  font-weight: 500;
  padding: 6px 12px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.chart-status.loading {
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
}

.chart-status.error {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.chart-status.success {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.chart-status.empty {
  background: rgba(156, 163, 175, 0.1);
  color: var(--text-muted);
}

/* 時間維度選擇器 */
.period-selector {
  display: flex;
  gap: 4px;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 10px;
  padding: 4px;
  flex-shrink: 0;
}

.period-btn {
  padding: 8px 16px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
  position: relative;
  white-space: nowrap;
}

.period-btn:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
}

.period-btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.period-btn.active {
  background: var(--brand-primary);
  color: white;
  box-shadow: 0 2px 4px rgba(59, 130, 246, 0.3);
}

.period-btn.loading {
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
}

.btn-loading-indicator {
  width: 12px;
  height: 12px;
  border: 2px solid transparent;
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.status-loading {
  display: flex;
  align-items: center;
  gap: 6px;
}

.loading-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  animation: pulse 1.5s infinite;
}

/* 圖表容器 */
.chart-container {
  height: 400px;
  position: relative;
  overflow: hidden;
  background: var(--bg-secondary);
  border-radius: 8px;
  margin: 24px;
}

.chart-canvas {
  width: 100%;
  height: 100%;
  background: transparent;
}

/* 圖表狀態 */
.chart-loading, .chart-error, .chart-empty {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  padding: 40px 20px;
  background: var(--bg-secondary);
}

.loading-spinner-large {
  width: 32px;
  height: 32px;
  border: 3px solid rgba(59, 130, 246, 0.3);
  border-top: 3px solid var(--brand-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

.loading-text {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.loading-text p {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
}

.loading-text small {
  font-size: 13px;
  font-weight: 400;
  color: var(--text-muted);
  opacity: 0.8;
}

.chart-error .error-icon,
.chart-empty .empty-icon {
  font-size: 32px;
  margin-bottom: 16px;
  opacity: 0.6;
}

.chart-error p,
.chart-empty p {
  margin: 0 0 20px 0;
  font-size: 16px;
  text-align: center;
  color: var(--text-primary);
}

.error-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
}

.retry-btn {
  padding: 8px 16px;
  background: var(--brand-primary);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.retry-btn:hover {
  background: var(--brand-secondary);
  transform: translateY(-1px);
}

.debug-btn {
  padding: 8px 16px;
  background: rgba(156, 163, 175, 0.1);
  color: var(--text-secondary);
  border: 1px solid rgba(156, 163, 175, 0.3);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.debug-btn:hover {
  background: rgba(156, 163, 175, 0.2);
  color: var(--text-primary);
}

/* Responsive Design */
@media (max-width: 1200px) {
  .stats-overview-section {
    gap: 20px;
  }
  

}

@media (max-width: 1024px) {
  .stats-overview-section {
    flex-direction: column;
    gap: 16px;
  }
  
  .overview-card {
    min-height: 120px;
    padding: 24px 20px;
  }
  
  .split-wallet-section,
  .reward-share-section,
  .wsteth-rewards-section,
  .historical-trends-section {
    padding: 0 20px;
  }
  
  .wsteth-stats-cards {
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 16px;
  }
  
  .operator-header-content {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .operator-info {
    margin-left: 0;
    margin-top: 16px;
  }

  .chart-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .period-selector {
    width: 100%;
    justify-content: flex-start;
  }

  .chart-container {
    height: 320px;
  }
}

@media (max-width: 768px) {
  .overview-card {
    padding: 20px 16px;
  }
  
  .overview-card .main-amount {
    font-size: 32px;
  }
  
  .split-wallet-section,
  .reward-share-section,
  .wsteth-rewards-section,
  .historical-trends-section {
    padding: 0 16px;
  }
  
  .wsteth-stats-cards {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .transaction-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .tx-amount {
    align-items: flex-start;
    min-width: auto;
  }
  
  .tx-hash {
    min-width: auto;
    justify-content: flex-start;
  }
  
  .split-wallet-section .address-header,
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .share-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .share-percentage {
    text-align: left;
    min-width: auto;
  }
  
  .operator-name {
    font-size: 24px;
  }
  
  .operator-subtitle {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .chart-header {
    padding: 16px;
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .period-selector {
    width: 100%;
    flex-wrap: wrap;
  }

  .period-btn {
    font-size: 13px;
    padding: 6px 12px;
  }

  .chart-container {
    height: 280px;
  }
}

@media (max-width: 480px) {
  .overview-section,
  .stats-overview-section,
  .contract-address-section,
  .split-wallet-section,
  .reward-share-section,
  .wsteth-rewards-section,
  .historical-trends-section {
    padding-left: 16px;
    padding-right: 16px;
  }
  
  .stats-card {
    padding: 16px;
    gap: 12px;
  }
  
  .stats-card .card-value {
    font-size: 20px;
  }
  
  .stats-card .card-icon {
    width: 40px;
    height: 40px;
  }
  
  .transaction-item {
    padding: 12px;
  }
  
  .transactions-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .split-wallet-section .address-description {
    font-size: 12px;
  }
  
  .share-item {
    padding: 12px;
  }
  
  .percentage-value {
    font-size: 16px;
  }
  
  .header-title {
    font-size: 14px;
  }
  
  .address-value {
    font-size: 12px;
    flex-wrap: wrap;
  }
  
  .share-address-info {
    margin-bottom: 8px;
  }

  .chart-header {
    padding: 12px;
    gap: 10px;
  }

  .chart-title {
    font-size: 20px;
  }

  .period-selector {
    width: 100%;
    justify-content: space-between;
  }

  .period-btn {
    flex: 1;
    font-size: 12px;
    padding: 6px 8px;
  }

  .chart-container {
    height: 250px;
    margin: 16px;
  }

  /* 載入畫面響應式 */
  .loading-container {
    padding: 32px 24px;
    max-width: 400px;
  }

  .loading-icon {
    width: 48px;
    height: 48px;
  }

  .loading-title {
    font-size: 20px;
  }

  .loading-subtitle {
    font-size: 14px;
  }

  .loading-steps {
    gap: 8px;
  }

  .loading-step {
    padding: 6px;
  }

  .cancel-loading-btn {
    padding: 8px 16px;
    font-size: 13px;
    margin-top: 16px;
  }

  .progress-hint {
    font-size: 11px;
  }
}

/* 深色模式 + 響應式 */
@media (prefers-color-scheme: dark) and (max-width: 480px) {
  .loading-container {
    box-shadow: 0 15px 30px rgba(0, 0, 0, 0.5);
    border-color: rgba(255, 255, 255, 0.15);
  }
  
  .cancel-loading-btn {
    --cancel-btn-hover-shadow: 0 3px 10px rgba(239, 68, 68, 0.4);
  }
}
</style> 