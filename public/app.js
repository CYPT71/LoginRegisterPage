const API_BASE = 'http://localhost:3000';

const registerBtn = document.getElementById('register-btn');
const loginBtn = document.getElementById('login-btn');
const logoutBtn = document.getElementById('logout-btn');

registerBtn.addEventListener('click', async () => {
    const username = document.getElementById('reg-username').value.trim();
    const password = document.getElementById('reg-password').value;
    if (!username || !password) return;

    const res = await fetch(`${API_BASE}/register/password/${encodeURIComponent(username)}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password })
    });
    const data = await res.json();
    if (res.ok) {
        localStorage.setItem('token', data.token);
        document.getElementById('register-msg').textContent = 'Registered!';
        loadUser();
    } else {
        document.getElementById('register-msg').textContent = data.err || 'Error';
    }
});

loginBtn.addEventListener('click', async () => {
    const username = document.getElementById('login-username').value.trim();
    const password = document.getElementById('login-password').value;
    if (!username || !password) return;

    const res = await fetch(`${API_BASE}/login/password/${encodeURIComponent(username)}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password })
    });
    const data = await res.json();
    if (res.ok) {
        localStorage.setItem('token', data.token);
        document.getElementById('login-msg').textContent = 'Logged in!';
        loadUser();
    } else {
        document.getElementById('login-msg').textContent = data.err || 'Error';
    }
});

logoutBtn.addEventListener('click', () => {
    localStorage.removeItem('token');
    document.getElementById('user-section').style.display = 'none';
});

async function loadUser() {
    const token = localStorage.getItem('token');
    if (!token) return;
    const res = await fetch(`${API_BASE}/user`, {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    if (res.ok) {
        const user = await res.json();
        document.getElementById('user-info').textContent = JSON.stringify(user, null, 2);
        document.getElementById('user-section').style.display = 'block';
    }
}

loadUser();
