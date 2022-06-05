package model

import (
	pb "pos-microservices/cashier/contract"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cashier struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Created  int64              `json:"created" bson:"created"`
	Updated  int64              `json:"updated" bson:"updated"`
}

type Cashiers []*Cashier

func (c *Cashier) ToPB() *pb.Cashier {
	return &pb.Cashier{
		Id:       c.ID.Hex(),
		Name:     c.Name,
		Email:    c.Email,
		Password: c.Password,
		Created:  c.Created,
		Updated:  c.Updated,
	}
}

func (c *Cashiers) ToPB() []*pb.Cashier {
	cashiers := make([]*pb.Cashier, len(*c))

	for _, cashier := range *c {
		cashiers = append(cashiers, cashier.ToPB())
	}

	return cashiers
}

func (c *Cashier) FromPB(cashier *pb.Cashier) {
	c.Name = cashier.Name
	c.Email = cashier.Email
	c.Password = cashier.Password
}

func (c *Cashiers) FromPB(cashiers []*pb.Cashier) {
	*c = make(Cashiers, len(cashiers))

	for i, cashier := range cashiers {
		(*c)[i] = &Cashier{}
		(*c)[i].FromPB(cashier)
	}
}
