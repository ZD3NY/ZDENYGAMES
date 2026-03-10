import { boot } from 'quasar/wrappers';
import axios, { type AxiosInstance, type AxiosError } from 'axios';

declare module 'vue' {
  interface ComponentCustomProperties {
    $axios: AxiosInstance;
    $api: AxiosInstance;
  }
}

const api = axios.create({ baseURL: '/api' });

export default boot(({ app, router }) => {
  app.config.globalProperties.$axios = axios;
  app.config.globalProperties.$api = api;

  // Response interceptor: on 401, try to refresh the access token
  api.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const originalRequest = error.config as typeof error.config & { _retry?: boolean };

      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;

        const refreshToken = localStorage.getItem('refreshToken');
        if (!refreshToken) {
          localStorage.removeItem('accessToken');
          localStorage.removeItem('refreshToken');
          localStorage.removeItem('user');
          await router.push('/auth/sign-in');
          return Promise.reject(error);
        }

        try {
          const { data } = await axios.post<{ accessToken: string; refreshToken: string }>(
            '/api/auth/refresh',
            { refreshToken },
          );

          localStorage.setItem('accessToken', data.accessToken);
          localStorage.setItem('refreshToken', data.refreshToken);
          api.defaults.headers.common['Authorization'] = `Bearer ${data.accessToken}`;

          // Update the auth store if it is already initialized
          try {
            const { useAuthStore } = await import('stores/auth');
            const authStore = useAuthStore();
            authStore.setAccessToken(data.accessToken, data.refreshToken);
          } catch {
            // store not yet initialized — localStorage update is enough
          }

          if (originalRequest.headers) {
            originalRequest.headers['Authorization'] = `Bearer ${data.accessToken}`;
          }
          return api(originalRequest);
        } catch {
          localStorage.removeItem('accessToken');
          localStorage.removeItem('refreshToken');
          localStorage.removeItem('user');
          await router.push('/auth/sign-in');
          return Promise.reject(error);
        }
      }

      return Promise.reject(error);
    },
  );
});

export { api };
