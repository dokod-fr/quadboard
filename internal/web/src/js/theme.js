export function initTheme() {
    const toggle = document.getElementById('theme-toggle');
    const html = document.documentElement;
    const favicon = document.getElementById('favicon');

    const apply = (theme) => {
        html.setAttribute('data-theme', theme);
        if (favicon) {
            favicon.href = theme === 'dark' ? '/assets/img/quadboard-bw.svg' : '/assets/img/quadboard-color.svg';
        }
    };

    const savedTheme = localStorage.getItem('quadboard-theme') || 'light';
    apply(savedTheme);

    if (toggle) {
        toggle.addEventListener('click', () => {
            const newTheme = html.getAttribute('data-theme') === 'light' ? 'dark' : 'light';
            apply(newTheme);
            localStorage.setItem('quadboard-theme', newTheme);
        });
    }
}