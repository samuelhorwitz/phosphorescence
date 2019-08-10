(async function() {
    if (!/\b(iPhone|iPod)\b/.test(navigator.userAgent)) {
        return;
    }

    const debug = false;
    const offset = 0;
    const parallaxMultiplier = 0.8;

    let debugEl;
    if (debug) {
        debugEl = document.createElement('div');
        debugEl.style.color = 'white';
        debugEl.style.backgroundColor = 'magenta';
        debugEl.style.fontSize = '16px';
        debugEl.style.position = 'fixed';
        debugEl.style.bottom = '0px';
        document.body.appendChild(debugEl);
    }
    let bgCanvas = document.createElement('canvas');
    bgCanvas.classList.add('bgIosCanvas');
    bgCanvas.width = innerWidth * devicePixelRatio;
    bgCanvas.height = innerHeight * devicePixelRatio;
    bgCanvas.style.width = `${innerWidth}px`;
    bgCanvas.style.height = `${innerHeight}px`;
    let ctx = bgCanvas.getContext('2d');
    ctx.scale(devicePixelRatio, devicePixelRatio);
    document.body.appendChild(bgCanvas);

    let bgImg = new Image();
    await new Promise(res => {
        bgImg.onload = res;
        bgImg.src = '/images/bg_small.jpg';
    });

    addEventListener('scroll', function() {
        if (scrollY >= 0) {
            requestAnimationFrame(repaint);
        }
    }, {passive: true});

    function repaint() {
        if (debug) {
            debugEl.innerText = scrollY;
        }
        ctx.save();
        if (debug) {
            ctx.fillStyle = 'magenta';
        }
        else {
            ctx.fillStyle = 'rgb(40, 27, 61)';
        }
        ctx.fillRect(0, 0, innerWidth, offset);
        ctx.drawImage(bgImg, (bgImg.naturalWidth - innerWidth) / 2, scrollY * parallaxMultiplier, innerWidth, innerHeight - offset, 0, offset, innerWidth, innerHeight - offset);
        ctx.restore();
    }
    requestAnimationFrame(repaint);
})();
