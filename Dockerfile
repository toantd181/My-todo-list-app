# Sử dụng image nền tảng Fedora
FROM fedora:latest

# Cập nhật danh sách gói và cài đặt cowsay
RUN dnf update -y && dnf install -y cowsay

# Thiết lập lệnh mặc định để chạy cowsay
CMD ["cowsay", "ditmethangKhoi!"]
