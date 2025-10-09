// 关系图展示功能
let graphData = {
    nodes: [],
    links: []
};

let svg, simulation, tooltip;
let width, height;

// 关系类型映射
const relationshipTypes = {
    'contains': { name: '包含关系', strength: 'strong', color: '#e74c3c' },
    'composes': { name: '组合关系', strength: 'strong', color: '#e74c3c' },
    'owns': { name: '拥有关系', strength: 'strong', color: '#e74c3c' },
    'relates_to': { name: '关联关系', strength: 'weak', color: '#3498db' },
    'depends_on': { name: '依赖关系', strength: 'weak', color: '#3498db' },
    'influences': { name: '影响关系', strength: 'weak', color: '#3498db' },
    'collaborates': { name: '协作关系', strength: 'weak', color: '#3498db' }
};

// 事物类型颜色映射
const thingTypeColors = {
    'person': '#e74c3c',
    'machine': '#f39c12',
    'object': '#9b59b6'
};

// 初始化图形
function initGraph() {
    const container = document.getElementById('graph');
    width = container.offsetWidth;
    height = 600;

    svg = d3.select('#graph')
        .append('svg')
        .attr('width', width)
        .attr('height', height);

    // 创建工具提示
    tooltip = d3.select('body')
        .append('div')
        .attr('class', 'tooltip')
        .style('opacity', 0);

    // 创建力导向图模拟
    simulation = d3.forceSimulation()
        .force('link', d3.forceLink().id(d => d.id).distance(100))
        .force('charge', d3.forceManyBody().strength(-300))
        .force('center', d3.forceCenter(width / 2, height / 2))
        .force('collision', d3.forceCollide().radius(30));

    // 绑定事件
    bindEvents();
    
    // 加载数据
    loadGraphData();
}

// 绑定事件
function bindEvents() {
    // 筛选控制
    document.getElementById('relationshipType').addEventListener('change', filterGraph);
    document.getElementById('thingType').addEventListener('change', filterGraph);
    document.getElementById('refreshGraph').addEventListener('click', loadGraphData);
    document.getElementById('resetView').addEventListener('click', resetView);
}

// 加载图形数据
async function loadGraphData() {
    try {
        // 加载事物数据
        const thingsResponse = await fetch('/api/v1/things');
        const thingsData = await thingsResponse.json();
        
        // 加载关系数据
        const relationshipsResponse = await fetch('/api/v1/relationships');
        const relationshipsData = await relationshipsResponse.json();

        // 构建节点数据
        // Things API 返回 {success: true, data: [...], count: N}
        const things = thingsData.data || [];
        graphData.nodes = things.map(thing => ({
            id: thing.id,
            name: thing.name,
            type: thing.type,
            description: thing.description
        }));

        // 构建连接数据
        // Relationships API 返回 {success: true, data: {data: [...], count: N}}
        const relationships = relationshipsData.data?.data || [];
        
        // 创建节点ID集合用于验证
        const nodeIds = new Set(graphData.nodes.map(n => n.id));
        
        // 只保留有效的关系（源和目标节点都存在）
        graphData.links = relationships
            .filter(rel => {
                if (!nodeIds.has(rel.sourceId) || !nodeIds.has(rel.targetId)) {
                    console.warn(`关系 ${rel.name} 引用了不存在的节点: source=${rel.sourceId}, target=${rel.targetId}`);
                    return false;
                }
                return true;
            })
            .map(rel => ({
                source: rel.sourceId,
                target: rel.targetId,
                type: rel.type,
                name: rel.name,
                description: rel.description,
                properties: rel.properties
            }));

        // 应用筛选
        applyFilters();
        
        // 更新图形
        updateGraph();
        
    } catch (error) {
        console.error('加载图形数据失败:', error);
        alert('加载数据失败: ' + error.message);
    }
}

// 应用筛选
function applyFilters() {
    const relationshipTypeFilter = document.getElementById('relationshipType').value;
    const thingTypeFilter = document.getElementById('thingType').value;

    let filteredLinks = graphData.links;
    let filteredNodes = graphData.nodes;

    // 按关系类型筛选
    if (relationshipTypeFilter) {
        filteredLinks = graphData.links.filter(link => link.type === relationshipTypeFilter);
    }

    // 按事物类型筛选
    if (thingTypeFilter) {
        filteredNodes = graphData.nodes.filter(node => node.type === thingTypeFilter);
        
        // 只保留与筛选节点相关的关系
        const nodeIds = new Set(filteredNodes.map(node => node.id));
        filteredLinks = filteredLinks.filter(link => 
            nodeIds.has(link.source) && nodeIds.has(link.target)
        );
    }

    // 更新数据
    graphData.filteredNodes = filteredNodes;
    graphData.filteredLinks = filteredLinks;
}

// 筛选图形
function filterGraph() {
    applyFilters();
    updateGraph();
}

// 更新图形
function updateGraph() {
    if (!graphData.filteredNodes || !graphData.filteredLinks) {
        return;
    }

    // 清除现有图形
    svg.selectAll('*').remove();

    // 创建连接线
    const link = svg.append('g')
        .selectAll('line')
        .data(graphData.filteredLinks)
        .enter().append('line')
        .attr('class', d => `link ${relationshipTypes[d.type]?.strength || 'weak'}`)
        .attr('stroke', d => relationshipTypes[d.type]?.color || '#999')
        .attr('stroke-width', d => relationshipTypes[d.type]?.strength === 'strong' ? 3 : 1)
        .attr('stroke-dasharray', d => relationshipTypes[d.type]?.strength === 'weak' ? '5,5' : null)
        .on('mouseover', showLinkTooltip)
        .on('mouseout', hideTooltip)
        .on('click', showLinkDetails);

    // 创建节点
    const node = svg.append('g')
        .selectAll('circle')
        .data(graphData.filteredNodes)
        .enter().append('circle')
        .attr('class', d => `node ${d.type}`)
        .attr('r', 15)
        .attr('fill', d => thingTypeColors[d.type] || '#95a5a6')
        .call(d3.drag()
            .on('start', dragstarted)
            .on('drag', dragged)
            .on('end', dragended))
        .on('mouseover', showNodeTooltip)
        .on('mouseout', hideTooltip)
        .on('click', showNodeDetails);

    // 创建节点标签
    const nodeLabel = svg.append('g')
        .selectAll('text')
        .data(graphData.filteredNodes)
        .enter().append('text')
        .attr('class', 'node-label')
        .text(d => d.name)
        .attr('dy', 25);

    // 创建连接标签
    const linkLabel = svg.append('g')
        .selectAll('text')
        .data(graphData.filteredLinks)
        .enter().append('text')
        .attr('class', 'link-label')
        .text(d => relationshipTypes[d.type]?.name || d.type)
        .attr('dy', -5);

    // 更新力导向图
    simulation
        .nodes(graphData.filteredNodes)
        .force('link')
        .links(graphData.filteredLinks);

    simulation.alpha(1).restart();

    // 更新位置
    simulation.on('tick', () => {
        link
            .attr('x1', d => d.source.x)
            .attr('y1', d => d.source.y)
            .attr('x2', d => d.target.x)
            .attr('y2', d => d.target.y);

        node
            .attr('cx', d => d.x)
            .attr('cy', d => d.y);

        nodeLabel
            .attr('x', d => d.x)
            .attr('y', d => d.y);

        linkLabel
            .attr('x', d => (d.source.x + d.target.x) / 2)
            .attr('y', d => (d.source.y + d.target.y) / 2);
    });
}

// 显示连接工具提示
function showLinkTooltip(event, d) {
    tooltip.transition()
        .duration(200)
        .style('opacity', .9);
    
    tooltip.html(`
        <h5>${relationshipTypes[d.type]?.name || d.type}</h5>
        <p><strong>源:</strong> ${d.source.name || d.source}</p>
        <p><strong>目标:</strong> ${d.target.name || d.target}</p>
        <p><strong>描述:</strong> ${d.description || '无描述'}</p>
    `)
    .style('left', (event.pageX + 10) + 'px')
    .style('top', (event.pageY - 10) + 'px');
}

// 显示节点工具提示
function showNodeTooltip(event, d) {
    tooltip.transition()
        .duration(200)
        .style('opacity', .9);
    
    tooltip.html(`
        <h5>${d.name}</h5>
        <p><strong>类型:</strong> ${d.type}</p>
        <p><strong>描述:</strong> ${d.description || '无描述'}</p>
    `)
    .style('left', (event.pageX + 10) + 'px')
    .style('top', (event.pageY - 10) + 'px');
}

// 隐藏工具提示
function hideTooltip() {
    tooltip.transition()
        .duration(500)
        .style('opacity', 0);
}

// 显示连接详情
function showLinkDetails(event, d) {
    const details = document.getElementById('relationshipDetails');
    details.innerHTML = `
        <h5>${relationshipTypes[d.type]?.name || d.type}</h5>
        <p><strong>源事物:</strong> ${d.source.name || d.source}</p>
        <p><strong>目标事物:</strong> ${d.target.name || d.target}</p>
        <p><strong>描述:</strong> ${d.description || '无描述'}</p>
        <p><strong>关系强度:</strong> ${relationshipTypes[d.type]?.strength === 'strong' ? '强关联' : '弱关联'}</p>
        ${d.properties ? `<p><strong>属性:</strong> ${JSON.stringify(d.properties, null, 2)}</p>` : ''}
    `;
}

// 显示节点详情
function showNodeDetails(event, d) {
    const details = document.getElementById('thingDetails');
    details.innerHTML = `
        <h5>${d.name}</h5>
        <p><strong>类型:</strong> ${d.type}</p>
        <p><strong>描述:</strong> ${d.description || '无描述'}</p>
        <p><strong>ID:</strong> ${d.id}</p>
    `;
}

// 拖拽事件
function dragstarted(event, d) {
    if (!event.active) simulation.alphaTarget(0.3).restart();
    d.fx = d.x;
    d.fy = d.y;
}

function dragged(event, d) {
    d.fx = event.x;
    d.fy = event.y;
}

function dragended(event, d) {
    if (!event.active) simulation.alphaTarget(0);
    d.fx = null;
    d.fy = null;
}

// 重置视图
function resetView() {
    if (simulation) {
        simulation.alpha(1).restart();
    }
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', initGraph);
