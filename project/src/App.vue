<template>
  <div id="app" :class="{ 'dark-mode': isDarkMode }">
    <!-- Main Container -->
    <div class="dashboard-container">
      <!-- Header -->
      <header class="dashboard-header">
        <div class="dashboard-title-row">
          <div class="dashboard-title">
            {{ currentDashboard === 'flow' ? 'Flow Reward Dashboard' : 'Obol Reward Dashboard' }}
          </div>
          <div class="dashboard-switcher">
            <button
              :class="['switch-btn', { active: currentDashboard === 'flow' }]"
              @click="currentDashboard = 'flow'"
            >Flow</button>
            <button
              :class="['switch-btn', { active: currentDashboard === 'obol' }]"
              @click="currentDashboard = 'obol'"
            >Obol</button>
          </div>
        </div>
        <div class="theme-toggle" @click="toggleTheme">
          <div class="toggle-switch" :class="{ 'active': isDarkMode }">
            <div class="toggle-icon">
              <svg v-if="isDarkMode" width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"/>
              </svg>
              <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                <path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/>
              </svg>
            </div>
          </div>
        </div>
      </header>

      <!-- Dashboard Content -->
      <Flow v-if="currentDashboard === 'flow'" />
      <Obol v-else />
    </div>
  </div>
</template>

<script>
import Flow from './components/Flow.vue'
import Obol from './components/Obol.vue'

export default {
  name: 'App',
  components: {
    Flow,
    Obol
  },
  data() {
    return {
      isDarkMode: true,
      currentDashboard: 'flow'
    }
  },
  methods: {
    toggleTheme() {
      this.isDarkMode = !this.isDarkMode;
      this.applyTheme();
    },
    applyTheme() {
      const root = document.documentElement;
      if (this.isDarkMode) {
        root.removeAttribute('data-theme');
      } else {
        root.setAttribute('data-theme', 'light');
      }
      // 保存到 localStorage
      localStorage.setItem('theme', this.isDarkMode ? 'dark' : 'light');
    }
  },
  mounted() {
    // 從 localStorage 讀取主題設置
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme) {
      this.isDarkMode = savedTheme === 'dark';
    } else {
      // 如果沒有保存的設置，使用系統偏好
      this.isDarkMode = window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    this.applyTheme();
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');

:root {
  --bg-primary: #0F1419;
  --bg-secondary: #1A1D24;
  --bg-card: #1E2228;
  --text-primary: #FFFFFF;
  --text-secondary: #B8C5D6;
  --text-muted: #64748B;
  --brand-primary: #3B82F6;
  --brand-secondary: #6366F1;
  --success: #10B981;
  --danger: #EF4444;
  --border-radius: 16px;
  --shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.1);
  --shadow-md: 0 8px 25px rgba(0, 0, 0, 0.15);
  --shadow-lg: 0 12px 40px rgba(0, 0, 0, 0.2);
}

:root[data-theme="light"] {
  --bg-primary: #F8FAFC;
  --bg-secondary: #F1F5F9;
  --bg-card: #FFFFFF;
  --text-primary: #1E293B;
  --text-secondary: #475569;
  --text-muted: #64748B;
  --brand-primary: #3B82F6;
  --brand-secondary: #6366F1;
  --success: #10B981;
  --danger: #EF4444;
  --shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 8px 25px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 12px 40px rgba(0, 0, 0, 0.15);
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body, #app {
  background: var(--bg-primary);
  color: var(--text-primary);
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  min-height: 100vh;
  overflow: hidden;
}

#app {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Main Container */
.dashboard-container {
  max-width: 1320px;
  margin: 0 auto;
  padding: 32px 48px;
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

/* Header */
.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40px;
}

.dashboard-title-row {
  display: flex;
  align-items: center;
  gap: 20px;
}

.dashboard-title {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}

.dashboard-switcher {
  display: flex;
  gap: 10px;
}

.switch-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.2s ease;
  font-size: 14px;
  font-weight: 600;
}

.switch-btn.active {
  background: linear-gradient(135deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
  color: var(--text-primary);
}

.theme-toggle {
  cursor: pointer;
  transition: all 0.3s ease;
}

.toggle-switch {
  width: 48px;
  height: 26px;
  background: var(--bg-secondary);
  border-radius: 13px;
  padding: 2px;
  transition: all 0.3s ease;
  position: relative;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.toggle-switch.active {
  background: linear-gradient(135deg, var(--brand-primary) 0%, var(--brand-secondary) 100%);
}

.toggle-icon {
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  transform: translateX(0);
  color: var(--bg-primary);
}

.toggle-switch.active .toggle-icon {
  transform: translateX(22px);
  color: var(--brand-primary);
}

/* Responsive Design */
@media (max-width: 1200px) {
  .dashboard-container {
    padding: 28px 20px;
  }
}

@media (max-width: 768px) {
  .dashboard-container {
    padding: 24px 16px;
  }
  
  .dashboard-header {
    margin-bottom: 32px;
  }
  
  .dashboard-title {
    font-size: 28px;
  }
}

@media (max-width: 480px) {
  .dashboard-container {
    padding: 20px 12px;
  }
  
  .dashboard-title {
    font-size: 24px;
  }
}
</style>
