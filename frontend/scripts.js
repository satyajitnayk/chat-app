let selectedChat = 'general';
let conn;

class Event {
  constructor(type, payload) {
    this.type = type;
    this.payload = payload;
  }
}

function routeEvent(event) {
  if (!event.type) {
    alert('no type field in the event');
  }

  switch (event.type) {
    case 'new_message':
      console.log('new message');
      break;
    default:
      alert('unsupported message type');
      break;
  }
}

function sendEvent(eventName, payload) {
  const event = new Event(eventName, payload);

  conn.send(JSON.stringify(event));
}

function changeChatRoom() {
  let newchat = document.getElementById('chatroom');
  if (newchat && newchat.value != selectedChat) {
    console.log(newchat);
  }
  return false; // this will prevent form to navigate to different URL
}

function sendMessage() {
  let newmessage = document.getElementById('message');
  if (newmessage) {
    sendEvent('send_message', newmessage.value);
  }
  return false;
}

window.onload = function () {
  document.getElementById('chatroom-selection').onsubmit = changeChatRoom;
  document.getElementById('chatroom-message').onsubmit = sendMessage;

  // check if user browser supports websocket
  if (window['WebSocket']) {
    console.log('browser supports websockets');
    // connect to websocket
    conn = new WebSocket(`ws://${document.location.host}/ws`);

    conn.onmessage = function (e) {
      const eventData = JSON.parse(e.data);

      const event = Object.assign(new Event(), eventData);
      routeEvent(event);
    };
  } else {
    console.log("browser doen't support websockets");
  }
};
