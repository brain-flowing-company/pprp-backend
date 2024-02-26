var conn;
var message_pane = document.getElementById("messages");
var message = document.getElementById("message");
var current_chat = "";
var chats = [];

message.addEventListener("keydown", (e) => {
  if (e.key === "Enter") {
    send();
  }
});

let sent_tags = {};
let messages = [];

function send_ws_msg(event, content, sent_at, value) {
  let tag = Math.random().toString(16).substring(2);

  let msg = {
    event,
    content,
    sent_at,
    tag,
  };

  sent_tags[tag] = value;
  conn.send(JSON.stringify(msg));
}

function send() {
  if (!conn || !message.value) return;

  let sent_at = new Date(Date.now()).toISOString();
  let node = push_message(
    {
      content: message.value,
      sent_at,
      sender_id: "sending",
    },
    true
  );

  send_ws_msg("MSG", message.value, sent_at, node);

  message.value = "";
}

function push_message(e, author) {
  console.log(e);
  let d = new Date(e.sent_at);
  let status = "";
  if (author) {
    if (status !== undefined) {
      status = e.read_at === null ? "sent" : "read";
    } else {
      status = "sending";
    }
  }

  let node = document.createElement("div");
  node.className = `mb-3 w-full ${author ? "text-right" : "text-left"}`;
  node.innerHTML = `
  <div class="text-xs">
    <span class=" font-semibold">${d.toDateString()}</span>
    &nbsp;<span id="sender-id">${e.sender_id}</span>
    &nbsp;<span id="status">${status}</span>
  </div>
  <div class="text-md">${e.content}</div>
  `;

  if (author && e.message_id !== undefined && e.read_at === null) {
    messages.push(node);
  }

  message_pane.appendChild(node);
  message_pane.scrollTo(0, message_pane.scrollHeight);

  return node;
}

function add() {
  open_chat(document.getElementById("user-id").value);
}

function open_chat(id) {
  let sent_at = new Date(Date.now()).toISOString();
  send_ws_msg("JOIN", id, sent_at, true);

  current_chat = id;

  fetch(`http://localhost:8000/api/v1/chats/${id}`, {
    method: "GET",
    credentials: "include",
  })
    .then((res) => res.json())
    .then((res) => {
      message_pane.innerHTML = "";

      res.forEach((e) => {
        push_message(e, e.receiver_id === current_chat);
      });
    });
}

function get_all_chats() {
  messages = [];
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
          <div class="text-xs w-full" id="content">${e.content}</div>
        </div>
        <div class="text-xs w-4 flex items-center justify-center" id="unread-messages">${e.unread_messages}</div>
        `;

        chats[e.user_id] = node;
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
      var msg = JSON.parse(evt.data);
      console.log(msg);

      switch (msg.event) {
        case "MSG":
          if (sent_tags[msg.tag] !== undefined) {
            console.log(msg.payload);

            sent_tags[msg.tag].querySelector("#status").innerText =
              msg.payload.read_at === null ? "sent" : "read";
            sent_tags[msg.tag].querySelector("#sender-id").innerText = msg.payload.sender_id;

            if (msg.payload.read_at === null) {
              messages.push(sent_tags[msg.tag]);
            }
            delete sent_tags[msg.tag];
          } else {
            push_message(msg.payload, msg.payload.receiver_id === current_chat);
          }

          break;

        case "READ":
          messages.forEach((e) => {
            e.querySelector("#status").innerText = "read";
          });

          messages = [];
          break;

        case "CHATS":
          node = chats[msg.payload.user_id];
          if (node !== undefined) {
            node.querySelector("#unread-messages").innerText =
              msg.payload.unread_messages === 0
                ? 0
                : parseInt(node.querySelector("#unread-messages").innerText) +
                  msg.payload.unread_messages;

            if (msg.payload.content !== "")
              node.querySelector("#content").innerText = msg.payload.content;
          }
          break;
      }
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
