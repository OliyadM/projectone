package mongo

import (
    "context"
    "fmt"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type mongoOrderRepository struct {
    collection *mongo.Collection
}

func NewMongoOrderRepository(db *mongo.Database) order.Repository {
    return &mongoOrderRepository{
        collection: db.Collection("orders"),
    }
}

func (r *mongoOrderRepository) CreateOrder(ctx context.Context, o *order.Order) error {
    _, err := r.collection.InsertOne(ctx, o)
    return err
}

func (r *mongoOrderRepository) GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, error) {
    fmt.Printf("🔍 Querying orders for consumer: %s\n", consumerID)
    
    var orders []*order.Order
    filter := bson.M{"consumerid": consumerID}
    fmt.Printf("🔍 Using filter: %+v\n", filter)
    
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        fmt.Printf("❌ Error querying orders: %v\n", err)
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var o order.Order
        if err := cursor.Decode(&o); err != nil {
            fmt.Printf("❌ Error decoding order: %v\n", err)
            return nil, err
        }
        orders = append(orders, &o)
    }

    if err := cursor.Err(); err != nil {
        fmt.Printf("❌ Cursor error: %v\n", err)
        return nil, err
    }

    fmt.Printf("✅ Found %d orders for consumer %s\n", len(orders), consumerID)
    return orders, nil
}

func (r *mongoOrderRepository) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
    var o order.Order
    // First try to find by _id
    err := r.collection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&o)
    if err == nil {
        return &o, nil
    }
    if err != mongo.ErrNoDocuments {
        return nil, err
    }
    
    // If not found by _id, try to find by id field
    err = r.collection.FindOne(ctx, bson.M{"id": orderID}).Decode(&o)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }
    return &o, nil
}

func (r *mongoOrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status order.OrderStatus) error {
    filter := bson.M{"_id": orderID} // Fixed: Corrected filter key from "order_id" to "_id"
    update := bson.M{"$set": bson.M{"status": status}}

    _, err := r.collection.UpdateOne(ctx, filter, update)
    return err
}

func (r *mongoOrderRepository) DeleteOrder(ctx context.Context, orderID string) error {
    _, err := r.collection.DeleteOne(ctx, bson.M{"_id": orderID})
    return err
}

func (r *mongoOrderRepository) GetOrdersBySupplier(ctx context.Context, supplierID string) ([]*order.Order, error) {
    fmt.Printf("🔍 Querying orders for supplier: %s\n", supplierID)
    
    var orders []*order.Order
    filter := bson.M{"supplierid": supplierID}
    fmt.Printf("🔍 Using filter: %+v\n", filter)
    
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        fmt.Printf("❌ Error querying orders: %v\n", err)
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var o order.Order
        if err := cursor.Decode(&o); err != nil {
            fmt.Printf("❌ Error decoding order: %v\n", err)
            return nil, err
        }
        orders = append(orders, &o)
    }

    if err := cursor.Err(); err != nil {
        fmt.Printf("❌ Cursor error: %v\n", err)
        return nil, err
    }

    fmt.Printf("✅ Found %d orders for supplier %s\n", len(orders), supplierID)
    return orders, nil
}

func (r *mongoOrderRepository) GetOrdersByReseller(ctx context.Context, resellerID string) ([]*order.Order, error) {
    fmt.Printf("🔍 Querying orders for reseller: %s\n", resellerID)
    
    var orders []*order.Order
    filter := bson.M{"resellerid": resellerID}
    fmt.Printf("🔍 Using filter: %+v\n", filter)
    
    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        fmt.Printf("❌ Error querying orders: %v\n", err)
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var o order.Order
        if err := cursor.Decode(&o); err != nil {
            fmt.Printf("❌ Error decoding order: %v\n", err)
            return nil, err
        }
        orders = append(orders, &o)
    }

    if err := cursor.Err(); err != nil {
        fmt.Printf("❌ Cursor error: %v\n", err)
        return nil, err
    }

    fmt.Printf("✅ Found %d orders for reseller %s\n", len(orders), resellerID)
    return orders, nil
}