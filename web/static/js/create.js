$(document).ready(function (e) {
    $('.createroom').on('click', function () {
        $('.form').stop().slideToggle();
    });
})

window.onload = function () {
    let conn;
    const msg = document.getElementById("in_form");
    const log = document.getElementById("scroll-inner");
    const alert = document.getElementById("alert");

    function scroll() {
        log.scrollIntoView(false);
    }

    function createHTML(msg, name) {
        let divChat = document.createElement("div");
        divChat.className = "chat_container"
        log.appendChild(divChat)
        let divFlex = document.createElement("div");
        divFlex.className = "flex";
        divChat.appendChild(divFlex);
        let icon = document.createElement("img");
        icon.src = "/static/image/user2.jpg";
        let pMsg = document.createElement("p");
        let text = document.createTextNode(msg);
        divFlex.appendChild(icon);
        divFlex.appendChild(pMsg);
        pMsg.appendChild(text);
        scroll();

    }

    function appendLog(item) {
        alert.appendChild(item)
    }

    scroll();

    var form = document.querySelector("#form");
    form.addEventListener("submit", function (event) {
        let jsonData = {}
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }
        jsonData["username"] = "username"
        jsonData["message"] = msg.value;
        let json = JSON.stringify(jsonData)
        conn.send(json)
        return false;
    });

    if (window["WebSocket"]) {
        const params = window.location.href.split("/");
        const roomId = params[params.length - 1];
        conn = new WebSocket("ws://" + document.location.host + "/ws/" + roomId);
        conn.onclose = function (evt) {
            console.log("failed")
        };
        conn.onopen = () => {
            console.log("Successfully connected")
        }

        conn.onmessage = json => {
            const msj = JSON.parse(json.data);
            let messages = msj.message;
            let username = msj.username;
            createHTML(messages, username);
        };


    } else {
        var text = document.createTextNode("Your browser does not support WebSockets.");
        appendLog(text);
    }
};