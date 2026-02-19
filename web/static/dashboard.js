function animateCountUp(el, target) {
    const duration = 750;
    const start = performance.now();
    const from = 0;

    function tick(now) {
        const t = Math.min(1, (now - start) / duration);
        const eased = 1 - Math.pow(1 - t, 3);
        const value = Math.round(from + (target - from) * eased);
        el.textContent = String(value);
        if (t < 1) requestAnimationFrame(tick);
    }

    requestAnimationFrame(tick);
}

function initCountUp() {
    const el = document.querySelector('.count-up[data-count]');
    if (!el) return;
    const target = Number(el.getAttribute('data-count') || '0');
    if (!Number.isFinite(target) || target < 0) return;
    animateCountUp(el, target);
}

function normalizeTypeName(name) {
    return String(name || '')
        .trim()
        .toLowerCase()
        .replace(/\s+/g, '-');
}

function initTypeBadges() {
    document.querySelectorAll('.card-types[data-types]').forEach((el) => {
        const raw = el.getAttribute('data-types') || el.textContent || '';
        const parts = raw
            .split(',')
            .map((s) => s.trim())
            .filter(Boolean);
        if (!parts.length) return;

        el.textContent = '';
        parts.forEach((t) => {
            const span = document.createElement('span');
            span.className = 'type-badge type-' + normalizeTypeName(t);
            span.textContent = t;
            el.appendChild(span);
        });
    });
}

function initPackOpenFX() {
    const form = document.querySelector('form.open-pack');
    if (!form) return;
    const btn = form.querySelector('button[type="submit"]');
    form.addEventListener('submit', () => {
        document.body.classList.add('is-opening-pack');
        if (btn) {
            btn.disabled = true;
            btn.classList.add('is-loading');
        }
    });
}

document.addEventListener('DOMContentLoaded', () => {
    initCountUp();
    initTypeBadges();
    initPackOpenFX();
});
