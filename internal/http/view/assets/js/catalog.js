document.addEventListener('DOMContentLoaded', () => {
    const grid = document.getElementById('catalog-grid');

    fetch('/api/v1/catalog')
        .then(response => {
            if (!response.ok) {
                throw new Error('Erreur réseau lors de la récupération du catalogue');
            }
            return response.json();
        })
        .then(resources => {
            grid.innerHTML = '';

            if (resources.length === 0) {
                grid.innerHTML = `
                    <div class="empty-state">
                        <i data-lucide="search-x"></i>
                        <p>Aucune application découverte pour le moment.</p>
                    </div>
                `;

                if (window.lucide) window.lucide.createIcons();
                return;
            }

            // Card generation
            resources.forEach(app => {
                const card = document.createElement('a');
                card.href = app.URL || '#';
                card.target = '_blank';
                card.rel = 'noopener noreferrer';
                card.className = 'card';

                let mediaHTML = '';
                if (app.Logo) {
                    mediaHTML = `<img src="${app.Logo}" alt="${app.Name}" class="card-logo">`;
                } else {
                    const iconName = app.Icon || 'box';
                    mediaHTML = `<i data-lucide="${iconName}" class="card-icon"></i>`;
                }

                const descriptionHTML = app.Description 
                    ? `<p class="card-desc">${app.Description}</p>` 
                    : '';

                const groupText = app.Group || 'Default';
                const openText = app.URL 
                    ? `Open <i data-lucide="external-link" style="width:12px; height:12px; display:inline-block; margin-left:2px;"></i>` 
                    : '<span style="color: var(--text-muted);">No URL</span>';

                card.innerHTML = `
                    <div class="card-content">
                        <div class="card-header">
                            ${mediaHTML}
                            <h3 class="card-title">${app.Name}</h3>
                        </div>
                        ${descriptionHTML}
                    </div>
                    <div class="card-footer">
                        <span class="card-provider">${groupText}</span>
                        <span class="card-open">${openText}</span>
                    </div>
                `;

                grid.appendChild(card);
            });

            // Draw icons only when the card is loaded
            if (window.lucide) {
                window.lucide.createIcons();
            }
        })
        .catch(err => {
            console.error(err);
            grid.innerHTML = `
                <div class="empty-state">
                    <i data-lucide="alert-circle" style="color: var(--text-muted);"></i>
                    <p>Erreur lors du chargement des applications : ${err.message}</p>
                </div>
            `;
            if (window.lucide) window.lucide.createIcons();
        });
});