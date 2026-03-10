<template>
  <q-page class="auth-page flex flex-center">
    <q-card class="auth-card" flat bordered>
      <q-card-section class="text-center q-pb-none">
        <div class="text-h5 text-weight-bold">Create account</div>
        <div class="text-body2 text-grey-6 q-mt-xs">Sign up for free</div>
      </q-card-section>

      <q-card-section>
        <q-form class="q-gutter-y-md" @submit.prevent="onSubmit">
          <q-input
            v-model="form.name"
            label="Name (optional)"
            outlined
            dense
            autocomplete="name"
            hide-bottom-space
          >
            <template #prepend>
              <q-icon name="person" size="sm" />
            </template>
          </q-input>

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
            autocomplete="new-password"
            :rules="[val => !!val || 'Password is required', val => val.length >= 8 || 'Minimum 8 characters']"
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

          <q-btn
            type="submit"
            label="Create account"
            color="primary"
            class="full-width"
            unelevated
            :loading="loading"
          />
        </q-form>

        <div v-if="waitingForQueue" class="q-mt-md text-center">
          <q-spinner color="primary" size="24px" />
          <div class="text-body2 text-grey-6 q-mt-sm">Creating your account...</div>
        </div>
      </q-card-section>

      <q-separator />

      <q-card-section class="text-center q-py-md">
        <span class="text-body2 text-grey-6">Already have an account? </span>
        <q-btn flat dense no-caps label="Sign in" color="primary" size="sm" to="/auth/sign-in" />
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
const waitingForQueue = ref(false);
const showPassword = ref(false);

const form = reactive({ name: '', email: '', password: '' });

function isValidEmail(val: string) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val);
}

async function onSubmit() {
  loading.value = true;
  try {
    const { jobId } = await authStore.signUp(form.email, form.password, form.name || undefined);
    loading.value = false;
    waitingForQueue.value = true;
    await pollStatus(jobId);
  } catch (err: unknown) {
    loading.value = false;
    const msg =
      (err as { response?: { data?: { message?: string } } })?.response?.data?.message ??
      'Sign-up failed. Please try again.';
    $q.notify({ type: 'negative', message: msg });
  }
}

async function pollStatus(jobId: string) {
  for (let i = 0; i < 30; i++) {
    await new Promise((r) => setTimeout(r, 1000));
    try {
      const result = await authStore.pollSignUpStatus(jobId);
      if (result.status === 'done') {
        waitingForQueue.value = false;
        $q.notify({ type: 'positive', message: 'Account created! Welcome.' });
        await router.push('/');
        return;
      }
      if (result.status === 'failed') {
        waitingForQueue.value = false;
        $q.notify({ type: 'negative', message: result.error ?? 'Sign-up failed.' });
        return;
      }
    } catch {
      // keep polling
    }
  }
  waitingForQueue.value = false;
  $q.notify({ type: 'negative', message: 'Sign-up timed out. Please try again.' });
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
