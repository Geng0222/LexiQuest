<template>
  <div class="sidebar" :class="{ expanded: expand }" @mouseover="expand = true" @mouseleave="expand = false">
    <div v-for="item in menuItems" :key="item.name" class="sidebar-item" @click="navigate(item.route)">
      <el-icon class="icon"><component :is="item.icon" /></el-icon>
      <span v-if="expand" class="text">{{ item.name }}</span>
    </div>
  </div>
</template>

<script>
import { useRouter } from "vue-router";
import { ref } from "vue";
import { House, List, User } from "@element-plus/icons-vue";

export default {
  setup() {
    const router = useRouter();
    const expand = ref(false);

    const menuItems = [
      { name: "首頁", route: "/", icon: House },
      { name: "個人設定", route: "/profile", icon: User },
      { name: "學習進度管理", route: "/manage-progress", icon: List },
    ];

    function navigate(route) {
      router.push(route);
    }

    return { menuItems, expand, navigate };
  }
};
</script>

<style scoped>
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  width: 60px;
  height: 100vh;
  background-color: #333;
  color: white;
  display: flex;
  flex-direction: column;
  align-items: center;
  transition: width 0.3s ease-in-out;
  overflow: hidden;
  padding-top: 20px;
}

.sidebar.expanded {
  width: 180px;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 12px;
  cursor: pointer;
  width: 100%;
  transition: background 0.3s;
}

.sidebar-item:hover {
  background: rgba(255, 255, 255, 0.1);
}

.icon {
  font-size: 22px;
}

.text {
  font-size: 16px;
  white-space: nowrap;
}
</style>
