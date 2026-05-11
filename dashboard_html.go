package main

const dashboardHTML = `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<title>KIMPO Chatbot Dashboard</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
<style>
  :root {
    --bg: #0f172a;
    --card: #1e293b;
    --card2: #273449;
    --text: #e2e8f0;
    --muted: #94a3b8;
    --primary: #22d3ee;
    --accent: #a78bfa;
    --success: #10b981;
    --danger: #ef4444;
    --warning: #f59e0b;
    --border: #334155;
  }
  * { box-sizing: border-box; }
  body {
    margin: 0;
    font-family: 'Segoe UI', system-ui, -apple-system, sans-serif;
    background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
    color: var(--text);
    min-height: 100vh;
  }
  header {
    padding: 24px 32px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid var(--border);
    background: rgba(15, 23, 42, 0.8);
    backdrop-filter: blur(10px);
    position: sticky;
    top: 0;
    z-index: 10;
  }
  header h1 {
    margin: 0;
    font-size: 22px;
    background: linear-gradient(90deg, var(--primary), var(--accent));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
  .status {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    padding: 8px 16px;
    border-radius: 999px;
    background: var(--card);
    border: 1px solid var(--border);
  }
  .dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: var(--danger);
    box-shadow: 0 0 10px currentColor;
  }
  .dot.on { background: var(--success); }
  main {
    padding: 24px 32px;
    max-width: 1400px;
    margin: 0 auto;
  }
  .grid {
    display: grid;
    gap: 20px;
  }
  .stats {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    margin-bottom: 24px;
  }
  .card {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 14px;
    padding: 20px;
    transition: transform 0.2s, box-shadow 0.2s;
  }
  .card:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 25px rgba(0,0,0,0.3);
  }
  .card h3 {
    margin: 0 0 8px 0;
    font-size: 13px;
    font-weight: 500;
    color: var(--muted);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  .card .value {
    font-size: 32px;
    font-weight: 700;
    color: var(--primary);
  }
  .card.accent .value { color: var(--accent); }
  .card.success .value { color: var(--success); }
  .card.warning .value { color: var(--warning); }
  .card.danger .value { color: var(--danger); }
  .card .sub {
    font-size: 12px;
    color: var(--muted);
    margin-top: 4px;
  }

  .two-col {
    grid-template-columns: 2fr 1fr;
  }
  @media (max-width: 900px) {
    .two-col { grid-template-columns: 1fr; }
  }

  .panel {
    background: var(--card);
    border: 1px solid var(--border);
    border-radius: 14px;
    padding: 20px;
  }
  .panel h2 {
    margin: 0 0 16px 0;
    font-size: 16px;
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .panel h2::before {
    content: '';
    width: 4px;
    height: 18px;
    background: linear-gradient(180deg, var(--primary), var(--accent));
    border-radius: 2px;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-size: 13px;
  }
  th, td {
    text-align: left;
    padding: 10px 8px;
    border-bottom: 1px solid var(--border);
  }
  th {
    color: var(--muted);
    font-weight: 500;
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }
  tr:last-child td { border-bottom: none; }

  .badge {
    display: inline-block;
    padding: 2px 8px;
    border-radius: 6px;
    font-size: 11px;
    font-weight: 600;
  }
  .badge.in  { background: rgba(34, 211, 238, 0.15); color: var(--primary); }
  .badge.out { background: rgba(167, 139, 250, 0.15); color: var(--accent); }

  .msg-list {
    max-height: 400px;
    overflow-y: auto;
  }
  .msg-list::-webkit-scrollbar { width: 6px; }
  .msg-list::-webkit-scrollbar-thumb { background: var(--border); border-radius: 3px; }

  .qr-box {
    text-align: center;
  }
  .qr-box img {
    width: 240px;
    height: 240px;
    background: white;
    padding: 8px;
    border-radius: 12px;
  }

  form.send-form {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  input, textarea {
    background: var(--card2);
    border: 1px solid var(--border);
    color: var(--text);
    padding: 10px 12px;
    border-radius: 8px;
    font-size: 14px;
    font-family: inherit;
    transition: border 0.2s;
  }
  input:focus, textarea:focus {
    outline: none;
    border-color: var(--primary);
  }
  textarea { min-height: 80px; resize: vertical; }
  button {
    background: linear-gradient(90deg, var(--primary), var(--accent));
    color: #0f172a;
    border: none;
    padding: 10px 16px;
    border-radius: 8px;
    font-weight: 600;
    cursor: pointer;
    transition: opacity 0.2s;
  }
  button:hover { opacity: 0.9; }
  button:disabled { opacity: 0.5; cursor: not-allowed; }

  .alert {
    padding: 10px 14px;
    border-radius: 8px;
    font-size: 13px;
    margin-top: 8px;
  }
  .alert.success { background: rgba(16, 185, 129, 0.15); color: var(--success); }
  .alert.error   { background: rgba(239, 68, 68, 0.15); color: var(--danger); }

  .chart-wrap { position: relative; height: 260px; }

  .top-users li {
    display: flex;
    justify-content: space-between;
    padding: 8px 0;
    border-bottom: 1px solid var(--border);
    font-size: 13px;
  }
  .top-users li:last-child { border-bottom: none; }
  .top-users { list-style: none; padding: 0; margin: 0; }
  .top-users .jid { color: var(--muted); font-family: monospace; font-size: 12px; }
</style>
</head>
<body>
<header>
  <h1>🤖 KIMPO Chatbot Dashboard</h1>
  <div class="status">
    <span class="dot" id="statusDot"></span>
    <span id="statusText">Memeriksa...</span>
    <span style="color: var(--muted);">|</span>
    <span id="uptime" style="color: var(--muted);">--</span>
  </div>
</header>

<main>
  <!-- Statistik Utama -->
  <div class="grid stats">
    <div class="card">
      <h3>Total Pesan</h3>
      <div class="value" id="totalMessages">0</div>
      <div class="sub">Sepanjang sesi aktif</div>
    </div>
    <div class="card accent">
      <h3>Pesan Masuk</h3>
      <div class="value" id="incomingMessages">0</div>
      <div class="sub">Dari pengguna</div>
    </div>
    <div class="card success">
      <h3>Pesan Terkirim</h3>
      <div class="value" id="outgoingMessages">0</div>
      <div class="sub">Balasan bot</div>
    </div>
    <div class="card warning">
      <h3>Pengguna Unik</h3>
      <div class="value" id="uniqueUsers">0</div>
      <div class="sub">Nomor berbeda</div>
    </div>
    <div class="card danger">
      <h3>Error</h3>
      <div class="value" id="errorCount">0</div>
      <div class="sub">Gagal kirim / API</div>
    </div>
  </div>

  <!-- Chart + QR -->
  <div class="grid two-col" style="margin-bottom: 24px;">
    <div class="panel">
      <h2>Statistik 7 Hari Terakhir</h2>
      <div class="chart-wrap"><canvas id="dailyChart"></canvas></div>
    </div>
    <div class="panel qr-box" id="qrPanel">
      <h2>Koneksi WhatsApp</h2>
      <div id="qrContent">
        <p style="color: var(--muted);">Memuat status koneksi...</p>
      </div>
    </div>
  </div>

  <!-- Hourly + Top Users -->
  <div class="grid two-col" style="margin-bottom: 24px;">
    <div class="panel">
      <h2>Distribusi Pesan per Jam</h2>
      <div class="chart-wrap"><canvas id="hourlyChart"></canvas></div>
    </div>
    <div class="panel">
      <h2>Pengguna Teraktif</h2>
      <ul class="top-users" id="topUsers">
        <li style="color: var(--muted);">Belum ada data</li>
      </ul>
    </div>
  </div>

  <!-- Kirim + Log -->
  <div class="grid two-col">
    <div class="panel">
      <h2>Log Pesan Terakhir</h2>
      <div class="msg-list">
        <table>
          <thead>
            <tr>
              <th>Waktu</th>
              <th>Arah</th>
              <th>Kontak</th>
              <th>Pesan</th>
            </tr>
          </thead>
          <tbody id="msgTable">
            <tr><td colspan="4" style="text-align:center; color: var(--muted);">Belum ada pesan</td></tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="panel">
      <h2>Kirim Pesan Manual</h2>
      <form class="send-form" id="sendForm">
        <input type="text" id="toInput" placeholder="Nomor (628xxx) atau JID lengkap" required />
        <textarea id="msgInput" placeholder="Tulis pesan..." required></textarea>
        <button type="submit" id="sendBtn">Kirim Pesan</button>
        <div id="sendAlert"></div>
      </form>
    </div>
  </div>
</main>

<script>
let dailyChart, hourlyChart;

function fmtTime(iso) {
  const d = new Date(iso);
  return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
}

function escapeHtml(s) {
  return (s || '').replace(/[&<>"']/g, c => ({
    '&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'
  }[c]));
}

async function loadStats() {
  try {
    const r = await fetch('/api/stats');
    const d = await r.json();

    document.getElementById('totalMessages').textContent = d.total_messages;
    document.getElementById('incomingMessages').textContent = d.incoming_messages;
    document.getElementById('outgoingMessages').textContent = d.outgoing_messages;
    document.getElementById('uniqueUsers').textContent = d.unique_users;
    document.getElementById('errorCount').textContent = d.error_count;
    document.getElementById('uptime').textContent = 'Uptime: ' + d.uptime;

    const dot = document.getElementById('statusDot');
    const text = document.getElementById('statusText');
    if (d.connected) {
      dot.classList.add('on');
      text.textContent = 'Terhubung ke WhatsApp';
    } else {
      dot.classList.remove('on');
      text.textContent = 'Tidak terhubung';
    }

    // daily chart
    const labels = d.daily_stats.map(x => x.date.slice(5));
    const values = d.daily_stats.map(x => x.count);
    if (!dailyChart) {
      dailyChart = new Chart(document.getElementById('dailyChart'), {
        type: 'line',
        data: {
          labels,
          datasets: [{
            label: 'Pesan Masuk',
            data: values,
            borderColor: '#22d3ee',
            backgroundColor: 'rgba(34, 211, 238, 0.15)',
            tension: 0.35,
            fill: true,
            pointRadius: 4,
            pointBackgroundColor: '#22d3ee'
          }]
        },
        options: chartOpts()
      });
    } else {
      dailyChart.data.labels = labels;
      dailyChart.data.datasets[0].data = values;
      dailyChart.update();
    }

    // hourly chart
    const hLabels = Array.from({length: 24}, (_, i) => i + ':00');
    if (!hourlyChart) {
      hourlyChart = new Chart(document.getElementById('hourlyChart'), {
        type: 'bar',
        data: {
          labels: hLabels,
          datasets: [{
            label: 'Pesan',
            data: d.hourly_stats,
            backgroundColor: 'rgba(167, 139, 250, 0.6)',
            borderColor: '#a78bfa',
            borderWidth: 1,
            borderRadius: 4
          }]
        },
        options: chartOpts()
      });
    } else {
      hourlyChart.data.datasets[0].data = d.hourly_stats;
      hourlyChart.update();
    }

    // top users
    const ul = document.getElementById('topUsers');
    if (d.top_users && d.top_users.length) {
      ul.innerHTML = d.top_users.map(u =>
        '<li><span class="jid">' + escapeHtml(u.jid) + '</span><strong>' + u.count + '</strong></li>'
      ).join('');
    } else {
      ul.innerHTML = '<li style="color: var(--muted);">Belum ada data</li>';
    }
  } catch (e) {
    console.error(e);
  }
}

function chartOpts() {
  return {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { labels: { color: '#e2e8f0' } }
    },
    scales: {
      x: { ticks: { color: '#94a3b8' }, grid: { color: 'rgba(148,163,184,0.1)' } },
      y: { ticks: { color: '#94a3b8' }, grid: { color: 'rgba(148,163,184,0.1)' }, beginAtZero: true }
    }
  };
}

async function loadMessages() {
  try {
    const r = await fetch('/api/messages');
    const list = await r.json();
    const tbody = document.getElementById('msgTable');
    if (!list || !list.length) {
      tbody.innerHTML = '<tr><td colspan="4" style="text-align:center; color: var(--muted);">Belum ada pesan</td></tr>';
      return;
    }
    tbody.innerHTML = list.map(m => {
      const dirBadge = m.direction === 'in'
        ? '<span class="badge in">MASUK</span>'
        : '<span class="badge out">KELUAR</span>';
      const text = m.direction === 'in'
        ? escapeHtml(m.message) + (m.reply ? '<br><em style="color: var(--muted);">↳ ' + escapeHtml(m.reply) + '</em>' : '')
        : escapeHtml(m.message);
      return '<tr><td>' + fmtTime(m.time) + '</td><td>' + dirBadge + '</td><td style="font-family:monospace;font-size:11px;">' + escapeHtml(m.from) + '</td><td>' + text + '</td></tr>';
    }).join('');
  } catch (e) {
    console.error(e);
  }
}

async function loadQR() {
  try {
    const r = await fetch('/api/qr');
    const d = await r.json();
    const box = document.getElementById('qrContent');
    if (d.connected) {
      box.innerHTML = '<div style="padding: 40px 0;"><div style="font-size:48px;">✅</div><p style="margin:12px 0 0; color: var(--success); font-weight:600;">Terhubung</p><p style="color: var(--muted); font-size:13px; margin-top:8px;">Bot siap menerima pesan</p></div>';
    } else if (d.qr_b64) {
      box.innerHTML = '<p style="color: var(--muted); font-size:13px;">Pindai QR berikut di WhatsApp:</p><img src="data:image/png;base64,' + d.qr_b64 + '" alt="QR Code" />';
    } else {
      box.innerHTML = '<div style="padding: 40px 0;"><div style="font-size:48px;">⏳</div><p style="color: var(--muted); margin-top:12px;">Menunggu QR Code...</p></div>';
    }
  } catch (e) {
    console.error(e);
  }
}

document.getElementById('sendForm').addEventListener('submit', async (e) => {
  e.preventDefault();
  const to = document.getElementById('toInput').value.trim();
  const msg = document.getElementById('msgInput').value.trim();
  const alertBox = document.getElementById('sendAlert');
  const btn = document.getElementById('sendBtn');

  btn.disabled = true;
  btn.textContent = 'Mengirim...';
  alertBox.innerHTML = '';

  try {
    const r = await fetch('/api/send', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ to, message: msg })
    });
    if (r.ok) {
      alertBox.innerHTML = '<div class="alert success">Pesan berhasil dikirim!</div>';
      document.getElementById('msgInput').value = '';
      loadMessages();
      loadStats();
    } else {
      const txt = await r.text();
      alertBox.innerHTML = '<div class="alert error">Gagal: ' + escapeHtml(txt) + '</div>';
    }
  } catch (err) {
    alertBox.innerHTML = '<div class="alert error">Error: ' + escapeHtml(err.message) + '</div>';
  } finally {
    btn.disabled = false;
    btn.textContent = 'Kirim Pesan';
    setTimeout(() => alertBox.innerHTML = '', 4000);
  }
});

function tick() {
  loadStats();
  loadMessages();
  loadQR();
}
tick();
setInterval(tick, 3000);
</script>
</body>
</html>`
