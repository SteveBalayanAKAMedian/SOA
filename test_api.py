import requests
import json

BASE_URL = "http://localhost:8080"

def register_user(username, password, first_name, last_name, email, phone_number, birthday):
    url = f"{BASE_URL}/register"
    payload = {
        "username": username,
        "password": password,
        "first_name": first_name,
        "last_name": last_name,
        "email": email,
        "phone_number": phone_number,
        "birthday": birthday
    }
    response = requests.post(url, headers={"Content-Type": "application/json"}, data=json.dumps(payload))
    print("Register User Response:", response.status_code, response.text)

def authenticate_user(username, password):
    url = f"{BASE_URL}/auth"
    payload = {
        "username": username,
        "password": password
    }
    response = requests.post(url, headers={"Content-Type": "application/json"}, data=json.dumps(payload))
    print("Authenticate User Response:", response.status_code, response.text)
    if response.status_code == 200:
        return response.cookies.get('session_token')
    return None

def update_user(token, first_name, last_name, email, phone_number, birthday):
    url = f"{BASE_URL}/update"
    payload = {
        "first_name": first_name,
        "last_name": last_name,
        "email": email,
        "phone_number": phone_number,
        "birthday": birthday
    }
    cookies = {"session_token": token}
    response = requests.put(url, headers={"Content-Type": "application/json"}, cookies=cookies, data=json.dumps(payload))
    print("Update User Response:", response.status_code, response.text)

def main():
    username = "testuser"
    password = "testpassword"
    first_name = "Test"
    last_name = "User"
    email = "test@example.com"
    phone_number = "1234567890"
    birthday = "1990-01-01"

    register_user(username, password, first_name, last_name, email, phone_number, birthday)

    token = authenticate_user(username, password)
    if token:
        print("Received token:", token)

        update_user(token, "UpdatedName", "UpdatedLastName", "updated@example.com", "0987654321", "1992-02-02")

    register_user(username, password, first_name, last_name, email, phone_number, birthday)

    bad_token = "badtoken"
    update_user(bad_token, "BadUpdate", "BadUser", "bad@example.com", "1111111111", "2000-01-01")

if __name__ == "__main__":
    main()
