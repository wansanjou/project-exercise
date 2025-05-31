1.SET ENV ตาม config
MONGO_URI=mongodb://localhost:27017
DB_NAME=user-api

2.รันคำสั่ง "make dev" ใน terminal เพื่อทำการ build && up docker 

3.set url ในการยิงผ่าน postman เช่น  localhost:8080/api/v1/users/

ตัวโปรเจคนั้น ออกแบบเป็น rest api ในรูปแบบ monolith โดยเขียนเป็นรูปแบบ hexagonal architecture โดยใช้ database เป็น mongo และ framework gin 


***อธิบายแนวทางในการ ขยายระบบให้รองรับโหลด 10 เท่า
ในกรณีที่มีผู้ใช้จำนวนมากขึ้น จะทำการเปลี่ยนจาก rest api รูปแบบ monolith เป็น รูปแบบ Microservices โดยใช่้ GRPC เพื่อคุยกับ service ต่างๆและใช้kafka ในการส่ง events ไปยัง service ต่างๆ
