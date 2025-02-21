//main.js
import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { createPinia } from 'pinia';
import PrimeVue from 'primevue/config';
import ToastService from 'primevue/toastservice';

import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';

// 創建 Vue 應用
const app = createApp(App);

// ✅ 註冊 Vue 插件
app.use(createPinia());
app.use(router);
app.use(PrimeVue);
app.use(ToastService);
app.use(ElementPlus);

// ✅ 掛載應用
app.mount('#app');

