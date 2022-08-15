
# Go-JWT-Authentication

## What makes up a Jwt
- header
- payload
- signature

## JWT
- Jwt has claims, by using it you are completely able to access fields such as accessUUID, accessToken, refreshUUID, refreshToken.(It requires you provide access token to extract to get claims and from that claim you get all things related to access)

## Requires
- Redis
- Psql
## Usage
1. Access postgresql and create database with name <strong>jwt</strong>
2. Run migration and project
    ```
   make migrate-up
   go run main.go 
   ```

## Api
1. Login
   ![image](https://user-images.githubusercontent.com/76799846/184618116-8aab193a-9d55-4697-9f39-f24a2021a42f.png)
2. Test Hello
   ![image](https://user-images.githubusercontent.com/76799846/184618201-9ef643d2-8103-4ccb-8209-e4d7e1c6b10d.png)
3. Logout
   ![image](https://user-images.githubusercontent.com/76799846/184618329-1ced73c4-2523-4195-a2c6-4fbae030ecca.png)
4. Refresh Token
   ![image](https://user-images.githubusercontent.com/76799846/184618414-49c5675b-6dea-4729-9185-3a7e675d89b4.png)
