// Authentication and team management JavaScript

const API_BASE = 'http://localhost:8080';
let authToken = localStorage.getItem('authToken');
let currentUser = null;
let currentTeamId = null;

// Authentication functions
async function login(email, password) {
    try {
        const response = await fetch(`${API_BASE}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Login failed');
        }

        const data = await response.json();
        setAuthToken(data.token);
        currentUser = data.user;
        
        // Load user's teams
        await loadTeams();
        
        return data;
    } catch (error) {
        throw error;
    }
}

async function register(username, email, password) {
    try {
        const response = await fetch(`${API_BASE}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, email, password })
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Registration failed');
        }

        const user = await response.json();
        return user;
    } catch (error) {
        throw error;
    }
}

async function logout() {
    authToken = null;
    currentUser = null;
    currentTeamId = null;
    localStorage.removeItem('authToken');
    localStorage.removeItem('currentTeamId');
    showAuthView();
    if (window.ws) {
        window.ws.close();
    }
}

function setAuthToken(token) {
    authToken = token;
    localStorage.setItem('authToken', token);
}

function getAuthToken() {
    return authToken || localStorage.getItem('authToken');
}

async function getCurrentUser() {
    if (currentUser) return currentUser;
    
    try {
        const token = getAuthToken();
        if (!token) return null;

        const response = await fetch(`${API_BASE}/auth/me`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            return null;
        }

        currentUser = await response.json();
        return currentUser;
    } catch (error) {
        return null;
    }
}

function isAuthenticated() {
    return !!getAuthToken();
}

// Team functions
async function loadTeams() {
    try {
        const token = getAuthToken();
        if (!token) return [];

        const response = await fetch(`${API_BASE}/teams`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            return [];
        }

        const teams = await response.json();
        
        // Update team selector
        const teamSelector = document.getElementById('teamSelector');
        if (teamSelector) {
            teamSelector.innerHTML = '<option value="">Select a team</option>' +
                teams.map(team => `<option value="${team.id}">${team.name}</option>`).join('');
            
            // Restore previous selection
            const savedTeamId = localStorage.getItem('currentTeamId');
            if (savedTeamId && teams.find(t => t.id == savedTeamId)) {
                teamSelector.value = savedTeamId;
                currentTeamId = parseInt(savedTeamId);
            }
        }
        
        return teams;
    } catch (error) {
        console.error('Failed to load teams:', error);
        return [];
    }
}

function selectTeam(teamId) {
    currentTeamId = teamId ? parseInt(teamId) : null;
    if (teamId) {
        localStorage.setItem('currentTeamId', teamId);
    } else {
        localStorage.removeItem('currentTeamId');
    }
    
    // Reload todos for the selected team
    if (window.loadTodos) {
        window.loadTodos();
    }
}

function getCurrentTeamId() {
    return currentTeamId || parseInt(localStorage.getItem('currentTeamId')) || null;
}

// UI functions
function showAuthView() {
    const authContainer = document.getElementById('authContainer');
    const todoContainer = document.getElementById('todoContainer');
    
    if (authContainer) authContainer.style.display = 'block';
    if (todoContainer) todoContainer.style.display = 'none';
}

function showTodoView() {
    const authContainer = document.getElementById('authContainer');
    const todoContainer = document.getElementById('todoContainer');
    
    if (authContainer) authContainer.style.display = 'none';
    if (todoContainer) todoContainer.style.display = 'block';
}

function updateUserMenu() {
    const userMenu = document.getElementById('userMenu');
    if (userMenu && currentUser) {
        userMenu.innerHTML = `
            <span>${escapeHtml(currentUser.username)}</span>
            <button onclick="logout()" class="btn btn-secondary btn-small">Logout</button>
        `;
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', async () => {
    const token = getAuthToken();
    if (token) {
        const user = await getCurrentUser();
        if (user) {
            currentUser = user;
            await loadTeams();
            showTodoView();
            updateUserMenu();
            
            // Initialize WebSocket connection
            if (window.connectWebSocket) {
                window.connectWebSocket();
            }
        } else {
            logout();
        }
    } else {
        showAuthView();
    }
});
