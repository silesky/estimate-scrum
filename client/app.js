const ESTIMATION_INPUT_ID = 'estimation';
const SUBMIT_BUTTON_ID = 'submit';
const WS_URL = 'ws://localhost:3333/ws';

const socket = new WebSocket(WS_URL);

socket.addEventListener('open', function() {
  console.log('Websocket opened @ ' + WS_URL);
});

socket.addEventListener('message', function(event) {
  console.log('Message from server ', event.data);
});

const createMessage = (username, estimate) => {
  const est = typeof estimate === 'string' ? parseInt(estimate, 10) : estimate;
  return JSON.stringify({ username, estimate: est });
};

document.getElementById(SUBMIT_BUTTON_ID).addEventListener('click', () => {
  const value = document.getElementById(ESTIMATION_INPUT_ID).value;
  const msg = createMessage('Seth', value);
  console.log('msg submitted:', msg);
  socket.send(msg);
});
