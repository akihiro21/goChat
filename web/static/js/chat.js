function viewChange1() {
    if (document.getElementById('sample1')) {
        id = document.getElementById('sample1').value;
        if (id == 'select1') {
            document.getElementById('conditions1').style.display = "";
            document.getElementById('conditions2').style.display = "none";
        } else if (id == 'select2') {
            document.getElementById('conditions1').style.display = "none";
            document.getElementById('conditions2').style.display = "";
        }
    }

    window.onload = viewChange1;

}
function viewChange2() {
    if (document.getElementById('sample2')) {
        id = document.getElementById('sample2').value;
        if (id == 'select1') {
            for (var i = 1; i < 13; ++i) {
                if (i == 1 || i == 7) {
                    document.getElementById("Box" + i).style.display = "";
                } else {
                    document.getElementById("Box" + i).style.display = "none";
                }
            }
        } else if (id == 'select2') {
            for (var i = 1; i < 13; ++i) {
                if (i == 2 || i == 8) {
                    document.getElementById("Box" + i).style.display = "";
                } else {
                    document.getElementById("Box" + i).style.display = "none";
                }
            }
        } else if (id == 'select3') {
            for (var i = 1; i < 13; ++i) {
                if (i == 3 || i == 9) {
                    document.getElementById("Box" + i).style.display = "";
                } else {
                    document.getElementById("Box" + i).style.display = "none";
                }
            }
        } else if (id == 'select4') {
            for (var i = 1; i < 13; ++i) {
                if (i == 4 || i == 10) {
                    document.getElementById("Box" + i).style.display = "";
                } else {
                    document.getElementById("Box" + i).style.display = "none";
                }
            }
        } else if (id == 'select5') {
            for (var i = 1; i < 13; ++i) {
                if (i == 5 || i == 11) {
                    document.getElementById("Box" + i).style.display = "";
                } else {
                    document.getElementById("Box" + i).style.display = "none";
                }
            }
        } else if (id == 'select6') {
            for (var i = 1; i < 13; ++i) {
                if (i == 6 || i == 12) {
                    document.getElementById("Box" + i).style.display = "";
                } else {
                    document.getElementById("Box" + i).style.display = "none";
                }
            }
        }
        window.onload = viewChange2;
    }
}
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
        if (name == "Orange" || name == "admin") {
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