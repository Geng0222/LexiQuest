// DictionaryFloat.vue
<template>
  <!-- 使用 transition 以便顯示/隱藏時有淡入淡出效果 -->
  <transition name="fade">
    <div class="dictionary-float" v-if="show && currentData">
      <div class="header">
        <h3>{{ currentData.word }}</h3>
        <!-- 關閉按鈕 -->
        <button class="close-btn" @click="$emit('close')">×</button>
      </div>
      <!-- 顯示音標 -->
      <p v-if="currentData.phonetic" class="phonetic">{{ currentData.phonetic }}</p>
      
      <!-- 音檔：透過 audio 標籤 + ref 來手動播放 -->
      <audio ref="audioPlayer" :src="currentData.audio" v-if="currentData.audio" />
      <el-button
        v-if="currentData.audio"
        type="primary"
        icon="el-icon-volume-up"
        class="play-btn"
        @click="playAudio"
      >
        播放音檔
      </el-button>
      
      <!-- 顯示多個定義 -->
      <div class="meanings" v-if="safeMeanings.length">
        <div v-for="(meaning, index) in safeMeanings" :key="index" class="meaning-item">
          <h4>{{ meaning.partOfSpeech }}</h4>
          <p>{{ meaning.definition }}</p>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
import { ref, watch, computed } from "vue";
export default {
  name: "DictionaryFloat",
  props: {
    // 是否顯示此浮動面板
    show: {
      type: Boolean,
      default: false,
    },
    // 要顯示的字典資料
    currentData: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const audioPlayer = ref(null);

    // 當 currentData 更新時自動播放語音（若有音檔）
    watch(
      () => props.currentData,
      (newVal) => {
        if (newVal && newVal.audio) {
          playAudio();
        }
      }
    );

    function playAudio() {
      if (props.currentData?.audio) {
        const audio = new Audio(props.currentData.audio);
        audio.play().catch(err => console.warn("🎵 無法播放音檔:", err));
      }
    }

    const safeMeanings = computed(() =>
      Array.isArray(props.currentData?.meanings) ? props.currentData.meanings : []
    );

    return { playAudio, safeMeanings, audioPlayer };
  },
};
</script>

<style scoped>
.dictionary-float {
  position: fixed;
  top: 100px;
  right: 30px;
  width: 300px;
  background-color: #fff;
  border: 1px solid #eaeaea;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border-radius: 8px;
  padding: 16px;
  z-index: 9999;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}
.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
}
.phonetic {
  font-style: italic;
  margin-bottom: 10px;
}
.meaning-item {
  margin-bottom: 10px;
}
.play-btn {
  margin-bottom: 16px;
}
/* 淡入淡出動畫 */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>
