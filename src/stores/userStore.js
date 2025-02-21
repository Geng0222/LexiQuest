import { defineStore } from 'pinia';

export const useUserStore = defineStore('user', {
  state: () => ({
    history: [],
    progress: {}
  }),
  actions: {
    saveProgress(category, score) {
      this.history.push({ category, score, date: new Date().toISOString() });
      this.progress[category] = score;
    }
  }
});
