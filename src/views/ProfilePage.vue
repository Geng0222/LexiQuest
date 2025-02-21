// ProfilePage.vue
<template>
  <div class="profile-container">
    <h1>👤 個人設定</h1>

    <!-- 顯示 API 模式或靜態模式 -->
    <el-tag :type="isUsingAPI ? 'success' : 'danger'" size="large">
      {{ isUsingAPI ? '🔗 API模式' : '🌐 靜態模式' }}
    </el-tag>

    <div v-if="isUsingAPI && userData">
      <el-divider content-position="left">基本資料</el-divider>
      <p>📌 使用者名稱：{{ userData.username }}</p>

      <!-- ✅ 讓名稱輸入框和按鈕排版更好 -->
      <el-form class="update-name-form">
        <el-form-item label="修改名稱：" class="update-name-item">
          <el-space>
            <el-input v-model="newUsername" placeholder="輸入新名稱" clearable class="username-input" />
            <el-button type="primary" @click="updateUsername">更新名稱</el-button>
          </el-space>
        </el-form-item>
      </el-form>

      <!-- 🔹 引入學習進度元件 -->
      <LearningProgress :progress="userData.progress" />
    </div>

    <div v-else>
      <el-alert title="目前處於靜態模式，部分功能可能無法使用" type="warning" show-icon />
    </div>

    <el-divider />
    <el-button type="info" icon="el-icon-back" @click="goBack">返回首頁</el-button>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useApiStatus } from "../utils/useApiStatus";
import { fetchUserProfile, updateUserProfile } from "../utils/useProfile";
import LearningProgress from "../components/LearningProgress.vue";

export default {
  components: { LearningProgress },
  setup() {
    const { isUsingAPI, checkAPIStatus } = useApiStatus();
    const router = useRouter();
    const userData = ref({ username: "", progress: {} });
    const newUsername = ref("");

    async function loadUserProfile() {
      if (!isUsingAPI.value) return;
      try {
        const data = await fetchUserProfile();
        userData.value = data;
        newUsername.value = data.username;
      } catch (error) {
        console.error("❌ 無法載入個人資料:", error);
      }
    }

    async function updateUsername() {
      if (!isUsingAPI.value) return;
      try {
        await updateUserProfile(newUsername.value);
        userData.value.username = newUsername.value;
        alert("✅ 名稱更新成功！");
      } catch (error) {
        console.error("❌ 更新名稱失敗:", error);
      }
    }

    function goBack() {
      router.push("/");
    }

    onMounted(async () => {
      await checkAPIStatus();
      await loadUserProfile();
    });

    return {
      isUsingAPI,
      userData,
      newUsername,
      updateUsername,
      goBack,
    };
  },
};
</script>

<style scoped>
.profile-container {
  max-width: 600px;
  margin: 20px auto;
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
}

.update-name-form {
  max-width: 400px;
  margin-top: 10px;
}

.update-name-item {
  display: flex;
  align-items: center;
}

.username-input {
  flex: 1;
  min-width: 200px;
}
</style>
