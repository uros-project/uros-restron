// 事物管理 JavaScript

let allThings = [];
let allTypes = [];

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', function() {
    loadTypes();
    loadThings();
});

// 加载所有类型
async function loadTypes() {
    try {
        const response = await fetch('/api/v1/thing-types');
        const result = await response.json();
        // API 返回格式: { success: true, data: { count: N, data: [...] } }
        allTypes = result.data?.data || [];
        
        // 填充类型选择框
        const typeSelect = document.getElementById('thingType');
        typeSelect.innerHTML = '<option value="">请选择类型</option>';
        allTypes.forEach(type => {
            const option = document.createElement('option');
            option.value = type.id;
            option.textContent = `${type.name} (${type.category})`;
            typeSelect.appendChild(option);
        });
    } catch (error) {
        console.error('加载类型失败:', error);
    }
}

// 加载所有事物
async function loadThings() {
    try {
        const response = await fetch('/api/v1/things');
        const result = await response.json();
        // API 返回格式: { success: true, data: [...], count: N }
        allThings = result.data || [];
        renderThings(allThings);
    } catch (error) {
        console.error('加载事物失败:', error);
        alert('加载事物失败: ' + error.message);
    }
}

// 渲染事物列表
function renderThings(things) {
    const grid = document.getElementById('thingsGrid');
    grid.innerHTML = '';

    things.forEach(thing => {
        const card = createThingCard(thing);
        grid.appendChild(card);
    });
}

// 创建事物卡片
function createThingCard(thing) {
    const card = document.createElement('div');
    card.className = 'bg-white rounded-xl shadow-md hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-200';
    
    const attributesHtml = thing.attributes && Object.keys(thing.attributes).length > 0 
        ? Object.entries(thing.attributes).map(([key, value]) => 
            `<div class="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                <span class="text-sm font-medium text-gray-700">${key}</span>
                <span class="text-sm text-gray-600">${value}</span>
            </div>`
          ).join('')
        : '<p class="text-gray-400 italic text-sm">暂无属性</p>';

    const featuresHtml = thing.features && Object.keys(thing.features).length > 0
        ? Object.entries(thing.features).map(([name, feature]) => {
            const properties = feature.properties || {};
            const propertiesList = Object.keys(properties).length > 0 
                ? Object.entries(properties).map(([propName, propValue]) => 
                    `<span class="text-xs text-gray-600">${propName}: ${propValue}</span>`
                  ).join(', ')
                : '无属性';
            return `<div class="flex items-start justify-between py-2 border-b border-gray-100 last:border-0">
                <span class="text-sm font-medium text-gray-700">${name}</span>
                <span class="text-xs text-gray-500">${propertiesList}</span>
            </div>`;
          }).join('')
        : '<p class="text-gray-400 italic text-sm">暂无功能</p>';

    card.innerHTML = `
        <div class="p-6">
            <div class="flex justify-between items-start mb-4">
                <h3 class="text-xl font-bold text-gray-900">${thing.name}</h3>
                <span class="px-3 py-1 bg-indigo-100 text-indigo-800 rounded-full text-xs font-semibold">${thing.type}</span>
            </div>
            <p class="text-gray-600 text-sm mb-4">${thing.description || '暂无描述'}</p>
            <div class="space-y-4">
                <div class="bg-gray-50 rounded-lg p-3">
                    <p class="text-sm font-semibold text-gray-700 mb-2">属性 (Attributes):</p>
                    ${attributesHtml}
                </div>
                <div class="bg-gray-50 rounded-lg p-3">
                    <p class="text-sm font-semibold text-gray-700 mb-2">功能 (Features):</p>
                    ${featuresHtml}
                </div>
            </div>
            <div class="flex gap-2 mt-6 pt-4 border-t border-gray-200">
                <button onclick="viewThingDetail('${thing.id}')" class="flex-1 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium text-sm">📋 查看详情</button>
                <button onclick="deleteThing('${thing.id}')" class="flex-1 px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition font-medium text-sm">🗑️ 删除</button>
            </div>
        </div>
    `;
    
    return card;
}

// 过滤事物
function filterThings() {
    const type = document.getElementById('typeFilter').value;
    const search = document.getElementById('searchInput').value.toLowerCase();
    
    let filtered = allThings;
    
    if (type) {
        filtered = filtered.filter(thing => thing.type === type);
    }
    
    if (search) {
        filtered = filtered.filter(thing => 
            thing.name.toLowerCase().includes(search) ||
            (thing.description && thing.description.toLowerCase().includes(search))
        );
    }
    
    renderThings(filtered);
}

// 搜索事物
function searchThings() {
    filterThings();
}

// 显示创建事物模态框
function showCreateThingModal() {
    const modal = document.getElementById('createThingModal');
    modal.classList.remove('hidden');
    modal.classList.add('flex');
}

// 关闭创建事物模态框
function closeCreateThingModal() {
    const modal = document.getElementById('createThingModal');
    modal.classList.add('hidden');
    modal.classList.remove('flex');
    document.getElementById('createThingForm').reset();
    document.getElementById('propertiesContainer').innerHTML = '';
}

// 加载类型模式
function loadTypeSchema() {
    const typeId = document.getElementById('thingType').value;
    const container = document.getElementById('propertiesContainer');
    container.innerHTML = '';
    
    if (!typeId) return;
    
    const type = allTypes.find(t => t.id === typeId);
    if (!type) return;
    
    // 处理 Attributes（静态属性）
    if (type.attributes && Object.keys(type.attributes).length > 0) {
        const attributesSection = document.createElement('div');
        attributesSection.className = 'form-section';
        attributesSection.innerHTML = '<h4>静态属性 (Attributes)</h4>';
        
        Object.entries(type.attributes).forEach(([name, schema]) => {
            const propertyDiv = document.createElement('div');
            propertyDiv.className = 'form-group';
            
            let inputHtml = '';
            switch (schema.type) {
                case 'number':
                    inputHtml = `<input type="number" name="attributes[${name}]" placeholder="请输入${name}">`;
                    break;
                case 'boolean':
                    inputHtml = `
                        <select name="attributes[${name}]">
                            <option value="true">是</option>
                            <option value="false">否</option>
                        </select>
                    `;
                    break;
                default:
                    inputHtml = `<input type="text" name="attributes[${name}]" placeholder="请输入${name}">`;
            }
            
            propertyDiv.innerHTML = `
                <label for="attributes[${name}]">${name}</label>
                ${inputHtml}
            `;
            
            attributesSection.appendChild(propertyDiv);
        });
        
        container.appendChild(attributesSection);
    }
    
    // 处理 Features（功能属性）
    if (type.features && Object.keys(type.features).length > 0) {
        const featuresSection = document.createElement('div');
        featuresSection.className = 'form-section';
        featuresSection.innerHTML = '<h4>功能属性 (Features)</h4>';
        
        Object.entries(type.features).forEach(([featureName, feature]) => {
            const featureDiv = document.createElement('div');
            featureDiv.className = 'feature-input-group';
            featureDiv.innerHTML = `<h5>${featureName}</h5>`;
            
            if (feature.properties && Object.keys(feature.properties).length > 0) {
                Object.entries(feature.properties).forEach(([propName, propSchema]) => {
                    const propertyDiv = document.createElement('div');
                    propertyDiv.className = 'form-group';
                    
                    let inputHtml = '';
                    switch (propSchema.type) {
                        case 'number':
                            inputHtml = `<input type="number" name="features[${featureName}][${propName}]" placeholder="请输入${propName}">`;
                            break;
                        case 'boolean':
                            inputHtml = `
                                <select name="features[${featureName}][${propName}]">
                                    <option value="true">是</option>
                                    <option value="false">否</option>
                                </select>
                            `;
                            break;
                        default:
                            inputHtml = `<input type="text" name="features[${featureName}][${propName}]" placeholder="请输入${propName}">`;
                    }
                    
                    propertyDiv.innerHTML = `
                        <label for="features[${featureName}][${propName}]">${propName}</label>
                        ${inputHtml}
                    `;
                    
                    featureDiv.appendChild(propertyDiv);
                });
            }
            
            featuresSection.appendChild(featureDiv);
        });
        
        container.appendChild(featuresSection);
    }
}

// 创建事物表单提交
document.getElementById('createThingForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const typeId = formData.get('thingTypeId');
    
    // 收集 Attributes（静态属性）
    const attributes = {};
    for (const [key, value] of formData.entries()) {
        if (key.startsWith('attributes[')) {
            const propName = key.match(/attributes\[(.*?)\]/)[1];
            attributes[propName] = value;
        }
    }
    
    // 收集 Features（功能属性）
    const features = {};
    for (const [key, value] of formData.entries()) {
        if (key.startsWith('features[')) {
            const match = key.match(/features\[(.*?)\]\[(.*?)\]/);
            if (match) {
                const featureName = match[1];
                const propName = match[2];
                
                if (!features[featureName]) {
                    features[featureName] = { properties: {} };
                }
                features[featureName].properties[propName] = value;
            }
        }
    }
    
    const thingData = {
        name: formData.get('name'),
        description: formData.get('description'),
        attributes: attributes,
        features: features
    };
    
    try {
        const response = await fetch(`/api/v1/thing-types/${typeId}/things`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(thingData)
        });
        
        if (response.ok) {
            alert('事物创建成功！');
            closeCreateThingModal();
            loadThings();
        } else {
            const error = await response.json();
            alert('创建失败: ' + (error.error || '未知错误'));
        }
    } catch (error) {
        console.error('创建事物失败:', error);
        alert('创建失败: ' + error.message);
    }
});

// 查看事物详情
async function viewThingDetail(id) {
    try {
        const response = await fetch(`/api/v1/things/${id}`);
        const thing = await response.json();
        
        const content = document.getElementById('thingDetailContent');
        content.innerHTML = `
            <div class="thing-detail">
                <h4>${thing.name}</h4>
                <p><strong>类型:</strong> ${thing.type}</p>
                <p><strong>描述:</strong> ${thing.description || '暂无描述'}</p>
                <p><strong>创建时间:</strong> ${new Date(thing.createdAt).toLocaleString()}</p>
                <p><strong>更新时间:</strong> ${new Date(thing.updatedAt).toLocaleString()}</p>
                
                <h5>静态属性 (Attributes)</h5>
                <div class="properties-list">
                    ${thing.attributes && Object.keys(thing.attributes).length > 0 
                        ? Object.entries(thing.attributes).map(([key, value]) => 
                            `<div class="property-item">
                                <span class="property-name">${key}</span>
                                <span class="property-type">${value}</span>
                            </div>`
                          ).join('')
                        : '<p>暂无属性</p>'}
                </div>
                
                <h5>功能 (Features)</h5>
                <div class="features-list">
                    ${thing.features && Object.keys(thing.features).length > 0
                        ? Object.entries(thing.features).map(([name, feature]) => {
                            const properties = feature.properties || {};
                            const propertiesHtml = Object.keys(properties).length > 0 
                                ? Object.entries(properties).map(([propName, propValue]) => 
                                    `<div class="property-item">
                                        <span class="property-name">${propName}</span>
                                        <span class="property-type">${propValue}</span>
                                    </div>`
                                  ).join('')
                                : '<p style="color: #7f8c8d; font-style: italic;">无属性</p>';
                            return `
                                <div class="feature-group">
                                    <h6>${name}</h6>
                                    <div class="feature-properties">${propertiesHtml}</div>
                                </div>
                            `;
                          }).join('')
                        : '<p>暂无功能</p>'}
                </div>
            </div>
        `;
        
        const modal = document.getElementById('thingDetailModal');
        modal.classList.remove('hidden');
        modal.classList.add('flex');
    } catch (error) {
        console.error('加载事物详情失败:', error);
        alert('加载详情失败: ' + error.message);
    }
}

// 关闭事物详情模态框
function closeThingDetailModal() {
    const modal = document.getElementById('thingDetailModal');
    modal.classList.add('hidden');
    modal.classList.remove('flex');
}

// 删除事物
async function deleteThing(id) {
    if (!confirm('确定要删除这个事物吗？')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/v1/things/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            alert('事物删除成功！');
            loadThings();
        } else {
            const error = await response.json();
            alert('删除失败: ' + (error.error || '未知错误'));
        }
    } catch (error) {
        console.error('删除事物失败:', error);
        alert('删除失败: ' + error.message);
    }
}

// 点击模态框外部关闭
window.onclick = function(event) {
    const createModal = document.getElementById('createThingModal');
    const detailModal = document.getElementById('thingDetailModal');
    
    if (event.target === createModal) {
        closeCreateThingModal();
    }
    if (event.target === detailModal) {
        closeThingDetailModal();
    }
}
