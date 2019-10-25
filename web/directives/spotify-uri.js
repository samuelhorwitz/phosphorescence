import Vue from 'vue';

const canvas = document.createElement('canvas');
const ctx = canvas.getContext('2d');
const canvasWidth = 200;
const logoSize = 20;
const padding = 5;
const spotifyLogo = new Image(logoSize * devicePixelRatio, logoSize * devicePixelRatio);
spotifyLogo.src = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADwAAAA8CAMAAAANIilAAAAAUVBMVEUAAAD///////////////////////////////////////////////////////////////////////////////////////////////////////8IN+deAAAAGnRSTlMA8RYLUTpJMimZg2fIv45y3tNCtaRbIOV7rSsbnPEAAAJnSURBVEjHlZZZFqMgEEUpZEYExxj2v9DuTpsKEDFy85Ej+rBODQ9JDcpXb4bB+JVT0gJleoSIwKjZ3Q2mBeIXsHDyG9bHCr34IZXu9ZqK3Mkr7RB/YKrSbYw/6WklURBvAPw0U/Em4krbrp5iA7zIFbSIIc/aGJvoL+vbUG+ZB1W5ylAodrXn6lu4sypBvxjLuFTbtinFxer1PMJFvXocOyu7c2Pwz0f56rzEA7mEm/Gk2Et8MRKky34f5AAR0a+YjoXhLcsoFjsf4sGjS9Ll8AEqBbPWW7sKriguH/8ekpTpeOD/BWZ1/4CsSOM8rIokKAw17cz+GWKFsKwUIzDYoxTiPeDJDrnF8eA3ZO8GC/4lno91SdayCYPTu2VimgSz2F4oXynFJDHiY8I4sI0U0MnP8NFDRCwxiXatnz8LoBDZ01F+kguo7WOEbAOTijU2At3UX7bCo0VfOEIaNihCOu61C3AsBLd4kWzBHlnYWcJA4+RlI244Bp/et1iqa8advrorczv23SRwbmKgt2N+EZm1J8oAwjiGAPl2hhKbuXdp2Q+nrZC0+z+bG1+N+3QIhNK8dfywqDMvmIZwZqBD7p2B1GAzyjP/7KAUUyVW6723jG9oInKOOY+uSKBWzOQ2Heb9XWSB6x8D/D3RYVCv17OYItH0M+D7cuZo8IVfiniDRcmsqNPnoGtmTo7YZhRBTKt2Jwl9m9blJtP6QZPBW8SSFIj72ol8IZq0JRxunTyyYs03cu4oqWGa6lui3GVPKnKNcLWPODeR33B9krmHluQenRh6SI+MQXSkBSqZ3Y3ZLZPVBP8BtYLcdvrDozQAAAAASUVORK5CYII=';
resizeCanvas(canvas, ctx, canvasWidth + (padding * 2), logoSize + (padding * 2));
document.body.appendChild(canvas);
canvas.style.position = 'absolute';
canvas.style.left = `-${canvasWidth + (padding * 2)}px`;
canvas.style.top = `-${logoSize + (padding * 2)}px`;
canvas.setAttribute('aria-hidden', 'true');

function resizeCanvas(canvas, ctx, width, height) {
    canvas.width = width * devicePixelRatio;
    canvas.height = height * devicePixelRatio;
    canvas.style.width = `${width}px`;
    canvas.style.height = `${height}px`;
    ctx.scale(devicePixelRatio, devicePixelRatio);
}

function getCanvas(title) {
    let width = canvas.width / devicePixelRatio;
    let height = canvas.height / devicePixelRatio;
    ctx.clearRect(0, 0, width, height);
    ctx.save();
    if (title) {
        let titleLines = title.split('\n');
        let textLeft = (spotifyLogo.width / devicePixelRatio) + 5;
        ctx.fillStyle = 'white';
        ctx.textBaseline = 'top';
        if (titleLines.length == 1) {
            ctx.font = '12pt sans-serif';
            let line = titleLines[0];
            let textWidth;
            while ((textWidth = ctx.measureText(line).width) > canvasWidth - textLeft) {
                if (line[line.length - 1] === '…') {
                    line = line.substring(0, line.length - 2) + '…';
                } else {
                    line = line.substring(0, line.length - 1) + '…';
                }
            }
            drawPill(ctx, textLeft + textWidth);
            ctx.translate(padding, padding);
            ctx.drawImage(spotifyLogo, 0, 0, logoSize, logoSize);
            ctx.fillText(line, textLeft, 2);
        } else if (titleLines.length == 2) {
            ctx.font = '7pt sans-serif';
            let offset = 0;
            let realWidth = 0;
            let lines = [];
            for (let line of titleLines) {
                let textWidth;
                while ((textWidth = ctx.measureText(line).width) > canvasWidth - textLeft) {
                    if (line[line.length - 1] === '…') {
                        line = line.substring(0, line.length - 2) + '…';
                    } else {
                        line = line.substring(0, line.length - 1) + '…';
                    }
                }
                realWidth = Math.max(realWidth, textWidth);
                lines.push({line, textLeft, offset});
                offset += 10;
            }
            drawPill(ctx, textLeft + realWidth);
            ctx.translate(padding, padding);
            ctx.drawImage(spotifyLogo, 0, 0, logoSize, logoSize);
            for (let {line, textLeft, offset} of lines) {
                ctx.fillText(line, textLeft, offset);
            }
        }
    }
    ctx.restore();
    return canvas;
}

function drawPill(ctx, width) {
    width += padding * 2;
    let radius = 5;
    let x = 0;
    let y = 0;
    let height = logoSize + (padding * 2);
    ctx.save();
    ctx.fillStyle = 'rgba(0, 0, 0, 0.8)';
    ctx.beginPath();
    ctx.moveTo(x + radius, y);
    ctx.lineTo(x + width - radius, y);
    ctx.quadraticCurveTo(x + width, y, x + width, y + radius);
    ctx.lineTo(x + width, y + height - radius);
    ctx.quadraticCurveTo(x + width, y + height, x + width - radius, y + height);
    ctx.lineTo(x + radius, y + height);
    ctx.quadraticCurveTo(x, y + height, x, y + height - radius);
    ctx.lineTo(x, y + radius);
    ctx.quadraticCurveTo(x, y, x + radius, y);
    ctx.closePath();
    ctx.fill();
    ctx.restore();
}

function handleDragStart(e) {
    if (!handle(e.target, e.dataTransfer)) {
        return;
    }
    e.dataTransfer.clearData('text/html');
    e.dataTransfer.effectAllowed = 'all';
    e.dataTransfer.setDragImage(getCanvas(e.target.getAttribute('data-spotify-uri-title')), 0, 0);
    e.stopPropagation();
}

function handleCopy(e) {
    if (!handle(document.activeElement, e.clipboardData)) {
        return;
    }
    e.preventDefault();
}

function handle(target, setDataInterface) {
    let ids = target.getAttribute('data-spotify-uri-ids').split(',');
    let type = target.getAttribute('data-spotify-uri-type');
    let uri;
    if (type === 'track') {
        uri = buildUriList('https://open.spotify.com/track/', ids);
        setDataInterface.setData('text/x-spotify-tracks', uri);
    }
    else if (type === 'artist') {
        uri = buildUriList('https://open.spotify.com/artist/', ids);
        setDataInterface.setData('text/x-spotify-artists', uri);
    }
    else if (type === 'album') {
        uri = buildUriList('https://open.spotify.com/album/', ids);
        setDataInterface.setData('text/x-spotify-albums', uri);
    }
    else {
        return false;
    }
    setDataInterface.setData('text/x-phosphor-origin', 'true');
    setDataInterface.setData('text/uri-list', uri);
    setDataInterface.setData('text/plain', uri);
    return true;
}

function buildUriList(prefix, ids) {
    let uris = [];
    for (let id of ids) {
        uris.push(`${prefix}${id}`);
    }
    return uris.join('\n');
}

Vue.directive('spotify-uri', {
    bind(el, {value, arg}) {
        if (!Array.isArray(value)) {
            value = [value];
        }
        el.setAttribute('data-spotify-uri-type', arg);
        el.setAttribute('data-spotify-uri-ids', value.join(','));
        el.setAttribute('draggable', 'true');
        el.addEventListener('dragstart', handleDragStart);
    },
    unbind(el) {
        el.setAttribute('draggable', 'false');
        el.removeAttribute('data-spotify-uri-type');
        el.removeAttribute('data-spotify-uri-ids');
        el.removeEventListener('dragstart', handleDragStart);
    },
    update(el, {value, arg}) {
        if (!Array.isArray(value)) {
            value = [value];
        }
        el.setAttribute('data-spotify-uri-type', arg);
        el.setAttribute('data-spotify-uri-ids', value.join(','));
    }
});

Vue.directive('spotify-uri-title', {
    bind(el, {value}) {
        el.setAttribute('data-spotify-uri-title', value);
    },
    unbind(el) {
        el.removeAttribute('data-spotify-uri-title');
    },
    update(el, {value}) {
        el.setAttribute('data-spotify-uri-title', value);
    }
});

document.body.addEventListener('copy', handleCopy);
