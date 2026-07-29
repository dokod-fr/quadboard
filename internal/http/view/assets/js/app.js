document.addEventListener('DOMContentLoaded', () => {
    // --- Theme ---
    const themeToggle = document.getElementById('theme-toggle');
    const htmlElement = document.documentElement;
    const savedTheme = localStorage.getItem('quadboard-theme') || 'light';
    htmlElement.setAttribute('data-theme', savedTheme);

    if (themeToggle) {
        themeToggle.addEventListener('click', () => {
            const currentTheme = htmlElement.getAttribute('data-theme');
            const newTheme = currentTheme === 'light' ? 'dark' : 'light';
            htmlElement.setAttribute('data-theme', newTheme);
            localStorage.setItem('quadboard-theme', newTheme);
        });
    }

    // --- Catalog & Search ---
    const container = document.getElementById('resource-container');
    const searchInput = document.getElementById('search-input');

    if (!container) return;

    renderSkeletons();

    fetch('/api/v1/catalog')
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            return response.json();
        })
        .then(resources => { renderResources(resources); })
        .catch(err => {
            container.innerHTML = `<div class="empty-state"><i data-lucide="alert-circle"></i><p>Erreur lors du chargement des applications : ${err.message}</p></div>`;
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

        const groups = {};
        resources.forEach(app => {
            const groupName = app.Group || 'Default';
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

            groups[groupName].forEach(app => {
                const card = document.createElement('a');
                card.href = app.URL || '#';
                card.target = '_blank';
                card.rel = 'noopener noreferrer';
                card.className = `card ${!app.URL ? 'disabled' : ''}`;

                // --- Logo management and fallback ---
                let mediaHTML = '';
                if (app.Logo) {
                    // If no image found (404 du CDN), show default image
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
                        <div class="card-header"><div class="card-icon-wrapper">${mediaHTML}</div><h3 class="card-title">${app.Name}</h3></div>
                        ${descriptionHTML}
                    </div>
                    <div class="card-footer"><span class="card-open">${openText}</span></div>
                `;
                grid.appendChild(card);
            });
            container.appendChild(section);
        });
        if (window.lucide) window.lucide.createIcons();
    }

    if (searchInput) {
        searchInput.addEventListener('input', (e) => {
            const searchTerm = e.target.value.toLowerCase();
            const groups = document.querySelectorAll('.resource-group');
            groups.forEach(group => {
                let visibleCards = 0;
                const cards = group.querySelectorAll('.card');
                cards.forEach(card => {
                    const text = card.textContent.toLowerCase();
                    if (text.includes(searchTerm)) { card.style.display = ''; visibleCards++; } 
                    else { card.style.display = 'none'; }
                });
                group.style.display = visibleCards > 0 ? '' : 'none';
            });
        });
    }
});