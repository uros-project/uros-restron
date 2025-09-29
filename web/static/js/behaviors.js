// 行为管理页面的JavaScript功能

class BehaviorManager {
    constructor() {
        this.currentTab = 'behaviors';
        this.behaviors = [];
        this.actors = [];
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.loadBehaviors();
        this.loadActors();
    }

    setupEventListeners() {
        // 标签页切换
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.switchTab(e.target.dataset.tab);
            });
        });

        // 按钮事件
        document.getElementById('loadBehaviors').addEventListener('click', () => {
            this.loadBehaviors();
        });

        document.getElementById('seedBehaviors').addEventListener('click', () => {
            this.seedBehaviors();
        });

        document.getElementById('refreshActors').addEventListener('click', () => {
            this.loadActors();
        });

        // 过滤器事件
        document.getElementById('categoryFilter').addEventListener('change', () => {
            this.filterBehaviors();
        });

        document.getElementById('typeFilter').addEventListener('change', () => {
            this.filterBehaviors();
        });

        // Actor选择事件
        document.getElementById('actorSelect').addEventListener('change', (e) => {
            this.loadActorFunctions(e.target.value);
        });

        // 函数调用
        document.getElementById('callFunction').addEventListener('click', () => {
            this.callActorFunction();
        });
    }

    switchTab(tabName) {
        // 更新标签页按钮状态
        document.querySelectorAll('.tab-btn').forEach(btn => {
            btn.classList.remove('active');
        });
        document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');

        // 更新标签页内容
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.remove('active');
        });
        document.getElementById(tabName).classList.add('active');

        this.currentTab = tabName;
    }

    async loadBehaviors() {
        try {
            console.log('Loading behaviors...');
            const response = await fetch('/api/v1/behaviors');
            const result = await response.json();
            console.log('Behaviors API response:', result);
            if (result.success && result.data) {
                this.behaviors = result.data;
                console.log('Behaviors loaded:', this.behaviors.length, 'items');
                this.renderBehaviors(result.data);
            } else {
                console.error('Failed to load behaviors:', result);
                this.showError('加载行为失败: ' + (result.error || '未知错误'));
            }
        } catch (error) {
            console.error('Failed to load behaviors:', error);
            this.showError('加载行为失败');
        }
    }

    async loadActors() {
        try {
            const response = await fetch('/api/v1/actors');
            const result = await response.json();
            if (result.success && result.data) {
                this.actors = result.data;
                this.renderActors(result.data);
                this.updateActorSelect(result.data);
            } else {
                this.showError('加载Actor状态失败: ' + (result.error || '未知错误'));
            }
        } catch (error) {
            console.error('Failed to load actors:', error);
            this.showError('加载Actor状态失败');
        }
    }

    async seedBehaviors() {
        try {
            const response = await fetch('/api/v1/behaviors/seed', {
                method: 'POST'
            });
            const result = await response.json();
            if (response.ok) {
                this.showSuccess('预定义行为填充成功');
                this.loadBehaviors();
            } else {
                this.showError('填充预定义行为失败: ' + result.message);
            }
        } catch (error) {
            console.error('Failed to seed behaviors:', error);
            this.showError('填充预定义行为失败');
        }
    }

    renderBehaviors(behaviors) {
        console.log('Rendering behaviors:', behaviors);
        const container = document.getElementById('behaviorsList');
        if (!container) {
            console.error('behaviorsList container not found');
            return;
        }
        container.innerHTML = '';

        if (behaviors.length === 0) {
            container.innerHTML = '<p class="no-data">暂无行为数据</p>';
            return;
        }

        console.log('Creating cards for', behaviors.length, 'behaviors');
        behaviors.forEach(behavior => {
            const behaviorCard = this.createBehaviorCard(behavior);
            container.appendChild(behaviorCard);
        });
    }

    createBehaviorCard(behavior) {
        const card = document.createElement('div');
        card.className = 'behavior-card';
        card.innerHTML = `
            <div class="behavior-header">
                <h3>${behavior.name}</h3>
                <span class="behavior-type">${behavior.type}</span>
            </div>
            <div class="behavior-content">
                <p class="behavior-description">${behavior.description}</p>
                <div class="behavior-meta">
                    <span class="behavior-category">分类: ${behavior.category}</span>
                </div>
                <div class="behavior-actions">
                    <button class="btn btn-sm btn-primary" onclick="behaviorManager.viewBehavior('${behavior.id}')">
                        查看详情
                    </button>
                    <button class="btn btn-sm btn-secondary" onclick="behaviorManager.editBehavior('${behavior.id}')">
                        编辑
                    </button>
                    <button class="btn btn-sm btn-danger" onclick="behaviorManager.deleteBehavior('${behavior.id}')">
                        删除
                    </button>
                </div>
            </div>
        `;
        return card;
    }

    renderActors(actors) {
        const container = document.getElementById('actorsList');
        container.innerHTML = '';

        if (actors.length === 0) {
            container.innerHTML = '<p class="no-data">暂无Actor数据</p>';
            return;
        }

        actors.forEach(actor => {
            const actorCard = this.createActorCard(actor);
            container.appendChild(actorCard);
        });
    }

    createActorCard(actor) {
        const card = document.createElement('div');
        card.className = 'actor-card';
        card.innerHTML = `
            <div class="actor-header">
                <h3>Actor ${actor.id.substring(0, 8)}...</h3>
                <span class="actor-status ${actor.status}">${actor.status}</span>
            </div>
            <div class="actor-content">
                <p class="actor-name">${actor.name || '未命名Actor'}</p>
                <div class="actor-meta">
                    <span class="actor-type">类型: ${actor.type}</span>
                    <span class="actor-last-active">最后活跃: ${new Date(actor.lastActive).toLocaleString()}</span>
                </div>
                <div class="actor-actions">
                    <button class="btn btn-sm btn-primary" onclick="behaviorManager.viewActor('${actor.id}')">
                        查看详情
                    </button>
                    <button class="btn btn-sm btn-secondary" onclick="behaviorManager.callActorFunction('${actor.id}')">
                        调用函数
                    </button>
                </div>
            </div>
        `;
        return card;
    }

    updateActorSelect(actors) {
        const select = document.getElementById('actorSelect');
        select.innerHTML = '<option value="">请选择Actor</option>';
        
        actors.forEach(actor => {
            const option = document.createElement('option');
            option.value = actor.id;
            option.textContent = `${actor.name || '未命名Actor'} (${actor.id.substring(0, 8)}...)`;
            select.appendChild(option);
        });
    }

    async loadActorFunctions(actorId) {
        if (!actorId) {
            document.getElementById('functionSelect').innerHTML = '<option value="">请选择函数</option>';
            return;
        }

        try {
            const response = await fetch(`/api/v1/actors/${actorId}/functions`);
            const data = await response.json();
            
            const select = document.getElementById('functionSelect');
            select.innerHTML = '<option value="">请选择函数</option>';
            
            if (data.functions && data.functions.length > 0) {
                data.functions.forEach(func => {
                    const option = document.createElement('option');
                    option.value = func.name;
                    option.textContent = func.name;
                    select.appendChild(option);
                });
            }
        } catch (error) {
            console.error('Failed to load actor functions:', error);
            this.showError('加载Actor函数失败');
        }
    }

    async callActorFunction() {
        const actorId = document.getElementById('actorSelect').value;
        const functionName = document.getElementById('functionSelect').value;
        const parametersText = document.getElementById('parameters').value;

        if (!actorId || !functionName) {
            this.showError('请选择Actor和函数');
            return;
        }

        let parameters = {};
        if (parametersText.trim()) {
            try {
                parameters = JSON.parse(parametersText);
            } catch (error) {
                this.showError('参数格式错误，请输入有效的JSON');
                return;
            }
        }

        try {
            const response = await fetch(`/api/v1/actors/${actorId}/functions/${functionName}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ parameters })
            });

            const result = await response.json();
            this.displayFunctionResult(result);
        } catch (error) {
            console.error('Failed to call function:', error);
            this.showError('函数调用失败');
        }
    }

    displayFunctionResult(result) {
        const container = document.getElementById('functionResult');
        container.innerHTML = `
            <h4>函数调用结果</h4>
            <pre class="result-json">${JSON.stringify(result, null, 2)}</pre>
        `;
    }

    filterBehaviors() {
        const categoryFilter = document.getElementById('categoryFilter').value;
        const typeFilter = document.getElementById('typeFilter').value;

        let filteredBehaviors = this.behaviors;

        if (categoryFilter) {
            filteredBehaviors = filteredBehaviors.filter(b => b.category === categoryFilter);
        }

        if (typeFilter) {
            filteredBehaviors = filteredBehaviors.filter(b => b.type === typeFilter);
        }

        this.renderBehaviors(filteredBehaviors);
    }

    viewBehavior(behaviorId) {
        // 实现查看行为详情的逻辑
        console.log('View behavior:', behaviorId);
    }

    editBehavior(behaviorId) {
        // 实现编辑行为的逻辑
        console.log('Edit behavior:', behaviorId);
    }

    async deleteBehavior(behaviorId) {
        if (!confirm('确定要删除这个行为吗？')) {
            return;
        }

        try {
            const response = await fetch(`/api/v1/behaviors/${behaviorId}`, {
                method: 'DELETE'
            });

            if (response.ok) {
                this.showSuccess('行为删除成功');
                this.loadBehaviors();
            } else {
                const result = await response.json();
                this.showError('删除失败: ' + result.message);
            }
        } catch (error) {
            console.error('Failed to delete behavior:', error);
            this.showError('删除行为失败');
        }
    }

    viewActor(actorId) {
        // 实现查看Actor详情的逻辑
        console.log('View actor:', actorId);
    }

    showSuccess(message) {
        // 实现成功消息显示
        alert('成功: ' + message);
    }

    showError(message) {
        // 实现错误消息显示
        alert('错误: ' + message);
    }
}

// 初始化行为管理器
const behaviorManager = new BehaviorManager();

// 测试函数
function testFunction() {
    console.log('Test function called');
    alert('测试按钮被点击了！');
    
    // 测试API调用
    fetch('/api/v1/behaviors')
        .then(response => response.json())
        .then(data => {
            console.log('Test API response:', data);
            alert('API调用成功，返回 ' + (data.data ? data.data.length : 0) + ' 个行为');
        })
        .catch(error => {
            console.error('Test API error:', error);
            alert('API调用失败: ' + error.message);
        });
}
