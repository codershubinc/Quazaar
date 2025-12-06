import { connect, disconnect, sendPing, sendCommand } from './websocket.js';
import { clearLog } from './ui.js';
import * as el from './elements.js';

// --- Event Listeners ---

// Connection
el.connectBtn.addEventListener('click', connect);
el.disconnectBtn.addEventListener('click', disconnect);
el.pingBtn.addEventListener('click', sendPing);
el.clearBtn.addEventListener('click', clearLog);

// Custom Command
el.sendBtn.addEventListener('click', () => {
    const text = el.messageInput.value;
    if (text) {
        try {
            // Try to parse as JSON first
            const json = JSON.parse(text);
            sendCommand(json);
        } catch (e) {
            // Send as raw string if not JSON
            sendCommand(text);
        }
        el.messageInput.value = '';
    }
});

el.messageInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        el.sendBtn.click();
    }
});

// Media Controls
el.prevBtn.addEventListener('click', () => sendCommand({ command: 'previous' }));
el.playBtn.addEventListener('click', () => sendCommand({ command: 'play' }));
el.pauseBtn.addEventListener('click', () => sendCommand({ command: 'pause' }));
el.nextBtn.addEventListener('click', () => sendCommand({ command: 'next' }));

// Volume Controls
el.volUpBtn.addEventListener('click', () => sendCommand({
    type: 'system',
    msg_of: 'volume',
    action: 'inc'
}));

el.volDownBtn.addEventListener('click', () => sendCommand({
    type: 'system',
    msg_of: 'volume',
    action: 'dec'
}));

el.volMuteBtn.addEventListener('click', () => sendCommand({
    type: 'system',
    msg_of: 'volume',
    action: 'mute'
}));

el.volSlider.addEventListener('change', (e) => {
    const val = parseInt(e.target.value, 10);
    el.volValue.textContent = val + '%';
    sendCommand({
        type: 'system',
        msg_of: 'volume',
        action: 'set',
        set_to_vol: val
    });
});

el.volSlider.addEventListener('input', (e) => {
    el.volValue.textContent = e.target.value + '%';
});

// Brightness Controls
el.brightUpBtn.addEventListener('click', () => sendCommand({
    type: 'system',
    msg_of: 'brightness',
    action: 'inc'
}));

el.brightDownBtn.addEventListener('click', () => sendCommand({
    type: 'system',
    msg_of: 'brightness',
    action: 'dec'
}));

el.brightSlider.addEventListener('change', (e) => {
    const val = parseInt(e.target.value, 10);
    el.brightValue.textContent = val + '%';
    sendCommand({
        type: 'system',
        msg_of: 'brightness',
        action: 'set',
        set_to: val
    });
});

el.brightSlider.addEventListener('input', (e) => {
    el.brightValue.textContent = e.target.value + '%';
});

// Spotify Controls
el.spotifyArtistBtn.addEventListener('click', () => {
    const artistId = el.spotifyArtistInput.value.trim();
    if (artistId) {
        sendCommand({
            message: 'spotify_artist',
            action: 'get',
            artistId: artistId
        });
    }
});

el.spotifyFollowBtn.addEventListener('click', () => {
    const artistId = el.spotifyArtistInput.value.trim();
    if (artistId) {
        sendCommand({
            message: 'spotify_artist',
            action: 'follow',
            artistId: artistId
        });
    }
});

el.spotifyArtistInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        el.spotifyArtistBtn.click();
    }
});

// Auto-connect on load
window.addEventListener('load', () => {
    // Optional: connect();
});
