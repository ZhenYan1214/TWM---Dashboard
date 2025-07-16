<template>
  <div class="obol-dashboard">
    <!-- Top Overview Cards -->
    <section class="overview-section">
      <div v-for="card in overviewCards" :key="card.label" class="overview-card" @mouseenter="addCardHover" @mouseleave="removeCardHover">
        <div class="card-header">
          <div class="card-icon">
            <svg :width="24" :height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path v-if="card.label === '活躍 Clusters'" d="M12 2L15.09 8.26L22 9L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9L8.91 8.26L12 2Z"/>
              <path v-else-if="card.label === '總驗證器數量'" d="M20 6L9 17L4 12"/>
              <path v-else d="M8 2V6M16 2V6M3 10H21M5 4H19C20.1046 4 21 4.89543 21 6V20C21 21.1046 20.1046 22 19 22H5C3.89543 22 3 21.1046 3 20V6C3 4.89543 3.89543 4 5 4Z"/>
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

    <!-- Module Status Indicator -->
    <section class="module-status-section">
      <div class="module-status-card">
        <div class="status-header">
          <div class="status-icon" :class="moduleStatus.status">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="M8 14s1.5-2 4-2 4 2 4 2"/>
              <path d="M9 9h.01M15 9h.01"/>
            </svg>
          </div>
          <div class="status-content">
            <h3 class="status-title">Obol 模組狀態</h3>
            <p class="status-description">{{ moduleStatus.description }}</p>
          </div>
          <div class="status-badge-large" :class="moduleStatus.status">
            {{ moduleStatus.statusText }}
          </div>
        </div>
      </div>
    </section>

    <!-- Divider -->
    <div class="divider"></div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading-section">
      <div class="loading-spinner"></div>
      <p>正在載入 Obol Clusters 資料...</p>
      <p v-if="loadingProgress" class="loading-progress">{{ loadingProgress }}</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-section">
      <div class="error-icon">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M15 9l-6 6"/>
          <path d="M9 9l6 6"/>
        </svg>
      </div>
      <h3>載入錯誤</h3>
      <p>{{ error }}</p>
      <button @click="fetchObolData" class="retry-button">重試</button>
    </div>

    <!-- Bottom Distribution Cards -->
    <section v-else class="distribution-section">
      <div class="distribution-title">Obol Node Operators</div>
      <div class="distribution-grid">
        <div v-for="(cluster, index) in clusterCards" :key="index" class="distribution-card" @click="goToOperatorDetail(cluster, index)" @mouseenter="addCardHover" @mouseleave="removeCardHover">
          <div class="card-header">
            <div class="card-title-section">
              <div class="operator-number">
                #{{ index }}
              </div>
              <span class="card-title">{{ cluster.name || `Operator #${index}` }}</span>
            </div>
          </div>
          
          <div class="card-content">
            <div class="validator-stats">
              <div class="stat-item">
                <div class="stat-label">活躍驗證器</div>
                <div class="stat-value active">{{ cluster.totalVettedValidators }}</div>
              </div>
              <div class="stat-item">
                <div class="stat-label">已退出</div>
                <div class="stat-value exited">{{ cluster.totalExitedValidators }}</div>
              </div>
            </div>
            
            <div class="card-details">
              <div class="detail-row">
                <div class="detail-label">總添加:</div>
                <div class="validator-badge">
                  {{ cluster.totalAddedValidators }}
                </div>
              </div>
              
              <div class="detail-row">
                <div class="detail-label">已存入:</div>
                <div class="deposited-badge">
                  {{ cluster.totalDepositedValidators }}
                </div>
              </div>
              
              <div class="detail-row">
                <div class="detail-label">獎勵地址:</div>
                <div class="address-badge" @click.stop="copyToClipboard(cluster.rewardAddress)" :title="cluster.rewardAddress">
                  {{ formatAddress(cluster.rewardAddress) }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Summary Info Card -->
    <div v-if="!isLoading && !error" class="summary-info-card">
      <div class="summary-title">
        <svg class="summary-icon" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <path d="M12 8v4M12 16h.01"/>
        </svg>
        Obol Network 摘要
      </div>
      <div class="summary-content">
        <div class="summary-row">
          <span class="summary-label">總節點操作者：</span>
          <span class="summary-value">{{ clusterCards.length }} 個</span>
        </div>
        <div class="summary-row">
          <span class="summary-label">活躍驗證器總數：</span>
          <span class="summary-value active">{{ totalActiveValidators }} 個</span>
        </div>
        <div class="summary-row">
          <span class="summary-label">表現最佳操作者：</span>
          <span class="summary-value">{{ mostActiveOperator }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ether_obol } from '@/utils/obol.js'

export default {
  name: 'Obol',
  data() {
    return {
      isLoading: true,
      error: null,
      clusterCards: [],
      moduleStatus: {
        status: 'healthy',
        statusText: '正常運行',
        description: 'Obol Network 正在正常運行，所有系統狀態良好'
      }
    }
  },
  computed: {
    overviewCards() {
      const activeCount = this.clusterCards.filter(cluster => cluster.active).length
      const totalValidators = this.clusterCards.reduce((sum, cluster) => sum + cluster.totalVettedValidators, 0)
      const avgValidators = this.clusterCards.length > 0 ? Math.round(totalValidators / this.clusterCards.length) : 0

      return [
        {
          label: '活躍 Clusters',
          amount: activeCount,
          unit: '個',
          status: activeCount > 0 ? 'healthy' : 'warning',
          statusText: activeCount > 0 ? '運行中' : '待啟動'
        },
        {
          label: '總驗證器數量',
          amount: totalValidators,
          unit: '個',
          status: totalValidators > 0 ? 'healthy' : 'error',
          statusText: totalValidators > 0 ? '活躍' : '無活動'
        },
        {
          label: '平均驗證器數',
          amount: avgValidators,
          unit: '個/操作者',
          status: 'info',
          statusText: '統計'
        }
      ]
    },
    totalActiveValidators() {
      return this.clusterCards.reduce((sum, cluster) => sum + cluster.totalVettedValidators, 0)
    },
    mostActiveOperator() {
      if (this.clusterCards.length === 0) return 'N/A'
      const mostActive = this.clusterCards.reduce((max, cluster, index) => {
        return cluster.totalVettedValidators > (max.totalVettedValidators || 0) ? 
          { ...cluster, index } : max
      }, {})
      return mostActive.name || `Operator #${mostActive.index}`
    }
  },
  async mounted() {
    await this.fetchObolData()
  },
  methods: {
    async fetchObolData() {
      this.isLoading = true
      this.error = null
      
      try {
        const nodeOperators = await ether_obol.getObolOperatorClustersRegistry()
        
        this.clusterCards = nodeOperators.map((operator, index) => ({
          active: operator[0], // active
          name: operator[1] || `Operator #${index}`, // name
          rewardAddress: operator[2], // rewardAddress
          totalVettedValidators: parseInt(operator[3]) || 0, // totalVettedValidators
          totalExitedValidators: parseInt(operator[4]) || 0, // totalExitedValidators
          totalAddedValidators: parseInt(operator[5]) || 0, // totalAddedValidators
          totalDepositedValidators: parseInt(operator[6]) || 0, // totalDepositedValidators
          status: this.getClusterStatus(operator),
          statusText: this.getClusterStatusText(operator)
        }))

        // 更新模組狀態
        this.updateModuleStatus()
        
      } catch (err) {
        console.error('Error fetching Obol data:', err)
        this.error = '無法載入 Obol Network 資料，請稍後再試'
        this.moduleStatus = {
          status: 'error',
          statusText: '連接錯誤',
          description: '無法連接到 Obol Network，請檢查網路連接'
        }
      } finally {
        this.isLoading = false
      }
    },
    getClusterStatus(operator) {
      const active = operator[0]
      const totalVetted = parseInt(operator[3]) || 0
      
      if (!active) return 'error'
      if (totalVetted > 50) return 'healthy'
      if (totalVetted > 10) return 'warning'
      return 'info'
    },
    getClusterStatusText(operator) {
      const active = operator[0]
      const totalVetted = parseInt(operator[3]) || 0
      
      if (!active) return '非活躍'
      if (totalVetted > 50) return '高活躍'
      if (totalVetted > 10) return '中活躍'
      return '低活躍'
    },
    updateModuleStatus() {
      const activeCount = this.clusterCards.filter(cluster => cluster.active).length
      const totalCount = this.clusterCards.length
      
      if (activeCount === 0) {
        this.moduleStatus = {
          status: 'error',
          statusText: '服務異常',
          description: '所有 Cluster 都處於非活躍狀態'
        }
      } else if (activeCount < totalCount * 0.5) {
        this.moduleStatus = {
          status: 'warning',
          statusText: '部分異常',
          description: `${activeCount}/${totalCount} 個 Cluster 處於活躍狀態`
        }
      } else {
        this.moduleStatus = {
          status: 'healthy',
          statusText: '正常運行',
          description: `${activeCount}/${totalCount} 個 Cluster 正在正常運行`
        }
      }
    },
    formatAddress(address) {
      if (!address) return 'N/A'
      return `${address.slice(0, 6)}...${address.slice(-4)}`
    },
    copyToClipboard(text) {
      navigator.clipboard.writeText(text).then(() => {
        console.log('Address copied to clipboard')
      }).catch(err => {
        console.error('Failed to copy address:', err)
      })
    },
    addCardHover(event) {
      event.currentTarget.style.transform = 'translateY(-2px)'
      event.currentTarget.style.boxShadow = '0 8px 25px rgba(0, 0, 0, 0.15)'
    },
    removeCardHover(event) {
      event.currentTarget.style.transform = 'translateY(0)'
      event.currentTarget.style.boxShadow = '0 2px 8px rgba(0, 0, 0, 0.1)'
    },
    goToOperatorDetail(cluster, index) {
      // 發送事件到父組件，傳遞操作者數據
      const operatorData = {
        operatorId: index,
        operatorData: cluster
      }
      
      console.log('Navigating to operator detail:', operatorData)
      this.$emit('go-to-operator-detail', operatorData)
    }
  }
}
</script>

<style scoped>
/* Overview Section */
.overview-section {
  display: flex;
  gap: 24px;
  margin-bottom: 32px;
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

/* Module Status Section */
.module-status-section {
  margin-bottom: 32px;
}

.module-status-card {
  background: var(--bg-card);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  padding: 24px;
}

.status-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-icon.healthy {
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
}

.status-icon.warning {
  background: rgba(245, 158, 11, 0.1);
  color: #F59E0B;
}

.status-icon.error {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
}

.status-content {
  flex: 1;
}

.status-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.status-description {
  color: var(--text-secondary);
  font-size: 14px;
}

.status-badge-large {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
}

.status-badge-large.healthy {
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
}

.status-badge-large.warning {
  background: rgba(245, 158, 11, 0.1);
  color: #F59E0B;
}

.status-badge-large.error {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
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

/* Loading Section */
.loading-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  color: var(--text-secondary);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(59, 130, 246, 0.1);
  border-top: 4px solid var(--brand-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Error Section */
.error-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  color: var(--text-secondary);
}

.error-icon {
  color: var(--danger);
  margin-bottom: 16px;
}

.retry-button {
  margin-top: 16px;
  padding: 8px 16px;
  background: var(--brand-primary);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.3s ease;
}

.retry-button:hover {
  background: var(--brand-secondary);
}

/* Distribution Section */
.distribution-section {
  flex: 1;
  overflow: visible; /* 讓 body 控制整頁滾動 */
}

.distribution-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--brand-primary);
  margin-bottom: 18px;
  margin-left: 8px;
}

.distribution-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
  padding-left: 12px;
  padding-right: 12px;
}

.distribution-card {
  background: var(--bg-card);
  border: 1.5px solid rgba(0, 0, 0, 0.12);
  border-radius: var(--border-radius);
  box-shadow: 0 6px 24px rgba(59, 130, 246, 0.10), 0 1.5px 6px rgba(0,0,0,0.08);
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  position: relative;
  min-height: 220px;
  cursor: pointer;
}

.distribution-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(59, 130, 246, 0.15), 0 3px 12px rgba(0,0,0,0.12);
  border-color: var(--brand-primary);
}

.distribution-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--brand-secondary) 0%, var(--brand-primary) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.distribution-card:hover::before {
  opacity: 1;
}

.distribution-card .card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.distribution-card .card-title-section {
  display: flex;
  align-items: center;
  gap: 10px;
}

.distribution-card .operator-number {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
  font-size: 14px;
  font-weight: 700;
  border: 2px solid rgba(59, 130, 246, 0.2);
}

.distribution-card .card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}



.validator-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.stat-item {
  flex: 1;
  text-align: center;
}

.stat-label {
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  line-height: 1;
}

.stat-value.active {
  color: var(--success);
}

.stat-value.exited {
  color: var(--text-muted);
}

.card-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.detail-label {
  color: var(--text-muted);
  font-weight: 500;
}

.validator-badge,
.deposited-badge,
.address-badge {
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  background: rgba(99, 102, 241, 0.1);
  color: var(--brand-secondary);
  border: 1px solid rgba(99, 102, 241, 0.2);
}

.address-badge {
  cursor: pointer;
  transition: all 0.3s ease;
}

.address-badge:hover {
  background: rgba(99, 102, 241, 0.2);
}

/* Summary Info Card */
.summary-info-card {
  width: 100%;
  max-width: 1180px;
  margin: 36px auto 0 auto;
  background: var(--bg-card);
  border-radius: 18px;
  box-shadow: 0 4px 24px rgba(59,130,246,0.10), 0 1.5px 6px rgba(0,0,0,0.08);
  padding: 22px 24px 16px 24px;
  display: flex;
  flex-direction: column;
  border: 1.5px solid rgba(0,0,0,0.10);
}

.summary-title {
  display: flex;
  align-items: center;
  font-size: 20px;
  font-weight: 700;
  color: var(--brand-primary);
  margin-bottom: 18px;
}

.summary-icon {
  margin-right: 10px;
  color: var(--brand-secondary);
}

.summary-content {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.summary-row {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 15px;
  font-weight: 500;
}

.summary-label {
  color: var(--text-muted);
  min-width: 140px;
}

.summary-value {
  color: var(--text-primary);
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 6px;
}

.summary-value.active {
  color: var(--success);
}



@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* Responsive Design */
@media (max-width: 1200px) {
  .overview-section {
    gap: 20px;
  }
  
  .distribution-grid {
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 20px;
  }
}

@media (max-width: 1024px) {
  .overview-section {
    flex-direction: column;
    gap: 16px;
  }
  
  .overview-card {
    min-height: 120px;
    padding: 24px 20px;
  }
  
  .distribution-grid {
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 16px;
  }
}

@media (max-width: 768px) {
  .overview-card {
    padding: 20px 16px;
  }
  
  .overview-card .main-amount {
    font-size: 32px;
  }
  
  .distribution-grid {
    grid-template-columns: 1fr;
    gap: 12px;
    padding-left: 0;
    padding-right: 0;
  }
  
  .summary-info-card {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .distribution-grid {
    gap: 10px;
  }
  
  .distribution-card {
    padding: 16px;
    min-height: 200px;
  }
}
</style> 