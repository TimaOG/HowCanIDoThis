var currentColorState = 0;
if (!sessionStorage.getItem('isLoadTheme')) {
  sessionStorage.clear();
  $.get('/account/getThemeNumber', function () {
  }).done(function (response) {
    var responseData = JSON.parse(response);
    if (responseData["status"] == "ok") {
      currentColorState = Number(responseData['body'])
      changeSiteTheme(true)
      sessionStorage.setItem('isLoadTheme', true);
      sessionStorage.setItem('Theme', currentColorState);
    }
  });
}
else {
  currentColorState = sessionStorage.getItem('Theme')
  changeSiteTheme()
}


function changeSiteTheme(isLoad = false) {
  currentColorState = Number(currentColorState)
  if (isLoad == false) {
    currentColorState += 1
    if (currentColorState > 2) {
      currentColorState = 0
    }
  }
  console.log(currentColorState)
  switch (currentColorState) {
    case 0:
      document.documentElement.style.setProperty('--firstColor', '##3c5bff');
      document.documentElement.style.setProperty('--secondColor', '#ffffff');
      document.documentElement.style.setProperty('--thirdColor', '#888888');
      document.documentElement.style.setProperty('--fourthColor', '#BBBBBB');
      document.documentElement.style.setProperty('--fontColor', '#000000');
      break;
    case 1:
      document.documentElement.style.setProperty('--firstColor', '#325660');
      document.documentElement.style.setProperty('--secondColor', '#737891');
      document.documentElement.style.setProperty('--thirdColor', '#AAAAAA');
      document.documentElement.style.setProperty('--fourthColor', '#BA1236');
      document.documentElement.style.setProperty('--fontColor', '#BBBBBB');
      break;
    case 2:
      document.documentElement.style.setProperty('--firstColor', '#E746E0');
      document.documentElement.style.setProperty('--secondColor', '#28293E');
      document.documentElement.style.setProperty('--thirdColor', '#2a2c39');
      document.documentElement.style.setProperty('--fourthColor', '#1B1C2C');
      document.documentElement.style.setProperty('--fontColor', '#fff');
      break;
  }
  if(!isLoad){
    $.get('/account/changeThemeNumber/' + currentColorState, function () {
    }).done(function (response) {
      var responseData = JSON.parse(response);
      if (responseData["status"] == "ok") {

      }
      if (responseData["status"] == "notAuth") {
        window.location.href = 'login?target=account'
      }
    });
  }
  sessionStorage.setItem('Theme', currentColorState);

}



const select = (el, all = false) => {
  el = el.trim()
  if (all) {
    return [...document.querySelectorAll(el)]
  } else {
    return document.querySelector(el)
  }
}

const on = (type, el, listener, all = false) => {
  let selectEl = select(el, all)
  if (selectEl) {
    if (all) {
      selectEl.forEach(e => e.addEventListener(type, listener))
    } else {
      selectEl.addEventListener(type, listener)
    }
  }
}

on('click', '.mobile-nav-toggle', function (e) {
  select('#navbar').classList.toggle('navbar-mobile')
  this.classList.toggle('bi-list')
  this.classList.toggle('bi-x')
})
$('.switch-btn').click(function () {
  $(this).toggleClass('switch-on');
  if ($(this).hasClass('switch-on')) {
    $(this).trigger('on.switch');
  } else {
    $(this).trigger('off.switch');
  }
});

$('.input-file input[type=file]').on('change', function () {
  if (this.id == 'new_file') {
    let file = this.files[0];
    $(this).closest('.input-file').find('.input-file-text').html(file.name);
  }
});

$('#switch').on('change', function () {
  if ($(this).is(":checked")) {
    $.post('/account/changeActiveStatus', { value: true }, function (data) {
    });
  } else {
    $.post('/account/changeActiveStatus', { value: false }, function (data) {
    });
  }
})

var readURL = function (input) {
  let formData = new FormData();
  if (input.files && input.files[0]) {
    var reader = new FileReader();

    reader.onload = function (e) {
      $('.profile-pic').attr('src', e.target.result);
    }

    reader.readAsDataURL(input.files[0]);
    formData.append("file", input.files[0]);
    fetch('/account/changeProfileImg', {
      method: "POST",
      body: formData
    }).then(function (response) {
      response.json().then((dataR) => {
        if (dataR["status"] == "notAuth") {
          window.location.href = 'login?target=account'
        }
      })
    });
  }
}
//Account
$(".file-upload").on('change', function () {
  readURL(this);
});

$(".upload-button").on('click', function () {
  $(".file-upload").click();
});

$("#saveSettingsButton").on('click', function () {
  var formData = new FormData(document.getElementById("saveSettingsForm"))
  fetch('/account/changeSettings', {
    method: "POST",
    body: formData
  }).then(function (response) {
    response.json().then((dataR) => {
      if (dataR["status"] == "notAuth") {
        window.location.href = 'login?target=account'
      }
      showNotify(dataR["body"])
    })

  });
});

$("#newOrderButton").on('click', function () {
  var formData = new FormData(document.getElementById("newOrderForm"))
  fetch('/account/newOrder', {
    method: "POST",
    body: formData
  }).then(function (response) {
    response.json().then((dataR) => {
      if (dataR["status"] == "ok") {
        addOrderToOrderList(formData, dataR["lastId"])
        showNotify(dataR["body"])
      }
      if (dataR["status"] == "notAuth") {
        window.location.href = 'login?target=account'
      }
    })
  });
});

$("#newOfferButton").on('click', function () {
  var formData = new FormData(document.getElementById("newOfferForm"))
  fetch('/account/newOffer', {
    method: "POST",
    body: formData
  }).then(function (response) {
    response.json().then((dataR) => {
      if (dataR["status"] == "ok") {
        addOfferToOfferList(formData, dataR["lastId"], dataR["fn"])
      }
      if (dataR["status"] == "notAuth") {
        window.location.href = 'login?target=account'
      }
      showNotify(dataR["body"])
    })
  });

});

$("#changeOfferButton").on('click', function () {
  var formData = new FormData(document.getElementById("changeOfferForm"))
  fetch('/account/changeOffer', {
    method: "POST",
    body: formData
  }).then(function (response) {
    response.json().then((dataR) => {
      if (dataR["status"] == "ok") {
        needId = "#offer" + formData.get("offerId")
        $(needId).remove();
        addOfferToOfferList(formData, formData.get("offerId"), dataR["fn"])
      }
      if (dataR["status"] == "notAuth") {
        window.location.href = 'login?target=account'
      }
      showNotify(dataR["body"])
    })
  });

});
//Offers
$("#showOffersWithFiltres").on('click', function () {
  var formData = new FormData(document.getElementById("filtersOffersFrom"))
  var newTitleNumber = formData.get('taskType')
  console.log(newTitleNumber)
  if (newTitleNumber != "0") {
    $('#mainTitle').html("")
    $('#mainTitle').html($('#typeFirst' + newTitleNumber).html())
    newTitleNumber = formData.get('secondType')
    if (newTitleNumber != "0") {
      $('#mainTitle').html($('#mainTitle').html() + ' > ' + $('#typeSecond' + newTitleNumber).html())
    }
  }
  formData.set('offset', 0)
  var data = {};
  formData.forEach((value, key) => (data[key] = value));
  var json = JSON.stringify(data);
  fetch('/offers/sort', {
    method: "POST",
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: json
  }).then(function (response) { return response.json(); })
    .then(function (json) {
      $("#offerListPage").empty()
      addOffersToPage(json)
    });

});


$("#loadMoreOffers").on('click', function () {
  var formData = new FormData(document.getElementById("filtersOffersFrom"))
  let page = $('#currentPageOffer').val()
  let pageN = Number(page) + 1
  $('#currentPageOffer').val(pageN)
  formData.set('offset', pageN)
  var data = {};
  formData.forEach((value, key) => (data[key] = value));
  var json = JSON.stringify(data);
  fetch('/offers/sort', {
    method: "POST",
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: json
  }).then(function (response) { return response.json(); })
    .then(function (json) {
      addOffersToPage(json)
    });

});
function addOffersToPage(data) {
  for (let index = 0; index < data.length; ++index) {
    const element = data[index];
    let wt = element['WorkType']; let wtype = "";
    if (wt == 1) { wtype = "Обучение" } else if (wt == 2) { wtype = "Выполнение" } else { wtype = "Временный работник" }
    $("#offerListPage").append(`
      <div class="col-lg-4 col-12" style="margin-bottom: 20px;">
        <div
          style="width: 99%; margin-left: auto;margin-right: auto;border-radius: 10px; border: 1px solid; background-color: #1B1C2C;">
          <div style="width: 90%; margin-left: auto;margin-right: auto;">
            <img src="/static/upload/offerSkin/` + element['CoverPath'] + `"
              style="width: 100%; height: 200px; border-radius: 10px; display: block; margin-left: auto; margin-right: auto; margin-top: 10px">
            <h4 style="color: #fff; margin-top: 10px" onclick="showCurrentOffer(` + element['Id'] + `)">` + element['Name'] + `</h4>
            <h5 style="color: #fff; margin-top: 10px">Рейтинг - ` + element['Rating'] + ` (` + element['HistoryCount'] + `) ` + `</h5>
            <div class="row">
              <div class="col-6">
                <p style="color: #fff"><a href="` + element['FkUserOwner'] + `">` + element['UserOwnerName'] + `</a></p>
              </div>
              <div class="col-6">
                <p style="color: #fff; text-align: right;">Цена: ` + element['Price'] + `р</p>
              </div>
            </div>
            <p style="color: #fff">Срок: ` + element['DaysToComplite'] + ` дня, ` + wtype + `</p>
            <p style="color: #fff">` + element['Tags'] + `</p>
            <button
                style="color: #fff; background-color: #E746E0; height: 35px; width: 100%; border-radius: 10px; margin-bottom: 10px;"
                onclick="sendMessageOpen(` + element['Id'] + `)">Написать</button>
          </div>
        </div>
      </div>
    `)
  }
}

$("#showOrdersWithFiltres").on('click', function () {
  var formData = new FormData(document.getElementById("filtersOrdersFrom"))
  formData.set('offset', 0)
  var data = {};
  formData.forEach((value, key) => (data[key] = value));
  var json = JSON.stringify(data);
  fetch('/orders/sort', {
    method: "POST",
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: json
  }).then(function (response) { return response.json(); })
    .then(function (json) {
      $("#orderListPage").empty()
      addOrdersToPage(json)
    });

});

$("#loadMoreOrders").on('click', function () {
  var formData = new FormData(document.getElementById("filtersOrdersFrom"))
  let page = $('#currentPageOrder').val()
  let pageN = Number(page) + 1
  $('#currentPageOrder').val(pageN)
  formData.set('offset', pageN)
  var data = {};
  formData.forEach((value, key) => (data[key] = value));
  var json = JSON.stringify(data);
  fetch('/orders/sort', {
    method: "POST",
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: json
  }).then(function (response) { return response.json(); })
    .then(function (json) {
      addOrdersToPage(json)
    });

});
function addOrdersToPage(data) {
  for (let index = 0; index < data.length; ++index) {
    const element = data[index];
    let wt = element['WorkType']; let wtype = "";
    let sr = element['Urgency']; let srtype = "";
    if (wt == 1) { wtype = "Обучение" } else if (wt == 2) { wtype = "Выполнение" } else { wtype = "Временный работник" }
    if (sr == 1) { srtype = "Низкая" } else if (sr == 2) { srtype = "Средняя" } else { srtype = "Высокая" }
    $("#orderListPage").append(`
      <div class="col-sm-5 col-12"
          style="margin-left: auto; border-radius: 10px; border: 1px solid; background-color: #1B1C2C; margin-bottom: 20px;">
          <h4 style="color: #fff">` + element['Name'] + `</h4>
          <p style="color: #fff">` + element['Discribtion'] + `... <a onclick="showCurrentOrder(` + element['Id'] + `)">Читать далее</a></p>
          <p style="color: #fff">Желаемая цена до: ` + element['Price'] + `р</p>
          <p style="color: #fff">Выполнить до: ` + element['Deadline'] + `</p>
          <p style="color: #fff">Срочность: ` + srtype + `</p>
          <p style="color: #fff">Тип задачи: ` + wtype + `</p>
          <p style="color: #fff">Теги:</p>
          <p style="color: #fff">` + element['Tags'] + `</p>
          <div class="row">
            <div class="col-5">
              <button
                style="border-radius: 10px; border: 1px solid #E746E0; background: #1B1C2C; width: 100%; color: #fff; margin-bottom: 10px;">Жалоба</button>
            </div>
            <div class="col-5" style="margin-left: auto;">
              <button
                style="border-radius: 10px; background: #E746E0; width: 100%; color: #fff; margin-bottom: 10px;">Отозваться</button>
            </div>
          </div>
        </div>
    `)
  }
}

$("#offerCategoryFirstSelect").on('change', function () {
  let first = $("#offerCategoryFirstSelect").val()
  let htmlIn = '<option value="0">Выберите подкаттегорию</option>';
  for (var i = 0; i < secondTypes.length; i++) {
    if (secondTypes[i]["FkFirstType"] == first) {
      htmlIn += '<option value="' + secondTypes[i]["Id"] + '" id="typeSecond' + secondTypes[i]["Id"] + '">' + secondTypes[i]["Name"] + '</option>'
    }
  }
  $("#offerCategorySecondSelect").html(htmlIn)
});

$("#orderCategoryFirstSelect").on('change', function () {
  let first = $("#orderCategoryFirstSelect").val()
  let htmlIn = '<option value="0">Выберите подкаттегорию</option>';
  for (var i = 0; i < secondTypes.length; i++) {
    if (secondTypes[i]["FkFirstType"] == first) {
      htmlIn += '<option value="' + secondTypes[i]["Id"] + '">' + secondTypes[i]["Name"] + '</option>'
    }
  }
  $("#orderCategorySecondSelect").html(htmlIn)
});

function openFiltersOffer() {
  $('#sortingFieldsOffer').css("display", "block");
  $('#sortingFieldsButtonOffer').css("display", "none");
  $('#hiddenFiltresButtonOffer').css("visibility", "visible")
}
function closeFiltersOffer() {
  $('#sortingFieldsOffer').css("display", "none");
  $('#sortingFieldsButtonOffer').css("display", "block");
  $('#hiddenFiltresButtonOffer').css("visibility", "hidden")
}

function openFiltersOrder() {
  $('#sortingFieldsOrder').css("display", "block");
  $('#sortingFieldsButtonOrder').css("display", "none");
  $('#hiddenFiltresButtonOrder').css("visibility", "visible")
}
function closeFiltersOrder() {
  $('#sortingFieldsOrder').css("display", "none");
  $('#sortingFieldsButtonOrder').css("display", "block");
  $('#hiddenFiltresButtonOrder').css("visibility", "hidden")
}

function addOfferToOfferList(formData, lastId, fn) {
  $("#rowOfferList").append(`
      <div class="col-lg-5" id="offer`+ lastId + `"
            style="text-align: center; background-color: #1B1C2C; border-radius: 10px; margin-left: auto; margin-right: auto; margin-top: 20px">
            <img src="/static/upload/offerSkin/` + fn + `"
              style="width: 100%; height: 150px; border-radius: 10px; display: block; margin-left: auto; margin-right: auto; margin-top: 10px">
            <a href="/offers/` + lastId + `"><h4 style="color: #fff; margin-top: 10px">` + formData.get('offerName') + `</h4></a>
            <p style="color: #fff">Цена:` + formData.get('offerPrice') + `р</p>
            <div class="row">
              <div class="col-5">
                <p style="color: #fff">Активный?</p>
              </div>
              <div class="col-5">
              <input type="checkbox" id="switchOffer" checked onclick="changeOfferStatus(` + lastId + `)"/><label for="switchOffer" class="swithBall">Toggle</label>
              </div>
            </div>
            <div class="row">
              <div class="col-5">
                <button
                  style="color: #fff; background-color: #ff0000; height: 45px; width: 100%; border-radius: 10px; margin-bottom: 20px"
                  id="deleteOfferButton" type="button" onclick="deleteOffer(` + lastId + `)">Удалить</button>
              </div>
              <div class="col-5">
              <button
                  style="color: #fff; background-color: #E746E0; height: 45px; width: 100%; border-radius: 10px; margin-bottom: 20px"
                  id="deleteOfferButton" type="button" onclick="changeOfferGetOldData(` + lastId + `)">Изменить</button>
              </div>
            </div>
      </div>
    `)
}

function addOrderToOrderList(formData, lastId) {
  let tmpD = ""
  if (formData.get("orderDiscribction").length > 20) {
    tmpD = formData.get("orderDiscribction").slice(20)
  } else {
    tmpD = formData.get("orderDiscribction")
  }
  $("#rowOrderList").append(`
      <div class="col-lg-5" id="order` + lastId + `"
        style="text-align: center; background-color: #1B1C2C; border-radius: 10px; margin-left: auto; margin-right: auto; margin-top: 20px">
        <a href="/orders/` + lastId + `">
          <h4 style="color: #fff; margin-top: 10px">` + formData.get("orderName") + `</h4>
        </a>
        <p style="color: #fff">` + tmpD + `... <a href="/orders/` + lastId + `">Читать далее</a></p>
        <p style="color: #fff">Желаемая цена до:` + formData.get("orderPrice") + `р</p>
        <button
            style="color: #fff; background-color: #ff0000; height: 45px; width: 100%; border-radius: 10px; margin-bottom: 20px"
            id="deleteOfferButton" type="button" onclick="deleteOrder(` + lastId + `)">Удалить</button>
      </div>
    `)
}

function changeOfferStatus(offerId) {
  $.post('/account/changeActiveStatusForOffer', { value: offerId }, function () {
  }).done(function (response) {
    response.json().then((dataR) => {
      if (dataR["status"] == "notAuth") {
        window.location.href = 'login?target=account'
      }
    });
    var tmpResp = JSON.parse(response);
    showNotify(tmpResp["body"])
  });
}

function changeOfferGetOldData(offerId) {
  $.get('/account/getOfferInfo/' + offerId, function () {
  }).done(function (response) {
    var responseData = JSON.parse(response);
    if (responseData["status"] == "ok") {
      $("#changeOfferForm input[name=offerId]").val(offerId)
      $("#changeOfferForm input[name=offerName]").val(responseData["name"])
      $("#changeOfferForm textarea[name=offerDiscribtion]").val(responseData["discribtion"])
      $("#changeOfferForm input[name=offerDaysToComplite]").val(responseData["daysToComplite"])
      $("#changeOfferForm input[name=offerPrice]").val(responseData["price"])
      $("#changeOfferForm select[name=offerCategoryFirst]").val(responseData["orderCategory"])
      $("#changeOfferForm select[name=offerCategorySecond]").val(responseData["orderCategorySecond"])
      $("#changeOfferForm select[name=offerWorkType]").val(responseData["workType"])
      $("#changeOfferForm input[name=offerTags]").val(responseData["tags"])
      $('#changeOffer').modal('show')
    }
    if (responseData["status"] == "notAuth") {
      window.location.href = 'login?target=account'
    }
    showNotify(responseData["body"])
  });
}

function deleteOffer(offerId) {
  $.post('/account/deleteOffer', { value: offerId }, function () {
  }).done(function (response) {
    var tmpResp = JSON.parse(response);
    if (tmpResp["status"] == "ok") {
      needId = "#offer" + offerId
      $(needId).remove();
    }
    if (tmpResp["status"] == "notAuth") {
      window.location.href = 'login?target=account'
    }
    showNotify(tmpResp["body"])
  });
}

function deleteOrder(orderId) {
  $.post('/account/deleteOrder', { value: orderId }, function () {
  }).done(function (response) {
    var tmpResp = JSON.parse(response);
    if (tmpResp["status"] == "ok") {
      needId = "#order" + orderId
      $(needId).remove();
    }
    if (tmpResp["status"] == "notAuth") {
      window.location.href = 'login?target=account'
    }
    showNotify(tmpResp["body"])
  });
}

function showCurrentOffer(offerId) {
  $.get('/offers/' + offerId, function () {
  }).done(function (response) {
    var responseData = JSON.parse(response);
    $("#offerClickId").val(offerId)
    let wt = responseData["workType"]; let wtype = "";
    if (wt == 1) { wtype = "Обучение" } else if (wt == 2) { wtype = "Выполнение" } else { wtype = "Временный работник" }
    $("#offerClickName").html(responseData["name"])
    $("#offerClickWorkType").html('Тип: ' + wtype)
    $("#offerClickDiscribtion").html(responseData["discribtion"])
    $("#offerClickDaysToCompite").html('Срок выполнения: ' + responseData["daysToComplite"] + ' дней')
    $("#offerClickPrice").html('Цена: ' + responseData["price"])
    $("#offerClickUser").html(responseData["userOwnerName"])
    $("#offerClickUser").attr('href', 'user/' + responseData["id"])
    $("#offerClickImg").attr('src', '/static/upload/offerSkin/' + responseData["coverPath"])
    $("#offerClickTags").html('Теги: ' + responseData["tags"])
    $('#offerClick').modal('show')
    $('#userOwnerId').val(offerId)
  });
}

function showCurrentOrder(offerId) {
  $.get('/orders/' + offerId, function () {
  }).done(function (response) {
    var responseData = JSON.parse(response);
    //$("#changeOfferForm input[name=offerId]").val(offerId)
    let wt = responseData["workType"]; let wtype = "";
    let sr = responseData['urgency']; let srtype = "";
    if (wt == 1) { wtype = "Обучение" } else if (wt == 2) { wtype = "Выполнение" } else { wtype = "Временный работник" }
    if (sr == 1) { srtype = "Низкая" } else if (sr == 2) { srtype = "Средняя" } else { srtype = "Высокая" }
    $("#orderClickName").html(responseData["name"])
    $("#orderClickWorkType").html('Тип: ' + wtype)
    $("#orderClickDiscribtion").html(responseData["discribtion"])
    $("#orderClickDeadline").html('Срок выполнения: ' + responseData["deadline"] + ' дней')
    $("#orderClickPrice").html('Цена: ' + responseData["price"])
    $("#orderClickPrice").html('Срочность : ' + srtype)
    $("#orderClickUserOwner").html(responseData["fkUserOwnerName"])
    $("#orderClickUserOwner").attr('href', 'user/' + responseData["fkUserOwner"])
    if (responseData["tzPath"] != "") {
      $("#orderClickTZ").show()
      $("#orderClickTZ").attr('src', '/static/upload/tz/' + responseData["tzPath"])
    } else {
      $("#orderClickTZ").css("visibility", "hidden")
    }
    $("#orderClickTags").html('Теги: ' + responseData["tags"])
    $('#orderClick').modal('show')
  });
}

function sendMessageOpen(offerId, ownerId) {
  $('#offerClickMessageId').val(offerId)
  $('#userOwnerId').val(ownerId)
  $('#offerClickMessage').modal('show')
}
function sendMessageOpenFromModal() {
  $('#offerClick').modal('hide')
  $('#offerClickMessageId').val($("#offerClickId").val())
  $('#offerClickMessage').modal('show')
}

function sendMessageSend() {
  $('#offerClickMessage').modal('hide')
  let offerId = $('#offerClickMessageId').val()
  let message = $('#offerClickMessageTextAria').val()
  let userOwnerId = $('#userOwnerId').val()
  $.post('/offers/' + offerId + '/message', { value: message, ownerId: userOwnerId }, function () {
  }).done(function (response) {
    var tmpResp = response;
    if (tmpResp["status"] == "notAuth") {
      window.location.href = 'login?target=offers'
      return
    }
    showNotify(tmpResp["body"])
  });
}

function showNotify(text) {
  $("#classicNotify").append('<div class="alert"><span class="closebtn" >×</span><strong>Уведомление!</strong>' + text + '</div>')
  var close = document.getElementsByClassName("closebtn");
  var i;

  for (i = 0; i < close.length; i++) {
    close[i].onclick = function () {
      var div = this.parentElement;
      div.style.opacity = "0";
      setTimeout(function () {
        div.style.display = "none";
      }, 6);
    }
  }
}

function getAdress() {
  var location = JSON.stringify(window.location.href).split('//')[1]
  location = location.split('/')[0]
  return location
}

$('body').append('<div class="upbtn"></div>');
$(window).scroll(function () {
  if ($(this).scrollTop() > 200) {
    $('.upbtn').css({
      bottom: '30px'
    });
  } else {
    $('.upbtn').css({
      bottom: '-80px'
    });
  }
});
$('.upbtn').on('click', function () {
  $('html, body').animate({
    scrollTop: 0
  }, 0);
  return false;
});






