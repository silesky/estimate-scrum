// Websocket URL
const WS_URL = 'ws://localhost:3333/ws';

// Form element IDs
const USERNAME = 'username';
const ESTIMATION = 'estimation';
const SUBMIT = 'submit';
const DISPLAY_AREA = 'display';
const SESSION_ID = window.location.pathname.split("/")[1]
const getValue = prop => document.getElementById(prop).value;

// utils
const renderMessage = msg =>
  `<div class="message">${msg.username} -> ${msg.estimate}</div>`;

const createMessage = (username, estimate, sessionID) => {
  const est = typeof estimate === 'string' ? parseInt(estimate, 10) : estimate;
  return JSON.stringify({ username, estimate: est, sessionID });
};

// globals
let messages = [];

// initialization
const socket = new WebSocket(WS_URL);

socket.addEventListener('open', function() {
  console.log('Websocket opened @ ' + WS_URL);
});

// user entering an estimation.
socket.addEventListener('message', function(event) {
  const data = JSON.parse(event.data);
  messages.push(data);
  console.log(messages.map(renderMessage))
  document.getElementById(DISPLAY_AREA).innerHTML = messages.map(renderMessage).join('');
  console.log('Message from server ', event.data);
});

// user submitting an estimation.
document.getElementById(SUBMIT).addEventListener('click', () => {
  const msg = createMessage(getValue(USERNAME), getValue(ESTIMATION), SESSION_ID);
  console.log('msg submitted:', msg);
  socket.send(msg);
});
