import { state } from './state.js';

let container = null;

export function initRenderer() {
    container = document.getElementById('resource-container');
}

export function renderSkeletons() {
    if (!container) return;
    container.innerHTML = '<div class="card-grid">' + 
        Array(8).fill('<div class="card skeleton"><div class="card-content"><div class="card-header"><div class="card-icon-wrapper skeleton-bg"><div class="skeleton-avatar"></div></div><div class="skeleton-line title"></div></div><div class="skeleton-line desc"></div><div class="skeleton-line desc short"></div></div><div class="card-footer"><div class="skeleton-line link"></div></div></div>').join('') + 
    '</div>';
}

export function render(resources) {
    if (!container) return;
    container.innerHTML = '';

    if (resources.length === 0) {
        container.innerHTML = '<div class="empty-state"><i data-lucide="search-x"></i><p>Aucune application découverte pour le moment.</p></div>';
        lucide.createIcons();
        return;
    }

    let resourcesToRender = resources;
    if (state.activeGroupFilter) {
        resourcesToRender = resources.filter(app => (app.Group || 'Default') === state.activeGroupFilter);
        renderFilterBanner();
    }

    if (state.isGrouped && !state.activeGroupFilter) {
        renderGrouped(resourcesToRender);
    } else {
        renderFlat(resourcesToRender);
    }
}

function renderFilterBanner() {
    const banner = document.createElement('div');
    banner.className = 'filter-banner';
    banner.innerHTML = `<span>Filtre actif : <strong>${state.activeGroupFilter}</strong></span><button id="clear-filter-btn">Effacer le filtre</button>`;
    container.appendChild(banner);
    banner.querySelector('#clear-filter-btn').addEventListener('click', () => {
        state.activeGroupFilter = null;
        render(state.allResources);
    });
}

function renderGrouped(resources) {
    const groups = {};
    resources.forEach(app => {
        const name = app.DisplayName || app.Group || 'Default';
        if (!groups[name]) groups[name] = [];
        groups[name].push(app);
    });

    Object.keys(groups).sort().forEach(groupName => {
        const section = document.createElement('section');
        section.className = 'resource-group';
        section.innerHTML = `<h2 class="group-title">${groupName}</h2><div class="card-grid"></div>`;
        const grid = section.querySelector('.card-grid');
        groups[groupName].forEach(app => grid.appendChild(createCard(app, groupName)));
        container.appendChild(section);
    });
    lucide.createIcons();
}

function renderFlat(resources) {
    const grid = document.createElement('div');
    grid.className = 'card-grid';
    resources.forEach(app => grid.appendChild(createCard(app, app.DisplayName || app.Group || 'Default')));
    container.appendChild(grid);
    lucide.createIcons();
}

function createCard(app, groupName) {
    const card = document.createElement('a');
    card.href = app.URL || '#';
    card.target = '_blank';
    card.rel = 'noopener noreferrer';
    
    if (!app.Authorized) {
        card.className = 'card disabled';
        card.innerHTML = `
            <div class="card-content"><div class="card-header"><div class="card-icon-wrapper"><i data-lucide="lock" class="card-icon"></i></div><h3 class="card-title">${app.Name}</h3></div></div>
            <div class="card-footer"><span class="card-group-badge" data-group="${app.Group || 'Default'}">${groupName}</span><span class="card-no-url">Restricted</span></div>
        `;
        return card;
    }

    card.className = `card ${!app.URL ? 'disabled' : ''}`;

    const mediaHTML = app.Logo 
        ? `<img src="${app.Logo}" alt="${app.Name}" class="card-logo" onerror="this.style.display='none'; this.nextElementSibling.style.display='inline-block';"><i data-lucide="${app.Icon || 'box'}" class="card-icon" style="display: none;"></i>`
        : `<i data-lucide="${app.Icon || 'box'}" class="card-icon"></i>`;

    card.innerHTML = `
        <div class="card-content">
            <div class="card-header"><div class="card-icon-wrapper">${mediaHTML}</div><h3 class="card-title">${app.Name}</h3></div>
            ${app.Description ? `<p class="card-desc">${app.Description}</p>` : ''}
        </div>
        <div class="card-footer">
            <span class="card-group-badge" data-group="${app.Group || 'Default'}">${groupName}</span>
            <span class="card-open">${app.URL ? 'Open <i data-lucide="external-link" style="width:12px;height:12px;display:inline-block;margin-left:2px;"></i>' : '<span class="card-no-url">No URL</span>'}</span>
        </div>
    `;
    return card;
}

export function initGroupToggle() {
    if (!container) return;

    state.isGrouped = document.body.dataset.groupByDefault === 'true';
    const groupToggle = document.getElementById('group-toggle');
    if (groupToggle) {
        groupToggle.classList.toggle('active', state.isGrouped);
        groupToggle.addEventListener('click', () => {
            state.isGrouped = !state.isGrouped;
            groupToggle.classList.toggle('active', state.isGrouped);
            render(state.allResources);
        });
    }

    container.addEventListener('click', (e) => {
        if (e.target.classList.contains('card-group-badge')) {
            e.preventDefault();
            state.activeGroupFilter = e.target.dataset.group;
            render(state.allResources);
        }
    });
}