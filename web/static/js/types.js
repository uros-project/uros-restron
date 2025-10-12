// 类型管理 JavaScript

let allTypes = [];

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', function() {
    loadTypes();
});

// 加载所有类型
async function loadTypes() {
    try {
        const response = await fetch('/api/v1/thing-types');
        const result = await response.json();
        // API 返回格式: { success: true, data: { count: N, data: [...] } }
        allTypes = result.data?.data || [];
        renderTypes(allTypes);
    } catch (error) {
        console.error('加载类型失败:', error);
        alert('加载类型失败: ' + error.message);
    }
}

// 渲染类型列表
function renderTypes(types) {
    const grid = document.getElementById('typesGrid');
    grid.innerHTML = '';

    types.forEach(type => {
        const card = createTypeCard(type);
        grid.appendChild(card);
    });
}

// 创建类型卡片
function createTypeCard(type) {
    const card = document.createElement('div');
    card.className = 'bg-white rounded-xl shadow-md hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-200';
    
    const categoryColors = {
        'person': 'bg-blue-100 text-blue-800',
        'machine': 'bg-purple-100 text-purple-800',
        'object': 'bg-green-100 text-green-800'
    };
    const categoryClass = categoryColors[type.category] || 'bg-gray-100 text-gray-800';
    
    const propertiesHtml = type.attributes && Object.keys(type.attributes).length > 0 
        ? Object.entries(type.attributes).map(([name, schema]) => 
            `<div class="flex items-center justify-between py-2 border-b border-gray-100 last:border-0">
                <span class="text-sm font-medium text-gray-700">${name}</span>
                <span class="text-xs px-2 py-1 bg-blue-50 text-blue-700 rounded">${schema.type || 'string'}</span>
            </div>`
          ).join('')
        : '<p class="text-gray-400 italic text-sm">暂无属性定义</p>';

    const featuresHtml = type.features && Object.keys(type.features).length > 0 
        ? Object.entries(type.features).map(([name, feature]) => {
            const properties = feature.properties || {};
            const propertiesList = Object.keys(properties).length > 0 
                ? Object.entries(properties).map(([propName, propDef]) => 
                    `<span class="text-xs text-gray-600">${propName}: ${propDef.type || 'any'}</span>`
                  ).join(', ')
                : '无属性';
            return `<div class="flex items-start justify-between py-2 border-b border-gray-100 last:border-0">
                <span class="text-sm font-medium text-gray-700">${name}</span>
                <span class="text-xs text-gray-500">${propertiesList}</span>
            </div>`;
          }).join('')
        : '<p class="text-gray-400 italic text-sm">暂无功能定义</p>';

    card.innerHTML = `
        <div class="p-6">
            <div class="flex justify-between items-start mb-4">
                <h3 class="text-xl font-bold text-gray-900">${type.name}</h3>
                <span class="px-3 py-1 ${categoryClass} rounded-full text-xs font-semibold">${type.category}</span>
            </div>
            <p class="text-gray-600 text-sm mb-4">${type.description || '暂无描述'}</p>
            <div class="space-y-4">
                <div class="bg-gray-50 rounded-lg p-3">
                    <p class="text-sm font-semibold text-gray-700 mb-2">属性模式:</p>
                    ${propertiesHtml}
                </div>
                <div class="bg-gray-50 rounded-lg p-3">
                    <p class="text-sm font-semibold text-gray-700 mb-2">功能定义:</p>
                    ${featuresHtml}
                </div>
            </div>
            <div class="flex gap-2 mt-6 pt-4 border-t border-gray-200">
                <button onclick="editType('${type.id}')" class="flex-1 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium text-sm">✏️ 编辑</button>
                <button onclick="deleteType('${type.id}')" class="flex-1 px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition font-medium text-sm">🗑️ 删除</button>
            </div>
        </div>
    `;
    
    return card;
}

// 过滤类型
function filterTypes() {
    const category = document.getElementById('categoryFilter').value;
    const filtered = category ? allTypes.filter(type => type.category === category) : allTypes;
    renderTypes(filtered);
}

// 显示创建类型模态框
function showCreateTypeModal() {
    const modal = document.getElementById('createTypeModal');
    modal.classList.remove('hidden');
    modal.classList.add('flex');
}

// 关闭创建类型模态框
function closeCreateTypeModal() {
    const modal = document.getElementById('createTypeModal');
    modal.classList.add('hidden');
    modal.classList.remove('flex');
    document.getElementById('createTypeForm').reset();
    // 重置属性列表
    const attributesSchema = document.getElementById('attributesSchema');
    attributesSchema.innerHTML = `
        <div class="flex gap-2">
            <input type="text" placeholder="属性名" class="flex-1 px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-blue-500 transition">
            <select class="px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-blue-500 transition">
                <option value="string">字符串</option>
                <option value="number">数字</option>
                <option value="boolean">布尔值</option>
            </select>
            <button type="button" onclick="removeProperty(this)" class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition">删除</button>
        </div>
    `;
    // 重置功能列表
    const featuresSchema = document.getElementById('featuresSchema');
    featuresSchema.innerHTML = `
        <div class="bg-gray-50 p-4 rounded-lg border border-gray-200">
            <div class="flex gap-2 mb-3">
                <input type="text" placeholder="功能名" class="flex-1 px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-blue-500 transition">
                <input type="text" placeholder="功能描述" class="flex-1 px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-blue-500 transition">
                <button type="button" onclick="removeFeature(this)" class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition">删除</button>
            </div>
            <div class="ml-4">
                <label class="block text-sm font-medium text-gray-600 mb-2">功能属性:</label>
                <div class="space-y-2 mb-2">
                    <div class="flex gap-2">
                        <input type="text" placeholder="属性名" class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:border-blue-500 transition text-sm">
                        <select class="px-3 py-2 border border-gray-300 rounded-lg focus:border-blue-500 transition text-sm">
                            <option value="string">字符串</option>
                            <option value="number">数字</option>
                            <option value="boolean">布尔值</option>
                        </select>
                        <button type="button" onclick="removeFeatureProperty(this)" class="px-3 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition text-sm">删除</button>
                    </div>
                </div>
                <button type="button" onclick="addFeatureProperty(this)" class="px-3 py-2 bg-gray-100 text-gray-600 rounded-lg hover:bg-gray-200 transition text-sm">+ 添加属性</button>
            </div>
        </div>
    `;
}

// 添加属性
function addProperty() {
    const container = document.getElementById('attributesSchema');
    const propertyItem = document.createElement('div');
    propertyItem.className = 'flex gap-2';
    propertyItem.innerHTML = `
        <input type="text" placeholder="属性名" class="flex-1 px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-blue-500 transition">
        <select class="px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-blue-500 transition">
            <option value="string">字符串</option>
            <option value="number">数字</option>
            <option value="boolean">布尔值</option>
        </select>
        <button type="button" onclick="removeProperty(this)" class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition">删除</button>
    `;
    container.appendChild(propertyItem);
}

// 删除属性
function removeProperty(button) {
    const container = document.getElementById('attributesSchema');
    if (container.children.length > 1) {
        button.parentElement.remove();
    }
}

// 添加功能
function addFeature() {
    const container = document.getElementById('featuresSchema');
    const featureItem = document.createElement('div');
    featureItem.className = 'bg-gray-50 p-4 rounded-lg border border-gray-200';
    featureItem.innerHTML = `
        <div class="flex gap-2 mb-3">
            <input type="text" placeholder="功能名" class="flex-1 px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-blue-500 transition">
            <input type="text" placeholder="功能描述" class="flex-1 px-4 py-2 border-2 border-gray-300 rounded-lg focus:border-blue-500 transition">
            <button type="button" onclick="removeFeature(this)" class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition">删除</button>
        </div>
        <div class="ml-4">
            <label class="block text-sm font-medium text-gray-600 mb-2">功能属性:</label>
            <div class="space-y-2 mb-2">
                <div class="flex gap-2">
                    <input type="text" placeholder="属性名" class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:border-blue-500 transition text-sm">
                    <select class="px-3 py-2 border border-gray-300 rounded-lg focus:border-blue-500 transition text-sm">
                        <option value="string">字符串</option>
                        <option value="number">数字</option>
                        <option value="boolean">布尔值</option>
                    </select>
                    <button type="button" onclick="removeFeatureProperty(this)" class="px-3 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition text-sm">删除</button>
                </div>
            </div>
            <button type="button" onclick="addFeatureProperty(this)" class="px-3 py-2 bg-gray-100 text-gray-600 rounded-lg hover:bg-gray-200 transition text-sm">+ 添加属性</button>
        </div>
    `;
    container.appendChild(featureItem);
}

// 删除功能
function removeFeature(button) {
    const container = document.getElementById('featuresSchema');
    if (container.children.length > 1) {
        button.closest('div').remove();
    }
}

// 添加功能属性
function addFeatureProperty(button) {
    const propertyList = button.previousElementSibling;
    const propertyItem = document.createElement('div');
    propertyItem.className = 'flex gap-2';
    propertyItem.innerHTML = `
        <input type="text" placeholder="属性名" class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:border-blue-500 transition text-sm">
        <select class="px-3 py-2 border border-gray-300 rounded-lg focus:border-blue-500 transition text-sm">
            <option value="string">字符串</option>
            <option value="number">数字</option>
            <option value="boolean">布尔值</option>
        </select>
        <button type="button" onclick="removeFeatureProperty(this)" class="px-3 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition text-sm">删除</button>
    `;
    propertyList.appendChild(propertyItem);
}

// 删除功能属性
function removeFeatureProperty(button) {
    const propertyList = button.closest('.feature-property-list');
    if (propertyList.children.length > 1) {
        button.parentElement.remove();
    }
}

// 创建类型表单提交
document.getElementById('createTypeForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const schema = {};
    
    // 收集属性模式
    const propertyItems = document.querySelectorAll('#attributesSchema > div');
    propertyItems.forEach(item => {
        const name = item.querySelector('input[placeholder="属性名"]').value.trim();
        const type = item.querySelector('select').value;
        if (name) {
            schema[name] = { type: type };
        }
    });
    
    // 收集功能定义
    const features = {};
    const featureItems = document.querySelectorAll('#featuresSchema > div');
    featureItems.forEach(item => {
        const name = item.querySelector('input[placeholder="功能名"]').value.trim();
        const description = item.querySelector('input[placeholder="功能描述"]').value.trim();
        if (name) {
            // 收集功能属性
            const properties = {};
            const propertyItems = item.querySelectorAll('.space-y-2 > div');
            propertyItems.forEach(propItem => {
                const propName = propItem.querySelector('input[placeholder="属性名"]').value.trim();
                const propType = propItem.querySelector('select').value;
                if (propName) {
                    properties[propName] = { type: propType };
                }
            });
            
            // 功能定义直接包含属性，符合 Ditto 标准
            features[name] = properties;
        }
    });
    
    const typeData = {
        name: formData.get('name'),
        description: formData.get('description'),
        category: formData.get('category'),
        attributes: schema,
        features: features
    };
    
    try {
        const response = await fetch('/api/v1/thing-types', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(typeData)
        });
        
        if (response.ok) {
            alert('类型创建成功！');
            closeCreateTypeModal();
            loadTypes();
        } else {
            const error = await response.json();
            alert('创建失败: ' + (error.error || '未知错误'));
        }
    } catch (error) {
        console.error('创建类型失败:', error);
        alert('创建失败: ' + error.message);
    }
});

// 删除类型
async function deleteType(id) {
    if (!confirm('确定要删除这个类型吗？')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/v1/thing-types/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            alert('类型删除成功！');
            loadTypes();
        } else {
            const error = await response.json();
            alert('删除失败: ' + (error.error || '未知错误'));
        }
    } catch (error) {
        console.error('删除类型失败:', error);
        alert('删除失败: ' + error.message);
    }
}

// 编辑类型（简化版，实际项目中可以添加编辑功能）
function editType(id) {
    alert('编辑功能待实现');
}

// 点击模态框外部关闭
window.onclick = function(event) {
    const modal = document.getElementById('createTypeModal');
    if (event.target === modal) {
        closeCreateTypeModal();
    }
}
