FROM golang:1.22.2-alpine

# Buat direktori app dan set sebagai WORKDIR
WORKDIR /app

# Copy file-file yang diperlukan
COPY go.mod go.sum ./

# Install dependencies
RUN go mod tidy

# Copy seluruh kode sumber
COPY . .

# Build aplikasi
RUN go build -o main ./main.go

# Berikan izin eksekusi pada binary
RUN chmod +x main

# Buat direktori untuk menyimpan token dan berikan izin yang sesuai
RUN mkdir -p /app/tokens && chmod 777 /app/tokens

# Expose port yang digunakan aplikasi
EXPOSE 4040

# Jalankan aplikasi
CMD ["./main"]