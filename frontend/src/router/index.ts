import { route } from 'quasar/wrappers';
import {
  createMemoryHistory,
  createRouter,
  createWebHashHistory,
  createWebHistory,
} from 'vue-router';
import routes from './routes';

export default route(function () {
  const createHistory = process.env.SERVER
    ? createMemoryHistory
    : process.env.VUE_ROUTER_MODE === 'history'
      ? createWebHistory
      : createWebHashHistory;

  const router = createRouter({
    scrollBehavior: () => ({ left: 0, top: 0 }),
    routes,
    history: createHistory(process.env.VUE_ROUTER_BASE),
  });

  router.beforeEach((to) => {
    const isAuthenticated = !!localStorage.getItem('accessToken');

    if (to.meta.requiresAuth && !isAuthenticated) {
      return { path: '/auth/sign-in' };
    }

    if (to.meta.guestOnly && isAuthenticated) {
      return { path: '/' };
    }
  });

  return router;
});
