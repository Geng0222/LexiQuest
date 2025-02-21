<template>
    <div>
      <el-divider content-position="left">📊 學習進度</el-divider>
  
      <div v-if="Object.keys(progress).length > 0">
        <el-collapse accordion>
          <el-collapse-item v-for="(categoryData, category) in progress" :key="category" :title="`📂 ${category}`">
            <el-card v-for="(subCategoryData, subCategory) in categoryData" :key="subCategory" class="progress-card">
              <template #header>
                <span class="sub-category-title">📌 {{ subCategory }}</span>
              </template>
  
              <el-table :data="formatProgressData(subCategoryData)" border style="width: 100%">
                <el-table-column prop="word" label="單字" width="150"></el-table-column>
                <el-table-column prop="score" label="學習權重(初始9)"></el-table-column>
              </el-table>
            </el-card>
          </el-collapse-item>
        </el-collapse>
      </div>
      
      <el-empty v-else description="目前沒有學習進度"></el-empty>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      progress: {
        type: Object,
        required: true,
      },
    },
    methods: {
      formatProgressData(subCategoryData) {
        return Object.keys(subCategoryData).map(word => ({
          word,
          score: subCategoryData[word],
        }));
      },
    },
  };
  </script>
  
  <style scoped>
  .progress-card {
    margin-bottom: 10px;
    border-radius: 8px;
  }
  
  .sub-category-title {
    font-weight: bold;
    font-size: 16px;
  }
  </style>
  