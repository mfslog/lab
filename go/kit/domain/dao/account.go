package dao

type Account struct{
	ID int64 `bson:"_id"`
	Name string `bson:"name"`
	Email string `bson:"email"`
	Password string `bson:"password"`
}