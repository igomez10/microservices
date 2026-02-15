import { useState } from "react";

// ─── Mock Data ───────────────────────────────────────────────
const MOCK_USERS = [
  { id: 1, username: "igomez", first_name: "Ignacio", last_name: "Gómez", email: "ignacio@socialapp.com", created_at: "2024-11-02T10:00:00Z", bio: "Building things on the internet." },
  { id: 2, username: "marialuz", first_name: "María Luz", last_name: "Restrepo", email: "mluz@socialapp.com", created_at: "2024-12-15T08:30:00Z", bio: "Designer & storyteller." },
  { id: 3, username: "camilodev", first_name: "Camilo", last_name: "Vargas", email: "camilo@socialapp.com", created_at: "2025-01-03T14:20:00Z", bio: "Full-stack developer. Coffee enthusiast." },
  { id: 4, username: "valentina_r", first_name: "Valentina", last_name: "Ríos", email: "val@socialapp.com", created_at: "2025-01-20T09:15:00Z", bio: "Exploring the world, one city at a time." },
  { id: 5, username: "andresf", first_name: "Andrés", last_name: "Franco", email: "andres@socialapp.com", created_at: "2025-02-01T11:45:00Z", bio: "Music producer & audio engineer." },
  { id: 6, username: "luciana.m", first_name: "Luciana", last_name: "Moreno", email: "luciana@socialapp.com", created_at: "2025-02-10T16:00:00Z", bio: "Data scientist. Cat person." },
];

const MOCK_COMMENTS = [
  { id: 1, content: "Just shipped v2.0 of the API — feels good to see it all come together after weeks of work.", username: "igomez", like_count: 24, created_at: "2026-02-14T09:30:00Z" },
  { id: 2, content: "The best code is the code you don't have to write.", username: "camilodev", like_count: 41, created_at: "2026-02-14T08:15:00Z" },
  { id: 3, content: "Currently in Medellín. The weather is perfect and the coffee is even better.", username: "valentina_r", like_count: 18, created_at: "2026-02-13T22:00:00Z" },
  { id: 4, content: "New track dropping this Friday. Been working on this one for months — can't wait for you all to hear it.", username: "andresf", like_count: 33, created_at: "2026-02-13T19:45:00Z" },
  { id: 5, content: "Redesigned the onboarding flow today. Sometimes the smallest changes make the biggest difference.", username: "marialuz", like_count: 15, created_at: "2026-02-13T17:30:00Z" },
  { id: 6, content: "TIL: you can use SVD for image compression. The math is beautiful.", username: "luciana.m", like_count: 29, created_at: "2026-02-13T14:10:00Z" },
  { id: 7, content: "Debugging OAuth flows on a Friday evening. Living the dream.", username: "igomez", like_count: 52, created_at: "2026-02-12T23:20:00Z" },
  { id: 8, content: "Hot take: Go error handling is actually good once you stop fighting it.", username: "camilodev", like_count: 67, created_at: "2026-02-12T16:00:00Z" },
];

const MOCK_ROLES = [
  { id: 1, name: "admin", description: "Full system access", created_at: "2024-11-02T10:00:00Z" },
  { id: 2, name: "moderator", description: "Content moderation privileges", created_at: "2024-11-02T10:00:00Z" },
  { id: 3, name: "user", description: "Standard user access", created_at: "2024-11-02T10:00:00Z" },
  { id: 4, name: "editor", description: "Content creation and editing", created_at: "2025-01-15T10:00:00Z" },
];

const MOCK_SCOPES = [
  { id: 1, name: "socialapp.users.list", description: "List users" },
  { id: 2, name: "socialapp.users.create", description: "Create users" },
  { id: 3, name: "socialapp.users.read", description: "Read a user" },
  { id: 4, name: "socialapp.comments.create", description: "Create comments" },
  { id: 5, name: "socialapp.comments.read", description: "Read comments" },
  { id: 6, name: "socialapp.feed.read", description: "Read user feed" },
  { id: 7, name: "socialapp.follower.create", description: "Follow users" },
  { id: 8, name: "socialapp.roles.list", description: "List roles" },
  { id: 9, name: "shortly.url.create", description: "Create short URLs" },
];

const MOCK_URLS = [
  { alias: "docs", url: "https://socialapp.gomezignacio.com/docs", created_at: "2025-12-01T10:00:00Z" },
  { alias: "github", url: "https://github.com/igomez/socialapp", created_at: "2025-12-15T10:00:00Z" },
  { alias: "blog", url: "https://blog.gomezignacio.com/socialapp-launch", created_at: "2026-01-20T10:00:00Z" },
];

const CURRENT_USER = MOCK_USERS[0];

// ─── Helpers ─────────────────────────────────────────────────
function timeAgo(dateStr) {
  const now = new Date("2026-02-14T12:00:00Z");
  const d = new Date(dateStr);
  const diff = Math.floor((now - d) / 1000);
  if (diff < 60) return "just now";
  if (diff < 3600) return `${Math.floor(diff / 60)}m`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`;
  if (diff < 604800) return `${Math.floor(diff / 86400)}d`;
  return d.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

function getInitials(user) {
  return (user.first_name[0] + user.last_name[0]).toUpperCase();
}

const AVATAR_COLORS = [
  "#1a1a2e", "#16213e", "#0f3460", "#533483",
  "#4a0e4e", "#2c3333", "#395B64", "#2C3639",
];

function avatarColor(username) {
  let hash = 0;
  for (let i = 0; i < username.length; i++) hash = username.charCodeAt(i) + ((hash << 5) - hash);
  return AVATAR_COLORS[Math.abs(hash) % AVATAR_COLORS.length];
}

// ─── Icons (inline SVG) ─────────────────────────────────────
const Icon = ({ d, size = 20, stroke = "currentColor", fill = "none" }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill={fill} stroke={stroke} strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <path d={d} />
  </svg>
);

const Icons = {
  home: (p) => <Icon {...p} d="M3 9.5L12 3l9 6.5V20a1 1 0 01-1 1H4a1 1 0 01-1-1V9.5z" />,
  search: (p) => <Icon {...p} d="M21 21l-4.35-4.35M11 19a8 8 0 100-16 8 8 0 000 16z" />,
  user: (p) => <Icon {...p} d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2M12 11a4 4 0 100-8 4 4 0 000 8z" />,
  users: (p) => <Icon {...p} d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75M9 11a4 4 0 100-8 4 4 0 000 8z" />,
  settings: (p) => <svg width={p?.size || 20} height={p?.size || 20} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 01-2.83 2.83l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>,
  shield: (p) => <Icon {...p} d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />,
  link: (p) => <Icon {...p} d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71" />,
  heart: (p) => <Icon {...p} d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z" />,
  edit: (p) => <Icon {...p} d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z" />,
  plus: (p) => <Icon {...p} d="M12 5v14M5 12h14" />,
  chevron: (p) => <Icon {...p} d="M9 18l6-6-6-6" />,
  logout: (p) => <Icon {...p} d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4M16 17l5-5-5-5M21 12H9" />,
  copy: (p) => <Icon {...p} d="M20 9h-9a2 2 0 00-2 2v9a2 2 0 002 2h9a2 2 0 002-2v-9a2 2 0 00-2-2zM5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1" />,
};

// ─── Styles ──────────────────────────────────────────────────
const globalStyles = `
  @import url('https://fonts.googleapis.com/css2?family=DM+Sans:ital,opsz,wght@0,9..40,300;0,9..40,400;0,9..40,500;0,9..40,600;1,9..40,400&family=JetBrains+Mono:wght@400;500&family=Instrument+Serif:ital@0;1&display=swap');

  :root {
    --bg: #FAFAF8;
    --bg-card: #FFFFFF;
    --bg-sidebar: #F5F4F0;
    --bg-hover: #F0EFEB;
    --bg-input: #F5F4F0;
    --border: #E8E6E1;
    --border-light: #F0EFEB;
    --text-primary: #1A1A1A;
    --text-secondary: #6B6B6B;
    --text-tertiary: #9B9B9B;
    --accent: #2D2D2D;
    --accent-soft: #E8E6E1;
    --accent-blue: #4A7C8A;
    --accent-green: #5A8A6A;
    --accent-warm: #8A6A4A;
    --danger: #C45D5D;
    --radius: 10px;
    --radius-sm: 6px;
    --radius-lg: 14px;
    --shadow-sm: 0 1px 2px rgba(0,0,0,0.04);
    --shadow-md: 0 2px 8px rgba(0,0,0,0.06);
    --shadow-lg: 0 4px 20px rgba(0,0,0,0.08);
    --font-body: 'DM Sans', sans-serif;
    --font-display: 'Instrument Serif', serif;
    --font-mono: 'JetBrains Mono', monospace;
  }

  * { margin: 0; padding: 0; box-sizing: border-box; }
  
  .app-root {
    font-family: var(--font-body);
    background: var(--bg);
    color: var(--text-primary);
    min-height: 100vh;
    display: flex;
    line-height: 1.55;
    font-size: 14px;
    -webkit-font-smoothing: antialiased;
  }

  /* ── Sidebar ── */
  .sidebar {
    width: 240px;
    min-width: 240px;
    background: var(--bg-sidebar);
    border-right: 1px solid var(--border);
    display: flex;
    flex-direction: column;
    height: 100vh;
    position: sticky;
    top: 0;
  }
  .sidebar-brand {
    padding: 28px 24px 24px;
    border-bottom: 1px solid var(--border);
  }
  .sidebar-brand h1 {
    font-family: var(--font-display);
    font-size: 26px;
    font-weight: 400;
    letter-spacing: -0.02em;
    color: var(--text-primary);
  }
  .sidebar-brand span {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--text-tertiary);
    letter-spacing: 0.02em;
  }
  .sidebar-nav {
    padding: 16px 12px;
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .sidebar-section-label {
    font-size: 10px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--text-tertiary);
    padding: 16px 12px 6px;
  }
  .nav-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 12px;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all 0.15s;
    font-size: 14px;
    color: var(--text-secondary);
    border: none;
    background: none;
    width: 100%;
    text-align: left;
    font-family: var(--font-body);
  }
  .nav-item:hover { background: var(--bg-hover); color: var(--text-primary); }
  .nav-item.active { background: var(--bg-card); color: var(--text-primary); font-weight: 500; box-shadow: var(--shadow-sm); }
  .nav-item .nav-icon { opacity: 0.6; flex-shrink: 0; }
  .nav-item.active .nav-icon { opacity: 1; }

  .sidebar-user {
    padding: 16px;
    border-top: 1px solid var(--border);
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .sidebar-user-info { flex: 1; min-width: 0; }
  .sidebar-user-name { font-weight: 500; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .sidebar-user-handle { font-size: 12px; color: var(--text-tertiary); font-family: var(--font-mono); }

  /* ── Main Content ── */
  .main-content {
    flex: 1;
    max-width: 900px;
    margin: 0 auto;
    padding: 0;
    min-height: 100vh;
  }
  .page-header {
    padding: 32px 40px 24px;
    border-bottom: 1px solid var(--border-light);
    position: sticky;
    top: 0;
    background: rgba(250, 250, 248, 0.92);
    backdrop-filter: blur(12px);
    z-index: 10;
  }
  .page-title {
    font-family: var(--font-display);
    font-size: 28px;
    font-weight: 400;
    letter-spacing: -0.01em;
  }
  .page-subtitle {
    color: var(--text-secondary);
    font-size: 13px;
    margin-top: 4px;
  }
  .page-body { padding: 0 40px 40px; }

  /* ── Avatar ── */
  .avatar {
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 500;
    color: #fff;
    flex-shrink: 0;
    font-family: var(--font-body);
    letter-spacing: 0.02em;
  }
  .avatar-sm { width: 32px; height: 32px; font-size: 11px; }
  .avatar-md { width: 40px; height: 40px; font-size: 13px; }
  .avatar-lg { width: 56px; height: 56px; font-size: 18px; }
  .avatar-xl { width: 80px; height: 80px; font-size: 24px; }

  /* ── Cards ── */
  .card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    overflow: hidden;
  }
  .card-flat {
    border-left: none;
    border-right: none;
    border-radius: 0;
  }

  /* ── Comment / Post ── */
  .comment-item {
    padding: 20px 40px;
    border-bottom: 1px solid var(--border-light);
    transition: background 0.15s;
    cursor: pointer;
  }
  .comment-item:hover { background: rgba(0,0,0,0.01); }
  .comment-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
  }
  .comment-meta { flex: 1; }
  .comment-author { font-weight: 500; font-size: 14px; }
  .comment-handle {
    font-family: var(--font-mono);
    font-size: 12px;
    color: var(--text-tertiary);
    margin-left: 6px;
  }
  .comment-time { font-size: 12px; color: var(--text-tertiary); }
  .comment-body {
    font-size: 15px;
    line-height: 1.6;
    color: var(--text-primary);
    margin-left: 52px;
  }
  .comment-actions {
    display: flex;
    gap: 24px;
    margin-top: 12px;
    margin-left: 52px;
  }
  .comment-action {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    color: var(--text-tertiary);
    cursor: pointer;
    background: none;
    border: none;
    font-family: var(--font-body);
    padding: 4px 0;
    transition: color 0.15s;
  }
  .comment-action:hover { color: var(--text-primary); }
  .comment-action.liked { color: #C45D5D; }
  .comment-action.liked svg { fill: #C45D5D; stroke: #C45D5D; }

  /* ── Compose ── */
  .compose-box {
    padding: 20px 40px;
    border-bottom: 1px solid var(--border);
    display: flex;
    gap: 12px;
  }
  .compose-input {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  .compose-textarea {
    width: 100%;
    border: none;
    resize: none;
    font-family: var(--font-body);
    font-size: 15px;
    line-height: 1.6;
    background: transparent;
    color: var(--text-primary);
    outline: none;
    min-height: 48px;
  }
  .compose-textarea::placeholder { color: var(--text-tertiary); }
  .compose-footer {
    display: flex;
    justify-content: flex-end;
  }

  /* ── Buttons ── */
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 8px 18px;
    border-radius: var(--radius-sm);
    font-family: var(--font-body);
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s;
    border: 1px solid transparent;
    white-space: nowrap;
  }
  .btn-primary {
    background: var(--text-primary);
    color: var(--bg);
    border-color: var(--text-primary);
  }
  .btn-primary:hover { opacity: 0.85; }
  .btn-secondary {
    background: transparent;
    color: var(--text-primary);
    border-color: var(--border);
  }
  .btn-secondary:hover { background: var(--bg-hover); }
  .btn-ghost {
    background: transparent;
    color: var(--text-secondary);
    border: none;
    padding: 6px 12px;
  }
  .btn-ghost:hover { color: var(--text-primary); background: var(--bg-hover); }
  .btn-danger {
    background: transparent;
    color: var(--danger);
    border-color: #E8D0D0;
  }
  .btn-danger:hover { background: #FDF5F5; }
  .btn-sm { padding: 5px 12px; font-size: 12px; }
  .btn-icon { padding: 6px; border-radius: var(--radius-sm); }

  /* ── Form elements ── */
  .input-group { margin-bottom: 20px; }
  .input-label {
    display: block;
    font-size: 12px;
    font-weight: 500;
    color: var(--text-secondary);
    margin-bottom: 6px;
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }
  .input {
    width: 100%;
    padding: 10px 14px;
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    font-family: var(--font-body);
    font-size: 14px;
    background: var(--bg-input);
    color: var(--text-primary);
    outline: none;
    transition: border-color 0.15s;
  }
  .input:focus { border-color: var(--text-tertiary); }
  .input-hint { font-size: 12px; color: var(--text-tertiary); margin-top: 4px; }
  .input-mono { font-family: var(--font-mono); font-size: 13px; }

  /* ── Tables ── */
  .table-wrap {
    border: 1px solid var(--border);
    border-radius: var(--radius);
    overflow: hidden;
    background: var(--bg-card);
  }
  table {
    width: 100%;
    border-collapse: collapse;
  }
  th {
    text-align: left;
    padding: 12px 16px;
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--text-tertiary);
    background: var(--bg-sidebar);
    border-bottom: 1px solid var(--border);
  }
  td {
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-light);
    font-size: 13px;
  }
  tr:last-child td { border-bottom: none; }
  tr:hover td { background: rgba(0,0,0,0.01); }

  /* ── Badge / Tag ── */
  .badge {
    display: inline-flex;
    align-items: center;
    padding: 3px 10px;
    border-radius: 20px;
    font-size: 11px;
    font-weight: 500;
    letter-spacing: 0.02em;
  }
  .badge-default { background: var(--accent-soft); color: var(--text-secondary); }
  .badge-blue { background: #E3EEF1; color: var(--accent-blue); }
  .badge-green { background: #E3F1E8; color: var(--accent-green); }
  .badge-warm { background: #F1EBE3; color: var(--accent-warm); }

  /* ── User card row ── */
  .user-row {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 14px 0;
    border-bottom: 1px solid var(--border-light);
    cursor: pointer;
    transition: opacity 0.15s;
  }
  .user-row:last-child { border-bottom: none; }
  .user-row:hover { opacity: 0.8; }
  .user-row-info { flex: 1; min-width: 0; }
  .user-row-name { font-weight: 500; font-size: 14px; }
  .user-row-handle { font-family: var(--font-mono); font-size: 12px; color: var(--text-tertiary); }
  .user-row-bio { font-size: 13px; color: var(--text-secondary); margin-top: 2px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

  /* ── Profile ── */
  .profile-header {
    padding: 32px 0;
    display: flex;
    gap: 24px;
    align-items: flex-start;
    border-bottom: 1px solid var(--border-light);
    margin-bottom: 24px;
  }
  .profile-info { flex: 1; }
  .profile-name {
    font-family: var(--font-display);
    font-size: 32px;
    font-weight: 400;
    line-height: 1.1;
  }
  .profile-handle {
    font-family: var(--font-mono);
    font-size: 14px;
    color: var(--text-tertiary);
    margin-top: 4px;
  }
  .profile-bio {
    font-size: 15px;
    color: var(--text-secondary);
    margin-top: 12px;
    line-height: 1.5;
  }
  .profile-stats {
    display: flex;
    gap: 24px;
    margin-top: 16px;
  }
  .profile-stat {
    font-size: 13px;
    color: var(--text-secondary);
  }
  .profile-stat strong {
    color: var(--text-primary);
    font-weight: 600;
    margin-right: 4px;
  }

  /* ── Tabs ── */
  .tabs {
    display: flex;
    border-bottom: 1px solid var(--border);
    margin-bottom: 0;
  }
  .tab {
    padding: 12px 20px;
    font-size: 13px;
    font-weight: 500;
    color: var(--text-tertiary);
    cursor: pointer;
    border-bottom: 2px solid transparent;
    transition: all 0.15s;
    background: none;
    border-top: none;
    border-left: none;
    border-right: none;
    font-family: var(--font-body);
  }
  .tab:hover { color: var(--text-primary); }
  .tab.active { color: var(--text-primary); border-bottom-color: var(--text-primary); }

  /* ── Stats grid ── */
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    margin-bottom: 32px;
  }
  .stat-card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 20px;
  }
  .stat-value {
    font-family: var(--font-display);
    font-size: 36px;
    font-weight: 400;
    line-height: 1;
    margin-bottom: 4px;
  }
  .stat-label {
    font-size: 12px;
    color: var(--text-tertiary);
    text-transform: uppercase;
    letter-spacing: 0.06em;
  }

  /* ── Login page ── */
  .login-page {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background: var(--bg);
    padding: 40px;
  }
  .login-container {
    width: 100%;
    max-width: 380px;
  }
  .login-brand {
    text-align: center;
    margin-bottom: 40px;
  }
  .login-brand h1 {
    font-family: var(--font-display);
    font-size: 42px;
    font-weight: 400;
    letter-spacing: -0.02em;
  }
  .login-brand p {
    color: var(--text-tertiary);
    font-size: 14px;
    margin-top: 8px;
  }
  .login-card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--radius-lg);
    padding: 32px;
    box-shadow: var(--shadow-lg);
  }
  .login-divider {
    display: flex;
    align-items: center;
    gap: 16px;
    margin: 24px 0;
    color: var(--text-tertiary);
    font-size: 12px;
  }
  .login-divider::before,
  .login-divider::after {
    content: '';
    flex: 1;
    height: 1px;
    background: var(--border);
  }
  .login-footer {
    text-align: center;
    margin-top: 24px;
    font-size: 13px;
    color: var(--text-secondary);
  }
  .login-footer a {
    color: var(--text-primary);
    font-weight: 500;
    text-decoration: none;
    cursor: pointer;
  }
  .login-footer a:hover { text-decoration: underline; }

  /* ── Section dividers ── */
  .section { margin-bottom: 32px; }
  .section-title {
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--text-tertiary);
    margin-bottom: 16px;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--border-light);
  }

  /* ── Misc ── */
  .empty-state {
    text-align: center;
    padding: 60px 40px;
    color: var(--text-tertiary);
  }
  .empty-state-icon { font-size: 48px; margin-bottom: 16px; opacity: 0.3; }
  .empty-state-text { font-size: 15px; }
  .flex { display: flex; }
  .flex-between { display: flex; justify-content: space-between; align-items: center; }
  .gap-8 { gap: 8px; }
  .gap-12 { gap: 12px; }
  .gap-16 { gap: 16px; }
  .mt-16 { margin-top: 16px; }
  .mt-24 { margin-top: 24px; }
  .mb-24 { margin-bottom: 24px; }

  /* ── Animations ── */
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(8px); }
    to { opacity: 1; transform: translateY(0); }
  }
  .animate-in {
    animation: fadeIn 0.3s ease forwards;
  }
  .delay-1 { animation-delay: 0.05s; opacity: 0; }
  .delay-2 { animation-delay: 0.1s; opacity: 0; }
  .delay-3 { animation-delay: 0.15s; opacity: 0; }
  .delay-4 { animation-delay: 0.2s; opacity: 0; }

  /* ── Search bar ── */
  .search-bar {
    display: flex;
    align-items: center;
    gap: 10px;
    background: var(--bg-input);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 10px 16px;
    margin: 20px 0;
  }
  .search-bar input {
    flex: 1;
    border: none;
    background: none;
    font-family: var(--font-body);
    font-size: 14px;
    outline: none;
    color: var(--text-primary);
  }
  .search-bar input::placeholder { color: var(--text-tertiary); }
  .search-bar svg { color: var(--text-tertiary); }

  /* ── URL item ── */
  .url-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 16px;
    border-bottom: 1px solid var(--border-light);
  }
  .url-item:last-child { border-bottom: none; }
  .url-alias {
    font-family: var(--font-mono);
    font-size: 14px;
    font-weight: 500;
    color: var(--accent-blue);
  }
  .url-target {
    font-size: 13px;
    color: var(--text-tertiary);
    margin-top: 2px;
    max-width: 400px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .page-indicator {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-family: var(--font-mono);
    font-size: 10px;
    color: var(--text-tertiary);
    background: var(--bg-input);
    padding: 2px 8px;
    border-radius: 4px;
    margin-left: 12px;
    vertical-align: middle;
    letter-spacing: 0.04em;
  }
`;

// ─── Components ──────────────────────────────────────────────

const Avatar = ({ user, size = "md" }) => (
  <div className={`avatar avatar-${size}`} style={{ background: avatarColor(user.username) }}>
    {getInitials(user)}
  </div>
);

// ─── PAGE: Login ─────────────────────────────────────────────
const LoginPage = ({ onLogin, onSwitch }) => {
  const [isSignup, setIsSignup] = useState(false);
  return (
    <div className="login-page">
      <div className="login-container animate-in">
        <div className="login-brand">
          <h1>Socialapp</h1>
          <p>A space for thoughtful conversation</p>
        </div>
        <div className="login-card">
          <h2 style={{ fontSize: 18, fontWeight: 500, marginBottom: 24 }}>
            {isSignup ? "Create account" : "Welcome back"}
          </h2>
          {isSignup && (
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12 }}>
              <div className="input-group">
                <label className="input-label">First name</label>
                <input className="input" placeholder="John" />
              </div>
              <div className="input-group">
                <label className="input-label">Last name</label>
                <input className="input" placeholder="Doe" />
              </div>
            </div>
          )}
          {isSignup && (
            <div className="input-group">
              <label className="input-label">Username</label>
              <input className="input input-mono" placeholder="johndoe" />
            </div>
          )}
          <div className="input-group">
            <label className="input-label">Email</label>
            <input className="input" type="email" placeholder="you@example.com" />
          </div>
          <div className="input-group">
            <label className="input-label">Password</label>
            <input className="input" type="password" placeholder="••••••••" />
          </div>
          <button className="btn btn-primary" style={{ width: "100%", padding: "12px", marginTop: 8 }} onClick={onLogin}>
            {isSignup ? "Create account" : "Sign in"}
          </button>
          <div className="login-divider">or</div>
          <button className="btn btn-secondary" style={{ width: "100%", padding: "12px" }}>
            Continue with OAuth
          </button>
          <div className="login-footer">
            {isSignup ? (
              <>Already have an account? <a onClick={() => setIsSignup(false)}>Sign in</a></>
            ) : (
              <>Don't have an account? <a onClick={() => setIsSignup(true)}>Create one</a></>
            )}
          </div>
        </div>
        <p style={{ textAlign: "center", fontSize: 12, color: "var(--text-tertiary)", marginTop: 24 }}>
          POST /v1/oauth/token · POST /v1/users
        </p>
      </div>
    </div>
  );
};

// ─── PAGE: Feed ──────────────────────────────────────────────
const FeedPage = ({ onUserClick }) => {
  const [liked, setLiked] = useState({});
  const [newComment, setNewComment] = useState("");
  const toggleLike = (id) => setLiked(p => ({ ...p, [id]: !p[id] }));

  return (
    <div>
      <div className="page-header">
        <div className="flex-between">
          <div>
            <h2 className="page-title">Feed</h2>
            <p className="page-subtitle">Posts from people you follow</p>
          </div>
          <span className="page-indicator">GET /v1/feed</span>
        </div>
      </div>

      <div className="compose-box animate-in">
        <Avatar user={CURRENT_USER} size="md" />
        <div className="compose-input">
          <textarea
            className="compose-textarea"
            placeholder="What's on your mind?"
            rows={2}
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
          />
          <div className="compose-footer">
            <button className="btn btn-primary btn-sm">Post</button>
          </div>
        </div>
      </div>

      {MOCK_COMMENTS.map((c, i) => {
        const user = MOCK_USERS.find(u => u.username === c.username) || MOCK_USERS[0];
        const isLiked = liked[c.id];
        return (
          <div key={c.id} className={`comment-item animate-in delay-${Math.min(i + 1, 4)}`}>
            <div className="comment-header">
              <Avatar user={user} size="md" />
              <div className="comment-meta">
                <span className="comment-author" onClick={() => onUserClick(user)} style={{cursor:'pointer'}}>
                  {user.first_name} {user.last_name}
                </span>
                <span className="comment-handle">@{c.username}</span>
                <span style={{color:'var(--text-tertiary)', fontSize: 12}}> · {timeAgo(c.created_at)}</span>
              </div>
            </div>
            <div className="comment-body">{c.content}</div>
            <div className="comment-actions">
              <button className={`comment-action ${isLiked ? 'liked' : ''}`} onClick={() => toggleLike(c.id)}>
                <Icons.heart size={15} fill={isLiked ? "#C45D5D" : "none"} stroke={isLiked ? "#C45D5D" : "currentColor"} />
                {c.like_count + (isLiked ? 1 : 0)}
              </button>
            </div>
          </div>
        );
      })}
    </div>
  );
};

// ─── PAGE: Explore (Users) ───────────────────────────────────
const ExplorePage = ({ onUserClick }) => {
  const [search, setSearch] = useState("");
  const filtered = MOCK_USERS.filter(u =>
    u.username !== CURRENT_USER.username &&
    (u.first_name.toLowerCase().includes(search.toLowerCase()) ||
     u.last_name.toLowerCase().includes(search.toLowerCase()) ||
     u.username.toLowerCase().includes(search.toLowerCase()))
  );

  return (
    <div>
      <div className="page-header">
        <div className="flex-between">
          <div>
            <h2 className="page-title">Explore</h2>
            <p className="page-subtitle">Discover people on the network</p>
          </div>
          <span className="page-indicator">GET /v1/users</span>
        </div>
      </div>
      <div className="page-body">
        <div className="search-bar animate-in">
          <Icons.search size={18} />
          <input placeholder="Search users…" value={search} onChange={e => setSearch(e.target.value)} />
        </div>

        <div className="animate-in delay-1">
          {filtered.map(user => (
            <div key={user.id} className="user-row" onClick={() => onUserClick(user)}>
              <Avatar user={user} size="md" />
              <div className="user-row-info">
                <div className="user-row-name">{user.first_name} {user.last_name}</div>
                <div className="user-row-handle">@{user.username}</div>
                {user.bio && <div className="user-row-bio">{user.bio}</div>}
              </div>
              <button className="btn btn-secondary btn-sm" onClick={e => e.stopPropagation()}>Follow</button>
            </div>
          ))}
        </div>

        {filtered.length === 0 && (
          <div className="empty-state">
            <div className="empty-state-icon">∅</div>
            <div className="empty-state-text">No users found</div>
          </div>
        )}
      </div>
    </div>
  );
};

// ─── PAGE: Profile ───────────────────────────────────────────
const ProfilePage = ({ user, isSelf }) => {
  const userComments = MOCK_COMMENTS.filter(c => c.username === user.username);
  const [activeTab, setActiveTab] = useState("posts");

  return (
    <div>
      <div className="page-header">
        <div className="flex-between">
          <h2 className="page-title">Profile</h2>
          <span className="page-indicator">GET /v1/users/{"{username}"}</span>
        </div>
      </div>
      <div className="page-body">
        <div className="profile-header animate-in">
          <Avatar user={user} size="xl" />
          <div className="profile-info">
            <div className="flex-between">
              <div>
                <h2 className="profile-name">{user.first_name} {user.last_name}</h2>
                <div className="profile-handle">@{user.username}</div>
              </div>
              {isSelf ? (
                <button className="btn btn-secondary btn-sm"><Icons.edit size={14} /> Edit profile</button>
              ) : (
                <button className="btn btn-primary btn-sm">Follow</button>
              )}
            </div>
            {user.bio && <p className="profile-bio">{user.bio}</p>}
            <div className="profile-stats">
              <span className="profile-stat"><strong>12</strong>following</span>
              <span className="profile-stat"><strong>48</strong>followers</span>
              <span className="profile-stat"><strong>{userComments.length}</strong>posts</span>
            </div>
            <p style={{fontSize: 12, color: 'var(--text-tertiary)', marginTop: 12}}>
              Joined {new Date(user.created_at).toLocaleDateString("en-US", { month: "long", year: "numeric" })}
            </p>
          </div>
        </div>

        <div className="tabs animate-in delay-1">
          <button className={`tab ${activeTab === 'posts' ? 'active' : ''}`} onClick={() => setActiveTab('posts')}>Posts</button>
          <button className={`tab ${activeTab === 'followers' ? 'active' : ''}`} onClick={() => setActiveTab('followers')}>Followers</button>
          <button className={`tab ${activeTab === 'following' ? 'active' : ''}`} onClick={() => setActiveTab('following')}>Following</button>
          {isSelf && <button className={`tab ${activeTab === 'roles' ? 'active' : ''}`} onClick={() => setActiveTab('roles')}>Roles</button>}
        </div>

        {activeTab === 'posts' && (
          <div className="animate-in delay-2">
            {userComments.length > 0 ? userComments.map(c => (
              <div key={c.id} className="comment-item" style={{paddingLeft: 0, paddingRight: 0}}>
                <div className="comment-header">
                  <Avatar user={user} size="md" />
                  <div className="comment-meta">
                    <span className="comment-author">{user.first_name} {user.last_name}</span>
                    <span className="comment-handle">@{c.username}</span>
                    <span style={{color:'var(--text-tertiary)', fontSize: 12}}> · {timeAgo(c.created_at)}</span>
                  </div>
                </div>
                <div className="comment-body">{c.content}</div>
                <div className="comment-actions">
                  <button className="comment-action">
                    <Icons.heart size={15} /> {c.like_count}
                  </button>
                </div>
              </div>
            )) : (
              <div className="empty-state">
                <div className="empty-state-text">No posts yet</div>
              </div>
            )}
          </div>
        )}

        {activeTab === 'followers' && (
          <div className="animate-in delay-2" style={{paddingTop: 16}}>
            {MOCK_USERS.slice(1, 4).map(u => (
              <div key={u.id} className="user-row">
                <Avatar user={u} size="md" />
                <div className="user-row-info">
                  <div className="user-row-name">{u.first_name} {u.last_name}</div>
                  <div className="user-row-handle">@{u.username}</div>
                </div>
                <button className="btn btn-secondary btn-sm">Following</button>
              </div>
            ))}
          </div>
        )}

        {activeTab === 'following' && (
          <div className="animate-in delay-2" style={{paddingTop: 16}}>
            {MOCK_USERS.slice(2, 6).map(u => (
              <div key={u.id} className="user-row">
                <Avatar user={u} size="md" />
                <div className="user-row-info">
                  <div className="user-row-name">{u.first_name} {u.last_name}</div>
                  <div className="user-row-handle">@{u.username}</div>
                </div>
                <button className="btn btn-secondary btn-sm">Unfollow</button>
              </div>
            ))}
          </div>
        )}

        {activeTab === 'roles' && isSelf && (
          <div className="animate-in delay-2" style={{paddingTop: 16}}>
            <div className="flex gap-8" style={{flexWrap:'wrap'}}>
              <span className="badge badge-blue">admin</span>
              <span className="badge badge-green">editor</span>
              <span className="badge badge-default">user</span>
            </div>
            <p style={{fontSize: 12, color: 'var(--text-tertiary)', marginTop: 12}}>
              GET /v1/users/{user.username}/roles
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

// ─── PAGE: Settings ──────────────────────────────────────────
const SettingsPage = () => {
  const [activeSection, setActiveSection] = useState("profile");

  return (
    <div>
      <div className="page-header">
        <h2 className="page-title">Settings</h2>
        <p className="page-subtitle">Manage your account</p>
      </div>
      <div className="page-body">
        <div className="tabs animate-in" style={{marginBottom: 32}}>
          <button className={`tab ${activeSection === 'profile' ? 'active' : ''}`} onClick={() => setActiveSection('profile')}>Profile</button>
          <button className={`tab ${activeSection === 'password' ? 'active' : ''}`} onClick={() => setActiveSection('password')}>Password</button>
          <button className={`tab ${activeSection === 'danger' ? 'active' : ''}`} onClick={() => setActiveSection('danger')}>Danger zone</button>
        </div>

        {activeSection === 'profile' && (
          <div className="animate-in delay-1" style={{maxWidth: 480}}>
            <div className="section-title">Profile information</div>
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12 }}>
              <div className="input-group">
                <label className="input-label">First name</label>
                <input className="input" defaultValue={CURRENT_USER.first_name} />
              </div>
              <div className="input-group">
                <label className="input-label">Last name</label>
                <input className="input" defaultValue={CURRENT_USER.last_name} />
              </div>
            </div>
            <div className="input-group">
              <label className="input-label">Username</label>
              <input className="input input-mono" defaultValue={CURRENT_USER.username} />
              <p className="input-hint">socialapp.gomezignacio.com/@{CURRENT_USER.username}</p>
            </div>
            <div className="input-group">
              <label className="input-label">Email</label>
              <input className="input" type="email" defaultValue={CURRENT_USER.email} />
            </div>
            <div className="flex gap-8 mt-16">
              <button className="btn btn-primary">Save changes</button>
              <button className="btn btn-ghost">Cancel</button>
            </div>
            <p style={{fontSize: 12, color: 'var(--text-tertiary)', marginTop: 16}}>
              PUT /v1/users/{CURRENT_USER.username}
            </p>
          </div>
        )}

        {activeSection === 'password' && (
          <div className="animate-in delay-1" style={{maxWidth: 480}}>
            <div className="section-title">Change password</div>
            <div className="input-group">
              <label className="input-label">Current password</label>
              <input className="input" type="password" placeholder="••••••••" />
            </div>
            <div className="input-group">
              <label className="input-label">New password</label>
              <input className="input" type="password" placeholder="••••••••" />
            </div>
            <div className="input-group">
              <label className="input-label">Confirm new password</label>
              <input className="input" type="password" placeholder="••••••••" />
            </div>
            <button className="btn btn-primary mt-16">Update password</button>
            <p style={{fontSize: 12, color: 'var(--text-tertiary)', marginTop: 16}}>
              POST /v1/password
            </p>

            <div style={{marginTop: 40}}>
              <div className="section-title">Reset password</div>
              <p style={{fontSize: 13, color: 'var(--text-secondary)', marginBottom: 16}}>
                Forgot your password? Enter your email to receive a reset link.
              </p>
              <div className="input-group">
                <label className="input-label">Email address</label>
                <input className="input" type="email" placeholder="you@example.com" />
              </div>
              <button className="btn btn-secondary">Send reset link</button>
              <p style={{fontSize: 12, color: 'var(--text-tertiary)', marginTop: 12}}>
                PUT /v1/password
              </p>
            </div>
          </div>
        )}

        {activeSection === 'danger' && (
          <div className="animate-in delay-1" style={{maxWidth: 480}}>
            <div className="card" style={{border: '1px solid #E8D0D0', padding: 24}}>
              <h3 style={{fontSize: 16, fontWeight: 500, color: 'var(--danger)', marginBottom: 8}}>Delete account</h3>
              <p style={{fontSize: 13, color: 'var(--text-secondary)', marginBottom: 16}}>
                Once you delete your account, there is no going back. All your data, posts, and connections will be permanently removed.
              </p>
              <button className="btn btn-danger">Delete my account</button>
              <p style={{fontSize: 12, color: 'var(--text-tertiary)', marginTop: 12}}>
                DELETE /v1/users/{CURRENT_USER.username}
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

// ─── PAGE: Admin ─────────────────────────────────────────────
const AdminPage = () => {
  const [activeTab, setActiveTab] = useState("roles");

  return (
    <div>
      <div className="page-header">
        <div className="flex-between">
          <div>
            <h2 className="page-title">Admin</h2>
            <p className="page-subtitle">Roles, scopes, and URL management</p>
          </div>
        </div>
      </div>
      <div className="page-body">
        <div className="tabs animate-in" style={{marginBottom: 32}}>
          <button className={`tab ${activeTab === 'roles' ? 'active' : ''}`} onClick={() => setActiveTab('roles')}>Roles</button>
          <button className={`tab ${activeTab === 'scopes' ? 'active' : ''}`} onClick={() => setActiveTab('scopes')}>Scopes</button>
          <button className={`tab ${activeTab === 'urls' ? 'active' : ''}`} onClick={() => setActiveTab('urls')}>Short URLs</button>
        </div>

        {activeTab === 'roles' && (
          <div className="animate-in delay-1">
            <div className="flex-between mb-24">
              <span className="page-indicator" style={{marginLeft: 0}}>GET /v1/roles</span>
              <button className="btn btn-primary btn-sm"><Icons.plus size={14} /> Create role</button>
            </div>
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Description</th>
                    <th>Created</th>
                    <th style={{width: 60}}></th>
                  </tr>
                </thead>
                <tbody>
                  {MOCK_ROLES.map(r => (
                    <tr key={r.id}>
                      <td><span className="badge badge-blue">{r.name}</span></td>
                      <td style={{color: 'var(--text-secondary)'}}>{r.description}</td>
                      <td style={{fontSize: 12, color: 'var(--text-tertiary)', fontFamily: 'var(--font-mono)'}}>
                        {new Date(r.created_at).toLocaleDateString("en-US", {month: "short", day: "numeric", year: "numeric"})}
                      </td>
                      <td>
                        <button className="btn btn-ghost btn-sm" style={{padding: '4px 8px'}}>
                          <Icons.chevron size={16} />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {activeTab === 'scopes' && (
          <div className="animate-in delay-1">
            <div className="flex-between mb-24">
              <span className="page-indicator" style={{marginLeft: 0}}>GET /v1/scopes</span>
              <button className="btn btn-primary btn-sm"><Icons.plus size={14} /> Create scope</button>
            </div>
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Scope</th>
                    <th>Description</th>
                    <th style={{width: 60}}></th>
                  </tr>
                </thead>
                <tbody>
                  {MOCK_SCOPES.map(s => (
                    <tr key={s.id}>
                      <td>
                        <code style={{fontFamily: 'var(--font-mono)', fontSize: 12, background: 'var(--bg-input)', padding: '3px 8px', borderRadius: 4}}>
                          {s.name}
                        </code>
                      </td>
                      <td style={{color: 'var(--text-secondary)'}}>{s.description}</td>
                      <td>
                        <button className="btn btn-ghost btn-sm" style={{padding: '4px 8px'}}>
                          <Icons.chevron size={16} />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}

        {activeTab === 'urls' && (
          <div className="animate-in delay-1">
            <div className="flex-between mb-24">
              <span className="page-indicator" style={{marginLeft: 0}}>GET /v1/urls · shortly</span>
              <button className="btn btn-primary btn-sm"><Icons.plus size={14} /> Shorten URL</button>
            </div>
            <div className="card">
              {MOCK_URLS.map((u, i) => (
                <div key={i} className="url-item">
                  <div>
                    <div className="url-alias">/{u.alias}</div>
                    <div className="url-target">{u.url}</div>
                  </div>
                  <div className="flex gap-8">
                    <button className="btn btn-ghost btn-sm btn-icon"><Icons.copy size={15} /></button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

// ─── Main App ────────────────────────────────────────────────
export default function App() {
  const [loggedIn, setLoggedIn] = useState(false);
  const [page, setPage] = useState("feed");
  const [profileUser, setProfileUser] = useState(null);

  const navItems = [
    { id: "feed", label: "Feed", icon: Icons.home },
    { id: "explore", label: "Explore", icon: Icons.search },
    { id: "profile", label: "Profile", icon: Icons.user },
    { id: "settings", label: "Settings", icon: Icons.settings },
    { id: "admin", label: "Admin", icon: Icons.shield },
  ];

  const handleUserClick = (user) => {
    setProfileUser(user);
    setPage("viewProfile");
  };

  if (!loggedIn) {
    return (
      <>
        <style>{globalStyles}</style>
        <LoginPage onLogin={() => setLoggedIn(true)} />
      </>
    );
  }

  const renderPage = () => {
    switch (page) {
      case "feed": return <FeedPage onUserClick={handleUserClick} />;
      case "explore": return <ExplorePage onUserClick={handleUserClick} />;
      case "profile": return <ProfilePage user={CURRENT_USER} isSelf={true} />;
      case "viewProfile": return <ProfilePage user={profileUser} isSelf={false} />;
      case "settings": return <SettingsPage />;
      case "admin": return <AdminPage />;
      default: return <FeedPage onUserClick={handleUserClick} />;
    }
  };

  return (
    <>
      <style>{globalStyles}</style>
      <div className="app-root">
        <aside className="sidebar">
          <div className="sidebar-brand">
            <h1>Socialapp</h1>
            <span>v1.0.0</span>
          </div>
          <nav className="sidebar-nav">
            <div className="sidebar-section-label">Navigation</div>
            {navItems.slice(0, 3).map(item => (
              <button
                key={item.id}
                className={`nav-item ${page === item.id || (item.id === 'profile' && page === 'viewProfile') ? 'active' : ''}`}
                onClick={() => {
                  if (item.id === 'profile') setProfileUser(null);
                  setPage(item.id);
                }}
              >
                <span className="nav-icon"><item.icon size={18} /></span>
                {item.label}
              </button>
            ))}
            <div className="sidebar-section-label">Account</div>
            {navItems.slice(3).map(item => (
              <button
                key={item.id}
                className={`nav-item ${page === item.id ? 'active' : ''}`}
                onClick={() => setPage(item.id)}
              >
                <span className="nav-icon"><item.icon size={18} /></span>
                {item.label}
              </button>
            ))}
          </nav>
          <div className="sidebar-user">
            <Avatar user={CURRENT_USER} size="sm" />
            <div className="sidebar-user-info">
              <div className="sidebar-user-name">{CURRENT_USER.first_name} {CURRENT_USER.last_name}</div>
              <div className="sidebar-user-handle">@{CURRENT_USER.username}</div>
            </div>
            <button className="btn btn-ghost btn-icon" style={{padding: 4}} onClick={() => setLoggedIn(false)}>
              <Icons.logout size={16} />
            </button>
          </div>
        </aside>

        <main className="main-content">
          {renderPage()}
        </main>
      </div>
    </>
  );
}
