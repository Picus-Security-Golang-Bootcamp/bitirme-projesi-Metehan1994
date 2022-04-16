package cart

import (
	"errors"
	"fmt"

	"github.com/Metehan1994/final-project/internal/models"
	"github.com/Metehan1994/final-project/internal/product"
	"github.com/google/uuid"
)

type CartService struct {
	Cartrepo     *CartRepository
	productRepo  *product.ProductRepository
	cartItemRepo *CartItemRepository
}

func InitializeCartService(Cartrepo *CartRepository, productRepo *product.ProductRepository, cartItemRepo *CartItemRepository) *CartService {
	return &CartService{
		Cartrepo:     Cartrepo,
		productRepo:  productRepo,
		cartItemRepo: cartItemRepo,
	}
}

func (cserv *CartService) GetOrCreateCart(userID uuid.UUID) (*models.Cart, string) {
	cart := cserv.Cartrepo.GetCartByUserID(userID)
	var s string
	if cart.ID != uuid.Nil {
		s = "You have already a cart. New item will be added to it."
		fmt.Println(s)
	} else {
		cart.UserID = userID
		cart = cserv.Cartrepo.CreateCart(cart)
	}
	return cart, s
}

func (cserv *CartService) GetProductByID(productId int) *models.Product {
	product, _ := cserv.productRepo.GetByID(productId)
	return product
}

func (cserv *CartService) AddItem(cart *models.Cart, product *models.Product, quantity int) (*models.Cart, error) {
	var foundCart bool = false
	for _, item := range cart.Items {
		if item.ProductID == product.ID {
			foundCart = true
		}
	}
	if foundCart {
		return cart, errors.New("the product is already available in the basket. Please update its quantity")
	} else {
		if product.Quantity < quantity {
			return nil, errors.New("product quantity is not enough to compansate your demand")
		}
		cartItem := models.CartItem{
			ProductID: product.ID,
			Product:   *product,
			Price:     product.Price * quantity,
			Amount:    quantity,
			CartID:    cart.ID,
		}
		cserv.cartItemRepo.CreateCartItem(&cartItem)
		cart.Items = append(cart.Items, cartItem)
		cart.TotalPrice += quantity * product.Price
		return cart, nil
	}
}

func (cserv *CartService) UpdateCartInDB(cart *models.Cart) {
	cserv.Cartrepo.Update(cart)
}

func (cserv *CartService) GetCartByUserID(userID uuid.UUID) *models.Cart {
	cart := cserv.Cartrepo.GetCartByUserID(userID)
	return cart
}

func (cserv *CartService) DeleteItem(userID uuid.UUID, id int) error {
	cart := cserv.Cartrepo.GetCartByUserID(userID)
	var cartItemFound bool = false
	for _, item := range cart.Items {
		if item.ID == uint(id) {
			cartItemFound = true
		}
	}
	if !cartItemFound {
		return errors.New("item not found")
	}
	cart, err := cserv.cartItemRepo.DeleteById(cart, uint(id))
	if err != nil {
		return err
	}
	cserv.Cartrepo.Update(cart)
	return nil
}

//UpdateQuantityById changes the quantity of item found in the cart
func (cserv *CartService) UpdateQuantityById(userID uuid.UUID, id int, quantity int) error {
	cart := cserv.Cartrepo.GetCartByUserID(userID)
	var cartItemFound bool = false
	for _, item := range cart.Items {
		if item.ID == uint(id) {
			cartItemFound = true
		}
	}
	if !cartItemFound {
		return errors.New("item not found")
	}
	cart, err := cserv.cartItemRepo.UpdateQuantityById(cart, id, quantity)
	if err != nil {
		return err
	}
	cserv.Cartrepo.db.Model(&cart).Preload("Items.Product")
	cserv.Cartrepo.Update(cart)
	return nil
}
