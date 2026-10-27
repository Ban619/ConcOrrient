const state = { actor: localStorage.getItem("concorrient-actor") || "mina", members: [], ideas: [], notices: [], filter: "all", expanded: "", snapshot: "", toastTimer: null };
let deferredPrompt = null;
const $ = (selector, root = document) => root.querySelector(selector);
const $$ = (selector, root = document) => [...root.querySelectorAll(selector)];
const colors = { lilac: "lilac", coral: "coral", mint: "mint", blue: "blue" };

if ("serviceWorker" in navigator) {
  window.addEventListener("load", () => {
    navigator.serviceWorker.register("/sw.js").catch(() => undefined);
  });
}

window.addEventListener("beforeinstallprompt", (event) => {
  event.preventDefault();
  deferredPrompt = event;
  const installButton = document.getElementById("install-button");
  if (installButton) installButton.classList.remove("hidden");
});

window.addEventListener("appinstalled", () => {
  const installButton = document.getElementById("install-button");
  if (installButton) installButton.classList.add("hidden");
  showToast("ConcOrrient is installed.");
});

const installButton = document.getElementById("install-button");
if (installButton) {
  installButton.addEventListener("click", async () => {
    if (!deferredPrompt) {
      showToast("This browser is already using the app or doesn’t support install prompts.");
      return;
    }
    deferredPrompt.prompt();
    const outcome = await deferredPrompt.userChoice;
    if (outcome.outcome === "accepted") {
      showToast("Installing ConcOrrient…");
    }
    deferredPrompt = null;
    installButton.classList.add("hidden");
  });
}

function esc(value = "") { return String(value).replace(/[&<>"']/g, char => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" })[char]); }
function member(id) { return state.members.find(item => item.id === id) || { id: "?", name: "Teammate", initials: "?", color: "mint" }; }
function avatar(id, extra = "") { const person = member(id); return `<span class="avatar avatar-${colors[person.color] || "mint"} ${extra}" title="${esc(person.name)}">${esc(person.initials)}</span>`; }
function timeAgo(date) {
  const seconds = Math.max(0, Math.floor((Date.now() - new Date(date).getTime()) / 1000));
  if (seconds < 60) return "just now";
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`;
  if (seconds < 604800) return `${Math.floor(seconds / 86400)}d ago`;
  return new Date(date).toLocaleDateString(undefined, { month: "short", day: "numeric" });
}
function sameDay(left, right) { const a = new Date(left), b = new Date(right); return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate(); }
function activeParticipants(idea) { return [...new Set([idea.authorId, ...(idea.collaborator ? [idea.collaborator] : []), ...idea.updates.map(update => update.authorId)])]; }

async function refresh(force = false) {
  try {
    const response = await fetch(`/api/state?actor=${encodeURIComponent(state.actor)}`, { cache: "no-store" });
    if (!response.ok) throw new Error("Couldn’t load this workspace.");
    const next = await response.json();
    state.members = next.members || [];
    state.ideas = next.ideas || [];
    state.notices = next.notices || [];
    const signature = JSON.stringify([state.members, state.ideas, state.notices, state.actor, state.filter, state.expanded]);
    if (!force && signature === state.snapshot) return;
    if (!force && ["INPUT", "TEXTAREA", "SELECT"].includes(document.activeElement?.tagName) && document.activeElement.closest(".thread-detail, .idea-dialog")) return;
    state.snapshot = signature;
    render();
  } catch (error) { showToast(error.message || "Couldn’t reach the workspace."); }
}

function render() {
  renderMembers();
  renderActor();
  renderIdeas();
  renderNotifications();
}

function renderActor() {
  const select = $("#actor-select");
  select.innerHTML = state.members.map(person => `<option value="${esc(person.id)}">Viewing as ${esc(person.name)}</option>`).join("");
  if (!state.members.some(person => person.id === state.actor)) state.actor = state.members[0]?.id || "mina";
  select.value = state.actor;
  $("#actor-avatar").className = `avatar avatar-current avatar-${colors[member(state.actor).color] || "mint"}`;
  $("#actor-avatar").textContent = member(state.actor).initials;
  const collaboratorSelect = $("#collaborator-select");
  collaboratorSelect.innerHTML = state.members.filter(person => person.id !== state.actor).map(person => `<option value="${esc(person.id)}">${esc(person.name)}</option>`).join("");
}

function renderMembers() {
  $("#member-total").textContent = state.members.length;
  $("#member-list").innerHTML = state.members.map(person => `<div class="member-row">${avatar(person.id)}<span>${esc(person.name.split(" ")[0])}${person.id === state.actor ? " · you" : ""}</span><i class="member-presence" aria-label="In workspace"></i></div>`).join("");
  $("#idea-count").textContent = state.ideas.length;
}

function renderIdeas() {
  const openCount = state.ideas.filter(idea => idea.status !== "resolved").length;
  const todayCount = state.ideas.filter(idea => sameDay(idea.updatedAt, Date.now())).length;
  $("#open-count").textContent = openCount;
  $("#today-count").textContent = todayCount;
  const ideas = state.ideas.filter(idea => state.filter === "all" || idea.status === state.filter);
  $("#visible-count").textContent = ideas.length;
  $("#empty-state").classList.toggle("hidden", ideas.length > 0);
  $("#thread-list").classList.toggle("hidden", ideas.length === 0);
  $("#thread-list").innerHTML = ideas.map(renderIdea).join("");
}

function renderIdea(idea) {
  const people = activeParticipants(idea);
  const names = people.map(id => member(id).name.split(" ")[0]);
  const updates = [...idea.updates].sort((a, b) => new Date(a.createdAt) - new Date(b.createdAt));
  const statusLabel = { open: "Open", exploring: "Exploring", resolved: "Settled" }[idea.status] || "Open";
  const modeLabel = idea.mode === "direct" ? "Direct" : "Open space";
  const tags = idea.tags || [];
  return `<article class="thread-card" id="card-${esc(idea.id)}">
    <div class="thread-main">
      <div class="thread-topline"><div class="thread-labels"><span class="status-pill status-${esc(idea.status)}">${statusLabel}</span><span class="mode-pill ${idea.mode === "open" ? "open" : ""}">${modeLabel}</span></div><time class="thread-time" datetime="${esc(idea.updatedAt)}">Updated ${timeAgo(idea.updatedAt)}</time></div>
      <h3 class="thread-title">${esc(idea.title)}</h3><p class="thread-summary">${esc(idea.summary)}</p>
      ${tags.length ? `<div class="tag-list">${tags.map(tag => `<span class="tag">#${esc(tag)}</span>`).join("")}</div>` : ""}
    </div>
    <div class="thread-footer"><div class="avatar-stack">${people.slice(0, 4).map(id => avatar(id)).join("")}</div><span class="participant-label">${esc(names.join(", "))}</span><span class="thread-footer-spacer"></span><span class="update-count"><svg viewBox="0 0 16 16" aria-hidden="true"><path d="M3 3.5h10v7H8l-3 2v-2H3z" /></svg>${updates.length} ${updates.length === 1 ? "thought" : "thoughts"}</span><button class="thread-action" data-expand="${esc(idea.id)}">${state.expanded === idea.id ? "Close thread −" : "Join thread +"}</button>
      <select class="status-control" data-status="${esc(idea.id)}" aria-label="Change thread status"><option value="open" ${idea.status === "open" ? "selected" : ""}>Open</option><option value="exploring" ${idea.status === "exploring" ? "selected" : ""}>Exploring</option><option value="resolved" ${idea.status === "resolved" ? "selected" : ""}>Settled</option></select>
    </div>
    ${state.expanded === idea.id ? `<div class="thread-detail"><div class="activity-feed">${updates.map(update => `<div class="activity-item">${avatar(update.authorId)}<div class="activity-copy"><div class="activity-byline"><strong>${esc(member(update.authorId).name)}</strong><time datetime="${esc(update.createdAt)}">${timeAgo(update.createdAt)}</time></div><p class="activity-message">${esc(update.message)}</p></div></div>`).join("")}</div><form class="reply-form" data-reply="${esc(idea.id)}"><textarea name="message" maxlength="1000" placeholder="Add a thought when you’re ready…" required></textarea><button type="submit">Add thought ↗</button></form></div>` : ""}
  </article>`;
}

function renderNotifications() {
  const unread = state.notices.filter(notice => !notice.read).length;
  $("#side-unread").classList.toggle("hidden", unread === 0);
  const list = $("#notification-list");
  if (!state.notices.length) {
    list.innerHTML = `<div class="no-notifications">No new threads have found you yet. Open a conversation or invite a teammate to get started.</div>`;
    return;
  }
  list.innerHTML = [...state.notices].sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt)).slice(0, 6).map(notice => `<div class="notification ${notice.read ? "" : "unread"}" tabindex="0" role="button" data-notice-idea="${esc(notice.ideaId)}"><p>${esc(notice.text)} <span class="notice-kind ${notice.kind === "open" ? "open" : ""}">${notice.kind === "direct" ? "· direct" : "· open"}</span></p><time datetime="${esc(notice.createdAt)}">${timeAgo(notice.createdAt)}</time></div>`).join("");
}

async function api(url, options = {}) {
  const response = await fetch(url, { ...options, headers: { "Content-Type": "application/json", ...(options.headers || {}) } });
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message.trim() || "That didn’t save. Please try again.");
  }
  if (response.status === 204) return null;
  return response.json();
}

function showToast(message) {
  const toast = $("#toast");
  toast.textContent = message;
  toast.classList.add("show");
  clearTimeout(state.toastTimer);
  state.toastTimer = setTimeout(() => toast.classList.remove("show"), 2600);
}

function openDialog() {
  $("#form-error").textContent = "";
  $("#idea-form").reset();
  $$(".mode-option").forEach(option => option.classList.toggle("selected", option.dataset.modeOption === "open"));
  $("#collaborator-field").classList.add("hidden");
  renderActor();
  $("#idea-dialog").showModal();
  setTimeout(() => $("#idea-title").focus(), 40);
}

$("#new-idea-button").addEventListener("click", openDialog);
$("#empty-create").addEventListener("click", openDialog);
$("#close-dialog").addEventListener("click", () => $("#idea-dialog").close());
$("#cancel-dialog").addEventListener("click", () => $("#idea-dialog").close());
$("#idea-dialog").addEventListener("click", event => { if (event.target === $("#idea-dialog")) $("#idea-dialog").close(); });
$("#idea-dialog").addEventListener("cancel", () => { $("#form-error").textContent = ""; });

$$('[name="mode"]').forEach(radio => radio.addEventListener("change", () => {
  $$(".mode-option").forEach(option => option.classList.toggle("selected", option.dataset.modeOption === radio.value));
  $("#collaborator-field").classList.toggle("hidden", radio.value !== "direct");
}));

$("#idea-form").addEventListener("submit", async event => {
  event.preventDefault();
  const form = new FormData(event.currentTarget);
  const payload = {
    title: form.get("title"), summary: form.get("summary"), mode: form.get("mode"),
    collaborator: form.get("collaborator") || "", authorId: state.actor,
    tags: String(form.get("tags") || "").split(",").map(tag => tag.trim()).filter(Boolean),
  };
  const submit = $("#idea-form button[type=submit]");
  submit.disabled = true;
  $("#form-error").textContent = "";
  try {
    const idea = await api("/api/ideas", { method: "POST", body: JSON.stringify(payload) });
    $("#idea-dialog").close();
    state.expanded = idea.id;
    await refresh(true);
    showToast(payload.mode === "direct" ? `Thread started with ${member(payload.collaborator).name.split(" ")[0]}.` : "Your open thread is ready.");
  } catch (error) { $("#form-error").textContent = error.message; }
  finally { submit.disabled = false; }
});

$("#actor-select").addEventListener("change", async event => {
  state.actor = event.target.value;
  localStorage.setItem("concorrient-actor", state.actor);
  await refresh(true);
});

$$('.filter-button').forEach(button => button.addEventListener("click", () => {
  state.filter = button.dataset.filter;
  $$(".filter-button").forEach(item => { const selected = item === button; item.classList.toggle("selected", selected); item.setAttribute("aria-selected", selected); });
  renderIdeas();
  state.snapshot = "";
}));

$("#thread-list").addEventListener("click", async event => {
  const expand = event.target.closest("[data-expand]");
  if (expand) { state.expanded = state.expanded === expand.dataset.expand ? "" : expand.dataset.expand; state.snapshot = ""; renderIdeas(); return; }
});

$("#thread-list").addEventListener("submit", async event => {
  const form = event.target.closest("[data-reply]");
  if (!form) return;
  event.preventDefault();
  const button = $("button", form);
  button.disabled = true;
  try {
    await api(`/api/ideas/${encodeURIComponent(form.dataset.reply)}/updates`, { method: "POST", body: JSON.stringify({ authorId: state.actor, message: new FormData(form).get("message") }) });
    await refresh(true);
    showToast("Thought added to the thread.");
  } catch (error) { showToast(error.message); }
  finally { button.disabled = false; }
});

$("#thread-list").addEventListener("change", async event => {
  const select = event.target.closest("[data-status]");
  if (!select) return;
  select.disabled = true;
  try {
    await api(`/api/ideas/${encodeURIComponent(select.dataset.status)}`, { method: "PATCH", body: JSON.stringify({ actor: state.actor, status: select.value }) });
    await refresh(true);
    showToast("Thread status updated.");
  } catch (error) { showToast(error.message); await refresh(true); }
  finally { select.disabled = false; }
});

$("#mark-read").addEventListener("click", async () => {
  try { await api("/api/notices/read", { method: "POST", body: JSON.stringify({ actor: state.actor }) }); await refresh(true); }
  catch (error) { showToast(error.message); }
});

$("#notification-list").addEventListener("click", async event => {
  const item = event.target.closest("[data-notice-idea]");
  if (!item) return;
  state.filter = "all"; state.expanded = item.dataset.noticeIdea;
  $$(".filter-button").forEach(button => { const selected = button.dataset.filter === "all"; button.classList.toggle("selected", selected); button.setAttribute("aria-selected", selected); });
  await refresh(true);
  document.querySelector(`#card-${CSS.escape(state.expanded)}`)?.scrollIntoView({ behavior: "smooth", block: "center" });
});
$("#notification-list").addEventListener("keydown", event => { if (event.key === "Enter" || event.key === " ") event.target.click(); });
$("#activity-nav").addEventListener("click", () => $("#notifications-section").scrollIntoView({ behavior: "smooth", block: "center" }));

refresh(true);
setInterval(() => refresh(false), 7000);
