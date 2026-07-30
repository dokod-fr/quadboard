document.addEventListener('DOMContentLoaded', () => {
    // --- Theme ---
    const themeToggle = document.getElementById('theme-toggle');
    const htmlElement = document.documentElement;
    const favicon = document.getElementById('favicon');

    function applyTheme(theme) {
        htmlElement.setAttribute('data-theme', theme);
        if (favicon) {
            favicon.href = theme === 'dark' ? '/assets/img/quadboard-bw.svg' : '/assets/img/quadboard-color.svg';
        }
    }

    const savedTheme = localStorage.getItem('quadboard-theme') || 'light';
    applyTheme(savedTheme);

    if (themeToggle) {
        themeToggle.addEventListener('click', () => {
            const currentTheme = htmlElement.getAttribute('data-theme');
            const newTheme = currentTheme === 'light' ? 'dark' : 'light';
            applyTheme(newTheme);
            localStorage.setItem('quadboard-theme', newTheme);
        });
    }

    // --- Catalog & Search ---
    const container = document.getElementById('resource-container');
    const searchInput = document.getElementById('search-input');

    if (!container) return;

    // Lecture de la config initiale depuis le HTML
    const groupByDefault = document.body.dataset.groupByDefault === 'true';
    
    let isGrouped = groupByDefault; 
    let activeGroupFilter = null;
    let allResources = [];

    const groupToggle = document.getElementById('group-toggle');
    if (groupToggle) {
        groupToggle.classList.toggle('active', isGrouped);
        groupToggle.addEventListener('click', () => {
            isGrouped = !isGrouped;
            groupToggle.classList.toggle('active', isGrouped);
            renderResources(allResources);
        });
    }

    // Affichage immédiat des skeletons pendant le fetch
    renderSkeletons();

    fetch('/api/v1/catalog')
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            return response.json();
        })
        .then(resources => {
            allResources = resources;
            renderResources(allResources);
        })
        .catch(err => {
            container.innerHTML = `<div class="empty-state"><i data-lucide="alert-circle"></i><p>Erreur lors du chargement : ${err.message}</p></div>`;
            if (window.lucide) window.lucide.createIcons();
        });

    function renderSkeletons() {
        container.innerHTML = '';
        const grid = document.createElement('div');
        grid.className = 'card-grid';
        for (let i = 0; i < 8; i++) {
            const card = document.createElement('div');
            card.className = 'card skeleton';
            card.innerHTML = `
                <div class="card-content">
                    <div class="card-header">
                        <div class="card-icon-wrapper skeleton-bg"><div class="skeleton-avatar"></div></div>
                        <div class="skeleton-line title"></div>
                    </div>
                    <div class="skeleton-line desc"></div>
                    <div class="skeleton-line desc short"></div>
                </div>
                <div class="card-footer"><div class="skeleton-line link"></div></div>
            `;
            grid.appendChild(card);
        }
        container.appendChild(grid);
    }

    function renderResources(resources) {
        container.innerHTML = '';
        if (resources.length === 0) {
            container.innerHTML = `<div class="empty-state"><i data-lucide="search-x"></i><p>Aucune application découverte pour le moment.</p></div>`;
            if (window.lucide) window.lucide.createIcons();
            return;
        }

        let resourcesToRender = resources;
        if (activeGroupFilter) {
            resourcesToRender = resources.filter(app => (app.Group || 'Default') === activeGroupFilter);
            
            const banner = document.createElement('div');
            banner.className = 'filter-banner';
            banner.innerHTML = `
                <span>Filtre actif : <strong>${activeGroupFilter}</strong></span>
                <button id="clear-filter-btn">Effacer le filtre</button>
            `;
            container.appendChild(banner);
            
            banner.querySelector('#clear-filter-btn').addEventListener('click', () => {
                activeGroupFilter = null;
                renderResources(allResources);
            });
        }

        if (isGrouped && !activeGroupFilter) {
            renderGroupedView(resourcesToRender);
        } else {
            renderFlatView(resourcesToRender);
        }
    }

    function renderGroupedView(resources) {
        const groups = {};
        resources.forEach(app => {
            const groupName = app.DisplayName || app.Group || 'Default';
            if (!groups[groupName]) groups[groupName] = [];
            groups[groupName].push(app);
        });

        const sortedGroupNames = Object.keys(groups).sort();

        sortedGroupNames.forEach(groupName => {
            const section = document.createElement('section');
            section.className = 'resource-group';
            section.dataset.groupName = groupName;
            section.innerHTML = `<h2 class="group-title">${groupName}</h2><div class="card-grid"></div>`;
            const grid = section.querySelector('.card-grid');
            
            groups[groupName].forEach(app => grid.appendChild(createCard(app, groupName)));
            container.appendChild(section);
        });
        if (window.lucide) window.lucide.createIcons();
    }

    function renderFlatView(resources) {
        const grid = document.createElement('div');
        grid.className = 'card-grid';
        resources.forEach(app => grid.appendChild(createCard(app, app.DisplayName || app.Group || 'Default')));
        container.appendChild(grid);
        if (window.lucide) window.lucide.createIcons();
    }

    function createCard(app, groupName) {
        // Si non autorisé, on crée une carte "désactivée" sans lien réel
        if (!app.Authorized) {
            const card = document.createElement('div');
            card.className = 'card disabled';
            card.innerHTML = `
                <div class="card-content">
                    <div class="card-header">
                        <div class="card-icon-wrapper"><i data-lucide="lock" class="card-icon"></i></div>
                        <h3 class="card-title">${app.Name}</h3>
                    </div>
                </div>
                <div class="card-footer">
                    <span class="card-group-badge" data-group="${app.Group || 'Default'}">${groupName}</span>
                    <span class="card-no-url">Restricted</span>
                </div>
            `;
            return card;
        }

        // Sinon, carte normale
        const card = document.createElement('a');
        card.href = app.URL || '#';
        card.target = '_blank';
        card.rel = 'noopener noreferrer';
        card.className = `card ${!app.URL ? 'disabled' : ''}`;

        let mediaHTML = '';
        if (app.Logo) {
            mediaHTML = `
                <img src="${app.Logo}" alt="${app.Name}" class="card-logo" onerror="this.style.display='none'; this.nextElementSibling.style.display='inline-block';">
                <i data-lucide="${app.Icon || 'box'}" class="card-icon" style="display: none;"></i>
            `;
        } else {
            mediaHTML = `<i data-lucide="${app.Icon || 'box'}" class="card-icon"></i>`;
        }

        let descriptionHTML = app.Description ? `<p class="card-desc">${app.Description}</p>` : '';
        let openText = app.URL 
            ? `Open <i data-lucide="external-link" style="width:12px; height:12px; display:inline-block; margin-left:2px;"></i>` 
            : '<span class="card-no-url">No URL</span>';

        card.innerHTML = `
            <div class="card-content">
                <div class="card-header">
                    <div class="card-icon-wrapper">${mediaHTML}</div>
                    <h3 class="card-title">${app.Name}</h3>
                </div>
                ${descriptionHTML}
            </div>
            <div class="card-footer">
                <span class="card-group-badge" data-group="${app.Group || 'Default'}">${groupName}</span>
                <span class="card-open">${openText}</span>
            </div>
        `;
        return card;
    }

    // Gestion du clic sur un badge de groupe
    container.addEventListener('click', (e) => {
        if (e.target.classList.contains('card-group-badge')) {
            e.preventDefault();
            activeGroupFilter = e.target.dataset.group;
            renderResources(allResources);
        }
    });

    // Recherche dynamique
    if (searchInput) {
        searchInput.addEventListener('input', (e) => {
            const searchTerm = e.target.value.toLowerCase();
            const filtered = allResources.filter(app => {
                const text = (app.Name + ' ' + app.Description + ' ' + (app.Group || '') + ' ' + (app.Icon || '')).toLowerCase();
                return text.includes(searchTerm);
            });
            renderResources(filtered);
        });
    }
});