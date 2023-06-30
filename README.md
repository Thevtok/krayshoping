# krayshoping
E Commerce Restfull API using Go,GIN,PostgreSQL


# Description
Krayshoping is an e-commerce application built using the go language.
Everyone can register as a user or a merchant. Users and merchants have different authentications. To top up users, this project uses Midtrans as a partner.Every transaction to send money, shop, or top up is guaranteed safe because it uses jwt tokens and some validation to overcome Cross-Site Scripting (XSS) and other attacks. If the transaction fails, your money will be returned.And every API hit is recorded in a log file to track the result of every failed or successful API hit



## Features

- Register for user or merchant
- Top up using Midtrans for user
- Send money to other user
- Buy merchant product
- Add Product for merchant
- Update Product for merchant
- Delete Product for merchant
- Sell Product for merchant



## Tech



- [GIN](https://github.com/gin-gonic/gin)
- [POSTGRESQL](https://github.com/lib/pq)
- [ENV](github.com/joho/godotenv)
- [Bcrypt](golang.org/x/crypto/bcrypt)
- [Logging](github.com/sirupsen/logrus)
- [JWT](github.com/dgrijalva/jwt-go)
- [Midtrans](github.com/midtrans/midtrans-go)

## API endpoint
- ##### User Register
```sh
METHOD POST
/user/register
```
- ##### User Login
```sh
METHOD POST
/user/login
```
- ##### User Logout
```sh
METHOD POST
/user/logout
```
- ##### Find UserByPhone
```sh
METHOD GET
/user/:phone_number
```
#### Required User Login
- ##### List Linked Bank Accout
```sh
METHOD GET
/bank/:user_id
```
- ##### Choose One Bank Accout
```sh
METHOD GET
/bank/:user_id/:account_number
```
- ##### Link Bank Account
```sh
METHOD POST
/bank/add/:user_id
```
- ##### Delete the Linked Bank Account
```sh
METHOD DELETE
/bank/delete/:user_id/:account_number
```
- ##### Top up User
```sh
METHOD POST
/transaction/deposit/:user_id/:account_number
```
- ##### Transfer User
```sh
METHOD POST
/transaction/transfer/:user_id
```
- ##### Buy Product
```sh
METHOD POST
/transaction/payment/:user_id/:merchant_id/:product_id
```
- ##### Update Payment Status
```sh
METHOD PUT
/transaction/payment/:user_id/:merchant_id/:tx_id
```
- ##### History All Transaction
```sh
METHOD GET
/transaction/:user_id
```

- ##### Merchant Register
```sh
METHOD POST
/merchant/register
```
- ##### Merchant Login
```sh
METHOD POST
/merchant/login
```
- ##### Merchant Logout
```sh
METHOD POST
/merchant/logout
```
- ##### Get All Product Every Mrchant
```sh
METHOD GET
/product
```
#### Required Merchant Login
- ##### Add Product
```sh
METHOD POST
/product/add/:merchant_id
```
- ##### Update Product
```sh
METHOD PUT
/product/update/:merchant_id/:product_id
```
- ##### Delete Product
```sh
METHOD DELETE
/product/delete/:merchant_id/:product_id
```
#### Created By Muhammad Fikri Alfarizi
