<template>
  <q-layout view="lHh Lpr lFf">
    <q-header class="tavern-header">
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menu" @click="toggleLeftDrawer" class="text-amber-4" />
        <q-toolbar-title class="tavern-title font-cinzel-deco flicker">
          ZdenyGames
        </q-toolbar-title>
        <div class="text-caption text-amber-6 q-mr-md gt-xs" style="font-family: 'Crimson Pro', serif; font-style: italic;">
          {{ user?.email }}
        </div>
        <q-btn flat dense round icon="logout" aria-label="Sign out" :loading="signingOut" @click="onSignOut" class="text-amber-5" />
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above class="drawer-cabin">
      <div class="drawer-brand q-pa-md">
        <div class="font-cinzel text-amber-5" style="font-size: 0.7rem; letter-spacing: 0.2em;">THE INN</div>
      </div>
      <q-separator class="separator-wood" />
      <q-list class="q-pt-sm">
        <q-item clickable v-ripple to="/" exact class="nav-item font-cinzel">
          <q-item-section avatar>
            <q-icon name="home" class="text-amber-6" />
          </q-item-section>
          <q-item-section class="text-amber-5 nav-label">The Commons</q-item-section>
        </q-item>
        <q-item clickable v-ripple tag="a" href="/tetris/" target="_blank" class="nav-item font-cinzel">
          <q-item-section avatar>
            <q-icon name="videogame_asset" class="text-amber-6" />
          </q-item-section>
          <q-item-section class="text-amber-5 nav-label">The Arena</q-item-section>
        </q-item>
      </q-list>

      <div class="drawer-footer q-pa-md">
        <div class="text-caption" style="font-style: italic; font-family: 'Crimson Pro', serif; font-size: 0.75rem; color: #7a5030;">
          "Beware the dark woods..."
        </div>
      </div>
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

<style scoped>
.tavern-header {
  background: #1e0e06;
  border-bottom: 2px solid #5a2e10;
  box-shadow: 0 2px 18px rgba(0, 0, 0, 0.6), 0 1px 0 rgba(220, 140, 20, 0.2);
}

.tavern-title {
  font-size: 1.2rem;
  font-weight: 700;
  color: #f0a830;
  letter-spacing: 0.08em;
}

.drawer-cabin {
  background: #1a0c05 !important;
  border-right: 2px solid #4a2010 !important;
  background-image:
    repeating-linear-gradient(
      89deg,
      transparent 0,
      transparent 3px,
      rgba(0,0,0,0.06) 3px,
      rgba(0,0,0,0.06) 6px
    ) !important;
}

.drawer-brand {
  background: rgba(40, 18, 8, 0.9);
}

.separator-wood {
  background: #4a2010;
  opacity: 1;
}

.nav-item {
  border-left: 3px solid transparent;
  transition: all 0.2s;
}

.nav-item:hover,
.nav-item.q-router-link--active {
  border-left-color: #e8960e;
  background: rgba(220, 140, 10, 0.12) !important;
}

.nav-label {
  font-size: 0.85rem;
  letter-spacing: 0.05em;
}

.drawer-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  border-top: 1px solid #3a1a08;
}

.font-cinzel      { font-family: 'Cinzel', serif; }
.font-cinzel-deco { font-family: 'Cinzel Decorative', serif; }

.flicker { animation: flicker 4s ease-in-out infinite; }

@keyframes flicker {
  0%, 100% { text-shadow: 0 0 10px #e8960e, 0 0 28px rgba(230, 140, 10, 0.6); }
  33%       { text-shadow: 0 0 7px  #d07808, 0 0 18px rgba(200, 110, 10, 0.4); }
  66%       { text-shadow: 0 0 14px #f0b020, 0 0 35px rgba(250, 170, 20, 0.7); }
}
</style>
