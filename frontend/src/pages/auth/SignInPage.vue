<template>
  <q-page class="auth-page flex flex-center">
    <q-card class="auth-card" flat bordered>
      <q-card-section class="text-center q-pb-none">
        <div class="text-h5 text-weight-bold">Welcome back</div>
        <div class="text-body2 text-grey-6 q-mt-xs">Sign in to your account</div>
      </q-card-section>

      <q-card-section>
        <q-form class="q-gutter-y-md" @submit.prevent="onSubmit">
          <q-input
            v-model="form.email"
            label="Email"
            type="email"
            outlined
            dense
            autocomplete="email"
            :rules="[val => !!val || 'Email is required', val => isValidEmail(val) || 'Invalid email']"
            hide-bottom-space
          >
            <template #prepend>
              <q-icon name="mail" size="sm" />
            </template>
          </q-input>

          <q-input
            v-model="form.password"
            label="Password"
            :type="showPassword ? 'text' : 'password'"
            outlined
            dense
            autocomplete="current-password"
            :rules="[val => !!val || 'Password is required']"
            hide-bottom-space
          >
            <template #prepend>
              <q-icon name="lock" size="sm" />
            </template>
            <template #append>
              <q-icon
                :name="showPassword ? 'visibility_off' : 'visibility'"
                class="cursor-pointer"
                size="sm"
                @click="showPassword = !showPassword"
              />
            </template>
          </q-input>

          <div class="flex justify-end">
            <q-btn flat dense no-caps label="Forgot password?" color="primary" size="sm" />
          </div>

          <q-btn
            type="submit"
            label="Sign in"
            color="primary"
            class="full-width"
            unelevated
            :loading="loading"
          />
        </q-form>
      </q-card-section>

      <q-separator />

      <q-card-section class="text-center q-py-md">
        <span class="text-body2 text-grey-6">Don't have an account? </span>
        <q-btn flat dense no-caps label="Sign up" color="primary" size="sm" to="/auth/sign-up" />
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { useAuthStore } from 'stores/auth';

const $q = useQuasar();
const router = useRouter();
const authStore = useAuthStore();

const loading = ref(false);
const showPassword = ref(false);

const form = reactive({
  email: '',
  password: '',
});

function isValidEmail(val: string) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val);
}

async function onSubmit() {
  loading.value = true;
  try {
    await authStore.signIn(form.email, form.password);
    await router.push('/');
  } catch {
    $q.notify({ type: 'negative', message: 'Invalid email or password.' });
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped lang="scss">
.auth-page {
  background: linear-gradient(135deg, $primary 0%, darken($primary, 15%) 100%);
  min-height: 100vh;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  border-radius: 16px;
  padding: 8px;
}
</style>
