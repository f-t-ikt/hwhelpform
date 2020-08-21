var room = document.getElementById("room");
var help = document.getElementById("help");
var call = document.getElementById("call");

socket.onopen = function() {
    init();
};

socket.onmessage = function(e) {
    var post = JSON.parse(e.data);
    if (post.Method == "help") {
        addCol(help, post.Id, "Help");
    } else if (post.Method == "call") {
        addCol(call, post.Id, "Call");
    } else if (post.Method == "deleteHelp" || "deleteCall") {
        var element = document.getElementById(post.Id);
        deleteCol(element);
    }
};

function init() {
    while (help.rows[1]) {
        help.deleteRow(1);
    }
    while (call.rows[1]) {
        call.deleteRow(1);
    }
}

function deletePost(id, field) {
    if (socket.readyState > WebSocket.OPEN) {
        document.getElementById("top").innerHTML += "<div style=\"color: red;\">切断されました。リロードしてください。</div>";
        return;
    }
    
    socket.send(JSON.stringify(
        {
            Method: "delete" + field,
            Id: id
        }
    ));
};

function addCol(table, id, field) {
    var row = table.insertRow(-1);
    row.id = id;
    var cell1 = row.insertCell(-1);
    var cell2 = row.insertCell(-1);
    cell1.innerHTML = id;
    var button = document.createElement("button");
    button.onclick = function(){deletePost(id, field)};
    button.innerHTML = "確認"
    cell2.appendChild(button);
};

function deleteCol(obj) {
    var table = obj.parentNode;
    table.deleteRow(obj.sectionRowIndex);
};
