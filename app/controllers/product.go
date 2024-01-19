package controllers

import (
	"fmt"
	"strings"
	"text/template"
	"webapp/app/config"
	"webapp/app/models"

	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Product row template
const productRowTemplate = `
<tr hx-target="this" hx-swap="outerHTML" class="text-left align-middle">
    <td class="p-4">
        <img src="{{.ImageUrl}}" alt="Product Image" class="w-[20rem] rounded-full">
    </td>
    <td class="p-4">{{.Sku}}</td>
    <td class="p-4">{{.Name}}</td>
    <td class="p-4">{{.Description}}</td>
    <td class="p-4">{{.UnitPrice}}$</td>
    <td class="p-4">{{if .Active}}Yes{{else}}No{{end}}</td>
    <td class="p-4">{{.UnitsInStock}}</td>
    <td class="p-4">{{.DateCreated}}</td>
    <td class="p-4">{{.LastUpdated}}</td>
    <td class="p-4">{{.CategoryID}}</td>
    <td class="p-4">{{.ProductCategory.CategoryName}}</td>
    <td class="p-4">
        <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                hx-get="/api/products/{{.ID}}/edit"
                hx-trigger="edit"
                onClick="let editing = document.querySelector('.editing')
							if(editing) {
							Swal.fire({title: 'Already Editing',
										showCancelButton: true,
										confirmButtonText: 'Yep, Edit This Row!',
										text:'Hey!  You are already editing a row!  Do you want to cancel that edit and continue?'})
							.then((result) => {
									if(result.isConfirmed) {
									htmx.trigger(editing, 'cancel')
									}
									htmx.trigger(this, 'edit')
								})
							} else {
								htmx.trigger(this, 'edit')
							}">
            Edit
        </button>
    </td>
</tr>
`

func GetAllProducts(c echo.Context) error {
	// Declare products
	var products []models.Product
	// Check if DB is initialized correctly
	if config.DB == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database connection is not initialized",
		})
	}
	// Select * from product
	result := config.DB.Preload("ProductCategory").Find(&products)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error fetching the products",
		})
	}
	// Render the table rows
	tmpl, err := template.New("productRow").Parse(productRowTemplate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error parsing the template"})
	}
	var htmlBuilder strings.Builder
	for _, product := range products {
		err := tmpl.Execute(&htmlBuilder, product)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error executing the template"})
		}
	}
	return c.HTML(http.StatusOK, htmlBuilder.String())
}

func GetProductByID(c echo.Context) error {
	// Declare product
	var product models.Product
	// Check DB init
	if config.DB == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database connection is not initialized",
		})
	}
	// Get the product ID from the path parameter
	productID := c.Param("id")
	// Select product_id from product
	result := config.DB.Preload("ProductCategory").First(&product, productID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error fetching the product",
		})
	}
	// Render the table rows
	tmpl, err := template.New("productRow").Parse(productRowTemplate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error parsing the template"})
	}
	var htmlBuilder strings.Builder
	err = tmpl.Execute(&htmlBuilder, product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error executing the template"})
	}
	return c.HTML(http.StatusOK, htmlBuilder.String())
}

func CreateProduct(c echo.Context) (err error) {
	p := new(models.Product)
	// Bind the input form into the model
	if err := c.Bind(p); err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	// Update the date that the product has been created
	now := time.Now()
	p.DateCreated = now
	p.LastUpdated = now
	// Create a new product in the database
	result := config.DB.Create(&p)
	if result.Error != nil {
		return c.String(http.StatusBadRequest, "Failed to create new product")
	}
	return c.JSON(http.StatusOK, "Product created successfully")
}

func UpdateProductByID(c echo.Context) (err error) {
	var product models.Product
	input := new(models.ProductInputDTO)
	productID := c.Param("id")
	// Bind the form input into the model input
	if err := c.Bind(input); err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	// Get the product by id from the database
	if result := config.DB.Preload("ProductCategory").First(&product, productID); result.Error != nil {
		return c.String(http.StatusBadRequest, "Product not found")
	}
	if input.Sku != nil {
		product.Sku = *input.Sku
	}
	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.UnitPrice != nil {
		product.UnitPrice = *input.UnitPrice
	}
	if input.ImageUrl != nil {
		product.ImageUrl = *input.ImageUrl
	}
	if input.UnitsInStock != nil {
		product.UnitsInStock = *input.UnitsInStock
	}
	if input.Active != nil {
		product.Active = *input.Active
	}
	if input.CategoryID != nil {
		product.CategoryID = *input.CategoryID
	}
	product.LastUpdated = time.Now()
	// Save the updated product
	if result := config.DB.Save(&product); result.Error != nil {
		return c.String(http.StatusInternalServerError, "Failed to update product")
	}
	// Render the table rows
	tmpl, err := template.New("productRow").Parse(productRowTemplate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error parsing the template"})
	}
	var htmlBuilder strings.Builder
	err = tmpl.Execute(&htmlBuilder, product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error executing the template"})
	}
	return c.HTML(http.StatusOK, htmlBuilder.String())
}

func DeleteProductByID(c echo.Context) (err error) {
	var product models.Product
	productID := c.Param("id")
	// Get the product by id from the database
	result := config.DB.Preload("ProductCategory").First(&product, productID)
	if result.Error != nil {
		return c.String(http.StatusBadRequest, "Product not found")
	}
	// Delete the product from the database
	done := config.DB.Delete(&product)
	if done.Error != nil {
		return c.String(http.StatusBadRequest, "Failed to delete product")
	}
	return c.JSON(http.StatusOK, "Product deleted successfully")
}

func GetProductToEdit(c echo.Context) (err error) {
	id := c.Param("id")
	fmt.Println(id)
	// Get the product by id from the database
	var product models.Product
	if err := config.DB.Preload("ProductCategory").First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}
	editableRowTemplate := `
	<tr class="editing bg-gray-100 align-middle" hx-trigger="cancel" hx-target="this" hx-swap="outerHTML" hx-get="/api/products/{{.ID}}">
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">
                <input type='text' class='form-input rounded-md w-full px-3 py-1.5' name='ImageUrl' value='{{.ImageUrl}}'>
            </td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">
                <input type='text' class='form-input rounded-md w-full px-3 py-1.5' name='Sku' value='{{.Sku}}'>
            </td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">
                <input type='text' class='form-input rounded-md w-full px-3 py-1.5' name='Name' value='{{.Name}}'>
            </td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">
                <textarea class='form-textarea rounded-md w-full h-32 px-3 py-1.5' name='Description'>{{.Description}}</textarea>
            </td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">
                <input type='number' class='form-input rounded-md w-full px-3 py-1.5' name='UnitPrice' value='{{.UnitPrice}}'>
            </td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">
                <select class='form-select rounded-md w-full px-7 py-1.5' name='active'">
					<option value='true' {{if .Active}}selected{{end}}>Yes</option>
					<option value='false' {{if not .Active}}selected{{end}}>No</option>
                </select>
            </td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">
                <input type='number' class='form-input rounded-md w-full px-3 py-1.5' name='UnitsInStock' value='{{.UnitsInStock}}'>
            </td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">---</td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">---</td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">
                <input type='number' class='form-input rounded-md w-full px-3 py-1.5' name='CategoryID' value='{{.CategoryID}}'>
            </td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">
                <input type='text' class='form-input rounded-md w-full px-3 py-1.5' name='CategoryName' value='{{.ProductCategory.CategoryName}}' readonly>
            </td>
            <td class="px-1 py-4 border-b border-gray-200 bg-white text-sm">
                <div class="flex space-x-2 justify-start">
                    <button class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
                            hx-get="/api/products/{{.ID}}">
                        Cancel
                    </button>
                    <button class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
                            hx-put="/api/products/{{.ID}}" hx-include="closest tr">
                        Save
                    </button>
                </div>
            </td>
        </tr>
	`
	// Render the editable table row
	tmpl, err := template.New("editableRow").Parse(editableRowTemplate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error parsing the template"})
	}
	var htmlBuilder strings.Builder
	if err := tmpl.Execute(&htmlBuilder, product); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error executing the template"})
	}
	return c.HTML(http.StatusOK, htmlBuilder.String())
}
