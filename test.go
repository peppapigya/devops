package k8s_platform_go

//// README: K8s 无限集群管理平台 - Go 语言骨架（示例）
//// 技术栈：Go 1.21, Gin, GORM, client-go
//// 目录结构示例：
////  - cmd/server/main.go
////  - internal/config/config.go
////  - internal/db/db.go
////  - internal/model/models.go
////  - internal/handler/clusters.go
////  - internal/middleware/jwt.go
////  - internal/middleware/rbac.go
////  - internal/k8s/client.go
////  - go.mod
//
//// file: go.mod
////module github.com/example/k8s-platform
//
////go 1.21
////
////require (
////github.com/gin-gonic/gin v1.9.0
////gorm.io/driver/mysql v1.5.1
////gorm.io/gorm v1.26.0
////github.com/golang-jwt/jwt/v5 v5.1.0
////k8s.io/client-go v0.29.13
////)
//
//// file: cmd/server/main.go
//package main
//
//import (
//	"fmt"
//	"log"
//	"os"
//
//	"github.com/example/k8s-platform/internal/config"
//	"github.com/example/k8s-platform/internal/db"
//	"github.com/example/k8s-platform/internal/handler"
//	"github.com/example/k8s-platform/internal/middleware"
//
//	"github.com/gin-gonic/gin"
//)
//
//func main() {
//	cfg := config.Load()
//	// 初始化数据库
//	dbConn, err := db.NewMySQL(cfg.DatabaseDSN)
//	if err != nil {
//		log.Fatalf("db init: %v", err)
//	}
//	defer func() { _ = dbConn.DB() }()
//
//	// 自动迁移（建议仅在开发/CI阶段）
//	dbConn.AutoMigrate()
//
//	r := gin.Default()
//
//	// 中间件
//	r.Use(middleware.RequestLogger())
//	r.Use(middleware.Recover())
//
//	api := r.Group("/api/v1")
//	api.POST("/auth/login", handler.LoginHandler)
//	// 受保护路由
//	auth := api.Group("")
//	auth.Use(middleware.JWTAuth())
//	// 集群管理
//	auth.POST("/clusters", handler.CreateCluster)
//	auth.GET("/clusters", handler.ListClusters)
//	auth.GET("/clusters/:id/health", handler.ClusterHealth)
//	auth.POST("/clusters/:id/toggle", handler.ToggleCluster)
//	// 更多接口按需添加...
//
//	addr := fmt.Sprintf(":%d", cfg.Port)
//	log.Printf("listening %s", addr)
//	if err := r.Run(addr); err != nil {
//		log.Fatalf("server run: %v", err)
//	}
//}
//
//// file: internal/config/config.go
//package config
//
//import (
//"os"
//"strconv"
//)
//type Config struct {
//	Port        int
//	DatabaseDSN string
//	JWTSecret   string
//}
//
//func Load() *Config {
//	port := 8080
//	if p := os.Getenv("PORT"); p != "" {
//		if v, err := strconv.Atoi(p); err == nil {
//			port = v
//		}
//	}
//	dsn := os.Getenv("DATABASE_DSN")
//	if dsn == "" {
//		dsn = "user:pass@tcp(127.0.0.1:3306)/k8sdb?parseTime=true&loc=Local"
//	}
//	jwt := os.Getenv("JWT_SECRET")
//	if jwt == "" {
//		jwt = "replace-with-secure-secret"
//	}
//	return &Config{Port: port, DatabaseDSN: dsn, JWTSecret: jwt}
//}
//
//// file: internal/db/db.go
//package db
//
//import (
//"gorm.io/driver/mysql"
//"gorm.io/gorm"
//"log"
//"time"
//
//"github.com/example/k8s-platform/internal/model"
//)
//
//var DB *gorm.DB
//
//func NewMySQL(dsn string) (*gorm.DB, error) {
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		return nil, err
//	}
//	sqlDB, _ := DB.DB()
//	sqlDB.SetMaxIdleConns(10)
//	sqlDB.SetMaxOpenConns(100)
//	sqlDB.SetConnMaxLifetime(time.Hour)
//	return DB, nil
//}
//
//func AutoMigrate() {
//	if DB == nil {
//		log.Fatal("db not initialized")
//	}
//	// 仅示例表，按需增加
//	if err := DB.AutoMigrate(
//		&model.User{}, &model.Org{}, &model.Project{}, &model.Role{},
//		&model.K8sCluster{}, &model.K8sClusterCredential{}, &model.K8sClusterStatus{},
//	); err != nil {
//		log.Fatalf("auto migrate: %v", err)
//	}
//}
//
//// file: internal/model/models.go
//package model
//
//import (
//"time"
//
//"gorm.io/datatypes"
//)
//
//// 用户与组织（简化字段）
//type User struct {
//	ID           uint   `gorm:"primaryKey"`
//	Username     string `gorm:"size:64;uniqueIndex"`
//	DisplayName  string `gorm:"size:128"`
//	Email        string `gorm:"size:128"`
//	PasswordHash string `gorm:"size:255"`
//	Status       int    `gorm:"default:1"`
//	CreatedAt    time.Time
//	UpdatedAt    time.Time
//}
//
//type Org struct {
//	ID        uint   `gorm:"primaryKey"`
//	Name      string `gorm:"size:128;uniqueIndex"`
//	Code      string `gorm:"size:64;uniqueIndex"`
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}
//
//type Project struct {
//	ID          uint `gorm:"primaryKey"`
//	OrgID       uint
//	Name        string `gorm:"size:128"`
//	Code        string `gorm:"size:64"`
//	Description string
//	CreatedAt   time.Time
//	UpdatedAt   time.Time
//}
//
//type Role struct {
//	ID          uint   `gorm:"primaryKey"`
//	Name        string `gorm:"size:64;uniqueIndex"`
//	DisplayName string
//	Level       int `gorm:"default:100"`
//	CreatedAt   time.Time
//	UpdatedAt   time.Time
//}
//
//// 集群相关
//type K8sCluster struct {
//	ID          uint `gorm:"primaryKey"`
//	OrgID       uint
//	ProjectID   *uint
//	Name        string `gorm:"size:128"`
//	Code        string `gorm:"size:64;index"`
//	ApiServer   string `gorm:"size:256"`
//	KubeVersion string `gorm:"size:64"`
//	Env         string `gorm:"size:32;default:prod"`
//	Region      string `gorm:"size:64"`
//	Labels      datatypes.JSON
//	Status      string `gorm:"size:32;default:Unknown"`
//	Enabled     bool   `gorm:"default:true"`
//	CreatedAt   time.Time
//	UpdatedAt   time.Time
//}
//
//type K8sClusterCredential struct {
//	ID         uint `gorm:"primaryKey"`
//	ClusterID  uint
//	Type       string `gorm:"size:32"`
//	ContentEnc string `gorm:"type:text"`
//	CreatedBy  uint
//	CreatedAt  time.Time
//}
//
//type K8sClusterStatus struct {
//	ID         uint `gorm:"primaryKey"`
//	ClusterID  uint
//	IsReady    bool
//	NodeCount  int
//	CpuTotal   float64
//	MemTotalGb float64
//	LastBeat   time.Time
//}
//
//// file: internal/middleware/jwt.go
//package middleware
//
//import (
//"net/http"
//"strings"
//"time"
//
//"github.com/gin-gonic/gin"
//"github.com/golang-jwt/jwt/v5"
//"github.com/example/k8s-platform/internal/config"
//)
//
//func JWTAuth() gin.HandlerFunc {
//	cfg := config.Load()
//	return func(c *gin.Context) {
//		auth := c.GetHeader("Authorization")
//		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "missing token"})
//			return
//		}
//		token := strings.TrimPrefix(auth, "Bearer ")
//		parsed, err := jwt.Parse(token, func(t *jwt.AccessToken) (interface{}, error) {
//			// HS256
//			return []byte(cfg.JWTSecret), nil
//		})
//		if err != nil || !parsed.Valid {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "invalid token"})
//			return
//		}
//		claims := parsed.Claims.(jwt.MapClaims)
//		c.Set("user", claims)
//		c.Next()
//	}
//}
//
//// helper: generate token (for demo)
//func GenerateToken(username string) (string, error) {
//	cfg := config.Load()
//	claims := jwt.MapClaims{
//		"sub": username,
//		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
//	}
//	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	return t.SignedString([]byte(cfg.JWTSecret))
//}
//
//// file: internal/middleware/rbac.go
//package middleware
//
//import (
//"github.com/gin-gonic/gin"
//)
//
//// Simple RBAC placeholder — integrate with DB role/permission checks
//func RequireRole(role string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 从上下文取 user claims,查库判断是否有权限
//		// 这里为占位符，生产需实现具体校验逻辑
//		c.Next()
//	}
//}
//
//// file: internal/handler/clusters.go
//package handler
//
//import (
//"encoding/base64"
//"net/http"
//"strconv"
//
//"github.com/gin-gonic/gin"
//"github.com/example/k8s-platform/internal/db"
//"github.com/example/k8s-platform/internal/model"
//"github.com/example/k8s-platform/internal/k8s"
//)
//
//// CreateCluster - 简化实现
//func CreateCluster(c *gin.Context) {
//	var req struct {
//		OrgID      uint              `json:"orgId" binding:"required"`
//		ProjectID  *uint             `json:"projectId"`
//		Name       string            `json:"name" binding:"required"`
//		Code       string            `json:"code" binding:"required"`
//		ApiServer  string            `json:"apiServer" binding:"required"`
//		Env        string            `json:"env"`
//		Region     string            `json:"region"`
//		Labels     map[string]string `json:"labels"`
//		Credential struct {
//			Type    string `json:"type"`
//			Content string `json:"content"`
//		} `json:"credential"`
//	}
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
//		return
//	}
//	cluster := model.K8sCluster{
//		OrgID: req.OrgID, ProjectID: req.ProjectID, Name: req.Name,
//		Code: req.Code, ApiServer: req.ApiServer, Env: req.Env, Region: req.Region,
//	}
//	if err := db.DB.Create(&cluster).Error; err != nil {
//		c.JSON(500, gin.H{"code": 500, "msg": "create failed"})
//		return
//	}
//	// 存凭据（演示：base64 存储，生产请加密）
//	if req.Credential.Content != "" {
//		enc := base64.StdEncoding.EncodeToString([]byte(req.Credential.Content))
//		cred := model.K8sClusterCredential{ClusterID: cluster.ID, Type: req.Credential.Type, ContentEnc: enc, CreatedBy: 0}
//		db.DB.Create(&cred)
//	}
//	c.JSON(200, gin.H{"code": 0, "msg": "OK", "data": cluster})
//}
//
//func ListClusters(c *gin.Context) {
//	var clusters []model.K8sCluster
//	q := db.DB
//	if org := c.Query("orgId"); org != "" {
//		if v, err := strconv.ParseUint(org, 10, 64); err == nil {
//			q = q.Where("org_id = ?", v)
//		}
//	}
//	q.Find(&clusters)
//	c.JSON(200, gin.H{"code": 0, "data": clusters})
//}
//
//func ClusterHealth(c *gin.Context) {
//	id := c.Param("id")
//	var status model.K8sClusterStatus
//	if err := db.DB.Where("cluster_id = ?", id).Order("last_beat desc").First(&status).Error; err == nil {
//		c.JSON(200, gin.H{"code": 0, "data": status})
//		return
//	}
//	c.JSON(404, gin.H{"code": 404, "msg": "not found"})
//}
//
//func ToggleCluster(c *gin.Context) {
//	id := c.Param("id")
//	var req struct {
//		Enabled bool `json:"enabled"`
//	}
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(400, gin.H{"code": 400, "msg": "bad request"})
//		return
//	}
//	if err := db.DB.Model(&model.K8sCluster{}).Where("id = ?", id).Update("enabled", req.Enabled).Error; err != nil {
//		c.JSON(500, gin.H{"code": 500, "msg": "update failed"})
//		return
//	}
//	c.JSON(200, gin.H{"code": 0, "msg": "OK"})
//}
//
//// file: internal/k8s/client.go
//package k8s
//
//import (
//"context"
//"encoding/base64"
//"fmt"
//
//"k8s.io/client-go/kubernetes"
//"k8s.io/client-go/rest"
//"k8s.io/client-go/tools/clientcmd"
//)
//
//// BuildClientFromKubeconfigBase64 - 从 base64(kubeconfig) 构建 clientset
//func BuildClientFromKubeconfigBase64(b64 string) (*kubernetes.Clientset, error) {
//	data, err := base64.StdEncoding.DecodeString(b64)
//	if err != nil {
//		return nil, err
//	}
//	config, err := clientcmd.RESTConfigFromKubeConfig(data)
//	if err != nil {
//		// 兼容：尝试用 kubeconfig loader
//		kcfg, err2 := clientcmd.NewClientConfigFromBytes(data)
//		if err2 != nil {
//			return nil, fmt.Errorf("kubeconfig parse failed: %v/%v", err, err2)
//		}
//		rc, err3 := kcfg.ClientConfig()
//		if err3 != nil {
//			return nil, err3
//		}
//		config = rc
//	}
//	// 如果还是 nil，尝试 InClusterConfig
//	if config == nil {
//		config, err = rest.InClusterConfig()
//		if err != nil {
//			return nil, err
//		}
//	}
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		return nil, err
//	}
//	return clientset, nil
//}
//
//// Example: 获取节点数（简单示例）
//func NodeCountFromClientset(clientset *kubernetes.Clientset) (int, error) {
//	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
//	if err != nil {
//		return 0, err
//	}
//	return len(nodes.Items), nil
//}
//
//// Note: 需要导入 metav1
//
//// file: internal/handler/auth.go
//package handler
//
//import (
//"net/http"
//
//"github.com/gin-gonic/gin"
//"github.com/example/k8s-platform/internal/middleware"
//)
//
//func LoginHandler(c *gin.Context) {
//	var req struct {
//		Username string `json:"username"`
//		Password string `json:"password"`
//	}
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(400, gin.H{"code": 400, "msg": "bad request"})
//		return
//	}
//	// TODO: 验证用户名密码（示例不做真实校验）
//	if req.Username == "admin" && req.Password == "admin" {
//		t, _ := middleware.GenerateToken(req.Username)
//		c.JSON(200, gin.H{"code": 0, "data": gin.H{"token": t}})
//		return
//	}
//	c.JSON(401, gin.H{"code": 401, "msg": "invalid credentials"})
//}
//
//// End of generated starter code.
//// 说明：
//// - 本骨架为最小可运行示例，着重于项目结构、模型、数据库自动迁移、认证中间件、集群 CRUD 与 Kubernetes 客户端接入点。
//// - 生产实战需要增加：配置加密（KMS）、错误处理、输入验证、权限校验、审计日志、测试用例、CI/CD、容器化部署等。
//// - 如需我：
////   1) 生成完整 CRUD 的更多模块（DNS/Registry/Git/Release/Events/Users）
////   2) 生成对应的 OpenAPI YAML 与 Swagger UI 集成
////   3) 把 DDL 导出为 schema.sql
//// 我可以基于当前骨架继续生成。
