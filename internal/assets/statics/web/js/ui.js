import { state } from './state.js';
import { getTime, formatTime } from './utils.js';
import * as el from './elements.js';

export function addMessage(type, content) {
    const message = document.createElement('div');
    message.className = `message ${type}`;

    const time = document.createElement('span');
    time.className = 'message-time';
    time.textContent = getTime();

    const typeLabel = document.createElement('span');
    typeLabel.className = 'message-type';
    typeLabel.textContent = type.toUpperCase();

    const contentSpan = document.createElement('span');
    contentSpan.className = 'message-content';
    contentSpan.textContent = typeof content === 'object' ? JSON.stringify(content, null, 2) : content;

    message.appendChild(time);
    message.appendChild(typeLabel);
    message.appendChild(contentSpan);

    el.messagesDiv.appendChild(message);
    el.messagesDiv.scrollTop = el.messagesDiv.scrollHeight;
}

export function clearLog() {
    el.messagesDiv.innerHTML = '';
    addMessage('info', 'Log cleared');
}

export function updateStats() {
    el.sentCountEl.textContent = state.sentCount;
    el.receivedCountEl.textContent = state.receivedCount;
    el.errorCountEl.textContent = state.errorCount;
}

export function updateStatus(connected) {
    if (connected) {
        el.statusDiv.className = 'status connected';
        el.statusDiv.innerHTML = '<span class="status-dot"></span><span>Connected</span>';
        el.connectBtn.disabled = true;
        el.disconnectBtn.disabled = false;
        el.pingBtn.disabled = false;
        el.sendBtn.disabled = false;
        el.messageInput.disabled = false;
        el.spotifyArtistInput.disabled = false;
        el.spotifyArtistBtn.disabled = false;
        el.spotifyFollowBtn.disabled = false;

        // Enable volume controls
        el.volDownBtn.disabled = false;
        el.volUpBtn.disabled = false;
        el.volMuteBtn.disabled = false;
        el.volSlider.disabled = false;

        // Enable brightness controls
        el.brightDownBtn.disabled = false;
        el.brightUpBtn.disabled = false;
        el.brightSlider.disabled = false;

        el.prevBtn.disabled = false;
        el.playBtn.disabled = false;
        el.pauseBtn.disabled = false;
        el.nextBtn.disabled = false;
        el.mediaPlayer.classList.remove('hidden');
    } else {
        el.statusDiv.className = 'status disconnected';
        el.statusDiv.innerHTML = '<span class="status-dot"></span><span>Disconnected</span>';
        el.connectBtn.disabled = false;
        el.disconnectBtn.disabled = true;
        el.pingBtn.disabled = true;
        el.sendBtn.disabled = true;
        el.messageInput.disabled = true;
        el.spotifyArtistInput.disabled = true;
        el.spotifyArtistBtn.disabled = true;
        el.spotifyFollowBtn.disabled = true;

        // Disable volume controls
        el.volDownBtn.disabled = true;
        el.volUpBtn.disabled = true;
        el.volMuteBtn.disabled = true;
        el.volSlider.disabled = true;

        // Disable brightness controls
        el.brightDownBtn.disabled = true;
        el.brightUpBtn.disabled = true;
        el.brightSlider.disabled = true;

        el.prevBtn.disabled = true;
        el.playBtn.disabled = true;
        el.pauseBtn.disabled = true;
        el.nextBtn.disabled = true;
        el.mediaPlayer.classList.add('hidden');
    }
}

export function updateMediaInfo(data) {
    // Check if media info exists
    if (!data || Object.keys(data).length === 0) {
        el.noMedia.classList.remove('hidden');
        el.mediaContent.classList.add('hidden');
        return;
    }

    el.noMedia.classList.add('hidden');
    el.mediaContent.classList.remove('hidden');

    // Update track info (handle both lowercase and Title case)
    el.trackTitle.textContent = data.Title || data.title || 'Unknown Track';
    el.trackArtist.textContent = data.Artist || data.artist || 'Unknown Artist';
    el.trackAlbum.textContent = data.Album || data.album || 'Unknown Album';

    // Update album art (handle Artwork field)
    const artworkUrl = data.Artwork || data.artwork || data.artUrl || data.album_art || data.albumArt;
    if (artworkUrl) {
        el.albumArt.innerHTML = `<img src="${artworkUrl}" alt="Album Art" onerror="this.style.display='none'; this.parentElement.innerHTML='🎵'">`;
    } else {
        el.albumArt.innerHTML = '🎵';
    }

    // Update playback status (handle Status field)
    const status = data.Status || data.status || data.playbackStatus || '';
    const isPlaying = status === 'Playing' || status === 'playing';
    if (isPlaying) {
        el.statusIcon.textContent = '▶️';
        el.statusText.textContent = 'Playing';
        el.playbackStatus.className = 'playback-status playing';
    } else {
        el.statusIcon.textContent = '⏸️';
        el.statusText.textContent = 'Paused';
        el.playbackStatus.className = 'playback-status paused';
    }

    // Update progress bar (handle Position and Length fields - they're strings in microseconds)
    const trackPosition = parseInt(data.Position || data.position || 0);
    const trackLength = parseInt(data.Length || data.length || data.duration_ms || 0);

    if (trackLength > 0) {
        const percentage = (trackPosition / trackLength) * 100;
        el.progressFill.style.width = `${percentage}%`;
        el.currentTime.textContent = formatTime(trackPosition);
        el.totalTime.textContent = formatTime(trackLength);
    } else {
        el.progressFill.style.width = '0%';
        el.currentTime.textContent = '0:00';
        el.totalTime.textContent = '0:00';
    }

    // Update metadata (handle Player field)
    el.volume.textContent = data.Volume ? data.Volume : (data.volume ? `${Math.round(data.volume * 100)}%` : 'N/A');
    el.player.textContent = data.Player || data.player || data.playerName || 'Unknown';
    el.format.textContent = data.Format || data.format || data.codec || 'N/A';
}
