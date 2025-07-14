<template>
  <div class="flow-dashboard">
    <!-- Top Overview Cards -->
    <section class="overview-section">
      <div v-for="card in overviewCards" :key="card.label" class="overview-card" @mouseenter="addCardHover" @mouseleave="removeCardHover">
        <div class="card-header">
          <div class="card-icon">
            <svg :width="24" :height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path v-if="card.label === '本週總獎勵'" d="M12 2L15.09 8.26L22 9L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9L8.91 8.26L12 2Z"/>
              <path v-else-if="card.label === '本月總獎勵'" d="M8 2V6M16 2V6M3 10H21M5 4H19C20.1046 4 21 4.89543 21 6V20C21 21.1046 20.1046 22 19 22H5C3.89543 22 3 21.1046 3 20V6C3 4.89543 3.89543 4 5 4Z"/>
              <path v-else d="M12 2L2 7L12 12L22 7L12 2Z"/>
            </svg>
          </div>
          <div class="change-badge" :class="{ 'positive': card.change > 0, 'negative': card.change < 0 }">
            <span class="change-amount">{{ formatChange(card.change) }}</span>
            <svg v-if="card.change > 0" class="arrow" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M7 14L12 9L17 14"/>
            </svg>
            <svg v-else-if="card.change < 0" class="arrow" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M7 10L12 15L17 10"/>
            </svg>
          </div>
        </div>
        
        <div class="card-content">
          <div class="card-label">{{ card.label }}</div>
          <div class="main-amount">
            {{ formatAmount(card.amount) }} <span class="unit">FLOW</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Divider -->
    <div class="divider"></div>

    <!-- Bottom Distribution Cards -->
    <section class="distribution-section">
      <div class="distribution-title">本周 Reward 變動</div>
      <div class="distribution-grid">
        <div v-for="card in distributionCards" :key="card.id" class="distribution-card" @mouseenter="addCardHover" @mouseleave="removeCardHover">
          <div class="card-header">
            <div class="card-title-section">
              <div class="card-icon">
                <svg :width="20" :height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path v-if="card.id === 'node'" d="M12 2L15.09 8.26L22 9L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9L8.91 8.26L12 2Z"/>
                  <path v-else d="M20 7L10 17L5 12"/>
                </svg>
              </div>
              <span class="card-title">{{ card.label }}</span>
            </div>
            <div class="change-indicator" :class="{ 'positive': card.change > 0, 'negative': card.change < 0 }">
              <span class="change-amount">{{ formatChange(card.change) }}</span>
              <svg v-if="card.change > 0" class="arrow" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M7 14L12 9L17 14"/>
              </svg>
              <svg v-else-if="card.change < 0" class="arrow" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M7 10L12 15L17 10"/>
              </svg>
            </div>
          </div>
          
          <div class="card-content">
            <div class="main-amount">
              {{ card.amount === 0 ? '—' : formatAmount(card.amount) }} <span class="unit">FLOW</span>
            </div>
            
            <div class="card-details">
              <div class="detail-row">
                <div class="detail-label">DelegatorID:</div>
                <div class="delegator-badge" :class="{ 'main-node': card.isMainNode }">
                  {{ card.delegatorId }}
                </div>
              </div>
              
              <div class="detail-row">
                <div class="detail-label">Epoch:</div>
                <div class="epoch-badge">
                  {{ card.epochCounter }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    <!-- Summary Info Card -->
    <div class="summary-info-card">
      <div class="summary-title">
        <svg class="summary-icon" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4M12 16h.01"/></svg>
        待更新
      </div>
      <div class="summary-content">
        <div class="summary-row">
          <span class="summary-label">此區塊功能將於日後完工，敬請期待。</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Flow',
  data() {
    return {
      overviewCards: [
        {
          label: 'Delegator #2 總獎勵',
          amount: 847.21,
          change: 27.01,
        },
        {
          label: 'Delegator #3 總獎勵',
          amount: 82.77,
          change: -1.23,
        },
        {
          label: 'Delegator #4 總獎勵',
          amount: 12055.09,
          change: 54.00,
        },
      ],
      distributionCards: [
        {
          id: 'node',
          label: 'Node',
          amount: 1247.0448,
          change: 123.88,
          lastWeek: 8708.34,
          idDisplay: 'e8f4bd64...e7982f92',
          fullId: 'e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92',
          nodeId: 'e8f4bd64...e7982f92',
          delegatorId: 'NODE',
          epochCounter: 186,
          isMainNode: true
        },
        {
          id: 'delegator2',
          label: 'Delegator #2',
          amount: 5.3658,
          change: 4.68,
          lastWeek: 2.1,
          idDisplay: 'e8f4bd64...22f92',
          fullId: 'e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92',
          nodeId: 'e8f4bd64...22f92',
          delegatorId: 'DEL-002',
          epochCounter: 188,
          isMainNode: false
        },
        {
          id: 'delegator3',
          label: 'Delegator #3',
          amount: 0.2083,
          change: -0.12,
          lastWeek: 2.1,
          idDisplay: 'e8f4bd64...32f92',
          fullId: 'e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92',
          nodeId: 'e8f4bd64...32f92',
          delegatorId: 'DEL-003',
          epochCounter: 189,
          isMainNode: false
        },
        {
          id: 'delegator4',
          label: 'Delegator #4',
          amount: 2905.6766,
          change: 2.34,
          lastWeek: 2.1,
          idDisplay: 'e8f4bd64...42f92',
          fullId: 'e8f4bd649d08ecb5afb7023a0c5e8bb10ce56659399665da8cc9d517e7982f92',
          nodeId: 'e8f4bd64...42f92',
          delegatorId: 'DEL-004',
          epochCounter: 190,
          isMainNode: false
        }
      ]
    }
  },
  methods: {
    addCardHover(event) {
      event.currentTarget.style.transform = 'translateY(-2px)';
      event.currentTarget.style.boxShadow = '0 8px 25px rgba(0, 0, 0, 0.15)';
    },
    removeCardHover(event) {
      event.currentTarget.style.transform = 'translateY(0)';
      event.currentTarget.style.boxShadow = '0 2px 8px rgba(0, 0, 0, 0.1)';
    },
    formatAmount(amount) {
      return new Intl.NumberFormat('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      }).format(amount);
    },
    formatChange(change) {
      const sign = change > 0 ? '+' : (change < 0 ? '' : '');
      return `${sign}${change.toFixed(2)}`;
    },
    copyToClipboard(text) {
      navigator.clipboard.writeText(text).then(() => {
        console.log('Copied to clipboard');
      });
    },
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

.overview-card .change-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
}

.overview-card .change-badge.negative {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
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

/* Distribution Section */
.distribution-section {
  flex: 1;
  overflow: hidden;
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
  grid-template-columns: repeat(4, 1fr);
  gap: 32px 28px;
  height: auto;
  padding-left: 12px;
  padding-right: 12px;
}

.distribution-card {
  width: 100%;
  min-width: 0;
  background: var(--bg-card);
  border: 1.5px solid rgba(0, 0, 0, 0.12);
  border-radius: var(--border-radius);
  box-shadow: 0 6px 24px rgba(59, 130, 246, 0.10), 0 1.5px 6px rgba(0,0,0,0.08);
  padding: 24px 20px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: stretch;
  transition: box-shadow 0.3s cubic-bezier(.4,2,.6,1), border 0.3s cubic-bezier(.4,2,.6,1);
  position: relative;
  overflow: visible;
  min-height: 180px;
  height: 180px;
  max-width: 260px;
  margin: 0 auto;
}

.distribution-card:hover {
  transform: none;
  box-shadow: 0 12px 32px rgba(59, 130, 246, 0.16), 0 2px 8px rgba(0,0,0,0.10);
}

.distribution-card:not(:last-child) {
  box-shadow: 0 6px 24px rgba(59, 130, 246, 0.10), 0 1.5px 6px rgba(0,0,0,0.08);
  border-bottom: 2px solid rgba(99, 102, 241, 0.10);
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
  min-height: 28px;
}

.distribution-card .card-title-section {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.distribution-card .card-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
}

.distribution-card .card-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
}

.distribution-card .change-indicator {
  display: flex;
  align-items: center;
  gap: 3px;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
}

.distribution-card .change-indicator.negative {
  background: rgba(239, 68, 68, 0.1);
  color: var(--danger);
}

.distribution-card .card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.distribution-card .main-amount {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.1;
  margin-bottom: 8px;
  margin-top: 8px;
  text-align: left;
  min-height: 38px;
  display: flex;
  align-items: center;
}

.distribution-card .unit {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-muted);
  margin-left: 0.8rem;
  opacity: 0.7;
}

.distribution-card .card-details {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  gap: 8px;
  margin-top: 8px;
}

.distribution-card .detail-row {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  font-size: 12px;
  min-width: 80px;
}

.distribution-card .detail-label {
  color: var(--text-muted);
  font-weight: 500;
  min-width: 70px;
  margin-bottom: 2px;
}

.distribution-card .delegator-badge,
.distribution-card .epoch-badge {
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  background: rgba(99, 102, 241, 0.1);
  color: var(--brand-secondary);
  border: 1px solid rgba(99, 102, 241, 0.2);
  margin-bottom: 2px;
}

.distribution-card .delegator-badge.main-node {
  background: rgba(16, 185, 129, 0.1);
  color: var(--success);
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.distribution-card .epoch-badge {
  background: rgba(59, 130, 246, 0.1);
  color: var(--brand-primary);
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.arrow {
  animation: arrowMove 1.5s infinite alternate ease-in-out;
}

@keyframes arrowMove {
  0% { transform: translateY(0); }
  100% { transform: translateY(-1px); }
}

/* Responsive Design */
@media (max-width: 1200px) {
  .overview-section {
    gap: 20px;
  }
  
  .distribution-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
  }
  .distribution-card {
    max-width: 100%;
    height: 160px;
    min-height: 140px;
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
    grid-template-columns: repeat(3, 1fr);
    gap: 14px;
  }
  
  .distribution-card {
    min-height: 160px;
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
  
  .distribution-card {
    max-width: 100%;
    height: 140px;
    min-height: 120px;
  }
  
  .distribution-card .main-amount {
    font-size: 28px;
  }
}

@media (max-width: 480px) {
  .distribution-grid {
    grid-template-columns: 1fr;
    gap: 10px;
  }
  
  .distribution-card {
    padding: 14px;
    min-height: 120px;
  }
  
  .distribution-card .main-amount {
    font-size: 24px;
  }
}
.summary-info-card {
  width: 100%;
  max-width: 1180px;
  min-width: 320px;
  margin: 36px auto 0 auto;
  background: var(--bg-card);
  border-radius: 18px;
  box-shadow: 0 4px 24px rgba(59,130,246,0.10), 0 1.5px 6px rgba(0,0,0,0.08);
  padding: 22px 24px 16px 24px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
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
  width: 100%;
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
  min-width: 120px;
}
.summary-value {
  color: var(--text-primary);
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 6px;
}
.summary-value.winner {
  color: var(--success);
}
.summary-value.loser {
  color: var(--danger);
}
.summary-value.growth {
  color: var(--brand-primary);
}
.winner-icon {
  color: var(--success);
}
.loser-icon {
  color: var(--danger);
}

@media (max-width: 1200px) {
  .summary-info-card {
    max-width: 100%;
    padding: 24px 16px 16px 16px;
  }
}

@media (max-width: 768px) {
  .summary-info-card {
    max-width: 100%;
    padding: 12px 4vw 10px 4vw;
  }
}
</style> 