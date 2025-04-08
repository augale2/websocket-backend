package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"websocket-backend/services/auth-service/internal/db"
	pb "websocket-backend/services/auth-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"

	"net/http"
)

var jwtSecret = []byte("abhishekisking")

func generateJWT(userID, username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

type authServer struct {
	pb.UnimplementedAuthServiceServer
	userDB *db.UserDBClient
}

func newAuthServer() *authServer {

	ctx := context.Background()

	userDB, err := db.NewUserDBClient(ctx, "Users")

	if err != nil {
		log.Fatalf("Failed to create UserDB client: %v", err)
	}

	return &authServer{
		userDB: userDB,
	}
}

func (s *authServer) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {

	_, err := s.userDB.GetUserByUsername(ctx, req.GetUsername())
	if err == nil {
		// If a user is found, enforce uniqueness by returning an error.
		return nil, fmt.Errorf("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	userID, err := s.userDB.CreateUser(ctx, req.GetUsername(), string(hashedPassword))

	if err != nil {
		return nil, fmt.Errorf("error registering user: %v", err)
	}

	log.Printf("Registered user: %s with ID: %s", req.GetUsername(), userID)

	return &pb.RegisterUserResponse{
		UserId:  userID,
		Message: "User registered successfully",
	}, nil

}

func (s *authServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	userItem, err := s.userDB.GetUserByUsername(ctx, req.GetUsername())

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	storedHash := userItem["Password"].(*types.AttributeValueMemberS).Value
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.GetPassword()))

	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	userID := userItem["UserId"].(*types.AttributeValueMemberS).Value
	token, err := generateJWT(userID, req.GetUsername())

	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)

	}
	log.Printf("User %s logged in with token: %s", req.GetUsername(), token)
	return &pb.LoginUserResponse{
		Token:   token,
		Message: "User logged in successfully",
	}, nil
}

func (s *authServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := validateJWT(req.GetToken())

	if err != nil {
		log.Printf("Token validation failed: %v", err)
		return &pb.ValidateTokenResponse{
			Valid:   false,
			Message: fmt.Sprintf("Invalid token: %v", err),
		}, nil
	}
	log.Printf("Token is valid. Claims: %v", claims)

	return &pb.ValidateTokenResponse{
		Valid:   true,
		Message: "Token is valid",
	}, nil
}

// HTTP handler functions
func (s *authServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	resp, err := s.RegisterUser(r.Context(), &pb.RegisterUserRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Registration failed",
			"details": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (s *authServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	resp, err := s.LoginUser(r.Context(), &pb.LoginUserRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Login failed",
			"details": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (s *authServer) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	resp, err := s.ValidateToken(r.Context(), &pb.ValidateTokenRequest{
		Token: req.Token,
	})
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Token validation failed",
			"details": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	// Create the gRPC server
	s := grpc.NewServer()
	authServer := newAuthServer()
	pb.RegisterAuthServiceServer(s, authServer)
	reflection.Register(s)

	// Create HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/v1/auth/register", authServer.handleRegister).Methods("POST", "OPTIONS")
	router.HandleFunc("/v1/auth/login", authServer.handleLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/v1/auth/validate", authServer.handleValidateToken).Methods("POST", "OPTIONS")

	// Add CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	wrappedGrpc := grpcweb.WrapServer(s,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			return true
		}),
	)

	httpServer := &http.Server{
		Addr: ":50052",
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			resp.Header().Set("Access-Control-Allow-Origin", "*")
			resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")

			if req.Method == "OPTIONS" {
				resp.WriteHeader(http.StatusOK)
				return
			}

			if wrappedGrpc.IsGrpcWebRequest(req) || wrappedGrpc.IsAcceptableGrpcCorsRequest(req) {
				wrappedGrpc.ServeHTTP(resp, req)
				return
			}

			corsMiddleware(router).ServeHTTP(resp, req)
		}),
	}

	log.Println("Auth Service server is running on port 50052")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
