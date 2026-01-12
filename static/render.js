const API_BASE = 'http://localhost:8080';
const WS_URL = `ws://localhost:8080/ws`;

let ws = null;
let todos = [];
let editingTodoId = null;
let wsReconnectAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 5;

// Initialize WebSocket connection
function connectWebSocket() {
    try {
        ws = new WebSocket(WS_URL);
        
        ws.onopen = () => {
            console.log('WebSocket connected');
            wsReconnectAttempts = 0;
            updateWSStatus(true);
        };
        
        ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            handleWebSocketMessage(message);
        };
        
        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            updateWSStatus(false);
        };
        
        ws.onclose = () => {
            console.log('WebSocket disconnected');
            updateWSStatus(false);
            
            // Attempt to reconnect
            if (wsReconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
                wsReconnectAttempts++;
                setTimeout(() => {
                    console.log(`Reconnecting... (attempt ${wsReconnectAttempts})`);
                    connectWebSocket();
                }, 2000 * wsReconnectAttempts);
            }
        };
    } catch (error) {
        console.error('Failed to connect WebSocket:', error);
        updateWSStatus(false);
    }
}

function updateWSStatus(connected) {
    const indicator = document.getElementById('wsStatus');
    const text = document.getElementById('wsStatusText');
    
    if (connected) {
        indicator.classList.remove('disconnected');
        text.textContent = 'Connected (Real-time updates enabled)';
    } else {
        indicator.classList.add('disconnected');
        text.textContent = 'Disconnected (Polling mode)';
    }
}

function handleWebSocketMessage(message) {
    console.log('WebSocket message received:', message);
    
    switch (message.type) {
        case 'create':
            todos.push(message.payload);
            showRealtimeIndicator('New todo created: ' + message.payload.name);
            renderTodos();
            break;
        case 'update':
            const updateIndex = todos.findIndex(t => t.id === message.payload.id);
            if (updateIndex !== -1) {
                todos[updateIndex] = message.payload;
                showRealtimeIndicator('Todo updated: ' + message.payload.name);
                renderTodos();
            }
            break;
        case 'delete':
            todos = todos.filter(t => t.id !== message.payload.id);
            showRealtimeIndicator('Todo deleted');
            renderTodos();
            break;
    }
}

function showRealtimeIndicator(text) {
    const indicator = document.getElementById('realtimeIndicator');
    indicator.textContent = text;
    indicator.classList.add('show');
    setTimeout(() => {
        indicator.classList.remove('show');
    }, 3000);
}

// Load todos from API
async function loadTodos() {
    try {
        const status = document.getElementById('filterStatus').value;
        const sortBy = document.getElementById('sortBy').value;
        const sortOrder = document.getElementById('sortOrder').value;
        
        let url = `${API_BASE}/todos?`;
        const params = new URLSearchParams();
        if (status) params.append('status', status);
        if (sortBy) params.append('sort_by', sortBy);
        if (sortOrder) params.append('order', sortOrder);
        
        url += params.toString();
        
        const response = await fetch(url);
        if (!response.ok) throw new Error('Failed to load todos');
        
        todos = await response.json();
        renderTodos();
    } catch (error) {
        showError('Failed to load todos: ' + error.message);
        document.getElementById('todoList').innerHTML = '<div class="error">Failed to load todos. Please refresh the page.</div>';
    }
}

// Render todos to the DOM
function renderTodos() {
    const container = document.getElementById('todoList');
    const searchTerm = document.getElementById('searchInput').value.toLowerCase();
    
    let filteredTodos = todos;
    
    // Apply search filter
    if (searchTerm) {
        filteredTodos = todos.filter(todo => 
            todo.name.toLowerCase().includes(searchTerm) ||
            (todo.description && todo.description.toLowerCase().includes(searchTerm))
        );
    }
    
    if (filteredTodos.length === 0) {
        container.innerHTML = '<div class="empty-state"><h2>No todos found</h2><p>Create your first todo to get started!</p></div>';
        return;
    }
    
    container.innerHTML = filteredTodos.map(todo => `
        <div class="todo-card ${todo.status === 'completed' ? 'completed' : ''}" data-id="${todo.id}">
            <div class="todo-header">
                <div>
                    <div class="todo-title">${escapeHtml(todo.name)}</div>
                    <span class="status-badge status-${todo.status.replace('_', '-')}">${todo.status.replace('_', ' ')}</span>
                </div>
                <div class="todo-actions">
                    <button class="btn btn-secondary btn-small" onclick="editTodo(${todo.id})">Edit</button>
                    <button class="btn btn-danger btn-small" onclick="deleteTodo(${todo.id})">Delete</button>
                </div>
            </div>
            ${todo.description ? `<p style="margin: 10px 0; color: #4b5563;">${escapeHtml(todo.description)}</p>` : ''}
            <div class="todo-meta">
                ${todo.due_date ? `<span>📅 Due: ${formatDate(todo.due_date)}</span>` : ''}
                ${todo.priority ? `<span>⭐ Priority: ${todo.priority}</span>` : ''}
            </div>
            ${todo.tags && todo.tags.length > 0 ? `
                <div class="todo-tags">
                    ${todo.tags.map(tag => `<span class="tag">${escapeHtml(tag)}</span>`).join('')}
                </div>
            ` : ''}
        </div>
    `).join('');
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
}

// Create todo
async function createTodo(todoData) {
    try {
        const response = await fetch(`${API_BASE}/todos`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(todoData)
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to create todo');
        }
        
        const created = await response.json();
        
        // WebSocket will handle the broadcast, but we can add locally for immediate feedback
        if (!todos.find(t => t.id === created.id)) {
            todos.push(created);
        }
        
        document.getElementById('todoForm').reset();
        showSuccess('Todo created successfully!');
        
        // Only re-render if WebSocket didn't already update
        setTimeout(() => renderTodos(), 100);
    } catch (error) {
        showError('Failed to create todo: ' + error.message);
    }
}

// Update todo
async function updateTodo(id, todoData) {
    try {
        const response = await fetch(`${API_BASE}/todos/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(todoData)
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to update todo');
        }
        
        const updated = await response.json();
        const index = todos.findIndex(t => t.id === id);
        if (index !== -1) {
            todos[index] = updated;
        }
        
        editingTodoId = null;
        document.getElementById('todoForm').reset();
        document.getElementById('cancelEdit').style.display = 'none';
        document.querySelector('button[type="submit"]').textContent = 'Create Todo';
        showSuccess('Todo updated successfully!');
        
        setTimeout(() => renderTodos(), 100);
    } catch (error) {
        showError('Failed to update todo: ' + error.message);
    }
}

// Delete todo
async function deleteTodo(id) {
    if (!confirm('Are you sure you want to delete this todo?')) {
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/todos/${id}`, {
            method: 'DELETE'
        });
        
        if (!response.ok) {
            throw new Error('Failed to delete todo');
        }
        
        todos = todos.filter(t => t.id !== id);
        showSuccess('Todo deleted successfully!');
        renderTodos();
    } catch (error) {
        showError('Failed to delete todo: ' + error.message);
    }
}

// Edit todo
function editTodo(id) {
    const todo = todos.find(t => t.id === id);
    if (!todo) return;
    
    editingTodoId = id;
    document.getElementById('todoName').value = todo.name;
    document.getElementById('todoDescription').value = todo.description || '';
    document.getElementById('todoStatus').value = todo.status;
    document.getElementById('todoPriority').value = todo.priority || 0;
    
    if (todo.due_date) {
        const date = new Date(todo.due_date);
        const localDateTime = new Date(date.getTime() - date.getTimezoneOffset() * 60000).toISOString().slice(0, 16);
        document.getElementById('todoDueDate').value = localDateTime;
    } else {
        document.getElementById('todoDueDate').value = '';
    }
    
    document.getElementById('cancelEdit').style.display = 'inline-block';
    document.querySelector('button[type="submit"]').textContent = 'Update Todo';
    document.getElementById('todoForm').scrollIntoView({ behavior: 'smooth' });
}

function cancelEdit() {
    editingTodoId = null;
    document.getElementById('todoForm').reset();
    document.getElementById('cancelEdit').style.display = 'none';
    document.querySelector('button[type="submit"]').textContent = 'Create Todo';
}

function showError(message) {
    const errorDiv = document.getElementById('errorMessage');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
    setTimeout(() => {
        errorDiv.style.display = 'none';
    }, 5000);
}

function showSuccess(message) {
    const successDiv = document.getElementById('successMessage');
    successDiv.textContent = message;
    successDiv.style.display = 'block';
    setTimeout(() => {
        successDiv.style.display = 'none';
    }, 3000);
}

// Event listeners
document.getElementById('todoForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const formData = {
        name: document.getElementById('todoName').value.trim(),
        description: document.getElementById('todoDescription').value.trim(),
        status: document.getElementById('todoStatus').value,
        priority: parseInt(document.getElementById('todoPriority').value) || 0,
    };
    
    const dueDate = document.getElementById('todoDueDate').value;
    if (dueDate) {
        formData.due_date = new Date(dueDate).toISOString();
    }
    
    if (!formData.name) {
        showError('Todo name is required');
        return;
    }
    
    if (editingTodoId) {
        await updateTodo(editingTodoId, formData);
    } else {
        await createTodo(formData);
    }
});

document.getElementById('filterStatus').addEventListener('change', loadTodos);
document.getElementById('sortBy').addEventListener('change', loadTodos);
document.getElementById('sortOrder').addEventListener('change', loadTodos);
document.getElementById('searchInput').addEventListener('input', renderTodos);

// Initialize
connectWebSocket();
loadTodos();

// Polling fallback if WebSocket is not available
setInterval(() => {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        loadTodos();
    }
}, 5000); // Poll every 5 seconds if WebSocket is down