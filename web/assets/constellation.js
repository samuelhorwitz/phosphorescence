export function initializeCanvas(canvasEl, component, fps) {
    let canvas = canvasEl;
    let ctx = canvas.getContext('2d');
    resizeCanvas(canvas, ctx, component);
    return beginLoop(ctx, component, fps);
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
    drawTracks(ctx, component);
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

function resizeCanvas(canvas, ctx, component, ratio = devicePixelRatio) {
    let outerSize = component.outerSize;
    canvas.width = outerSize * ratio;
    canvas.height = outerSize * ratio;
    canvas.style.width = `${outerSize}px`;
    canvas.style.height = `${outerSize}px`;
    ctx.scale(ratio, ratio);
}

function beginLoop(ctx, component, fps) {
    let keepLooping = true;
    (async () => {
        while (keepLooping) {
            requestAnimationFrame(() => {
                let outerSize = component.outerSize;
                let innerSize = component.innerSize;
                let padding = (outerSize - innerSize) / 2;
                ctx.clearRect(0, 0, outerSize, outerSize);
                ctx.save();
                ctx.translate(padding, padding);
                drawGrid(ctx, component);
                drawTracks(ctx, component);
                ctx.restore();
            });
            await new Promise(resolve => setTimeout(resolve, 1000 / fps));
        }
    })();
    return () => {
        keepLooping = false;
        ctx.clearRect(0, 0, component.outerSize, component.outerSize);
    };
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
    ctx.stroke();
    ctx.beginPath();
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

function drawTracks(ctx, component) {
    for (let i = 0; i < component.tracks.length - 1; i++) {
        let track = component.tracks[i];
        let nextTrack = component.tracks[i+1];
        drawConnector(ctx, component, component.edges[i], track.evocativeness.aetherealness, track.evocativeness.primordialness,
            nextTrack.evocativeness.aetherealness, nextTrack.evocativeness.primordialness);
    }
    let haloX, haloY;
    let hoverHaloX, hoverHaloY;
    let detailsHaloX, detailsHaloY;
    let activeId, hoverId, detailsId;
    for (let track of component.tracks) {
        drawTrack(ctx, component, track.evocativeness.aetherealness, track.evocativeness.primordialness);
        if (component.currentTrack && component.currentTrack.track.id === track.track.id) {
            activeId = track.track.id;
            haloX = track.evocativeness.aetherealness;
            haloY = track.evocativeness.primordialness;
        }
        if (typeof component.hoverTrack != 'undefined' && component.hoverTrack !== null && track.track.id === component.tracks[component.hoverTrack].track.id) {
            hoverId = track.track.id;
            hoverHaloX = track.evocativeness.aetherealness;
            hoverHaloY = track.evocativeness.primordialness;
            continue;
        }
        if (!!component.detailsTrack && track.track.id === component.detailsTrack.track.id) {
            detailsId = track.track.id;
            detailsHaloX = track.evocativeness.aetherealness;
            detailsHaloY = track.evocativeness.primordialness;
        }
    }
    drawBeatPulses(ctx, component, activeId, hoverId, detailsId, component.beatPulses);
    if (haloX && haloY) {
        drawActiveHalo(ctx, component, haloX, haloY);
    }
    if (hoverHaloX && hoverHaloY) {
        drawHoverHalo(ctx, component, hoverHaloX, hoverHaloY);
    }
    if (detailsHaloX && detailsHaloY) {
        drawHoverHalo(ctx, component, detailsHaloX, detailsHaloY);
    }
}

function drawConnector(ctx, component, edge, x1, y1, x2, y2) {
    let innerSize = component.innerSize;
    ctx.save();
    ctx.beginPath();
    ctx.moveTo(x1 * innerSize, y1 * innerSize);
    ctx.lineTo(x2 * innerSize, y2 * innerSize);
    ctx.globalCompositeOperation = 'screen';
    ctx.strokeStyle = 'lightcyan';
    ctx.lineWidth = edge;
    ctx.shadowBlur = 30;
    ctx.shadowColor = 'cyan';
    ctx.stroke();
    ctx.restore();
}

function drawTrack(ctx, component, x, y) {
    let innerSize = component.innerSize;
    ctx.save();
    ctx.beginPath();
    ctx.arc(x * innerSize, y * innerSize, 6, 0, 2 * Math.PI);
    ctx.closePath();
    ctx.fillStyle = 'white';
    ctx.strokeStyle = 'lightcyan';
    ctx.shadowColor = 'cyan';
    ctx.lineWidth = 6;
    ctx.shadowBlur = 30;
    ctx.globalCompositeOperation = 'screen';
    ctx.fill();
    ctx.stroke();
    ctx.restore();
}

function drawBeatPulses(ctx, component, activeId, hoverId, detailsId, beatPulses) {
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
        ctx.beginPath();
        ctx.arc(beatPulse.x * innerSize, beatPulse.y * innerSize, 6 * beatPulse.radiusMultiplier * beatPulse.radiusMultiplier * beatPulse.radiusMultiplier, 0, 2 * Math.PI);
        ctx.closePath();
        ctx.strokeStyle = `rgba(255, 255, 255, ${beatPulse.opacity * 0.5})`;
        ctx.stroke();
    }
    ctx.restore();
}

function drawActiveHalo(ctx, component, x, y) {
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

function drawHoverHalo(ctx, component, x, y) {
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
