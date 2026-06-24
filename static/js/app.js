/* =========================================================
   app.js — QuestDo (api.js + hud.js + auth.js +
            dashboard.js + tasks.js + badges.js)
   ========================================================= */

/* ── api.js ─────────────────────────────────────────────── */

const TOKEN_KEY = 'questly_token';

const Auth = {
  getToken()       { return localStorage.getItem(TOKEN_KEY); },
  setToken(token)  {
    localStorage.setItem(TOKEN_KEY, token);
    document.cookie = `token=${token}; path=/; max-age=${60 * 60 * 24}`;
  },
  clear()          {
    localStorage.removeItem(TOKEN_KEY);
    document.cookie = 'token=; path=/; max-age=0';
  },
  isLoggedIn()     { return !!Auth.getToken(); },
};

function requireAuth() {
  if (!Auth.isLoggedIn()) window.location.href = '/login';
}

async function apiFetch(path, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) };
  const token = Auth.getToken();
  if (token) headers.Authorization = `Bearer ${token}`;

  let res;
  try {
    res = await fetch(path, { ...options, headers });
  } catch (err) {
    throw new Error('Gagal menghubungi server. Pastikan backend Go sudah jalan.');
  }

  let body = null;
  try { body = await res.json(); } catch (_) {}

  if (res.status === 401) {
    Auth.clear();
    if (!window.location.pathname.endsWith('/login')) window.location.href = '/login';
  }

  if (!res.ok) throw new Error(body?.message || `Terjadi kesalahan (${res.status})`);
  return body;
}

function showToast(message, type = 'success') {
  let root = document.getElementById('toast-root');
  if (!root) {
    root = document.createElement('div');
    root.id = 'toast-root';
    root.className = 'toast-root';
    document.body.appendChild(root);
  }
  const toast = document.createElement('div');
  toast.className = `toast toast-${type}`;
  toast.textContent = message;
  root.appendChild(toast);
  setTimeout(() => toast.remove(), 3500);
}

function formatDeadline(dateStr) {
  if (!dateStr) return null;
  const d = new Date(dateStr);
  if (Number.isNaN(d.getTime())) return null;
  const today = new Date();
  const diffDays = Math.ceil((d.setHours(0,0,0,0) - today.setHours(0,0,0,0)) / 86400000);
  const label = d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' });
  if (diffDays < 0)  return { label, text: `Lewat ${label}`, overdue: true };
  if (diffDays === 0) return { label, text: 'Hari ini', overdue: false };
  return { label, text: `H-${diffDays} (${label})`, overdue: false };
}

function escapeHtml(str) {
  const div = document.createElement('div');
  div.textContent = str ?? '';
  return div.innerHTML;
}

// Konversi "YYYY-MM-DD" dari <input type="date"> ke format RFC3339
// yang bisa dibaca Go. Pakai T00:00:00 biar tidak geser akibat timezone.
function toRFC3339(dateStr) {
  if (!dateStr) return null;
  return new Date(dateStr + 'T00:00:00').toISOString();
}

/* ── hud.js ─────────────────────────────────────────────── */

async function initHud(activePage) {
  const root = document.getElementById('hud-root');
  if (!root) return null;

  root.innerHTML = `
    <header class="hud-bar">
      <div class="hud-logo">QUEST<span>LY</span></div>
      <nav class="hud-nav">
        <a href="/dashboard" data-page="dashboard">Dashboard</a>
        <a href="/tasks"     data-page="tasks">Papan Quest</a>
        <a href="/badges"    data-page="badges">Badge</a>
      </nav>
      <div class="hud-rank">
        <div class="level-hex"><span id="hud-level">-</span></div>
        <div class="xp-meta">
          <span class="xp-label" id="hud-xp-label">XP --/100</span>
          <div class="xp-track"><div class="xp-fill" id="hud-xp-fill" style="width:0%"></div></div>
        </div>
      </div>
      <div class="hud-user">
        <span class="hud-email" id="hud-email">...</span>
        <div class="hud-avatar" id="hud-avatar">?</div>
        <button class="btn btn-ghost btn-sm" id="hud-logout">Keluar</button>
      </div>
    </header>
  `;

  root.querySelectorAll(`.hud-nav a[data-page="${activePage}"]`)
      .forEach(a => a.classList.add('active'));

  document.getElementById('hud-logout').addEventListener('click', () => {
    Auth.clear();
    window.location.href = '/login';
  });

  try {
    const res = await apiFetch('/api/dashboard');
    renderHudData(res.data);
    return res.data;
  } catch (err) {
    showToast(err.message, 'error');
    return null;
  }
}

function renderHudData(user) {
  const xpInLevel = user.total_xp % 100;
  document.getElementById('hud-level').textContent     = user.current_level;
  document.getElementById('hud-xp-label').textContent  = `XP ${xpInLevel}/100`;
  document.getElementById('hud-xp-fill').style.width   = `${xpInLevel}%`;
  document.getElementById('hud-email').textContent      = user.email;
  document.getElementById('hud-avatar').textContent     = (user.email || '?').charAt(0).toUpperCase();
}

/* ── auth.js ─────────────────────────────────────────────── */

(function initAuth() {
  const page = document.body.dataset.page;
  if (page !== 'login' && page !== 'register') return;

  if (Auth.isLoggedIn()) { window.location.href = '/dashboard'; return; }

  const loginForm = document.getElementById('login-form');
  if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      hideMsg();
      const email    = document.getElementById('email').value.trim();
      const password = document.getElementById('password').value;
      const btn      = loginForm.querySelector('button[type="submit"]');
      btn.disabled = true; btn.textContent = 'Memproses...';
      try {
        const res = await apiFetch('/api/auth/login', {
          method: 'POST',
          body: JSON.stringify({ email, password }),
        });
        Auth.setToken(res.token);
        window.location.href = '/dashboard';
      } catch (err) {
        showMsg(err.message, 'error');
        btn.disabled = false; btn.textContent = 'Masuk';
      }
    });
  }

  const registerForm = document.getElementById('register-form');
  if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      hideMsg();
      const email    = document.getElementById('email').value.trim();
      const password = document.getElementById('password').value;
      const confirm  = document.getElementById('confirm').value;
      const btn      = registerForm.querySelector('button[type="submit"]');
      if (password !== confirm) { showMsg('Password dan konfirmasi tidak sama.', 'error'); return; }
      btn.disabled = true; btn.textContent = 'Memproses...';
      try {
        await apiFetch('/api/auth/register', {
          method: 'POST',
          body: JSON.stringify({ email, password }),
        });
        showMsg('Akun berhasil dibuat. Mengarahkan ke halaman login...', 'success');
        setTimeout(() => { window.location.href = '/login'; }, 1200);
      } catch (err) {
        showMsg(err.message, 'error');
        btn.disabled = false; btn.textContent = 'Daftar';
      }
    });
  }

  function showMsg(text, type) {
    const box = document.getElementById('form-message');
    box.textContent = text;
    box.className   = type === 'error' ? 'form-error' : 'form-success';
    box.hidden      = false;
  }
  function hideMsg() {
    const box = document.getElementById('form-message');
    if (box) box.hidden = true;
  }
})();

/* ── dashboard.js ────────────────────────────────────────── */

(async function initDashboard() {
  if (document.body.dataset.page !== 'dashboard') return;
  requireAuth();
  const user = await initHud('dashboard');
  if (!user) return;

  renderShards(user);
  renderBadgePreview(user.badges || []);
  loadQuestPreview();
})();

function renderShards(user) {
  const xpToNext = 100 - (user.total_xp % 100);
  document.getElementById('shards-grid').innerHTML = `
    <div class="shard">
      <div class="shard-icon">⬡</div>
      <div>
        <div class="shard-value">${user.current_level}</div>
        <div class="shard-label">Level Sekarang</div>
      </div>
    </div>
    <div class="shard">
      <div class="shard-icon">✦</div>
      <div>
        <div class="shard-value">${user.total_xp}</div>
        <div class="shard-label">Total XP</div>
      </div>
    </div>
    <div class="shard">
      <div class="shard-icon">↑</div>
      <div>
        <div class="shard-value">${xpToNext}</div>
        <div class="shard-label">XP ke Level ${user.current_level + 1}</div>
      </div>
    </div>
  `;
}

function renderBadgePreview(badges) {
  const list = document.getElementById('badges-preview-list');
  if (!badges.length) {
    list.innerHTML = '<p class="empty-note">Belum ada badge. Selesaikan quest buat naik level dan membukanya.</p>';
    return;
  }
  list.innerHTML = badges.slice(-4).reverse()
    .map(b => `<div class="badge-mini" title="${escapeHtml(b.name)}">🏅</div>`)
    .join('');
}

async function loadQuestPreview() {
  const list = document.getElementById('quest-preview-list');
  list.innerHTML = '<li class="empty-note">Memuat quest...</li>';
  try {
    const res = await apiFetch('/api/tasks');
    const pending = (res.data || []).filter(t => t.status === 'pending').slice(0, 4);
    if (!pending.length) {
      list.innerHTML = '<li class="empty-note">Gak ada quest aktif. Saatnya catat quest baru!</li>';
      return;
    }
    list.innerHTML = pending.map(t => `
      <li>
        <span style="display:flex;align-items:center;gap:.6rem;min-width:0;">
          <span class="gem gem-${t.difficulty}"></span>
          <span style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">${escapeHtml(t.title)}</span>
        </span>
      </li>`).join('');
  } catch (err) {
    list.innerHTML = `<li class="empty-note">${escapeHtml(err.message)}</li>`;
  }
}

/* ── tasks.js ────────────────────────────────────────────── */

let allTasks = [];
let allCategories = [];
let activeCategoryId = 'all';

(async function initTasks() {
  if (document.body.dataset.page !== 'tasks') return;
  requireAuth();
  const user = await initHud('tasks');
  if (!user) return;
  bindModalControls();
  await loadCategories();
  await loadTasks();
})();

async function loadCategories() {
  try {
    const res = await apiFetch('/api/categories');
    allCategories = res.data || [];
    renderFilters();
    renderCategoryOptions();
  } catch (err) { showToast(err.message, 'error'); }
}

async function loadTasks() {
  const listEl = document.getElementById('quest-list');
  listEl.innerHTML = '<div class="center-loading">Memuat papan quest...</div>';
  try {
    const res = await apiFetch('/api/tasks');
    allTasks = res.data || [];
    renderTasks();
  } catch (err) {
    listEl.innerHTML = `<div class="center-loading">${escapeHtml(err.message)}</div>`;
  }
}

function renderFilters() {
  const row = document.getElementById('category-filters');
  const chips = [{ id: 'all', name: 'Semua Quest' }, ...allCategories];
  row.innerHTML = chips.map(c =>
    `<button class="chip ${activeCategoryId == c.id ? 'active' : ''}" data-cat="${c.id}">${escapeHtml(c.name)}</button>`
  ).join('') + '<button class="chip chip-add" id="add-category-chip">+ Kategori</button>';

  row.querySelectorAll('.chip[data-cat]').forEach(btn => {
    btn.addEventListener('click', () => {
      activeCategoryId = btn.dataset.cat === 'all' ? 'all' : Number(btn.dataset.cat);
      renderFilters();
      renderTasks();
    });
  });
  document.getElementById('add-category-chip').addEventListener('click', () => toggleModal('category-modal', true));
}

function renderCategoryOptions() {
  const select = document.getElementById('task-category');
  select.innerHTML = '<option value="">Tanpa kategori</option>' +
    allCategories.map(c => `<option value="${c.id}">${escapeHtml(c.name)}</option>`).join('');
}

function renderTasks() {
  const listEl = document.getElementById('quest-list');
  let tasks = [...allTasks];
  if (activeCategoryId !== 'all') tasks = tasks.filter(t => t.category_id === activeCategoryId);
  tasks.sort((a, b) => {
    if (a.status !== b.status) return a.status === 'pending' ? -1 : 1;
    return new Date(b.created_at) - new Date(a.created_at);
  });
  if (!tasks.length) {
    listEl.innerHTML = '<div class="center-loading">Belum ada quest di kategori ini. Klik "+ Catat Quest Baru" buat mulai.</div>';
    return;
  }
  listEl.innerHTML = tasks.map(renderQuestCard).join('');
  listEl.querySelectorAll('[data-complete]').forEach(btn => {
    btn.addEventListener('click', () => completeTask(btn.dataset.complete, btn));
  });
}

function renderQuestCard(t) {
  const deadline = formatDeadline(t.deadline);
  const isDone   = t.status === 'completed';
  return `
    <div class="quest-card ${isDone ? 'is-done' : ''}">
      ${isDone ? '<span class="done-stamp">SELESAI</span>' : ''}
      <span class="gem gem-${t.difficulty}"></span>
      <div class="quest-body">
        <h3 class="quest-title">${escapeHtml(t.title)}</h3>
        ${t.description ? `<p class="quest-desc">${escapeHtml(t.description)}</p>` : ''}
        <div class="quest-meta">
          <span class="tag">${labelDifficulty(t.difficulty)}</span>
          ${t.category?.name ? `<span class="tag">${escapeHtml(t.category.name)}</span>` : ''}
          ${deadline ? `<span class="tag tag-deadline">${deadline.text}</span>` : ''}
          ${!isDone ? `<span class="quest-actions"><button class="btn btn-primary btn-sm" data-complete="${t.id}">Tandai Selesai</button></span>` : ''}
        </div>
      </div>
    </div>`;
}

function labelDifficulty(d) {
  return { easy: 'Easy · +10 XP', medium: 'Medium · +20 XP', hard: 'Hard · +30 XP' }[d] || d;
}

async function completeTask(id, btn) {
  btn.disabled = true; btn.textContent = 'Memproses...';
  try {
    const res = await apiFetch(`/api/tasks/${id}/complete`, { method: 'PUT' });
    showToast(`+${res.reward.gained_xp} XP! Quest selesai.`, 'success');
    await loadTasks();
    const hudUser = await apiFetch('/api/dashboard');
    renderHudData(hudUser.data);
    if (res.reward.level_up) showLevelUp(res.reward.level);
  } catch (err) {
    showToast(err.message, 'error');
    btn.disabled = false; btn.textContent = 'Tandai Selesai';
  }
}

function showLevelUp(level) {
  document.getElementById('levelup-number').innerHTML = `<span>${level}</span>`;
  document.getElementById('levelup-text').textContent = `Selamat, kamu naik ke Level ${level}! Cek koleksi badge buat lihat hadiahnya.`;
  toggleModal('levelup-modal', true);
}

function bindModalControls() {
  document.getElementById('open-add-task').addEventListener('click', () => toggleModal('task-modal', true));
  document.getElementById('cancel-task').addEventListener('click',    () => toggleModal('task-modal', false));
  document.getElementById('cancel-category').addEventListener('click',() => toggleModal('category-modal', false));
  document.getElementById('levelup-close').addEventListener('click',  () => toggleModal('levelup-modal', false));

  document.getElementById('task-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const payload = {
      title:       document.getElementById('task-title').value.trim(),
      description: document.getElementById('task-desc').value.trim(),
      difficulty:  document.getElementById('task-difficulty').value,
      deadline:    toRFC3339(document.getElementById('task-deadline').value),
    };
    const catId = document.getElementById('task-category').value;
    if (catId) payload.category_id = Number(catId);
    const btn = e.target.querySelector('button[type="submit"]');
    btn.disabled = true;
    try {
      await apiFetch('/api/tasks', { method: 'POST', body: JSON.stringify(payload) });
      showToast('Quest baru tercatat!', 'success');
      e.target.reset();
      toggleModal('task-modal', false);
      await loadTasks();
    } catch (err) { showToast(err.message, 'error'); }
    finally { btn.disabled = false; }
  });

  document.getElementById('category-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const name = document.getElementById('category-name').value.trim();
    const btn  = e.target.querySelector('button[type="submit"]');
    btn.disabled = true;
    try {
      await apiFetch('/api/categories', { method: 'POST', body: JSON.stringify({ name }) });
      showToast('Kategori ditambahkan.', 'success');
      e.target.reset();
      toggleModal('category-modal', false);
      await loadCategories();
    } catch (err) { showToast(err.message, 'error'); }
    finally { btn.disabled = false; }
  });
}

function toggleModal(id, show) {
  document.getElementById(id).hidden = !show;
}

/* ── badges.js ───────────────────────────────────────────── */

(async function initBadges() {
  if (document.body.dataset.page !== 'badges') return;
  requireAuth();
  const user = await initHud('badges');
  if (!user) return;
  const ownedIds = new Set((user.badges || []).map(b => b.id));
  loadAllBadges(ownedIds, user.current_level);
})();

async function loadAllBadges(ownedIds, currentLevel) {
  const grid = document.getElementById('badge-grid');
  grid.innerHTML = '<div class="center-loading">Memuat koleksi badge...</div>';
  try {
    const res    = await apiFetch('/api/badges');
    const badges = (res.data || []).sort((a, b) => a.required_level - b.required_level);
    if (!badges.length) {
      grid.innerHTML = '<div class="center-loading">Belum ada badge yang terdaftar di sistem.</div>';
      return;
    }
    grid.innerHTML = badges.map(b => {
      const unlocked = ownedIds.has(b.id);
      return `
        <div class="badge-tile">
          <div class="medallion ${unlocked ? 'unlocked' : 'locked'}">${unlocked ? '🏅' : '🔒'}</div>
          <div class="badge-name">${escapeHtml(b.name)}</div>
          <div class="badge-req ${currentLevel >= b.required_level ? 'met' : ''}">Lv. ${b.required_level}</div>
        </div>`;
    }).join('');
  } catch (err) {
    grid.innerHTML = `<div class="center-loading">${escapeHtml(err.message)}</div>`;
  }
}
