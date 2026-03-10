<template>
  <q-page class="q-pa-md">
    <div class="row q-mb-lg items-center">
      <div class="col">
        <div class="text-h5 text-weight-bold">Dashboard</div>
        <div class="text-body2 text-grey-6">Welcome back, {{ user?.name ?? user?.email }}</div>
      </div>
    </div>

    <div class="row q-col-gutter-md q-mb-lg">
      <div class="col-12 col-sm-4">
        <q-card flat bordered>
          <q-card-section>
            <div class="row items-center no-wrap">
              <div class="col">
                <div class="text-overline text-grey-6">Account</div>
                <div class="text-h6 text-weight-bold ellipsis">{{ user?.email }}</div>
              </div>
              <div class="col-auto">
                <q-icon name="person" size="2rem" color="primary" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4">
        <q-card flat bordered>
          <q-card-section>
            <div class="row items-center no-wrap">
              <div class="col">
                <div class="text-overline text-grey-6">Status</div>
                <div class="text-h6 text-weight-bold text-positive">Active</div>
              </div>
              <div class="col-auto">
                <q-icon name="verified" size="2rem" color="positive" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4">
        <q-card flat bordered>
          <q-card-section>
            <div class="row items-center no-wrap">
              <div class="col">
                <div class="text-overline text-grey-6">Session</div>
                <div class="text-h6 text-weight-bold">Authenticated</div>
              </div>
              <div class="col-auto">
                <q-icon name="lock_open" size="2rem" color="secondary" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Tetris Leaderboard -->
    <q-card flat bordered>
      <q-card-section class="row items-center q-pb-none">
        <div class="col">
          <div class="text-subtitle1 text-weight-medium">Tetris Leaderboard</div>
          <div class="text-caption text-grey-6">Top 10 all-time scores</div>
        </div>
        <div class="col-auto">
          <q-btn flat round dense icon="refresh" color="grey-6" :loading="loadingScores" @click="fetchLeaderboard" />
          <q-btn flat no-caps label="Play" icon="videogame_asset" color="primary" size="sm" tag="a" href="/tetris/" target="_blank" class="q-ml-sm" />
        </div>
      </q-card-section>

      <q-card-section>
        <q-markup-table flat bordered>
          <thead>
            <tr>
              <th class="text-left" style="width:48px">#</th>
              <th class="text-left">Player</th>
              <th class="text-right">Score</th>
              <th class="text-right">Lines</th>
              <th class="text-right">Date</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loadingScores">
              <td colspan="5" class="text-center q-py-md">
                <q-spinner color="primary" size="24px" />
              </td>
            </tr>
            <tr v-else-if="!leaderboard.length">
              <td colspan="5" class="text-center text-grey-6 q-py-md">No scores yet — be the first to play!</td>
            </tr>
            <tr v-for="entry in leaderboard" :key="entry.rank" :class="entry.rank === 1 ? 'bg-amber-1' : ''">
              <td>
                <span v-if="entry.rank === 1">🥇</span>
                <span v-else-if="entry.rank === 2">🥈</span>
                <span v-else-if="entry.rank === 3">🥉</span>
                <span v-else class="text-grey-6">{{ entry.rank }}</span>
              </td>
              <td class="text-weight-medium">{{ entry.name }}</td>
              <td class="text-right text-weight-bold text-primary">{{ entry.score.toLocaleString() }}</td>
              <td class="text-right text-grey-7">{{ entry.lines }}</td>
              <td class="text-right text-grey-6 text-caption">{{ formatDate(entry.date) }}</td>
            </tr>
          </tbody>
        </q-markup-table>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from 'stores/auth';
import { storeToRefs } from 'pinia';
import { api } from 'boot/axios';

interface LeaderboardEntry {
  rank: number;
  name: string;
  score: number;
  lines: number;
  date: string;
}

const authStore = useAuthStore();
const { user } = storeToRefs(authStore);

const leaderboard = ref<LeaderboardEntry[]>([]);
const loadingScores = ref(false);

async function fetchLeaderboard() {
  loadingScores.value = true;
  try {
    const { data } = await api.get<LeaderboardEntry[]>('/scores/leaderboard');
    leaderboard.value = data;
  } catch {
    // silently fail
  } finally {
    loadingScores.value = false;
  }
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

onMounted(fetchLeaderboard);
</script>
