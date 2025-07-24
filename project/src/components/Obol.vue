<template>
  <div class="obol-dashboard">
    <!-- Top Overview Cards -->
    <section class="overview-section">
      <div v-for="card in overviewCards" :key="card.label" class="overview-card" @mouseenter="addCardHover" @mouseleave="removeCardHover">
        <div class="card-header">
          <div class="card-icon">
            <svg :width="24" :height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path v-if="card.label === 'Clusters'" d="M12 2L15.09 8.26L22 9L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9L8.91 8.26L12 2Z"/>
              <path v-else-if="card.label === 'Total Validators'" d="M20 6L9 17L4 12"/>
              <path v-else-if="card.label === 'Lido Protocol APR'" d="M3 3v18h18M8 17l4-4 4 4 4-4M7 8h3v3H7zM17 12h3v3h-3z"/>
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
      <button @click="refreshData" class="retry-button">重試</button>
    </div>

    <!-- Bottom Distribution Cards -->
    <section v-else class="distribution-section">
      <div class="distribution-header">
        <div class="distribution-title">Obol Node Operators</div>
        <button @click="refreshData" class="refresh-button" title="重新載入 Obol 數據">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M23 4v6h-6M1 20v-6h6"/>
            <path d="M20.49 9A9 9 0 0 0 5.64 5.64L1 10m22 4l-4.64 4.36A9 9 0 0 1 3.51 15"/>
          </svg>
          重新載入
        </button>
      </div>
      
      <!-- 搜尋框 -->
      <div class="search-section">
        <div class="search-container">
          <div class="search-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <path d="m21 21-4.35-4.35"/>
            </svg>
          </div>
          <input 
            type="text" 
            v-model="searchQuery" 
            placeholder="搜尋操作者名稱..." 
            class="search-input"
            @input="handleSearch"
          />
          <button 
            v-if="searchQuery" 
            @click="clearSearch" 
            class="clear-search-btn"
            title="清除搜尋"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="m15 9-6 6"/>
              <path d="m9 9 6 6"/>
            </svg>
          </button>
        </div>
      </div>
      
      <!-- 篩選標籤 -->
      <div class="filter-section">
        <div class="filter-tabs">
          <button 
            @click="setActiveFilter('all')" 
            :class="['filter-tab', { active: activeFilter === 'all' }]"
          >
            全部 ({{ totalOperators }})
          </button>
          <button 
            @click="setActiveFilter('lido-obol')" 
            :class="['filter-tab', { active: activeFilter === 'lido-obol' }]"
          >
            Lido x Obol ({{ lidoObolCount }})
          </button>
          <button 
            @click="setActiveFilter('lido-ssv')" 
            :class="['filter-tab', { active: activeFilter === 'lido-ssv' }]"
          >
            Lido x SSV ({{ lidoSsvCount }})
          </button>
          <button 
            @click="setActiveFilter('clever-chameleon')" 
            :class="['filter-tab', { active: activeFilter === 'clever-chameleon' }]"
          >
            Lido x Obol: Clever Chameleon ({{ cleverChameleonCount }})
          </button>
          <button 
            @click="setActiveFilter('zippy-zorilla')" 
            :class="['filter-tab', { active: activeFilter === 'zippy-zorilla' }]"
          >
            Lido x SSV: Zippy Zorilla ({{ zippyZorillaCount }})
          </button>
        </div>
      </div>
      
      <div class="distribution-grid">
        <div v-for="(cluster, index) in filteredClusterCards" :key="index" 
             class="distribution-card" 
             :class="{ 'error-card': cluster.error }"
             @click="cluster.error ? handleErrorCardClick() : goToOperatorDetail(cluster, index)" 
             @mouseenter="cluster.error ? null : addCardHover" 
             @mouseleave="cluster.error ? null : removeCardHover">
          <div class="card-header">
            <div class="card-title-section">
              <div class="operator-number">
                #{{ cluster.originalIndex }}
              </div>
              <span class="card-title">{{ cluster.name || `Operator #${cluster.originalIndex}` }}</span>
            </div>
          </div>
          
          <div class="card-content">
            <div class="validator-stats">
              <div class="stat-item">
                <div class="stat-label">Active Validators</div>
                <div class="stat-value active">{{ cluster.totalVettedValidators }}</div>
              </div>
              <div class="stat-item">
                <div class="stat-label">Exited Validators</div>
                <div class="stat-value exited">{{ cluster.totalExitedValidators }}</div>
              </div>
            </div>
            
            <div class="card-details">
              <div class="detail-row">
                <div class="detail-label">Total:</div>
                <div class="validator-badge">
                  {{ cluster.totalAddedValidators }}
                </div>
              </div>
              
              <div class="detail-row">
                <div class="detail-label">Active:</div>
                <div class="deposited-badge">
                  {{ cluster.totalDepositedValidators }}
                </div>
              </div>
              
              <div class="detail-row">
                <div class="detail-label">Address:</div>
                <div class="address-badge" 
                     @click.stop="copyToClipboard(cluster.rewardAddress)" 
                     :title="cluster.rewardAddress">
                  {{ formatAddress(cluster.rewardAddress) }}
                </div>
              </div>
            </div>
          </div>
          
          <!-- 連線逾時遮罩 -->
          <div v-if="cluster.error" class="connection-timeout-overlay">
            <div class="timeout-content">
              <div class="timeout-icon">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <path d="M12 6v6l4 2"/>
                </svg>
              </div>
              <div class="timeout-text">連線逾時</div>
              <div class="timeout-subtext">無法載入操作者資料</div>
            </div>
          </div>
        </div>
      </div>
    </section>


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
      loadingProgress: '',
      // Lido APR 相關數據
      lidoAPR: null,
      lidoAPRLoading: false,
      lidoAPRError: null,
      // 快取配置
      cacheKey: 'obol_clusters_data',
      cacheExpiryTime: 60 * 60 * 1000, // 1小時快取過期時間
      lastFetchTime: null,
      // 篩選狀態
      activeFilter: 'all', // 'all', 'lido-obol', 'lido-ssv'
      // 搜尋狀態
      searchQuery: '',
    }
  },
  computed: {
    overviewCards() {
      const activeCount = this.clusterCards.filter(cluster => cluster.active && !cluster.error).length
      const totalValidators = this.clusterCards.reduce((sum, cluster) => 
        cluster.error ? sum : sum + cluster.totalVettedValidators, 0
      )
      const errorCount = this.clusterCards.filter(cluster => cluster.error).length

      return [
        {
          label: 'Clusters',
          amount: this.clusterCards.length,
          unit: '個',
          status: errorCount === 0 ? (activeCount > 0 ? 'healthy' : 'warning') : 'error',
          statusText: errorCount === 0 ? (activeCount > 0 ? '正常' : '載入中') : `${errorCount} 個連線逾時`
        },
        {
          label: 'Total Validators',
          amount: totalValidators,
          unit: '個',
          status: totalValidators > 0 ? 'healthy' : 'error',
          statusText: totalValidators > 0 ? '正常' : '載入中'
        },
        {
          label: 'Lido Protocol APR',
          amount: this.lidoAPRLoading ? '載入中' : this.lidoAPRError ? 'N/A' : this.lidoAPR ? `${(this.lidoAPR * 100).toFixed(2)}` : 'N/A',
          unit: this.lidoAPRLoading || this.lidoAPRError || !this.lidoAPR ? '' : '%',
          status: this.lidoAPRLoading ? 'info' : this.lidoAPRError ? 'error' : this.lidoAPR ? 'healthy' : 'warning',
          statusText: this.lidoAPRLoading ? '載入中' : this.lidoAPRError ? '載入失敗' : this.lidoAPR ? '即時 APR' : '無數據'
        }
      ]
    },
    totalActiveValidators() {
      return this.clusterCards.reduce((sum, cluster) => 
        cluster.error ? sum : sum + cluster.totalVettedValidators, 0
      )
    },
    mostActiveOperator() {
      if (this.clusterCards.length === 0) return 'N/A'
      const validClusters = this.clusterCards.filter(cluster => !cluster.error)
      if (validClusters.length === 0) return 'N/A'
      
      const mostActive = validClusters.reduce((max, cluster, index) => {
        return cluster.totalVettedValidators > (max.totalVettedValidators || 0) ? 
          { ...cluster, index } : max
      }, {})
      return mostActive.name || `Operator #${mostActive.index}`
    },
    totalOperators() {
      return this.clusterCards.length;
    },
    lidoObolCount() {
      return this.clusterCards.filter(cluster => cluster.name && cluster.name.includes('Lido x Obol')).length;
    },
    lidoSsvCount() {
      return this.clusterCards.filter(cluster => cluster.name && cluster.name.includes('Lido x SSV')).length;
    },
    cleverChameleonCount() {
      return this.clusterCards.filter(cluster => cluster.name && cluster.name.includes('Lido x Obol: Clever Chameleon')).length;
    },
    zippyZorillaCount() {
      return this.clusterCards.filter(cluster => cluster.name && cluster.name.includes('Lido x SSV: Zippy Zorilla')).length;
    },
    filteredClusterCards() {
      let filtered = this.clusterCards.map((cluster, index) => ({
        ...cluster,
        originalIndex: index
      }));
      
      // 先進行篩選
      if (this.activeFilter !== 'all') {
        filtered = filtered.filter(cluster => {
          if (this.activeFilter === 'lido-obol') {
            return cluster.name && cluster.name.includes('Lido x Obol');
          } else if (this.activeFilter === 'lido-ssv') {
            return cluster.name && cluster.name.includes('Lido x SSV');
          } else if (this.activeFilter === 'clever-chameleon') {
            return cluster.name && cluster.name.includes('Lido x Obol: Clever Chameleon');
          } else if (this.activeFilter === 'zippy-zorilla') {
            return cluster.name && cluster.name.includes('Lido x SSV: Zippy Zorilla');
          }
          return false;
        });
      }
      
      // 再進行搜尋
      if (this.searchQuery.trim()) {
        const query = this.searchQuery.toLowerCase().trim();
        filtered = filtered.filter(cluster => 
          cluster.name && cluster.name.toLowerCase().includes(query)
        );
      }
      
      return filtered;
    }
  },
  async mounted() {
    // 並行載入 Obol 數據和 Lido APR
    await Promise.all([
      this.loadObolData(),
      this.loadLidoAPR()
    ])
  },
      methods: {
      // 載入 Lido Protocol APR
      async loadLidoAPR() {
        this.lidoAPRLoading = true
        this.lidoAPRError = null
        
        try {
          console.log('🚀 開始載入 Lido Protocol APR')
          const apr = await ether_obol.getLidoProtocolAPR()
          this.lidoAPR = apr
          console.log('✅ Lido APR 載入成功:', `${(apr * 100).toFixed(2)}%`)
        } catch (error) {
          console.error('❌ Lido APR 載入失敗:', error)
          this.lidoAPRError = '載入失敗'
        } finally {
          this.lidoAPRLoading = false
        }
      },

      // 檢查快取是否有效
      isCacheValid() {
      try {
        const cachedData = sessionStorage.getItem(this.cacheKey)
        const cacheTimestamp = sessionStorage.getItem(`${this.cacheKey}_timestamp`)
        
        if (!cachedData || !cacheTimestamp) {
          return false
        }
        
        const now = Date.now()
        const cacheAge = now - parseInt(cacheTimestamp)
        
        return cacheAge < this.cacheExpiryTime
      } catch (error) {
        console.error('Error checking cache:', error)
        return false
      }
    },
    
    // 從快取載入資料
    loadFromCache() {
      try {
        const cachedData = sessionStorage.getItem(this.cacheKey)
        const cacheTimestamp = sessionStorage.getItem(`${this.cacheKey}_timestamp`)
        
                 if (cachedData && cacheTimestamp) {
           const data = JSON.parse(cachedData)
           this.clusterCards = data.clusterCards || []
           this.lastFetchTime = parseInt(cacheTimestamp)
           console.log('資料已從快取載入')
           return true
         }
      } catch (error) {
        console.error('Error loading from cache:', error)
        this.clearCache()
      }
      return false
    },
    
    // 儲存資料到快取
    saveToCache() {
      try {
        const dataToCache = {
          clusterCards: this.clusterCards
        }
        const timestamp = Date.now()
        
        sessionStorage.setItem(this.cacheKey, JSON.stringify(dataToCache))
        sessionStorage.setItem(`${this.cacheKey}_timestamp`, timestamp.toString())
        this.lastFetchTime = timestamp
        console.log('資料已儲存到快取')
      } catch (error) {
        console.error('Error saving to cache:', error)
      }
    },
    
    // 清除快取
    clearCache() {
      try {
        sessionStorage.removeItem(this.cacheKey)
        sessionStorage.removeItem(`${this.cacheKey}_timestamp`)
        this.lastFetchTime = null
        console.log('快取已清除')
      } catch (error) {
        console.error('Error clearing cache:', error)
      }
    },
    
    // 主要載入函數（檢查快取或拉取新資料）
    async loadObolData() {
      this.isLoading = true
      this.error = null
      this.loadingProgress = '檢查快取資料...'
      
      // 檢查快取
      if (this.isCacheValid() && this.loadFromCache()) {
        this.isLoading = false
        this.loadingProgress = ''
        return
      }
      
      // 快取無效或不存在，重新拉取資料
      await this.fetchObolData()
    },

    async fetchObolData(forceRefresh = false) {
      this.isLoading = true
      this.error = null
      
      // 如果不是強制重新整理且有有效快取，直接載入快取
      if (!forceRefresh && this.isCacheValid() && this.loadFromCache()) {
        this.isLoading = false
        return
      }
      
      this.loadingProgress = '正在連接 Obol Network...'
      
      try {
              this.loadingProgress = '正在拉取節點操作者資料...'
      const result = await ether_obol.getObolOperatorClustersRegistry()
      
      this.loadingProgress = '正在處理資料...'
      this.clusterCards = result.validOperators.map((operator, index) => {
        // 檢查是否為錯誤結果
        if (operator.error) {
          return {
            active: false,
            name: `Operator #${operator.index}`,
            rewardAddress: '連線逾時',
            totalVettedValidators: 0,
            totalExitedValidators: 0,
            totalAddedValidators: 0,
            totalDepositedValidators: 0,
            status: 'error',
            statusText: '連線逾時',
            error: true,
            errorMessage: operator.errorMessage
          }
        }
        
        // 正常數據處理
        return {
          active: operator[0], // active
          name: operator[1] || `Operator #${index}`, // name
          rewardAddress: operator[2], // rewardAddress
          totalVettedValidators: parseInt(operator[3]) || 0, // totalVettedValidators
          totalExitedValidators: parseInt(operator[4]) || 0, // totalExitedValidators
          totalAddedValidators: parseInt(operator[5]) || 0, // totalAddedValidators
          totalDepositedValidators: parseInt(operator[6]) || 0, // totalDepositedValidators
          status: this.getClusterStatus(operator),
          statusText: this.getClusterStatusText(operator),
          error: false
        }
      })
        
        // 儲存到快取
        this.saveToCache()
        this.loadingProgress = '資料載入完成'
        
      } catch (err) {
        console.error('Error fetching Obol data:', err)
        this.error = '無法載入 Obol Network 資料，請稍後再試'
        this.loadingProgress = ''
      } finally {
        this.isLoading = false
      }
    },
    
    // 強制重新整理資料
    async refreshData() {
      this.clearCache()
      await Promise.all([
        this.fetchObolData(true),
        this.loadLidoAPR()
      ])
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
      // 如果是錯誤狀態，不允許進入詳情頁
      if (cluster.error) {
        console.log('Cannot navigate to operator detail - connection timeout')
        return
      }
      
      // 發送事件到父組件，傳遞操作者數據
      const operatorData = {
        operatorId: cluster.originalIndex,
        operatorData: cluster
      }
      
      console.log('Navigating to operator detail:', operatorData)
      this.$emit('go-to-operator-detail', operatorData)
    },

    // 處理錯誤卡片的點擊
    handleErrorCardClick() {
      // 錯誤卡片不執行任何操作
      console.log('Error card clicked - no action taken')
    },
    setActiveFilter(filter) {
      this.activeFilter = filter;
    },
    
    // 處理搜尋
    handleSearch() {
      // 搜尋是即時的，不需要額外處理
      console.log('搜尋查詢:', this.searchQuery);
    },
    
    // 清除搜尋
    clearSearch() {
      this.searchQuery = '';
    }
  }
}
</script>

<style scoped>
/* Overview Section */
.overview-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
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

.distribution-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
  margin-left: 8px;
  margin-right: 8px;
}

.distribution-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--brand-primary);
}

.search-section {
  margin-bottom: 16px;
  margin-left: 8px;
  margin-right: 8px;
}

.search-container {
  position: relative;
  display: flex;
  align-items: center;
  background: var(--bg-card);
  border: 1px solid rgba(59, 130, 246, 0.15);
  border-radius: 12px;
  padding: 4px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.05);
}

.search-container:focus-within {
  border-color: rgba(59, 130, 246, 0.3);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.1);
}

.search-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  color: var(--text-muted);
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-size: 14px;
  color: var(--text-primary);
  padding: 8px 12px;
  min-width: 0;
}

.search-input::placeholder {
  color: var(--text-muted);
}

.clear-search-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.clear-search-btn:hover {
  background: rgba(239, 68, 68, 0.15);
  transform: scale(1.05);
}

.clear-search-btn:active {
  transform: scale(0.95);
}

.filter-section {
  margin-bottom: 18px;
  margin-left: 8px;
  margin-right: 8px;
}

.filter-tabs {
  display: flex;
  gap: 8px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  padding-bottom: 8px;
  flex-wrap: wrap;
}

.filter-tab {
  padding: 8px 12px;
  border: 1px solid rgba(59, 130, 246, 0.15);
  border-radius: 18px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s ease;
  background: rgba(59, 130, 246, 0.03);
  position: relative;
  overflow: hidden;
  white-space: nowrap;
}

.filter-tab::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.05) 0%, rgba(99, 102, 241, 0.05) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.filter-tab:hover::before {
  opacity: 1;
}

.filter-tab.active {
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.3);
  color: var(--brand-primary);
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.15);
}

.filter-tab.active::before {
  opacity: 1;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(99, 102, 241, 0.1) 100%);
}

.filter-tab:hover {
  border-color: rgba(59, 130, 246, 0.25);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.1);
}

.filter-tab:active {
  transform: translateY(0);
  box-shadow: 0 2px 6px rgba(59, 130, 246, 0.1);
}

.refresh-button {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.refresh-button:hover {
  background: rgba(59, 130, 246, 0.15);
  border-color: rgba(59, 130, 246, 0.3);
  transform: translateY(-1px);
}

.refresh-button:active {
  transform: translateY(0);
}

.refresh-button svg {
  transition: transform 0.3s ease;
}

.refresh-button:hover svg {
  transform: rotate(180deg);
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
  overflow: hidden;
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

/* 移除舊的錯誤樣式，因為現在使用遮罩 */

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

/* 連線逾時遮罩樣式 */
.connection-timeout-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  border-radius: var(--border-radius);
}

.timeout-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  text-align: center;
  color: white;
}

.timeout-icon {
  color: var(--danger);
  opacity: 0.9;
}

.timeout-text {
  font-size: 18px;
  font-weight: 700;
  color: var(--danger);
}

.timeout-subtext {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
  opacity: 0.9;
}

/* 錯誤狀態的卡片樣式 */
.distribution-card.error-card {
  cursor: not-allowed;
}

.distribution-card.error-card:hover {
  transform: none;
  box-shadow: 0 6px 24px rgba(59, 130, 246, 0.10), 0 1.5px 6px rgba(0,0,0,0.08);
  border-color: rgba(0, 0, 0, 0.12);
}

.address-badge {
  cursor: pointer;
  transition: all 0.3s ease;
}

.address-badge:hover {
  background: rgba(99, 102, 241, 0.2);
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
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
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
  .overview-section {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .overview-card {
    padding: 20px 16px;
  }
  
  .overview-card .main-amount {
    font-size: 32px;
  }
  
  .distribution-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
    margin-left: 0;
    margin-right: 0;
  }
  
  .refresh-button {
    align-self: flex-end;
  }

  .search-section {
    margin-left: 0;
    margin-right: 0;
  }
  
  .search-container {
    border-radius: 10px;
  }
  
  .search-input {
    font-size: 13px;
    padding: 6px 10px;
  }
  
  .filter-section {
    margin-left: 0;
    margin-right: 0;
  }

  .filter-tabs {
    flex-wrap: wrap;
    justify-content: center;
    gap: 6px;
  }
  
  .filter-tab {
    font-size: 12px;
    padding: 6px 10px;
    border-radius: 16px;
  }
  
  .distribution-grid {
    grid-template-columns: 1fr;
    gap: 12px;
    padding-left: 0;
    padding-right: 0;
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