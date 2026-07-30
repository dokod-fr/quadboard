import { initTheme } from './theme.js';
import { initRenderer, renderSkeletons, render, initGroupToggle } from './renderer.js';
import { initSearch } from './search.js';
import { state } from './state.js';

function init() {
    initTheme();
    initRenderer();
    initGroupToggle();
    initSearch();

    if (!document.getElementById('resource-container')) return;

    renderSkeletons();

    fetch('/api/v1/catalog')
        .then(res => res.ok ? res.json() : Promise.reject('Network error'))
        .then(data => {
            state.allResources = data;
            render(data);
        })
        .catch(err => {
            const container = document.getElementById('resource-container');
            if (container) {
                container.innerHTML = `<div class="empty-state"><i data-lucide="alert-circle"></i><p>Erreur de chargement : ${err}</p></div>`;
                lucide.createIcons();
            }
        });
}

document.addEventListener('DOMContentLoaded', init);