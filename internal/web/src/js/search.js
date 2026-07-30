import { state } from './state.js';
import { render } from './renderer.js';

export function initSearch() {
    const input = document.getElementById('search-input');
    if (!input) return;

    input.addEventListener('input', (e) => {
        const term = e.target.value.toLowerCase();
        const filtered = state.allResources.filter(app => {
            const text = (app.Name + ' ' + app.Description + ' ' + (app.Group || '') + ' ' + (app.Icon || '')).toLowerCase();
            return text.includes(term);
        });
        render(filtered);
    });
}