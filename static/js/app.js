let currentUser = null;
let token = localStorage.getItem('jwt');

const API_URL = 'http://localhost:8080/api';

async function handleResponse(response) {
    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Request failed');
    }
    return response.json();
}

async function login() {
    const email = document.getElementById('loginEmail').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        const data = await handleResponse(response);
        token = data.token;
        localStorage.setItem('jwt', token);
        currentUser = data.user;
        updateUI();
    } catch (error) {
        alert(error.message);
    }
}

async function register() {
    const name = document.getElementById('regName').value;
    const email = document.getElementById('regEmail').value;
    const password = document.getElementById('regPassword').value;

    try {
        const response = await fetch(`${API_URL}/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, password })
        });

        const data = await handleResponse(response);
        alert('Registration successful! Please login.');
    } catch (error) {
        alert(error.message);
    }
}

async function loadUsers() {
    try {
        const response = await fetch(`${API_URL}/users`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        const users = await handleResponse(response);
        renderUsers(users);
    } catch (error) {
        alert(error.message);
    }
}

function renderUsers(users) {
    const container = document.getElementById('usersList');
    container.innerHTML = users.map(user => `
        <div class="user-card">
            <h3>${user.name} (${user.role})</h3>
            <p>Email: ${user.email}</p>
            ${currentUser.role === 'admin' ? `
                <button onclick="deleteUser(${user.id})">Delete</button>
            ` : ''}
        </div>
    `).join('');
}

async function deleteUser(userId) {
    if (!confirm('Are you sure?')) return;

    try {
        await fetch(`${API_URL}/users/${userId}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        loadUsers();
    } catch (error) {
        alert(error.message);
    }
}

function logout() {
    localStorage.removeItem('jwt');
    token = null;
    currentUser = null;
    updateUI();
}

function updateUI() {
    document.querySelector('.auth-section').style.display = token ? 'none' : 'block';
    document.querySelector('.user-info').style.display = token ? 'block' : 'none';
    document.querySelector('.admin-section').style.display =
        (currentUser && currentUser.role === 'admin') ? 'block' : 'none';

    if (currentUser) {
        document.getElementById('userName').textContent = currentUser.name;
        document.getElementById('userEmail').textContent = currentUser.email;
        document.getElementById('userRole').textContent = currentUser.role;
    }

    if (currentUser?.role === 'admin') {
        loadUsers();
    }
}
function renderUsers(users) {
    const container = document.getElementById('usersList');
    container.innerHTML = users.map(user => `
        <div class="user-card">
            <h3>${user.name} (${user.role})</h3>
            <p>Email: ${user.email}</p>
            ${currentUser.role === 'admin' && user.id !== currentUser.id ? `
                <button onclick="deleteUser(${user.id})">Delete</button>
            ` : ''}
        </div>
    `).join('');
}
// Проверка токена при загрузке
if (token) {
    fetch(`${API_URL}/users/me`, {
        headers: { 'Authorization': `Bearer ${token}` }
    })
        .then(handleResponse)
        .then(user => {
            currentUser = user;
            updateUI();
        })
        .catch(() => {
            localStorage.removeItem('jwt');
            token = null;
        });
}