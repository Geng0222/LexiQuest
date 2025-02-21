<template>
  <div class="app-container">
    <Sidebar /> <!-- ✅ 側邊欄 -->

    <div class="content">
      <router-view /> 
    </div>
  </div>
</template>

<script>
import Sidebar from "./components/Sidebar.vue"; // 

export default {
  name: "App",
  components: { Sidebar }
};
</script>

<style>
/* 🌟 主要佈局 */
.app-container {
  display: flex;
  min-height: 100vh;
}

/* 📌 內容區 */
.content {
  flex: 1;
  padding: 40px;
  margin-left: 60px; /* 預設 Sidebar 關閉時 */
  transition: margin-left 0.3s;
}

/* 📌 Sidebar 展開時調整 */
.sidebar.expanded ~ .content {
  margin-left: 180px;
}

/* 🖥️ 全局樣式 */
body {
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
  text-align: center;
  background-color: #f5f5f5;
}
</style>
