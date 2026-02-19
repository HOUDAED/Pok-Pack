function animateCountUp(el, target) {
    const duration = 850;
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

function initCountUps() {
    document.querySelectorAll('.count-up[data-count]').forEach((el) => {
        const target = Number(el.getAttribute('data-count') || '0');
        if (!Number.isFinite(target) || target < 0) return;
        animateCountUp(el, target);
    });
}

function normalizeTypeName(name) {
    return String(name || '')
        .trim()
        .toLowerCase()
        .replace(/\s+/g, '-')
        .replace(/[^a-z0-9-]/g, '');
}

function renderTypeBadges(el) {
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
}

function initTypeBadges() {
    document.querySelectorAll('.types[data-types]').forEach((el) => renderTypeBadges(el));
}

function initPackOpenFX() {
    document.querySelectorAll('form.open-pack').forEach((form) => {
        const btn = form.querySelector('button[type="submit"]');
        form.addEventListener('submit', () => {
            document.body.classList.add('is-opening-pack');
            if (btn) {
                btn.disabled = true;
                btn.setAttribute('aria-busy', 'true');
                btn.style.filter = 'saturate(1.08)';
            }
        });
    });
}

function parseSeries(el) {
    try {
        const raw = el.getAttribute('data-series') || '[]';
        const arr = JSON.parse(raw);
        if (!Array.isArray(arr)) return [];
        return arr
            .map((x) => ({ type: String(x.type || ''), count: Number(x.count || 0) }))
            .filter((x) => x.type && Number.isFinite(x.count) && x.count > 0);
    } catch {
        return [];
    }
}

function typeClass(type) {
    return 'type-' + normalizeTypeName(type);
}

function drawBars(chartEl) {
    const svg = chartEl.querySelector('svg');
    const g = svg ? svg.querySelector('g.bars') : null;
    if (!svg || !g) return;

    const series = parseSeries(chartEl);
    g.textContent = '';
    if (!series.length) {
        const txt = document.createElementNS('http://www.w3.org/2000/svg', 'text');
        txt.setAttribute('x', '18');
        txt.setAttribute('y', '46');
        txt.setAttribute('fill', 'rgba(255,255,255,0.65)');
        txt.setAttribute('font-size', '14');
        txt.textContent = 'Aucune donnée';
        g.appendChild(txt);
        return;
    }

    const w = 600;
    const h = 220;
    const padding = { left: 18, right: 18, top: 18, bottom: 42 };
    const max = Math.max(...series.map((s) => s.count));
    const innerW = w - padding.left - padding.right;
    const innerH = h - padding.top - padding.bottom;
    const gap = 10;
    const barW = Math.max(12, Math.floor((innerW - gap * (series.length - 1)) / series.length));

    series.forEach((s, i) => {
        const x = padding.left + i * (barW + gap);
        const bh = Math.max(2, Math.round((s.count / max) * innerH));
        const y = padding.top + (innerH - bh);

        const rect = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
        rect.setAttribute('x', String(x));
        rect.setAttribute('y', String(padding.top + innerH));
        rect.setAttribute('width', String(barW));
        rect.setAttribute('height', '0');
        rect.setAttribute('rx', '10');
        rect.setAttribute('fill', 'rgba(255,255,255,0.10)');
        rect.setAttribute('stroke', 'rgba(255,255,255,0.12)');

            const glow = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
            glow.setAttribute('x', String(x));
            glow.setAttribute('y', String(y));
            glow.setAttribute('width', String(barW));
            glow.setAttribute('height', String(bh));
            glow.setAttribute('rx', '10');
            glow.setAttribute('class', typeClass(s.type));

        const label = document.createElementNS('http://www.w3.org/2000/svg', 'text');
        label.setAttribute('x', String(x + barW / 2));
        label.setAttribute('y', String(h - 26));
        label.setAttribute('text-anchor', 'middle');
        label.setAttribute('fill', 'rgba(255,255,255,0.68)');
        label.setAttribute('font-size', '12');
        label.setAttribute('font-weight', '800');
        label.textContent = s.type;

        const count = document.createElementNS('http://www.w3.org/2000/svg', 'text');
        count.setAttribute('x', String(x + barW / 2));
        count.setAttribute('y', String(h - 10));
        count.setAttribute('text-anchor', 'middle');
        count.setAttribute('fill', 'rgba(255,255,255,0.90)');
        count.setAttribute('font-size', '12');
        count.setAttribute('font-weight', '900');
        count.textContent = String(s.count);

        g.appendChild(rect);
        g.appendChild(glow);
        g.appendChild(label);
        g.appendChild(count);

        // Animate
        glow.animate(
            [
                { transform: 'translate(0px,' + (innerH - bh) + 'px) scale(1,0.01)', opacity: 0.1 },
                { transform: 'translate(0px,0px) scale(1,1)', opacity: 1 },
            ],
            { duration: 640, easing: 'cubic-bezier(0.2, 0.8, 0.2, 1)', fill: 'forwards', delay: 40 * i }
        );
    });
}

function polarToCartesian(cx, cy, r, a) {
    const rad = ((a - 90) * Math.PI) / 180;
    return { x: cx + r * Math.cos(rad), y: cy + r * Math.sin(rad) };
}

function arcPath(cx, cy, r, startAngle, endAngle) {
    const start = polarToCartesian(cx, cy, r, endAngle);
    const end = polarToCartesian(cx, cy, r, startAngle);
    const largeArcFlag = endAngle - startAngle <= 180 ? '0' : '1';
    return `M ${start.x} ${start.y} A ${r} ${r} 0 ${largeArcFlag} 0 ${end.x} ${end.y}`;
}

function drawRadial(chartEl) {
    const svg = chartEl.querySelector('svg');
    const g = svg ? svg.querySelector('g.radial') : null;
    if (!svg || !g) return;

    const series = parseSeries(chartEl);
    g.textContent = '';
    if (!series.length) {
        const txt = document.createElementNS('http://www.w3.org/2000/svg', 'text');
        txt.setAttribute('x', '20');
        txt.setAttribute('y', '40');
        txt.setAttribute('fill', 'rgba(255,255,255,0.65)');
        txt.setAttribute('font-size', '14');
        txt.textContent = 'Aucune donnée';
        g.appendChild(txt);
        return;
    }

    const total = series.reduce((a, s) => a + s.count, 0);
    const cx = 130;
    const cy = 130;
    const r = 92;
    const strokeW = 18;

    let angle = 0;
    series.forEach((s, i) => {
        const frac = s.count / total;
        const sweep = Math.max(4, frac * 360);
        const start = angle;
        const end = angle + sweep;
        angle += sweep;

        const path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
        path.setAttribute('d', arcPath(cx, cy, r, start, end));
        path.setAttribute('fill', 'none');
        path.setAttribute('stroke-linecap', 'round');
        path.setAttribute('stroke-width', String(strokeW));
        path.setAttribute('class', typeClass(s.type));

        const base = document.createElementNS('http://www.w3.org/2000/svg', 'path');
        base.setAttribute('d', arcPath(cx, cy, r, start, end));
        base.setAttribute('fill', 'none');
        base.setAttribute('stroke-linecap', 'round');
        base.setAttribute('stroke-width', String(strokeW));
        base.setAttribute('stroke', 'rgba(255,255,255,0.08)');

        g.appendChild(base);
        g.appendChild(path);

        const len = path.getTotalLength();
        path.style.strokeDasharray = String(len);
        path.style.strokeDashoffset = String(len);
        path.animate(
            [{ strokeDashoffset: len }, { strokeDashoffset: 0 }],
            { duration: 760, easing: 'cubic-bezier(0.2, 0.8, 0.2, 1)', fill: 'forwards', delay: 45 * i }
        );
    });

    const center = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
    center.setAttribute('cx', String(cx));
    center.setAttribute('cy', String(cy));
    center.setAttribute('r', String(r - strokeW));
    center.setAttribute('fill', 'rgba(15, 23, 42, 0.52)');
    center.setAttribute('stroke', 'rgba(255,255,255,0.10)');
    g.appendChild(center);

    const txt = document.createElementNS('http://www.w3.org/2000/svg', 'text');
    txt.setAttribute('x', String(cx));
    txt.setAttribute('y', String(cy + 6));
    txt.setAttribute('text-anchor', 'middle');
    txt.setAttribute('fill', 'rgba(255,255,255,0.90)');
    txt.setAttribute('font-size', '18');
    txt.setAttribute('font-weight', '900');
    txt.textContent = String(total);
    g.appendChild(txt);
}

function initCharts() {
    document.querySelectorAll('[data-chart="bars"]').forEach(drawBars);
    document.querySelectorAll('[data-chart="radial"]').forEach(drawRadial);
}

function initChartLegends() {
    document.querySelectorAll('.chart[data-series]').forEach((chart) => {
        const legend = chart.querySelector('[data-chart-legend]');
        if (!legend) return;

        const series = parseSeries(chart);
        legend.textContent = '';
        if (!series.length) return;

        series.forEach((s) => {
            const item = document.createElement('div');
            item.className = 'legend-item';

            const dot = document.createElement('span');
            dot.className = 'legend-dot ' + typeClass(s.type);

            const label = document.createElement('span');
            label.textContent = String(s.type || '').toLowerCase();

            item.appendChild(dot);
            item.appendChild(label);
            legend.appendChild(item);
        });
    });
}

function initCollectionFilters() {
    const grid = document.querySelector('[data-collection-grid]');
    if (!grid) return;

    const search = document.querySelector('[data-collection-search]');
    const type = document.querySelector('[data-collection-type]');

    function apply() {
        const q = (search && search.value ? search.value : '').trim().toLowerCase();
        const t = (type && type.value ? type.value : '').trim().toLowerCase();

        grid.querySelectorAll('.poke-card').forEach((card) => {
            const name = String(card.getAttribute('data-name') || '').toLowerCase();
            const types = String(card.getAttribute('data-types') || '').toLowerCase();
            const okName = !q || name.includes(q);
            const okType = !t || types.split(',').map((s) => s.trim()).includes(t) || types.includes(t);
            card.style.display = okName && okType ? '' : 'none';
        });
    }

    if (search) search.addEventListener('input', apply);
    if (type) type.addEventListener('input', apply);
}

function initAddPreview() {
    const idInput = document.querySelector('[data-add-card-id]');
    const img = document.querySelector('[data-preview-img]');
    const nameEl = document.querySelector('[data-preview-name]');
    const typesEl = document.querySelector('[data-preview-types]');
    if (!idInput || !img || !nameEl || !typesEl) return;

    let lastReq = 0;

    async function load(id) {
        const reqId = ++lastReq;
        nameEl.textContent = 'Chargement…';
        typesEl.textContent = '';
        img.style.opacity = '0';
        img.removeAttribute('src');

        try {
            const res = await fetch('https://pokeapi.co/api/v2/pokemon/' + encodeURIComponent(String(id)));
            if (!res.ok) throw new Error('HTTP ' + res.status);
            const p = await res.json();
            if (reqId !== lastReq) return;

            const types = Array.isArray(p.types) ? p.types.map((t) => t && t.type && t.type.name).filter(Boolean) : [];
            const typesStr = types.join(', ');

            const oa = p && p.sprites && p.sprites.other && p.sprites.other['official-artwork'] && p.sprites.other['official-artwork'].front_default;
            const fallback = p && p.sprites && p.sprites.front_default;
            const url = oa || fallback || '';

            nameEl.textContent = p.name || '—';
            typesEl.setAttribute('data-types', typesStr);
            typesEl.textContent = typesStr || '—';
            renderTypeBadges(typesEl);

            if (url) {
                img.setAttribute('src', url);
                img.setAttribute('alt', p.name || 'Pokemon');
                img.onload = () => { img.style.opacity = '1'; };
            }
        } catch {
            if (reqId !== lastReq) return;
            nameEl.textContent = 'Introuvable';
            typesEl.textContent = 'Vérifie l’ID (1 à 1025).';
        }
    }

    let timer = 0;
    idInput.addEventListener('input', () => {
        const v = Number(idInput.value);
        window.clearTimeout(timer);
        if (!Number.isFinite(v) || v <= 0) return;
        timer = window.setTimeout(() => load(v), 220);
    });
}

function initNewCardsEntrance() {
    const cards = document.querySelectorAll('[data-new-card]');
    if (!cards.length) return;
    cards.forEach((el, i) => {
        el.animate(
            [
                { opacity: 0, transform: 'translateY(10px) scale(0.98)' },
                { opacity: 1, transform: 'translateY(0px) scale(1)' },
            ],
            { duration: 420, easing: 'cubic-bezier(0.2, 0.8, 0.2, 1)', fill: 'forwards', delay: 90 + 80 * i }
        );
    });
}

document.addEventListener('DOMContentLoaded', () => {
    initCountUps();
    initTypeBadges();
    initPackOpenFX();
    initCharts();
    initChartLegends();
    initCollectionFilters();
    initAddPreview();
    initNewCardsEntrance();
});
