import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      redirect: "/platform/knowledge-bases",
    },
    {
      path: "/knowledgeBase",
      name: "home",
      component: () => import("../views/knowledge/KnowledgeBase.vue"),
    },
    {
      path: "/platform",
      name: "Platform",
      redirect: "/platform/knowledge-bases",
      component: () => import("../views/platform/index.vue"),
      children: [
        {
          path: "tenant",
          redirect: "/platform/settings"
        },
        {
          path: "settings",
          name: "settings",
          component: () => import("../views/settings/Settings.vue"),
        },
        {
          path: "model-settings",
          name: "modelSettings",
          component: () => import("../views/settings/ModelSettingsV2.vue"),
        },
        {
          path: "knowledge-bases",
          name: "knowledgeBaseList",
          component: () => import("../views/knowledge/KnowledgeBaseList.vue"),
        },
        {
          path: "knowledge-bases/:kbId",
          name: "knowledgeBaseDetail",
          component: () => import("../views/knowledge/KnowledgeBase.vue"),
        },
        {
          path: "agents",
          name: "agentList",
          component: () => import("../views/agent/AgentList.vue"),
        },
        {
          path: "creatChat",
          name: "globalCreatChat",
          component: () => import("../views/creatChat/creatChat.vue"),
        },
        {
          path: "knowledge-bases/:kbId/creatChat",
          name: "kbCreatChat",
          component: () => import("../views/creatChat/creatChat.vue"),
        },
        {
          path: "chat/:chatid",
          name: "chat",
          component: () => import("../views/chat/index.vue"),
        },
      ],
    },
  ],
});

// 开发环境启用路由调试
if (import.meta.env.DEV) {
  import('./debug-router').then(({ enableRouterDebug }) => {
    enableRouterDebug(router);
  });
}

export default router
