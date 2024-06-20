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

def create_task(token, title, description):
    url = f"{BASE_URL}/task"
    payload = {
        "title": title,
        "description": description,
        "status": "open"
    }
    cookies = {"session_token": token}
    response = requests.post(url, headers={"Content-Type": "application/json"}, cookies=cookies, data=json.dumps(payload))
    print("Create Task Response:", response.status_code, response.text)
    return response

def update_task(token, task_id, title, description, status):
    url = f"{BASE_URL}/task"
    payload = {
        "id": task_id,
        "title": title,
        "description": description,
        "status": status
    }
    cookies = {"session_token": token}
    response = requests.put(url, headers={"Content-Type": "application/json"}, cookies=cookies, data=json.dumps(payload))
    print("Update Task Response:", response.status_code, response.text)
    return response

def delete_task(token, task_id):
    url = f"{BASE_URL}/task"
    params = {"id": task_id}
    cookies = {"session_token": token}
    response = requests.delete(url, headers={"Content-Type": "application/json"}, cookies=cookies, params=params)
    print("Delete Task Response:", response.status_code, response.text)
    return response

def get_task(token, task_id):
    url = f"{BASE_URL}/task/{task_id}"
    cookies = {"session_token": token}
    response = requests.get(url, headers={"Content-Type": "application/json"}, cookies=cookies)
    print("Get Task Response:", response.status_code, response.text)
    return response

def list_tasks(token, page, page_size):
    url = f"{BASE_URL}/tasks"
    params = {"page": page, "pageSize": page_size}
    cookies = {"session_token": token}
    response = requests.get(url, headers={"Content-Type": "application/json"}, cookies=cookies, params=params)
    print("List Tasks Response:", response.status_code, response.text)
    return response

def main():
    username = "testuser22"
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

        # Create a task
        create_response = create_task(token, "Test Task", "This is a test task")
        task_id = json.loads(create_response.text)["id"]
        print(task_id)

        # Update the task
        update_response = update_task(token, task_id, "Updated Task", "Updated description", "completed")

        # Get the task by ID
        get_response = get_task(token, task_id)

        # List tasks with pagination
        list_response = list_tasks(token, page=1, page_size=10)

        # Delete the task
        delete_response = delete_task(token, task_id)

    print("\nTrying with a bad token:")
    bad_token = "badtoken"
    update_user(bad_token, "BadUpdate", "BadUser", "bad@example.com", "1111111111", "2000-01-01")
    create_response = create_task(bad_token, "Bad Task", "This task should fail")
    get_response = get_task(bad_token, task_id)
    delete_response = delete_task(bad_token, task_id)
    list_response = list_tasks(bad_token, page=1, page_size=10)

    print("\nTrying to create task without auth:")
    create_response = requests.post(f"{BASE_URL}/task", headers={"Content-Type": "application/json"}, data=json.dumps({
        "title": "Unauth Task",
        "description": "This task should fail due to no auth"
    }))
    print("Create Task Response:", create_response.status_code, create_response.text)

if __name__ == "__main__":
    main()
