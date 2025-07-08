const KEY = 'flow_rewards'
export function loadRewards() {
  return JSON.parse(localStorage.getItem(KEY) || '[]')
}
export function saveRewards(data) {
  localStorage.setItem(KEY, JSON.stringify(data))
} 