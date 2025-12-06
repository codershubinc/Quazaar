import { state } from './state.js';
import { addMessage, updateStatus, updateStats, updateMediaInfo } from './ui.js';
import * as el from './elements.js';

export function connect() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host; // Use current host (e.g., localhost:8765)
    const url = `${protocol}//${host}/ws?deviceId=$2a$10$jWT5DfCYez7vSyrR2NiBg.REJDNvP5dxy8Pr0uyuJXqGgg3XHpqv2`;
    addMessage('info', `Connecting to ${url}...`);

    state.ws = new WebSocket(url);

    state.ws.onopen = function (event) {
        addMessage('info', 'Connected successfully!');
        updateStatus(true);
    };

    state.ws.onmessage = function (event) {
        handleWebSocketMessage(event);
    };

    state.ws.onerror = function (error) {
        state.errorCount++;
        updateStats();
        addMessage('error', 'WebSocket error occurred');
        console.error('WebSocket error:', error);
    };

    state.ws.onclose = function (event) {
        addMessage('info', `Connection closed (code: ${event.code})`);
        updateStatus(false);
        state.ws = null;
    };
}

export function disconnect() {
    if (state.ws) {
        state.ws.close();
        addMessage('info', 'Disconnecting...');
    }
}

export function sendPing() {
    sendCommand({ command: 'ping' });
}

export function sendCommand(command) {
    if (state.ws && state.ws.readyState === WebSocket.OPEN) {
        const message = typeof command === 'string' ? command : JSON.stringify(command);
        state.ws.send(message);
        state.sentCount++;
        updateStats();
        addMessage('sent', message);
    } else {
        addMessage('error', 'Not connected to WebSocket');
    }
}

function handleWebSocketMessage(event) {
    try {
        const message = JSON.parse(event.data);

        // Handle Spotify artist messages
        if (message.message === 'spotify_artist_get' || message.message === 'spotify_artist_follow') {
            console.log('🎵 Spotify Artist Response:', message);
            state.receivedCount++;
            updateStats();

            if (message.status === 'success') {
                if (message.message === 'spotify_artist_get') {
                    addMessage('received', `Artist Info: ${JSON.stringify(message.data, null, 2)}`);
                } else if (message.message === 'spotify_artist_follow') {
                    addMessage('received', 'Successfully followed artist!');
                }
            } else {
                addMessage('error', `Spotify Error: ${message.data?.error || 'Unknown error'}`);
            }
            return;
        }

        // Handle system messages (volume)
        if (message.message === 'system' || message.type === 'system') {
            console.log('🔊 System Response:', message);
            state.receivedCount++;
            updateStats();

            if (message.status === 'success') {
                const data = message.data || {};
                let text = `System Action: ${data.action || 'Unknown'}`;

                if (data.current_volume !== undefined) {
                    text += ` | Volume: ${data.current_volume}%`;
                    // Update slider if we have volume
                    el.volSlider.value = data.current_volume;
                    el.volValue.textContent = data.current_volume + '%';
                }

                if (data.current_brightness !== undefined) {
                    text += ` | Brightness: ${data.current_brightness}%`;
                    // Update slider if we have brightness
                    el.brightSlider.value = data.current_brightness;
                    el.brightValue.textContent = data.current_brightness + '%';
                }

                if (data.meta && data.meta.is_muted !== undefined) {
                    text += ` | Muted: ${data.meta.is_muted}`;
                    el.volMuteBtn.style.backgroundColor = data.meta.is_muted ? '#dc3545' : '';
                }

                addMessage('received', text);
            } else {
                addMessage('error', `System Error: ${message.data?.error || 'Unknown error'}`);
            }
            return;
        }

        // Handle media_info messages
        if (message.message === 'media_info' || message.Message === 'media_info') {
            const mediaData = message.data || message.Data || message;
            console.log('Media info received:', mediaData);
            updateMediaInfo(mediaData);
            return; // Don't log media info in messages
        }

        // Still show in log
        state.receivedCount++;
        updateStats();
        addMessage('received', event.data);
    } catch (e) {
        // If not JSON, just show as text
        state.receivedCount++;
        updateStats();
        addMessage('received', event.data);
    }
}
