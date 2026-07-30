const QuadBoard = {
    // --- État global de l'application ---
    state: {
        isGrouped: false,
        activeGroupFilter: null,
        allResources: []
    },

    // --- Module de Thème (Dark/Light) ---
    theme: {
        init() {
            this.toggle = document.getElementById('theme-toggle');
            this.html = document.documentElement;
            this.favicon = document.getElementById('favicon');

            const savedTheme = localStorage.getItem('quadboard-theme') || 'light';
            this.apply(savedTheme);

            if (this.toggle) {
                this.toggle.addEventListener('click', () => {
                    const newTheme = this.html.getAttribute('data-theme') === 'light' ? 'dark' : 'light';
                    this.apply(newTheme);
                    localStorage.setItem('quadboard-theme', newTheme);
                });
            }
        },
        apply(theme) {
            this.html.setAttribute('data-theme', theme);
            if (this.favicon) {
                this.favicon.href = theme === 'dark' ? '/assets/img/quadboard-bw.svg' : '/assets/img/quadboard-color.svg';
            }
        }
    },

    // --- Module de Rendu UI ---
    renderer: {
        container: null,

        init() {
            this.container = document.getElementById('resource-container');
        },

        renderSkeletons() {
            if (!this.container) return;
            this.container.innerHTML = '<div class="card-grid">' + 
                Array(8).fill('<div class="card skeleton"><div class="card-content"><div class="card-header"><div class="card-icon-wrapper skeleton-bg"><div class="skeleton-avatar"></div></div><div class="skeleton-line title"></div></div><div class="skeleton-line desc"></div><div class="skeleton-line desc short"></div></div><div class="card-footer"><div class="skeleton-line link"></div></div></div>').join('') + 
            '</div>';
        },

        render(resources) {
            if (!this.container) return;
            this.container.innerHTML = '';

            if (resources.length === 0) {
                this.container.innerHTML = '<div class="empty-state"><i data-lucide="search-x"></i><p>Aucune application découverte pour le moment.</p></div>';
                lucide.createIcons();
                return;
            }

            let resourcesToRender = resources;
            if (QuadBoard.state.activeGroupFilter) {
                resourcesToRender = resources.filter(app => (app.Group || 'Default') === QuadBoard.state.activeGroupFilter);
                this.renderFilterBanner();
            }

            if (QuadBoard.state.isGrouped && !QuadBoard.state.activeGroupFilter) {
                this.renderGrouped(resourcesToRender);
            } else {
                this.renderFlat(resourcesToRender);
            }
        },

        renderFilterBanner() {
            const banner = document.createElement('div');
            banner.className = 'filter-banner';
            banner.innerHTML = `<span>Filtre actif : <strong>${QuadBoard.state.activeGroupFilter}</strong></span><button id="clear-filter-btn">Effacer le filtre</button>`;
            this.container.appendChild(banner);
            banner.querySelector('#clear-filter-btn').addEventListener('click', () => {
                QuadBoard.state.activeGroupFilter = null;
                this.render(QuadBoard.state.allResources);
            });
        },

        renderGrouped(resources) {
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
                groups[groupName].forEach(app => grid.appendChild(this.createCard(app, groupName)));
                this.container.appendChild(section);
            });
            lucide.createIcons();
        },

        renderFlat(resources) {
            const grid = document.createElement('div');
            grid.className = 'card-grid';
            resources.forEach(app => grid.appendChild(this.createCard(app, app.DisplayName || app.Group || 'Default')));
            this.container.appendChild(grid);
            lucide.createIcons();
        },

        createCard(app, groupName) {
            const card = document.createElement('a');
            card.href = app.URL || '#';
            card.target = '_blank';
            card.rel = 'noopener noreferrer';
            
            if (!app.Authorized) {
                card.className = 'card disabled';
                card.innerHTML = `
                    <div class="card-content"><div class="card-header"><div class="card-icon-wrapper"><i data-lucide="lock" class="card-icon"></i></div><h3 class="card-title">${app.Name}</h3></div></div>
                    <div class="card-footer"><span class="card-group-badge">${groupName}</span><span class="card-no-url">Restricted</span></div>
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
    },

    // --- Module de gestion du Catalogue (API & Events) ---
    catalog: {
        init() {
            if (!QuadBoard.renderer.container) return;

            QuadBoard.state.isGrouped = document.body.dataset.groupByDefault === 'true';
            const groupToggle = document.getElementById('group-toggle');
            if (groupToggle) {
                groupToggle.classList.toggle('active', QuadBoard.state.isGrouped);
                groupToggle.addEventListener('click', () => {
                    QuadBoard.state.isGrouped = !QuadBoard.state.isGrouped;
                    groupToggle.classList.toggle('active', QuadBoard.state.isGrouped);
                    QuadBoard.renderer.render(QuadBoard.state.allResources);
                });
            }

            // Délégation de clic pour les badges de groupe
            QuadBoard.renderer.container.addEventListener('click', (e) => {
                if (e.target.classList.contains('card-group-badge')) {
                    e.preventDefault();
                    QuadBoard.state.activeGroupFilter = e.target.dataset.group;
                    QuadBoard.renderer.render(QuadBoard.state.allResources);
                }
            });

            this.fetch();
        },

        fetch() {
            QuadBoard.renderer.renderSkeletons();
            fetch('/api/v1/catalog')
                .then(res => res.ok ? res.json() : Promise.reject('Network error'))
                .then(data => {
                    QuadBoard.state.allResources = data;
                    QuadBoard.renderer.render(data);
                })
                .catch(err => {
                    QuadBoard.renderer.container.innerHTML = `<div class="empty-state"><i data-lucide="alert-circle"></i><p>Erreur de chargement : ${err}</p></div>`;
                    lucide.createIcons();
                });
        }
    },

    // --- Module de Recherche ---
    search: {
        init() {
            const input = document.getElementById('search-input');
            if (!input) return;

            input.addEventListener('input', (e) => {
                const term = e.target.value.toLowerCase();
                const filtered = QuadBoard.state.allResources.filter(app => {
                    const text = (app.Name + ' ' + app.Description + ' ' + (app.Group || '') + ' ' + (app.Icon || '')).toLowerCase();
                    return text.includes(term);
                });
                QuadBoard.renderer.render(filtered);
            });
        }
    },

    // --- Initialisation globale ---
    init() {
        this.theme.init();
        this.renderer.init();
        this.catalog.init();
        this.search.init();
    }
};

document.addEventListener('DOMContentLoaded', () => QuadBoard.init());