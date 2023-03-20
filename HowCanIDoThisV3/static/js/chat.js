var socket = null;
function openChat(chatId, avatarImg, userName) {
    if (window.screen.width <= 1000) {
        $("#header").css("display", "none")
        $("#plist").css("display", "none")
        $("#currentChatId").css("display", "block")
        $('#backButton').css('display', 'block')
    }
    $('#messagePlace').empty()
    $('#currentChatAvatar').attr('src', '/static/upload/avatars/' + avatarImg)
    $('#currentChatName').html(userName)
    
    $.get('/chat/' + chatId, function () {
    }).done(function (response) {
        var responseData = response;
        for (var i = 0; i < responseData['mesFor'].length; i++) {
            let fileNameInsert = responseData['mesFor'][i]['filePath'].split('||')[0];
            if (fileNameInsert.length > 25) {
                fileNameInsert = fileNameInsert.slice(0, 25) + '...'
            }
            let sendTimeArr = responseData['mesFor'][i]['sendTime'].split('T')
            let sendTime = sendTimeArr[1].split(':')[0] + ':' + sendTimeArr[1].split(':')[1] + " " + sendTimeArr[0]
            if (responseData['mesFor'][i]['userSenderId'] == responseData['userHowSend']) {
                if (responseData['mesFor'][i]['filePath'] != "") {
                    $('#messagePlace').append(`
                        <li class="clearfix" style="z-index: -10">
                        <div class="message-data text-right" style="text-align: right;">
                            <span class="message-data-time ">` + sendTime + `</span>
                        </div>
                        <div>
                            <div class="messageImg other-message float-right" onclick="window.open('/static/upload/messeges/` + responseData['mesFor'][i]['userSenderId'] + responseData['mesFor'][i]['filePath'] + `', '_blank')">
                            <div style="width: 50px; height: 50px;vertical-align: middle; text-align: center; display: inline-block;">
                            <i class="fa fa-file-o" aria-hidden="true" style="font-size:40px; color: #E746E0; top: 10%; position: relative;"></i>
                            </div>
                            <p class="chatFile">` + fileNameInsert + `</p>
                            </div>
                        </div>0
                        </li>
                    `)
                }
                if (responseData['mesFor'][i]['messageText'] != "") {
                    $('#messagePlace').append(`
                        <li class="clearfix" style="z-index: -10">
                            <div class="message-data text-right" style="text-align: right;">
                            <span class="message-data-time ">` + sendTime + `</span>
                            </div>
                            <div class="message other-message float-right">` + responseData['mesFor'][i]['messageText'] + `</div>
                        </li>
                    `)
                }
            } else {
                if (responseData['mesFor'][i]['filePath'] != "") {
                    $('#messagePlace').append(`
                        <li class="clearfix" style="z-index: -10">
                        <div class="message-data">
                            <span class="message-data-time ">` + sendTime + `</span>
                        </div>
                        <div>
                            <div class="messageImg my-message float-right" onclick="window.open('/static/upload/messeges/` + responseData['mesFor'][i]['userSenderId'] + responseData['mesFor'][i]['filePath'] + `', '_blank')">
                            <img src="/static/img/file.png" alt="" style="width: 50px; height: 50px; display: inline-block;">
                            <p
                                class="chatFile">` +
                                    fileNameInsert + `</p>
                            </div>
                        </div>0
                        </li>
                    `)
                }
                if (responseData['mesFor'][i]['messageText'] != "") {
                    $('#messagePlace').append(`
                        <li class="clearfix" style="z-index: -10">
                            <div class="message-data">
                            <span class="message-data-time">` + sendTime + `</span>
                            </div>
                            <div class="message my-message" >` + responseData['mesFor'][i]['messageText'] + `</div>
                        </li>
                    `)
                }
            }
        }
        scrollChat();
    });
    if (socket != null){
        socket.close();
    }
    socket = new WebSocket("ws://" + getAdress() + "/chat/ws/" + chatId)
    socket.onopen = function () {
        console.log("Status: Connected\n");
    };
    socket.onmessage = function (responseData) {
        responseData = JSON.parse(responseData.data)
        let fileNameInsert = responseData['filePath'].split('||')[0];
        if (fileNameInsert.length > 25) {
            fileNameInsert = fileNameInsert.slice(0, 25) + '...'
        }
        let sendTimeArr = responseData['sendTime'].split(' ')
        let sendTime = sendTimeArr[1].split(':')[0] + ':' + sendTimeArr[1].split(':')[1] + " " + sendTimeArr[0]
        if (responseData['userSenderId'] == document.getElementById("currentUserId").value) {
            if (responseData['filePath'] != "") {
                document.getElementById("messagePlace").innerHTML += `
                <li class="clearfix" style="z-index: -10">
                    <div class="message-data text-right" style="text-align: right;">
                    <span class="message-data-time ">` + sendTime + `</span>
                    </div>
                    <div class="messageImg other-message float-right" onclick="window.open('/static/upload/messeges/` + responseData['userSenderId'] + responseData['filePath'] + `', '_blank')">
                    <img src="/static/img/file.png" alt="" style="width: 50px; height: 50px; display: inline-block;">
                    <p
                    class="chatFile">` +
                        fileNameInsert + `</p>
                    </div>
                </li>
                `
            }
            else if (responseData['messageText'] != "") {
                document.getElementById("messagePlace").innerHTML += `
                <li class="clearfix" style="z-index: -10">
                    <div class="message-data text-right" style="text-align: right;">
                    <span class="message-data-time ">` + sendTime + `</span>
                    </div>
                    <div class="message other-message float-right">` + responseData['messageText'] + `</div>
                </li>
                    `
            }
        } else {
            if (responseData['filePath'] != "") {
                document.getElementById("messagePlace").innerHTML += `
                <li class="clearfix" style="z-index: -10">
                    <div class="message-data">
                        <span class="message-data-time ">` + sendTime + `</span>
                    </div>
                    <div>
                        <div class="messageImg my-message float-right" onclick="window.open('/static/upload/messeges/` + responseData['userSenderId'] + responseData['filePath'] + `', '_blank')">
                        <img src="/static/img/file.png" alt="" style="width: 50px; height: 50px; display: inline-block; ">
                        <p
                        class="chatFile">` +
                            fileNameInsert + `</p>
                        </div>
                    </div>0
                    </li>
                `
            }
            else if (responseData['messageText'] != "") {
                document.getElementById("messagePlace").innerHTML += `
                <li class="clearfix" style="z-index: -10">
                    <div class="message-data">
                    <span class="message-data-time">` + sendTime + `</span>
                    </div>
                    <div class="message my-message" >` + responseData['messageText'] + `</div>
                </li>
                    `
            }
        }
        scrollChat()
    };
    document.getElementById("buttonSendMessage").addEventListener("click", function () {
        var input = document.getElementById("inputSendValue").value;
        let currentUserId = document.getElementById("currentUserId").value
        var sendFileOrNot = true
        var fileName = ""
        if (document.getElementById("fileInputChat").files.length == 0) {
            sendFileOrNot = false;
        }
        else if (document.getElementById("fileInputChat").files[0].name == "") {
            sendFileOrNot = false
        }
        else {
            fileName = document.getElementById("fileInputChat").files[0].name
        }
        if (sendFileOrNot) {
            sendFileChat()
        }
        if (input != "" || fileName != "") {
            closeFileChat()
            socket.send(JSON.stringify({ "Id": currentUserId, "Message": input, "FileName": fileName }));
            document.getElementById("inputSendValue").value = "";
        }
    });

}

function sendFileChat() {
    let fn = document.getElementById("fileInputChat").files[0]
    let formData = new FormData()
    formData.append('document', fn)
    fetch('/chat/saveChatFile', {
        method: "POST",
        body: formData
    })
}

function scrollChat() {
    var chat = document.getElementById("chatPlace")
    if (window.screen.width <= 1000) {
        document.getElementById('linkToEnd').click();
    }
    else {
        chat.scrollTop = chat.scrollHeight; 
    }
    console.log('gg')
}
$("#fileInputChat").on('change', function (event) {
    var fileName = event.target.files[0].name;
    $("#choosenFile").html(fileName)
    $("#choosenFile").css('display', 'inline-block')
    $("#closeFile").css('display', 'inline-block')

});
const textarea = document.getElementById("sendMessageInput");
if (textarea) {
    textarea.addEventListener("input", function (event) {
        this.style.height = "";
        this.style.height = this.scrollHeight + "px";
    });
}

function closeFileChat() {
    $("#fileInputChat").val("")
    $("#choosenFile").css('display', 'none')
    $("#closeFile").css('display', 'none')
}

let touchstartX = 0;
let touchendX = 0;

function handleGesture() {
    if (touchendX - touchstartX >= 80) {
        $("#header").css("display", "block")
        $("#plist").css("display", "block")
        $("#currentChatId").css("display", "none")
    }
}
function fuckGoBack() {
    $("#header").css("display", "block")
    $("#plist").css("display", "block")
    $("#currentChatId").css("display", "none")
}
document.addEventListener("touchstart", function (event) {
    touchstartX = event.touches[0].clientX;
}, false);

document.addEventListener("touchend", function (event) {
    touchendX = event.changedTouches[0].clientX;
    handleGesture();
}, false);