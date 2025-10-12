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
        const result = await response.json();
        // API 返回格式: { success: true, data: { count: N, data: [...] } }
        allActors = result.data?.data || [];
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
    card.className = 'bg-white rounded-xl shadow-md hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-200';
    
    const functionsHtml = actor.functions && actor.functions.length > 0 
        ? actor.functions.map(func => `<span class="inline-block px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-xs font-medium mr-2 mb-2">${func}</span>`).join('')
        : '<span class="text-gray-400 italic text-sm">暂无函数</span>';

    const lastActive = actor.lastActive ? new Date(actor.lastActive).toLocaleString() : '未知';
    
    // 状态颜色映射
    const statusColors = {
        'running': 'bg-green-100 text-green-800',
        'stopped': 'bg-gray-100 text-gray-800',
        'error': 'bg-red-100 text-red-800',
        'idle': 'bg-yellow-100 text-yellow-800'
    };
    const statusClass = statusColors[actor.status] || 'bg-gray-100 text-gray-800';
    
    // 状态图标
    const statusIcons = {
        'running': '🟢',
        'stopped': '⚫',
        'error': '🔴',
        'idle': '🟡'
    };
    const statusIcon = statusIcons[actor.status] || '⚪';

    card.innerHTML = `
        <div class="p-6">
            <div class="flex justify-between items-start mb-4">
                <h3 class="text-xl font-bold text-gray-900">${actor.name || '未命名Actor'}</h3>
                <span class="px-3 py-1 ${statusClass} rounded-full text-xs font-semibold flex items-center gap-1">
                    ${statusIcon} ${actor.status}
                </span>
            </div>
            <div class="space-y-3">
                <p class="text-xs text-gray-500 font-mono bg-gray-50 px-3 py-2 rounded break-all">
                    ID: ${actor.id}
                </p>
                <p class="text-sm text-gray-600">
                    <span class="font-medium">最后活跃:</span> ${lastActive}
                </p>
                <div>
                    <p class="text-sm font-medium text-gray-700 mb-2">可用函数:</p>
                    <div class="flex flex-wrap">${functionsHtml}</div>
                </div>
            </div>
            <div class="mt-6 pt-4 border-t border-gray-200 flex gap-3">
                <button onclick="viewActorDetail('${actor.id}')" class="flex-1 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium text-sm">
                    📋 查看详情
                </button>
                <button onclick="selectActorForCall('${actor.id}')" class="flex-1 px-4 py-2 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-lg hover:shadow-md transition font-medium text-sm">
                    ▶ 调用
                </button>
            </div>
        </div>
    `;
    
    return card;
}

// 选择Actor用于调用
function selectActorForCall(actorId) {
    document.getElementById('actorSelect').value = actorId;
    loadActorFunctions();
    // 滚动到调用区域
    const callSection = document.querySelector('.bg-white.rounded-lg.shadow-lg');
    if (callSection) {
        callSection.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }
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
        const result = await response.json();
        // API 返回格式: { success: true, data: [...] }
        const functions = result.data || [];
        
        if (functions.length > 0) {
            functions.forEach(func => {
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
            // 直接发送参数对象，不包装在 params 中
            body: JSON.stringify(parameters)
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
    
    const success = result.success || false;
    const bgColor = success ? 'bg-green-50 border-green-300' : 'bg-red-50 border-red-300';
    const iconColor = success ? 'text-green-600' : 'text-red-600';
    const icon = success ? '✅' : '❌';
    
    container.innerHTML = `
        <div class="border-2 ${bgColor} rounded-lg p-6">
            <h4 class="text-lg font-bold ${iconColor} mb-4 flex items-center gap-2">
                <span class="text-2xl">${icon}</span>
                函数调用结果
            </h4>
            <div class="bg-white rounded-lg p-4 border border-gray-200">
                <pre class="text-sm font-mono overflow-x-auto">${JSON.stringify(result, null, 2)}</pre>
            </div>
        </div>
    `;
}

// 查看Actor详情
async function viewActorDetail(id) {
    try {
        const response = await fetch(`/api/v1/actors/${id}`);
        const result = await response.json();
        // API 返回格式: { success: true, data: {...} }
        const actor = result.data;
        
        const content = document.getElementById('actorDetailContent');
        
        const statusColors = {
            'running': 'bg-green-100 text-green-800',
            'stopped': 'bg-gray-100 text-gray-800',
            'error': 'bg-red-100 text-red-800',
            'idle': 'bg-yellow-100 text-yellow-800'
        };
        const statusClass = statusColors[actor.status] || 'bg-gray-100 text-gray-800';
        
        content.innerHTML = `
            <div class="space-y-6">
                <div>
                    <h4 class="text-2xl font-bold text-gray-900 mb-4">${actor.name || '未命名Actor'}</h4>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div class="bg-gray-50 p-4 rounded-lg">
                            <p class="text-sm text-gray-600 mb-1">Actor ID</p>
                            <p class="font-mono text-xs text-gray-800 break-all">${actor.id}</p>
                        </div>
                        <div class="bg-gray-50 p-4 rounded-lg">
                            <p class="text-sm text-gray-600 mb-1">状态</p>
                            <span class="inline-block px-3 py-1 ${statusClass} rounded-full text-sm font-semibold">${actor.status}</span>
                        </div>
                        <div class="bg-gray-50 p-4 rounded-lg md:col-span-2">
                            <p class="text-sm text-gray-600 mb-1">最后活跃时间</p>
                            <p class="text-gray-800">${actor.lastActive ? new Date(actor.lastActive).toLocaleString() : '未知'}</p>
                        </div>
                    </div>
                </div>
                
                <div>
                    <h5 class="text-lg font-bold text-gray-900 mb-3">可用函数</h5>
                    <div class="flex flex-wrap gap-2">
                        ${actor.functions && actor.functions.length > 0 
                            ? actor.functions.map(func => `<span class="px-4 py-2 bg-blue-100 text-blue-700 rounded-lg text-sm font-medium">${func}</span>`).join('')
                            : '<p class="text-gray-400 italic">暂无函数</p>'}
                    </div>
                </div>
                
                <div>
                    <h5 class="text-lg font-bold text-gray-900 mb-3">完整状态信息</h5>
                    <div class="bg-gray-900 text-green-400 rounded-lg p-4 overflow-x-auto">
                        <pre class="text-sm font-mono">${JSON.stringify(actor, null, 2)}</pre>
                    </div>
                </div>
            </div>
        `;
        
        const modal = document.getElementById('actorDetailModal');
        modal.classList.remove('hidden');
        modal.classList.add('flex');
    } catch (error) {
        console.error('加载Actor详情失败:', error);
        alert('加载详情失败: ' + error.message);
    }
}

// 关闭Actor详情模态框
function closeActorDetailModal() {
    const modal = document.getElementById('actorDetailModal');
    modal.classList.add('hidden');
    modal.classList.remove('flex');
}

// 显示健康检查
async function showActorHealthCheck() {
    try {
        const response = await fetch('/api/v1/actors/health');
        const result = await response.json();
        // API 返回格式: { success: true, data: {...} }
        const healthData = result.data;
        
        const content = document.getElementById('healthCheckContent');
        
        const isHealthy = healthData.status === 'healthy';
        const statusColor = isHealthy ? 'text-green-600' : 'text-red-600';
        const statusBg = isHealthy ? 'bg-green-100' : 'bg-red-100';
        const statusIcon = isHealthy ? '💚' : '❤️';
        
        content.innerHTML = `
            <div class="space-y-6">
                <div class="text-center">
                    <div class="text-6xl mb-4">${statusIcon}</div>
                    <h4 class="text-2xl font-bold ${statusColor} mb-2">系统${isHealthy ? '健康' : '异常'}</h4>
                </div>
                
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div class="bg-blue-50 p-4 rounded-lg text-center border border-blue-200">
                        <p class="text-3xl font-bold text-blue-600">${healthData.actors || 0}</p>
                        <p class="text-sm text-gray-600 mt-1">活跃Actor数量</p>
                    </div>
                    <div class="${statusBg} p-4 rounded-lg text-center border ${isHealthy ? 'border-green-200' : 'border-red-200'}">
                        <p class="text-3xl font-bold ${statusColor}">${healthData.status}</p>
                        <p class="text-sm text-gray-600 mt-1">系统状态</p>
                    </div>
                    <div class="bg-purple-50 p-4 rounded-lg text-center border border-purple-200">
                        <p class="text-sm font-medium text-purple-600 break-all">${healthData.timestamp || '未知'}</p>
                        <p class="text-sm text-gray-600 mt-1">检查时间</p>
                    </div>
                </div>
                
                <div>
                    <h5 class="text-lg font-bold text-gray-900 mb-3">详细状态信息</h5>
                    <div class="bg-gray-900 text-green-400 rounded-lg p-4 overflow-x-auto">
                        <pre class="text-sm font-mono">${JSON.stringify(healthData, null, 2)}</pre>
                    </div>
                </div>
            </div>
        `;
        
        const modal = document.getElementById('healthCheckModal');
        modal.classList.remove('hidden');
        modal.classList.add('flex');
    } catch (error) {
        console.error('健康检查失败:', error);
        alert('健康检查失败: ' + error.message);
    }
}

// 关闭健康检查模态框
function closeHealthCheckModal() {
    const modal = document.getElementById('healthCheckModal');
    modal.classList.add('hidden');
    modal.classList.remove('flex');
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
