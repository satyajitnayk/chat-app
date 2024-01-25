let selectedChat = 'general';
let conn;

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
    conn.send(newmessage.value);
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

    conn.onmessage = function (event) {
      console.log(event);
    };
  } else {
    console.log("browser doen't support websockets");
  }
};
