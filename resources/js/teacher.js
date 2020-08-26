var help = document.getElementById("help");
var call = document.getElementById("call");

var timeago = new timeago();
timeago.setLocale("ja");

socket.onopen = function() {
    init();
};

socket.onmessage = function(e) {
    var post = JSON.parse(e.data);
    if (post.Method == "help") {
        addCol(help, post.Id, post.Date, "Help");
    } else if (post.Method == "call") {
        addCol(call, post.Id, post.Date, "Call");
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

function addCol(table, id, date, field) {
    var row = table.insertRow(-1);
    row.id = id;
    var idCell = row.insertCell(-1);
    var timeCell = row.insertCell(-1);
    var buttonCell = row.insertCell(-1);
    
    idCell.className = "roomid";
    idCell.innerHTML = id;
    
    var time = document.createElement("div");
    time.setAttribute("datetime", date);
    time.innerHTML = date;
    timeCell.className = "time";
    timeCell.appendChild(time);
    
    var button = document.createElement("button");
    button.onclick = function(){deletePost(id, field)};
    button.innerHTML = "確認"
    buttonCell.className = "button";
    buttonCell.appendChild(button);
    
    timeago.render(time)
};

function deleteCol(obj) {
    var table = obj.parentNode;
    table.deleteRow(obj.sectionRowIndex);
};
