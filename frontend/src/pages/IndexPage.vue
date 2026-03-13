<template>
  <q-page class="page-cabin q-pa-md q-pa-sm-lg">

    <!-- Page header -->
    <div class="row q-mb-xl items-end">
      <div class="col">
        <div class="page-eyebrow font-cinzel">The Commons</div>
        <h1 class="page-heading font-cinzel-deco flicker q-my-xs">
          Traveller's Board
        </h1>
        <div class="page-sub font-crimson">
          Hail, <em>{{ user?.name ?? user?.email }}</em>. The fire still burns.
        </div>
      </div>
      <div class="col-auto">
        <div class="ember-wrap row items-center q-gutter-xs">
          <div class="ember-dot" />
          <span class="font-cinzel ember-label">Alive</span>
        </div>
      </div>
    </div>

    <!-- Info cards row -->
    <div class="row q-col-gutter-md q-mb-xl">
      <div class="col-12 col-sm-4">
        <div class="wood-card">
          <div class="wood-card__notch" />
          <div class="wood-card__label font-cinzel">Traveller</div>
          <div class="row items-center no-wrap q-mt-xs">
            <div class="col wood-card__value font-crimson ellipsis">{{ user?.email }}</div>
            <q-icon name="person" size="1.5rem" class="card-icon q-ml-sm" />
          </div>
        </div>
      </div>

      <div class="col-12 col-sm-4">
        <div class="wood-card">
          <div class="wood-card__notch" />
          <div class="wood-card__label font-cinzel">Standing</div>
          <div class="row items-center no-wrap q-mt-xs">
            <div class="col wood-card__value font-crimson text-positive-cabin">In Good Standing</div>
            <q-icon name="local_fire_department" size="1.5rem" class="card-icon-fire q-ml-sm" />
          </div>
        </div>
      </div>

      <div class="col-12 col-sm-4">
        <div class="wood-card">
          <div class="wood-card__notch" />
          <div class="wood-card__label font-cinzel">Passage</div>
          <div class="row items-center no-wrap q-mt-xs">
            <div class="col wood-card__value font-crimson">Gates Open</div>
            <q-icon name="vpn_key" size="1.5rem" class="card-icon q-ml-sm" />
          </div>
        </div>
      </div>
    </div>

    <!-- Divider with rune -->
    <div class="rune-divider q-mb-xl">
      <div class="rune-line" />
      <span class="rune-symbol font-cinzel">⚔</span>
      <div class="rune-line" />
    </div>

    <!-- Leaderboard -->
    <div class="wood-card wood-card--wide">
      <div class="wood-card__notch" />
      <div class="row items-start q-mb-lg">
        <div class="col">
          <div class="wood-card__label font-cinzel">Hall of Heroes</div>
          <div class="section-heading font-cinzel-deco flicker q-mt-xs">Tetris Leaderboard</div>
          <div class="font-crimson text-brown-4" style="font-style: italic; font-size: 1rem;">
            The ten mightiest block-stackers of the realm
          </div>
        </div>
        <div class="col-auto row items-center q-gutter-sm q-mt-xs">
          <q-btn flat round dense icon="refresh" class="text-amber-8" :loading="loadingScores" @click="fetchLeaderboard" />
          <a href="/tetris/" target="_blank" class="play-btn font-cinzel">
            <q-icon name="videogame_asset" size="14px" class="q-mr-xs" />
            Enter the Arena
          </a>
        </div>
      </div>

      <table class="ledger-table">
        <thead>
          <tr>
            <th class="font-cinzel">#</th>
            <th class="font-cinzel text-left">Name</th>
            <th class="font-cinzel text-right">Score</th>
            <th class="font-cinzel text-right">Lines</th>
            <th class="font-cinzel text-right">Date</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loadingScores">
            <td colspan="5" class="text-center q-py-lg">
              <q-spinner color="warning" size="28px" />
            </td>
          </tr>
          <tr v-else-if="!leaderboard.length">
            <td colspan="5" class="empty-row font-crimson">
              No soul has yet proven their worth — be the first.
            </td>
          </tr>
          <tr
            v-for="entry in leaderboard"
            :key="entry.rank"
            class="ledger-row"
            :class="entry.rank <= 3 ? `ledger-row--top${entry.rank}` : ''"
          >
            <td class="rank-cell font-cinzel">
              <span v-if="entry.rank === 1" class="rank-1">I</span>
              <span v-else-if="entry.rank === 2" class="rank-2">II</span>
              <span v-else-if="entry.rank === 3" class="rank-3">III</span>
              <span v-else class="rank-n">{{ toRoman(entry.rank) }}</span>
            </td>
            <td class="name-cell font-crimson">{{ entry.name }}</td>
            <td class="score-cell font-cinzel text-right">{{ entry.score.toLocaleString() }}</td>
            <td class="lines-cell font-crimson text-right">{{ entry.lines }}</td>
            <td class="date-cell font-crimson text-right">{{ formatDate(entry.date) }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Rune divider -->
    <div class="rune-divider q-mb-xl q-mt-xl">
      <div class="rune-line" />
      <span class="rune-symbol font-cinzel">🐺</span>
      <div class="rune-line" />
    </div>

    <!-- Wolfpack -->
    <div class="wood-card wood-card--wide">
      <div class="wood-card__notch" />
      <div class="row items-start q-mb-lg">
        <div class="col">
          <div class="wood-card__label font-cinzel">The Dark Wood</div>
          <div class="section-heading font-cinzel-deco flicker q-mt-xs">Wolfpack Leaderboard</div>
          <div class="font-crimson text-brown-4" style="font-style: italic; font-size: 1rem;">
            Those who survived the longest in the forest
          </div>
        </div>
        <div class="col-auto row items-center q-gutter-sm q-mt-xs">
          <q-btn flat round dense icon="refresh" class="text-amber-8" :loading="loadingWolfpack" @click="fetchWolfpackLeaderboard" />
          <a href="/wolfpack/" target="_blank" class="play-btn font-cinzel">
            <q-icon name="forest" size="14px" class="q-mr-xs" />
            Enter the Forest
          </a>
        </div>
      </div>

      <table class="ledger-table">
        <thead>
          <tr>
            <th class="font-cinzel">#</th>
            <th class="font-cinzel text-left">Name</th>
            <th class="font-cinzel text-right">Score</th>
            <th class="font-cinzel text-right">Waves</th>
            <th class="font-cinzel text-right">Date</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loadingWolfpack">
            <td colspan="5" class="text-center q-py-lg">
              <q-spinner color="warning" size="28px" />
            </td>
          </tr>
          <tr v-else-if="!wolfpackLeaderboard.length">
            <td colspan="5" class="empty-row font-crimson">
              No soul has yet survived the wolfpack — be the first.
            </td>
          </tr>
          <tr
            v-for="entry in wolfpackLeaderboard"
            :key="entry.rank"
            class="ledger-row"
            :class="entry.rank <= 3 ? `ledger-row--top${entry.rank}` : ''"
          >
            <td class="rank-cell font-cinzel">
              <span v-if="entry.rank === 1" class="rank-1">I</span>
              <span v-else-if="entry.rank === 2" class="rank-2">II</span>
              <span v-else-if="entry.rank === 3" class="rank-3">III</span>
              <span v-else class="rank-n">{{ toRoman(entry.rank) }}</span>
            </td>
            <td class="name-cell font-crimson">{{ entry.name }}</td>
            <td class="score-cell font-cinzel text-right">{{ entry.score.toLocaleString() }}</td>
            <td class="lines-cell font-crimson text-right">{{ entry.lines }}</td>
            <td class="date-cell font-crimson text-right">{{ formatDate(entry.date) }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Footer flavor text -->
    <div class="page-footer font-crimson q-mt-xl">
      <span>"Deep in the forest, something watches. Keep your fire burning."</span>
    </div>

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
    const { data } = await api.get<LeaderboardEntry[]>('/scores/tetris/leaderboard');
    leaderboard.value = data;
  } catch {
    // silently fail
  } finally {
    loadingScores.value = false;
  }
}

const wolfpackLeaderboard = ref<LeaderboardEntry[]>([]);
const loadingWolfpack = ref(false);

async function fetchWolfpackLeaderboard() {
  loadingWolfpack.value = true;
  try {
    const { data } = await api.get<LeaderboardEntry[]>('/scores/wolfpack/leaderboard');
    wolfpackLeaderboard.value = data;
  } catch {
    // silently fail
  } finally {
    loadingWolfpack.value = false;
  }
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

const ROMAN = ['', 'I','II','III','IV','V','VI','VII','VIII','IX','X'];
function toRoman(n: number) {
  return ROMAN[n] ?? String(n);
}

onMounted(() => {
  fetchLeaderboard();
  fetchWolfpackLeaderboard();
});
</script>

<style scoped>
/* Fonts */
.font-cinzel      { font-family: 'Cinzel', serif; }
.font-cinzel-deco { font-family: 'Cinzel Decorative', serif; }
.font-crimson     { font-family: 'Crimson Pro', Georgia, serif; }

/* Page */
.page-cabin {
  min-height: 100vh;
  background: transparent;
}

/* Header */
.page-eyebrow {
  font-size: 0.65rem;
  letter-spacing: 0.25em;
  color: #b07848;
  text-transform: uppercase;
}

.page-heading {
  font-size: 2rem;
  font-weight: 700;
  color: #f0b030;
  margin: 0;
  line-height: 1.2;
}

.page-sub {
  font-size: 1.15rem;
  color: #c09060;
  margin-top: 0.3rem;
}

/* Ember status */
.ember-wrap { opacity: 0.9; }

.ember-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: #ff6020;
  animation: ember 2.8s ease-in-out infinite;
}

@keyframes ember {
  0%, 100% { opacity: 1;   box-shadow: 0 0 6px #ff6020, 0 0 14px rgba(220, 80, 0, 0.7); }
  50%       { opacity: 0.6; box-shadow: 0 0 3px #d04010, 0 0  8px rgba(180, 50, 0, 0.4); }
}

.ember-label {
  font-size: 0.65rem;
  letter-spacing: 0.18em;
  color: #c07848;
  text-transform: uppercase;
}

/* Cards */
.wood-card {
  position: relative;
  background:
    repeating-linear-gradient(
      89.5deg,
      transparent 0,
      transparent 5px,
      rgba(0,0,0,0.05) 5px,
      rgba(0,0,0,0.05) 10px
    ),
    linear-gradient(165deg, #2e1608 0%, #221005 60%, #281408 100%);
  border: 1px solid #6a3018;
  border-radius: 3px;
  padding: 1.1rem 1.25rem 1rem;
  box-shadow:
    0 4px 24px rgba(0,0,0,0.5),
    inset 0 1px 0 rgba(220,150,30,0.12),
    inset 0 -1px 0 rgba(0,0,0,0.3);
  overflow: hidden;
  transition: border-color 0.25s, box-shadow 0.25s;
}

.wood-card:hover {
  border-color: #a05030;
  box-shadow:
    0 4px 28px rgba(0,0,0,0.55),
    0 0 22px rgba(220, 120, 10, 0.09),
    inset 0 1px 0 rgba(220,150,30,0.18);
}

.wood-card--wide {
  margin-top: 0;
}

/* Decorative top notch */
.wood-card__notch {
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 2px;
  background: linear-gradient(to right, transparent, rgba(230, 150, 30, 0.5), transparent);
}

.wood-card__label {
  font-size: 0.6rem;
  letter-spacing: 0.22em;
  color: #9a6030;
  text-transform: uppercase;
  margin-bottom: 0.2rem;
}

.wood-card__value {
  font-size: 1.1rem;
  color: #e0b880;
  line-height: 1.3;
}

.card-icon {
  color: #9a6838;
  opacity: 0.85;
}

.card-icon-fire {
  color: #ff6020;
  text-shadow: 0 0 8px rgba(255, 96, 32, 0.7);
  animation: ember 2.8s ease-in-out infinite;
}

.text-positive-cabin {
  color: #7acc80;
}

/* Rune divider */
.rune-divider {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.rune-line {
  flex: 1;
  height: 1px;
  background: linear-gradient(to right, transparent, #6a3018, transparent);
}

.rune-symbol {
  font-size: 1rem;
  color: #9a5828;
  letter-spacing: 0;
}

/* Leaderboard table */
.section-heading {
  font-size: 1.4rem;
  color: #f0b030;
  line-height: 1.2;
}

.ledger-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.95rem;
}

.ledger-table thead tr {
  border-bottom: 1px solid #6a3018;
}

.ledger-table thead th {
  padding: 0.4rem 0.75rem 0.6rem;
  font-size: 0.6rem;
  letter-spacing: 0.2em;
  color: #a06838;
  font-weight: 600;
  text-transform: uppercase;
}

.ledger-table tbody td {
  padding: 0.65rem 0.75rem;
  border-bottom: 1px solid rgba(90, 38, 14, 0.5);
  color: #c09060;
}

.ledger-row:hover td {
  background: rgba(150, 80, 20, 0.12);
}

.ledger-row--top1 td { background: rgba(220, 150, 10, 0.09); }
.ledger-row--top2 td { background: rgba(170, 140, 80, 0.06); }
.ledger-row--top3 td { background: rgba(150, 90, 30, 0.05); }

.rank-cell { width: 48px; }

.rank-1 {
  color: #f0b020;
  text-shadow: 0 0 8px rgba(250, 170, 20, 0.7);
  font-weight: 700;
  font-size: 0.85rem;
}

.rank-2 {
  color: #c8b888;
  font-weight: 700;
  font-size: 0.85rem;
}

.rank-3 {
  color: #c08858;
  font-weight: 700;
  font-size: 0.85rem;
}

.rank-n {
  color: #8a5830;
  font-size: 0.75rem;
}

.name-cell {
  color: #ddb878;
  font-size: 1rem;
}

.score-cell {
  color: #f0a830;
  text-shadow: 0 0 6px rgba(230, 150, 20, 0.45);
  font-size: 0.85rem;
  letter-spacing: 0.08em;
}

.lines-cell { color: #a07848; font-size: 0.95rem; }
.date-cell  { color: #8a6038; font-size: 0.85rem; font-style: italic; }

.empty-row {
  text-align: center;
  padding: 2rem;
  color: #8a5830;
  font-style: italic;
  font-size: 1rem;
}

/* Play button */
.play-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.35rem 0.9rem;
  font-size: 0.75rem;
  letter-spacing: 0.1em;
  color: #e8960e;
  border: 1px solid #7a4020;
  border-radius: 2px;
  text-decoration: none;
  background: rgba(60, 25, 8, 0.7);
  transition: all 0.2s;
}

.play-btn:hover {
  background: rgba(100, 45, 12, 0.85);
  border-color: #e8960e;
  color: #f0b030;
  text-shadow: 0 0 8px rgba(230, 140, 10, 0.55);
}

/* Flicker animation */
.flicker {
  animation: flicker 4s ease-in-out infinite;
}

@keyframes flicker {
  0%, 100% { text-shadow: 0 0 10px #e8960e, 0 0 28px rgba(230, 140, 10, 0.6); }
  33%       { text-shadow: 0 0 7px  #d07808, 0 0 18px rgba(200, 110, 10, 0.4); }
  66%       { text-shadow: 0 0 14px #f0b020, 0 0 35px rgba(250, 170, 20, 0.7); }
}

/* Footer */
.page-footer {
  text-align: center;
  font-style: italic;
  font-size: 0.95rem;
  color: #7a5030;
  padding-bottom: 2rem;
}
</style>
