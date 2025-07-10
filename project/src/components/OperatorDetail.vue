<template>
  <div class="operator-detail-dashboard">
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

  </div>
</template>

<script>
import { ether_obol } from '../utils/obol.js'

export default {
  name: 'OperatorDetail',
  props: {
    operatorData: {
      type: Object,
      required: true
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
      rewardShareError: null
    }
  },
  computed: {
    operatorName() {
      return this.operatorInfo.name || `Operator #${this.operatorId}`
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
          this.initializeData(newData)
          // 自動載入 Split Wallet 資料
          this.$nextTick(() => {
            this.fetchSplitWalletData()
          })
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
    
    // Split Wallet 相關方法
    async fetchSplitWalletData() {
      if (!this.operatorInfo.rewardAddress) {
        this.splitWalletError = '沒有獎勵地址，無法查詢 Split Wallet 資料'
        return
      }

      this.splitWalletLoading = true
      this.splitWalletError = null

      try {
        console.log('Fetching split wallet address for:', this.operatorInfo.rewardAddress)
        const address = await ether_obol.getObolOperatorSplitWallets(this.operatorInfo.rewardAddress)
        this.splitWalletAddress = address
        console.log('Split wallet address received:', address)
        
        // 自動載入 Reward Share 資料
        if (address) {
          this.fetchRewardShareData(address)
        }
      } catch (error) {
        console.error('Error fetching split wallet address:', error)
        this.splitWalletError = `載入失敗: ${error.message || '未知錯誤'}`
      } finally {
        this.splitWalletLoading = false
      }
    },
    
    // Reward Share 相關方法
    async fetchRewardShareData(splitWalletAddress) {
      this.rewardShareLoading = true
      this.rewardShareError = null
      
      try {
        console.log('Fetching reward share data for split wallet:', splitWalletAddress)
        const data = await ether_obol.getObolOperatorRewardshare(splitWalletAddress)
        this.rewardShareData = data
        console.log('Reward share data received:', data)
      } catch (error) {
        console.error('Error fetching reward share data:', error)
        this.rewardShareError = `載入失敗: ${error.message || '未知錯誤'}`
      } finally {
        this.rewardShareLoading = false
      }
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
    }
  }
}
</script>

<style scoped>
/* Main Container */
.operator-detail-dashboard {
  min-height: 100vh;
  background: var(--bg-primary);
  padding: 0;
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
  .reward-share-section {
    padding: 0 20px;
  }
  
  .operator-header-content {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .operator-info {
    margin-left: 0;
    margin-top: 16px;
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
  .reward-share-section {
    padding: 0 16px;
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
}

@media (max-width: 480px) {
  .overview-section,
  .stats-overview-section,
  .contract-address-section,
  .split-wallet-section,
  .reward-share-section {
    padding-left: 16px;
    padding-right: 16px;
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
}
</style> 