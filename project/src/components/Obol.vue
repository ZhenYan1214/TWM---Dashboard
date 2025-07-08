<template>
  <div class="obol-dashboard">
    <div class="obol-placeholder">
      <h2>Obol Reward Dashboard</h2>
      <p>Coming Soon...</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Obol',
  data() {
    return {
      // Obol 相關數據將在這裡添加
    }
  },
  methods: {
    // Obol 相關方法將在這裡添加
  }
}
</script>

<style scoped>
.obol-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  background: var(--bg-card);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  margin-top: 48px;
  color: var(--text-primary);
}

.obol-placeholder h2 {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 16px;
}

.obol-placeholder p {
  color: var(--text-secondary);
  font-size: 16px;
}
</style> 