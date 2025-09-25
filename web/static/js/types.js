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
        const data = await response.json();
        allTypes = data.data || [];
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
    card.className = 'type-card';
    
    const propertiesHtml = type.attributes && Object.keys(type.attributes).length > 0 
        ? Object.entries(type.attributes).map(([name, schema]) => 
            `<div class="property-item">
                <span class="property-name">${name}</span>
                <span class="property-type">${schema.type || 'string'}</span>
            </div>`
          ).join('')
        : '<p style="color: #7f8c8d; font-style: italic;">暂无属性定义</p>';

    const featuresHtml = type.features && Object.keys(type.features).length > 0 
        ? Object.entries(type.features).map(([name, feature]) => {
            const properties = feature.properties || {};
            const propertiesList = Object.keys(properties).length > 0 
                ? Object.entries(properties).map(([propName, propDef]) => 
                    `<span class="feature-property">${propName}: ${propDef.type || 'any'}</span>`
                  ).join(', ')
                : '无属性';
            return `<div class="property-item">
                <span class="property-name">${name}</span>
                <span class="property-type">${propertiesList}</span>
            </div>`;
          }).join('')
        : '<p style="color: #7f8c8d; font-style: italic;">暂无功能定义</p>';

    card.innerHTML = `
        <div class="card-header">
            <h3 class="card-title">${type.name}</h3>
            <span class="card-category">${type.category}</span>
        </div>
        <p class="card-description">${type.description || '暂无描述'}</p>
        <div class="card-properties">
            <strong>属性模式:</strong>
            ${propertiesHtml}
        </div>
        <div class="card-properties">
            <strong>功能定义:</strong>
            ${featuresHtml}
        </div>
        <div class="card-actions">
            <button class="btn btn-secondary" onclick="editType('${type.id}')">编辑</button>
            <button class="btn btn-danger" onclick="deleteType('${type.id}')">删除</button>
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
    document.getElementById('createTypeModal').style.display = 'block';
}

// 关闭创建类型模态框
function closeCreateTypeModal() {
    document.getElementById('createTypeModal').style.display = 'none';
    document.getElementById('createTypeForm').reset();
    // 重置属性列表
    const attributesSchema = document.getElementById('attributesSchema');
    attributesSchema.innerHTML = `
        <div class="property-item">
            <input type="text" placeholder="属性名" class="property-name">
            <select class="property-type">
                <option value="string">字符串</option>
                <option value="number">数字</option>
                <option value="boolean">布尔值</option>
            </select>
            <button type="button" onclick="removeProperty(this)">删除</button>
        </div>
    `;
    // 重置功能列表
    const featuresSchema = document.getElementById('featuresSchema');
    featuresSchema.innerHTML = `
        <div class="feature-item">
            <div class="feature-header">
                <input type="text" placeholder="功能名" class="feature-name">
                <input type="text" placeholder="功能描述" class="feature-description">
                <button type="button" onclick="removeFeature(this)">删除</button>
            </div>
            <div class="feature-properties">
                <label class="feature-properties-label">功能属性:</label>
                <div class="feature-property-list">
                    <div class="feature-property-item">
                        <input type="text" placeholder="属性名" class="feature-property-name">
                        <select class="feature-property-type">
                            <option value="string">字符串</option>
                            <option value="number">数字</option>
                            <option value="boolean">布尔值</option>
                        </select>
                        <button type="button" onclick="removeFeatureProperty(this)">删除</button>
                    </div>
                </div>
                <button type="button" onclick="addFeatureProperty(this)" class="btn btn-small">添加属性</button>
            </div>
        </div>
    `;
}

// 添加属性
function addProperty() {
    const container = document.getElementById('attributesSchema');
    const propertyItem = document.createElement('div');
    propertyItem.className = 'property-item';
    propertyItem.innerHTML = `
        <input type="text" placeholder="属性名" class="property-name">
        <select class="property-type">
            <option value="string">字符串</option>
            <option value="number">数字</option>
            <option value="boolean">布尔值</option>
        </select>
        <button type="button" onclick="removeProperty(this)">删除</button>
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
    featureItem.className = 'feature-item';
    featureItem.innerHTML = `
        <div class="feature-header">
            <input type="text" placeholder="功能名" class="feature-name">
            <input type="text" placeholder="功能描述" class="feature-description">
            <button type="button" onclick="removeFeature(this)">删除</button>
        </div>
        <div class="feature-properties">
            <label class="feature-properties-label">功能属性:</label>
            <div class="feature-property-list">
                <div class="feature-property-item">
                    <input type="text" placeholder="属性名" class="feature-property-name">
                    <select class="feature-property-type">
                        <option value="string">字符串</option>
                        <option value="number">数字</option>
                        <option value="boolean">布尔值</option>
                    </select>
                    <button type="button" onclick="removeFeatureProperty(this)">删除</button>
                </div>
            </div>
            <button type="button" onclick="addFeatureProperty(this)" class="btn btn-small">添加属性</button>
        </div>
    `;
    container.appendChild(featureItem);
}

// 删除功能
function removeFeature(button) {
    const container = document.getElementById('featuresSchema');
    if (container.children.length > 1) {
        button.closest('.feature-item').remove();
    }
}

// 添加功能属性
function addFeatureProperty(button) {
    const propertyList = button.previousElementSibling;
    const propertyItem = document.createElement('div');
    propertyItem.className = 'feature-property-item';
    propertyItem.innerHTML = `
        <input type="text" placeholder="属性名" class="feature-property-name">
        <select class="feature-property-type">
            <option value="string">字符串</option>
            <option value="number">数字</option>
            <option value="boolean">布尔值</option>
        </select>
        <button type="button" onclick="removeFeatureProperty(this)">删除</button>
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
    const propertyItems = document.querySelectorAll('#attributesSchema .property-item');
    propertyItems.forEach(item => {
        const name = item.querySelector('.property-name').value.trim();
        const type = item.querySelector('.property-type').value;
        if (name) {
            schema[name] = { type: type };
        }
    });
    
    // 收集功能定义
    const features = {};
    const featureItems = document.querySelectorAll('#featuresSchema .feature-item');
    featureItems.forEach(item => {
        const name = item.querySelector('.feature-name').value.trim();
        const description = item.querySelector('.feature-description').value.trim();
        if (name) {
            // 收集功能属性
            const properties = {};
            const propertyItems = item.querySelectorAll('.feature-property-item');
            propertyItems.forEach(propItem => {
                const propName = propItem.querySelector('.feature-property-name').value.trim();
                const propType = propItem.querySelector('.feature-property-type').value;
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
