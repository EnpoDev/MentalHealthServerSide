# Mental Health Companion API

Bu proje, Mental Health Companion uygulamasının backend API'sini içerir. Go programlama dili ve modern web teknolojileri kullanılarak geliştirilmiştir.

## Güvenlik Önlemleri

1. **Şifre Güvenliği**
   - Şifreler bcrypt ile hashlenerek saklanır
   - Minimum 8 karakter zorunluluğu
   - En az bir büyük harf, bir küçük harf, bir rakam ve bir özel karakter zorunluluğu
   - Şifre karşılaştırmaları timing attack'lere karşı güvenli

2. **Kimlik Doğrulama**
   - JWT (JSON Web Token) tabanlı kimlik doğrulama
   - Token'lar 24 saat geçerli
   - Her istek için token doğrulaması

3. **Veri Güvenliği**
   - Email adresleri benzersiz
   - Email formatı doğrulaması
   - Hassas veriler JSON response'larında gizlenir
   - SQL injection koruması (GORM ORM kullanımı)
   - Prepared statement'lar

4. **API Güvenliği**
   - CORS politikası yapılandırması
   - Rate limiting koruması (TODO)
   - OPTIONS request handling
   - Güvenli HTTP başlıkları

## Kurulum

1. PostgreSQL veritabanını kurun ve bir veritabanı oluşturun
2. `.env` dosyasını oluşturun:
   ```
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=123123
   DB_NAME=mental_health_companion
   DB_PORT=5432
   JWT_SECRET=your-super-secret-key-change-in-production
   PORT=8080
   ```
3. Bağımlılıkları yükleyin:
   ```bash
   go mod tidy
   ```
4. Hot Reload için Air'i yükleyin:
   ```bash
   go install github.com/cosmtrek/air@latest
   ```
5. Uygulamayı geliştirme modunda çalıştırın:
   ```bash
   air
   ```
   veya normal modda çalıştırın:
   ```bash
   go run main.go
   ```

## API Endpointleri

### Public Endpointler

#### POST /register
Yeni kullanıcı kaydı için kullanılır.

Request body:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "name": "John",
  "surname": "Doe"
}
```

#### POST /login
Kullanıcı girişi için kullanılır.

Request body:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

### Protected Endpointler

#### GET /me
Kullanıcı bilgilerini görüntülemek için kullanılır.

Header:
```
Authorization: Bearer <token>
```

## Performans Optimizasyonları

1. **Veritabanı**
   - Connection pooling
   - Prepared statement caching
   - İndekslenmiş sütunlar

2. **API**
   - GORM logger optimizasyonu
   - Gereksiz SQL sorgularının önlenmesi
   - Lazy loading

## Geliştirme Standartları

1. Clean Code prensipleri
2. SOLID prensipleri
3. Go best practices
4. Güvenlik odaklı geliştirme
5. Performans odaklı geliştirme

## Geliştirici Araçları

1. **Hot Reload**
   - Air kütüphanesi ile otomatik yeniden başlatma
   - Kod değişikliklerini anında görme
   - Hızlı geliştirme döngüsü

## TODO

- [ ] Rate limiting implementasyonu
- [ ] Detaylı loglama sistemi
- [ ] API dokümantasyonu (Swagger)
- [ ] Test coverage artırımı 