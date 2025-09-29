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
        const data = await response.json();
        allThings = data.data || [];
        
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
        const data = await response.json();
        allRelationships = data.data || [];
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
    card.className = 'relationship-card';
    
    const sourceThing = allThings.find(t => t.id === relationship.sourceId);
    const targetThing = allThings.find(t => t.id === relationship.targetId);
    
    const propertiesHtml = relationship.properties && Object.keys(relationship.properties).length > 0 
        ? Object.entries(relationship.properties).map(([key, value]) => 
            `<div class="property-item">
                <span class="property-name">${key}</span>
                <span class="property-type">${value}</span>
            </div>`
          ).join('')
        : '<p style="color: #7f8c8d; font-style: italic;">暂无属性</p>';

    card.innerHTML = `
        <div class="card-header">
            <h3 class="card-title">${relationship.name}</h3>
            <span class="card-category">${relationship.type}</span>
        </div>
        <div class="relationship-content">
            <div class="relationship-connection">
                <div class="relationship-source">
                    <strong>源:</strong> ${sourceThing ? sourceThing.name : '未知事物'}
                </div>
                <div class="relationship-arrow">→</div>
                <div class="relationship-target">
                    <strong>目标:</strong> ${targetThing ? targetThing.name : '未知事物'}
                </div>
            </div>
            <p class="card-description">${relationship.description || '暂无描述'}</p>
            <div class="card-properties">
                <strong>关系属性:</strong>
                ${propertiesHtml}
            </div>
        </div>
        <div class="card-actions">
            <button class="btn btn-secondary" onclick="viewRelationshipDetail('${relationship.id}')">查看详情</button>
            <button class="btn btn-danger" onclick="deleteRelationship('${relationship.id}')">删除</button>
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
    document.getElementById('createRelationshipModal').style.display = 'block';
}

// 关闭创建关系模态框
function closeCreateRelationshipModal() {
    document.getElementById('createRelationshipModal').style.display = 'none';
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
        
        document.getElementById('relationshipDetailModal').style.display = 'block';
    } catch (error) {
        console.error('加载关系详情失败:', error);
        alert('加载详情失败: ' + error.message);
    }
}

// 关闭关系详情模态框
function closeRelationshipDetailModal() {
    document.getElementById('relationshipDetailModal').style.display = 'none';
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
