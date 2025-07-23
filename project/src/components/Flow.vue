<template>
  <div class="flow-dashboard">
    <!-- Top Overview Cards -->
    <div class="overview-title-block overview-title-flex">
      <span class="overview-title-text">總 Reward 變動</span>
      <div class="delegator-filter-popwrap" @click.stop>
        <button class="delegator-filter-btn-pro" @click="showOverviewFilter = !showOverviewFilter">
          <svg class="filter-icon" width="18" height="18" viewBox="0 0 20 20" fill="none">
            <path d="M3 5h14M6 9h8M9 13h2" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"/>
          </svg>
          FILTER
        </button>
        <div v-if="showOverviewFilter" class="delegator-popover">
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
          <div class="change-badge" :class="{ 'positive': card.change > 0, 'negative': card.change < 0, 'neutral': card.change === 0 }">
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
      <!-- 輪播按鈕已移除，超過三個自動換行 -->
    </section>

    <!-- Divider -->
    <div class="divider"></div>

    <!-- Bottom Distribution Cards -->
    <section class="distribution-section">
      <div class="distribution-title-block overview-title-flex">
        <span class="distribution-title-text">本週 Reward 變動</span>
      </div>
      <div class="distribution-grid">
        <div v-for="card in filteredDistributionCards" :key="card.id" class="distribution-card" @mouseenter="addCardHover" @mouseleave="removeCardHover">
          <div class="card-header">
            <div class="card-title-section">
              <div class="card-icon">
                <svg v-if="card.id === 'node'" width="28" height="28" viewBox="0 0 28 28" fill="none">
                  <defs>
                    <linearGradient id="star-blue" x1="0" y1="0" x2="1" y2="1">
                      <stop offset="0%" stop-color="#7eb6ff"/>
                      <stop offset="100%" stop-color="#b3aaff"/>
                    </linearGradient>
                  </defs>
                  <path d="M14 3L17.09 10.26L25 11.27L19 17.14L20.18 25.02L14 21.77L7.82 25.02L9 17.14L3 11.27L10.91 10.26L14 3Z"
                    stroke="url(#star-blue)" stroke-width="2.2" fill="none"/>
                </svg>
                <svg v-else width="42" height="42" viewBox="0 0 28 28" fill="none">
                  <defs>
                    <linearGradient id="check-blue" x1="0" y1="0" x2="1" y2="1">
                      <stop offset="0%" stop-color="#7eb6ff"/>
                      <stop offset="100%" stop-color="#b3aaff"/>
                    </linearGradient>
                  </defs>
                  <circle cx="14" cy="14" r="12" fill="url(#check-blue)" opacity="0.13"/>
                  <path d="M8 15L13 20L20 9" stroke="url(#check-blue)" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </div>
              <span class="card-title">{{ card.label }}</span>
            </div>
            <div class="change-indicator" :class="{ 'positive': card.change > 0, 'negative': card.change < 0, 'neutral': card.change === 0 }">
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
        <!-- 輪播按鈕已移除，超過四個自動換行 -->
      </div>
    </section>
    <!-- Divider between distribution and summary -->
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
              <th @click="sortLog('timestamp')">Timestamp <span v-if="logSort==='timestamp'">{{ logSortDir==='asc'?'▲':'▼' }}</span></th>
              <th @click="sortLog('name')">ID <span v-if="logSort==='name'">{{ logSortDir==='asc'?'▲':'▼' }}</span></th>
              <th @click="sortLog('amount')">Reward <span v-if="logSort==='amount'">{{ logSortDir==='asc'?'▲':'▼' }}</span></th>
              <th @click="sortLog('epoch_counter')">Epoch <span v-if="logSort==='epoch_counter'">{{ logSortDir==='asc'?'▲':'▼' }}</span></th>
            </tr>
          </thead>
          <transition-group name="log-row-fade" tag="tbody">
            <tr v-for="(row, idx) in displayedLog" :key="row._logKey" :style="logFadeStyle(idx)">
              <td>{{ row.timestamp.split('T')[0] }}</td>
              <td>{{ row.type==='Node' ? 'Node' : ('Delegator #' + row.delegator_id) }}</td>
              <td :class="{'log-amount-pos': row.amount>0, 'log-amount-neg': row.amount<0, 'log-amount-zero': row.amount===0}" >{{ row.amount>0?'+':'' }}{{ formatAmount(row.amount) }}</td>
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
      distributionCards: [ // 顯示本周新增內容
        {
          id: 'node',
          label: 'Node',
          amount: 0,
          change: 123.88,
          delegatorId: -1,
          epochCounter: 0,
        },
        {
          id: 'delegator2',
          label: 'Delegator #2',
          amount: 0,
          change: 4.68,
          delegatorId: 'DEL-002',
          epochCounter: 0,
        },
        {
          id: 'delegator3',
          label: 'Delegator #3',
          amount: 0,
          change: -0.12,
          delegatorId: 'DEL-003',
          epochCounter: 0,
        },
        {
          id: 'delegator4',
          label: 'Delegator #4',
          amount: 0,
          change: 2.34,
          delegatorId: 'DEL-004',
          epochCounter: 0,
        }
      ],
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
      showDistributionFilter: false,
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
    filteredDistributionCards() {
      // 只顯示被選中的 delegator
      return this.distributionCards.filter(card => this.selectedDelegators.includes(card.delegatorId));
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

    if (thisWeekNode) {
      delegatorTotals.push({
        label: `Node 總獎勵`,
        amount: thisWeekNode.node_total,
        change: lastWeekNode ? (thisWeekNode.node_total - lastWeekNode.node_total) : 0,
        delegatorId: -1
      });
    }
    thisWeekDelegators.forEach(now => {
      const last = lastWeekDelegators.find(lw => lw.delegator_id === now.delegator_id);
      delegatorTotals.push({
        label: `Delegator #${now.delegator_id} 總獎勵`,
        amount: now.delegator_total,
        change: last ? (now.delegator_total - last.delegator_total) : 0,
        delegatorId: now.delegator_id
      });
    });
    this.allDelegatorTotals = delegatorTotals;
    // 設定所有可選 delegatorId（含 Node）
    this.allDelegatorIdList = delegatorTotals.map(card => card.delegatorId);
    // 預設只選 Node、Delegator2、Delegator3
    this.selectedDelegators = this.allDelegatorIdList.filter(id => id === -1 || id === 2 || id === 3);

    // 3. 動態產生 distributionCards，支援所有 delegator（含 Node）
    const distCards = [];
    if (thisWeekNode) {
      distCards.push({
        id: 'node',
        label: 'Node',
        amount: thisWeekNode.amount,
        change: lastWeekNode ? (thisWeekNode.amount - lastWeekNode.amount) : 0,
        delegatorId: -1,
        epochCounter: thisWeekNode.epoch_counter,
      });
    }
    thisWeekDelegators.forEach(now => {
      const last = lastWeekDelegators.find(lw => lw.delegator_id === now.delegator_id);
      distCards.push({
        id: `delegator${now.delegator_id}`,
        label: `Delegator #${now.delegator_id}`,
        amount: now.amount,
        change: last ? (now.amount - last.amount) : 0,
        delegatorId: now.delegator_id,
        epochCounter: now.epoch_counter,
      });
    });
    this.distributionCards = distCards;
    // distributionCards 也要有 delegatorId 屬性（確保是數字）
    this.distributionCards.forEach(card => {
      if (typeof card.delegatorId === 'string' && card.delegatorId.startsWith('DEL-')) {
        card.delegatorId = parseInt(card.delegatorId.replace('DEL-', ''));
      }
    });

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

.overview-card .change-badge.neutral {
  background: #23263a;
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

.overview-card .main-amount,
.distribution-card .main-amount {
  font-size: 36px;
  font-weight: 950;
  color: var(--text-primary);
  letter-spacing: 1.2px;
  text-shadow: 0 2px 12px rgba(126,182,255,0.08), 0 1px 2px rgba(179,170,255,0.06);
  transition: transform 0.18s cubic-bezier(.4,2,.6,1);
}
.overview-card:hover .main-amount,
.distribution-card:hover .main-amount {
  transform: scale(1.06);
}

.overview-card .unit {
  font-size: 18px;
  font-weight: 700;
  color: var(--brand-primary);
  margin-left: 0.8rem;
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

/* Distribution Section */
.distribution-section {
  flex: 1;
  overflow: hidden;
}

.distribution-title {
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

.distribution-title-block {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 20px;
}

.distribution-title-icon {
  color: var(--brand-primary);
}

.distribution-title-text {
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

/* 改為 grid，每行四個卡片，自動換行 */
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
  padding: 26px 22px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: stretch;
  transition: box-shadow 0.3s cubic-bezier(.4,2,.6,1), border 0.3s cubic-bezier(.4,2,.6,1);
  position: relative;
  overflow: visible;
  min-height: 200px;
  height: 200px;
  max-width: 280px;
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

.distribution-card .change-indicator.neutral {
  background: #23263a;
  color: var(--brand-primary);
  font-style: normal;
  font-weight: 700;
  letter-spacing: 0.5px;
  border: 1.5px solid var(--brand-primary);
  box-shadow: none;
}

.distribution-card .card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.distribution-card .main-amount {  
  font-size: 32px;
  font-weight: 900;
  color: var(--text-primary);
  line-height: 1.1;
  margin-bottom: -12px;  
  margin-top: 28px;
  text-align: left;
  min-height: 38px;
  display: flex;
  align-items: center;
}

.distribution-card .unit {
  margin-top: 12px;
  font-size: 16px;
  font-weight: 900;
  color: var(--brand-primary);
  margin-left: 0.8rem;
  opacity: 1;
}

.distribution-card .card-details {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  gap: 8px;
  margin-top: 36px;
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
  font-weight: 600;
  min-width: 70px;
  margin-bottom: 2px;
}

.distribution-card .delegator-badge,
.distribution-card .epoch-badge {
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  background: rgba(99, 102, 241, 0.1);
  color: var(--brand-secondary);
  border: 1px solid rgba(99, 102, 241, 0.2);
  margin-bottom: 2px;
  animation: change-breath 1.8s infinite cubic-bezier(.4,2,.6,1);
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

@keyframes change-breath {
  0% { filter: brightness(1); }
  50% { filter: brightness(1.35); }
  100% { filter: brightness(1); }
}
.overview-card .change-badge .change-amount,
.distribution-card .change-indicator .change-amount {
  animation: change-breath 1.8s infinite cubic-bezier(.4,2,.6,1);
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
  .overview-card .main-amount,
  .distribution-card .main-amount {
    font-size: 32px;
  }
}

@media (max-width: 768px) {
  .overview-card {
    padding: 20px 16px;
  }
  
  .overview-card .main-amount,
  .distribution-card .main-amount {
    font-size: 28px;
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
  
  .overview-card .main-amount,
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
}
.log-amount-neg {
  color: var(--danger);
  font-weight: 700;
}
.log-amount-zero {
  color: #b0b4ba;
  font-weight: 600;
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
.distribution-next-btn {
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
.distribution-next-btn:hover {
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
</style> 