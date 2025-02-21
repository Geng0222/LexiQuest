//TopicCategory.vue
<template>
  <div class="topic-card">
    <el-card shadow="hover" class="category-header" @click="openCategory">
      <h2>
        <el-icon><Notebook /></el-icon>
        {{ title }}
        <el-icon><ArrowDown /></el-icon>
      </h2>
    </el-card>
  </div>
</template>

<script setup>
import { Notebook, ArrowDown } from "@element-plus/icons-vue";

const props = defineProps({
  title: String,
  books: Array,
  category: String,
});

const emit = defineEmits(["open-category"]);

// 點擊後通知父組件打開類別視窗
const openCategory = () => {
  emit("open-category", {
    label: props.title,
    books: props.books,
    category: props.category
  });
};
</script>

<style scoped>
.topic-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 220px;
  margin: 10px;
}

/* 題庫分類按鈕 */
.category-header {
  cursor: pointer;
  text-align: center;
  width: 100%;
  transition: all 0.3s ease;
  border-radius: 10px;
  background: #f8f9fa;
}

.category-header:hover {
  transform: scale(1.03);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.15);
}

h2 {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 18px;
  padding: 10px;
}
</style>
