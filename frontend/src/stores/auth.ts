import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { api } from 'boot/axios';

interface User {
  id: string;
  email: string;
  name: string | null;
}

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(localStorage.getItem('accessToken'));
  const refreshToken = ref<string | null>(localStorage.getItem('refreshToken'));
  const user = ref<User | null>(JSON.parse(localStorage.getItem('user') ?? 'null'));

  const isAuthenticated = computed(() => !!accessToken.value);

  function setSession(newAccessToken: string, newRefreshToken: string, newUser: User) {
    accessToken.value = newAccessToken;
    refreshToken.value = newRefreshToken;
    user.value = newUser;
    localStorage.setItem('accessToken', newAccessToken);
    localStorage.setItem('refreshToken', newRefreshToken);
    localStorage.setItem('user', JSON.stringify(newUser));
    api.defaults.headers.common['Authorization'] = `Bearer ${newAccessToken}`;
  }

  function setAccessToken(newAccessToken: string, newRefreshToken: string) {
    accessToken.value = newAccessToken;
    refreshToken.value = newRefreshToken;
    localStorage.setItem('accessToken', newAccessToken);
    localStorage.setItem('refreshToken', newRefreshToken);
    api.defaults.headers.common['Authorization'] = `Bearer ${newAccessToken}`;
  }

  function clearSession() {
    accessToken.value = null;
    refreshToken.value = null;
    user.value = null;
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('user');
    delete api.defaults.headers.common['Authorization'];
  }

  async function signIn(email: string, password: string) {
    const { data } = await api.post<{ accessToken: string; refreshToken: string; user: User }>(
      '/auth/sign-in',
      { email, password },
    );
    setSession(data.accessToken, data.refreshToken, data.user);
  }

  async function signUp(email: string, password: string, name?: string) {
    const { data } = await api.post<{ jobId: string }>('/auth/sign-up', { email, password, name });
    return data;
  }

  async function pollSignUpStatus(jobId: string) {
    const { data } = await api.get<{
      status: 'pending' | 'done' | 'failed';
      error?: string;
      accessToken?: string;
      refreshToken?: string;
      user?: User;
    }>(`/auth/sign-up/status/${jobId}`);

    if (data.status === 'done' && data.accessToken && data.refreshToken && data.user) {
      setSession(data.accessToken, data.refreshToken, data.user);
    }

    return data;
  }

  async function signOut() {
    try {
      await api.post('/auth/sign-out', { refreshToken: refreshToken.value });
    } catch {
      // ignore errors on sign-out
    }
    clearSession();
  }

  // Restore auth header on app start
  if (accessToken.value) {
    api.defaults.headers.common['Authorization'] = `Bearer ${accessToken.value}`;
  }

  return { accessToken, refreshToken, user, isAuthenticated, signIn, signUp, pollSignUpStatus, signOut, setAccessToken, clearSession };
});
