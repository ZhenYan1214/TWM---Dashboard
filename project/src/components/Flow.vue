<template>
  <div class="flow-dashboard">
    <!-- Top Overview Cards -->
    <div class="overview-title-block overview-title-flex">
      <span class="overview-title-text">Reward 概覽</span>
      <div class="delegator-filter-popwrap" @click.stop>
        <button class="delegator-filter-btn-pro" @mouseenter="showOverviewFilter = true">
          <svg class="filter-icon" width="18" height="18" viewBox="0 0 20 20" fill="none">
            <path d="M3 5h14M6 9h8M9 13h2" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"/>
          </svg>
          FILTER
        </button>
        <div v-if="showOverviewFilter" class="delegator-popover" @mouseenter="showOverviewFilter = true" @mouseleave="showOverviewFilter = false">
          <div v-for="id in allDelegatorIdList" :key="id" class="delegator-pop-option"
               @click.stop="toggleDelegator(id)">
            <span class="delegator-pop-label">{{ id === -1 ? 'Node' : `Delegator${id}` }}</span>
            <span v-if="selectedDelegators.includes(id)" class="delegator-pop-check">
              <svg width="18" height="18" viewBox="0 0 20 20" fill="none">
                <circle cx="10" cy="10" r="9" fill="#1ecb7b" fill-opacity="0.18"/>
                <path d="M6 10.5L9 13.5L14 7.5" stroke="#1ecb7b" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </span>
          </div>
        </div>
      </div>
    </div>
    <section class="overview-section">
      <div v-for="card in filteredDelegatorTotals" :key="card.label" class="overview-card" @mouseenter="addCardHover" @mouseleave="removeCardHover">
        <div class="card-header">
          <div class="overview-title-section">
            <div class="card-icon">
              <svg width="52" height="52" viewBox="0 0 32 32" fill="none">
                <defs>
                  <linearGradient id="user-blue" x1="0" y1="0" x2="1" y2="1">
                    <stop offset="0%" stop-color="#7eb6ff"/>
                    <stop offset="100%" stop-color="#b3aaff"/>
                  </linearGradient>
                </defs>
                <circle cx="16" cy="16" r="14" fill="url(#user-blue)" opacity="0.18"/>
                <circle cx="16" cy="13" r="5" stroke="url(#user-blue)" stroke-width="2.2" fill="#fff"/>
                <ellipse cx="16" cy="22.5" rx="7" ry="4.5" stroke="url(#user-blue)" stroke-width="2.2" fill="#fff"/>
              </svg>
            </div>
            <span class="overview-card-title">{{ card.label }}</span>
          </div>
          <div class="change-badge" :class="{ 'positive': card.change > 0, 'neutral': card.change === 0 }">
            <span class="change-amount">{{ formatRewardAmount(card.change) }}</span>
            <svg v-if="card.change > 0" class="arrow" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M7 14L12 9L17 14"/>
            </svg>
          </div>
        </div>
        
        <div class="card-content">
          <div class="main-amount">
            {{ formatAmount(card.amount) }} <span class="unit">FLOW</span>
          </div>
          <div class="overview-staked-divider"></div>
          <div class="staked-row">
            <svg class="staked-icon" width="16" height="16" viewBox="0 0 20 20" fill="none">
              <circle cx="10" cy="10" r="9" stroke="#7eb6ff" stroke-width="2"/>
              <path d="M7 10.5L10 13.5L15 8.5" stroke="#7eb6ff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span class="staked-label">Staked</span>
            <span class="staked-amount">{{ formatAmount(card.staked) }}</span>
          </div>
        </div>
      </div>
      <!-- 輪播按鈕已移除，超過三個自動換行 -->
    </section>

    <!-- Divider -->
    <div class="divider"></div>
    <!-- Summary Info Card -->
    <div class="summary-info-card reward-log-card">
      <div class="summary-title">
        <svg class="summary-icon" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 8v4M12 16h.01"/></svg>
        歷史紀錄查詢
      </div>
      <div class="log-toolbar">
        <select v-model="selectedWeek" class="log-week-select">
          <option v-for="week in weekOptions" :key="week.value" :value="week.value">{{ week.label }}</option>
        </select>
        <input v-model="logSearch" class="log-search" placeholder="搜尋名稱/ID/Hash..." />
        <button class="log-export-btn" @click="exportCSV">匯出 CSV</button>
      </div>
      <div class="log-table-wrap" style="position:relative;">
        <table class="reward-log-table">
          <thead>
            <tr>
              <th @click="sortLog('timestamp')">獎勵發放時間 <span v-if="logSort==='timestamp'">{{ logSortDir==='asc'?'▲':'▼' }}</span></th>
              <th @click="sortLog('name')">ID <span v-if="logSort==='name'">{{ logSortDir==='asc'?'▲':'▼' }}</span></th>
              <th @click="sortLog('amount')">本週Reward <span v-if="logSort==='amount'">{{ logSortDir==='asc'?'▲':'▼' }}</span></th>
              <th @click="sortLog('epoch_counter')">Epoch <span v-if="logSort==='epoch_counter'">{{ logSortDir==='asc'?'▲':'▼' }}</span></th>
            </tr>
          </thead>
          <transition-group name="log-row-fade" tag="tbody">
            <tr v-for="(row, idx) in displayedLog" :key="row._logKey" :style="logFadeStyle(idx)">
              <td>{{ row.timestamp.split('T')[0] }}</td>
              <td>{{ row.type==='Node' ? 'Node' : ('Delegator #' + row.delegator_id) }}</td>
              <td :class="{'log-amount-pos': row.amount>0, 'log-amount-neg': row.amount<0, 'log-amount-zero': row.amount===0}" >{{ formatAmount(row.amount) }}</td>
              <td>{{ row.epoch_counter }}</td>
            </tr>
          </transition-group>
          <tr v-if="filteredLog.length===0"><td colspan="4" class="log-empty">查無資料</td></tr>
        </table>
        <div v-if="!logExpanded && filteredLog.length > 3" class="log-fade-mask"></div>
        <transition name="fade">
          <div v-if="filteredLog.length > 3" class="log-expand-wrap">
            <button class="log-expand-btn" @click="toggleLogExpand">
              {{ logExpanded ? '收合' : '展開更多' }}
            </button>
          </div>
        </transition>
      </div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts';

const MultiSelectButton = {
  name: 'MultiSelectButton',
  props: {
    options: Array,
    modelValue: Array,
    labelFn: Function
  },
  emits: ['update:modelValue'],
  data() {
    return {
      open: false,
      localValue: [...(this.modelValue || [])]
    }
  },
  watch: {
    modelValue(val) {
      this.localValue = [...val];
    }
  },
  methods: {
    toggle() { this.open = !this.open; },
    close() { this.open = false; },
    isChecked(opt) { return this.localValue.includes(opt); },
    toggleOption(opt) {
      if (this.isChecked(opt)) {
        this.localValue = this.localValue.filter(v => v !== opt);
      } else {
        this.localValue = [...this.localValue, opt];
      }
      this.$emit('update:modelValue', this.localValue);
    },
    label(opt) { return this.labelFn ? this.labelFn(opt) : opt; },
    handleClickOutside(e) {
      if (!this.$el.contains(e.target)) this.open = false;
    }
  },
  mounted() { document.addEventListener('click', this.handleClickOutside); },
  beforeUnmount() { document.removeEventListener('click', this.handleClickOutside); },
  template: `
    <div class='msb-wrap'>
      <button class='msb-btn' @click.stop='toggle' type='button'>
        <span v-if='localValue.length === 0' class='msb-placeholder'>選擇 Delegator</span>
        <span v-else class='msb-selected'>
          <span v-for='(id, idx) in localValue' :key='id' class='msb-tag'>{{ label(id) }}</span>
        </span>
        <svg class='msb-arrow' width='16' height='16' viewBox='0 0 24 24'><path d='M7 10l5 5 5-5' stroke='currentColor' stroke-width='2' fill='none'/></svg>
      </button>
      <div v-if='open' class='msb-dropdown'>
        <div v-for='opt in options' :key='opt' class='msb-option' @click.stop='toggleOption(opt)'>
          <input type='checkbox' :checked='isChecked(opt)' readonly />
          <span>{{ label(opt) }}</span>
        </div>
      </div>
    </div>
  `
};

export default {
  name: 'Flow',
  components: { MultiSelectButton },
  data() {
    return {
      allDelegatorTotals: [
        // { label: 'Delegator #2 總獎勵', amount: 0, change: 0, delegatorId: 2 }, ...
      ],
      selectedDelegators: [], // 多選篩選器
      allDelegatorIdList: [], // 所有可選 delegatorId

      summaryStats: {
        thisWeekTotal: 0,
        lastWeekTotal: 0,
        growthRate: 0,
      },
      rewardHistory: [], // 歷史資料
      chart: null,
      logSearch: '',
      logStartDate: '',
      logEndDate: '',
      logSort: 'timestamp',
      logSortDir: 'desc',
      selectedWeek: '',
      weekOptions: [],
      logExpanded: false,
      showOverviewFilter: false,
    }
  },
  computed: {
    filteredLog() {
      let arr = (this.rawLog || []);
      // 週篩選
      if (this.selectedWeek && this.selectedWeek !== 'all') {
        const [start, end] = this.selectedWeek.split('|');
        arr = arr.filter(row => row.timestamp >= start && row.timestamp <= end+'T23:59:59');
      }
      // 關鍵字搜尋
      if (this.logSearch) {
        const kw = this.logSearch.toLowerCase();
        arr = arr.filter(row =>
          (row.type==='Node' ? 'node' : ('delegator #'+row.delegator_id).toLowerCase()).includes(kw)
        );
      }
      // 排序
      arr = arr.slice().sort((a,b)=>{
        let v1 = a[this.logSort], v2 = b[this.logSort];
        if (this.logSort==='name') {
          v1 = a.type==='Node' ? 'Node' : ('Delegator #'+a.delegator_id);
          v2 = b.type==='Node' ? 'Node' : ('Delegator #'+b.delegator_id);
        }
        if (v1===undefined) v1=''; if (v2===undefined) v2='';
        if (v1<v2) return this.logSortDir==='asc'?-1:1;
        if (v1>v2) return this.logSortDir==='asc'?1:-1;
        return 0;
      });
      return arr;
    },


    filteredDelegatorTotals() {
      // 只顯示被選中的 delegator
      return this.allDelegatorTotals.filter(card => this.selectedDelegators.includes(card.delegatorId));
    },

    displayedLog() {
      return this.logExpanded ? this.filteredLog : this.filteredLog.slice(0, 3);
    }
  },
  async mounted() { // 跟Api server進行交互
    const res = await fetch('http://localhost:8081/api/rewards/all');
    const data = await res.json();

    // 1. 依照日期分組
    const groupedByDate = {};
    data.forEach(item => {
      const date = item.timestamp.split('T')[0];
      if (!groupedByDate[date]) groupedByDate[date] = [];
      groupedByDate[date].push(item);
    });
    const allDates = Object.keys(groupedByDate).sort().reverse();
    const thisWeekDate = allDates[0];
    const lastWeekDate = allDates[1];
    const thisWeekData = thisWeekDate ? groupedByDate[thisWeekDate] : [];
    const lastWeekData = lastWeekDate ? groupedByDate[lastWeekDate] : [];

    // 2. 動態產生所有 Delegator 卡片（含 Node）
    const delegatorTotals = [];
    const thisWeekDelegators = thisWeekData.filter(item => item.type === 'Delegator');
    thisWeekDelegators.sort((a, b) => a.delegator_id - b.delegator_id);
    const lastWeekDelegators = lastWeekData.filter(item => item.type === 'Delegator');
    const thisWeekNode = thisWeekData.find(item => item.type === 'Node');
    const lastWeekNode = lastWeekData.find(item => item.type === 'Node');

    if (thisWeekNode) { //Node
      delegatorTotals.push({
        label: `Node `,
        amount: thisWeekNode.node_total,
        change: thisWeekNode.amount, // 本周reward
        delegatorId: -1,
        staked: thisWeekNode.node_staked
      });
    }
    thisWeekDelegators.forEach(now => { //Delegator
      delegatorTotals.push({
        label: `Delegator #${now.delegator_id} `,
        amount: now.delegator_total,
        change: now.amount, // 本周reward金額
        delegatorId: now.delegator_id,
        staked: now.delegator_staked
      });
    });
    this.allDelegatorTotals = delegatorTotals;
    // 設定所有可選 delegatorId（含 Node）
    this.allDelegatorIdList = delegatorTotals.map(card => card.delegatorId);
    // 預設只選 Node、Delegator2、Delegator3
    this.selectedDelegators = this.allDelegatorIdList.filter(id => id === -1 || id === 2 || id === 3);



    // 整理 rewardHistory
    const grouped = {};
    data.forEach(item => {
      const date = item.timestamp.split('T')[0];
      if (!grouped[date]) grouped[date] = {};
      if (item.type === 'Node') grouped[date].node_total = item.node_total;
      if (item.type === 'Delegator') {
        grouped[date]['delegator_total' + item.delegator_id] = item['delegator_total' + item.delegator_id];
      }
    });
    this.rewardHistory = Object.entries(grouped).map(([date, vals]) => ({
      date,
      node_total: vals.node_total || 0,
      delegator_total2: vals.delegator_total2 || 0,
      delegator_total3: vals.delegator_total3 || 0,
      delegator_total4: vals.delegator_total4 || 0,
    })).sort((a, b) => a.date.localeCompare(b.date));
    this.$nextTick(this.renderChart);
    // 取得完整 log
    this.rawLog = data.map((row, idx) => ({...row, _logKey: row.timestamp+'-'+(row.type==='Node'?row.node_id:row.delegator_id)+'-'+idx}));
    // 產生週選單
    const allLogDates = Array.from(new Set(this.rawLog.map(r => r.timestamp.split('T')[0]))).sort().reverse();
    const weekList = [];
    // 新增 All 選項
    weekList.push({ value: 'all', label: '全部' });
    for (let i = 0; i < Math.min(6, allLogDates.length); i++) {
      const start = allLogDates[i];
      const end = allLogDates[i];
      weekList.push({
        value: `${start}|${end}`,
        label: i === 0 ? '本週' : (i === 1 ? '上週' : `上上週${i>2?`(-${i-1})`:''}`) + `（${start}）`
      });
    }
    this.weekOptions = weekList;
    // 預設選擇「本週」
    const thisWeek = weekList.find(w => w.label === '本週');
    this.selectedWeek = thisWeek ? thisWeek.value : (weekList[0]?.value || '');
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
      // 完整顯示小數點，不限制位數
      return new Intl.NumberFormat('en-US', {
        minimumFractionDigits: 0,
        maximumFractionDigits: 20
      }).format(amount);
    },
    formatChange(change) {
      const sign = change > 0 ? '+' : (change < 0 ? '' : '');
      // 完整顯示小數點
      return `${sign}${change.toString()}`;
    },
    formatRewardAmount(amount) {
      // 專門用於可提領餘額卡片的 reward 金額顯示，完整顯示小數點
      return `+${amount.toString()}`;
    },
    copyToClipboard(text) {
      navigator.clipboard.writeText(text).then(() => {
        console.log('Copied to clipboard');
      });
    },
    renderChart() {
      if (!this.chart) {
        this.chart = echarts.init(this.$refs.rewardChart);
      }
      const option = {
        tooltip: { trigger: 'axis' },
        legend: { data: ['Node', 'Delegator #2', 'Delegator #3', 'Delegator #4'] },
        xAxis: { type: 'category', data: this.rewardHistory.map(d => d.date) },
        yAxis: { type: 'value' },
        series: [
          { name: 'Node', type: 'line', data: this.rewardHistory.map(d => d.node_total) },
          { name: 'Delegator #2', type: 'line', data: this.rewardHistory.map(d => d.delegator_total2) },
          { name: 'Delegator #3', type: 'line', data: this.rewardHistory.map(d => d.delegator_total3) },
          { name: 'Delegator #4', type: 'line', data: this.rewardHistory.map(d => d.delegator_total4) },
        ]
      };
      this.chart.setOption(option);
    },
    sortLog(field) {
      if (this.logSort === field) {
        this.logSortDir = this.logSortDir === 'asc' ? 'desc' : 'asc';
      } else {
        this.logSort = field;
        this.logSortDir = 'desc';
      }
    },
    exportCSV() {
      // 取得目前週期參數
      const week = this.selectedWeek;
      const url = `http://localhost:8081/api/rewards/export?week=${encodeURIComponent(week)}`;
      fetch(url)
        .then(res => res.blob())
        .then(blob => {
          const a = document.createElement('a');
          a.href = URL.createObjectURL(blob);
          a.download = 'reward_log.csv';
          document.body.appendChild(a);
          a.click();
          setTimeout(() => {
            document.body.removeChild(a);
            URL.revokeObjectURL(a.href);
          }, 100);
        });
    },
    toggleLogExpand() {
      this.logExpanded = !this.logExpanded;
    },
    logFadeStyle(idx) {
      if (this.logExpanded || this.filteredLog.length <= 3) return {};
      // 只顯示3筆時，讓第2、3筆漸隱
      if (this.displayedLog.length === 3) {
        if (idx === 1) return { opacity: 0.8 };
        if (idx === 2) return { opacity: 0.5 };
      }
      // 只顯示2筆時，讓第2筆漸隱
      if (this.displayedLog.length === 2 && idx === 1) return { opacity: 0.7 };
      return {};
    },
    // updateOverviewCards/nextOverview 已移除
    // nextDistribution 已移除
    handleGlobalClick() {
      this.showOverviewFilter = false;
    },
    toggleDelegator(id) {
      if (this.selectedDelegators.includes(id)) {
        this.selectedDelegators = this.selectedDelegators.filter(v => v !== id);
      } else {
        this.selectedDelegators = [...this.selectedDelegators, id];
      }
    },
  }
}
</script>

<style scoped>
</style>

<style>
:root {
  --change-badge-bg: rgba(16, 185, 129, 0.13);
  --change-badge-color: #16a34a;
  --change-badge-bg-neg: rgba(239, 68, 68, 0.13);
  --change-badge-color-neg: #ef4444;
  --change-badge-bg-neutral: #f3f4f6;
}
.dark {
  --change-badge-bg: rgba(16, 185, 129, 0.1);
  --change-badge-color: var(--success);
  --change-badge-bg-neg: rgba(239, 68, 68, 0.1);
  --change-badge-color-neg: var(--danger);
  --change-badge-bg-neutral: #23263a;
}
@media (prefers-color-scheme: dark) {
  :root {
    --change-badge-bg: rgba(16, 185, 129, 0.1);
    --change-badge-color: var(--success);
    --change-badge-bg-neg: rgba(239, 68, 68, 0.1);
    --change-badge-color-neg: var(--danger);
    --change-badge-bg-neutral: #23263a;
  }
}
/* Overview Section */
/* 改為 grid，每行三個卡片，自動換行 */
.overview-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-bottom: 32px;
}

.overview-card {
  background: var(--bg-card);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 18px;
  box-shadow: 0 4px 24px rgba(59,130,246,0.10), 0 1.5px 6px rgba(0,0,0,0.08), 0 0.5px 2px rgba(126,182,255,0.08);
  padding: 12px 10px;
  display: flex;
  flex-direction: column;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  min-height: 120px;
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
  padding: 6px 10px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  background: var(--change-badge-bg, rgba(16, 185, 129, 0.1));
  color: var(--change-badge-color, var(--success));
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.overview-card .change-badge.negative {
  background: var(--change-badge-bg-neg, rgba(239, 68, 68, 0.1));
  color: var(--change-badge-color-neg, var(--danger));
}

.overview-card .change-badge.neutral {
  background: var(--change-badge-bg-neutral, #f3f4f6);
  color: var(--brand-primary);
  font-style: normal;
  font-weight: 700;
  letter-spacing: 0.5px;
  border: 1.5px solid var(--brand-primary);
  box-shadow: none;
}

.overview-card .card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  position: static;
  padding-bottom: 0;
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
  padding-left: 14px;
  text-align: auto;
  font-size: 28px;
  font-weight: 950;
  color: var(--text-primary);
  letter-spacing: 0.8px;
  text-shadow: 0 2px 12px rgba(126,182,255,0.08), 0 1px 2px rgba(179,170,255,0.06);
  transition: transform 0.18s cubic-bezier(.4,2,.6,1);
  word-break: break-all;
  line-height: 1.1;
}
.overview-card:hover .main-amount {
  transform: scale(1.06);
}

.overview-card .unit {
  font-size: 16px;
  font-weight: 700;
  color: var(--brand-primary);
  margin-left: 0.6rem;
  opacity: 1;
}

/* Overview Title 樣式，與 distribution-title 一致 */
.overview-title {
  font-size: 1.5rem;
  font-weight: 800;
  color: var(--brand-primary);
  margin-bottom: 18px;
  margin-left: 4px;
  letter-spacing: 1px;
  line-height: 1.18;
  text-shadow: 0 2px 8px rgba(126,182,255,0.08), 0 1px 2px rgba(179,170,255,0.06);
  display: inline-block;
  background: linear-gradient(90deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.overview-title-block {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 20px;
}

.overview-title-icon {
  color: var(--brand-primary);
}

.overview-title-text {
  font-size: 1.8rem;
  font-weight: 900;
  color: var(--brand-primary);
  margin-left: 4px;
  letter-spacing: 2px;
  line-height: 1.13;
  text-shadow: 0 4px 18px rgba(126,182,255,0.13), 0 1.5px 6px rgba(179,170,255,0.10);
  display: inline-block;
  background: linear-gradient(90deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-transform: uppercase;
  font-stretch: expanded;
  filter: brightness(1.08) drop-shadow(0 2px 8px rgba(126,182,255,0.10));
}

.delegator-filter-wrap {
  margin-left: 12px;
  display: flex;
  align-items: center;
}

.delegator-multiselect {
  min-width: 120px;
  padding: 6px 12px;
  border-radius: 8px;
  border: 1.5px solid var(--brand-primary);
  background: var(--bg-card);
  color: var(--brand-primary);
  font-size: 15px;
  font-weight: 700;
  outline: none;
  box-shadow: 0 1px 4px rgba(126,182,255,0.08);
}

.delegator-multiselect:focus {
  border: 1.5px solid var(--brand-secondary);
  box-shadow: 0 0 0 2px rgba(126,182,255,0.15);
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



.arrow {
  animation: arrowMove 1.5s infinite alternate ease-in-out;
}

@keyframes arrowMove {
  0% { transform: translateY(0); }
  100% { transform: translateY(-1px); }
}

@keyframes change-breath {
  0% { filter: brightness(1); }
  50% { filter: brightness(1.35); }
  100% { filter: brightness(1); }
}
.overview-card .change-badge .change-amount {
  animation: change-breath 1.8s infinite cubic-bezier(.4,2,.6,1);
}

/* Responsive Design */
@media (max-width: 1200px) {
  .overview-section {
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
  

  

  .overview-card .main-amount {
    font-size: 32px;
  }
}

@media (max-width: 768px) {
  .overview-card {
    padding: 20px 16px;
  }
  
  .overview-card .main-amount {
    font-size: 28px;
  }

}

@media (max-width: 480px) {

  
  .overview-card .main-amount {
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
  font-weight: 900;
  color: var(--brand-primary);
  margin-bottom: 18px;
  letter-spacing: 2px;
  line-height: 1.13;
  text-shadow: 0 4px 18px rgba(126,182,255,0.13), 0 1.5px 6px rgba(179,170,255,0.10);
  background: linear-gradient(90deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-transform: uppercase;
  font-stretch: expanded;
  filter: brightness(1.08) drop-shadow(0 2px 8px rgba(126,182,255,0.10));
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
.reward-log-card {
  background: var(--bg-card);
  border-radius: 18px;
  box-shadow: 0 4px 24px rgba(59,130,246,0.10), 0 1.5px 6px rgba(0,0,0,0.08);
  padding: 22px 24px 16px 24px;
  margin-top: 36px;
  margin-bottom: 36px;
  border: 1.5px solid rgba(0,0,0,0.10);
  max-width: 1400px;
  min-width: 320px;
  width: 100%;
}
.log-toolbar {
  display: flex;
  gap: 18px;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 18px;
  justify-content: space-between;
}
.log-week-select {
  padding: 7px 14px;
  border-radius: 8px;
  border: 1.5px solid var(--brand-primary);
  background: var(--bg-card);
  color: var(--brand-primary);
  font-size: 15px;
  min-width: 120px;
  font-weight: 700;
  transition: border 0.2s;
}
.log-week-select:focus {
  outline: none;
  border: 1.5px solid var(--brand-secondary);
  box-shadow: 0 0 0 2px rgba(126,182,255,0.15);
}
.log-search {
  padding: 7px 14px;
  border-radius: 8px;
  border: 1.5px solid var(--brand-primary);
  background: var(--bg-card);
  color: var(--text-primary);
  font-size: 15px;
  min-width: 200px;
  transition: border 0.2s;
}
.log-search:focus {
  outline: none;
  border: 1.5px solid var(--brand-secondary);
  box-shadow: 0 0 0 2px rgba(126,182,255,0.15);
}
.log-export-btn {
  padding: 8px 22px;
  border-radius: 8px;
  background: linear-gradient(90deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  color: #fff;
  font-size: 15px;
  font-weight: 700;
  border: none;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(126,182,255,0.10);
  transition: background 0.18s, box-shadow 0.18s, transform 0.12s;
  margin-left: auto;
  letter-spacing: 1px;
}
.log-export-btn:hover {
  background: linear-gradient(90deg, var(--brand-secondary) 0%, var(--brand-primary) 100%);
  box-shadow: 0 4px 16px rgba(126,182,255,0.18);
  transform: translateY(-2px) scale(1.04);
}
.log-table-wrap {
  width: 100%;
  overflow-x: auto;
  border-radius: 12px;
  background: var(--bg-card);
  box-shadow: 0 2px 12px rgba(126,182,255,0.06);
}
.reward-log-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  background: var(--bg-card);
  color: var(--text-primary);
  font-size: 15px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1.5px 6px rgba(59,130,246,0.08);
}
.reward-log-table thead tr {
  background: var(--bg-card);
}
.log-export-btn {
  /* ...原有樣式... */
  opacity: 0.85; /* 你可以調整這個值，0.7~0.9都很自然 */
}
.reward-log-table th {
  padding: 13px 10px;
  font-size: 15px;
  font-weight: 700;
  color: var(--brand-primary);
  letter-spacing: 0.5px;
  text-align: left;
  cursor: pointer;
  user-select: none;
  border-bottom: 2.5px solid var(--brand-secondary);
  transition: background 0.2s;
  background: var(--bg-card);
}
.reward-log-table th:last-child {
  border-top-right-radius: 12px;
}
.reward-log-table th:first-child {
  border-top-left-radius: 12px;
}
.reward-log-table tbody tr {
  background: var(--bg-card);
  transition: background 0.18s;
}
.reward-log-table tbody tr:hover {
  background: rgba(126,182,255,0.10);
}
.reward-log-table td {
  padding: 12px 10px;
  border-bottom: 1.2px solid #23263a;
  font-size: 15px;
  font-weight: 500;
  vertical-align: middle;
  transition: background 0.18s;
}
.reward-log-table td:last-child {
  word-break: break-all;
  font-size: 13px;
  color: var(--text-secondary);
}
.log-amount-pos {
  color: var(--success);
  font-weight: 700;
  font-size: 13px;
  word-break: break-all;
  line-height: 1.3;
}
.log-amount-neg {
  color: var(--danger);
  font-weight: 700;
  font-size: 13px;
  word-break: break-all;
  line-height: 1.3;
}
.log-amount-zero {
  color: #b0b4ba;
  font-weight: 600;
  font-size: 13px;
  word-break: break-all;
  line-height: 1.3;
}
.log-hash {
  font-family: 'Fira Mono', 'Consolas', monospace;
  font-size: 13px;
  color: var(--brand-secondary);
}
.log-empty {
  text-align: center;
  color: #8a8f98;
  font-size: 16px;
  font-style: italic;
  background: transparent;
  padding: 24px 0;
}
@media (max-width: 900px) {
  .reward-log-card {
    padding: 12px 4vw 10px 4vw;
  }
  .log-table-wrap {
    border-radius: 8px;
  }
  .reward-log-table th, .reward-log-table td {
    font-size: 14px;
    padding: 8px 6px;
  }
  .overview-card .main-amount {
    font-size: 24px;
  }
  .overview-card .unit {
    font-size: 14px;
  }
  .staked-amount {
    font-size: 14px;
  }
}
@media (max-width: 600px) {
  .reward-log-card {
    padding: 6px 2vw 6px 2vw;
  }
  .log-toolbar {
    flex-direction: column;
    gap: 8px;
    align-items: stretch;
  }
  .reward-log-table th, .reward-log-table td {
    font-size: 13px;
    padding: 6px 4px;
  }
  .overview-card .main-amount {
    font-size: 20px;
  }
  .overview-card .unit {
    font-size: 12px;
  }
  .staked-amount {
    font-size: 12px;
  }
  .overview-card .change-badge {
    font-size: 10px;
    max-width: 100px;
  }
}
.log-expand-wrap {
  display: flex;
  justify-content: center;
  margin-top: 10px;
}
.log-expand-btn {
  padding: 8px 32px;
  border-radius: 8px;
  background: linear-gradient(90deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  color: #fff;
  font-size: 15px;
  font-weight: 700;
  border: none;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(126,182,255,0.10);
  transition: background 0.18s, box-shadow 0.18s, transform 0.12s;
  letter-spacing: 1px;
  margin: 0 auto;
  opacity: 1;
}
.log-expand-btn:hover {
  background: linear-gradient(90deg, var(--brand-secondary) 0%, var(--brand-primary) 100%);
  box-shadow: 0 4px 16px rgba(126,182,255,0.18);
  transform: translateY(-2px) scale(1.04);
  opacity: 1;
}
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.25s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}
.log-fade-mask {
  position: absolute;
  left: 0; right: 0; bottom: 0;
  height: 48px;
  pointer-events: none;
  background: linear-gradient(180deg, rgba(36,40,54,0) 0%, var(--bg-card) 90%);
  border-radius: 0 0 12px 12px;
  z-index: 2;
}
.overview-next-btn {
  height: 48px;
  min-width: 48px;
  border-radius: 50%;
  background: linear-gradient(90deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  color: #fff;
  font-size: 20px;
  font-weight: 700;
  border: none;
  margin-left: 16px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(126,182,255,0.10);
  transition: background 0.18s, box-shadow 0.18s, transform 0.12s;
}
.overview-next-btn:hover {
  background: linear-gradient(90deg, var(--brand-secondary) 0%, var(--brand-primary) 100%);
  box-shadow: 0 4px 16px rgba(126,182,255,0.18);
  transform: translateY(-2px) scale(1.04);
}

.msb-wrap {
  position: relative;
  display: inline-block;
}
.msb-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: var(--bg-card);
  border: 1.5px solid var(--brand-primary);
  border-radius: 8px;
  color: var(--brand-primary);
  font-size: 15px;
  font-weight: 700;
  padding: 7px 18px 7px 12px;
  cursor: pointer;
  box-shadow: 0 1px 4px rgba(126,182,255,0.08);
  transition: border 0.18s, box-shadow 0.18s;
  min-width: 120px;
}
.msb-btn:focus, .msb-btn:hover {
  border: 1.5px solid var(--brand-secondary);
  box-shadow: 0 0 0 2px rgba(126,182,255,0.15);
}
.msb-placeholder {
  color: #b0b4ba;
  font-weight: 500;
}
.msb-selected {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}
.msb-tag {
  background: linear-gradient(90deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  color: #fff;
  border-radius: 6px;
  padding: 2px 8px;
  font-size: 13px;
  font-weight: 700;
  margin-right: 2px;
}
.msb-arrow {
  margin-left: 4px;
  color: var(--brand-primary);
  transition: transform 0.18s;
}
.msb-btn[aria-expanded="true"] .msb-arrow {
  transform: rotate(180deg);
}
.msb-dropdown {
  position: absolute;
  top: 110%;
  left: 0;
  min-width: 160px;
  background: var(--bg-card);
  border: 1.5px solid var(--brand-primary);
  border-radius: 10px;
  box-shadow: 0 4px 24px rgba(126,182,255,0.13);
  z-index: 20;
  padding: 8px 0;
  margin-top: 4px;
}
.msb-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 7px 18px 7px 14px;
  font-size: 15px;
  color: var(--brand-primary);
  cursor: pointer;
  transition: background 0.15s;
}
.msb-option:hover {
  background: rgba(126,182,255,0.08);
}
.msb-option input[type="checkbox"] {
  accent-color: var(--brand-primary);
  margin-right: 2px;
}
.delegator-filter-btn {
  background: linear-gradient(90deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  color: #fff;
  font-size: 15px;
  font-weight: 700;
  border: none;
  border-radius: 8px;
  padding: 8px 22px;
  margin-left: 12px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(126,182,255,0.10);
  transition: background 0.18s, box-shadow 0.18s, transform 0.12s;
}
.delegator-filter-btn:hover {
  background: linear-gradient(90deg, var(--brand-secondary) 0%, var(--brand-primary) 100%);
  box-shadow: 0 4px 16px rgba(126,182,255,0.18);
  transform: translateY(-2px) scale(1.04);
}
.delegator-filter-popwrap {
  position: relative;
  display: inline-block;
}
.delegator-popover {
  position: absolute;
  top: 110%;
  left: 0;
  min-width: 160px;
  background: var(--bg-card);
  border: 1.5px solid var(--brand-primary);
  border-radius: 10px;
  box-shadow: 0 4px 24px rgba(126,182,255,0.13);
  z-index: 20;
  padding: 8px 0;
  margin-top: 4px;
}
.delegator-pop-option {
  padding: 7px 18px 7px 14px;
  font-size: 15px;
  color: var(--brand-primary);
  cursor: pointer;
  transition: background 0.15s, color 0.15s, font-weight 0.15s;
  display: flex;
  align-items: center;
  border-radius: 8px;
}
.delegator-pop-option.selected {
  /* 不再高亮，只用勾勾 */
}
.delegator-pop-option:hover {
  background: rgba(126,182,255,0.13);
}
.delegator-pop-check {
  margin-left: 12px;
  display: flex;
  align-items: center;
}
.delegator-pop-label {
  font-weight: 700;
  letter-spacing: 0.5px;
}
.delegator-filter-btn-pro {
  display: flex;
  align-items: center;
  gap: 8px;
  background: linear-gradient(90deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  color: #fff;
  font-size: 15px;
  font-weight: 800;
  letter-spacing: 1.2px;
  border: none;
  border-radius: 14px;
  padding: 7px 20px;
  margin-left: 12px;
  cursor: pointer;
  box-shadow: 0 2px 12px rgba(126,182,255,0.13);
  transition: background 0.18s, box-shadow 0.18s, transform 0.12s, filter 0.18s;
  text-transform: uppercase;
  position: relative;
  overflow: hidden;
}
.delegator-filter-btn-pro:hover {
  background: linear-gradient(90deg, var(--brand-secondary) 0%, var(--brand-primary) 100%);
  box-shadow: 0 6px 24px rgba(126,182,255,0.18);
  transform: translateY(-2px) scale(1.045);
  filter: brightness(1.08);
}
.filter-icon {
  color: #fff;
  opacity: 0.92;
  margin-right: 2px;
}
.staked-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  color: rgba(59,130,246,0.55);
  font-weight: 500;
  letter-spacing: 0.2px;
  margin-bottom: 2px;
  margin-top: 7px;
  justify-content: flex-start;
  animation: change-breath 1.8s infinite cubic-bezier(.4,2,.6,1);
}
.staked-icon {
  vertical-align: middle;
  margin-right: 2px;
  opacity: 0.7;
  margin-left: 10px;
  
}
.staked-label {
  font-weight: 600;
  margin-right: 4px;
  color: rgba(59,130,246,0.55);

  letter-spacing: 0.5px;
}
.staked-amount {
  font-weight: 700;
  color: rgba(59,130,246,0.75);
  font-size: 16px;
  margin-left: 2px;
  word-break: break-all;
  line-height: 1.2;
}
/* CSS: 新增 divider 與 staked-row 專業樣式 */
.overview-staked-divider {
  width: 100%;
  height: 1.5px;
  background: linear-gradient(90deg, rgba(126,182,255,0.10) 0%, rgba(59,130,246,0.13) 50%, rgba(126,182,255,0.10) 100%);
  border-radius: 2px;
  margin: 18px 0 10px 0;
}
/* CSS: overview-title-section 與 overview-card-title */
.overview-title-section {
  display: flex;
  align-items: center;
  gap: 10px;
  padding-left: 8px;
  margin-top: 8px;
}
.overview-card-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  white-space: nowrap;
  letter-spacing: 0.5px;
}
</style> 