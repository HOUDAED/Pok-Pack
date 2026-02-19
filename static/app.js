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

function initCountUps() {
    document.querySelectorAll('.count-up[data-count]').forEach((el) => {
        const target = Number(el.getAttribute('data-count') || '0');
        if (!Number.isFinite(target) || target < 0) return;
        animateCountUp(el, target);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    initCountUps();
});
