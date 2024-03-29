(async function() {
    const stopped = 0;
    const paused = 1;
    const playing = 2;
    const apiOrigin = 'https://api.phosphor.me';
    let blockClick = false;
    let currentTouch = 0;
    let tracks = [];
    let trackTitles = [];
    let trackTitlesWithLinks = [];
    let trackImages = [];
    let trackAlbumUrls = [];
    let previewUrls = {};
    let progressEl, cursorEl, audioEl, playPlaylistEl, pausePlaylistEl, subtitleEl, playlistImageEl;
    let state = {
        playState: stopped,
        currentTrack: 0,
        percent: 0,
        tracksPlayed: 0
    };

    function getCaptchaToken(action) {
        return grecaptcha.execute('6LdfBboUAAAAAFv0977A1dWeer-eTy0IBmynzHcS', {action});
    }

    function showCover(coverName) {
        let nav = document.getElementById('nav');
        let cover = document.getElementById(coverName);
        if (nav && cover) {
            nav.classList.add('cover');
            cover.hidden = false;
        }
    }

    function hideCover(coverName) {
        let nav = document.getElementById('nav');
        let cover = document.getElementById(coverName);
        if (nav && cover) {
            nav.classList.remove('cover');
            cover.hidden = true;
        }
    }

    function msToPrettyTime(ms) {
        let s = ms / 1000;
        let m = Math.floor(s / 60);
        let rs = `${Math.round(s % 60)}`.padStart(2, '0');
        return `${m}:${rs}`;
    }

    function togglePlayback() {
        if (state.playState === playing) {
            gtag('event', 'click-pause', {
                event_category: 'twitter-player'
            });
            pause();
        } else if (state.playState === stopped) {
            gtag('event', 'click-play', {
                event_category: 'twitter-player',
                event_label: 'stopped'
            });
            incTracksPlayed();
            play();
        } else if (state.playState === paused) {
            gtag('event', 'click-play', {
                event_category: 'twitter-player',
                event_label: 'paused'
            });
            resume();
        }
    }

    function play() {
        if (!previewUrls[tracks[state.currentTrack]]) {
            return;
        }
        state.playState = playing;
        audioEl.src = previewUrls[tracks[state.currentTrack]];
        audioEl.play();
        pausePlaylistEl.hidden = false;
        playPlaylistEl.hidden = true;
        clearPlaying();
        document.title = `phosphor.me - ${trackTitles[state.currentTrack]}`;
        audioEl.title = `phosphor.me - ${trackTitles[state.currentTrack]}`;
        subtitleEl.innerHTML = trackTitlesWithLinks[state.currentTrack];
        playlistImageEl.setAttribute('data-old-background-image', playlistImageEl.style.backgroundImage);
        playlistImageEl.style.backgroundImage = `url('${trackImages[state.currentTrack]}')`;
        playlistImageLinkEl.setAttribute('data-old-href', playlistImageLinkEl.href);
        playlistImageLinkEl.href = trackAlbumUrls[state.currentTrack];
        let shouldPlay = document.querySelector(`tr[data-track-index="${state.currentTrack}"]`);
        if (shouldPlay) {
            shouldPlay.classList.add('isPlaying');
        }
    }

    function resume() {
        state.playState = playing;
        audioEl.play();
        pausePlaylistEl.hidden = false;
        playPlaylistEl.hidden = true;
    }

    function pause() {
        state.playState = paused;
        audioEl.pause();
        pausePlaylistEl.hidden = true;
        playPlaylistEl.hidden = false;
    }

    function previous() {
        if (state.currentTrack === 0) {
            return;
        }
        incTracksPlayed();
        state.currentTrack--;
        gtag('event', 'click-previous', {
            event_category: 'twitter-player',
            event_label: 'track',
            value: state.currentTrack
        });
        play();
    }

    function next() {
        if (state.currentTrack === tracks.length - 1) {
            return;
        }
        incTracksPlayed();
        state.currentTrack++;
        gtag('event', 'click-next', {
            event_category: 'twitter-player',
            event_label: 'track',
            value: state.currentTrack
        });
        play();
    }

    function handleTrackClick({target}) {
        while (target && target.nodeName !== 'TR') {
            target = target.parentNode;
        }
        if (!target) {
            return;
        }
        incTracksPlayed();
        let i = parseInt(target.getAttribute('data-track-index'), 10);
        gtag('event', 'click-track', {
            event_category: 'twitter-player',
            event_label: 'track',
            value: i
        });
        if (i < 0 || i >= tracks.length) {
            return;
        }
        state.currentTrack = i;
        play();
    }

    function handleScrubberClick({offsetX}) {
        if (state.playState === stopped) {
            return;
        }
        let frac = offsetX / progressEl.offsetWidth;
        state.percent = frac * 100;
        gtag('event', 'click-scrubber', {
            event_category: 'twitter-player',
            event_label: 'percent',
            value: Math.round(state.percent)
        });
        redrawScrubber();
        audioEl.currentTime = frac * audioEl.duration;
    }

    function handlePlaybackTic({target}) {
        state.percent = !target.currentTime || !target.duration ? 0 : (target.currentTime / target.duration) * 100;
        redrawScrubber();
    }

    function handlePlaybackEnd() {
        gtag('event', 'preview-completion', {
            event_category: 'twitter-player',
            event_label: 'track',
            value: state.currentTrack
        });
        state.percent = 0;
        redrawScrubber();
        if (state.currentTrack === tracks.length - 1) {
            state.currentTrack = 0;
            state.playState = stopped;
            pausePlaylistEl.hidden = true;
            playPlaylistEl.hidden = false;
            document.title = `phosphor.me - Web Player`;
            audioEl.removeAttribute('title');
            subtitleEl.innerHTML = 'Created with <a href="https://phosphor.me" target="_blank">phosphor.me</a>.';
            playlistImageEl.style.backgroundImage = playlistImageEl.getAttribute('data-old-background-image');
            playlistImageEl.removeAttribute('data-old-background-image');
            playlistImageLinkEl.href = playlistImageLinkEl.getAttribute('data-old-href');
            playlistImageLinkEl.removeAttribute('data-old-href');
            clearPlaying();
        } else {
            next();
        }
    }

    function clearPlaying() {
        let isPlaying = document.querySelectorAll('tr.isPlaying');
        for (let el of isPlaying) {
            el.classList.remove('isPlaying');
        }
    }

    function redrawScrubber() {
        let percentString = `${state.percent}%`;
        progressEl.value = state.percent;
        progressEl.innerText = percentString;
        cursorEl.style.left = percentString;
    }

    function incTracksPlayed() {
        state.tracksPlayed++;
        if (state.tracksPlayed > 0 && state.tracksPlayed % 3 === 0) {
            showCover('playFullCover');
        }
    }

    function updateLinks(playlistId) {
        let spotifyUrl = `https://open.spotify.com/playlist/${playlistId}?utm_campaign=me.phosphor`;
        let phosphorescenceUrl = `https://phosphor.me/${playlistId}`;
        let spotifyLinks = document.querySelectorAll('[data-link-spotify]');
        spotifyLinks.forEach(el => {
            el.href = spotifyUrl;
            el.target = '_blank';
        });
        let phosphorescenceLinks = document.querySelectorAll('[data-link-phosphor]');
        phosphorescenceLinks.forEach(el => {
            el.href = phosphorescenceUrl;
            el.target = '_blank';
        });
    }

    async function init() {
        let playlistId = location.search.slice(1);
        updateLinks(playlistId);
        let playlistResponse = await fetch(`${apiOrigin}/player/playlist/${playlistId}?captcha=${await getCaptchaToken('api/player/playlist')}`);
        let {playlist} = await playlistResponse.json();
        if (playlist.owner.id !== 'spv1hpk4dj1qwneuxwg4yg0cn') {
            throw new Error('This is not a Phosphorescence playlist!');
        }
        for (let {id, preview_url} of playlist.tracks) {
            previewUrls[id] = preview_url;
        }
        let playlistName = playlist.name.replace('phosphor.me | ', '');
        let playlistTitleEl = document.getElementById('playlistTitle');
        playlistTitleEl.innerText = playlistName;
        playlistTitleEl.setAttribute('data-text', playlistName);
        document.body.addEventListener('touchstart', () => {
            playlistTitleEl.classList.add('touch');
            currentTouch = Date.now();
            blockClick = true;
        });
        document.body.addEventListener('touchend', async () => {
            let oldCurrentTouch = currentTouch;
            await new Promise(resolve => setTimeout(resolve, 1000));
            if (oldCurrentTouch === currentTouch) {
                playlistTitleEl.classList.remove('touch');
            }
        });
        document.body.addEventListener('touchcancel', async () => {
            let oldCurrentTouch = currentTouch;
            await new Promise(resolve => setTimeout(resolve, 1000));
            if (oldCurrentTouch === currentTouch) {
                playlistTitleEl.classList.remove('touch');
            }
        });
        document.body.addEventListener('click', async () => {
            if (blockClick) {
                return;
            }
            currentTouch = Date.now();
            playlistTitleEl.classList.add('touch');
            let oldCurrentTouch = currentTouch;
            await new Promise(resolve => setTimeout(resolve, 1000));
            if (oldCurrentTouch === currentTouch) {
                playlistTitleEl.classList.remove('touch');
            }
        });
        playlistImageLinkEl = document.getElementById('playlistImageLink');
        playlistImageEl = document.getElementById('playlistImage');
        if (playlist.images && playlist.images.length && playlist.images[0]) {
            playlistImageEl.style.backgroundImage = `url("${playlist.images[0].url}")`;
        }
        let playlistTableEl = document.getElementById('playlistTable');
        let trackTemplate = document.getElementById('trackRow');
        for (let [i, track] of Object.entries(playlist.tracks)) {
            let trackEl = document.importNode(trackTemplate.content, true);
            let tr = trackEl.querySelector('tr');
            tr.setAttribute('data-track-index', i);
            tr.addEventListener('click', handleTrackClick);
            let td = trackEl.querySelectorAll('td');
            let trackIndex = td[0].querySelector('.trackIndex');
            trackIndex.innerText = parseInt(i, 10) + 1;
            let trackData = td[1].querySelectorAll('div');
            trackData[0].innerText = track.name;
            trackData[1].innerText = track.artists.map(a => a.name).join(', ');
            td[2].innerText = msToPrettyTime(track.duration_ms);
            playlistTableEl.appendChild(trackEl);
        }
        playPlaylistEl = document.getElementById('playPlaylist');
        playPlaylistEl.addEventListener('click', togglePlayback);
        pausePlaylistEl = document.getElementById('pausePlaylist');
        pausePlaylistEl.addEventListener('click', togglePlayback);
        let previousEl = document.getElementById('previous');
        let nextEl = document.getElementById('next');
        previousEl.addEventListener('click', previous);
        nextEl.addEventListener('click', next);
        audioEl = document.getElementById('audio');
        audioEl.addEventListener('timeupdate', handlePlaybackTic);
        audioEl.addEventListener('ended', handlePlaybackEnd);
        let scrubberContainerEl = document.getElementById('scrubberContainer');
        progressEl = document.getElementById('progress');
        cursorEl = document.getElementById('cursor');
        scrubberContainerEl.addEventListener('click', handleScrubberClick);
        let coverCloseEl = document.getElementById('coverClose');
        coverCloseEl.addEventListener('click', () => hideCover('playFullCover'));
        let infoButtonEl = document.getElementById('infoButton');
        let legalFooterEl = document.getElementById('legalFooter');
        infoButtonEl.addEventListener('click', () => legalFooterEl.classList.toggle('touched'));
        subtitleEl = document.getElementById('subtitle');
        tracks = playlist.tracks.map(t => t.id);
        trackTitles = playlist.tracks.map(t => `${t.name} - ${t.artists.map(a => a.name).join(', ')}`);
        trackTitlesWithLinks = playlist.tracks.map(t => `
            <a href="https://open.spotify.com/track/${t.id}?utm_campaign=me.phosphor" target="_blank" rel="external noopener">${t.name}</a>
            -
            ${t.artists.map(a => `<a href="https://open.spotify.com/artist/${a.id}?utm_campaign=me.phosphor" target="_blank" rel="external noopener">${a.name}</a>`).join(', ')}
        `);
        trackImages = playlist.tracks.map(t => {
            if (t.album.images.length) {
                return t.album.images[0].url;
            }
            return null;
        });
        trackAlbumUrls = playlist.tracks.map(t => `https://open.spotify.com/album/${t.album.id}?utm_campaign=me.phosphor`);
    }

    await new Promise(resolve => grecaptcha.ready(resolve));
    let timeoutEl = document.getElementById('timeoutCover');
    try {
        await init();
    }
    catch (e) {
        if (timeoutEl) {
            timeoutEl.hidden = true;
        }
        showCover('errorCover');
        console.error('Could not initialize Phosphorescence player', e);
        gtag('event', 'exception', {
            'description': e,
            'fatal': true
        });
    }
    if (timeoutEl) {
        timeoutEl.hidden = true;
    }
})();