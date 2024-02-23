function connect() {
  var conn;
  var msg = document.getElementById("msg");
  var log = document.getElementById("log");

  function appendLog(item) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
      log.scrollTop = log.scrollHeight - log.clientHeight;
    }
  }

  document.getElementById("form").onsubmit = function () {
    if (!conn) {
      return false;
    }

    if (!msg.value) {
      return false;
    }

    conn.send(msg.value);
    msg.value = "";

    return false;
  };

  if (window["WebSocket"]) {
    conn = new WebSocket(`ws://localhost:8000/ws/chats`);

    conn.onopen = function (evt) {
      var item = document.createElement("div");
      item.innerHTML = "<b>Connected.</b>";
      appendLog(item);
    };

    conn.onclose = function (evt) {
      console.log(evt.reason);
      var item = document.createElement("div");
      item.innerHTML = "<b>Connection closed.</b>";
      appendLog(item);
    };

    conn.onmessage = function (evt) {
      console.log(evt);
      var msg = JSON.parse(evt.data);
      var item = document.createElement("div");
      item.innerHTML = `<b>${new Date(msg.created_at).toLocaleTimeString()} ${
        msg.sender.username
      }</b>: ${msg.content}`;
      appendLog(item);
    };
  } else {
    var item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    appendLog(item);
  }
}

function login() {
  fetch("http://localhost:8000/api/v1/login", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email: "johnd@email.com",
      password: "abcdefg",
    }),
  }).then((res) => console.log(res));
}
