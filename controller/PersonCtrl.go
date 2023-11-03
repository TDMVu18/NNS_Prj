package controller

import (
	"GoAPI/initializer"
	"GoAPI/model"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// api tim 1 item bang id
//func GetPersonById(ctx *gin.Context) {
//	id := ctx.Param("id")
//	result := model.ModelGet(id)
//	ctx.JSON(http.StatusOK, gin.H{
//		"data": result,
//	})
//}

func ListPerson(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	search := ctx.DefaultQuery("search", "")
	results := model.ModelList(search)
	if err := ctx.ShouldBind(&results); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Không lấy được dữ liệu nhân viên từ CSDL",
			"error":   err.Error(),
		})
	}
	// phân trang
	total := len(results)
	rowsPerPage := 6
	startIndex := (page - 1) * rowsPerPage
	endIndex := startIndex + rowsPerPage
	if endIndex > total {
		endIndex = total
	}
	currentPageData := results[startIndex:endIndex]

	// tạo danh sách trang
	totalPages := int(math.Ceil(float64(len(results)) / float64(rowsPerPage)))
	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	var nextPage int
	var isLastPage bool

	if page < totalPages {
		nextPage = page + 1
	} else {
		isLastPage = true
	}

	// Render template
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"data":        currentPageData,
		"prevPage":    page - 1,
		"currentPage": page,
		"total":       total,
		"nextPage":    nextPage,
		"startIndex":  startIndex + 1,
		"endIndex":    endIndex,
		"isLastPage":  isLastPage,
		"pages":       pages,
		"search":      search,
	})
}

func AddPerson(ctx *gin.Context) {
	var person model.Person
	if err := ctx.ShouldBind(&person); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Xử lý request dữ liệu nhân viên từ form không thành công",
			"error":   err.Error(),
		})
		return
	}
	person.ID = primitive.NewObjectID()
	person.Appearance = false
	now := time.Now()
	person.CreatedAt = &now
	person.UpdatedAt = &now
	message := model.ModelCreate(person)
	fmt.Println(message)
	ctx.Redirect(http.StatusFound, "/person/info/web")
}

func DeletePersonById(ctx *gin.Context) {
	id := ctx.Query("id")
	page := ctx.PostForm("page")
	search := ctx.PostForm("search")
	message := model.ModelDelete(id)
	fmt.Println(message)
	redirectURL := fmt.Sprintf("/person/info/web?page=%s&search=%s", page, search)
	ctx.Redirect(http.StatusFound, redirectURL)
}

func UpdatePersonById(ctx *gin.Context) {
	id := ctx.Query("id")
	var person model.Person
	if err := ctx.ShouldBind(&person); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Xử lý request dữ liệu nhân viên từ form không thành công",
			"error":   err.Error(),
		})
		return
	}
	page := ctx.PostForm("page")
	search := ctx.PostForm("search")
	person.ID, _ = primitive.ObjectIDFromHex(id)
	now := time.Now()
	person.UpdatedAt = &now
	message := model.ModelUpdate(person)
	fmt.Println(message)
	fmt.Printf("page is %s and search is %s", page, search)
	redirectURL := fmt.Sprintf("/person/info/web?page=%s&search=%s", page, search)
	ctx.Redirect(http.StatusFound, redirectURL)
}

func ToggleAppearance(ctx *gin.Context) {
	id := ctx.Query("id")
	var person model.Person
	if err := ctx.ShouldBind(&person); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Xử lý request dữ liệu nhân viên từ form không thành công",
			"error":   err.Error(),
		})
		return
	}
	page := ctx.PostForm("page")
	search := ctx.PostForm("search")
	person.ID, _ = primitive.ObjectIDFromHex(id)
	//update true sang false
	person.Appearance = !person.Appearance
	person.Name = ctx.PostForm("name")
	person.Major = ctx.PostForm("major")
	person.Level = ctx.PostForm("level")
	now := time.Now()
	person.UpdatedAt = &now
	message := model.ModelUpdate(person)
	fmt.Println(message)
	fmt.Printf("page is %s and search is %s", page, search)
	redirectURL := fmt.Sprintf("/person/info/web?page=%s&search=%s", page, search)
	ctx.Redirect(http.StatusFound, redirectURL)
}

func ShowProfile(ctx *gin.Context) {
	// Lấy id thông qua url để truyền vào hàm model, trả dữ liệu ra
	id := ctx.Query("id")
	var person model.Person

	person = *model.ModelGet(id)
	personLevel := person.Level
	officeName := person.Office
	salaryCollection := initializer.ConnectDB("salary_info")
	var salary model.Salary
	err := salaryCollection.FindOne(context.Background(), bson.M{"level": personLevel}).Decode(&salary)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Không tìm thấy thông tin lương của nhân viên",
			"error":   err.Error(),
		})
	}
	initializer.DisconnectDB()

	officeCollection := initializer.ConnectDB("office_info")
	var office model.Office
	err = officeCollection.FindOne(context.Background(), bson.M{"name": officeName}).Decode(&office)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Không tìm thấy thông tin văn phòng của nhân viên",
			"error":   err.Error(),
		})
	}
	initializer.DisconnectDB()
	ctx.HTML(http.StatusOK, "profile.html", gin.H{
		"person":        person,
		"salaryValue":   salary.Value,
		"officeAddress": office.Address,
		"id":            id,
	})
}

func Upload(ctx *gin.Context) {
	id := ctx.PostForm("id")
	formFile, _ := ctx.FormFile("profileImage")
	uploadPath := "./uploads/"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Tạo đường dẫn cho file upload không thành công",
			"error":   err.Error(),
		})
		return
	}
	if string(filepath.Ext(formFile.Filename)) != ".jpg" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Định dạng tệp tin không được hỗ trợ",
		})
		return
	}
	fileName := id + filepath.Ext(formFile.Filename)
	filePath := filepath.Join(uploadPath, fileName)
	if err := ctx.SaveUploadedFile(formFile, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Lưu file không thành công",
			"error":   err.Error(),
		})
		return
	}
	redirectURL := fmt.Sprintf("/person/info/web/profile?id=%s", id)
	ctx.Redirect(http.StatusFound, redirectURL)
}

func GetSalaryLevels(ctx *gin.Context) {
	collection := initializer.ConnectDB("salary_info")
	defer initializer.DisconnectDB()

	var salaryLevels []string

	filter := bson.M{}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Không lấy được dữ liệu lương nhân viên",
			"error":   err.Error(),
		})
		return
	}

	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Không decode được dữ liệu lương nhân viên",
				"error":   err.Error(),
			})
			return
		}

		level, ok := result["level"].(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Dữ liệu lương nhân viên sai định dạng",
				"error":   err.Error(),
			})
			return
		}

		salaryLevels = append(salaryLevels, level)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"salaryLevels": salaryLevels,
	})
}

func SalaryAdd(ctx *gin.Context) {
	var salary model.Salary
	if err := ctx.ShouldBind(&salary); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Xử lý request dữ liệu lương từ form không thành công",
			"error":   err.Error(),
		})
		return
	}
	salary.ID = primitive.NewObjectID()
	salary.Value = salary.Value + " $"
	message := model.ModelSalaryCreate(salary)
	fmt.Println(message)
	ctx.Redirect(http.StatusFound, "/person/salary")
}

func ListSalary(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	results := model.ModelSalaryList() // Sử dụng hàm ModelListSalary để lấy dữ liệu mức lương
	total := len(results)
	rowsPerPage := 6
	startIndex := (page - 1) * rowsPerPage
	endIndex := startIndex + rowsPerPage
	if endIndex > total {
		endIndex = total
	}
	currentPageData := results[startIndex:endIndex]
	totalPages := int(math.Ceil(float64(len(results)) / float64(rowsPerPage)))
	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	var nextPage int
	var isLastPage bool

	if page < totalPages {
		nextPage = page + 1
	} else {
		isLastPage = true
	}

	// Render Template
	ctx.HTML(http.StatusOK, "salary.html", gin.H{
		"data":        currentPageData,
		"prevPage":    page - 1,
		"currentPage": page,
		"total":       total,
		"nextPage":    nextPage,
		"startIndex":  startIndex + 1,
		"endIndex":    endIndex,
		"isLastPage":  isLastPage,
		"pages":       pages,
	})
}

func DeleteSalaryById(ctx *gin.Context) {
	id := ctx.Query("id")
	page := ctx.PostForm("page")
	message := model.ModelDeleteSalary(id)
	fmt.Println(message)
	redirectURL := fmt.Sprintf("/person/salary?page=%s", page)
	ctx.Redirect(http.StatusFound, redirectURL)
}

func UpdateSalaryById(ctx *gin.Context) {
	id := ctx.Query("id")
	var salary model.Salary
	if err := ctx.ShouldBind(&salary); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Xử lý request dữ liệu lương từ form không thành công",
			"error":   err.Error(),
		})
		return
	}
	page := ctx.PostForm("page")
	salary.ID, _ = primitive.ObjectIDFromHex(id)
	message := model.ModelUpdateSalary(salary)
	fmt.Println(message)
	redirectURL := fmt.Sprintf("/person/salary?page=%s", page)
	ctx.Redirect(http.StatusFound, redirectURL)
}

func OfficeAdd(ctx *gin.Context) {
	var office model.Office
	if err := ctx.ShouldBind(&office); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Xử lý request dữ liệu văn phòng từ form không thành công",
			"error":   err.Error(),
		})
		return
	}
	office.ID = primitive.NewObjectID()
	message := model.ModelOfficeCreate(office)
	fmt.Println(message)
	ctx.Redirect(http.StatusFound, "/person/office")
}

func ListOffice(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	results := model.ModelOfficeList()
	total := len(results)

	rowsPerPage := 6
	startIndex := (page - 1) * rowsPerPage
	endIndex := startIndex + rowsPerPage
	if endIndex > total {
		endIndex = total
	}
	currentPageData := results[startIndex:endIndex]
	totalPages := int(math.Ceil(float64(len(results)) / float64(rowsPerPage)))
	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	var nextPage int
	var isLastPage bool

	if page < totalPages {
		nextPage = page + 1
	} else {
		isLastPage = true
	}

	// Render Template
	ctx.HTML(http.StatusOK, "office.html", gin.H{
		"data":        currentPageData,
		"prevPage":    page - 1,
		"currentPage": page,
		"total":       total,
		"nextPage":    nextPage,
		"startIndex":  startIndex + 1,
		"endIndex":    endIndex,
		"isLastPage":  isLastPage,
		"pages":       pages,
	})
}

func DeleteOfficeById(ctx *gin.Context) {
	id := ctx.Query("id")
	page := ctx.PostForm("page")
	message := model.ModelDeleteOffice(id)
	fmt.Println(message)
	redirectURL := fmt.Sprintf("/person/office?page=%s", page)
	ctx.Redirect(http.StatusFound, redirectURL)
}

func UpdateOfficeById(ctx *gin.Context) {
	id := ctx.Query("id")
	var office model.Office
	if err := ctx.ShouldBind(&office); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Xử lý request dữ liệu văn phòng từ form không thành công",
			"error":   err.Error(),
		})
		return
	}
	page := ctx.PostForm("page")
	office.ID, _ = primitive.ObjectIDFromHex(id)
	message := model.ModelUpdateOffice(office)
	fmt.Println(message)
	redirectURL := fmt.Sprintf("/person/office?page=%s", page)
	ctx.Redirect(http.StatusFound, redirectURL)
}

func GetOfficeName(ctx *gin.Context) {
	collection := initializer.ConnectDB("office_info")
	defer initializer.DisconnectDB()

	var officeNames []string

	filter := bson.M{}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Không lấy được dữ liệu văn phòng",
			"error":   err.Error(),
		})
		return
	}

	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Không decode được dữ liệu văn phòng",
				"error":   err.Error(),
			})
			return
		}

		name, ok := result["name"].(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Dữ liệu văn phòng sai định dạng",
				"error":   err.Error(),
			})
			return
		}

		officeNames = append(officeNames, name)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"officeNames": officeNames,
	})
}
