1.SET ENV ตาม config
MONGO_URI=mongodb://localhost:27017
DB_NAME=user-api

2.รันคำสั่ง "make dev" ใน terminal เพื่อทำการ build && up docker 

3.set url ในการยิงผ่าน postman เช่น  localhost:8080/api/v1/users/


ตัวโปรเจคนั้น ออกแบบเป็น rest api ในรูปแบบ monolith

backend-exercise/
├── cmd/                     # Entry point หรือไฟล์ main.go สำหรับรันแอป
├── config/                  # การจัดการ configuration เช่น env config
├── infrastructures/        # การเชื่อมต่อกับ database (mongo.go) หรือ kafka , redis
├── internal/               # แยก logic ภายในแบบ clean architecture
│   ├── core/               # business logic เช่น service layer
│   ├── handlers/           # HTTP handlers (controller layer)
│   └── repositories/       # data access layer (เช่น MongoRepository)
├── middleware/             # Middleware เช่น auth, logging ฯลฯ
├── utils/                  # ฟังก์ชันช่วยเหลือ เช่น hashing password
├── .env                    # Environment variables
├── config.yaml             # ไฟล์ config หลักของระบบ
├── docker-compose.yaml     # Setup แบบ container ด้วย MongoDB
├── dockerfile              # สำหรับ build container image
├── go.mod / go.sum         # Go module dependencies
├── makefile                # คำสั่งช่วยในการ build/run/test
└── README.md               # คำอธิบายโปรเจกต์

***อธิบายแนวทางในการ ขยายระบบให้รองรับโหลด 10 เท่า
ในกรณีที่มีผู้ใช้จำนวนมากขึ้น จะทำการเปลี่ยนจาก rest api รูปแบบ monolith เป็น รูปแบบ Microservices โดยใช่้ GRPC เพื่อคุยกับ service ต่างๆและใช้kafka ในการส่ง events ไปยัง service ต่างๆ