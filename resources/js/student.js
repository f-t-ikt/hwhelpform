var room = document.getElementById("room");
var help = document.getElementById("help");
var call = document.getElementById("call");

socket.onopen = function() {
    init();
};

socket.onmessage = function(e) {
    var post = JSON.parse(e.data);
    if (post.Method == "help") {
        addCol(help, post.Id);
    } else if (post.Method == "call") {
        addCol(call, post.Id);
    } else if (post.Method == "deleteHelp" || post.Method == "deleteCall") {
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

function iso8601(date) {
  return date.getUTCFullYear()
    + "-" + (date.getUTCMonth()+1)
    + "-" + date.getUTCDate()
    + "T" + date.getUTCHours()
    + ":" + date.getUTCMinutes()
    + ":" + date.getUTCSeconds() + "Z";
}

function send(method) {
    const roomId = room.value;
    if (roomId == "") return;
    if (socket.readyState > WebSocket.OPEN) {
        document.getElementById("top").innerHTML += "<div style=\"color: red;\">切断されました。リロードしてください。</div>";
        return;
    }
    
    socket.send(JSON.stringify(
        {
            Method: method,
            Id: Number(roomId),
            Date: iso8601(new Date())
        }
    ));
};

function addCol(table, id) {
    var row = table.insertRow(-1);
    row.id = id;
    var cell = row.insertCell(-1);
    cell.innerHTML = id;
};

function deleteCol(obj) {
    var table = obj.parentNode;
    table.deleteRow(obj.sectionRowIndex);
};