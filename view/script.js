// Konfigurasi axios untuk CORS
axios.defaults.withCredentials = true;
axios.interceptors.request.use(function (config) {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, function (error) {
    return Promise.reject(error);
});

// Dark mode toggle
const darkModeSwitch = document.getElementById('darkModeSwitch');
darkModeSwitch.addEventListener('change', () => {
    document.documentElement.setAttribute('data-bs-theme', darkModeSwitch.checked ? 'dark' : 'light');
});

// Tooltip initialization
const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
const tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
    return new bootstrap.Tooltip(tooltipTriggerEl);
});

// close and open sidebar
const sidebar = document.getElementById('sidebar');
const toggleSidebarBtn = document.getElementById('toggleSidebar');
const openSidebarBtn = document.getElementById('openSidebarBtn');
const content = document.querySelector('main');

function toggleSidebar() {
    sidebar.classList.toggle('active');
    content.classList.toggle('active');
    openSidebarBtn.classList.toggle('d-none');
    
    setTimeout(() => {
        renderBoard();
        updateSidebarHeight();
    }, 300);
}

toggleSidebarBtn.addEventListener('click', toggleSidebar);
openSidebarBtn.addEventListener('click', toggleSidebar);

// Global variables
let boards = [];
let currentBoard = null;
let customRoleNames = {
    Owner: 'Owner',
    Manager: 'Manager',
    Employee: 'Employee'
};

// DOM Elements
const boardsList = document.getElementById('boardsList');
const boardContent = document.getElementById('boardContent');
const addBoardBtn = document.getElementById('addBoardBtn');
const addBoardModal = new bootstrap.Modal(document.getElementById('addBoardModal'));
const saveNewBoardBtn = document.getElementById('saveNewBoard');
const newBoardNameInput = document.getElementById('name_task');
const addPersonModal = new bootstrap.Modal(document.getElementById('addPersonModal'));
const saveNewPersonBtn = document.getElementById('saveNewPerson');
const searchBoardInput = document.getElementById('searchBoard');


// Event Listeners
addBoardBtn.addEventListener('click', () => addBoardModal.show());
saveNewBoardBtn.addEventListener('click', createNewBoard);
saveNewPersonBtn.addEventListener('click', addNewPerson);
searchBoardInput.addEventListener('input', searchBoards);
addBoardBtn.addEventListener('keydown', handleAddBoardKeydown);
newBoardNameInput.addEventListener('keydown', handleNewBoardNameKeydown);
addBoardBtn.addEventListener('click', openAddBoardModal);

// Fungsi untuk memeriksa autentikasi
function checkAuth() {
    const token = localStorage.getItem('token');
    const userEmail = localStorage.getItem('userEmail');
    if (!token || !userEmail) {
        window.location.href = 'login.html';
    }
    return { token, userEmail };
}

// Fungsi untuk logout
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('userEmail');
    window.location.href = 'login.html';
}

// Fungsi untuk mengambil task berdasarkan ID
async function fetchTask(taskId) {
    const cachedTask = getCachedData(`task_${taskId}`);
    if (cachedTask) {
        return cachedTask;
    }

    const { token } = checkAuth();
    try {
        const response = await axios.get(`http://127.0.0.1:8080/task/${taskId}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const taskData = response.data.data;
        cacheData(`task_${taskId}`, taskData);
        return taskData;
    } catch (error) {
        console.error('Error fetching task:', error);
        if (error.response && error.response.status === 401) {
            logout();
        }
        return null;
    }
}

// Fungsi untuk memuat board
async function loadBoards() {
    const { token } = checkAuth();
    try {
        const boardIDs = getCookie("BoardIDs");
        if (boardIDs) {
            const ids = JSON.parse(decodeURIComponent(boardIDs));
            boards = [];
            for (const id of ids) {
                const response = await axios.get(`http://127.0.0.1:8080/board/${id}`, {
                    headers: { 'Authorization': `Bearer ${token}` },
                    withCredentials: true
                });
                if (response.data.code === 200) {
                    boards.push({
                        id: response.data.data.board_id,
                        name: response.data.data.name_board,
                    });
                }
            }
            renderBoardsList();
        }
    } catch (error) {
        console.error('Error loading boards:', error);
    }
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

function openAddBoardModal() {
    addBoardModal.show();
    addBoardModal._element.addEventListener('shown.bs.modal', function onModalShown() {
        newBoardNameInput.focus();
        addBoardModal._element.removeEventListener('shown.bs.modal', onModalShown);
    });
}

function handleAddBoardKeydown(event) {
    if (event.key === 'Enter') {
        addBoardModal.show();
        setTimeout(() => newBoardNameInput.focus(), 300);
    }
}

function handleNewBoardNameKeydown(event) {
    if (event.key === 'Enter') {
        event.preventDefault();
        createNewBoard();
    }
}

// if unauthorized
function showSessionExpiredOverlay() {
    document.getElementById('sessionExpiredOverlay').style.display = 'flex';
}

function redirectToLogin() {
    window.location.href = "login.html";
}

// Fungsi untuk memeriksa response dan menampilkan overlay jika unauthorized
function checkAuthorization(error) {
    if (error.response && (error.response.status === 401 || error.response.status === 403)) {
        showSessionExpiredOverlay();
    }
}

// Menambahkan interceptor untuk semua request Axios
axios.interceptors.response.use(
    response => response,
    error => {
        checkAuthorization(error);
        return Promise.reject(error);
    }
);

// Fungsi untuk membuat board baru
async function createNewBoard() {
    const { token, userEmail } = checkAuth();
    const boardName = newBoardNameInput.value.trim();
    if (boardName) {
        try {
            // Buat objek FormData
            const formData = new FormData();
            formData.append('name_task', boardName);

            const response = await axios.post('http://127.0.0.1:8080/board', 
                formData,
                { 
                    headers: { 
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'multipart/form-data'
                    } 
                }
            );
            if (response.data.code === 200) {
                const newBoard = {
                    id: response.data.data.task_id,
                    name: response.data.data.name_task,
                    ownerEmail: userEmail,
                    todo: [],
                    completed: []
                };
                boards.push(newBoard);
                
                // Simpan ke localStorage
                const storedBoards = JSON.parse(localStorage.getItem('boards')) || [];
                storedBoards.push({ taskId: newBoard.id, ownerEmail: userEmail });
                localStorage.setItem('boards', JSON.stringify(storedBoards));
                
                renderBoardsList();
                addBoardModal.hide();
                newBoardNameInput.value = '';
                loadBoard(newBoard.id);
            }
        } catch (error) {
            console.error('Error creating board:', error);
            if (error.response) {
                console.error('Response data:', error.response.data);
                console.error('Response status:', error.response.status);
                console.error('Response headers:', error.response.headers);
            }
        }
    }
}

window.addEventListener('load', loadBoards);
addBoardBtn.addEventListener('click', () => addBoardModal.show());
saveNewBoardBtn.addEventListener('click', createNewBoard);

function addInitialItem(boardId) {
    const board = boards.find(b => b.id === boardId);
    if (board && board.owner) {
        const newItem = {
            persons: [{
                email: board.owner.email,
                role: 'Owner',
                roleDisplay: 'Owner'
            }],
            planningDescription: '',
            planningFile: null,
            planningDueDate: '',
            planningStatus: 'Not Approved',
            priority: 'Low',
            projectFile: null,
            projectComment: '',
            projectStatus: 'Undone',
            projectDueDate: ''
        };
        board.todo.push(newItem);
        renderBoard();
    }
}

function clearActiveBoardLinks() {
    document.querySelectorAll('#boardsList .nav-link').forEach(link => {
        link.classList.remove('active');
    });
}

// Fungsi untuk merender daftar board
function renderBoardsList() {
    boardsList.innerHTML = '';
    boards.forEach((board, index) => {
        const li = document.createElement('li');
        li.className = 'nav-item d-flex justify-content-between align-items-center';
        li.innerHTML = `
            <a class="nav-link ${currentBoard && currentBoard.id === board.id ? 'active' : ''}" href="#" data-board-id="${board.id}" title="${board.name}">${board.name}</a>
            <div class="board-actions">
                <button class="btn btn-sm btn-primary rename-board-btn" data-board-id="${board.id}">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="btn btn-sm btn-danger delete-board-btn" data-board-id="${board.id}">
                    <i class="fas fa-trash"></i>
                </button>
            </div>
        `;
        li.querySelector('a').addEventListener('click', (e) => {
            e.preventDefault();
            loadBoard(board.id);
        });
        li.querySelector('.rename-board-btn').addEventListener('click', (e) => {
            e.stopPropagation();
            showRenameBoardModal(board.id, board.name);
        });
        li.querySelector('.delete-board-btn').addEventListener('click', (e) => {
            e.stopPropagation();
            showDeleteBoardConfirmModal(board.id, board.name);
        });
        boardsList.appendChild(li);
    });
}

function showRenameBoardModal(boardId, currentName) {
    const renameBoardModal = new bootstrap.Modal(document.getElementById('renameBoardModal'));
    const renameBoardInput = document.getElementById('renameBoardInput');
    const saveRenamedBoardBtn = document.getElementById('saveRenamedBoard');
    
    renameBoardInput.value = currentName;
    
    function renameBoard() {
        const newName = renameBoardInput.value.trim();
        if (newName && newName !== currentName) {
            const board = boards.find(b => b.id === boardId);
            if (board) {
                board.name = newName;
                renderBoardsList();
                if (currentBoard && currentBoard.id === boardId) {
                    currentBoard.name = newName;
                    renderBoard();
                }
            }
        }
        renameBoardModal.hide();
    }
    
    saveRenamedBoardBtn.onclick = renameBoard;
    
    renameBoardModal._element.addEventListener('shown.bs.modal', function () {
        renameBoardInput.focus();
        renameBoardInput.select();
    });
    
    renameBoardModal._element.addEventListener('keydown', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            renameBoard();
        }
    });
    
    renameBoardModal.show();
}

// Variabel untuk pagination
let currentPage = 1;
const itemsPerPage = 10;

// Fungsi untuk memuat board dengan pagination
async function loadBoardsWithPagination() {
    const { userEmail } = checkAuth();
    const storedBoards = JSON.parse(localStorage.getItem('boards')) || [];
    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const paginatedBoards = storedBoards.slice(startIndex, endIndex);

    boards = [];
    for (const storedBoard of paginatedBoards) {
        if (storedBoard.ownerEmail === userEmail) {
            const taskData = await fetchTask(storedBoard.taskId);
            if (taskData) {
                boards.push({
                    id: taskData.id,
                    name: taskData.name_task,
                    ownerEmail: taskData.owner.email,
                    todo: [],
                    completed: []
                });
            }
        }
    }
    renderBoardsList();
    renderPaginationControls(storedBoards.length);
}

// Fungsi untuk merender kontrol pagination
function renderPaginationControls(totalItems) {
    const totalPages = Math.ceil(totalItems / itemsPerPage);
    const paginationContainer = document.getElementById('paginationContainer');
    paginationContainer.innerHTML = '';

    for (let i = 1; i <= totalPages; i++) {
        const button = document.createElement('button');
        button.textContent = i;
        button.classList.add('btn', 'btn-sm', 'btn-outline-primary', 'm-1');
        if (i === currentPage) {
            button.classList.add('active');
        }
        button.addEventListener('click', () => {
            currentPage = i;
            loadBoardsWithPagination();
        });
        paginationContainer.appendChild(button);
    }
}

// Fungsi untuk memuat lebih banyak item
async function loadMoreItems(listType) {
    try {
        const taskData = await fetchTask(currentBoard.id);
        if (taskData && taskData[listType]) {
            const currentItems = currentBoard[listType].length;
            const newItems = taskData[listType].slice(currentItems, currentItems + 10);
            currentBoard[listType] = [...currentBoard[listType], ...newItems];
            renderBoard(taskData);
        } else {
            console.error('No task data or list type not found');
        }
    } catch (error) {
        console.error('Error loading more items:', error);
    }
}

function getCachedBoard(boardId) {
    const cachedBoard = localStorage.getItem(`board_${boardId}`);
    return cachedBoard ? JSON.parse(cachedBoard) : null;
}

function cacheBoard(board) {
    localStorage.setItem(`board_${board.id}`, JSON.stringify(board));
}

async function loadBoard(boardId) {
    let board = getCachedBoard(boardId);
    if (!board) {
        try {
            const taskData = await fetchTask(boardId);
            if (taskData) {
                board = {
                    id: taskData.id,
                    name: taskData.name_task,
                    todo: taskData.todo || [],
                    completed: taskData.completed || [],
                    userRole: determineUserRole(taskData)
                };
                cacheBoard(board);
            }
        } catch (error) {
            console.error('Error loading board:', error);
        }
    }

    if (board) {
        currentBoard = board;
        renderBoard();
        renderBoardsList();
        document.querySelector('h2').textContent = currentBoard.name;
        saveToLocalStorage();
    } else {
        console.error('Board not found');
    }
}

function determineUserRole(taskData) {
    const userEmail = localStorage.getItem('userEmail');
    if (taskData.owner && taskData.owner.email === userEmail) return 'Owner';
    if (taskData.manager && taskData.manager.email === userEmail) return 'Manager';
    if (taskData.employee && taskData.employee.email === userEmail) return 'Employee';
    return 'Viewer';
}

function renderBoard(taskData) {
    if (!currentBoard) {
        console.error('No current board');
        return;
    }

    boardContent.innerHTML = `
    <h2 class="fade-in-up">${currentBoard.name}</h2>
        <div class="board-content fade-in-up">
            <div class="mb-4 animate-fade-in">
                <h3>TO-DO</h3>
                <div class="table-responsive">
                    <table class="table table-bordered table-hover">
                        <thead>
                            <tr>
                                <th>Person</th>
                                <th>Planning Description</th>
                                <th>Planning File</th>
                                <th>Planning Due Date</th>
                                <th>Planning Status</th>
                                <th>Priority</th>
                                <th>Project File</th>
                                <th>Project Comment</th>
                                <th>Project Status</th>
                                <th>Project Due Date</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody id="todoTable">
                            ${renderTodoItems()}
                        </tbody>
                    </table>
                </div>
                <button class="btn btn-success add-item-btn" onclick="addNewItem('todo')">
                    <i class="fas fa-plus"></i> Add Item
                </button>
            </div>
            <div class="animate-fade-in">
                <h3>COMPLETED</h3>
                <div class="table-responsive">
                    <table class="table table-bordered table-hover">
                        <thead>
                            <tr>
                                <th>Person</th>
                                <th>Planning Description</th>
                                <th>Planning File</th>
                                <th>Planning Due Date</th>
                                <th>Planning Status</th>
                                <th>Priority</th>
                                <th>Project File</th>
                                <th>Project Comment</th>
                                <th>Project Status</th>
                                <th>Project Due Date</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody id="completedTable">
                            ${renderCompletedItems()}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    `;

    updateSidebarHeight();
}

// Fungsi untuk menyimpan data ke cache
function cacheData(key, data, expirationInMinutes = 5) {
    const now = new Date();
    const item = {
        value: data,
        expiry: now.getTime() + expirationInMinutes * 60000
    };
    localStorage.setItem(key, JSON.stringify(item));
}

// Fungsi untuk mengambil data dari cache
function getCachedData(key) {
    const itemStr = localStorage.getItem(key);
    if (!itemStr) {
        return null;
    }
    const item = JSON.parse(itemStr);
    const now = new Date();
    if (now.getTime() > item.expiry) {
        localStorage.removeItem(key);
        return null;
    }
    return item.value;
}

function updateSidebarHeight() {
    const sidebar = document.getElementById('sidebar');
    const main = document.querySelector('main');
    const darkModeSwitch = document.querySelector('#sidebar .form-check.form-switch');
    
    // Reset sidebar height
    sidebar.style.height = '100vh';
    
    // Calculate the available height for the sidebar content
    const availableHeight = sidebar.clientHeight - darkModeSwitch.offsetHeight;
    
    // Set the height of the boardsList
    const boardsList = document.getElementById('boardsList');
    boardsList.style.maxHeight = `${availableHeight - boardsList.offsetTop}px`;
    
    // Adjust sidebar height if main content is taller
    if (main.scrollHeight > sidebar.clientHeight) {
        sidebar.style.height = `${main.scrollHeight}px`;
    }
}

window.addEventListener('resize', updateSidebarHeight);

function showDeleteBoardConfirmModal(boardId, boardName) {
    const deleteBoardConfirmModal = new bootstrap.Modal(document.getElementById('deleteBoardConfirmModal'));
    const deleteBoardNameSpan = document.getElementById('deleteBoardName');
    const confirmDeleteBoardBtn = document.getElementById('confirmDeleteBoard');
    
    deleteBoardNameSpan.textContent = boardName;
    
    confirmDeleteBoardBtn.onclick = () => {
        deleteBoard(boardId);
        deleteBoardConfirmModal.hide();
    };
    
    const handleEnterKey = (event) => {
        if (event.key === 'Enter') {
            event.preventDefault();
            confirmDeleteBoardBtn.click();
        }
    };
    
    deleteBoardConfirmModal._element.addEventListener('shown.bs.modal', () => {
        document.addEventListener('keydown', handleEnterKey);
    });
    
    deleteBoardConfirmModal._element.addEventListener('hidden.bs.modal', () => {
        document.removeEventListener('keydown', handleEnterKey);
    });
    
    deleteBoardConfirmModal.show();
}

function deleteBoard(boardId) {
    const boardIndex = boards.findIndex(board => board.id === boardId);
    if (boardIndex !== -1) {
        boards.splice(boardIndex, 1);
        renderBoardsList();
        if (currentBoard && currentBoard.id === boardId) {
            currentBoard = null;
            boardContent.innerHTML = '<h2>No board selected</h2>';
        }
    }
}

function renderTodoItems() {
    return currentBoard.todo.map((item, index) => `
        <tr>
        <td>${item.editMode ? `<input type="text" class="form-control name-task" value="${item.nameTask || ''}" onchange="updateItem(${index}, 'todo', 'nameTask', this.value)">` : (item.nameTask || '')}</td>
        <td>
            <div class="avatar-container">
                ${renderPersons(item.persons, 'todo', index)}
                ${item.editMode ? `<span class="add-person-btn" onclick="showAddPersonModal(${index}, 'todo')"><i class="fas fa-plus"></i></span>` : ''}
            </div>
        </td>
            <td>${item.editMode ? `<input type="text" class="form-control planning-description" value="${item.planningDescription || ''}" onchange="updateItem(${index}, 'todo', 'planningDescription', this.value)">` : (item.planningDescription || '')}</td>
            <td>
                ${item.editMode ? `
                    <div class="custom-file-upload">
                        <label for="planningFile-${index}-input" class="btn btn-outline-secondary btn-sm">Choose Files</label>
                        <input type="file" id="planningFile-${index}-input" class="form-control" multiple onchange="updateFiles(${index}, 'todo', 'planningFile', this.files)" style="display: none;">
                    </div>
                ` : ''}
                <div class="file-list" id="planningFile-list-${index}">
                    ${renderFileList(item.planningFile, index, 'todo', 'planningFile', item.editMode)}
                </div>
            </td>
            <td>${item.editMode ? `<input type="date" class="form-control" value="${item.planningDueDate || ''}" onchange="updateItem(${index}, 'todo', 'planningDueDate', this.value)">` : (item.planningDueDate || '')}</td>
            <td>${item.editMode ? `
                <select class="form-select" onchange="updateItem(${index}, 'todo', 'planningStatus', this.value)">
                    <option value="" ${!item.planningStatus ? 'selected' : ''}>Select Status</option>
                    <option value="Not Approved" ${item.planningStatus === 'Not Approved' ? 'selected' : ''}>Not Approved</option>
                    <option value="Approved" ${item.planningStatus === 'Approved' ? 'selected' : ''}>Approved</option>
                </select>
            ` : (item.planningStatus || '')}</td>
            <td>${item.editMode ? `
                <select class="form-select priority" onchange="updateItem(${index}, 'todo', 'priority', this.value)">
                    <option value="" ${!item.priority ? 'selected' : ''}>Select Priority</option>
                    <option value="Low" ${item.priority === 'Low' ? 'selected' : ''}>Low</option>
                    <option value="Medium" ${item.priority === 'Medium' ? 'selected' : ''}>Medium</option>
                    <option value="High" ${item.priority === 'High' ? 'selected' : ''}>High</option>
                </select>
            ` : (item.priority || '')}</td>
            <td>
                ${item.editMode ? `
                    <div class="custom-file-upload">
                        <label for="projectFile-${index}-input" class="btn btn-outline-secondary btn-sm">Choose Files</label>
                        <input type="file" id="projectFile-${index}-input" class="form-control" multiple onchange="updateFiles(${index}, 'todo', 'projectFile', this.files)" style="display: none;">
                    </div>
                ` : ''}
                <div class="file-list" id="projectFile-list-${index}">
                    ${renderFileList(item.projectFile, index, 'todo', 'projectFile', item.editMode)}
                </div>
            </td>
            <td>${item.editMode ? `<input type="text" class="form-control project-comment" value="${item.projectComment || ''}" onchange="updateItem(${index}, 'todo', 'projectComment', this.value)">` : (item.projectComment || '')}</td>
            <td>${item.editMode ? `
                <select class="form-select projectStatus" onchange="updateItem(${index}, 'todo', 'projectStatus', this.value)">
                    <option value="" ${!item.projectStatus ? 'selected' : ''}>Select Status</option>
                    <option value="Undone" ${item.projectStatus === 'Undone' ? 'selected' : ''}>Undone</option>
                    <option value="Working" ${item.projectStatus === 'Working' ? 'selected' : ''}>Working</option>
                    <option value="Done" ${item.projectStatus === 'Done' ? 'selected' : ''}>Done</option>
                </select>
            ` : (item.projectStatus || '')}</td>
            <td>${item.editMode ? `<input type="date" class="form-control" value="${item.projectDueDate || ''}" onchange="updateItem(${index}, 'todo', 'projectDueDate', this.value)">` : (item.projectDueDate || '')}</td>
            <td>
                <div class="action-buttons">
                    ${item.editMode ? 
                        `<button class="btn btn-primary btn-sm" onclick="saveItem(${index}, 'todo')"><i class="fas fa-save"></i> Save</button>` :
                        `<button class="btn btn-warning btn-sm" onclick="editItem(${index}, 'todo')"><i class="fas fa-edit"></i> Edit</button>`
                    }
                    <button class="btn btn-danger btn-sm" onclick="deleteItem(${index}, 'todo')"><i class="fas fa-trash"></i> Delete</button>
                </div>
            </td>
        </tr>
    `).join('');
}

function handleInputKeydown(event, index, listType, field) {
    if (event.key === 'Enter') {
        event.preventDefault();
        updateItem(index, listType, field, event.target.value);
        event.target.blur();
    }
}

async function updateFiles(index, listType, field, newFiles) {
    for (const file of newFiles) {
        const formData = new FormData();
        formData.append(field, file);

        try {
            const response = await axios.put(`http://127.0.0.1:8080/task/${currentBoard.id}`, formData, {
                headers: { 'Content-Type': 'multipart/form-data' }
            });

            if (response.data.code === 200) {
                updateItemFromResponse(index, listType, response.data.data);
            }
        } catch (error) {
            if (error.response && error.response.data && error.response.data.error === "File already exist") {
                alert(`File ${file.name} already exists.`);
            } else {
                alert('An error occurred while uploading the file.');
            }
            console.error('Error uploading file:', error);
        }
    }
    renderBoard();
}

function renderFileList(files, itemIndex, listType, field, showRemoveButton = true) {
    if (!files || files.length === 0) return '';
    return files.map((file, idx) => `
        <div>
            ${file.name}
            ${showRemoveButton ? `
                <button type="button" class="btn btn-sm btn-danger" onclick="removeFile(${itemIndex}, '${listType}', '${field}', ${idx})">x</button>
            ` : ''}
        </div>
    `).join('');
}

function removeFile(itemIndex, listType, field, fileIndex) {
    currentBoard[listType][itemIndex][field].splice(fileIndex, 1);
    document.getElementById(`${field}-list-${itemIndex}`).innerHTML = renderFileList(currentBoard[listType][itemIndex][field], itemIndex, listType, field);
}

function renderCompletedItems() {
    return currentBoard.completed.map((item, index) => `
        <tr>
            <td>
                <div class="avatar-container">
                    ${renderPersons(item.persons, 'completed')}
                </div>
            </td>
            <td>${item.planningDescription || ''}</td>
            <td>${renderFileList(item.planningFile, index, 'completed', 'planningFile', false)}</td>
            <td>${item.planningDueDate || ''}</td>
            <td><span class="status-pill status-${item.planningStatus.toLowerCase().replace(' ', '-')}">${item.planningStatus}</span></td>
            <td><span class="status-pill priority-${item.priority.toLowerCase()}">${item.priority}</span></td>
            <td>${renderFileList(item.projectFile, index, 'completed', 'projectFile', false)}</td>
            <td>${item.projectComment || ''}</td>
            <td><span class="status-pill project-status-${item.projectStatus.toLowerCase()}">${item.projectStatus}</span></td>
            <td>${item.projectDueDate || ''}</td>
            <td>
                <button class="btn btn-primary btn-sm" onclick="moveToTodo(${index})"><i class="fas fa-arrow-left"></i> Move to TO-DO</button>
            </td>
        </tr>
    `).join('');
}

function canEditField(field) {
    const role = currentBoard.userRole;
    const ownerFields = ['planningDescription', 'planningDueDate', 'planningStatus', 'projectStatus', 'projectDueDate', 'manager'];
    const managerFields = ['planningFile', 'priority', 'employee'];
    const employeeFields = ['projectFile', 'projectComment'];

    if (role === 'Owner') return ownerFields.includes(field);
    if (role === 'Manager') return managerFields.includes(field);
    if (role === 'Employee') return employeeFields.includes(field);
    return false;
}

function renderPersons(persons, listType, itemIndex) {
    if (!persons || !Array.isArray(persons) || persons.length === 0) {
        // Tambahkan owner secara otomatis jika tidak ada person
        const ownerEmail = localStorage.getItem('userEmail');
        persons = [{ email: ownerEmail || 'unknown@example.com', role: 'Owner', roleDisplay: 'Owner' }];
    }
    return persons.map((person, personIndex) => `
        <span class="avatar ${person.role.toLowerCase()}" 
              data-bs-toggle="tooltip" 
              data-bs-placement="top" 
              title="${person.email} (${person.roleDisplay || person.role})"
              onmouseover="showPersonInfo(event, '${person.email}', '${person.roleDisplay || person.role}', ${personIndex}, '${listType}', ${itemIndex})"
              onmouseout="hidePersonInfo()">
            ${person.email ? person.email.charAt(0).toUpperCase() : 'U'}
            ${listType === 'todo' ? `
                <div class="person-actions" style="display: none;">
                    <button class="edit-person-btn" onclick="showEditPersonModal(${itemIndex}, ${personIndex}, '${listType}')">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="delete-person-btn" onclick="showDeletePersonModal('${person.email}', ${itemIndex}, ${personIndex}, '${listType}')">
                        <i class="fas fa-times"></i>
                    </button>
                </div>
            ` : ''}
        </span>
    `).join('');
}

function showDeletePersonModal(email, itemIndex, personIndex, listType) {
    const modal = document.getElementById('deletePersonModal');
    const modalContent = modal.querySelector('.modal-body');
    modalContent.innerHTML = `Are you sure you want to delete ${email}?`;
    
    function deletePerson() {
        currentBoard[listType][itemIndex].persons.splice(personIndex, 1);
        bootstrap.Modal.getInstance(modal).hide();
        renderBoard();
    }
    
    const confirmBtn = modal.querySelector('#confirmDeletePerson');
    confirmBtn.onclick = deletePerson;
    
    // Add keydown event listener
    modal.addEventListener('keydown', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            deletePerson();
        }
    });
    
    new bootstrap.Modal(modal).show();
}

function showPersonInfo(event, email, roleDisplay, personIndex, listType, itemIndex) {
    const avatar = event.target;
    if (listType === 'todo') {
        const actions = avatar.querySelector('.person-actions');
        if (actions) {
            actions.style.display = 'flex';
        }
    }

    const infoDiv = document.createElement('div');
    infoDiv.className = 'person-info';
    infoDiv.innerHTML = `
        <p><strong>Email:</strong> ${email}</p>
        <p><strong>Role:</strong> ${roleDisplay}</p>
    `;
    infoDiv.style.left = `${event.pageX}px`;
    infoDiv.style.top = `${event.pageY}px`;
    document.body.appendChild(infoDiv);
}

function showEditPersonModal(itemIndex, personIndex, listType) {
    const editPersonModal = new bootstrap.Modal(document.getElementById('editPersonModal'));
    const item = currentBoard[listType][itemIndex];
    if (item && item.persons && item.persons[personIndex]) {
        const person = item.persons[personIndex];
        
        const emailInput = document.getElementById('editPersonEmail');
        const roleInput = document.getElementById('editPersonRole');
        const roleDisplayInput = document.getElementById('editPersonRoleDisplay');
        
        emailInput.value = person.email;
        roleInput.value = person.role;
        roleDisplayInput.value = person.roleDisplay || '';
        
        function savePerson() {
            saveEditPerson(itemIndex, personIndex, listType);
        }
        
        const saveEditPersonBtn = document.getElementById('saveEditPerson');
        saveEditPersonBtn.onclick = savePerson;
        
        // Add keydown event listener
        editPersonModal._element.addEventListener('keydown', function(event) {
            if (event.key === 'Enter') {
                event.preventDefault();
                savePerson();
            }
        });
        
        editPersonModal.show();
        
        // Focus on the email input when the modal is shown
        editPersonModal._element.addEventListener('shown.bs.modal', function() {
            emailInput.focus();
        });
    } else {
        console.error('Person not found');
    }
}

function saveEditPerson(itemIndex, personIndex, listType) {
    const email = document.getElementById('editPersonEmail').value;
    const role = document.getElementById('editPersonRole').value;
    const roleDisplay = document.getElementById('editPersonRoleDisplay').value;
    
    if (email && validateEmail(email) && role) {
        if (!currentBoard[listType][itemIndex].persons) {
            currentBoard[listType][itemIndex].persons = [];
        }
        
        if (!currentBoard[listType][itemIndex].persons[personIndex]) {
            currentBoard[listType][itemIndex].persons[personIndex] = {};
        }
        
        const person = currentBoard[listType][itemIndex].persons[personIndex];
        
        person.email = email;
        person.role = role;
        
        if (roleDisplay) {
            person.roleDisplay = roleDisplay;
            customRoleNames[role] = roleDisplay;
        } else {
            delete person.roleDisplay;
        }
        
        bootstrap.Modal.getInstance(document.getElementById('editPersonModal')).hide();
        renderBoard();
    }
}

document.getElementById('editPersonEmail').addEventListener('input', function() {
    const emailInput = this;
    const saveButton = document.getElementById('saveEditPerson');
    const errorDiv = document.getElementById('editEmailError');

    if (emailInput.value === '' || validateEmail(emailInput.value)) {
        emailInput.classList.remove('is-invalid');
        emailInput.classList.add('is-valid');
        errorDiv.textContent = '';
        saveButton.disabled = false;
    } else {
        emailInput.classList.remove('is-valid');
        emailInput.classList.add('is-invalid');
        errorDiv.textContent = 'Please enter a valid email address.';
        saveButton.disabled = true;
    }
});

function hidePersonInfo() {
    const infoDiv = document.querySelector('.person-info');
    if (infoDiv) {
        infoDiv.remove();
    }

    const actionButtons = document.querySelectorAll('.person-actions');
    actionButtons.forEach(btn => {
        btn.style.display = 'none';
    });
}

// remove person
function confirmDeletePerson(event, email, personIndex) {
    event.stopPropagation(); // Prevent the avatar's mouseover event from firing
    const confirmed = confirm(`Are you sure you want to delete ${email}?`);
    if (confirmed) {
        deletePerson(personIndex);
    }
}

function deletePerson(personIndex) {
    const itemIndex = findItemIndexByPersonIndex(personIndex);
    if (itemIndex !== -1) {
        currentBoard.todo[itemIndex].persons.splice(personIndex, 1);
        renderBoard();
    }
}

function findItemIndexByPersonIndex(personIndex) {
    let currentIndex = 0;
    for (let i = 0; i < currentBoard.todo.length; i++) {
        const item = currentBoard.todo[i];
        if (currentIndex + item.persons.length > personIndex) {
            return i;
        }
        currentIndex += item.persons.length;
    }
    return -1;
}

async function addNewItem() {
    const { token } = checkAuth();
    const taskName = document.getElementById('newItemName').value.trim();
    if (taskName) {
        try {
            const formData = new FormData();
            formData.append('newItemName', taskName);

            const response = await axios.post(`http://127.0.0.1:8080/task`, 
                formData,
                { 
                    headers: { 
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'multipart/form-data'
                    },
                    withCredentials: true
                }
            );
            if (response.data.code === 200) {
                const newItem = {
                    id: response.data.data.task_id,
                    nameTask: response.data.data.name_task,
                    persons: [{
                        email: response.data.data.user_email,
                        role: 'Owner',
                        roleDisplay: 'Owner'
                    }],
                    // ... (other properties)
                };
                currentBoard.todo.push(newItem);
                renderBoard();
                updateSidebarHeight();
                bootstrap.Modal.getInstance(document.getElementById('addItemModal')).hide();
            }
        } catch (error) {
            console.error('Error adding new item:', error);
        }
    }
}

document.getElementById('saveNewItem').addEventListener('click', addNewItem);

function saveToLocalStorage() {
    if (currentBoard) {
        const boardToSave = {...currentBoard};
        boardToSave.todo = boardToSave.todo.map(item => ({...item, editMode: false}));
        boardToSave.completed = boardToSave.completed.map(item => ({...item, editMode: false}));
        localStorage.setItem(`board_${currentBoard.id}`, JSON.stringify(boardToSave));
    }
}

function loadFromLocalStorage() {
    const boardId = new URLSearchParams(window.location.search).get('boardId');
    if (boardId) {
        const savedBoard = localStorage.getItem(`board_${boardId}`);
        if (savedBoard) {
            currentBoard = JSON.parse(savedBoard);
            currentBoard.todo = currentBoard.todo.map(item => ({...item, editMode: false}));
            currentBoard.completed = currentBoard.completed.map(item => ({...item, editMode: false}));
            renderBoard();
        } else {
            loadBoard(boardId);
        }
    }
}

window.addEventListener('load', loadFromLocalStorage);

function updateItem(index, listType, field, value) {
    if (field === 'planningDueDate' || field === 'projectDueDate') {
        // Konversi format tanggal ke ISO string
        value = new Date(value).toISOString().split('T')[0];
    }

    if (field !== 'planningFile' && field !== 'projectFile') {
        currentBoard[listType][index][field] = value;
    } else {
        // Handle file uploads separately
        currentBoard[listType][index][field] = value;
    }
    checkAndMoveItem(index, listType);
    sortTodoItems();
    renderBoard();
    saveToLocalStorage();
}

function checkAndMoveItem(index, listType) {
    const item = currentBoard[listType][index];
    if (item.planningStatus === 'Approved' && item.projectStatus === 'Done') {
        currentBoard.completed.push({...item});
        currentBoard.todo.splice(index, 1);
    }
}

function sortTodoItems() {
    currentBoard.todo.sort((a, b) => {
        const priorityOrder = { 'High': 0, 'Medium': 1, 'Low': 2 };
        return priorityOrder[a.priority] - priorityOrder[b.priority];
    });
}

async function saveItem(index, listType) {
    const item = currentBoard[listType][index];
    const changedData = getChangedData(item);
    const changedPersons = getChangedPersons(item.persons || [], originalItemData.persons || []);
    
    if (Object.keys(changedData).length === 0 && changedPersons.added.length === 0 && changedPersons.removed.length === 0) {
        console.log("No changes detected");
        return;
    }

    const formData = new FormData();

    for (const [key, value] of Object.entries(changedData)) {
        if (value !== null && value !== undefined && value !== '') {
            if (key === 'planningFile' || key === 'projectFile') {
                if (Array.isArray(value) && value.length > 0) {
                    value.forEach(file => {
                        if (file instanceof File) {
                            formData.append(key === 'planningFile' ? 'planning_file' : 'project_file', file);
                        }
                    });
                }
            } else {
                const backendKey = key.replace(/([A-Z])/g, '_$1').toLowerCase();
                formData.append(backendKey, value);
            }
        }
    }

    changedPersons.added.forEach(person => {
        if (person.role === 'Manager') {
            formData.append('manager', person.email);
        } else if (person.role === 'Employee') {
            formData.append('employee', person.email);
        }
    });

    try {
        const response = await axios.put(`http://127.0.0.1:8080/board/${currentBoard.id}/task/${item.id}`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
            withCredentials: true
        });

        if (response.data.code === 200) {
            updateItemFromResponse(index, listType, response.data.data);
            item.editMode = false;
            renderBoard();
            saveToLocalStorage();
            alert('Item saved successfully!');
        }
    } catch (error) {
        console.error('Error saving item:', error);
        if (error.response) {
            console.error('Response data:', error.response.data);
            console.error('Response status:', error.response.status);
            if (error.response.data.error === "Only for owner") {
                alert("Only the owner can modify this field.");
            } else if (error.response.data.error === "File already exist") {
                alert("This file already exists.");
            } else if (error.response.data.error === "User is already assigned as manager to a task") {
                alert("This user is already assigned as a manager to this task.");
            } else if (error.response.data.error === "User is already assigned as employee to a task") {
                alert("This user is already assigned as an employee to this task.");
            } else if (error.response.data.error === "user not found") {
                alert("One or more users were not found. Please check the email addresses and try again.");
            } else {
                alert('An error occurred while saving the item.');
            }
        }
    }
}

function updateItemFromResponse(index, listType, responseData) {
    const item = currentBoard[listType][index];
    for (const [key, value] of Object.entries(responseData)) {
        if (key === 'planning_file' || key === 'project_file') {
            if (value && value.file_name) {
                item[key] = [{ name: value.file_name, url: value.file_url }];
            }
        } else if (key === 'manager' || key === 'employee') {
            if (value && value.email) {
                const existingPersonIndex = item.persons.findIndex(p => p.role === (key === 'manager' ? 'Manager' : 'Employee'));
                if (existingPersonIndex !== -1) {
                    item.persons[existingPersonIndex] = {
                        email: value.email,
                        role: key === 'manager' ? 'Manager' : 'Employee',
                        roleDisplay: key === 'manager' ? 'Manager' : 'Employee'
                    };
                } else {
                    item.persons.push({
                        email: value.email,
                        role: key === 'manager' ? 'Manager' : 'Employee',
                        roleDisplay: key === 'manager' ? 'Manager' : 'Employee'
                    });
                }
            }
        } else {
            const camelKey = key.replace(/_([a-z])/g, (g) => g[1].toUpperCase());
            item[camelKey] = value === 'Undone' ? item[camelKey] : value;
        }
    }
}

let originalItemData = {};

function editItem(index, listType) {
    const item = currentBoard[listType][index];
    originalItemData = JSON.parse(JSON.stringify(item)); // Deep copy
    item.editMode = true;
    renderBoard();
}

function getChangedData(item) {
    const changedData = {};
    for (const [key, value] of Object.entries(item)) {
        if (key !== 'persons' && key !== 'editMode') {
            if (JSON.stringify(value) !== JSON.stringify(originalItemData[key])) {
                if (key === 'planningDueDate' || key === 'projectDueDate') {
                    changedData[key] = new Date(value).toISOString().split('T')[0];
                } else {
                    changedData[key] = value;
                }
            }
        }
    }
    return changedData;
}

function getChangedPersons(newPersons, originalPersons) {
    const changedPersons = {
        added: [],
        removed: []
    };

    newPersons.forEach(person => {
        if (!originalPersons.some(p => p.email === person.email && p.role === person.role)) {
            changedPersons.added.push(person);
        }
    });

    originalPersons.forEach(person => {
        if (!newPersons.some(p => p.email === person.email && p.role === person.role)) {
            changedPersons.removed.push(person);
        }
    });

    return changedPersons;
}

function showDeleteConfirmModal(index, listType) {
    const deleteConfirmModal = new bootstrap.Modal(document.getElementById('deleteConfirmModal'));
    const confirmDeleteBtn = document.getElementById('confirmDelete');
    
    function deleteItem() {
        currentBoard[listType].splice(index, 1);
        deleteConfirmModal.hide();
        renderBoard();
    }
    
    confirmDeleteBtn.onclick = deleteItem;
    
    // Add keydown event listener
    deleteConfirmModal._element.addEventListener('keydown', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            deleteItem();
        }
    });
    
    deleteConfirmModal.show();
}

function deleteItem(index, listType) {
    showDeleteConfirmModal(index, listType);
    renderBoard();
    updateSidebarHeight();
}

function moveToTodo(index) {
    const item = currentBoard.completed[index];
    item.planningStatus = 'Not Approved';
    item.projectStatus = 'Undone';
    currentBoard.todo.push(item);
    currentBoard.completed.splice(index, 1);
    sortTodoItems();
    renderBoard();
}

document.getElementById('newPersonEmail').addEventListener('keydown', handleNewPersonKeydown);
document.getElementById('newPersonRole').addEventListener('keydown', handleNewPersonKeydown);
document.getElementById('newPersonRoleDisplay').addEventListener('keydown', handleNewPersonKeydown);

function handleNewPersonKeydown(event) {
    if (event.key === 'Enter') {
        event.preventDefault();
        const saveNewPersonBtn = document.getElementById('saveNewPerson');
        if (!saveNewPersonBtn.disabled) {
            saveNewPersonBtn.click();
        }
    }
}

function showAddPersonModal(itemIndex, listType) {
    const addPersonModal = document.getElementById('addPersonModal');
    addPersonModal.dataset.itemIndex = itemIndex;
    addPersonModal.dataset.listType = listType;

    const newPersonEmailInput = document.getElementById('newPersonEmail');
    const newPersonRoleInput = document.getElementById('newPersonRole');
    const newPersonRoleDisplayInput = document.getElementById('newPersonRoleDisplay');

    // Reset form fields
    newPersonEmailInput.value = '';
    newPersonRoleInput.value = 'Manager'; // Change default to Manager
    newPersonRoleDisplayInput.value = '';

    // Remove the 'Owner' option from the role select
    const ownerOption = newPersonRoleInput.querySelector('option[value="Owner"]');
    if (ownerOption) {
        ownerOption.remove();
    }

    newPersonRoleInput.addEventListener('change', function() {
        newPersonRoleDisplayInput.value = customRoleNames[this.value] !== this.value ? customRoleNames[this.value] : '';
    });

    new bootstrap.Modal(addPersonModal).show();
    addPersonModal.addEventListener('shown.bs.modal', function onModalShown() {
        newPersonEmailInput.focus();
        addPersonModal.removeEventListener('shown.bs.modal', onModalShown);
    });
}

function addNewPerson() {
    const addPersonModal = document.getElementById('addPersonModal');
    const itemIndex = parseInt(addPersonModal.dataset.itemIndex);
    const listType = addPersonModal.dataset.listType;

    const email = document.getElementById('newPersonEmail').value;
    const role = document.getElementById('newPersonRole').value;
    const roleDisplay = document.getElementById('newPersonRoleDisplay').value || role;
    
    if (validateEmail(email) && role) {
        if (!isNaN(itemIndex) && listType) {
            if (!currentBoard[listType][itemIndex].persons) {
                currentBoard[listType][itemIndex].persons = [];
            }
            
            // Check for existing person with same role
            const existingPerson = currentBoard[listType][itemIndex].persons.find(p => p.role === role);
            if (existingPerson) {
                alert(`A ${role} is already assigned to this task.`);
                return;
            }

            currentBoard[listType][itemIndex].persons.push({
                email: email,
                role: role,
                roleDisplay: roleDisplay
            });

            bootstrap.Modal.getInstance(addPersonModal).hide();
            renderBoard();
        }
    }
}

addPersonModal._element.addEventListener('hidden.bs.modal', function () {
    document.getElementById('newPersonEmail').value = '';
    document.getElementById('newPersonRole').value = 'Owner';
    document.getElementById('newPersonRoleDisplay').value = '';
    document.getElementById('emailError').textContent = '';
    document.getElementById('saveNewPerson').disabled = true;
});

document.getElementById('newPersonEmail').addEventListener('input', function() {
    const emailInput = this;
    const saveButton = document.getElementById('saveNewPerson');
    const errorDiv = document.getElementById('emailError');

    if (validateEmail(emailInput.value)) {
        emailInput.classList.remove('is-invalid');
        emailInput.classList.add('is-valid');
        errorDiv.textContent = '';
        saveButton.disabled = false;
    } else {
        emailInput.classList.remove('is-valid');
        emailInput.classList.add('is-invalid');
        errorDiv.textContent = 'Please enter a valid email address.';
        saveButton.disabled = true;
    }
});

addPersonModal._element.addEventListener('show.bs.modal', function() {
    const emailInput = document.getElementById('newPersonEmail');
    const errorDiv = document.getElementById('emailError');
    const saveButton = document.getElementById('saveNewPerson');

    emailInput.value = '';
    emailInput.classList.remove('is-invalid', 'is-valid');
    errorDiv.textContent = '';
    saveButton.disabled = true;
});

function validateEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

function searchBoards() {
    const searchTerm = searchBoardInput.value.toLowerCase();
    const filteredBoards = boards.filter(board => board.name.toLowerCase().includes(searchTerm));
    renderFilteredBoards(filteredBoards);
}

function renderFilteredBoards(filteredBoards) {
    boardsList.innerHTML = '';
    filteredBoards.forEach(board => {
        const li = document.createElement('li');
        li.className = 'nav-item';
        li.innerHTML = `<a class="nav-link" href="#" data-board-id="${board.id}">${board.name}</a>`;
        li.querySelector('a').addEventListener('click', () => loadBoard(board.id));
        boardsList.appendChild(li);
    });
}

// Initial render
renderBoardsList();