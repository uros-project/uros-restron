// 关系管理 JavaScript

let allRelationships = [];
let allThings = [];

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', function() {
    loadThings();
    loadRelationships();
});

// 加载所有事物
async function loadThings() {
    try {
        const response = await fetch('/api/v1/things');
        const result = await response.json();
        // API 返回格式: { success: true, data: [...], count: N }
        allThings = result.data || [];
        
        // 填充源事物和目标事物选择框
        const sourceSelect = document.getElementById('sourceId');
        const targetSelect = document.getElementById('targetId');
        
        sourceSelect.innerHTML = '<option value="">请选择源事物</option>';
        targetSelect.innerHTML = '<option value="">请选择目标事物</option>';
        
        allThings.forEach(thing => {
            const option1 = document.createElement('option');
            option1.value = thing.id;
            option1.textContent = `${thing.name} (${thing.type})`;
            sourceSelect.appendChild(option1);
            
            const option2 = document.createElement('option');
            option2.value = thing.id;
            option2.textContent = `${thing.name} (${thing.type})`;
            targetSelect.appendChild(option2);
        });
    } catch (error) {
        console.error('加载事物失败:', error);
    }
}

// 加载所有关系
async function loadRelationships() {
    try {
        const response = await fetch('/api/v1/relationships');
        const result = await response.json();
        // API 返回格式: { success: true, data: {data: [...], count: N} }
        allRelationships = result.data?.data || [];
        renderRelationships(allRelationships);
    } catch (error) {
        console.error('加载关系失败:', error);
        alert('加载关系失败: ' + error.message);
    }
}

// 渲染关系列表
function renderRelationships(relationships) {
    const grid = document.getElementById('relationshipsGrid');
    grid.innerHTML = '';

    relationships.forEach(relationship => {
        const card = createRelationshipCard(relationship);
        grid.appendChild(card);
    });
}

// 创建关系卡片
function createRelationshipCard(relationship) {
    const card = document.createElement('div');
    card.className = 'bg-white rounded-xl shadow-md hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-200';
    
    const sourceThing = allThings.find(t => t.id === relationship.sourceId);
    const targetThing = allThings.find(t => t.id === relationship.targetId);
    
    const typeColors = {
        'contains': 'bg-blue-100 text-blue-800',
        'owns': 'bg-purple-100 text-purple-800',
        'collaborates': 'bg-green-100 text-green-800',
        'depends_on': 'bg-yellow-100 text-yellow-800',
        'composes': 'bg-pink-100 text-pink-800',
        'influences': 'bg-orange-100 text-orange-800',
        'relates_to': 'bg-gray-100 text-gray-800'
    };
    const typeClass = typeColors[relationship.type] || 'bg-gray-100 text-gray-800';
    
    const propertiesHtml = relationship.properties && Object.keys(relationship.properties).length > 0 
        ? Object.entries(relationship.properties).map(([key, value]) => 
            `<div class="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                <span class="text-sm font-medium text-gray-700">${key}</span>
                <span class="text-sm text-gray-600">${value}</span>
            </div>`
          ).join('')
        : '<p class="text-gray-400 italic text-sm">暂无属性</p>';

    card.innerHTML = `
        <div class="p-6">
            <div class="flex justify-between items-start mb-4">
                <h3 class="text-xl font-bold text-gray-900">${relationship.name}</h3>
                <span class="px-3 py-1 ${typeClass} rounded-full text-xs font-semibold">${relationship.type}</span>
            </div>
            <div class="bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg p-4 mb-4">
                <div class="flex items-center justify-between gap-2">
                    <div class="flex-1 bg-white rounded-lg p-3 shadow-sm">
                        <p class="text-xs text-gray-500 mb-1">源事物</p>
                        <p class="text-sm font-semibold text-gray-900">${sourceThing ? sourceThing.name : '未知事物'}</p>
                    </div>
                    <div class="text-2xl text-blue-600 font-bold">→</div>
                    <div class="flex-1 bg-white rounded-lg p-3 shadow-sm">
                        <p class="text-xs text-gray-500 mb-1">目标事物</p>
                        <p class="text-sm font-semibold text-gray-900">${targetThing ? targetThing.name : '未知事物'}</p>
                    </div>
                </div>
            </div>
            <p class="text-gray-600 text-sm mb-4">${relationship.description || '暂无描述'}</p>
            <div class="bg-gray-50 rounded-lg p-3">
                <p class="text-sm font-semibold text-gray-700 mb-2">关系属性:</p>
                ${propertiesHtml}
            </div>
            <div class="flex gap-2 mt-6 pt-4 border-t border-gray-200">
                <button onclick="viewRelationshipDetail('${relationship.id}')" class="flex-1 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium text-sm">📋 查看详情</button>
                <button onclick="deleteRelationship('${relationship.id}')" class="flex-1 px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition font-medium text-sm">🗑️ 删除</button>
            </div>
        </div>
    `;
    
    return card;
}

// 过滤关系
function filterRelationships() {
    const type = document.getElementById('typeFilter').value;
    const search = document.getElementById('searchInput').value.toLowerCase();
    
    let filtered = allRelationships;
    
    if (type) {
        filtered = filtered.filter(relationship => relationship.type === type);
    }
    
    if (search) {
        filtered = filtered.filter(relationship => 
            relationship.name.toLowerCase().includes(search) ||
            (relationship.description && relationship.description.toLowerCase().includes(search))
        );
    }
    
    renderRelationships(filtered);
}

// 搜索关系
function searchRelationships() {
    filterRelationships();
}

// 显示创建关系模态框
function showCreateRelationshipModal() {
    const modal = document.getElementById('createRelationshipModal');
    modal.classList.remove('hidden');
    modal.classList.add('flex');
}

// 关闭创建关系模态框
function closeCreateRelationshipModal() {
    const modal = document.getElementById('createRelationshipModal');
    modal.classList.add('hidden');
    modal.classList.remove('flex');
    document.getElementById('createRelationshipForm').reset();
}

// 创建关系表单提交
document.getElementById('createRelationshipForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    
    // 解析属性JSON
    let properties = {};
    const propertiesText = formData.get('properties');
    if (propertiesText && propertiesText.trim()) {
        try {
            properties = JSON.parse(propertiesText);
        } catch (error) {
            alert('属性格式错误，请输入有效的JSON');
            return;
        }
    }
    
    const relationshipData = {
        sourceId: formData.get('sourceId'),
        targetId: formData.get('targetId'),
        type: formData.get('type'),
        name: formData.get('name'),
        description: formData.get('description'),
        properties: properties
    };
    
    try {
        const response = await fetch('/api/v1/relationships', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(relationshipData)
        });
        
        if (response.ok) {
            alert('关系创建成功！');
            closeCreateRelationshipModal();
            loadRelationships();
        } else {
            const error = await response.json();
            alert('创建失败: ' + (error.error || '未知错误'));
        }
    } catch (error) {
        console.error('创建关系失败:', error);
        alert('创建失败: ' + error.message);
    }
});

// 查看关系详情
async function viewRelationshipDetail(id) {
    try {
        const response = await fetch(`/api/v1/relationships/${id}`);
        const relationship = await response.json();
        
        const sourceThing = allThings.find(t => t.id === relationship.sourceId);
        const targetThing = allThings.find(t => t.id === relationship.targetId);
        
        const content = document.getElementById('relationshipDetailContent');
        content.innerHTML = `
            <div class="relationship-detail">
                <h4>${relationship.name}</h4>
                <p><strong>关系类型:</strong> ${relationship.type}</p>
                <p><strong>描述:</strong> ${relationship.description || '暂无描述'}</p>
                <p><strong>创建时间:</strong> ${new Date(relationship.createdAt).toLocaleString()}</p>
                <p><strong>更新时间:</strong> ${new Date(relationship.updatedAt).toLocaleString()}</p>
                
                <div class="relationship-connection-detail">
                    <h5>关系连接</h5>
                    <div class="connection-info">
                        <div class="source-info">
                            <strong>源事物:</strong> ${sourceThing ? sourceThing.name : '未知事物'}
                            <br><small>ID: ${relationship.sourceId}</small>
                        </div>
                        <div class="connection-arrow">→</div>
                        <div class="target-info">
                            <strong>目标事物:</strong> ${targetThing ? targetThing.name : '未知事物'}
                            <br><small>ID: ${relationship.targetId}</small>
                        </div>
                    </div>
                </div>
                
                <h5>关系属性</h5>
                <div class="properties-list">
                    ${relationship.properties && Object.keys(relationship.properties).length > 0 
                        ? Object.entries(relationship.properties).map(([key, value]) => 
                            `<div class="property-item">
                                <span class="property-name">${key}</span>
                                <span class="property-type">${value}</span>
                            </div>`
                          ).join('')
                        : '<p>暂无属性</p>'}
                </div>
            </div>
        `;
        
        const modal = document.getElementById('relationshipDetailModal');
        if (modal) {
            modal.classList.remove('hidden');
            modal.classList.add('flex');
        }
    } catch (error) {
        console.error('加载关系详情失败:', error);
        alert('加载详情失败: ' + error.message);
    }
}

// 关闭关系详情模态框
function closeRelationshipDetailModal() {
    const modal = document.getElementById('relationshipDetailModal');
    if (modal) {
        modal.classList.add('hidden');
        modal.classList.remove('flex');
    }
}

// 删除关系
async function deleteRelationship(id) {
    if (!confirm('确定要删除这个关系吗？')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/v1/relationships/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            alert('关系删除成功！');
            loadRelationships();
        } else {
            const error = await response.json();
            alert('删除失败: ' + (error.error || '未知错误'));
        }
    } catch (error) {
        console.error('删除关系失败:', error);
        alert('删除失败: ' + error.message);
    }
}

// 点击模态框外部关闭
window.onclick = function(event) {
    const createModal = document.getElementById('createRelationshipModal');
    const detailModal = document.getElementById('relationshipDetailModal');
    
    if (event.target === createModal) {
        closeCreateRelationshipModal();
    }
    if (event.target === detailModal) {
        closeRelationshipDetailModal();
    }
}
