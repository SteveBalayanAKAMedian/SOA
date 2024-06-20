package user

import (
    "context"
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "time"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	//"github.com/google/uuid"
	//"github.com/gorilla/mux"

	"google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "task-service/proto"
	"github.com/gorilla/mux"
)

type Handler struct {
    db         *sql.DB
    grpcClient proto.TaskServiceClient
}

func NewHandler(db *sql.DB, grpcAddress string) (*Handler, error) {
    var opts []grpc.DialOption

    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

    conn, err := grpc.NewClient(grpcAddress, opts...)
    if err != nil {
        log.Fatalf("fail to dial: %v", err)
        return nil, err
    }

    client := proto.NewTaskServiceClient(conn)
    return &Handler{db: db, grpcClient: client}, nil
}

// CreateTask handles creating a new task
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
    tmp := r.Context().Value("userID")
	userID := tmp.(int)
	userIDStr := strconv.Itoa(userID)

    var req struct {
        Title       string `json:"title"`
        Description string `json:"description"`
    }
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    createReq := &proto.CreateTaskRequest{
        UserId:      userIDStr,
        Title:       req.Title,
        Description: req.Description,
    }

    createRes, err := h.grpcClient.CreateTask(context.Background(), createReq)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(createRes.Task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    tmp := r.Context().Value("userID")
	userID := tmp.(int)
	userIDStr := strconv.Itoa(userID)

    var req struct {
        Id          string `json:"id"`
        Title       string `json:"title"`
        Description string `json:"description"`
        Status      string `json:"status"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    grpcReq := &proto.UpdateTaskRequest{
        Id:          req.Id,
        UserId:      userIDStr,
        Title:       req.Title,
        Description: req.Description,
        Status:      req.Status,
    }

    grpcResp, err := h.grpcClient.UpdateTask(context.Background(), grpcReq)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(grpcResp.Task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    tmp := r.Context().Value("userID")
	userID := tmp.(int)
	userIDStr := strconv.Itoa(userID)

    taskID := r.URL.Query().Get("id")

    grpcReq := &proto.DeleteTaskRequest{
        Id:     taskID,
        UserId: userIDStr,
    }

    grpcResp, err := h.grpcClient.DeleteTask(context.Background(), grpcReq)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

   json.NewEncoder(w).Encode(grpcResp)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    taskID := vars["id"]
    log.Printf("Task ID from URL: %s", taskID)

    if taskID == "" {
        http.Error(w, "Task ID is required", http.StatusBadRequest)
        return
    }

    grpcReq := &proto.GetTaskRequest{Id: taskID}

    grpcResp, err := h.grpcClient.GetTask(context.Background(), grpcReq)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(grpcResp.Task)
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

    grpcReq := &proto.ListTasksRequest{
        Page:     int32(page),
        PageSize: int32(pageSize),
    }

    grpcResp, err := h.grpcClient.ListTasks(context.Background(), grpcReq)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(grpcResp.Tasks)
}



func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	if err := user.Create(h.db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Генерация токена
	token, expiresAt, err := generateToken(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохранение сессии в базу данных
	if err := CreateSession(h.db, user.ID, token, expiresAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: expiresAt,
	})

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    username, ok := r.Context().Value("username").(string)
    if !ok {
        http.Error(w, "Missing username in header", http.StatusUnauthorized)
        return
    }

    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    user.Username = username

    if err := user.Update(h.db); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := Authenticate(h.db, credentials.Username, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, expiresAt, err := generateToken(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохранение сессии в базу данных
	if err := CreateSession(h.db, user.ID, token, expiresAt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: expiresAt,
	})

	w.WriteHeader(http.StatusOK)
}

func generateToken(user *User) (string, time.Time, error) {
	signingKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	expiresAt := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
