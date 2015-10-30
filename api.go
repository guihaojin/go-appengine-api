package count

import(
	"appengine"
	"appengine/datastore"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Count struct {
	Value int
}

func init() {
	router := gin.Default()

	router.GET("/", func (c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Gin API.")
	})

	v1 := router.Group("/v1")
	v1.POST("/count/set/:number", setCount)
	v1.GET("/count", getCount)
	v1.POST("/count/add/:number", addCount)
	v1.POST("/count/subtract/:number", subtractCount)
	v1.POST("/count/multiply/:number", multiplyCount)
	v1.DELETE("/count", deleteCount)

	// Handle all requests using net/http
	http.Handle("/", router)
}

func countKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Count", "default_count", 0, nil)
}

func setCount(c *gin.Context) {
	n, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		c.String(http.StatusForbidden, "invalid number.")
	}

	if count, err := getCountHelper(c); err == nil {
		count.Value = n
		_, err = saveCounterHelper(c, count)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.String(http.StatusOK, strconv.Itoa(count.Value))
	} else {
		val := Count{Value: n}
		_, err = saveCounterHelper(c, &val)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.String(http.StatusOK, strconv.Itoa(val.Value))
	}
}

func getCount(c *gin.Context) {
	count, err := getCountHelper(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.String(http.StatusOK, strconv.Itoa(count.Value))
}

func addCount(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		c.String(http.StatusForbidden, "invalid number.")
	}

	count, err := getCountHelper(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	count.Value += i
	_, err = saveCounterHelper(c, count)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.String(http.StatusOK, strconv.Itoa(count.Value))
}

func subtractCount(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		c.String(http.StatusForbidden, "invalid number.")
	}

	count, err := getCountHelper(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	count.Value -= i
	_, err = saveCounterHelper(c, count)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.String(http.StatusOK, strconv.Itoa(count.Value))
}

func multiplyCount(c *gin.Context) {
	i, err := strconv.Atoi(c.Param("number"))
	if err != nil {
		c.String(http.StatusForbidden, "invalid number.")
	}

	count, err := getCountHelper(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	count.Value *= i
	_, err = saveCounterHelper(c, count)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.String(http.StatusOK, strconv.Itoa(count.Value))
}

func deleteCount(c *gin.Context) {
	err := deleteCountHelper(c)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.String(http.StatusOK, "delete succeed.")
}

func getCountHelper(c *gin.Context) (*Count, error) {
	context := appengine.NewContext(c.Request)
	key := countKey(context)
	var count Count
	err := datastore.Get(context, key, &count)
	return &count, err
}

func saveCounterHelper(c *gin.Context, count *Count) (*datastore.Key, error) {
	context := appengine.NewContext(c.Request)
	key := countKey(context)
	return datastore.Put(context, key, count)
}

func deleteCountHelper(c *gin.Context) error {
	context := appengine.NewContext(c.Request)
	key := countKey(context)
	return datastore.Delete(context, key)
}
