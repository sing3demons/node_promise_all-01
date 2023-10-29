package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RequestData struct {
	Id         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Href       string    `json:"href,omitempty"`
	CreateDate time.Time `json:"createDate,omitempty"`
	UpdateDate time.Time `json:"updateDate,omitempty"`
}

type deleteData struct {
	Id         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Href       string    `json:"href,omitempty"`
	CreateDate time.Time `json:"createDate,omitempty"`
	UpdateDate time.Time `json:"updateDate,omitempty"`
	DeleteDate time.Time `json:"deleteDate,omitempty"`
}

type DB struct {
	col *mongo.Collection
}

func (h *DB) getTotal(filter primitive.M, ch chan int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	total, err := h.col.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	ch <- total
	return nil
}

func (h *DB) getCount(filter primitive.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	total, err := h.col.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return total, nil
}

type ResponseData struct {
	Total int64         `json:"total"`
	Data  []RequestData `json:"data"`
}

func TimeToBangkok(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return t.In(loc)
}

func (h *DB) getData(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response := []RequestData{}
	opts := options.FindOptions{}
	opts.SetSort(bson.D{{Key: "name", Value: 1}})

	filter := bson.M{"deleteDate": nil}
	cursor, err := h.col.Find(ctx, filter, &opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	for cursor.Next(ctx) {
		result := RequestData{}
		if err := cursor.Decode(&result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		response = append(response, RequestData{
			Id:         result.Id,
			Name:       result.Name,
			Href:       fmt.Sprintf("http://localhost:2566/example/%s", result.Id),
			CreateDate: TimeToBangkok(result.CreateDate),
			UpdateDate: TimeToBangkok(result.UpdateDate),
		})
	}

	count, err := h.getCount(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	r := ResponseData{Total: count, Data: response}

	c.JSON(http.StatusOK, r)
	end := time.Now()
	fmt.Printf("Time: %s\n", end.Sub(start))
}

func (h *DB) getDataMulti(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response := []RequestData{}
	opts := options.FindOptions{}
	opts.SetSort(bson.D{{Key: "name", Value: 1}})

	filter := bson.M{"deleteDate": nil}
	total := make(chan int64, 1)
	go h.getTotal(filter, total)
	cursor, err := h.col.Find(ctx, filter, &opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	for cursor.Next(ctx) {
		result := RequestData{}
		if err := cursor.Decode(&result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		response = append(response, RequestData{
			Id:         result.Id,
			Name:       result.Name,
			Href:       fmt.Sprintf("http://localhost:2566/example/%s", result.Id),
			CreateDate: TimeToBangkok(result.CreateDate),
			UpdateDate: TimeToBangkok(result.UpdateDate),
		})
	}

	r := ResponseData{Total: <-total, Data: response}
	c.JSON(http.StatusOK, r)
	end := time.Now()
	fmt.Printf("Time: %s\n", end.Sub(start))
}

func main() {
	col, err := ConnectMonoDB()
	if err != nil {
		panic(err)
	}

	data := []RequestData{}
	for i := 0; i < 900; i++ {
		data = append(data, RequestData{Name: fmt.Sprintf("Name_%d", i)})
	}
	r := gin.Default()

	h := &DB{col: col}
	r.GET("/example", h.getData)
	r.GET("/example/multi", h.getDataMulti)
	r.GET("/example/:id", h.getDataById)
	r.DELETE("/example/:id", h.deleteDataById)
	r.POST("/example/1", func(c *gin.Context) {
		start := time.Now()
		result := createAsyncData(data) // Time: 5.04395879s
		for _, data := range result {
			fmt.Println(data)
		}
		c.JSON(http.StatusOK, result)
		end := time.Now()
		fmt.Printf("Time: %s\n", end.Sub(start))
	})
	r.POST("/example/2", func(c *gin.Context) {
		start := time.Now()
		result := createData(data) // Time: 10.354932109s
		for _, data := range result {
			fmt.Println(data)
		}
		c.JSON(http.StatusOK, result)
		end := time.Now()
		fmt.Printf("Time: %s\n", end.Sub(start))
	})

	r.Run(":2566")

}


func (h *DB) getDataById(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")

	result := RequestData{}

	filter := bson.M{"id": id, "deleteDate": nil}

	if err := h.col.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := RequestData{
		Id:   result.Id,
		Name: result.Name,
		Href: fmt.Sprintf("http://localhost:3000/example/%s", result.Id),
		CreateDate: TimeToBangkok(result.CreateDate),
		UpdateDate: TimeToBangkok(result.UpdateDate),
	}

	c.JSON(http.StatusOK, response)
	end := time.Now()
	fmt.Printf("Time: %s\n", end.Sub(start))
}

func (h *DB) deleteDataById(c *gin.Context) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")

	result := deleteData{}

	filter := bson.M{"id": id, "deleteDate": nil}

	updateDoc := bson.M{"$set": bson.M{"deleteDate": TimeToBangkok(time.Now())}}

	opts := options.FindOneAndUpdate()
	opts.SetReturnDocument(options.After)
	opts.SetReturnDocument(options.After)
	if err := h.col.FindOneAndUpdate(ctx, filter, updateDoc, opts).Decode(&result); err != nil {
			if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := deleteData{
		Id:         result.Id,
		Name:       result.Name,
		Href:       fmt.Sprintf("http://localhost:3000/example/%s", result.Id),
		CreateDate: TimeToBangkok(result.CreateDate),
		UpdateDate: TimeToBangkok(result.UpdateDate),
		DeleteDate: TimeToBangkok(result.DeleteDate),
	}

	c.JSON(http.StatusOK, response)
	end := time.Now()
	fmt.Printf("Time: %s\n", end.Sub(start))
}
