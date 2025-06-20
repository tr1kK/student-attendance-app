// src/utils/broadcast.ts
const CHANNEL_NAME = 'attendance_channel';
const channel = new BroadcastChannel(CHANNEL_NAME);

export const postMessage = (message: any) => {
  channel.postMessage(message);
};

export const addMessageListener = (callback: (event: MessageEvent) => void) => {
  channel.addEventListener('message', callback);
};

export const removeMessageListener = (callback: (event: MessageEvent) => void) => {
  channel.removeEventListener('message', callback);
}; 