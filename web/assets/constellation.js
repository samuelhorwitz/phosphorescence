export function initializeCanvas(canvasEl, component, fps) {
    let canvas = canvasEl;
    let ctx = canvas.getContext('2d');
    resizeCanvas(canvas, ctx, component);
    return beginLoop(ctx, drawGridOffscreen(component), component, fps);
}

export async function renderOffscreen(component, bgFn, fgFn) {
    let canvas = document.createElement('canvas');
    let ctx = canvas.getContext('2d');
    let outerSize = component.outerSize;
    let innerSize = component.innerSize;
    let padding = (outerSize - innerSize) / 2;
    resizeCanvas(canvas, ctx, component, 2);
    bgFn && await bgFn(ctx, outerSize);
    ctx.save();
    ctx.translate(padding, padding);
    drawGrid(ctx, component);
    drawConstellation(ctx, component, getStarCanvas());
    ctx.restore();
    fgFn && await fgFn(ctx, outerSize);
    return canvas;
}

export function saneOffscreenOptions(innerSize, gridSegments, tracks) {
    let outerSize = innerSize + 50;
    let gridInterval = innerSize / gridSegments;
    let totalGridlines = (innerSize / gridInterval) - 1;
    let centerGridlineIndex = (gridSegments / 2) - 1;
    let edges = Array(Math.max(0, tracks.length - 1)).fill(5);
    return {
        innerSize,
        outerSize,
        gridSegments,
        gridInterval,
        totalGridlines,
        centerGridlineIndex,
        tracks,
        edges
    }
}

function getStarCanvas() {
    let offscreenStarCanvas = document.createElement('canvas');
    let offscreenStarCtx = offscreenStarCanvas.getContext('2d');
    let starSize = 100;
    resizeStarCanvas(offscreenStarCanvas, offscreenStarCtx, starSize);
    drawStarOffscreen(offscreenStarCtx, starSize);
    return {canvas: offscreenStarCanvas, size: starSize};
}

function resizeCanvas(canvas, ctx, component, ratio = devicePixelRatio) {
    let outerSize = component.outerSize;
    canvas.width = outerSize * ratio;
    canvas.height = outerSize * ratio;
    canvas.style.width = `${outerSize}px`;
    canvas.style.height = `${outerSize}px`;
    ctx.scale(ratio, ratio);
}

function resizeStarCanvas(canvas, ctx, size, ratio = devicePixelRatio) {
    canvas.width = size * ratio;
    canvas.height = size * ratio;
    canvas.style.width = `${size}px`;
    canvas.style.height = `${size}px`;
    ctx.scale(ratio, ratio);
}

function beginLoop(ctx, grid, component, fps) {
    let keepLooping = true;
    let redrawConstellation = true;
    let constellation;
    (async () => {
        while (keepLooping) {
            requestAnimationFrame(() => {
                if (redrawConstellation) {
                    constellation = drawConstellationOffscreen(component);
                    redrawConstellation = false;
                }
                let outerSize = component.outerSize;
                let innerSize = component.innerSize;
                let padding = (outerSize - innerSize) / 2;
                ctx.clearRect(0, 0, outerSize, outerSize);
                ctx.save();
                ctx.drawImage(grid, 0, 0, component.outerSize, component.outerSize);
                ctx.globalCompositeOperation = 'screen';
                ctx.drawImage(constellation, 0, 0, component.outerSize, component.outerSize);
                ctx.translate(padding, padding);
                drawInteractive(ctx, component);
                ctx.restore();
            });
            await new Promise(resolve => setTimeout(resolve, 1000 / fps));
        }
    })();
    return {
        destroyer() {
            keepLooping = false;
            ctx.clearRect(0, 0, component.outerSize, component.outerSize);
        },
        redrawConstellation() {
            redrawConstellation = true;
        }
    };
}

function drawStarOffscreen(ctx, size) {
    ctx.moveTo(size / 2, size / 2);
    ctx.arc(size / 2, size / 2, 6, 0, 2 * Math.PI);
    ctx.fillStyle = 'white';
    ctx.strokeStyle = 'lightcyan';
    ctx.shadowColor = 'cyan';
    ctx.lineWidth = 6;
    ctx.shadowBlur = 30;
    ctx.globalCompositeOperation = 'screen';
    ctx.fill();
    ctx.stroke();
}

function drawGridOffscreen(component) {
    let canvas = document.createElement('canvas');
    let ctx = canvas.getContext('2d');
    resizeCanvas(canvas, ctx, component);
    let outerSize = component.outerSize;
    let innerSize = component.innerSize;
    let padding = (outerSize - innerSize) / 2;
    ctx.translate(padding, padding);
    drawGrid(ctx, component);
    return canvas;
}

function drawConstellationOffscreen(component) {
    let canvas = document.createElement('canvas');
    let ctx = canvas.getContext('2d');
    resizeCanvas(canvas, ctx, component);
    let outerSize = component.outerSize;
    let innerSize = component.innerSize;
    let padding = (outerSize - innerSize) / 2;
    ctx.translate(padding, padding);
    drawConstellation(ctx, component, getStarCanvas());
    return canvas;
}

function drawGrid(ctx, component) {
    let innerSize = component.innerSize;
    let totalGridlines = component.totalGridlines;
    let centerGridlineIndex = component.centerGridlineIndex;
    let gridInterval = component.gridInterval;
    let gridSegments = component.gridSegments;
    ctx.save();
    ctx.lineWidth = 4;
    ctx.strokeStyle = 'white';
    ctx.shadowBlur = 20;
    ctx.shadowColor = 'magenta';
    ctx.beginPath();
    ctx.moveTo(0, 0);
    ctx.lineTo(innerSize, 0);
    ctx.lineTo(innerSize, innerSize);
    ctx.lineTo(0, innerSize);
    ctx.closePath();
    for (let i = 0; i < totalGridlines; i++) {
        if (i === centerGridlineIndex) {
            // we draw the middle ones separately
            continue;
        }
        ctx.moveTo(0, (i + 1) * gridInterval);
        ctx.lineTo(innerSize, (i + 1) * gridInterval);
        ctx.moveTo((i + 1) * gridInterval, 0);
        ctx.lineTo((i + 1) * gridInterval, innerSize);
    }
    ctx.closePath();
    ctx.stroke();
    ctx.lineWidth = 7;
    ctx.strokeStyle = 'white';
    ctx.shadowBlur = 20;
    ctx.shadowColor = 'cyan';
    ctx.beginPath();
    ctx.moveTo(0, (gridSegments / 2) * gridInterval);
    ctx.lineTo(innerSize, (gridSegments / 2) * gridInterval);
    ctx.moveTo((gridSegments / 2) * gridInterval, 0);
    ctx.lineTo((gridSegments / 2) * gridInterval, innerSize);
    ctx.closePath();
    ctx.stroke();
    ctx.restore();
}

function drawConstellation(ctx, component, star) {
    drawConnectors(ctx, cb => {
        for (let i = 0; i < component.tracks.length - 1; i++) {
            let track = component.tracks[i];
            let nextTrack = component.tracks[i+1];
            drawConnector(cb, component.innerSize, component.edges[i], track.evocativeness.aetherealness, track.evocativeness.primordialness, nextTrack.evocativeness.aetherealness, nextTrack.evocativeness.primordialness);
        }
    });
    drawTracks(ctx, ctx => {
        for (let track of component.tracks) {
            drawTrack(ctx, star, component.innerSize, track.evocativeness.aetherealness, track.evocativeness.primordialness);
        }
    });
}

function drawInteractive(ctx, component) {
    drawBeatPulses(ctx, component);
    drawActiveHalo(ctx, component);
    drawHoverHalo(ctx, component);
    drawDetailsHalo(ctx, component);
}

function drawConnectors(ctx, draw) {
    ctx.save();
    ctx.beginPath();
    let bucket = {};
    draw((edge, fn) => {
        if (!bucket[edge]) {
            bucket[edge] = [];
        }
        bucket[edge].push(fn);
    });
    ctx.globalCompositeOperation = 'screen';
    ctx.strokeStyle = 'lightcyan';
    ctx.shadowBlur = 30;
    ctx.shadowColor = 'cyan';
    for (let [edge, fns] of Object.entries(bucket)) {
        ctx.lineWidth = parseInt(edge, 10);
        for (let fn of fns) {
            fn(ctx);
        }
        ctx.stroke();
    }
    ctx.restore();
}

function drawConnector(cb, innerSize, edge, x1, y1, x2, y2) {
    cb(edge, ctx => {
        ctx.moveTo(x1 * innerSize, y1 * innerSize);
        ctx.lineTo(x2 * innerSize, y2 * innerSize);
    });
}

function drawTracks(ctx, draw) {
    ctx.save();
    ctx.globalCompositeOperation = 'screen';
    draw(ctx);
    ctx.restore();
}

function drawTrack(ctx, {canvas, size}, innerSize, x, y) {
    ctx.drawImage(canvas, (x * innerSize) - (size / 2), (y * innerSize) - (size / 2), size, size);
}

function drawBeatPulses(ctx, component) {
    let activeId, hoverId, detailsId;
    if (component.currentTrack) {
        activeId = component.currentTrack.id;
    }
    if (typeof component.hoverTrack != 'undefined' && component.hoverTrack !== null) {
        hoverId = component.tracks[component.hoverTrack].id;
    }
    if (component.detailsTrack) {
        detailsId = component.detailsTrack.id;
    }
    let beatPulses = component.beatPulses;
    if (!beatPulses) {
        return;
    }
    let innerSize = component.innerSize;
    ctx.save();
    ctx.globalCompositeOperation = 'screen';
    ctx.shadowColor = 'white';
    ctx.shadowBlur = 30;
    ctx.lineWidth = 5;
    for (let beatPulse of beatPulses) {
        if (beatPulse.trackId !== activeId && beatPulse.trackId !== hoverId && beatPulse.trackId !== detailsId) {
            continue;
        }
        let percentDone = (performance.now() - beatPulse.start) / beatPulse.ttl;
        if (percentDone > 1) {
            continue;
        }
        let opacity = 1 - percentDone;
        let radiusMultiplier = 1 + percentDone;
        ctx.beginPath();
        ctx.arc(beatPulse.x * innerSize, beatPulse.y * innerSize, 6 * radiusMultiplier * radiusMultiplier * radiusMultiplier, 0, 2 * Math.PI);
        ctx.closePath();
        ctx.strokeStyle = `rgba(255, 255, 255, ${opacity * 0.5})`;
        ctx.stroke();
    }
    ctx.restore();
}

function drawActiveHalo(ctx, component) {
    if (!component.currentTrack) {
        return;
    }
    let x = component.currentTrack.evocativeness.aetherealness;
    let y = component.currentTrack.evocativeness.primordialness;
    let innerSize = component.innerSize;
    ctx.save();
    ctx.translate(x * innerSize, y * innerSize);
    ctx.rotate((component.haloRotation * 0.5) * Math.PI / 180);
    ctx.translate(-x * innerSize, -y * innerSize);
    ctx.beginPath();
    ctx.arc(x * innerSize, y * innerSize, 22, 0, 2 * Math.PI);
    ctx.closePath();
    ctx.globalCompositeOperation = 'source-over';
    ctx.strokeStyle = 'magenta';
    ctx.lineWidth = 5;
    ctx.setLineDash([17, 5]);
    ctx.stroke();
    ctx.restore();
}

function drawHoverHalo(ctx, component) {
    if (typeof component.hoverTrack == 'undefined' || component.hoverTrack == null) {
        return;
    }
    let hoverTrack = component.tracks[component.hoverTrack];
    let x = hoverTrack.evocativeness.aetherealness;
    let y = hoverTrack.evocativeness.primordialness;
    drawHoverDetailsHalo(ctx, component, x, y);
}

function drawDetailsHalo(ctx, component) {
    if (!component.detailsTrack) {
        return;
    }
    let x = component.detailsTrack.evocativeness.aetherealness;
    let y = component.detailsTrack.evocativeness.primordialness;
    drawHoverDetailsHalo(ctx, component, x, y);
}

function drawHoverDetailsHalo(ctx, component, x, y) {
    let innerSize = component.innerSize;
    ctx.save();
    ctx.translate(x * innerSize, y * innerSize);
    ctx.rotate(-(component.haloRotation * 0.65) * Math.PI / 180);
    ctx.translate(-x * innerSize, -y * innerSize);
    ctx.beginPath();
    ctx.arc(x * innerSize, y * innerSize, 15, 0, 2 * Math.PI);
    ctx.closePath();
    ctx.globalCompositeOperation = 'source-over';
    ctx.strokeStyle = 'cyan';
    ctx.lineWidth = 5;
    ctx.setLineDash([7, 2]);
    ctx.stroke();
    ctx.restore();
}
