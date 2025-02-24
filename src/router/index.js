// index.js
import { createRouter, createWebHashHistory } from "vue-router";
import HomePage from "../views/HomePage.vue";
import QuizPage from "../views/QuizPage.vue";
import ResultPage from "../views/ResultPage.vue";
import WordlistPage from "../views/WordlistPage.vue";
import ProfilePage from "../views/ProfilePage.vue";
import ManageProgressPage from "../views/ManageProgressPage.vue";

const routes = [
  { path: "/", component: HomePage },
  { path: "/quiz/:category/:filename", component: QuizPage },
  { path: "/result/:category/:filename", component: ResultPage, props: true },
  { path: "/wordlist/:category/:filename", component: WordlistPage },
  { path: "/profile", component: ProfilePage },
  { path: "/manage-progress", component: ManageProgressPage },
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL), // 改為 hash 模式
  routes,
});

export default router;
