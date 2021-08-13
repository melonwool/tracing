# tracing
gin.hander, http.hander, http request tool

## init global jaeger tracing 
```go
    import "github.com/melonwool/tracing/jaegerc"

    tracer, _ := jaegerc.NewJaegerTracer("127.0.0.1:6831", "user-center", jaegerc.JaegerType("const"), 1)
	
```

## jaeger tracing in gin
```go

import (
	"github.com/gin-gonic/gin"
	"github.com/melonwool/tracing/ginh"
)

func main(){
    engine := gin.New()
    engine.Use(ginh.OpenTracingHandler())
    engine.GET("/v1/user/get", func(ctx *gin.Context) {
    	// 通过 http.Request.Context 传递
        models.UserGet(ctx.Request.Context(), "userid")
    })
}
```

## jaeger tracing in gorm
```go
import (
    "github.com/melonwool/tracing/gormplugin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)
var DB = NewDB()

func UserGet(ctx context.Context, userid string) User{
	user := User{}
    DB.WithContext(ctx).Tabe("user").Where("userid = ?",userid).Find(&user)
	return user
}
// NewDB 使用 OpentracingPlugin
func NewDB()*gorm.DB{
    var err error
    var db *gorm.DB
    var dbLogger logger.Interface
    if c.Mysql.LogMode {
        dbLogger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
            SlowThreshold: 200 * time.Millisecond,
            LogLevel:      logger.Info,
            Colorful:      true,
        })
    } else {
        dbLogger = logger.Discard
    }
    db, err = gorm.Open(mysql.Open(c.Mysql.DSN), &gorm.Config{
        Logger: dbLogger,
    })
    if err != nil {
        sentry.CaptureException(err)
        panic(err)
    }
    sqlDB, err := db.DB()
    if err != nil {
        sentry.CaptureException(err)
        panic(err)
    }
    sqlDB.SetConnMaxLifetime(time.Second * time.Duration(c.Mysql.MaxLife))
    sqlDB.SetMaxIdleConns(c.Mysql.MaxIdle)
    sqlDB.SetMaxOpenConns(c.Mysql.MaxOpen)
    _ = db.Use(&gormplugin.OpentracingPlugin{})
    return db
}
```