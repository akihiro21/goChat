window.onload = function () {
    let conn;
    const msg = document.querySelectorAll("#in_form");
    const usr = document.getElementById("user");
    const log = document.getElementById("scroll-inner");
    const alert = document.getElementById("alert");

    function scroll() {
        log.scrollIntoView(false);
    }

    function createHTML(message, name) {
        let divChat = document.createElement("div");
        divChat.className = "chat_container"
        log.appendChild(divChat)
        let divFlex = document.createElement("div");
        divFlex.className = "flex";
        divChat.appendChild(divFlex);
        let divIcon = document.createElement("div");
        divIcon.className = "icon";
        divFlex.appendChild(divIcon);
        let pMsg = document.createElement("p");
        pMsg.id = "mes"
        let text = document.createTextNode(message);
        divFlex.appendChild(pMsg);
        pMsg.appendChild(text);
        let icon = document.createElement("img");
        if (name == "OrangeBot" || name == "admin") {
            icon.src = "/static/image/bot.jpg";
        } else {
            icon.src = "/static/image/user2.jpg";
        }
        let user = document.createElement("p");
        let userText = document.createTextNode(name);
        divIcon.appendChild(icon);
        divIcon.appendChild(user);
        user.appendChild(userText);
        scroll();
    }

    scroll();

    function DeleteNewLine(myLen) {
        var newLen = '';
        for (var i = 0; i < myLen.length; i++) {
            text = escape(myLen.substring(i, i + 1));
            if (text != "%0D" && text != "%0A") {
                newLen += myLen.substring(i, i + 1);
            }
        }
        return (newLen);
    }

    var form = document.querySelectorAll("#form");

    for (let i = 0; i < form.length; i++) {
        form[i].addEventListener("submit", function (event) {
            let jsonData = {}
            msg[i].value = DeleteNewLine(msg[i].value);
            if (!conn) {
                return false;
            }
            if (!msg[i].value || msg[i].value == "") {
                return false;
            }
            jsonData["username"] = usr.value;
            jsonData["message"] = msg[i].value;
            let json = JSON.stringify(jsonData)
            conn.send(json)
            return false;
        });
    };


    if (window["WebSocket"]) {
        const params = window.location.href.split("/");
        const roomId = params[params.length - 1];
        conn = new WebSocket("ws://" + document.location.host + "/ws/" + roomId);
        conn.onclose = function (evt) {
            window.alert("failed please reload");
        };
        conn.onopen = () => {
            console.log("Successfully connected");
        }

        conn.onmessage = json => {
            const reserved = JSON.parse(json.data);
            let messages = reserved.message;
            let username = reserved.username;
            createHTML(messages, username);
        };


    } else {
        window.alert("Your browser does not support WebSockets.");
    }
};