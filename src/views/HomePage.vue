// HomePage.vue
<template>
  <el-container>
    <el-main>
      <header class="header">
        <h1>📚 LexiQuest</h1>
        <p>選擇你的題庫</p>
        <el-tag :type="isUsingAPI ? 'success' : 'danger'" size="small">
          {{ isUsingAPI ? '🔗 API模式' : '🌐 靜態模式' }}
        </el-tag>
      </header>

      <el-skeleton v-if="loading" animated :rows="3" />

      <section v-else class="categories-container">
        <TopicCategory
          v-for="category in sortedCategories"
          :key="category.category"
          :title="category.label"
          :books="category.books"
          :category="category.category"
          @open-category="openCategory"
        />
      </section>

      <!-- 漂浮視窗（模態窗口） -->
      <el-dialog
        v-model="isDialogVisible"
        :title="selectedCategory?.label"
        width="600px"
        destroy-on-close
        @closed="clearCategory"
      >
        <div class="books-container">
          <el-card
            v-for="(book, index) in selectedCategory?.books"
            :key="index"
            class="book-item"
            @click="goToWordlist(book)"
          >
            <p>{{ book }}</p>
          </el-card>
        </div>
      </el-dialog>
    </el-main>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useWordlists } from "../utils/useWordlists";
import TopicCategory from "../components/TopicCategory.vue";
import { useRouter } from "vue-router";

const { wordlistsWithLabels, loading, loadWordlists, isUsingAPI } = useWordlists();
const router = useRouter();

const isDialogVisible = ref(false);
const selectedCategory = ref(null);

onMounted(loadWordlists);

// 類別名稱排序，確保從左到右排列
const sortedCategories = computed(() => {
  return [...wordlistsWithLabels.value].sort((a, b) => a.label.localeCompare(b.label, 'zh-Hant'));
});

// 打開類別的模態窗口
const openCategory = (category) => {
  selectedCategory.value = category;
  isDialogVisible.value = true;
};

// 清除選中的類別
const clearCategory = () => {
  selectedCategory.value = null;
};

// 跳轉到單字列表
const goToWordlist = (filename) => {
  router.push(`/wordlist/${selectedCategory.value.category}/${filename}`);
  isDialogVisible.value = false;
};
</script>

<style scoped>
/* 標題區 */
.header {
  text-align: center;
  margin-bottom: 20px;
}

/* 類別列表容器 */
.categories-container {
  display: flex;
  flex-wrap: wrap; /* 自動換行 */
  gap: 20px;
  padding: 15px;
  justify-content: flex-start; /* 保持左對齊 */
  align-items: flex-start;
}

/* 書籍容器 */
.books-container {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
  padding: 10px;
}

/* 書籍卡片 */
.book-item {
  min-width: 150px;
  text-align: center;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  border-radius: 8px;
  background: white;
}

.book-item:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

/* 小螢幕適應 */
@media (max-width: 768px) {
  .categories-container {
    justify-content: center; /* 小螢幕居中 */
  }
}
</style>
