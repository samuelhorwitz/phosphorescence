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
        debugEl.style.left = '40px';
        document.body.appendChild(debugEl);
    }
    let bgCanvas = document.createElement('canvas');
    bgCanvas.classList.add('bgIosCanvas');
    let ctx = bgCanvas.getContext('2d');
    function resetCanvasSize() {
        setTimeout(() => {
            let width = innerWidth;
            let height = innerHeight;
            if (navigator.standalone) {
                width = outerWidth;
                height = outerHeight;
            }
            bgCanvas.width = width * devicePixelRatio;
            bgCanvas.height = height * devicePixelRatio;
            bgCanvas.style.width = `${width}px`;
            bgCanvas.style.height = `${height}px`;
            ctx.scale(devicePixelRatio, devicePixelRatio);
        }, 200);
    }
    resetCanvasSize();
    document.body.appendChild(bgCanvas);

    document.body.addEventListener('resize', function() {
        resetCanvasSize();
        requestAnimationFrame(repaint);
    });

    window.addEventListener('orientationchange', function() {
        resetCanvasSize();
        requestAnimationFrame(repaint);
    })

    let bgImg = new Image();
    await new Promise(res => {
        bgImg.onload = res;
        bgImg.src = '/images/bg_small.jpg';
    });

    let repaintLocked = false;

    addEventListener('scroll', function() {
        if (scrollY >= 0) {
            requestAnimationFrame(repaint);
        }
    }, {passive: true});

    addEventListener('touchend', function() {
        if (scrollY >= 0) {
            return;
        }
        repaintLocked = true;
        setTimeout(() => repaintLocked = false, 100);
    }, {passive: true});

    function repaint() {
        if (repaintLocked) {
            return;
        }
        let width = innerWidth;
        let height = innerHeight;
        if (navigator.standalone) {
            width = outerWidth;
            height = outerHeight;
        }
        if (debug) {
            debugEl.innerText = `${width}, ${height}, ${scrollY}`;
        }
        ctx.save();
        if (debug) {
            ctx.fillStyle = 'magenta';
        }
        else {
            ctx.fillStyle = 'rgb(40, 27, 61)';
        }
        if (offset > 0) {
            ctx.fillRect(0, 0, width, offset);
        }
        let neededPixels = height - offset;
        let yOffset = scrollY * parallaxMultiplier;
        let pixelsLeft = bgImg.naturalHeight - yOffset;
        if (yOffset + neededPixels > bgImg.naturalHeight) {
            yOffset = bgImg.naturalHeight - neededPixels;
        }
        ctx.drawImage(bgImg,
            // source
            (bgImg.naturalWidth - width) / 2,
            yOffset,
            width,
            neededPixels,

            // dest
            0,
            offset,
            width,
            neededPixels
        );
        ctx.restore();
    }
    requestAnimationFrame(repaint);
    document.body.classList.add('iosRubberband');
    document.body.classList.add('ios');
})();
