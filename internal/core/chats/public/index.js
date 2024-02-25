var conn;
var message_pane = document.getElementById("messages");
var message = document.getElementById("message");

message.addEventListener("keydown", (e) => {
  if (e.key === "Enter") {
    send();
  }
});

function send() {
  if (!conn || !message.value) return;

  conn.send(
    JSON.stringify({
      content: message.value,
      sent_at: new Date(Date.now()).toISOString(),
      tag: `tag=${Math.random().toString(16).substring(2)}`,
    })
  );

  message.value = "";
}

function push_message(e) {
  let d = new Date(e.sent_at);
  let node = document.createElement("div");
  node.className = "mb-3";
  node.innerHTML = `
  <div class="text-xs">
    <span class=" font-semibold">${d.toDateString()}</span>&nbsp;${e.sender_id}
  </div>
  <div class="text-md">${e.content}</div>
  `;

  message_pane.appendChild(node);
  message_pane.scrollTo(0, message_pane.scrollHeight);
}

function add() {
  open_chat(document.getElementById("user-id").value);
}

function open_chat(id) {
  conn.send(
    JSON.stringify({
      sent_at: new Date(Date.now()).toISOString(),
      tag: `join=${id}`,
    })
  );

  fetch(`http://localhost:8000/api/v1/chats/${id}`, {
    method: "GET",
    credentials: "include",
  })
    .then((res) => res.json())
    .then((res) => {
      message_pane.innerHTML = "";

      res.forEach((e) => {
        push_message(e);
      });
    });
}

function get_all_chats() {
  message_pane.innerHTML = "";

  let users = document.getElementById("users");
  users.innerHTML = "";

  fetch("http://localhost:8000/api/v1/chats", {
    method: "GET",
    credentials: "include",
  })
    .then((res) => res.json())
    .then((res) => {
      res.forEach((e) => {
        let node = document.createElement("div");
        node.className =
          "border-t border-gray-400 flex flex-row select-none cursor-pointer hover:bg-gray-200";
        node.onclick = () => open_chat(e.user_id);
        node.innerHTML = `
        <div>
          <div class="text-xs w-full">${e.user_id}</div>
          <div class="text-xs w-full">${e.content}</div>
        </div>
        <div class="text-xs w-4 flex items-center justify-center">${e.unread_messages}</div>
        `;

        users.appendChild(node);
      });
    })
    .catch((err) => console.error(err));
}

function connect() {
  if (window["WebSocket"]) {
    conn = new WebSocket(`ws://localhost:8000/ws/chats`);

    conn.onopen = function (evt) {
      get_all_chats();

      let node = document.createElement("div");
      node.innerHTML = "<b>Connected.</b>";
      message_pane.appendChild(node);
    };

    conn.onclose = function (evt) {
      console.log(evt.reason);
      var item = document.createElement("div");
      item.innerHTML = "<b>Connection closed.</b>";
      message_pane.appendChild(item);
    };

    conn.onmessage = function (evt) {
      console.log(evt.data);
      var msg = JSON.parse(evt.data);
      push_message(msg);
    };
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
      email: document.getElementById("user").value,
      password: document.getElementById("pass").value,
    }),
  })
    .then((res) => res.json())
    .then((res) => console.log(res));
}
