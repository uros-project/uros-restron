// Actor管理 JavaScript

let allActors = [];

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', function() {
    loadActors();
});

// 加载所有Actor
async function loadActors() {
    try {
        const response = await fetch('/api/v1/actors');
        const data = await response.json();
        allActors = data.data || [];
        renderActors(allActors);
        updateActorSelect(allActors);
    } catch (error) {
        console.error('加载Actor失败:', error);
        alert('加载Actor失败: ' + error.message);
    }
}

// 渲染Actor列表
function renderActors(actors) {
    const grid = document.getElementById('actorsGrid');
    grid.innerHTML = '';

    actors.forEach(actor => {
        const card = createActorCard(actor);
        grid.appendChild(card);
    });
}

// 创建Actor卡片
function createActorCard(actor) {
    const card = document.createElement('div');
    card.className = 'actor-card';
    
    const functionsHtml = actor.functions && actor.functions.length > 0 
        ? actor.functions.map(func => `<span class="function-tag">${func}</span>`).join('')
        : '<span style="color: #7f8c8d; font-style: italic;">暂无函数</span>';

    const lastActive = new Date(actor.lastActive).toLocaleString();

    card.innerHTML = `
        <div class="card-header">
            <h3 class="card-title">${actor.name || '未命名Actor'}</h3>
            <span class="card-status status-${actor.status}">${actor.status}</span>
        </div>
        <div class="actor-content">
            <p class="actor-id">ID: ${actor.id}</p>
            <p class="actor-last-active">最后活跃: ${lastActive}</p>
            <div class="actor-functions">
                <strong>可用函数:</strong>
                <div class="functions-list">${functionsHtml}</div>
            </div>
        </div>
        <div class="card-actions">
            <button class="btn btn-secondary" onclick="viewActorDetail('${actor.id}')">查看详情</button>
            <button class="btn btn-primary" onclick="callActorFunction('${actor.id}')">调用函数</button>
        </div>
    `;
    
    return card;
}

// 更新Actor选择框
function updateActorSelect(actors) {
    const select = document.getElementById('actorSelect');
    select.innerHTML = '<option value="">请选择Actor</option>';
    
    actors.forEach(actor => {
        const option = document.createElement('option');
        option.value = actor.id;
        option.textContent = `${actor.name || '未命名Actor'} (${actor.id.substring(0, 8)}...)`;
        select.appendChild(option);
    });
}

// 过滤Actor
function filterActors() {
    const status = document.getElementById('statusFilter').value;
    const search = document.getElementById('searchInput').value.toLowerCase();
    
    let filtered = allActors;
    
    if (status) {
        filtered = filtered.filter(actor => actor.status === status);
    }
    
    if (search) {
        filtered = filtered.filter(actor => 
            (actor.name && actor.name.toLowerCase().includes(search)) ||
            actor.id.toLowerCase().includes(search)
        );
    }
    
    renderActors(filtered);
}

// 搜索Actor
function searchActors() {
    filterActors();
}

// 刷新Actor状态
function refreshActors() {
    loadActors();
}

// 加载Actor函数
async function loadActorFunctions() {
    const actorId = document.getElementById('actorSelect').value;
    const functionSelect = document.getElementById('functionSelect');
    functionSelect.innerHTML = '<option value="">请选择函数</option>';
    
    if (!actorId) {
        return;
    }
    
    try {
        const response = await fetch(`/api/v1/actors/${actorId}/functions`);
        const data = await response.json();
        
        if (data.functions && data.functions.length > 0) {
            data.functions.forEach(func => {
                const option = document.createElement('option');
                option.value = func;
                option.textContent = func;
                functionSelect.appendChild(option);
            });
        }
    } catch (error) {
        console.error('加载Actor函数失败:', error);
        alert('加载Actor函数失败: ' + error.message);
    }
}

// 调用Actor函数
async function callActorFunction(actorId = null) {
    const selectedActorId = actorId || document.getElementById('actorSelect').value;
    const functionName = document.getElementById('functionSelect').value;
    const parametersText = document.getElementById('functionParams').value;
    
    if (!selectedActorId || !functionName) {
        alert('请选择Actor和函数');
        return;
    }
    
    let parameters = {};
    if (parametersText && parametersText.trim()) {
        try {
            parameters = JSON.parse(parametersText);
        } catch (error) {
            alert('参数格式错误，请输入有效的JSON');
            return;
        }
    }
    
    try {
        const response = await fetch(`/api/v1/actors/${selectedActorId}/functions/${functionName}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ params: parameters })
        });
        
        const result = await response.json();
        displayFunctionResult(result);
    } catch (error) {
        console.error('调用Actor函数失败:', error);
        alert('调用函数失败: ' + error.message);
    }
}

// 显示函数调用结果
function displayFunctionResult(result) {
    const container = document.getElementById('functionResult');
    container.innerHTML = `
        <h4>函数调用结果</h4>
        <pre class="result-json">${JSON.stringify(result, null, 2)}</pre>
    `;
}

// 查看Actor详情
async function viewActorDetail(id) {
    try {
        const response = await fetch(`/api/v1/actors/${id}`);
        const actor = await response.json();
        
        const content = document.getElementById('actorDetailContent');
        content.innerHTML = `
            <div class="actor-detail">
                <h4>${actor.name || '未命名Actor'}</h4>
                <p><strong>ID:</strong> ${actor.id}</p>
                <p><strong>状态:</strong> <span class="status-${actor.status}">${actor.status}</span></p>
                <p><strong>最后活跃:</strong> ${new Date(actor.lastActive).toLocaleString()}</p>
                
                <h5>可用函数</h5>
                <div class="functions-list">
                    ${actor.functions && actor.functions.length > 0 
                        ? actor.functions.map(func => `<span class="function-tag">${func}</span>`).join('')
                        : '<p>暂无函数</p>'}
                </div>
                
                <h5>Actor状态详情</h5>
                <pre class="actor-status-json">${JSON.stringify(actor, null, 2)}</pre>
            </div>
        `;
        
        document.getElementById('actorDetailModal').style.display = 'block';
    } catch (error) {
        console.error('加载Actor详情失败:', error);
        alert('加载详情失败: ' + error.message);
    }
}

// 关闭Actor详情模态框
function closeActorDetailModal() {
    document.getElementById('actorDetailModal').style.display = 'none';
}

// 显示健康检查
async function showActorHealthCheck() {
    try {
        const response = await fetch('/api/v1/actors/health');
        const healthData = await response.json();
        
        const content = document.getElementById('healthCheckContent');
        content.innerHTML = `
            <div class="health-check">
                <h4>系统健康状态</h4>
                <div class="health-status">
                    <p><strong>状态:</strong> <span class="status-${healthData.status}">${healthData.status}</span></p>
                    <p><strong>活跃Actor数量:</strong> ${healthData.actors}</p>
                    <p><strong>检查时间:</strong> ${healthData.timestamp}</p>
                </div>
                <h5>详细状态</h5>
                <pre class="health-json">${JSON.stringify(healthData, null, 2)}</pre>
            </div>
        `;
        
        document.getElementById('healthCheckModal').style.display = 'block';
    } catch (error) {
        console.error('健康检查失败:', error);
        alert('健康检查失败: ' + error.message);
    }
}

// 关闭健康检查模态框
function closeHealthCheckModal() {
    document.getElementById('healthCheckModal').style.display = 'none';
}

// 点击模态框外部关闭
window.onclick = function(event) {
    const detailModal = document.getElementById('actorDetailModal');
    const healthModal = document.getElementById('healthCheckModal');
    
    if (event.target === detailModal) {
        closeActorDetailModal();
    }
    if (event.target === healthModal) {
        closeHealthCheckModal();
    }
}
