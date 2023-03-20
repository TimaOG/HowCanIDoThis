import requests
#pytest.exe -s -v tests.py
s = requests.Session()

def test_get_site():
    response = s.get("http://localhost:8080")
    assert response.status_code == 200

def test_login():
    response = s.post("http://localhost:8080/login", data={'email':'bookgun@mail.ru', 'password': "123"})
    assert response.status_code == 200

def test_openAccount():
    response = s.get("http://localhost:8080/account")
    assert response.status_code == 200

def test_changeActiveStatus():
    response = s.post("http://localhost:8080/account/changeActiveStatus")
    assert response.status_code == 200 and response.json()["status"] == 'ok'
    response = s.post("http://localhost:8080/account/changeActiveStatus")
    assert response.status_code == 200 and response.json()["status"] == 'ok'

def test_changeProfileImg():
    files = {'file': open('wTestImg/1.jpg','rb')}
    response = s.post("http://localhost:8080/account/changeProfileImg", files=files)
    assert response.status_code == 200 and response.json()["status"] == 'ok'

def test_changeAccountSettings():
    response = s.post("http://localhost:8080/account/changeSettings", data={'newEmail': 'bookgun1@mail.ru', 'oldPassword': "123", 'newPassword': "", 'newPassword2': ""})
    assert response.status_code == 200 and response.json()["status"] == 'ok'
    response = s.post("http://localhost:8080/account/changeSettings", data={'newEmail': 'bookgun1@mail.ru', 'oldPassword': "123", 'newPassword': "111", 'newPassword2': "111"})
    assert response.status_code == 200 and response.json()["status"] == 'ok'
    response = s.post("http://localhost:8080/account/changeSettings", data={'newEmail': 'bookgun@mail.ru', 'oldPassword': "111", 'newPassword': "123", 'newPassword2': "123"})
    assert response.status_code == 200 and response.json()["status"] == 'ok'

def test_createAndDeleteOrder():
    response = s.post("http://localhost:8080/account/newOrder", files={'fileTZ': None}, data={'orderName': 'Доработать сайт', 'orderDiscribction': "Есть сайт вот ссылка ЭТО ССЫЛКА, у меня не работает то то то то, надо исправить", 
    'orderDeadline': "2025-04-01", 'orderPrice': "5000", 'orderCategoryFirst': "1", 'orderCategorySecond': "7", 'orderUrgency': '3', 'orderWorkType': '2', 'orderTags': 'lol;kek;suka'})
    assert response.status_code == 200 and response.json()["status"] == 'ok'
    lastOrderId = response.json()["lastId"]
    response = s.post("http://localhost:8080/account/deleteOrder", data={'value': lastOrderId})
    assert response.status_code == 200 and response.json()["status"] == 'ok'

def test_createChangeAndDeleteOffer():
    response = s.post("http://localhost:8080/account/newOffer", files={'cover': open("wTestImg/1.jpg", "rb")}, data={'offerName': 'Помогу доработать сайт', 'offerDiscribtion': "Есть сайт вот ссылка скинь ССЫЛКу, у меня не работает то то то то, надо исправить", 
    'offerDaysToComplite': "5", 'offerPrice': "5000", 'offerCategoryFirst': "1", 'offerCategorySecond': "7", 'offerWorkType': '2', 'offerTags': 'lol;kek;suka'})
    assert response.status_code == 200 and response.json()["status"] == 'ok'
    lastOfferId = response.json()["lastId"]
    response = s.post("http://localhost:8080/account/changeOffer", files={'cover': None}, data={'offerId': lastOfferId,'offerName': 'Помогу доработать сайт дебил', 'offerDiscribtion': "Есть сайт вот ссылка скинь ССЫЛКу, у меня не работает то то то то, надо исправить", 
    'offerDaysToComplite': "5", 'offerPrice': "5000", 'offerCategoryFirst': "1", 'offerCategorySecond': "7", 'offerWorkType': '2', 'offerTags': 'lol;kek;suka'})
    assert response.status_code == 200 and response.json()["status"] == 'ok'
    response = s.post("http://localhost:8080/account/deleteOffer", data={'value': lastOfferId})
    assert response.status_code == 200 and response.json()["status"] == 'ok'