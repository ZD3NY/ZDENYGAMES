<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" />
        <q-toolbar-title>App</q-toolbar-title>
        <div class="text-body2 text-grey-4 q-mr-md gt-xs">{{ user?.email }}</div>
        <q-btn flat dense round icon="logout" aria-label="Sign out" :loading="signingOut" @click="onSignOut" />
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>
      <q-list>
        <q-item-label header>Navigation</q-item-label>
        <q-item clickable v-ripple to="/" exact>
          <q-item-section avatar>
            <q-icon name="home" />
          </q-item-section>
          <q-item-section>Dashboard</q-item-section>
        </q-item>
        <q-item clickable v-ripple tag="a" href="/tetris/" target="_blank">
          <q-item-section avatar>
            <q-icon name="videogame_asset" />
          </q-item-section>
          <q-item-section>Tetris</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useAuthStore } from 'stores/auth';

const router = useRouter();
const authStore = useAuthStore();
const { user } = storeToRefs(authStore);

const leftDrawerOpen = ref(false);
const signingOut = ref(false);

function toggleLeftDrawer() {
  leftDrawerOpen.value = !leftDrawerOpen.value;
}

async function onSignOut() {
  signingOut.value = true;
  await authStore.signOut();
  await router.push('/auth/sign-in');
  signingOut.value = false;
}
</script>
