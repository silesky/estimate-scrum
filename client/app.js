// Websocket URL
const WS_URL = 'ws://localhost:3333/ws';

// Form element IDs
const USERNAME = 'username';
const ESTIMATION = 'estimation';
const SUBMIT = 'submit';
const DISPLAY_AREA = 'display';

const getValue = prop => document.getElementById(prop).value;

// utils
const renderMessage = msg =>
  `<div class="message">${msg.username} -> ${msg.estimate}</div>`;

const createMessage = (username, estimate) => {
  const est = typeof estimate === 'string' ? parseInt(estimate, 10) : estimate;
  return JSON.stringify({ username, estimate: est });
};

// globals
let messages = [];

// initialization
const socket = new WebSocket(WS_URL);

socket.addEventListener('open', function() {
  console.log('Websocket opened @ ' + WS_URL);
});

socket.addEventListener('message', function(event) {
  const data = JSON.parse(event.data);
  messages.push(data);
  console.log(messages.map(renderMessage))
  document.getElementById(DISPLAY_AREA).innerHTML = messages.map(renderMessage).join('');
  console.log('Message from server ', event.data);
});

document.getElementById(SUBMIT).addEventListener('click', () => {
  const msg = createMessage(getValue(USERNAME), getValue(ESTIMATION));
  console.log('msg submitted:', msg);
  socket.send(msg);
});
