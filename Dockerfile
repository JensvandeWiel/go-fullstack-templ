FROM node:20 AS build-frontend
WORKDIR /build

COPY ./frontend/package.json ./frontend/package-lock.json .
RUN npm install --frozen-lockfile

COPY ./frontend .
RUN npm run build

FROM golang:1.22 AS build

# Set the Current Working Directory inside the container
WORKDIR /build

# Copy code
COPY . .

# Download the deps
RUN go mod download

# Copy the frontend build into the expected folder
COPY --from=build-frontend /build/dist ./frontend/dist

RUN CGO_ENABLED=0 ENV=prod go build -buildvcs=false -o ./bin/app .

FROM alpine:3.14

COPY --from=build /build/bin/app /usr/bin/app

EXPOSE 8080

CMD ["/usr/bin/app"]