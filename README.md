# Final Project | Basket Service

## Overview

In this project, a basic basket service app is created. It provides a connection between a go gin server and postgreSQL. It includes user-product-cart-order relations.

## How to Use the App ?

The program includes some server endpoints which can be grouped into three categories by considering customers, admins and other users.

Non-registered users can only list the products, categories and search the products by name and sku. They can register into the system whereas customers and admins can login.

Customers can make some processes in the app which are:

* Adding some products to their carts,
* Listing the products in their carts,
* Updating/Deleting items in their carts,
* Completing Order with their items
* Listing orders up to date
* Canceling order until 14 days after their shopping

Admins can do things listed above and can additionally perform internal server processes like:

* Creating bulk category by uploading csv file
* Creating, Deleting and updating products

### Some Notes for Usage

1. Registered users can be online with their JWT tokens until it expires). Roles of users differentiate the accessibility into the system as it is defined above. So, authentication check is performed.

2. Category and product lists are paginated.

3. Product quantities are not changed in the database until order is submitted. So, users can add the products to cart even other users keep the same products in their carts.

### Endpoint Examples

1. User Processes:

```Postman
Login:
http://localhost:8080/api/v1/basket-service/user/login + Post Method
Body:
{
    "email":"xxxx",
    "password":"xxxx"
}

Sign Up:
http://localhost:8080/api/v1/basket-service/user/signup + Post Method
Body:
{
    "email":"xxxx",
    "password":"xxxx",
    "firstName":"xxxx",
    "lastName":"xxxx",
    "username":"xxxx",
    "passwordConfirm":"xxxx"
}
```

2. Category Processes

```Postman
http://localhost:8080/api/v1/basket-service/category/admin/addBulkCategory + Post Method
Form Data:
    Add csv file

http://localhost:8080/api/v1/basket-service/category/list?page=x&pageSize=x + Get Method
```

3. Cart Processes

```Postman
Add to Cart:
http://localhost:8080/api/v1/basket-service/cart/addToCart/productId/{id}/quantity/{quantity} + Post Method

List Cart Items:
http://localhost:8080/api/v1/basket-service/cart/listCartItems + Get Method

Delete Cart Items:
http://localhost:8080/api/v1/basket-service/cart/deleteItem/{id} + Delete Method

Update Cart Item Quantity:
http://localhost:8080/api/v1/basket-service/cart/updateItem/{id}/quantity/{quantity} + Put Method
```

4. Order Processes

```Postman
Complete Order:
http://localhost:8080/api/v1/basket-service/order/completeOrder + Post Method

List Orders:
http://localhost:8080/api/v1/basket-service/order/listOrder + Get Method

Cancel Order:
http://localhost:8080/api/v1/basket-service/order/cancelOrder/orderId/{orderId} + Put Method
```

5. Product Processes

```Postman
List Products:
http://localhost:8080/api/v1/basket-service/product/list?page=x&pageSize=x + Get Method

Search Products By Name:
http://localhost:8080/api/v1/basket-service/product/search/name/{name} + Get Method

Search Products By Sku:
http://localhost:8080/api/v1/basket-service/product/search/sku/{sku} + Get Method

Create Product:
http://localhost:8080/api/v1/basket-service/product/admin/createProduct + Post Method
Body:
{
    "name":"xxxx",
    "price":1111,
    "quantity":1111,
    "sku":"xxxx",
    "category": {
        "id":1111
    }
}

Update Product:
http://localhost:8080/api/v1/basket-service/product/admin/updateProduct/{productID} + Put Method
{
    "name":"xxxx"
}

Delete Product:
http://localhost:8080/api/v1/basket-service/product/admin/deleteProduct/{productID} + Delete Method

```

## Packages Used

Some significant third party packages used in project are:

* **GORM**
* **PostgreSQL**
* **Gin**
* **JWT**
* **Viper**
* **Swagger**
* **Zap**
