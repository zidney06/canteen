package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"server/src/config"
	"server/src/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

type MyClaims struct {
	ID         string `json:"id"`
	ClientName string `json:"username"`
	ClientKey  string `json:"clientKey"`
	jwt.RegisteredClaims
}

type AuthSession struct {
	ClientID   string
	ClientName string
	ClientKey  string
}

// IPLimiter menyimpan limiter untuk setiap IP
type IPLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.Mutex
	r   rate.Limit // Seberapa cepat token diisi ulang ke dalam ember (misal: 1 token per detik).
	b   int        // Kapasitas maksimal ember. Jika pengguna tidak melakukan request dalam waktu lama, mereka bisa
	//melakukan $b$ request sekaligus secara instan sebelum akhirnya dibatasi kembali oleh rate $r$.
}

var ctxName string = "clientCtx"

// constructor
func NewIPLimiter(r rate.Limit, b int) *IPLimiter {
	return &IPLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   r,
		b:   b,
	}
}

func (i *IPLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

func RateLimitMiddleware(limiter *IPLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)

		if !l.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many request, please try again later.",
			})
			return
		}

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminSecret := c.GetHeader("Admin-Secret")
		hashedAdminSecret := os.Getenv("ADMIN_SECRET")

		hashedSecret := utils.HashingSecret(adminSecret)

		if hashedSecret != hashedAdminSecret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		c.Next()
	}
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		tokenString := c.GetHeader("Authorization")
		var jwtToken string

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token required",
			})
			return
		} else if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			jwtToken = tokenString[7:] // Ambil setelah karakter ke-7
		}

		// Parse token
		token, err := jwt.ParseWithClaims(jwtToken, &MyClaims{}, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalid"})
			return
		}

		if claims, ok := token.Claims.(*MyClaims); ok {
			// Bungkus data ke dalam struct
			session := AuthSession{
				ClientID:   claims.ID,
				ClientName: claims.ClientName,
				ClientKey:  claims.ClientKey,
			}

			// Cek apakah token ada di Redis (Blacklist)
			redisKey := fmt.Sprintf("blacklist:%s", claims.ID)
			badToken, _ := config.Rdb.Get(config.Ctx, redisKey).Result()

			if badToken != "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token logged out"})
				return
			}

			// Simpan seluruh struct ke context
			c.Set(ctxName, session)
		}

		c.Next()
	}
}

func GetCtxFromReq(c *gin.Context) (AuthSession, error) {
	val, isExist := c.Get(ctxName)
	if !isExist {
		return AuthSession{}, errors.New("Current client not found")
	}

	// karena data yang didapat dari c.Get() berupa any, maka kita perlu melakukan
	// type asserton terlebih dahulu.
	client, ok := val.(AuthSession)
	if !ok {
		return AuthSession{}, errors.New("Failed to convert session to AuthSession")
	}
	return client, nil
}
