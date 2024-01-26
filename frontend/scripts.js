let selectedChat = 'general';
let conn;

class Event {
  constructor(type, payload) {
    this.type = type;
    this.payload = payload;
  }
}

class SendMessageEvent {
  constructor(message, from) {
    this.message = message;
    this.from = from;
  }
}

class ReceiveMessageEvent {
  constructor(message, from, sent) {
    this.message = message;
    this.from = from;
    this.sent = sent;
  }
}

class ChangeChatRoomEvent {
  constructor(name) {
    this.name = name;
  }
}

function changeChatRoom() {
  let newchat = document.getElementById('chatroom');
  if (newchat && newchat.value != selectedChat) {
    selectedChat = newchat.value;
    header = document.getElementById('chat-header').innerHTML =
      'Currently in chatroom: ' + selectedChat;

    let changeEvent = new ChangeChatRoomEvent(selectedChat);
    sendEvent('change_chat_room', changeEvent);

    // clean the chat area as changing chatroom
    textarea = document.getElementById('chatmessages');
    textarea.innerHTML = `You changed room into: ${selectedChat}`;
  }
  return false; // to avoid redirection of form
}

function routeEvent(event) {
  if (!event.type) {
    alert('no type field in the event');
  }

  switch (event.type) {
    case 'receive_message':
      const messageEvent = Object.assign(
        new ReceiveMessageEvent(),
        event.payload
      );
      appendChatMessage(messageEvent);
      break;
    default:
      alert('unsupported message type');
      break;
  }
}

function appendChatMessage(messageEvent) {
  const date = new Date(messageEvent.sent);
  const formattedMsg = `${messageEvent.from} \n${date.toLocaleString()}: ${
    messageEvent.message
  }`;

  const textarea = document.getElementById('chatmessages');
  // append new msg to text area
  textarea.innerHTML = textarea.innerHTML + '\n' + formattedMsg;
  // scroll to height
  textarea.scrollTop = textarea.scrollHeight;
}

function sendEvent(eventName, payload) {
  const event = new Event(eventName, payload);

  conn.send(JSON.stringify(event));
}

function sendMessage() {
  let newmessage = document.getElementById('message');
  if (newmessage) {
    // TODO: pick username from client side
    let outgoingEvent = new SendMessageEvent(newmessage.value, 'user1');
    sendEvent('send_message', outgoingEvent);
  }
  return false;
}

function login() {
  let formData = {
    username: document.getElementById('username').value,
    password: document.getElementById('password').value,
  };
  fetch('login', {
    method: 'POST',
    body: JSON.stringify(formData),
    mode: 'cors',
  })
    .then((res) => {
      if (res.ok) {
        return res.json();
      } else {
        throw 'unauthorized';
      }
    })
    .then((data) => {
      // connect to websocket
      connectWebsocket(data.otp);
    })
    .catch((err) => {
      alert(err);
    });
  return false; // to avoid redirection
}

function connectWebsocket(otp) {
  // check if user browser supports websocket
  if (window['WebSocket']) {
    console.log('browser supports websockets');
    // connect to secure websocket
    conn = new WebSocket(`wss://${document.location.host}/ws?otp=${otp}`);

    conn.onopen = function (e) {
      document.getElementById('connection-header').innerHTML =
        'connected to websocket: true';
    };

    conn.onclose = function (e) {
      document.getElementById('connection-header').innerHTML =
        'connected to websocket: false';
    };

    conn.onmessage = function (e) {
      const eventData = JSON.parse(e.data);

      const event = Object.assign(new Event(), eventData);
      routeEvent(event);
    };
  } else {
    console.log("browser doen't support websockets");
  }
}

window.onload = function () {
  document.getElementById('chatroom-selection').onsubmit = changeChatRoom;
  document.getElementById('chatroom-message').onsubmit = sendMessage;
  document.getElementById('login-form').onsubmit = login;
};
