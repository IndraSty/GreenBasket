package repository

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type productRepository struct {
	Collection *mongo.Collection
}

func NewProductRepository(client *mongo.Client) domain.ProductRepository {
	return &productRepository{
		Collection: db.OpenCollection(client, "Products"),
	}
}

// CheckNameExists implements domain.ProductRepository.
func (repo *productRepository) CheckNameExists(ctx context.Context, name string) (bool, error) {
	count, err := repo.Collection.CountDocuments(ctx, bson.M{"product_name": name})
	return count > 0, err
}

// GetAllProductWithNoPage implements domain.ProductRepository.
func (repo *productRepository) GetAllProductWithNoPage(ctx context.Context, storeID string) (*[]domain.ProductWithSalesData, error) {
	filter := bson.M{"store_id": storeID}

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$lookup": bson.M{
				"from": "Sales_Report",
				"let":  bson.M{"product_id": "$product_id"},
				"pipeline": []bson.M{
					{
						"$unwind": "$products",
					},
					{
						"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$$product_id", "$products.product_id"}}},
					},
					{
						"$project": bson.M{
							"_id":            0,
							"average_rating": "$products.average_rating",
							"total_sales":    "$products.total_sales",
						},
					},
				},
				"as": "sales_data",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$sales_data",
				"preserveNullAndEmptyArrays": true,
			},
		},

		{
			"$addFields": bson.M{
				"average_rating": "$sales_data.average_rating",
				"total_sales":    "$sales_data.total_sales",
			},
		},
	}

	cur, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []domain.ProductWithSalesData
	for cur.Next(ctx) {
		var product domain.ProductWithSalesData
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &products, nil

}

func (repo *productRepository) CreateProduct(ctx context.Context, product domain.Products) (primitive.ObjectID, error) {
	insertResult, err := repo.Collection.InsertOne(ctx, product)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertResult.InsertedID.(primitive.ObjectID), nil
}

func (repo *productRepository) UpdateProduct(ctx context.Context, storeID, productID string, update bson.D) (*mongo.UpdateResult, error) {
	filter := bson.M{"store_id": storeID, "product_id": productID}
	return repo.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
}

// UpdateStockProduct implements domain.ProductRepository.
func (repo *productRepository) UpdateStockProduct(ctx context.Context, storeID string, productID string, stock int, updatedAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"store_id": storeID, "product_id": productID}
	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "stock", Value: stock}}},
		{Key: "$set", Value: bson.D{{Key: "updated_at", Value: updatedAt}}},
	}
	return repo.Collection.UpdateOne(ctx, filter, update)
}

func (repo *productRepository) DeleteProductById(ctx context.Context, storeID, productID string) (*mongo.DeleteResult, error) {
	filter := bson.M{"store_id": storeID, "product_id": productID}
	result, err := repo.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("no product was deleted")
	}

	return result, nil
}

func (repo *productRepository) GetProductById(ctx context.Context, productID string, storeID ...string) (*domain.ProductWithSalesData, error) {
	filter := bson.M{"product_id": productID}

	if len(storeID) > 0 {
		filter["store_id"] = storeID[0]
	}

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$lookup": bson.M{
				"from": "Sales_Report",
				"let":  bson.M{"product_id": "$product_id"},
				"pipeline": []bson.M{
					{
						"$unwind": "$products",
					},
					{
						"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$$product_id", "$products.product_id"}}},
					},
					{
						"$project": bson.M{
							"_id":            0,
							"average_rating": "$products.average_rating",
							"total_sales":    "$products.total_sales",
						},
					},
				},
				"as": "sales_data",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$sales_data",
				"preserveNullAndEmptyArrays": true,
			},
		},

		{
			"$addFields": bson.M{
				"average_rating": "$sales_data.average_rating",
				"total_sales":    "$sales_data.total_sales",
			},
		},
	}

	cur, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var product domain.ProductWithSalesData
	for cur.Next(ctx) {
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	if product.Product_id == "" {
		return nil, errors.New("product not found")
	}

	return &product, nil
}

// GetAllProduct implements domain.ProductRepository.
func (repo *productRepository) GetAllProduct(ctx context.Context, page int, storeID ...string) (*domain.PagedProducts, error) {
	filter := bson.M{}
	if len(storeID) > 0 {
		filter["store_id"] = storeID[0]
	}

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$lookup": bson.M{
				"from": "Sales_Report",
				"let":  bson.M{"product_id": "$product_id"},
				"pipeline": []bson.M{
					{
						"$unwind": "$products",
					},
					{
						"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$$product_id", "$products.product_id"}}},
					},
					{
						"$project": bson.M{
							"_id":            0,
							"average_rating": "$products.average_rating",
							"total_sales":    "$products.total_sales",
						},
					},
				},
				"as": "sales_data",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$sales_data",
				"preserveNullAndEmptyArrays": true,
			},
		},

		{
			"$addFields": bson.M{
				"average_rating": "$sales_data.average_rating",
				"total_sales":    "$sales_data.total_sales",
			},
		},
	}

	var limit int = 9
	skip := (page - 1) * limit

	totalCount, err := repo.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	lastPage := int(math.Ceil(float64(totalCount) / float64(limit)))
	pipeline = append(pipeline, bson.M{"$skip": skip}, bson.M{"$limit": limit})

	cur, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []domain.ProductWithSalesData
	for cur.Next(ctx) {
		var product domain.ProductWithSalesData
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &domain.PagedProducts{
		Products:  products,
		Page:      page,
		TotalItem: int(totalCount),
		LastPage:  lastPage,
	}, nil
}

// GetAllProductByQuery implements domain.ProductRepository.
func (repo *productRepository) GetAllProductByQuery(ctx context.Context, query string, page int, storeID ...string) (*domain.PagedProducts, error) {
	filter := bson.M{}
	if query != "" {
		filter["name"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: query,
				Options: "i",
			},
		}
	}
	if len(storeID) > 0 {
		filter["store_id"] = storeID[0]
	}

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$lookup": bson.M{
				"from": "Sales_Report",
				"let":  bson.M{"product_id": "$product_id"},
				"pipeline": []bson.M{
					{
						"$unwind": "$products",
					},
					{
						"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$$product_id", "$products.product_id"}}},
					},
					{
						"$project": bson.M{
							"_id":            0,
							"average_rating": "$products.average_rating",
							"total_sales":    "$products.total_sales",
						},
					},
				},
				"as": "sales_data",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$sales_data",
				"preserveNullAndEmptyArrays": true,
			},
		},

		{
			"$addFields": bson.M{
				"average_rating": "$sales_data.average_rating",
				"total_sales":    "$sales_data.total_sales",
			},
		},
	}

	var limit int = 9
	skip := (page - 1) * limit

	totalCount, err := repo.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	lastPage := int(math.Ceil(float64(totalCount) / float64(limit)))
	pipeline = append(pipeline, bson.M{"$skip": skip}, bson.M{"$limit": limit})

	cur, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []domain.ProductWithSalesData
	for cur.Next(ctx) {
		var product domain.ProductWithSalesData
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &domain.PagedProducts{
		Products:  products,
		Page:      page,
		TotalItem: int(totalCount),
		LastPage:  lastPage,
	}, nil
}

// GetAllProductByQueryForCust implements domain.ProductRepository.
func (repo *productRepository) GetAllProductByQueryForCust(ctx context.Context, page int, query ...string) (*domain.PagedProducts, error) {
	filter := bson.M{}

	if len(query) > 0 && query[0] != "" {
		filter["name"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: query[0],
				Options: "i",
			},
		}
	}

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from": "Sales_Report",
				"let":  bson.M{"product_id": "$product_id"},
				"pipeline": []bson.M{
					{
						"$unwind": "$products",
					},
					{
						"$match": filter,
					},
					{
						"$project": bson.M{
							"_id":            0,
							"average_rating": "$products.average_rating",
							"total_sales":    "$products.total_sales",
						},
					},
				},
				"as": "sales_data",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$sales_data",
				"preserveNullAndEmptyArrays": true,
			},
		},

		{
			"$addFields": bson.M{
				"average_rating": "$sales_data.average_rating",
				"total_sales":    "$sales_data.total_sales",
			},
		},
	}

	if len(query) > 1 && query[1] != "" {
		if query[1] == "asc" {
			pipeline = append(pipeline, bson.M{"$sort": bson.D{{Key: "price", Value: 1}}})
		} else if query[1] == "desc" {
			pipeline = append(pipeline, bson.M{"$sort": bson.D{{Key: "price", Value: -1}}})
		}
	}

	var limit int = 9
	skip := (page - 1) * limit

	totalCount, err := repo.Collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	lastPage := int(math.Ceil(float64(totalCount) / float64(limit)))
	pipeline = append(pipeline, bson.M{"$skip": skip}, bson.M{"$limit": limit})

	cur, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []domain.ProductWithSalesData
	for cur.Next(ctx) {
		var product domain.ProductWithSalesData
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &domain.PagedProducts{
		Products:  products,
		Page:      page,
		TotalItem: int(totalCount),
		LastPage:  lastPage,
	}, nil

}

// GetAllByCategory implements domain.ProductRepository.
func (repo *productRepository) GetAllByCategory(ctx context.Context, category string, page int, storeID ...string) (*domain.PagedProducts, error) {
	filter := bson.M{
		"category": category,
	}

	if len(storeID) > 0 {
		filter["store_id"] = storeID[0]
	}

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$lookup": bson.M{
				"from": "Sales_Report",
				"let":  bson.M{"product_id": "$product_id"},
				"pipeline": []bson.M{
					{
						"$unwind": "$products",
					},
					{
						"$match": bson.M{"$expr": bson.M{"$eq": []interface{}{"$$product_id", "$products.product_id"}}},
					},
					{
						"$project": bson.M{
							"_id":            0,
							"average_rating": "$products.average_rating",
							"total_sales":    "$products.total_sales",
						},
					},
				},
				"as": "sales_data",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$sales_data",
				"preserveNullAndEmptyArrays": true,
			},
		},

		{
			"$addFields": bson.M{
				"average_rating": "$sales_data.average_rating",
				"total_sales":    "$sales_data.total_sales",
			},
		},
	}

	var limit int = 9
	skip := (page - 1) * limit

	totalCount, err := repo.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	lastPage := int(math.Ceil(float64(totalCount) / float64(limit)))
	pipeline = append(pipeline, bson.M{"$skip": skip}, bson.M{"$limit": limit})

	cur, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []domain.ProductWithSalesData
	for cur.Next(ctx) {
		var product domain.ProductWithSalesData
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &domain.PagedProducts{
		Products:  products,
		Page:      page,
		TotalItem: int(totalCount),
		LastPage:  lastPage,
	}, nil
}

func (repo *productRepository) GetAllProductSorted(ctx context.Context, sortParams map[string]string, page int, storeID ...string) (*domain.PagedProducts, error) {
	filter := bson.M{}

	if len(storeID) > 0 {
		filter["store_id"] = storeID[0]
	}

	pipeline := []bson.M{
		{
			"$match": filter,
		},
		{
			"$lookup": bson.M{
				"from": "Sales_Report",
				"let":  bson.M{"product_id": "$product_id"},
				"pipeline": []bson.M{
					{
						"$unwind": "$products",
					},
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": []string{"$$product_id", "$products.product_id"},
							},
						},
					},
					{
						"$project": bson.M{
							"_id":            0,
							"average_rating": "$products.average_rating",
							"total_sales":    "$products.total_sales",
						},
					},
				},
				"as": "sales_data",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$sales_data",
				"preserveNullAndEmptyArrays": true,
			},
		},

		{
			"$addFields": bson.M{
				"average_rating": "$sales_data.average_rating",
				"total_sales":    "$sales_data.total_sales",
			},
		},
	}

	var limit int = 9
	skip := (page - 1) * limit

	totalCount, err := repo.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	lastPage := int(math.Ceil(float64(totalCount) / float64(limit)))
	pipeline = append(pipeline, bson.M{"$skip": skip}, bson.M{"$limit": limit})

	sortFields := bson.D{}
	for param, direction := range sortParams {
		sortOrder := 1
		if direction == "desc" {
			sortOrder = -1
		} else if direction == "asc" {
			sortOrder = 1
		}
		sortFields = append(sortFields, bson.E{Key: param, Value: sortOrder})
	}
	if len(sortFields) > 0 {
		pipeline = append(pipeline, bson.M{"$sort": sortFields})
	}

	cur, err := repo.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []domain.ProductWithSalesData
	for cur.Next(ctx) {
		var product domain.ProductWithSalesData
		err := cur.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	productsres := domain.PagedProducts{
		Products:  products,
		Page:      page,
		TotalItem: int(totalCount),
		LastPage:  lastPage,
	}

	return &productsres, nil
}
